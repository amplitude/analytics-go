package amplitude

type timeline struct {
	configuration      Config
	logger             Logger
	middlewarePlugins  []MiddlewarePlugin
	destinationPlugins []DestinationPlugin
}

func (t *timeline) process(event *Event) {
	if t.configuration.OptOut {
		t.logger.Info("Skipped event for opt out config")

		return
	}

	event = t.applyMiddlewarePlugins(event)
	if event != nil {
		t.applyDestinationPlugins(event)
	}
}

func (t *timeline) applyMiddlewarePlugins(event *Event) *Event {
	result := event

	for priority := MiddlewarePriorityBefore; priority <= MiddlewarePriorityEnrichment; priority++ {
		for _, plugin := range t.middlewarePlugins {
			if plugin.Priority() == priority {
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
	case MiddlewarePlugin:
		t.middlewarePlugins = append(t.middlewarePlugins, plugin)
	case DestinationPlugin:
		t.destinationPlugins = append(t.destinationPlugins, plugin)
	default:
		panic("unknown plugin type")
	}
}

func (t *timeline) remove(plugin Plugin) {
	switch plugin := plugin.(type) {
	case MiddlewarePlugin:
		for i, p := range t.middlewarePlugins {
			if p == plugin {
				t.middlewarePlugins = append(t.middlewarePlugins[:i], t.middlewarePlugins[i+1:]...)
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
		plugin.Flush()
	}
}

func (t *timeline) shutdown() {
	for _, plugin := range t.destinationPlugins {
		plugin.Shutdown()
	}
}
