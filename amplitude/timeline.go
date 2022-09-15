package amplitude

import (
	"sync"
)

type timeline struct {
	logger             Logger
	beforePlugins      []BeforePlugin
	enrichmentPlugins  []EnrichmentPlugin
	destinationPlugins []DestinationPlugin
	mu                 sync.RWMutex
}

func (t *timeline) Process(event *EventPayload) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	event = t.applyBeforePlugins(event)
	if event == nil {
		return
	}

	event = t.applyEnrichmentPlugins(event)
	if event == nil {
		return
	}

	t.applyDestinationPlugins(event)
}

func (t *timeline) applyBeforePlugins(event *EventPayload) *EventPayload {
	result := event

	for _, plugin := range t.beforePlugins {
		result = plugin.Execute(result)
		if result == nil {
			return nil
		}
	}

	return result
}

func (t *timeline) applyEnrichmentPlugins(event *EventPayload) *EventPayload {
	result := event

	for _, plugin := range t.enrichmentPlugins {
		result = plugin.Execute(result)
		if result == nil {
			return nil
		}
	}

	return result
}

func (t *timeline) applyDestinationPlugins(event *EventPayload) {
	var wg sync.WaitGroup
	for _, plugin := range t.destinationPlugins {
		clone := event.Clone()
		wg.Add(1)
		go func(plugin DestinationPlugin, event *EventPayload) {
			defer wg.Done()
			plugin.Execute(event)
		}(plugin, &clone)
	}
	wg.Wait()
}

func (t *timeline) AddPlugin(plugin Plugin) {
	t.mu.Lock()
	defer t.mu.Unlock()

	switch plugin.Type() {
	case PluginTypeBefore:
		plugin, ok := plugin.(BeforePlugin)
		if !ok {
			t.logger.Errorf("Plugin %s doesn't implement Before interface", plugin.Name())
		}
		t.beforePlugins = append(t.beforePlugins, plugin)
	case PluginTypeEnrichment:
		plugin, ok := plugin.(EnrichmentPlugin)
		if !ok {
			t.logger.Errorf("Plugin %s doesn't implement Enrichment interface", plugin.Name())
		}
		t.enrichmentPlugins = append(t.enrichmentPlugins, plugin)
	case PluginTypeDestination:
		plugin, ok := plugin.(DestinationPlugin)
		if !ok {
			t.logger.Errorf("Plugin %s doesn't implement Destination interface", plugin.Name())
		}
		t.destinationPlugins = append(t.destinationPlugins, plugin)
	default:
		t.logger.Errorf("Plugin %s - unknown type %s", plugin.Name(), plugin.Type())
	}
}

func (t *timeline) RemovePlugin(pluginName string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i := len(t.beforePlugins) - 1; i >= 0; i-- {
		if t.beforePlugins[i].Name() == pluginName {
			t.beforePlugins = append(t.beforePlugins[:i], t.beforePlugins[i+1:]...)
		}
	}
	for i := len(t.enrichmentPlugins) - 1; i >= 0; i-- {
		if t.enrichmentPlugins[i].Name() == pluginName {
			t.enrichmentPlugins = append(t.enrichmentPlugins[:i], t.enrichmentPlugins[i+1:]...)
		}
	}
	for i := len(t.destinationPlugins) - 1; i >= 0; i-- {
		if t.destinationPlugins[i].Name() == pluginName {
			t.destinationPlugins = append(t.destinationPlugins[:i], t.destinationPlugins[i+1:]...)
		}
	}
}

func (t *timeline) Flush() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var wg sync.WaitGroup
	for _, plugin := range t.destinationPlugins {
		if plugin, ok := plugin.(ExtendedDestinationPlugin); ok {
			wg.Add(1)
			go func(plugin ExtendedDestinationPlugin) {
				defer wg.Done()
				plugin.Flush()
			}(plugin)
		}
	}
	wg.Wait()
}

func (t *timeline) Shutdown() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var wg sync.WaitGroup
	for _, plugin := range t.destinationPlugins {
		if plugin, ok := plugin.(ExtendedDestinationPlugin); ok {
			wg.Add(1)
			go func(plugin ExtendedDestinationPlugin) {
				defer wg.Done()
				plugin.Shutdown()
			}(plugin)
		}
	}
	wg.Wait()
}
