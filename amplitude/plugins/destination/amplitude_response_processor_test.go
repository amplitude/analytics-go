package destination_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude/loggers"
	"github.com/amplitude/analytics-go/amplitude/plugins/destination"
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

func TestAmplitudeResponseProcessorSuite(t *testing.T) {
	suite.Run(t, new(AmplitudeResponseProcessorSuite))
}

type AmplitudeResponseProcessorSuite struct {
	suite.Suite
}

func (t *AmplitudeResponseProcessorSuite) TestSuccess() {
	p := destination.AmplitudeResponseProcessor{
		Logger: loggers.NewDefaultLogger(),
	}
	events := t.cloneOriginalEvents()

	result := p.Process(events, destination.AmplitudeResponse{
		Status: http.StatusOK,
		Code:   222,
	})

	require := t.Require()
	require.Equal(222, result.Code)
	require.Equal("Event sent successfully.", result.Message)
	require.Equal(len(originalEvents), len(result.Events))
	for i, event := range result.Events {
		require.Equal(originalEvents[i], *event)
	}
}

func (t *AmplitudeResponseProcessorSuite) TestTimeout() {
	events := t.cloneOriginalEvents()
	events[1].RetryCount = 2

	now := time.Now()
	retryBaseInterval := time.Second * 3
	storage := &mockEventStorage{}
	storage.On("ReturnBack", []*types.Event{events[0], events[2]}).Once()

	p := destination.AmplitudeResponseProcessor{
		EventStorage:      storage,
		MaxRetries:        2,
		RetryBaseInterval: retryBaseInterval,
		Now:               func() time.Time { return now },
		Logger:            loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
		Status: http.StatusRequestTimeout,
		Code:   408,
		Error:  "too busy",
	})

	require := t.Require()
	require.Equal(408, result.Code)
	require.Equal("Event reached max retry times 2", result.Message)
	require.Equal(1, len(result.Events))
	require.Equal(2, result.Events[0].RetryCount)
	result.Events[0].RetryCount = 0
	require.Equal(originalEvents[1], *result.Events[0])

	for _, event := range []*types.Event{events[0], events[2]} {
		require.Equal(1, event.RetryCount)
		require.Equal(now.Add(retryBaseInterval), event.RetryAt)
	}

	storage.AssertExpectations(t.T())
}

func (t *AmplitudeResponseProcessorSuite) TestTooLargeRequest_OneEvent() {
	events := t.cloneOriginalEvents()
	events = events[:1]

	p := destination.AmplitudeResponseProcessor{
		Logger: loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
		Status: http.StatusRequestEntityTooLarge,
		Code:   413,
		Error:  "too large",
	})

	require := t.Require()
	require.Equal(413, result.Code)
	require.Equal("too large", result.Message)
	require.Equal(1, len(result.Events))
	require.Equal(originalEvents[0], *result.Events[0])
}

func (t *AmplitudeResponseProcessorSuite) TestTooLargeRequest() {
	events := t.cloneOriginalEvents()

	now := time.Now()
	storage := &mockEventStorage{}
	storage.On("ReturnBack", events).Once()
	storage.On("ReduceChunkSize").Once()

	p := destination.AmplitudeResponseProcessor{
		EventStorage: storage,
		Now:          func() time.Time { return now },
		Logger:       loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
		Status: http.StatusRequestEntityTooLarge,
		Code:   413,
		Error:  "too large",
	})

	require := t.Require()
	require.Equal(413, result.Code)
	require.Equal("too large", result.Message)
	require.Empty(result.Events)

	storage.AssertExpectations(t.T())
}

func (t *AmplitudeResponseProcessorSuite) TestBadRequest_InvalidAPIKey() {
	events := t.cloneOriginalEvents()

	p := destination.AmplitudeResponseProcessor{
		Logger: loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
		Status: http.StatusBadRequest,
		Code:   400,
		Error:  "Invalid API key: info",
	})

	require := t.Require()
	require.Equal(400, result.Code)
	require.Equal("Invalid API key", result.Message)
	require.Equal(len(originalEvents), len(result.Events))
	for i, event := range result.Events {
		require.Equal(originalEvents[i], *event)
	}
}

