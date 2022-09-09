package types

type (
	PluginType int
)

const (
	PluginTypeBefore PluginType = iota
	PluginTypeEnrichment
	PluginTypeDestination
	PluginTypeObserve
)

type Plugin interface {
	Name() string
	Type() PluginType
	Setup(config Config)
}

type EnrichmentPlugin interface {
	Plugin
	Execute(event *Event) *Event
}

type DestinationPlugin interface {
	Plugin
	Execute(event *Event)
}

type ExtendedDestinationPlugin interface {
	DestinationPlugin
	Flush()
	Shutdown()
}

type ExecuteResult struct {
	PluginName string
	Event      *Event
	Code       int
	Message    string
}
