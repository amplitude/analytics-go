package amplitude

type Revenue struct {
	Price       float64
	Quantity    int
	ProductID   string
	RevenueType string
	Receipt     string
	ReceiptSig  string
	Properties  map[string]interface{}
	Revenue     float64
}

// IsValid checks if a Revenue object has a positive integer Quantity,
// and either Revenue or Price is set.
func (r Revenue) IsValid() bool {
	return r.Revenue != 0 || (r.Revenue == 0 && r.Price != 0) && r.Quantity > 0
}

func (r Revenue) ToRevenueEvent(eventOptions EventOptions) Event {
	return Event{
		EventType:       RevenueEventType,
		EventOptions:    eventOptions,
		EventProperties: r.GetEventProperties(),
	}
}

func (r Revenue) GetEventProperties() map[string]interface{} {
	eventProperties := make(map[string]interface{})
	eventProperties[RevenueProductID] = r.ProductID
	if r.Quantity == 0 {
		eventProperties[RevenueQuantity] = 1
	}
	eventProperties[RevenueQuantity] = r.Quantity
	eventProperties[RevenuePrice] = r.Price
	eventProperties[RevenueType] = r.RevenueType
	eventProperties[RevenueReceipt] = r.Receipt
	eventProperties[RevenueReceiptSig] = r.ReceiptSig
	if r.Revenue == 0 {
		eventProperties[DefaultRevenue] = r.Price * float64(r.Quantity)
	}
	eventProperties[DefaultRevenue] = r.Revenue

	return eventProperties
}
