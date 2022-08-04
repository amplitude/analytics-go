// A basic example of using Amplitude Go SDK to set user property

package main

// Import amplitude package
import (
	"github.com/amplitude/Amplitude-Go/amplitude"
)

func main() {

	config := amplitude.NewConfig("your_api_key")

	client := amplitude.NewClient(config)

	// Identify struct provides controls over setting user properties.
	identifyObj := amplitude.Identify{}

	// Set the value of a user property
	identifyObj.Set("location", "LAX")

	// Call Identify method of client
	client.Identify(identifyObj, amplitude.EventOptions{UserID: "identify-user-id"})

	// Create a BaseEvent instance
	event := amplitude.Event{
		EventOptions: amplitude.EventOptions{DeviceID: "identify-device-id", UserID: "identify-user-id"},
		EventType:    "identify-event-type",
	}

	// Track an event
	client.Track(event)

	// Flush the event buffer
	client.Flush()

	// Shutdown the client
	client.Shutdown()
}
