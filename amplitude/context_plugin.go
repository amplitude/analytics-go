package amplitude

// ContextPlugin is the default plugin that add library info to event.
// It also sets event default timestamp and insert_id if not set elsewhere
type ContextPlugin struct {
	BasePlugin
}
