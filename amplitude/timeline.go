package amplitude

type timeline struct {
	configuration Config
	logger        Logger
	plugins       map[PluginType][]Plugin
}

func (t *timeline) process(event Event) Event {
	if t.configuration.OptOut {
		t.logger.Info("Skipped event for opt out config")

		return event
	}

	beforeResult := t.applyPlugins(BEFORE, event)
	enrichResult := t.applyPlugins(ENRICHMENT, beforeResult)
	t.applyPlugins(DESTINATION, enrichResult)

	return enrichResult
}

func (t *timeline) applyPlugins(pluginType PluginType, event Event) Event {
	result := event

	for _, plugin := range t.plugins[pluginType] {
		result = plugin.Execute(result)
	}

	return result
}

func (t *timeline) add(pluginType PluginType, plugin Plugin) {
	//	TO-DO stop current thread
	t.plugins[pluginType] = append(t.plugins[pluginType], plugin)
}

func (t *timeline) remove(plugin Plugin) {
	for pluginsType, plugins := range t.plugins {
		for i, p := range plugins {
			if p == plugin {
				t.plugins[pluginsType] = append(t.plugins[pluginsType][:i], t.plugins[pluginsType][i+1:]...)
			}
		}
	}
}

func (t *timeline) flush() {
	for _, destinationPlugin := range t.plugins[DESTINATION] {
		destinationPlugin.flush()
	}
}

func (t *timeline) shutdown() {
	for _, destinationPlugin := range t.plugins[DESTINATION] {
		destinationPlugin.shutdown()
	}
}
