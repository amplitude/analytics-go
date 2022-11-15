package destination_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude/plugins/destination"
	"github.com/amplitude/analytics-go/amplitude/plugins/destination/internal"
	"github.com/amplitude/analytics-go/amplitude/types"
)

type AmplitudePlugin interface {
	types.ExtendedDestinationPlugin
	SetHTTPClient(client internal.AmplitudeHTTPClient)
	SetResponseProcessor(responseProcessor internal.AmplitudeResponseProcessor)
}

func TestAmplitudePlugin(t *testing.T) {
	suite.Run(t, new(AmplitudePluginSuite))
}

type AmplitudePluginSuite struct {
	suite.Suite
}

func (t *AmplitudePluginSuite) TestIsValidAmplitudeEvent() {
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

	assert := t.Assert()

	for _, ev := range validEvents {
		assert.True(destination.IsValidAmplitudeEvent(&ev))
	}

	for _, ev := range invalidEvents {
		assert.False(destination.IsValidAmplitudeEvent(&ev))
	}
}

func (t *AmplitudePluginSuite) TestAmplitudePlugin_Basic() {
	plugin := destination.NewAmplitudePlugin()

	require := t.Require()
	require.Equal("amplitude", plugin.Name())
	require.Equal(types.PluginTypeDestination, plugin.Type())
}

func (t *AmplitudePluginSuite) TestAmplitudePlugin_FlushInterval() {
	plugin := destination.NewAmplitudePlugin().(AmplitudePlugin)

	flushInterval := time.Millisecond * 100
	flushQueueSize := 3
	event1 := t.createEvent(1)
	event2 := t.createEvent(2)
	storageEvent1 := &types.StorageEvent{Event: event1}
	storageEvent2 := &types.StorageEvent{Event: event2}

	storage := &mockStorage{}
	storage.On("PushNew", storageEvent1).Once()
	storage.On("Count", mock.Anything).Return(1).Once()
	storage.On("PushNew", storageEvent2).Once()
	storage.On("Count", mock.Anything).Return(2).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return([]*types.StorageEvent{storageEvent1, storageEvent2}).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return(nil).Once()

	httpClient := &mockHTTPClient{}
	httpClient.On("Send", internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{event1, event2},
	}).Return(internal.AmplitudeResponse{Status: 200}).Once()

	responseProcessor := &mockResponseProcessor{}
	responseProcessor.On("Process", []*types.StorageEvent{storageEvent1, storageEvent2}, internal.AmplitudeResponse{Status: 200}).Return(internal.AmplitudeProcessorResult{
		Code:              200,
		EventsForCallback: []*types.StorageEvent{storageEvent1, storageEvent2},
	})

	plugin.SetHTTPClient(httpClient)
	plugin.SetResponseProcessor(responseProcessor)

	plugin.Setup(types.Config{
		APIKey:             "my-api-key",
		MaxStorageCapacity: 10,
		FlushInterval:      flushInterval,
		FlushQueueSize:     flushQueueSize,
		FlushSizeDivider:   1,
		StorageFactory: func() types.EventStorage {
			return storage
		},
		Logger: noopLogger{},
	})

	plugin.Execute(event1)
	plugin.Execute(event2)

	time.Sleep(flushInterval + time.Millisecond*50)

	httpClient.AssertExpectations(t.T())
	responseProcessor.AssertExpectations(t.T())
	storage.AssertExpectations(t.T())

	storage.On("Pull", flushQueueSize, mock.Anything).Return(nil).Once()

	plugin.Shutdown()
}

func (t *AmplitudePluginSuite) TestAmplitudePlugin_FlushQueueSize() {
	plugin := destination.NewAmplitudePlugin().(AmplitudePlugin)

	flushInterval := time.Second * 100
	flushQueueSize := 2
	event1 := t.createEvent(1)
	event2 := t.createEvent(2)
	storageEvent1 := &types.StorageEvent{Event: event1}
	storageEvent2 := &types.StorageEvent{Event: event2}

	storage := &mockStorage{}
	storage.On("PushNew", storageEvent1).Once()
	storage.On("Count", mock.Anything).Return(1).Once()
	storage.On("PushNew", storageEvent2).Once()
	storage.On("Count", mock.Anything).Return(2).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return([]*types.StorageEvent{storageEvent1, storageEvent2}).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return(nil).Once()

	httpClient := &mockHTTPClient{}
	httpClient.On("Send", internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{event1, event2},
	}).Return(internal.AmplitudeResponse{Status: 200}).Once()

	responseProcessor := &mockResponseProcessor{}
	responseProcessor.On("Process", []*types.StorageEvent{storageEvent1, storageEvent2}, internal.AmplitudeResponse{Status: 200}).Return(internal.AmplitudeProcessorResult{
		Code:              200,
		EventsForCallback: []*types.StorageEvent{storageEvent1, storageEvent2},
	})

	plugin.SetHTTPClient(httpClient)
	plugin.SetResponseProcessor(responseProcessor)

	plugin.Setup(types.Config{
		APIKey:             "my-api-key",
		MaxStorageCapacity: 10,
		FlushInterval:      flushInterval,
		FlushQueueSize:     flushQueueSize,
		FlushSizeDivider:   1,
		StorageFactory: func() types.EventStorage {
			return storage
		},
		Logger: noopLogger{},
	})

	plugin.Execute(event1)
	plugin.Execute(event2)

	time.Sleep(time.Millisecond * 50)

	httpClient.AssertExpectations(t.T())
	responseProcessor.AssertExpectations(t.T())
	storage.AssertExpectations(t.T())

	storage.On("Pull", flushQueueSize, mock.Anything).Return(nil).Once()

	plugin.Shutdown()
}

