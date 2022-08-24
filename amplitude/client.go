package amplitude

type Client interface {
	Track(event Event)
	Identify(identify Identify, eventOptions EventOptions)
	GroupIdentify(groupType string, groupName []string, identify Identify,
		eventOptions EventOptions)
	Revenue(revenue Revenue, eventOptions EventOptions)
	SetGroup(groupType string, groupName []string, eventOptions EventOptions)
	Flush()
	Add(plugin Plugin)
	Remove(plugin Plugin)
	Shutdown()
}

func NewClient(config Config) Client {
	config.Logger.Debugf("Client initialized")

	client := &client{configuration: config}

	client.Add(&AmplitudePlugin{})
	client.Add(NewContextPlugin())

	return client
}

type client struct {
	configuration Config
	timeline      timeline
}

// Track processes and sends the given event object.
func (a *client) Track(event Event) {
	a.configuration.Logger.Debugf("Track event: \n\t%+v", event)
	a.timeline.process(&event)
}

// Identify sends an identify event to update user Properties.
func (a *client) Identify(identify Identify, eventOptions EventOptions) {
	if !identify.IsValid() {
		a.configuration.Logger.Errorf("Empty Identify Properties: \n\t%+v", identify)
	} else {
		identifyEvent := Event{
			EventType:      IdentifyEventType,
			EventOptions:   eventOptions,
			UserProperties: identify.Properties,
		}

		a.Track(identifyEvent)
	}
}

// GroupIdentify sends a group identify event to update group Properties.
func (a *client) GroupIdentify(groupType string, groupName []string, identify Identify,
	eventOptions EventOptions,
) {
	if !identify.IsValid() {
		a.configuration.Logger.Errorf("Empty group identify Properties: \n\t%+v", identify)
	} else {
		groupIdentifyEvent := Event{
			EventType:       GroupIdentifyEventType,
			EventOptions:    eventOptions,
			Groups:          map[string][]string{groupType: groupName},
			GroupProperties: identify.Properties,
		}

		a.Track(groupIdentifyEvent)
	}
}

// Revenue sends a revenue event with revenue info in eventProperties.
func (a *client) Revenue(revenue Revenue, eventOptions EventOptions) {
	if !revenue.IsValid() {
		a.configuration.Logger.Errorf("Either Revenue or Price should be set. Invalid Revenue object: \n\t%+v", revenue)
	} else {
		revenueEvent := revenue.ToRevenueEvent(eventOptions)
		a.Track(revenueEvent)
	}
}

// SetGroup sends an identify event to put a user in group(s)
// by setting group type and group name as user property for a user.
func (a *client) SetGroup(groupType string, groupName []string, eventOptions EventOptions) {
	identify := Identify{}
	identify.Set(groupType, groupName)
	a.Identify(identify, eventOptions)
}

// Flush flushes all events waiting to be sent in the buffer.
func (a *client) Flush() {
	a.timeline.flush()
}

// Add adds the plugin object to client instance.
// Events tracked bby this client instance will be processed by instances' plugins.
func (a *client) Add(plugin Plugin) {
	a.timeline.add(plugin)
	plugin.Setup(a.configuration)
}

// Remove removes the plugin object from client instance.
func (a *client) Remove(plugin Plugin) {
	a.timeline.remove(plugin)
}

// Shutdown shuts the client instance down from accepting new events
// flushes all events in the buffer.
func (a *client) Shutdown() {
	a.configuration.Logger.Debugf("Client shutdown")
	a.configuration.OptOut = true
	a.timeline.shutdown()
}
