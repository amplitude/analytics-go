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
