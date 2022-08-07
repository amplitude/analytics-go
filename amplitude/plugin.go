package amplitude

type Plugin interface {
	Setup(config Config)
	Type() PluginType
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
