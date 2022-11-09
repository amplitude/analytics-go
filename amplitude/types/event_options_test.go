package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude/types"
)

func TestEventOptions(t *testing.T) {
	suite.Run(t, new(EventOptionsSuite))
}

type EventOptionsSuite struct {
	suite.Suite
}

func (t *EventOptionsSuite) TestClone() {
	original := types.EventOptions{
		UserID:      "my-user",
		DeviceID:    "my-device",
		Time:        123,
		LocationLat: 5.67,
		Plan: &types.Plan{
			Branch: "my-branch",
		},
		IngestionMetadata: &types.IngestionMetadata{
			SourceName: "my-source",
		},
	}

	clone := *original.Clone()

	require := t.Require()

	require.Equal(original, clone)

	require.NotSame(original, clone)
	require.NotSame(original.Plan, clone.Plan)
	require.NotSame(original.IngestionMetadata, clone.IngestionMetadata)
}

func (t *EventOptionsSuite) TestSetTime() {
	options := types.EventOptions{}
	options.SetTime(time.Date(2022, 11, 12, 13, 14, 15, 0, time.UTC))

	require := t.Require()
	require.Equal(int64(1668258855000), options.Time)
}
