package amplitude

type EventOptions struct {
	UserId   string `json:"user_id"`
	DeviceId string `json:"device_id"`
	Time     int64  `json:"time,omitempty"`
	InsertID string `json:"insert_id,omitempty"`
	Library  string `json:"library,omitempty"`
}

type Event struct {
	EventType string `json:"event_type"`
	EventOptions
	EventProperties map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties  map[string]interface{} `json:"user_properties,omitempty"`
	Groups          map[string][]string    `json:"groups,omitempty"`
	GroupProperties map[string]interface{} `json:"group_properties,omitempty"`
}

func (e Event) Clone() Event {
	return Event{
		EventType:       e.EventType,
		EventOptions:    e.EventOptions,
		EventProperties: cloneProperties(e.EventProperties),
		UserProperties:  cloneProperties(e.UserProperties),
		Groups:          cloneGroups(e.Groups),
		GroupProperties: cloneProperties(e.GroupProperties),
	}
}

// TODO: deep copy
func cloneProperties(properties map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(properties))
	for k, v := range properties {
		result[k] = v
	}

	return result
}

// TODO: deep copy
func cloneGroups(properties map[string][]string) map[string][]string {
	result := make(map[string][]string, len(properties))
	for k, v := range properties {
		result[k] = v
	}

	return result
}
