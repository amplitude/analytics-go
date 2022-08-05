package controllers

import (
	"github.com/amplitude/Amplitude-Go/amplitude"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

type payload struct {
	eventType      string                                          `json:"event_type"`
	userProperties map[amplitude.IdentityOp]map[string]interface{} `json:"user_properties,omitempty"`
	deviceID       string                                          `json:"device_id,omitempty"`
	userID         string                                          `json:"user_id,omitempty"`
	sessionID      int                                             `json:"session_id,omitempty"`
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Analytics() revel.Result {
	config := amplitude.NewConfig("c253b75dce3e593c44ea5eb95999f92a")
	client := amplitude.NewClient(config)
	defer client.Shutdown()

	var jsonData payload
	c.Params.BindJSON(&jsonData)

	event := amplitude.Event{
		EventOptions: amplitude.EventOptions{
			UserID:    jsonData.userID,
			DeviceID:  jsonData.deviceID,
			SessionID: jsonData.sessionID,
		},
		EventType:      jsonData.eventType,
		UserProperties: jsonData.userProperties,
	}
	client.Track(event)

	return c.Render()
}
