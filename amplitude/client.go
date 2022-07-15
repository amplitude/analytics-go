package amplitude

type Client interface {
	Track(event BaseEvent)
	Identify(identifyObj Identify, eventOptions EventOptions, eventProperties map[string]interface{})
	GroupIdentify(groupType, groupName string, identify Identify,
		eventOptions EventOptions, eventProperties, userProperties map[string]interface{},
	)
	Revenue(revenue Revenue, eventOptions EventOptions)
	SetGroup(groupType string, groupName string, eventOptions EventOptions)
	Flush()
	Add(plugin Plugin)
	Remove(plugin Plugin)
	Shutdown()
}

func NewClient(config Config) Client {
	return &client{configuration: config}
}

type client struct {
	configuration Config
	timeline      Timeline
}

// Track processes and sends the given event object.
func (a *client) Track(event BaseEvent) {
	a.timeline.Process(event)
}

// Identify sends an identify event to update user Properties.
func (a *client) Identify(identifyObj Identify, eventOptions EventOptions, eventProperties map[string]interface{}) {
	if !identifyObj.IsValid() {
		a.configuration.Logger.Error("Empty Identify Properties")
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
func (a *client) GroupIdentify(groupType, groupName string, identifyObj Identify,
	eventOptions EventOptions, eventProperties, userProperties map[string]interface{},
) {
	if !identifyObj.IsValid() {
		a.configuration.Logger.Error("Empty group identify Properties")
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
func (a *client) Revenue(revenue Revenue, eventOptions EventOptions) {
	if !revenue.IsValid() {
		a.configuration.Logger.Error("Invalid revenue quantity")
	} else {
		revenueEvent := revenue.ToRevenueEvent()
		revenueEvent.loadEventOptions(eventOptions)
		a.Track(revenueEvent.BaseEvent)
	}
}

// SetGroup sends an identify event to put a user in group(s)
// by setting group type and group name as user property for a user.
func (a *client) SetGroup(groupType string, groupName string, eventOptions EventOptions) {
	identifyObj := Identify{}
	identifyObj.Set(groupType, groupName)
	a.Identify(identifyObj, eventOptions, map[string]interface{}{})
}

// Flush flushes all events waiting to be sent in the buffer.
func (a *client) Flush() {
	a.timeline.Flush()
}

// Add adds the plugin object to client instance.
// Events tracked bby this client instance will be processed by instances' plugins.
func (a *client) Add(plugin Plugin) {
	a.timeline.Add(plugin)
	plugin.Setup(a)
}

// Remove removes the plugin object from client instance.
func (a *client) Remove(plugin Plugin) {
	a.timeline.Remove(plugin)
}

// Shutdown shuts the client instance down from accepting new events
// flushes all events in the buffer.
func (a *client) Shutdown() {
	a.configuration.OptOut = false
	a.timeline.Shutdown()
}