func (t *AmplitudeResponseProcessorSuite) TestBadRequest_MissingField() {
	events := t.cloneOriginalEvents()

	p := destination.AmplitudeResponseProcessor{
		Logger: loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
		Status:       http.StatusBadRequest,
		Code:         400,
		Error:        "some error",
		MissingField: "ABC",
	})

	require := t.Require()
	require.Equal(400, result.Code)
	require.Equal("Request missing required field ABC", result.Message)
	require.Equal(len(originalEvents), len(result.Events))
	for i, event := range result.Events {
		require.Equal(originalEvents[i], *event)
	}
}

func (t *AmplitudeResponseProcessorSuite) TestBadRequest_InvalidEvents() {
	events := t.cloneOriginalEvents()

	now := time.Now()
	storage := &mockEventStorage{}
	storage.On("ReturnBack", []*types.Event{events[1]}).Once()

	p := destination.AmplitudeResponseProcessor{
		EventStorage: storage,
		Now:          func() time.Time { return now },
		Logger:       loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
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
	require.Equal(2, len(result.Events))
	require.Equal(originalEvents[0], *result.Events[0])
	require.Equal(originalEvents[2], *result.Events[1])

	storage.AssertExpectations(t.T())
}

func (t *AmplitudeResponseProcessorSuite) TestTooManyRequests() {
	events := t.cloneOriginalEvents()

	now := time.Now()
	storage := &mockEventStorage{}
	storage.On("ReturnBack", []*types.Event{events[2]}).Once()
	storage.On("ReturnBack", []*types.Event{events[0]}).Once()

	retryThrottledInterval := time.Second * 7
	p := destination.AmplitudeResponseProcessor{
		EventStorage:           storage,
		Now:                    func() time.Time { return now },
		RetryThrottledInterval: retryThrottledInterval,
		Logger:                 loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
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
	require.Equal(1, len(result.Events))
	require.Equal(originalEvents[1], *result.Events[0])

	require.Equal(0, events[0].RetryCount)
	require.Equal(time.Time{}, events[0].RetryAt)

	require.Equal(0, events[2].RetryCount)
	require.Equal(now.Add(retryThrottledInterval), events[2].RetryAt)

	storage.AssertExpectations(t.T())
}

func (t *AmplitudeResponseProcessorSuite) TestProcessUnknownError_Err() {
	events := t.cloneOriginalEvents()

	p := destination.AmplitudeResponseProcessor{
		Logger: loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
		Status: http.StatusOK,
		Code:   202,
		Err:    errors.New("some error"),
	})

	require := t.Require()
	require.Equal(202, result.Code)
	require.Equal("some error", result.Message)
	require.Equal(len(originalEvents), len(result.Events))
	for i, event := range result.Events {
		require.Equal(originalEvents[i], *event)
	}
}

func (t *AmplitudeResponseProcessorSuite) TestProcessUnknownError_ResponseError() {
	events := t.cloneOriginalEvents()

	p := destination.AmplitudeResponseProcessor{
		Logger: loggers.NewDefaultLogger(),
	}

	result := p.Process(events, destination.AmplitudeResponse{
		Status: 100,
		Code:   100,
		Error:  "some error",
	})

	require := t.Require()
	require.Equal(100, result.Code)
	require.Equal("some error", result.Message)
	require.Equal(len(originalEvents), len(result.Events))
	for i, event := range result.Events {
		require.Equal(originalEvents[i], *event)
	}
}

func (t *AmplitudeResponseProcessorSuite) cloneOriginalEvents() []*types.Event {
	events := make([]*types.Event, len(originalEvents))
	for i, originalEvent := range originalEvents {
		event := originalEvent
		events[i] = &event
	}

	return events
}

type mockEventStorage struct {
	mock.Mock
}

func (s *mockEventStorage) PushNew(event *types.Event) {
	s.Called(event)
}

func (s *mockEventStorage) ReturnBack(events ...*types.Event) {
	s.Called(events)
}

func (s *mockEventStorage) PullChunk() []*types.Event {
	args := s.Called()
	if args[0] == nil {
		return []*types.Event(nil)
	}
	return args[0].([]*types.Event)
}

func (s *mockEventStorage) HasFullChunk() bool {
	args := s.Called()

	return args.Bool(0)
}

func (s *mockEventStorage) ReduceChunkSize() {
	s.Called()
}
