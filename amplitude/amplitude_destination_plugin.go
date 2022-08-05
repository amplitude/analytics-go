package amplitude

import (
	"sync"
	"time"
)

type messageType byte

const (
	flushMessage messageType = iota
	eventMessage
	flushQueueSizeFactor = 10
)

type message struct {
	messageType
	event *Event
	wg    *sync.WaitGroup
}

type AmplitudeDestinationPlugin struct {
	config         Config
	storage        *InMemoryStorage
	messageChannel chan message
	httpClient     httpClient
}

func (a *AmplitudeDestinationPlugin) Setup(config Config) {
	a.config = config
	a.storage = &InMemoryStorage{}
	a.messageChannel = make(chan message, a.config.FlushQueueSize*flushQueueSizeFactor)
	a.httpClient = httpClient{logger: config.Logger, serverURL: config.ServerURL}

	autoFlushTicker := time.NewTicker(a.config.FlushInterval)
	defer autoFlushTicker.Stop()

	go func() {
	Loop:
		for {
			select {
			case <-autoFlushTicker.C:
				a.sendEventsFromStorage(nil)
			case message, ok := <-a.messageChannel:
				a.config.Logger.Debug("Message received from messageChannel: ", message, message.event)
				if !ok {
					a.sendEventsFromStorage(nil)

					break Loop
				}
				if message.messageType == flushMessage {
					a.sendEventsFromStorage(message.wg)
				} else {
					a.storage.Push(message.event)

					if a.storage.Len() >= a.config.FlushQueueSize {
						a.sendEventsFromStorage(nil)
					}
				}
			}
		}
	}()
}

// Execute processes the event with plugins added to the destination plugin.
// Then pushed the event to storage waiting to be sent.
func (a *AmplitudeDestinationPlugin) Execute(event *Event) {
	if !isValidEvent(event) {
		a.config.Logger.Error("Invalid event, EventType and either UserID or DeviceID cannot be empty.", event)
	}

	if a.messageChannel == nil {
		return
	}

	a.messageChannel <- message{
		messageType: eventMessage,
		event:       event,
	}
}

func (a *AmplitudeDestinationPlugin) Flush() {
	var flushWaitGroup sync.WaitGroup

	flushWaitGroup.Add(1)

	a.messageChannel <- message{
		messageType: flushMessage,
		event:       nil,
		wg:          &flushWaitGroup,
	}

	flushWaitGroup.Wait()
}

func (a *AmplitudeDestinationPlugin) sendEventsFromStorage(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	events := a.storage.Pull()

	chunks := a.chunk(events)
	for _, chunk := range chunks {
		a.httpClient.send(payload{
			APIKey: a.config.APIKey,
			Events: chunk,
		})
	}
}

func (a *AmplitudeDestinationPlugin) Shutdown() {
	a.Flush()
	close(a.messageChannel)
}

func isValidEvent(event *Event) bool {
	return event.EventType != "" && (event.UserID != "" || event.DeviceID != "")
}

func (a *AmplitudeDestinationPlugin) chunk(events []*Event) [][]*Event {
	chunkNum := len(events)/a.config.FlushQueueSize + 1
	chunks := make([][]*Event, chunkNum)

	for index := range chunks[:chunkNum-1] {
		chunks[index] = make([]*Event, a.config.FlushQueueSize)
		copy(chunks[index], events[index*a.config.FlushQueueSize:(index+1)*a.config.FlushQueueSize])
	}

	chunks[chunkNum-1] = make([]*Event, len(events)-(chunkNum-1)*a.config.FlushQueueSize)
	copy(chunks[chunkNum-1], events[(chunkNum-1)*a.config.FlushQueueSize:])

	return chunks
}
