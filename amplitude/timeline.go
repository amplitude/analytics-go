package amplitude

type timeline struct {
	configuration      Config
	logger             Logger
	enrichmentPlugins  []EnrichmentPlugin
	destinationPlugins []DestinationPlugin
}

func (t *timeline) process(event *Event) {
	if t.configuration.OptOut {
		t.logger.Infof("Skipped event for opt out config: \n\t%+v", event)

		return
	}

	event = t.applyEnrichmentPlugins(event)
	if event != nil {
		t.applyDestinationPlugins(event)
	}
}

func (t *timeline) applyEnrichmentPlugins(event *Event) *Event {
	result := event

	for pluginType := BEFORE; pluginType <= ENRICHMENT; pluginType++ {
		for _, plugin := range t.enrichmentPlugins {
			if plugin.Type() == pluginType {
				result = plugin.Execute(result)
				if result == nil {
					return nil
				}
			}
		}
	}

	return result
}

func (t *timeline) applyDestinationPlugins(event *Event) {
	for _, plugin := range t.destinationPlugins {
		clone := event.Clone()
		plugin.Execute(&clone)
	}
}

func (t *timeline) add(plugin Plugin) {
	//	TO-DO stop current thread

	switch plugin := plugin.(type) {
	case EnrichmentPlugin:
		t.enrichmentPlugins = append(t.enrichmentPlugins, plugin)
	case DestinationPlugin:
		t.destinationPlugins = append(t.destinationPlugins, plugin)
	default:
		panic("unknown plugin type")
	}
}

func (t *timeline) remove(plugin Plugin) {
	switch plugin := plugin.(type) {
	case EnrichmentPlugin:
		for i, p := range t.enrichmentPlugins {
			if p == plugin {
				t.enrichmentPlugins = append(t.enrichmentPlugins[:i], t.enrichmentPlugins[i+1:]...)
			}
		}
	case DestinationPlugin:
		for i, p := range t.destinationPlugins {
			if p == plugin {
				t.destinationPlugins = append(t.destinationPlugins[:i], t.destinationPlugins[i+1:]...)
			}
		}
	default:
		panic("unknown plugin type")
	}
}

func (t *timeline) flush() {
	for _, plugin := range t.destinationPlugins {
		if plugin, ok := plugin.(ExtendedDestinationPlugin); ok {
			plugin.Flush()
		}
	}
}

func (t *timeline) shutdown() {
	for _, plugin := range t.destinationPlugins {
		if plugin, ok := plugin.(ExtendedDestinationPlugin); ok {
			plugin.Shutdown()
		}
	}
}
