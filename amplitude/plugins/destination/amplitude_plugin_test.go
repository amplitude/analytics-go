package destination

import (
	"testing"

	"github.com/amplitude/analytics-go/amplitude/types"
	"github.com/stretchr/testify/assert"
)

func Test_isValidEvent(t *testing.T) {

	validEvents := []types.Event{
		{
			EventType: "type",
			UserID:    "id",
		}, {
			EventType: "type",
			DeviceID:  "id",
		}, {
			EventType:    "type",
			EventOptions: types.EventOptions{UserID: "id"},
		}, {
			EventType:    "type",
			EventOptions: types.EventOptions{DeviceID: "id"},
		}, {
			EventType:    "type",
			UserID:       "id",
			EventOptions: types.EventOptions{UserID: "id"},
		},
	}

	invalidEvents := []types.Event{
		{},
		{
			UserID: "id",
		},
		{
			DeviceID: "id",
		},
		{
			EventOptions: types.EventOptions{UserID: "id"},
		},
		{
			EventOptions: types.EventOptions{DeviceID: "id"},
		},
	}

	for _, ev := range validEvents {
		assert.True(t, isValidEvent(&ev))
	}

	for _, ev := range invalidEvents {
		assert.False(t, isValidEvent(&ev))
	}
}
