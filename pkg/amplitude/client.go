package amplitude

import "github.com/amplitude/Amplitude-Go/internal"

type Amplitude struct {
	Configuration Config
	timeline      internal.TimeLine
}

func (a Amplitude) track(event BaseEvent) {

}

func (a Amplitude) identify(identityObj Identity, eventOptions EventOptions, eventProperties map[string]string) {

}

func (a Amplitude) groupIdentify(groupType string, groupName string, identifyObj Identity, EventOptions EventOptions, eventProperties map[string]string, userProperties map[string]string) {

}

func (a Amplitude) revenue(revenueObj Revenue, eventOptions EventOptions) {

}

func (a Amplitude) setGroup(groupType string, groupName string, eventOptions EventOptions) {

}

func (a Amplitude) flush() {

}

func (a Amplitude) add(plugin internal.Plugin) {

}

func (a Amplitude) remove(plugin internal.Plugin) {

}

func (a Amplitude) shutdown() {

}
