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

	var config = amplitude.Config{ApiKey: "your_api_key"}

	// Config callback function (optional)
	var client = amplitude.Amplitude{Configuration: config}

	// Create a BaseEvent instance
	var event = amplitude.BaseEvent{}

	// Track an event
	client.Track(event)

	// Flush the event buffer
	client.Flush()

	// Shutdown the client
	client.Shutdown()
}
