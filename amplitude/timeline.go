package amplitude

import (
	"sync"
)

type timeline struct {
	beforePlugins      []BeforePlugin
	destinationPlugins []DestinationPlugin
	mu                 sync.RWMutex
}

func (t *timeline) Process(event *Event) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	event = t.applyBeforePlugins(event)
	if event != nil {
		t.applyDestinationPlugins(event)
	}
}

func (t *timeline) applyBeforePlugins(event *Event) *Event {
	result := event

	for _, plugin := range t.beforePlugins {
		result = plugin.Execute(result)
		if result == nil {
			return nil
		}
	}

	return result
}

func (t *timeline) applyDestinationPlugins(event *Event) {
	var wg sync.WaitGroup
	for _, plugin := range t.destinationPlugins {
		clone := event.Clone()
		wg.Add(1)
		go func(plugin DestinationPlugin, event *Event) {
			defer wg.Done()
			plugin.Execute(event)
		}(plugin, &clone)
	}
	wg.Wait()
}

func (t *timeline) AddPlugin(plugin Plugin) {
	t.mu.Lock()
	defer t.mu.Unlock()

	switch plugin := plugin.(type) {
	case BeforePlugin:
		t.beforePlugins = append(t.beforePlugins, plugin)
	case DestinationPlugin:
		t.destinationPlugins = append(t.destinationPlugins, plugin)
	default:
		panic("unknown plugin type")
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
