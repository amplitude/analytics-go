package amplitude

type EventOptions struct{}

type Event struct {
	EventType string
	EventOptions
	EventProperties map[string]interface{}
	UserProperties  map[string]interface{}
	Groups          map[string][]string
	GroupProperties map[string]interface{}
}
