package amplitude

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type payload struct {
	APIKey string   `json:"api_key"`
	Events []*Event `json:"events"`
}

type AmplitudeDestinationPlugin struct {
	config    Config
	scheduled bool
	storage   chan *Event
}

func (a *AmplitudeDestinationPlugin) Setup(config Config) {
	a.config = config
	a.scheduled = false
	a.storage = make(chan *Event, a.config.FlushQueueSize)
}

// Execute processes the event with plugins added to the destination plugin.
// Then pushed the event to storage waiting to be sent.
func (a *AmplitudeDestinationPlugin) Execute(event *Event) {
	if !isValidEvent(event) {
		a.config.Logger.Error("Invalid event, EventType, UserID, and DeviceID cannot be empty.", event)
	}

	if len(a.storage) == a.config.FlushQueueSize {
		a.Flush()
	}

	a.storage <- event

	if !a.scheduled {
		time.AfterFunc(a.config.FlushInterval, func() { go a.Flush() })
	}
}

func (a *AmplitudeDestinationPlugin) Flush() {
	if len(a.storage) == 0 {
		return
	}

	events := make([]*Event, len(a.storage))
	currentStorageSize := len(a.storage)

	for i := 0; i < currentStorageSize; i++ {
		events[i] = <-a.storage
	}

	go a.send(events)

	a.scheduled = false
}

func (a *AmplitudeDestinationPlugin) send(chunk []*Event) {
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

	a.config.Logger.Info("HTTP request response", response)

	defer response.Body.Close()
}

func (a *AmplitudeDestinationPlugin) Shutdown() {
	a.Flush()
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