func (t *AmplitudePluginSuite) TestAmplitudePlugin_ExplicitFlush() {
	plugin := destination.NewAmplitudePlugin().(AmplitudePlugin)

	flushInterval := time.Second * 100
	flushQueueSize := 10
	event1 := t.createEvent(1)
	event2 := t.createEvent(2)
	storageEvent1 := &types.StorageEvent{Event: event1}
	storageEvent2 := &types.StorageEvent{Event: event2}

	storage := &mockStorage{}
	storage.On("PushNew", storageEvent1).Once()
	storage.On("Count", mock.Anything).Return(1).Once()
	storage.On("PushNew", storageEvent2).Once()
	storage.On("Count", mock.Anything).Return(2).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return([]*types.StorageEvent{storageEvent1, storageEvent2}).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return(nil).Once()

	httpClient := &mockHTTPClient{}
	httpClient.On("Send", internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{event1, event2},
	}).Return(internal.AmplitudeResponse{Status: 200}).Once()

	responseProcessor := &mockResponseProcessor{}
	responseProcessor.On("Process", []*types.StorageEvent{storageEvent1, storageEvent2}, internal.AmplitudeResponse{Status: 200}).Return(internal.AmplitudeProcessorResult{
		Code:              200,
		EventsForCallback: []*types.StorageEvent{storageEvent1, storageEvent2},
	})

	plugin.SetHTTPClient(httpClient)
	plugin.SetResponseProcessor(responseProcessor)

	plugin.Setup(types.Config{
		APIKey:             "my-api-key",
		MaxStorageCapacity: 10,
		FlushInterval:      flushInterval,
		FlushQueueSize:     flushQueueSize,
		FlushSizeDivider:   1,
		StorageFactory: func() types.EventStorage {
			return storage
		},
		Logger: noopLogger{},
	})

	plugin.Execute(event1)
	plugin.Execute(event2)
	plugin.Flush()

	time.Sleep(time.Millisecond * 50)

	httpClient.AssertExpectations(t.T())
	responseProcessor.AssertExpectations(t.T())
	storage.AssertExpectations(t.T())

	storage.On("Pull", flushQueueSize, mock.Anything).Return(nil).Once()

	plugin.Shutdown()
}

func (t *AmplitudePluginSuite) TestAmplitudePlugin_ExplicitShutdown() {
	plugin := destination.NewAmplitudePlugin().(AmplitudePlugin)

	flushInterval := time.Second * 100
	flushQueueSize := 10
	event1 := t.createEvent(1)
	event2 := t.createEvent(2)
	storageEvent1 := &types.StorageEvent{Event: event1}
	storageEvent2 := &types.StorageEvent{Event: event2}

	storage := &mockStorage{}
	storage.On("PushNew", storageEvent1).Once()
	storage.On("Count", mock.Anything).Return(1).Once()
	storage.On("PushNew", storageEvent2).Once()
	storage.On("Count", mock.Anything).Return(2).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return([]*types.StorageEvent{storageEvent1, storageEvent2}).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return(nil).Once()

	httpClient := &mockHTTPClient{}
	httpClient.On("Send", internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{event1, event2},
	}).Return(internal.AmplitudeResponse{Status: 200}).Once()

	responseProcessor := &mockResponseProcessor{}
	responseProcessor.On("Process", []*types.StorageEvent{storageEvent1, storageEvent2}, internal.AmplitudeResponse{Status: 200}).Return(internal.AmplitudeProcessorResult{
		Code:              200,
		EventsForCallback: []*types.StorageEvent{storageEvent1, storageEvent2},
	})

	plugin.SetHTTPClient(httpClient)
	plugin.SetResponseProcessor(responseProcessor)

	plugin.Setup(types.Config{
		APIKey:             "my-api-key",
		MaxStorageCapacity: 10,
		FlushInterval:      flushInterval,
		FlushQueueSize:     flushQueueSize,
		FlushSizeDivider:   1,
		StorageFactory: func() types.EventStorage {
			return storage
		},
		Logger: noopLogger{},
	})

	plugin.Execute(event1)
	plugin.Execute(event2)
	plugin.Shutdown()

	time.Sleep(time.Millisecond * 50)

	httpClient.AssertExpectations(t.T())
	responseProcessor.AssertExpectations(t.T())
	storage.AssertExpectations(t.T())
}

