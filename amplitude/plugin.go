package amplitude

type Plugin struct {
}

func (p Plugin) Setup(client *client) {

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
