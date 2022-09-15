// A basic example of using Amplitude Go SDK to track an event.

package main

// Import amplitude package
import (
	"github.com/amplitude/analytics-go/amplitude"
)

func main() {

	config := amplitude.NewConfig("your-api-key")

	client := amplitude.NewClient(config)

	// Track a basic event
	// One of UserID and DeviceID is required
	event := amplitude.Event{
		UserID:    "user-id",
		EventType: "Button Clicked",
	}
	client.Track(event)

	// Track events with optional properties
	client.Track(amplitude.Event{
		EventType:  "type-of-event",
		UserID:     "user-id",
		Properties: map[string]interface{}{"source": "notification"},
	}, amplitude.EventOptions{DeviceID: "device-id"})

	// Flush the event buffer
	client.Flush()

	// Shutdown the client
	client.Shutdown()
}
