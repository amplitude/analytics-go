package amplitude

type Timeline struct {
	Configuration Config
}

func (t Timeline) setup(client Amplitude) {
	t.Configuration = client.Configuration
}

func (t Timeline) add(plugin Plugin) {

}

func (t Timeline) remove(plugin Plugin) {

}

func (t Timeline) flush() {

}

func (t Timeline) process(event BaseEvent) {

}

func (t Timeline) applyPlugins(pluginType PluginType, event BaseEvent) {

}

func (t Timeline) shutdown() {

}