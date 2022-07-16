package amplitude

type Revenue struct {
	Price    float64
	Quantity int
}

// IsValid checks if a revenue instance has a positive integer quantity.
func (r Revenue) IsValid() bool {
	return r.Quantity > 0
}

func (r Revenue) ToRevenueEvent(eventOptions EventOptions) Event {
	return Event{
		EventType:       "$revenue",
		EventOptions:    eventOptions,
		EventProperties: r.GetEventProperties(),
	}
}

func (r Revenue) GetEventProperties() map[string]interface{} {
	return map[string]interface{}{}
}
