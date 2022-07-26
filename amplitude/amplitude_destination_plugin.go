package amplitude

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type payload struct {
	api_key string
	events  []*Event
}

type AmplitudeDestinationPlugin struct {
	config Config
}

func (a *AmplitudeDestinationPlugin) Setup(config Config) {
	a.config = config
}

// Execute processes the event with plugins added to the destination plugin.
// Then pushed the event to storage waiting to be sent.
func (a *AmplitudeDestinationPlugin) Execute(event *Event) {
	if !isValidEvent(event) {
		a.config.Logger.Error("Invalid event, EventType, UserID, and DeviceID cannot be empty.", event)
	}

	a.config.StorageProvider.Push(event)
}

func (a *AmplitudeDestinationPlugin) Flush() {
	events := a.config.StorageProvider.Pull()
	payload := payload{api_key: a.config.APIKey}

	if len(events) != 0 {
		payload.events = append(payload.events, events...)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		a.config.Logger.Error("Events encoding failed", err)
	}

	request, err := http.NewRequest("POST", a.config.ServerURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		a.config.Logger.Error("Building new request failed", err)
	}

	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)
	if err != nil {
		a.config.Logger.Error("HTTP request failed", err)
	}

	defer response.Body.Close()
}

func (a *AmplitudeDestinationPlugin) Shutdown() {
	a.Flush()
}

func isValidEvent(event *Event) bool {
	return event.EventType != "" && event.userId != "" && event.DeviceId != ""
}
