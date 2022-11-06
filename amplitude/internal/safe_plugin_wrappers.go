package internal

import (
	"github.com/amplitude/analytics-go/amplitude/types"
)

type SafeBeforePluginWrapper struct {
	Plugin types.BeforePlugin
	Logger types.Logger
}

func (w *SafeBeforePluginWrapper) Name() string {
	return w.Plugin.Name()
}

func (w *SafeBeforePluginWrapper) Type() types.PluginType {
	return w.Plugin.Type()
}

func (w *SafeBeforePluginWrapper) Setup(config types.Config) {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Setup: %s", w.Plugin.Name(), r)
		}
	}()

	w.Plugin.Setup(config)
}

func (w *SafeBeforePluginWrapper) Execute(event *types.Event) (result *types.Event) {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Execute: %s", w.Plugin.Name(), r)

			result = event
		}
	}()

	return w.Plugin.Execute(event)
}

type SafeEnrichmentPluginWrapper struct {
	Plugin types.EnrichmentPlugin
	Logger types.Logger
}

func (w *SafeEnrichmentPluginWrapper) Name() string {
	return w.Plugin.Name()
}

func (w *SafeEnrichmentPluginWrapper) Type() types.PluginType {
	return w.Plugin.Type()
}

func (w *SafeEnrichmentPluginWrapper) Setup(config types.Config) {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Setup: %s", w.Plugin.Name(), r)
		}
	}()

	w.Plugin.Setup(config)
}

func (w *SafeEnrichmentPluginWrapper) Execute(event *types.Event) (result *types.Event) {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Execute: %s", w.Plugin.Name(), r)

			result = event
		}
	}()

	return w.Plugin.Execute(event)
}

type SafeDestinationPluginWrapper struct {
	Plugin types.DestinationPlugin
	Logger types.Logger
}

func (w *SafeDestinationPluginWrapper) Name() string {
	return w.Plugin.Name()
}

func (w *SafeDestinationPluginWrapper) Type() types.PluginType {
	return w.Plugin.Type()
}

func (w *SafeDestinationPluginWrapper) Setup(config types.Config) {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Setup: %s", w.Plugin.Name(), r)
		}
	}()

	w.Plugin.Setup(config)
}

func (w *SafeDestinationPluginWrapper) Execute(event *types.Event) {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Execute: %s", w.Plugin.Name(), r)
		}
	}()

	w.Plugin.Execute(event)
}

type SafeExtendedDestinationPluginWrapper struct {
	Plugin types.ExtendedDestinationPlugin
	Logger types.Logger
}

func (w *SafeExtendedDestinationPluginWrapper) Name() string {
	return w.Plugin.Name()
}

func (w *SafeExtendedDestinationPluginWrapper) Type() types.PluginType {
	return w.Plugin.Type()
}

func (w *SafeExtendedDestinationPluginWrapper) Setup(config types.Config) {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Setup: %s", w.Plugin.Name(), r)
		}
	}()

	w.Plugin.Setup(config)
}

func (w *SafeExtendedDestinationPluginWrapper) Execute(event *types.Event) {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Execute: %s", w.Plugin.Name(), r)
		}
	}()

	w.Plugin.Execute(event)
}

func (w *SafeExtendedDestinationPluginWrapper) Flush() {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Flush: %s", w.Plugin.Name(), r)
		}
	}()

	w.Plugin.Flush()
}

func (w *SafeExtendedDestinationPluginWrapper) Shutdown() {
	defer func() {
		if r := recover(); r != nil {
			w.Logger.Errorf("Panic in plugin %s.Shutdown: %s", w.Plugin.Name(), r)
		}
	}()

	w.Plugin.Shutdown()
}
