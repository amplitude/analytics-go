package destination

import (
	"sync"
	"time"

	"github.com/amplitude/analytics-go/amplitude/types"
)

func NewAmplitudePlugin() types.ExtendedDestinationPlugin {
	return &amplitudePlugin{}
}

type amplitudePlugin struct {
	config           types.Config
	storage          types.EventStorage
	client           *amplitudeHTTPClient
	messageChannel   chan amplitudeMessage
	messageChannelMu sync.RWMutex
}

type amplitudeMessage struct {
	event *types.EventPayload
	wg    *sync.WaitGroup
}

func (p *amplitudePlugin) Name() string {
	return "amplitude"
}

func (p *amplitudePlugin) Type() types.PluginType {
	return types.PluginTypeDestination
}

func (p *amplitudePlugin) Setup(config types.Config) {
	p.config = config
	p.storage = config.StorageFactory()
	p.messageChannel = make(chan amplitudeMessage, config.MaxStorageCapacity)
	p.client = newAmplitudeHTTPClient(
		config.ServerURL,
		clientPayloadOptions{MinIDLength: config.MinIDLength},
		config.Logger,
		config.ConnectionTimeout,
	)

	messageChannel := p.messageChannel
	go func() {
		autoFlushTicker := time.NewTicker(p.config.FlushInterval)
		defer autoFlushTicker.Stop()

		for {
			select {
			case <-autoFlushTicker.C:
				p.sendEventsFromStorage(nil)
			case message, ok := <-messageChannel:
				if !ok {
					p.sendEventsFromStorage(nil)

					return
				}

				if message.wg != nil {
					p.sendEventsFromStorage(message.wg)
					autoFlushTicker.Reset(p.config.FlushInterval)
				} else {
					p.storage.Push(message.event)

					if p.storage.Len() >= p.config.FlushQueueSize {
						p.sendEventsFromStorage(nil)
						autoFlushTicker.Reset(p.config.FlushInterval)
					}
				}
			}
		}
	}()
}

// Execute processes the event with plugins added to the destination plugin.
// Then pushed the event to storage waiting to be sent.
func (p *amplitudePlugin) Execute(event *types.EventPayload) {
	if !isValidEvent(event) {
		p.config.Logger.Errorf("Invalid event, EventType and either UserID or DeviceID cannot be empty: \n\t%+v", event)
	}

	p.messageChannelMu.RLock()
	defer p.messageChannelMu.RUnlock()

	select {
	case p.messageChannel <- amplitudeMessage{
		event: event,
		wg:    nil,
	}:
	default:
	}
}

func (p *amplitudePlugin) Flush() {
	p.messageChannelMu.RLock()
	defer p.messageChannelMu.RUnlock()

	p.flush(p.messageChannel)
}

func (p *amplitudePlugin) flush(messageChannel chan<- amplitudeMessage) {
	var flushWaitGroup sync.WaitGroup
	flushWaitGroup.Add(1)

	select {
	case messageChannel <- amplitudeMessage{
		event: nil,
		wg:    &flushWaitGroup,
	}:
	default:
		flushWaitGroup.Done()
	}

	flushWaitGroup.Wait()
}

func (p *amplitudePlugin) sendEventsFromStorage(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	events := p.storage.Pull()
	if len(events) == 0 {
		return
	}

	result := p.client.send(clientPayload{
		APIKey: p.config.APIKey,
		Events: events,
	})

	executeCallback := p.config.ExecuteCallback
	if executeCallback != nil {
		go func() {
			for _, event := range events {
				executeCallback(types.ExecuteResult{
					PluginName: p.Name(),
					Event:      event,
					Code:       result.Code,
					Message:    result.Message,
				})
			}
		}()
	}
}

func (p *amplitudePlugin) Shutdown() {
	p.messageChannelMu.Lock()
	messageChannel := p.messageChannel
	p.messageChannel = nil
	p.messageChannelMu.Unlock()

	p.flush(messageChannel)
	close(messageChannel)
}

func isValidEvent(event *types.EventPayload) bool {
	return event.EventType != "" && (event.UserID != "" || event.DeviceID != "")
}
