package controllers

import (
	"github.com/amplitude/Amplitude-Go/amplitude"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

type payload struct {
	EventType      string                                          `json:"event_type"`
	UserProperties map[amplitude.IdentityOp]map[string]interface{} `json:"user_properties,omitempty"`
	DeviceID       string                                          `json:"device_id,omitempty"`
	UserID         string                                          `json:"user_id,omitempty"`
	SessionID      int                                             `json:"session_id,omitempty"`
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Analytics() revel.Result {
	config := amplitude.NewConfig("your-api-key")
	client := amplitude.NewClient(config)
	defer client.Shutdown()

	jsonData := payload{}
	//config.Logger.Debug("string(c.Params.JSON): ", string(c.Params.JSON))
	c.Params.BindJSON(&jsonData)
	//config.Logger.Debug("jsonData: ", jsonData)

	event := amplitude.Event{
		EventOptions: amplitude.EventOptions{
			UserID:    jsonData.UserID,
			DeviceID:  jsonData.DeviceID,
			SessionID: jsonData.SessionID,
		},
		EventType:      jsonData.EventType,
		UserProperties: jsonData.UserProperties,
	}
	client.Track(event)

	return c.Render()
}
