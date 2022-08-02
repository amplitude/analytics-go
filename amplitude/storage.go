package amplitude

type Storage interface {
	Push(event *Event)
	Pull() []*Event
}

type InMemoryStorage struct {
	eventsBuffer []*Event
}

// Push pushes an event to default InMemoryStorage.
func (i *InMemoryStorage) Push(event *Event) {
	i.eventsBuffer = append(i.eventsBuffer, event)
}

// Pull returns all events in default InMemoryStorage and empties InMemoryStorage.
func (i *InMemoryStorage) Pull() []*Event {
	events := i.eventsBuffer
	i.eventsBuffer = []*Event{}

	return events
}
