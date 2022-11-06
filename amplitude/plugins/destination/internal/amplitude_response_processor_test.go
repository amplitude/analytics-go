package internal_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/amplitude/analytics-go/amplitude/plugins/destination/internal"

	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude/loggers"
	"github.com/amplitude/analytics-go/amplitude/types"
)

var originalEvents = []types.Event{
	{
		EventType:       "event-A",
		EventProperties: map[string]interface{}{"property-1": "value-1"},
		UserID:          "user-1",
	},
	{
		EventType:       "event-B",
		EventProperties: map[string]interface{}{"property-2": 2},
		UserID:          "user-2",
	},
	{
		EventType:       "event-C",
		EventProperties: map[string]interface{}{"property-3": true},
		UserID:          "user-1",
	},
}

func TestAmplitudeResponseProcessor(t *testing.T) {
	suite.Run(t, new(AmplitudeResponseProcessorSuite))
}

type AmplitudeResponseProcessorSuite struct {
	suite.Suite
}

func (t *AmplitudeResponseProcessorSuite) TestSuccess() {
	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Logger: loggers.NewDefaultLogger(),
	})
	events := t.cloneOriginalEvents()

	result := p.Process(events, internal.AmplitudeResponse{
		Status: http.StatusOK,
		Code:   222,
	})

	require := t.Require()
	require.Equal(222, result.Code)
	require.Equal("Event sent successfully.", result.Message)
	require.Equal(len(originalEvents), len(result.EventsForCallback))
	require.Equal(0, len(result.EventsForRetry))

	for i, event := range result.EventsForCallback {
		require.Equal(originalEvents[i], *event.Event)
	}
}

func (t *AmplitudeResponseProcessorSuite) TestTimeout() {
	events := t.cloneOriginalEvents()
	events[1].RetryCount = 2

	now := time.Now()
	retryBaseInterval := time.Second * 3

	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		MaxRetries:        2,
		RetryBaseInterval: retryBaseInterval,
		Now:               func() time.Time { return now },
		Logger:            loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status: http.StatusRequestTimeout,
		Code:   408,
		Error:  "too busy",
	})

	require := t.Require()
	require.Equal(408, result.Code)
	require.Equal("Event reached max retry times 2", result.Message)
	require.Equal(1, len(result.EventsForCallback))
	require.Equal(2, result.EventsForCallback[0].RetryCount)
	require.Equal(originalEvents[1], *result.EventsForCallback[0].Event)

	require.Equal(2, len(result.EventsForRetry))

	for i, originalEvent := range []types.Event{originalEvents[0], originalEvents[2]} {
		event := *result.EventsForRetry[i]
		require.Equal(originalEvent, *event.Event)
		require.Equal(1, event.RetryCount)
		require.Equal(now.Add(retryBaseInterval), event.RetryAt)
	}
}

func (t *AmplitudeResponseProcessorSuite) TestTooLargeRequest_OneEvent() {
	events := t.cloneOriginalEvents()
	events = events[:1]

	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Logger: loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status: http.StatusRequestEntityTooLarge,
		Code:   413,
		Error:  "too large",
	})

	require := t.Require()
	require.Equal(413, result.Code)
	require.Equal("too large", result.Message)
	require.Equal(1, len(result.EventsForCallback))
	require.Equal(originalEvents[0], *result.EventsForCallback[0].Event)

	require.Equal(0, len(result.EventsForRetry))
}

func (t *AmplitudeResponseProcessorSuite) TestTooLargeRequest() {
	events := t.cloneOriginalEvents()

	now := time.Now()

	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Now:    func() time.Time { return now },
		Logger: loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status: http.StatusRequestEntityTooLarge,
		Code:   413,
		Error:  "too large",
	})

	require := t.Require()
	require.Equal(413, result.Code)
	require.Equal("too large", result.Message)
	require.Empty(result.EventsForCallback)

	require.Equal(len(originalEvents), len(result.EventsForRetry))

	for i, originalEvent := range originalEvents {
		event := *result.EventsForRetry[i]
		require.Equal(originalEvent, *event.Event)
		require.Equal(0, event.RetryCount)
		require.Empty(event.RetryAt)
	}
}

func (t *AmplitudeResponseProcessorSuite) TestBadRequest_InvalidAPIKey() {
	events := t.cloneOriginalEvents()

	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Logger: loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status: http.StatusBadRequest,
		Code:   400,
		Error:  "Invalid API key: info",
	})

	require := t.Require()
	require.Equal(400, result.Code)
	require.Equal("Invalid API key", result.Message)
	require.Equal(len(originalEvents), len(result.EventsForCallback))

	for i, event := range result.EventsForCallback {
		require.Equal(originalEvents[i], *event.Event)
	}

	require.Equal(0, len(result.EventsForRetry))
}

