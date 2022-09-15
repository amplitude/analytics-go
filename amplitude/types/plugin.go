package types

type (
	PluginType int
)

const (
	PluginTypeBefore PluginType = iota
	PluginTypeEnrichment
	PluginTypeDestination
)

type Plugin interface {
	Name() string
	Type() PluginType
	Setup(config Config)
}

type BeforePlugin interface {
	Plugin
	Execute(event *EventPayload) *EventPayload
}

type EnrichmentPlugin interface {
	Plugin
	Execute(event *EventPayload) *EventPayload
}

type DestinationPlugin interface {
	Plugin
	Execute(event *EventPayload)
}

type ExtendedDestinationPlugin interface {
	DestinationPlugin
	Flush()
	Shutdown()
}

type ExecuteResult struct {
	PluginName string
	Event      *EventPayload
	Code       int
	Message    string
}
