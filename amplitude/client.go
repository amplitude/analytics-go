package amplitude

type Amplitude struct {
	Configuration Config
	timeline      Timeline
}

// Track processes and sends the given event object
func (a Amplitude) Track(event BaseEvent) {
	a.timeline.process(event)
}

// Identify sends an identify event to update user properties
func (a Amplitude) Identify(identityObj Identity, eventOptions EventOptions, eventProperties ...map[string]string) {
	if !identityObj.isValid() {
		a.Configuration.Logger.Println("Error: Empty identify properties")
	} else {
		var event = IdentifyEvent{}
		event.loadEventOptions(eventOptions)
		a.Track(event.BaseEvent)
	}

}

// GroupIdentify sends a group identify event to update group properties
func (a Amplitude) GroupIdentify(groupType string, groupName string, identifyObj Identity, EventOptions EventOptions, eventProperties map[string]string) {
	if !identifyObj.isValid() {
		a.Configuration.Logger.Println("Error: Empty identify properties")
	} else {
		var event = GroupIdentifyEvent{}
		event.loadEventOptions(EventOptions)
		a.Track(event.BaseEvent)
	}
}

// Revenue sends a revenue event with revenue info in eventProperties
func (a Amplitude) Revenue(revenueObj Revenue, eventOptions EventOptions) {
	if !revenueObj.isValid() {
		a.Configuration.Logger.Println("Error: Empty identify properties")
	} else {
		var event = revenueObj.toRevenueEvent()
		event.loadEventOptions(eventOptions)
		a.Track(event.BaseEvent)
	}
}

// SetGroup sends an identify event to put a user in group(s)
// by setting group type and group name as user property for a user
func (a Amplitude) SetGroup(groupType string, groupName string, eventOptions EventOptions) {
	var identifyObj = Identity{}
	identifyObj.set(groupType, groupName)
	a.Identify(identifyObj, eventOptions)
}

// Flush flushes all event waiting to be sent in the buffer
func (a Amplitude) Flush() {
	a.timeline.flush()
}

// Add adds the plugin object to client instance.
// Events tracked by this client instance will be processed by instance's plugins.
func (a Amplitude) Add(plugin Plugin) {
	a.timeline.add(plugin)
	plugin.setup(a)
}

// Remove removes the plugin object from client instance
func (a Amplitude) Remove(plugin Plugin) {
	a.timeline.remove(plugin)
}

// Shutdown shuts the client instance down from accepting new events
// flushes all events in the buffer
func (a Amplitude) Shutdown() {
	a.Configuration.OptOut = true
	a.timeline.shutdown()
}
