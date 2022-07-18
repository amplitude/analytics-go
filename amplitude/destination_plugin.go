package amplitude

// DestinationPlugin  is the base class to implement plugins that
// send events to customized destinations.
type DestinationPlugin struct {
	EventPlugin
}

// AmplitudeDestinationPlugin is the default destination plugin that
// sends events to Amplitude.
type AmplitudeDestinationPlugin struct {
	DestinationPlugin
}
