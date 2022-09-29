package storages_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude/storages"
	"github.com/amplitude/analytics-go/amplitude/types"
)

func TestInMemoryEventStorage(t *testing.T) {
	suite.Run(t, new(InMemoryEventStorageSuite))
}

type InMemoryEventStorageSuite struct {
	suite.Suite
}

func (t *InMemoryEventStorageSuite) TestSimple() {
	event1 := &types.Event{EventType: "event-A"}
	event2 := &types.Event{EventType: "event-B"}
	event3 := &types.Event{EventType: "event-C"}
	event4 := &types.Event{EventType: "event-D"}

	require := t.Require()

	s := storages.NewInMemoryEventStorage(3, 1)
	s.PushNew(event1)
	require.False(s.HasFullChunk())
	s.PushNew(event2)
	require.False(s.HasFullChunk())
	s.PushNew(event3)
	require.True(s.HasFullChunk())
	s.PushNew(event4)
	require.True(s.HasFullChunk())

	chunk := s.PullChunk()
	require.False(s.HasFullChunk())
	require.Equal([]*types.Event{event1, event2, event3}, chunk)

	s.PushNew(event2)
	require.False(s.HasFullChunk())
	s.PushNew(event3)
	require.True(s.HasFullChunk())
	s.PushNew(event1)
	require.True(s.HasFullChunk())

	chunk = s.PullChunk()
	require.False(s.HasFullChunk())
	require.Equal([]*types.Event{event4, event2, event3}, chunk)

	chunk = s.PullChunk()
	require.False(s.HasFullChunk())
	require.Equal([]*types.Event{event1}, chunk)

	chunk = s.PullChunk()
	require.False(s.HasFullChunk())
	require.Empty(chunk)
}

func (t *InMemoryEventStorageSuite) TestReturnBack() {
	event1 := &types.Event{EventType: "event-A"}
	event2 := &types.Event{EventType: "event-B"}
	event3 := &types.Event{EventType: "event-C"}
	event4 := &types.Event{EventType: "event-D"}
	event5 := &types.Event{EventType: "event-E"}

	require := t.Require()

	s := storages.NewInMemoryEventStorage(4, 1)
	s.PushNew(event1)
	s.PushNew(event2)
	s.PushNew(event3)
	s.PushNew(event4)
	s.PushNew(event5)

	chunk := s.PullChunk()
	require.Equal([]*types.Event{event1, event2, event3, event4}, chunk)

	event3.RetryAt = time.Now().Add(time.Millisecond * 300)
	event4.RetryAt = time.Now().Add(time.Millisecond * 100)
	s.ReturnBack(event2, event3, event4, event1)

	chunk = s.PullChunk()
	require.Equal([]*types.Event{event2, event1, event5}, chunk)

	chunk = s.PullChunk()
	require.Empty(chunk)

	time.Sleep(time.Millisecond * 150)

	chunk = s.PullChunk()
	require.Equal([]*types.Event{event4}, chunk)

	chunk = s.PullChunk()
	require.Empty(chunk)

	time.Sleep(time.Millisecond * 200)

	chunk = s.PullChunk()
	require.Equal([]*types.Event{event3}, chunk)

	chunk = s.PullChunk()
	require.Empty(chunk)
}

func (t *InMemoryEventStorageSuite) TestReduceChunk() {
	event1 := &types.Event{EventType: "event-A"}
	event2 := &types.Event{EventType: "event-B"}
	event3 := &types.Event{EventType: "event-C"}
	event4 := &types.Event{EventType: "event-D"}
	event5 := &types.Event{EventType: "event-E"}
	event6 := &types.Event{EventType: "event-F"}
	event7 := &types.Event{EventType: "event-G"}
	event8 := &types.Event{EventType: "event-H"}

	require := t.Require()

	s := storages.NewInMemoryEventStorage(4, 1)
	s.PushNew(event1)
	s.PushNew(event2)
	s.PushNew(event3)
	s.PushNew(event4)
	s.PushNew(event5)
	s.PushNew(event6)
	s.PushNew(event7)
	s.PushNew(event8)

	chunk := s.PullChunk()
	require.Equal([]*types.Event{event1, event2, event3, event4}, chunk)

	s.ReduceChunkSize()
	chunk = s.PullChunk()
	require.Equal([]*types.Event{event5, event6}, chunk)

	s.ReduceChunkSize()
	chunk = s.PullChunk()
	require.Equal([]*types.Event{event7}, chunk)

	s.ReduceChunkSize()
	chunk = s.PullChunk()
	require.Equal([]*types.Event{event8}, chunk)

	s.ReduceChunkSize()
	chunk = s.PullChunk()
	require.Empty(chunk)
}
