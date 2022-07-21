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

	config := amplitude.NewConfig("your_api_key")

	// Config callback function (optional)
	client := amplitude.NewClient(config)

	client.Add(amplitude.NewContextPlugin())

	// Create a BaseEvent instance
	event := amplitude.Event{}

	// Track an event
	client.Track(event)

	// Flush the event buffer
	client.Flush()

	// Shutdown the client
	client.Shutdown()
}
