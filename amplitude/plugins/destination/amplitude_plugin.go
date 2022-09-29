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
	config            types.Config
	storage           types.EventStorage
	client            *amplitudeHTTPClient
	responseProcessor *AmplitudeResponseProcessor
	messageChannel    chan amplitudeMessage
	messageChannelMu  sync.RWMutex
}

type amplitudeMessage struct {
	event *types.Event
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
		amplitudePayloadOptions{MinIDLength: config.MinIDLength},
		config.Logger,
		config.ConnectionTimeout,
	)
	p.responseProcessor = &AmplitudeResponseProcessor{
		EventStorage:           p.storage,
		MaxRetries:             config.FlushMaxRetries,
		RetryBaseInterval:      config.RetryBaseInterval,
		RetryThrottledInterval: config.RetryThrottledInterval,
		Now:                    time.Now,
		Logger:                 config.Logger,
	}

	messageChannel := p.messageChannel
	go func() {
		defer func() {
			if r := recover(); r != nil {
				p.config.Logger.Errorf("Panic in AmplitudePlugin: %s", r)
			}
		}()

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
					p.storage.PushNew(message.event)

					if p.storage.HasFullChunk() {
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
func (p *amplitudePlugin) Execute(event *types.Event) {
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

	for {
		events := p.storage.PullChunk()
		if len(events) == 0 {
			break
		}

		response := p.client.Send(amplitudePayload{
			APIKey: p.config.APIKey,
			Events: events,
		})
		result := p.responseProcessor.Process(events, response)

		if len(result.Events) > 0 && p.config.ExecuteCallback != nil {
			go func() {
				for _, event := range result.Events {
					p.executeCallback(types.ExecuteResult{
						PluginName: p.Name(),
						Event:      event,
						Code:       result.Code,
						Message:    result.Message,
					})
				}
			}()
		}
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

func (p *amplitudePlugin) executeCallback(result types.ExecuteResult) {
	defer func() {
		if r := recover(); r != nil {
			p.config.Logger.Errorf("Panic in callback: %s", r)
		}
	}()
	p.config.ExecuteCallback(result)
}

func isValidEvent(event *types.Event) bool {
	userID := event.EventOptions.UserID
	if userID == "" {
		userID = event.UserID
	}
	deviceID := event.EventOptions.DeviceID
	if deviceID == "" {
		deviceID = event.DeviceID
	}

	return event.EventType != "" && (userID != "" || deviceID != "")
}
