package amplitude

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type payload struct {
	ApiKey string   `json:"api_key"`
	Events []*Event `json:"events"`
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

	a.config.Storage.Push(event)
}

func (a *AmplitudeDestinationPlugin) Flush() {
	events := a.config.Storage.Pull()
	eventPayload := &payload{
		ApiKey: a.config.APIKey,
		Events: events,
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
	return event.EventType != "" && event.UserId != "" && event.DeviceId != ""
}
