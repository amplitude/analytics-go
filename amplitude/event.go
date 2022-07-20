package amplitude

import "time"

type EventOptions struct {
	time     time.Time
	insertID string
	library  string
}

type Event struct {
	EventType string
	EventOptions
	EventProperties map[string]interface{}
	UserProperties  map[string]interface{}
	Groups          map[string][]string
	GroupProperties map[string]interface{}
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
