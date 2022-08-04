// A basic example of using Amplitude Go SDK to track an event.

package main

// Import amplitude package
import (
	"github.com/amplitude/Amplitude-Go/amplitude"
)

// Define your callback function (optional)
func callbackFunc(e string, code int, message string) {
	println(e)
	println(code, message)
}

func main() {

	config := amplitude.NewConfig("your-api-key")

	// Config callback function (optional)
	client := amplitude.NewClient(config)

	// Track a basic event
	// One of UserID and DeviceID is required
	event := amplitude.Event{
		EventOptions: amplitude.EventOptions{UserID: "user-id"},
		EventType:    "Button Clicked",
	}
	client.Track(event)

	// Track events with optional properties
	client.Track(amplitude.Event{
		EventType: "type-of-event",
		EventOptions: amplitude.EventOptions{
			UserID:   "user-id",
			DeviceID: "device-id",
		},
		EventProperties: map[string]interface{}{"source": "notification"},
	})

	// Flush the event buffer
	client.Flush()

	// Shutdown the client
	client.Shutdown()
}
