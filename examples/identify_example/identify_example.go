// A basic example of using Amplitude Go SDK to set user property

package main

// Import amplitude package
import (
	"github.com/amplitude/analytics-go/amplitude"
)

func main() {

	config := amplitude.NewConfig("your-api-key")

	client := amplitude.NewClient(config)

	// Identify struct provides controls over setting user properties.
	identifyObj := amplitude.Identify{}

	// Set the value of a user property
	identifyObj.Set("location", "LAX")

	// Call Identify method of client
	// Event here will not display in your Amplitude Analytics
	//client.Identify(identifyObj, amplitude.EventOptions{UserID: "identify-user-id"})
	client.Identify(identifyObj, amplitude.EventOptions{})

	// To see identify actually works
	// Let's track another event
	// Then you can see that user properties of this event has location set to LAX
	//event := amplitude.Event{
	//	EventOptions: amplitude.EventOptions{DeviceID: "identify-device-id", UserID: "identify-user-id"},
	//	EventType:    "identify-event-type",
	//}
	//client.Track(event)

	// Flush the event buffer
	client.Flush()

	// Shutdown the client
	client.Shutdown()
}
