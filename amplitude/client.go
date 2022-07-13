package amplitude

type Amplitude struct {
	Configuration Config
	timeline      Timeline
}

func (a Amplitude) Track(event BaseEvent) {

}

func (a Amplitude) Identify(identityObj Identity, eventOptions EventOptions, eventProperties map[string]string) {

}

func (a Amplitude) GroupIdentify(groupType string, groupName string, identifyObj Identity, eventOptions EventOptions, eventProperties map[string]string) {

}

func (a Amplitude) Revenue(revenueObj Revenue, eventOptions EventOptions) {

}

func (a Amplitude) SetGroup(groupType string, groupName string, eventOptions EventOptions) {

}

func (a Amplitude) Flush() {

}

func (a Amplitude) Add(plugin Plugin) {

}

func (a Amplitude) Remove(plugin Plugin) {

}

func (a Amplitude) Shutdown() {

}
