package amplitude

type EventOptions struct{}

type BaseEvent struct {
	EventOptions
	EventProperties map[string]interface{}
	UserProperties  map[string]interface{}
	Groups          map[string]string
	GroupProperties map[string]interface{}
}

func (b BaseEvent) loadEventOptions(options EventOptions) {
}

type GroupIdentifyEvent struct {
	BaseEvent
}

type IdentifyEvent struct {
	BaseEvent
}

type RevenueEvent struct {
	BaseEvent
}
