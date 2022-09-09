package types

type EventStorage interface {
	Push(event *Event)
	Pull() []*Event
	Len() int
}
