package types

type EventStorage interface {
	Push(event *EventPayload)
	Pull() []*EventPayload
	Len() int
}
