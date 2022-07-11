package amplitude

type Plugin struct {
	pluginType string
}

func (p Plugin) setup(client Amplitude) {

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
