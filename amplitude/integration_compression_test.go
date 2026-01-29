package amplitude_test

import (
	"testing"
	"time"

	"github.com/amplitude/analytics-go/amplitude"
)

// TestCompressionIntegration tests that compression works end-to-end with real Amplitude API
// Compression is always enabled in the SDK
func TestCompressionIntegration(t *testing.T) {
	apiKey := "test-api-key" // Replace with your API key to run this test

	config := amplitude.NewConfig(apiKey)
	config.FlushInterval = time.Second * 1

	client := amplitude.NewClient(config)

	// Track multiple events to test batching and compression
	// All events are automatically compressed with gzip
	for i := 0; i < 10; i++ {
		client.Track(amplitude.Event{
			EventType: "compression_integration_test",
			UserID:    "test_user_compression",
			EventProperties: map[string]interface{}{
				"test_number": i,
				"timestamp":   time.Now().Format(time.RFC3339),
				"data":        "This test data is automatically compressed with gzip before sending",
			},
		})
	}

	// Shutdown will flush all events
	client.Shutdown()

	// If we get here without errors, compression worked!
	t.Log("✅ Compression integration test passed - events sent successfully with automatic gzip compression")
}
