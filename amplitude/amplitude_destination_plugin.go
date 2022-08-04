package amplitude

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type payload struct {
	APIKey string   `json:"api_key"`
	Events []*Event `json:"events"`
}

type eventType byte

const (
	flushEvent eventType = iota
	userEvent
	flushQueueSizeFactor = 10
)

type destinationEvent struct {
	eventType
	event *Event
	wg    *sync.WaitGroup
}

type AmplitudeDestinationPlugin struct {
	config       Config
	storage      *InMemoryStorage
	eventChannel chan destinationEvent
}

func (a *AmplitudeDestinationPlugin) Setup(config Config) {
	a.config = config
	a.storage = &InMemoryStorage{}
	a.eventChannel = make(chan destinationEvent, a.config.FlushQueueSize*flushQueueSizeFactor)

	autoFlushTicker := time.NewTicker(a.config.FlushInterval)
	defer autoFlushTicker.Stop()

	go func() {
	Loop:
		for {
			select {
			case <-autoFlushTicker.C:
				a.flush(nil)
			case event, ok := <-a.eventChannel:
				a.config.Logger.Debug("Event received from eventChannel: ", event, event.event)
				if !ok {
					a.flush(nil)

					break Loop
				}
				if event.eventType == flushEvent {
					a.flush(event.wg)
				} else {
					if a.storage.Len() >= a.config.FlushQueueSize {
						a.flush(nil)
					}
					a.storage.Push(event.event)
				}
			}
		}
	}()
}

// Execute processes the event with plugins added to the destination plugin.
// Then pushed the event to storage waiting to be sent.
func (a *AmplitudeDestinationPlugin) Execute(event *Event) {
	if !isValidEvent(event) {
		a.config.Logger.Error("Invalid event, EventType, UserID, and DeviceID cannot be empty.", event)
	}

	if a.eventChannel == nil {
		return
	}

	a.eventChannel <- destinationEvent{
		eventType: userEvent,
		event:     event,
	}
}

func (a *AmplitudeDestinationPlugin) Flush() {
	var flushWaitGroup sync.WaitGroup
	flushWaitGroup.Add(1)

	a.eventChannel <- destinationEvent{
		eventType: flushEvent,
		event:     nil,
		wg:        &flushWaitGroup,
	}

	flushWaitGroup.Wait()
}

func (a *AmplitudeDestinationPlugin) flush(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	events := a.storage.Pull()

	chunks := a.chunk(events)
	for _, chunk := range chunks {
		a.sendChunk(chunk)
	}
}

func (a *AmplitudeDestinationPlugin) sendChunk(chunk []*Event) {
	eventPayload := &payload{
		APIKey: a.config.APIKey,
		Events: chunk,
	}

	eventPayloadBytes, err := json.Marshal(eventPayload)
	if err != nil {
		a.config.Logger.Error("Events encoding failed", err)
	}

	a.config.Logger.Debug("eventPayloadBytes: ", string(eventPayloadBytes))

	request, err := http.NewRequest("POST", a.config.ServerURL, bytes.NewReader(eventPayloadBytes))
	if err != nil {
		a.config.Logger.Error("Building new request failed", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)
	if err != nil {
		a.config.Logger.Error("HTTP request failed", err)
	}
	defer response.Body.Close()

	a.config.Logger.Info("HTTP request response", response)
}

func (a *AmplitudeDestinationPlugin) Shutdown() {
	a.Flush()
	close(a.eventChannel)
}

func isValidEvent(event *Event) bool {
	return event.EventType != "" && event.UserID != "" && event.DeviceID != ""
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
