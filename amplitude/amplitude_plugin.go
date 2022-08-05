package amplitude

import (
	"sync"
	"time"
)

const (
	flushQueueSizeFactor = 10
)

type message struct {
	event *Event
	wg    *sync.WaitGroup
}

type AmplitudePlugin struct {
	config         Config
	storage        *InMemoryStorage
	messageChannel chan message
	httpClient     httpClient
}

func (a *AmplitudePlugin) Setup(config Config) {
	a.config = config
	a.storage = &InMemoryStorage{}
	a.messageChannel = make(chan message, MaxBufferCapacity)
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
				if message.wg != nil {
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
func (a *AmplitudePlugin) Execute(event *Event) {
	if !isValidEvent(event) {
		a.config.Logger.Error("Invalid event, EventType and either UserID or DeviceID cannot be empty.", event)
	}

	if a.messageChannel == nil {
		return
	}

	a.messageChannel <- message{
		event: event,
		wg:    nil,
	}
}

func (a *AmplitudePlugin) Flush() {
	var flushWaitGroup sync.WaitGroup

	flushWaitGroup.Add(1)

	a.messageChannel <- message{
		event: nil,
		wg:    &flushWaitGroup,
	}

	flushWaitGroup.Wait()
}

func (a *AmplitudePlugin) sendEventsFromStorage(wg *sync.WaitGroup) {
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

func (a *AmplitudePlugin) Shutdown() {
	a.Flush()
	close(a.messageChannel)
}

func isValidEvent(event *Event) bool {
	return event.EventType != "" && (event.UserID != "" || event.DeviceID != "")
}

func (a *AmplitudePlugin) chunk(events []*Event) [][]*Event {
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