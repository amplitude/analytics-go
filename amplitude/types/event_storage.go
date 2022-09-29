package types

type EventStorage interface {
	PushNew(event *Event)
	ReturnBack(events ...*Event)
	PullChunk() []*Event
	HasFullChunk() bool
	ReduceChunkSize()
}
