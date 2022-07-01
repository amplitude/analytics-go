package amplitude

type Plan struct {
	branch    string
	source    string
	version   string
	versionId string
}

func (p *Plan) getPlanBody() {

}

type EventOptions struct {
}

type BaseEvent struct {
	EventOptions
}

type Identity struct {
}

type GroupIdentifyEvent struct {
	BaseEvent
}

type IdentifyEvent struct {
	BaseEvent
}

type Revenue struct {
}

type RevenueEvent struct {
	Revenue
}
