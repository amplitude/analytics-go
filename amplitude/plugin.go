package amplitude

type Plugin struct {
}

type EventPlugin struct {
	Plugin
}

type DestinationPlugin struct {
	EventPlugin
}

type AmplitudeDestinationPlugin struct {
	DestinationPlugin
}

type ContextPlugin struct {
	Plugin
}