func (t *AmplitudeResponseProcessorSuite) TestBadRequest_MissingField() {
	events := t.cloneOriginalEvents()

	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Logger: loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status:       http.StatusBadRequest,
		Code:         400,
		Error:        "some error",
		MissingField: "ABC",
	})

	require := t.Require()
	require.Equal(400, result.Code)
	require.Equal("Request missing required field ABC", result.Message)
	require.Equal(len(originalEvents), len(result.EventsForCallback))

	for i, event := range result.EventsForCallback {
		require.Equal(originalEvents[i], *event.Event)
	}

	require.Equal(0, len(result.EventsForRetry))
}

func (t *AmplitudeResponseProcessorSuite) TestBadRequest_InvalidEvents() {
	events := t.cloneOriginalEvents()

	now := time.Now()

	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Now:    func() time.Time { return now },
		Logger: loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status:         http.StatusBadRequest,
		Code:           400,
		Error:          "some error",
		SilencedEvents: []int{2},
		EventsWithInvalidFields: map[string][]int{
			"abc": {0},
		},
	})

	require := t.Require()
	require.Equal(400, result.Code)
	require.Equal("some error", result.Message)
	require.Equal(2, len(result.EventsForCallback))
	require.Equal(originalEvents[0], *result.EventsForCallback[0].Event)
	require.Equal(originalEvents[2], *result.EventsForCallback[1].Event)

	require.Equal(1, len(result.EventsForRetry))

	for i, originalEvent := range []types.Event{originalEvents[1]} {
		event := *result.EventsForRetry[i]
		require.Equal(originalEvent, *event.Event)
		require.Equal(0, event.RetryCount)
		require.Empty(event.RetryAt)
	}
}

func (t *AmplitudeResponseProcessorSuite) TestTooManyRequests() {
	events := t.cloneOriginalEvents()

	now := time.Now()

	retryThrottledInterval := time.Second * 7
	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Now:                    func() time.Time { return now },
		RetryThrottledInterval: retryThrottledInterval,
		Logger:                 loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status:          http.StatusTooManyRequests,
		Code:            429,
		Error:           "some error",
		ThrottledEvents: []int{1, 2},
		ExceededDailyQuotaUsers: map[string]int{
			"user-2": 100,
		},
	})

	require := t.Require()
	require.Equal(429, result.Code)
	require.Equal("Exceeded daily quota", result.Message)
	require.Equal(1, len(result.EventsForCallback))
	require.Equal(originalEvents[1], *result.EventsForCallback[0].Event)

	require.Equal(2, len(result.EventsForRetry))

	for i, originalEvent := range []types.Event{originalEvents[0], originalEvents[2]} {
		event := *result.EventsForRetry[i]
		require.Equal(originalEvent, *event.Event)
	}

	require.Equal(0, events[0].RetryCount)
	require.Empty(events[0].RetryAt)

	require.Equal(0, events[2].RetryCount)
	require.Equal(now.Add(retryThrottledInterval), events[2].RetryAt)
}

func (t *AmplitudeResponseProcessorSuite) TestProcessUnknownError_Err() {
	events := t.cloneOriginalEvents()

	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Logger: loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status: http.StatusOK,
		Code:   202,
		Err:    errors.New("some error"),
	})

	require := t.Require()
	require.Equal(202, result.Code)
	require.Equal("some error", result.Message)
	require.Equal(len(originalEvents), len(result.EventsForCallback))

	for i, event := range result.EventsForCallback {
		require.Equal(originalEvents[i], *event.Event)
	}

	require.Equal(0, len(result.EventsForRetry))
}

func (t *AmplitudeResponseProcessorSuite) TestProcessUnknownError_ResponseError() {
	events := t.cloneOriginalEvents()

	p := internal.NewAmplitudeResponseProcessor(internal.AmplitudeResponseProcessorOptions{
		Logger: loggers.NewDefaultLogger(),
	})

	result := p.Process(events, internal.AmplitudeResponse{
		Status: 100,
		Code:   100,
		Error:  "some error",
	})

	require := t.Require()
	require.Equal(100, result.Code)
	require.Equal("some error", result.Message)
	require.Equal(len(originalEvents), len(result.EventsForCallback))

	for i, event := range result.EventsForCallback {
		require.Equal(originalEvents[i], *event.Event)
	}

	require.Equal(0, len(result.EventsForRetry))
}

func (t *AmplitudeResponseProcessorSuite) cloneOriginalEvents() []*types.StorageEvent {
	events := make([]*types.StorageEvent, len(originalEvents))

	for i, originalEvent := range originalEvents {
		event := originalEvent
		events[i] = &types.StorageEvent{Event: &event}
	}

	return events
}
