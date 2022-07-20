package amplitude

type MiddlewarePriority byte

const (
	MiddlewarePriorityBefore MiddlewarePriority = iota
	MiddlewarePriorityEnrichment
)

type Plugin interface {
	Setup(config Config)
}

type MiddlewarePlugin interface {
	Plugin
	Priority() MiddlewarePriority
	Execute(event *Event) *Event
}

type DestinationPlugin interface {
	Plugin
	Execute(event *Event)
	Flush()
	Shutdown()
}
