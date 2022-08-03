package amplitude

type Plugin interface {
	Setup(config Config)
	Type() PluginType
}

type EventPlugin interface {
	Plugin
	Execute(event *Event) *Event
}

type DestinationPlugin interface {
	Plugin
	Execute(event *Event)
	Flush()
	Shutdown()
}
