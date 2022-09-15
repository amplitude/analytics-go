package types

type Event struct {
	EventType  string
	Properties map[string]interface{}
	UserID     string
}
