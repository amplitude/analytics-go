package amplitude

import "github.com/amplitude/Amplitude-Go/internal"

type Amplitude struct {
	Configuration Config
	timeline      internal.Timeline
}

func (a Amplitude) Track(event BaseEvent) {

}

func (a Amplitude) Identify(identityObj Identity, eventOptions EventOptions, eventProperties map[string]string) {

}

func (a Amplitude) GroupIdentify(groupType string, groupName string, identifyObj Identity, EventOptions EventOptions, eventProperties map[string]string) {

}

func (a Amplitude) Revenue(revenueObj Revenue, eventOptions EventOptions) {

}

func (a Amplitude) SetGroup(groupType string, groupName string, eventOptions EventOptions) {

}

func (a Amplitude) Flush() {

}

func (a Amplitude) Add(plugin internal.Plugin) {

}

func (a Amplitude) Remove(plugin internal.Plugin) {

}

func (a Amplitude) Shutdown() {

}
