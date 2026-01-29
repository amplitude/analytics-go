// Example demonstrating gzip compression feature in Amplitude Go SDK

package main

import (
	"fmt"

	"github.com/amplitude/analytics-go/amplitude"
)

func main() {
	// Create config with your API key
	config := amplitude.NewConfig("your-api-key")

	// Gzip compression is ALWAYS ENABLED automatically
	// This reduces bandwidth usage by ~70-90% with no configuration needed
	fmt.Println("The Amplitude Go SDK automatically compresses all events with gzip")

	client := amplitude.NewClient(config)

	// Track events as usual - compression happens automatically
	for i := 0; i < 10; i++ {
		client.Track(amplitude.Event{
			EventType: "page_view",
			UserID:    "user-123",
			EventProperties: map[string]interface{}{
				"page":       fmt.Sprintf("/page-%d", i),
				"session_id": "abc-123",
				"referrer":   "https://example.com",
			},
		})
	}

	fmt.Println("Tracked 10 events - they are automatically compressed before sending to Amplitude")

	// Flush and shutdown
	client.Shutdown()

	fmt.Println("Done! Events sent with automatic gzip compression.")
	fmt.Println("\nBenefits of automatic compression:")
	fmt.Println("  - 70-90% bandwidth reduction")
	fmt.Println("  - Faster uploads on slow connections")
	fmt.Println("  - Lower network costs")
	fmt.Println("  - Supported by both HTTP V2 and Batch APIs")
	fmt.Println("  - No configuration needed - always works!")
}
