package storages

import (
	"github.com/amplitude/analytics-go/amplitude/types"
)

func NewInMemoryEventStorage(capacity int) types.EventStorage {
	return &inMemoryEventStorage{
		capacity: capacity,
		events:   make([]*types.Event, 0, capacity),
	}
}

type inMemoryEventStorage struct {
	capacity int
	events   []*types.Event
}

// Push pushes an event to the storage.
func (i *inMemoryEventStorage) Push(event *types.Event) {
	i.events = append(i.events, event)
}

// Pull returns all Events and empties EventStorage.
func (i *inMemoryEventStorage) Pull() []*types.Event {
	events := i.events
	i.events = make([]*types.Event, 0, i.capacity)

	return events
}

func (i *inMemoryEventStorage) Len() int {
	return len(i.events)
}