func (t *AmplitudePluginSuite) TestAmplitudePlugin_ReduceChunkSize() {
	plugin := destination.NewAmplitudePlugin().(AmplitudePlugin)

	flushInterval := time.Second * 100
	flushQueueSize := 2
	event1 := t.createEvent(1)
	event2 := t.createEvent(2)
	storageEvent1 := &types.StorageEvent{Event: event1}
	storageEvent2 := &types.StorageEvent{Event: event2}

	storage := &mockStorage{}
	storage.On("PushNew", storageEvent1).Once()
	storage.On("Count", mock.Anything).Return(1).Once()
	storage.On("PushNew", storageEvent2).Once()
	storage.On("Count", mock.Anything).Return(2).Once()
	storage.On("Pull", flushQueueSize, mock.Anything).Return([]*types.StorageEvent{storageEvent1, storageEvent2}).Once()
	storage.On("Pull", flushQueueSize/2, mock.Anything).Return([]*types.StorageEvent{storageEvent1}).Once()
	storage.On("Pull", flushQueueSize/2, mock.Anything).Return([]*types.StorageEvent{storageEvent2}).Once()
	storage.On("Pull", flushQueueSize/2, mock.Anything).Return(nil).Once()

	httpClient := &mockHTTPClient{}
	httpClient.On("Send", internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{event1, event2},
	}).Return(internal.AmplitudeResponse{Status: http.StatusRequestEntityTooLarge}).Once()
	httpClient.On("Send", internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{event1},
	}).Return(internal.AmplitudeResponse{Status: 200}).Once()
	httpClient.On("Send", internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{event2},
	}).Return(internal.AmplitudeResponse{Status: 200}).Once()

	responseProcessor := &mockResponseProcessor{}
	responseProcessor.On("Process", []*types.StorageEvent{storageEvent1, storageEvent2}, internal.AmplitudeResponse{Status: http.StatusRequestEntityTooLarge}).Return(internal.AmplitudeProcessorResult{
		Code:           http.StatusRequestEntityTooLarge,
		EventsForRetry: []*types.StorageEvent{storageEvent1, storageEvent2},
	})
	responseProcessor.On("Process", []*types.StorageEvent{storageEvent1}, internal.AmplitudeResponse{Status: 200}).Return(internal.AmplitudeProcessorResult{
		Code:              200,
		EventsForCallback: []*types.StorageEvent{storageEvent1},
	})
	responseProcessor.On("Process", []*types.StorageEvent{storageEvent2}, internal.AmplitudeResponse{Status: 200}).Return(internal.AmplitudeProcessorResult{
		Code:              200,
		EventsForCallback: []*types.StorageEvent{storageEvent2},
	})

	plugin.SetHTTPClient(httpClient)
	plugin.SetResponseProcessor(responseProcessor)

	plugin.Setup(types.Config{
		APIKey:             "my-api-key",
		MaxStorageCapacity: 10,
		FlushInterval:      flushInterval,
		FlushQueueSize:     flushQueueSize,
		FlushSizeDivider:   1,
		StorageFactory: func() types.EventStorage {
			return storage
		},
		Logger: noopLogger{},
	})

	plugin.Execute(event1)
	plugin.Execute(event2)

	time.Sleep(time.Millisecond * 50)

	httpClient.AssertExpectations(t.T())
	responseProcessor.AssertExpectations(t.T())
	storage.AssertExpectations(t.T())

	storage.On("Pull", flushQueueSize/2, mock.Anything).Return(nil).Once()

	plugin.Shutdown()
}

func (t *AmplitudePluginSuite) createEvent(index int) *types.Event {
	postfix := fmt.Sprintf("-%d", index)

	return &types.Event{
		EventType: "event" + postfix,
		EventOptions: types.EventOptions{
			UserID:   "user" + postfix,
			InsertID: "insert" + postfix,
			Time:     int64(index),
		},
		EventProperties: map[string]interface{}{
			"prop" + postfix: index,
		},
	}
}

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Send(payload internal.AmplitudePayload) internal.AmplitudeResponse {
	args := m.Called(payload)

	return args[0].(internal.AmplitudeResponse)
}

type mockResponseProcessor struct {
	mock.Mock
}

func (m *mockResponseProcessor) Process(events []*types.StorageEvent, response internal.AmplitudeResponse) internal.AmplitudeProcessorResult {
	args := m.Called(events, response)

	return args[0].(internal.AmplitudeProcessorResult)
}

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) PushNew(event *types.StorageEvent) {
	m.Called(event)
}

func (m *mockStorage) ReturnBack(events ...*types.StorageEvent) {
	m.Called(events)
}

func (m *mockStorage) Pull(count int, before time.Time) []*types.StorageEvent {
	args := m.Called(count, before)
	if args[0] == nil {
		return []*types.StorageEvent(nil)
	}

	return args[0].([]*types.StorageEvent)
}

func (m *mockStorage) Count(before time.Time) int {
	args := m.Called(before)

	return args.Int(0)
}

type noopLogger struct{}

func (l noopLogger) Debugf(string, ...interface{}) {
}

func (l noopLogger) Infof(string, ...interface{}) {
}

func (l noopLogger) Warnf(string, ...interface{}) {
}

func (l noopLogger) Errorf(string, ...interface{}) {
}
