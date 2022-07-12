package amplitude

type Plan struct {
	branch    string
	source    string
	version   string
	versionID string
}

func (p *Plan) getPlanBody() {

}

type EventOptions struct {
}

type BaseEvent struct {
	EventOptions
}

func (b BaseEvent) loadEventOptions(options EventOptions) {

}

// An Identity helps to generate an IdentifyEvent or a GroupIdentifyEvent instance
// with special eventType and userProperties/ groupProperties
type Identity struct {
	propertiesSet map[string]string
	properties    map[string]string
}

func (i Identity) set(key string, value string) {

}

func (i Identity) append(key string, value string) {

}

func (i Identity) prepend(key string, value string) {

}

func (i Identity) preInsert(key string, value string) {

}

func (i Identity) postInsert(key string, value string) {

}

func (i Identity) remove(key string, value string) {

}

func (i Identity) add(key string, value string) {

}

func (i Identity) unset(key string) {

}

func (i Identity) clearAll() {

}

func (i Identity) isValid() bool {
	return true
}

type GroupIdentifyEvent struct {
	BaseEvent
}

type IdentifyEvent struct {
	BaseEvent
}

// A Revenue helps to generate a RevenueEvent
// with special event type
// and revenue information like price, quantity, product id, receipt,etc.
type Revenue struct {
	price       float64
	quantity    int
	productId   string
	revenueType string
	receipt     string
	receiptSig  string
	properties  map[string]string
	revenue     float64
}

func (r Revenue) setReceipt(receipt string, receiptSignature string) {

}

func (r Revenue) isValid() bool {
	return true
}

// creates and returns a RevenueEvent instance, sets revenue information as eventProperties
func (r Revenue) toRevenueEvent() RevenueEvent {
	return RevenueEvent{}
}

func (r Revenue) getEventProperties() {

}

type RevenueEvent struct {
	BaseEvent
}
