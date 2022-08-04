// A basic example of using Amplitude Go SDK to set user property

package main

// Import amplitude package
import (
	"github.com/amplitude/Amplitude-Go/amplitude"
)

func main() {

	config := amplitude.NewConfig("c253b75dce3e593c44ea5eb95999f92a")

	client := amplitude.NewClient(config)

	// Revenue struct is passed into Revenue method
	// to send as a revenue event
	revenueObj := amplitude.Revenue{
		Price:     3.99,
		Quantity:  3,
		ProductID: "com.company.productID",
	}
	client.Revenue(revenueObj, amplitude.EventOptions{DeviceID: "revenue-device-id", UserID: "revenue-user-id"})

	// Flush the event buffer
	client.Flush()

	// Shutdown the client
	client.Shutdown()
}
