package amplitude

import (
	"time"

	"github.com/google/uuid"
)

// ContextPlugin is the default enrichment plugin that add library info to event.
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

func (c *ContextPlugin) Type() PluginType {
	return BEFORE
}

// Execute sets default timestamp and insertID if not set elsewhere
// It also adds SDK name and version to event library.
func (c *ContextPlugin) Execute(event *Event) *Event {
	if event.Time == 0 {
		event.Time = time.Now().UnixMilli()
	}

	if event.InsertID == "" {
		event.InsertID = uuid.NewString()
	}

	event.Library = c.contextString

	return event
}
