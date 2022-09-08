package amplitude

import (
	"github.com/amplitude/analytics-go/amplitude/internal"
)

type Client interface {
	Track(event Event)
	Identify(identify Identify, eventOptions EventOptions)
	GroupIdentify(groupType string, groupName string, identify Identify, eventOptions EventOptions)
	SetGroup(groupType string, groupName []string, eventOptions EventOptions)
	Revenue(revenue Revenue, eventOptions EventOptions)

	Flush()
	Shutdown()

	AddPlugin(plugin Plugin)
	RemovePlugin(pluginName string)

	Config() Config
}

func NewClient(config Config) Client {
	config = config.setDefaultValues()
	config.Logger.Debugf("Client initialized")

	client := &client{
		config:   config,
		optOut:   internal.NewAtomicBool(config.OptOut),
		timeline: &timeline{},
	}

	client.AddPlugin(&AmplitudePlugin{})
	client.AddPlugin(NewContextPlugin())

	return client
}

type client struct {
	config   Config
	timeline *timeline
	optOut   *internal.AtomicBool
}

func (c *client) Config() Config {
	return c.config
}

// Track processes and sends the given event object.
func (c *client) Track(event Event) {
	if !c.enabled() {
		return
	}

	if event.Plan == (Plan{}) {
		event.Plan = c.config.Plan
	}

	c.config.Logger.Debugf("Track event: \n\t%+v", event)
	c.timeline.Process(&event)
}

// Identify sends an identify event to update user Properties.
func (c *client) Identify(identify Identify, eventOptions EventOptions) {
	if !c.enabled() {
		return
	}

	validateErrors, validateWarnings := identify.Validate()

	for _, validateWarning := range validateWarnings {
		c.config.Logger.Warnf("Identify: %s", validateWarning)
	}

	if len(validateErrors) > 0 {
		for _, validateError := range validateErrors {
			c.config.Logger.Errorf("Identify: %s", validateError)
		}
	} else {
		identifyEvent := Event{
			EventType:      IdentifyEventType,
			EventOptions:   eventOptions,
			UserProperties: identify.Properties,
		}

		c.Track(identifyEvent)
	}
}

// GroupIdentify sends a group identify event to update group Properties.
func (c *client) GroupIdentify(groupType string, groupName string, identify Identify, eventOptions EventOptions) {
	if !c.enabled() {
		return
	}

	validateErrors, validateWarnings := identify.Validate()

	for _, validateWarning := range validateWarnings {
		c.config.Logger.Warnf("Identify: %s", validateWarning)
	}

	if len(validateErrors) > 0 {
		for _, validateError := range validateErrors {
			c.config.Logger.Errorf("Invalid Identify: %s", validateError)
		}
	} else {
		groupIdentifyEvent := Event{
			EventType:       GroupIdentifyEventType,
			EventOptions:    eventOptions,
			Groups:          map[string][]string{groupType: {groupName}},
			GroupProperties: identify.Properties,
		}

		c.Track(groupIdentifyEvent)
	}
}

// Revenue sends a revenue event with revenue info in eventProperties.
func (c *client) Revenue(revenue Revenue, eventOptions EventOptions) {
	if !c.enabled() {
		return
	}

	if validateErrors := revenue.Validate(); len(validateErrors) > 0 {
		for _, validateError := range validateErrors {
			c.config.Logger.Errorf("Invalid Revenue: %s", validateError)
		}
	} else {
		revenueEvent := Event{
			EventType:    RevenueEventType,
			EventOptions: eventOptions,
			EventProperties: map[string]interface{}{
				RevenueProductID:  revenue.ProductID,
				RevenueQuantity:   revenue.Quantity,
				RevenuePrice:      revenue.Price,
				RevenueType:       revenue.RevenueType,
				RevenueReceipt:    revenue.Receipt,
				RevenueReceiptSig: revenue.ReceiptSig,
				DefaultRevenue:    revenue.Revenue,
			},
		}
		c.Track(revenueEvent)
	}
}

// SetGroup sends an identify event to put a user in group(s)
// by setting group type and group name as user property for a user.
func (c *client) SetGroup(groupType string, groupName []string, eventOptions EventOptions) {
	if !c.enabled() {
		return
	}

	identify := Identify{}
	identify.Set(groupType, groupName)
	c.Identify(identify, eventOptions)
}

// Flush flushes all events waiting to be sent in the buffer.
func (c *client) Flush() {
	c.timeline.Flush()
}

// AddPlugin adds the plugin object to client instance.
// Events tracked by this client instance will be processed by instances' plugins.
func (c *client) AddPlugin(plugin Plugin) {
	c.timeline.AddPlugin(plugin)
	plugin.Setup(c.config)
}

// RemovePlugin removes the plugin object from client instance.
func (c *client) RemovePlugin(pluginName string) {
	c.timeline.RemovePlugin(pluginName)
}

// Shutdown shuts the client instance down from accepting new events.
func (c *client) Shutdown() {
	c.optOut.Set()

	c.config.Logger.Debugf("Client shutdown")
	c.timeline.Shutdown()
}

func (c *client) enabled() bool {
	return !c.optOut.IsSet()
}
