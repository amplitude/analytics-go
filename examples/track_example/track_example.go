// A basic example of using Amplitude Go SDK to track an event.

package main

// Import amplitude package
import (
	"fmt"
	"github.com/amplitude/Amplitude-Go/amplitude"
	"time"
)

// Define your callback function (optional)
func callbackFunc(e string, code int, message string) {
	println(e)
	println(code, message)
}

func main() {

	config := amplitude.NewConfig("your-api-key")
	config.FlushQueueSize = 3

	// Config callback function (optional)
	client := amplitude.NewClient(config)

	client.Add(amplitude.NewContextPlugin())

	// Create and track events
	for i := 0; i < 10; i++ {
		event := amplitude.Event{
			EventType: "go-event-type",
			EventOptions: amplitude.EventOptions{
				UserID:   "go-user-id-" + fmt.Sprint(i),
				DeviceID: "go-device-id-" + fmt.Sprint(i),
			},
		}
		client.Track(event)
	}

	// Flush the event buffer
	client.Flush()

	for i := 10; i < 20; i++ {
		event := amplitude.Event{
			EventType: "go-event-type",
			EventOptions: amplitude.EventOptions{
				UserID:   "go-user-id-" + fmt.Sprint(i),
				DeviceID: "go-device-id-" + fmt.Sprint(i),
			},
		}
		client.Track(event)
	}
	// Shutdown the client
	client.Shutdown()
}
