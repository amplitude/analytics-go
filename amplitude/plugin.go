package amplitude

type EnrichmentPriority byte

const (
	EnrichmentPriorityBefore EnrichmentPriority = iota
	EnrichmentPriorityEnrichment
)

type Plugin interface {
	Setup(config Config)
}

type EnrichmentPlugin interface {
	Plugin
	Priority() EnrichmentPriority
	Execute(event *Event) *Event
}

type DestinationPlugin interface {
	Plugin
	Execute(event *Event)
	Flush()
	Shutdown()
}
