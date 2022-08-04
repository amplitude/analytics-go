package controllers

import (
	"fmt"
	"github.com/amplitude/Amplitude-Go/amplitude"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Analytics() revel.Result {
	config := amplitude.NewConfig("c253b75dce3e593c44ea5eb95999f92a")
	client := amplitude.NewClient(config)
	defer client.Shutdown()

	// Track a basic event
	// One of UserID and DeviceID is required
	for i := 0; i < 5; i++ {
		event := amplitude.Event{
			EventOptions: amplitude.EventOptions{UserID: "revel-user-id-" + fmt.Sprint(i), DeviceID: "revel-device-id"},
			EventType:    "Open Analytics",
		}
		client.Track(event)
	}

	return c.Render()
}
