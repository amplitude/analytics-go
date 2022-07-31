package amplitude

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChunk(t *testing.T) {
	amplitudeDestinationPlugin := AmplitudeDestinationPlugin{
		config: Config{
			FlushQueueSize: 3,
		},
		scheduled: false,
	}
	events := make([]*Event, 10)
	for index, _ := range events {
		events[index] = &Event{
			EventOptions: EventOptions{
				UserID: "user-" + fmt.Sprint(index),
			},
		}
	}

	chunks := amplitudeDestinationPlugin.chunk(events)
	assert.Equal(t, 4, len(chunks))
	assert.Equal(t, 1, len(chunks[len(chunks)-1]))

}
