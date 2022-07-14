package amplitude

type Amplitude struct {
	Configuration Config
	timeline      Timeline
}

// Track processes and sends the given event object.
func (a Amplitude) Track(event BaseEvent) {
	a.timeline.Process(event)
}

// Identify sends an identify event to update user Properties.
func (a Amplitude) Identify(identifyObj Identify, eventOptions EventOptions, eventProperties map[string]interface{}) {
	if !identifyObj.IsValid() {
		a.Configuration.Logger.Error("Empty Identify Properties")
	} else {
		identifyEvent := IdentifyEvent{
			BaseEvent{
				EventProperties: eventProperties,
				UserProperties:  identifyObj.Properties,
			},
		}

		identifyEvent.loadEventOptions(eventOptions)
		a.Track(identifyEvent.BaseEvent)
	}
}

// GroupIdentify sends a group identify event to update group Properties.
func (a Amplitude) GroupIdentify(groupType string, groupName string, identifyObj Identify,
	eventOptions EventOptions, eventProperties map[string]interface{}, userProperties map[string]interface{},
) {
	if !identifyObj.IsValid() {
		a.Configuration.Logger.Error("Empty group identify Properties")
	} else {
		groupIdentifyEvent := GroupIdentifyEvent{
			BaseEvent{
				EventProperties: eventProperties,
				UserProperties:  userProperties,
				Groups:          map[string]string{groupType: groupName},
				GroupProperties: identifyObj.Properties,
			},
		}

		groupIdentifyEvent.loadEventOptions(eventOptions)
		a.Track(groupIdentifyEvent.BaseEvent)
	}
}

// Revenue sends a revenue event with revenue info in eventProperties.
func (a Amplitude) Revenue(revenueObj Revenue, eventOptions EventOptions) {
	if !revenueObj.IsValid() {
		a.Configuration.Logger.Error("Invalid revenue quantity")
	} else {
		revenueEvent := revenueObj.ToRevenueEvent()
		revenueEvent.loadEventOptions(eventOptions)
		a.Track(revenueEvent.BaseEvent)
	}
}

// SetGroup sends an identify event to put a user in group(s)
// by setting group type and group name as user property for a user.
func (a Amplitude) SetGroup(groupType string, groupName string, eventOptions EventOptions) {
	identifyObj := Identify{}
	identifyObj.Set(groupType, groupName)
	a.Identify(identifyObj, eventOptions, map[string]interface{}{})
}

// Flush flushes all events waiting to be sent in the buffer.
func (a Amplitude) Flush() {
	a.timeline.Flush()
}

// Add adds the plugin object to client instance.
// Events tracked bby this client instance will be processed by instances' plugins.
func (a *Amplitude) Add(plugin Plugin) *Amplitude {
	a.timeline.Add(plugin)
	plugin.Setup(a)

	return a
}

// Remove removes the plugin object from client instance.
func (a *Amplitude) Remove(plugin Plugin) *Amplitude {
	a.timeline.Remove(plugin)

	return a
}

// Shutdown shuts the client instance down from accepting new events
// flushes all events in the buffer.
func (a *Amplitude) Shutdown() {
	a.Configuration.OptOut = false
	a.timeline.Shutdown()
}
