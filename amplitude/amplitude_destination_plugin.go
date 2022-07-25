package amplitude

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type payLoad struct {
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
		a.config.Logger.Error("Invalid event.")
	}

	a.config.StorageProvider.Push(event)
}

func (a *AmplitudeDestinationPlugin) Flush() {
	events := a.config.StorageProvider.Pull()
	payLoad := payLoad{api_key: a.config.APIKey}

	if len(events) != 0 {
		payLoad.events = append(payLoad.events, events...)
	}

	payLoadBytes, err := json.Marshal(payLoad)
	if err != nil {
		a.config.Logger.Error("Events encoding failed")
	}

	request, err := http.NewRequest("POST", a.config.ServerURL, bytes.NewBuffer(payLoadBytes))
	if err != nil {
		a.config.Logger.Error("New request failed")
	}

	request.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)
	if err != nil {
		a.config.Logger.Error("HTTP request failed")
	}

	defer response.Body.Close()
}

func (a *AmplitudeDestinationPlugin) Shutdown() {
	a.Flush()
}

func isValidEvent(event *Event) bool {
	return event.EventType != "" && event.userId != "" && event.DeviceId != ""
}
