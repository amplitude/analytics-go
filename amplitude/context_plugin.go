package amplitude

import (
	"time"

	"github.com/google/uuid"
)

// ContextPlugin is the default plugin that add library info to event.
// It also sets event default timestamp and insertID if not set elsewhere.
type ContextPlugin struct {
	contextString string
}

func NewContextPlugin() *ContextPlugin {
	return &ContextPlugin{
		contextString: SdkLibrary + "/" + SdkVersion,
	}
}

func (c *ContextPlugin) Setup(config Config) {
}

func (c *ContextPlugin) Priority() MiddlewarePriority {
	return MiddlewarePriorityBefore
}

// Execute sets default timestamp and insertID if not set elsewhere
// It also adds SDK name and version to event library.
func (c *ContextPlugin) Execute(event *Event) *Event {
	if event.time.IsZero() {
		event.time = time.Now()
	}

	if event.insertID == "" {
		event.insertID = uuid.NewString()
	}

	event.library = c.contextString

	return event
}
