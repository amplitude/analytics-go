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

func NewRevenue(price float64) Revenue {
	return Revenue{
		Quantity: 1,
		Price:    price,
	}
}

// IsValid checks if a revenue instance has a positive integer quantity.
func (r Revenue) IsValid() bool {
	return r.Quantity > 0
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
	eventProperties[RevenueQuantity] = r.Quantity
	eventProperties[RevenuePrice] = r.Price
	eventProperties[RevenueType] = r.RevenueType
	eventProperties[RevenueReceipt] = r.Receipt
	eventProperties[RevenueReceiptSig] = r.ReceiptSig
	eventProperties[DefaultRevenue] = r.Revenue

	return eventProperties
}
