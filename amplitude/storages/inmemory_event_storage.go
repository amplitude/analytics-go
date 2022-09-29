package storages

import (
	"sort"
	"sync"
	"time"

	"github.com/amplitude/analytics-go/amplitude/types"
)

func NewInMemoryEventStorage(chunkSize, sizeDivider int) types.EventStorage {
	if sizeDivider < 1 {
		sizeDivider = 1
	}

	currentChunkSize := chunkSize / sizeDivider
	if currentChunkSize < 1 {
		currentChunkSize = 1
	}

	return &inMemoryEventStorage{
		chunkSize:        chunkSize,
		currentChunkSize: currentChunkSize,
		sizeDivider:      sizeDivider,
		events:           make([]*types.Event, 0, currentChunkSize),
	}
}

type inMemoryEventStorage struct {
	chunkSize        int
	currentChunkSize int
	sizeDivider      int
	events           []*types.Event
	retriedEvents    []*types.Event

	mu sync.RWMutex
}

func (s *inMemoryEventStorage) PushNew(event *types.Event) {
	s.push(false, event)
}

func (s *inMemoryEventStorage) ReturnBack(events ...*types.Event) {
	s.push(true, events...)
}

func (s *inMemoryEventStorage) push(prepend bool, events ...*types.Event) {
	if len(events) == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	prependIndex := 0
	for _, event := range events {
		if event.RetryAt.IsZero() {
			s.addNonRetriedEvent(event, prepend, &prependIndex)
		} else {
			s.addRetriedEvent(event)
		}
	}
}

// PullChunk returns a chunk of events.
func (s *inMemoryEventStorage) PullChunk() []*types.Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.events) >= s.currentChunkSize {
		events := make([]*types.Event, s.currentChunkSize)
		copy(events, s.events)
		copy(s.events, s.events[s.currentChunkSize:])
		s.events = s.events[:len(s.events)-s.currentChunkSize]

		return events
	}

	events := make([]*types.Event, len(s.events), s.currentChunkSize)
	copy(events, s.events)
	s.events = s.events[:0]

	now := time.Now()
	retriedCount := 0
	for retriedCount < len(s.retriedEvents) && len(s.events)+retriedCount < s.chunkSize && s.retriedEvents[retriedCount].RetryAt.Before(now) {
		retriedCount++
	}
	if retriedCount == 0 {
		return events
	}
	events = append(events, s.retriedEvents[:retriedCount]...)
	copy(s.retriedEvents, s.retriedEvents[retriedCount:])
	s.retriedEvents = s.retriedEvents[:len(s.retriedEvents)-retriedCount]

	return events
}

func (s *inMemoryEventStorage) HasFullChunk() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.events) >= s.currentChunkSize
}

func (s *inMemoryEventStorage) ReduceChunkSize() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sizeDivider++
	s.currentChunkSize = s.chunkSize / s.sizeDivider
	if s.currentChunkSize < 1 {
		s.currentChunkSize = 1
	}
}

func (s *inMemoryEventStorage) addNonRetriedEvent(event *types.Event, prepend bool, prependIndex *int) {
	if prepend {
		s.events = append(s.events, nil)
		copy(s.events[*prependIndex+1:], s.events[*prependIndex:])
		s.events[*prependIndex] = event
		*prependIndex++
	} else {
		s.events = append(s.events, event)
	}
}

func (s *inMemoryEventStorage) addRetriedEvent(event *types.Event) {
	index := sort.Search(len(s.retriedEvents), func(i int) bool {
		return s.retriedEvents[i].RetryAt.After(event.RetryAt)
	})
	if index == len(s.retriedEvents) {
		s.retriedEvents = append(s.retriedEvents, event)
	} else {
		s.retriedEvents = append(s.retriedEvents, nil)
		copy(s.retriedEvents[index+1:], s.retriedEvents[index:])
		s.retriedEvents[index] = event
	}
}
