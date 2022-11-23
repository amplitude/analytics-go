package types_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude/types"
)

func TestEvent(t *testing.T) {
	suite.Run(t, new(EventSuite))
}

type EventSuite struct {
	suite.Suite
}

func (t *EventSuite) TestClone() {
	original := types.Event{
		EventType: "my-event",
		EventOptions: types.EventOptions{
			UserID:      "my-user",
			DeviceID:    "my-device",
			Time:        123,
			LocationLat: 5.67,
		},
		EventProperties: map[string]interface{}{
			"prop-string":  "string value",
			"prop-int":     1,
			"prop-float":   2.3,
			"prop-boolean": true,
			"prop-object": map[string]interface{}{
				"prop-string":        "string value",
				"prop-int":           1,
				"prop-float":         2.3,
				"prop-boolean":       true,
				"prop-array-string":  []string{"string", "value"},
				"prop-array-int":     []int{1, 2, 3},
				"prop-array-float":   []float64{1.23, 4.56},
				"prop-array-boolean": []bool{true, false},
				"prop-array-object": []interface{}{
					map[string]interface{}{"prop-1": 123},
				},
			},
			"prop-array-string":  []string{"string", "value"},
			"prop-array-int":     []int{1, 2, 3},
			"prop-array-float":   []float64{1.23, 4.56},
			"prop-array-boolean": []bool{true, false},
			"prop-array-object": []interface{}{
				map[string]interface{}{"prop-1": 123},
			},
		},
		UserProperties: map[types.IdentityOp]map[string]interface{}{
			types.IdentityOpSet: {
				"prop-string":        "string value",
				"prop-int":           1,
				"prop-float":         2.3,
				"prop-boolean":       true,
				"prop-array-string":  []string{"string", "value"},
				"prop-array-int":     []int{1, 2, 3},
				"prop-array-float":   []float64{1.23, 4.56},
				"prop-array-boolean": []bool{true, false},
				"prop-array-object": []interface{}{
					map[string]interface{}{"prop-1": 123},
				},
			},
		},
		Groups: map[string][]string{
			"group-1": {"group-A", "group-B"},
			"group-2": {"group-C"},
		},
		GroupProperties: map[types.IdentityOp]map[string]interface{}{
			types.IdentityOpAdd: {
				"prop-array-string": []string{"string", "value"},
				"prop-array-object": []interface{}{
					map[string]interface{}{"prop-1": 123},
				},
			},
			types.IdentityOpClearAll: {
				"prop-string": "string value",
			},
		},
		UserID:   "user-1",
		DeviceID: "device-1",
	}

	clone := original.Clone()

	require := t.Require()

	require.Equal(original, clone)

	require.NotSame(original, clone)
	require.NotSame(original.EventOptions, clone.EventOptions)

	require.NotSame(original.EventProperties, clone.EventProperties)
	require.NotSame(original.EventProperties["prop-array-string"], clone.EventProperties["prop-array-string"])
	require.NotSame(original.EventProperties["prop-array-int"], clone.EventProperties["prop-array-int"])
	require.NotSame(original.EventProperties["prop-array-float"], clone.EventProperties["prop-array-float"])
	require.NotSame(original.EventProperties["prop-array-boolean"], clone.EventProperties["prop-array-boolean"])
	require.NotSame(original.EventProperties["prop-array-object"], clone.EventProperties["prop-array-object"])
	require.NotSame(original.EventProperties["prop-object"], clone.EventProperties["prop-objectt"])

	require.NotSame(original.UserProperties, clone.UserProperties)
	require.NotSame(original.UserProperties[types.IdentityOpSet], clone.UserProperties[types.IdentityOpSet])

	require.NotSame(original.Groups, clone.Groups)
	require.NotSame(original.Groups["group-1"], clone.Groups["group-1"])
	require.NotSame(original.Groups["group-2"], clone.Groups["group-2"])

	require.NotSame(original.GroupProperties, clone.GroupProperties)
	require.NotSame(original.GroupProperties[types.IdentityOpAdd], clone.GroupProperties[types.IdentityOpAdd])
	require.NotSame(original.GroupProperties[types.IdentityOpClearAll], clone.GroupProperties[types.IdentityOpClearAll])
}
