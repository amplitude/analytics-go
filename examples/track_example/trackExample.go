package track_example

import "github.com/amplitude/Amplitude-Go/pkg/amplitude"

func callbackFunc(e string, code int, message string) {
	println(e)
	println(code, message)
}

var config amplitude.Config = amplitude.Config{ApiKey: "your_api_key"}
var client amplitude.Amplitude = amplitude.Amplitude{Configuration: config}

var event amplitude.BaseEvent = amplitude.BaseEvent{}
