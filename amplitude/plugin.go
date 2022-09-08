package amplitude

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
