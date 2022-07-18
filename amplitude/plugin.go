package amplitude

type Plugin interface {
	Setup(client *client)
	Execute(event Event) Event
	flush()
	shutdown()
}

// BasePlugin is the base class of all plugins.
type BasePlugin struct {
	pluginType PluginType
}

func (p *BasePlugin) Setup(client *client) {

}

func (p *BasePlugin) Execute(event Event) Event {
	return Event{}
}

func (p *BasePlugin) flush() {

}

func (p *BasePlugin) shutdown() {

}
