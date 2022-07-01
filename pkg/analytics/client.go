package analytics

import "github.com/amplitude/Amplitude-Go/internal"

type Amplitude struct {
	configuration Config
	timeline      internal.TimeLine
}

func (a Amplitude) track(event internal.BaseEvent) {

}

func (a Amplitude) identify(identityObj internal.Identity, eventOptions internal.EventOptions, eventProperties map[string]string) {

}

func (a Amplitude) groupIdentify(groupType string, groupName string, identifyObj internal.Identity, EventOptions internal.EventOptions, eventProperties map[string]string, userProperties map[string]string) {

}

func (a Amplitude) revenue(revenueObj internal.Revenue, eventOptions internal.EventOptions) {

}

func (a Amplitude) setGroup(groupType string, groupName string, eventOptions internal.EventOptions) {

}

func (a Amplitude) flush() {

}

func (a Amplitude) add(plugin internal.Plugin) {

}

func (a Amplitude) remove(plugin internal.Plugin) {

}

func (a Amplitude) shutdown() {

}