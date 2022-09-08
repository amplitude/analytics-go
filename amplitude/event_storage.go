package amplitude

type EventStorage interface {
	Push(event *Event)
	Pull() []*Event
	Len() int
}

func NewInMemoryEventStorage(capacity int) *InMemoryEventStorage {
	return &InMemoryEventStorage{
		capacity:     capacity,
		eventsBuffer: make([]*Event, 0, capacity),
	}
}

type InMemoryEventStorage struct {
	capacity     int
	eventsBuffer []*Event
}

// Push pushes an event to the storage.
func (i *InMemoryEventStorage) Push(event *Event) {
	i.eventsBuffer = append(i.eventsBuffer, event)
}

// Pull returns all Events and empties InMemoryEventStorage.
func (i *InMemoryEventStorage) Pull() []*Event {
	events := i.eventsBuffer
	i.eventsBuffer = make([]*Event, 0, i.capacity)

	return events
}

func (i *InMemoryEventStorage) Len() int {
	return len(i.eventsBuffer)
}
