<p align="center">
  <a href="https://amplitude.com" target="_blank" align="center">
    <img src="https://static.amplitude.com/lightning/46c85bfd91905de8047f1ee65c7c93d6fa9ee6ea/static/media/amplitude-logo-with-text.4fb9e463.svg" width="280">
  </a>
  <br />
</p>

<div align="center">

[![Build](https://github.com/amplitude/analytics-go/actions/workflows/build.yml/badge.svg)](https://github.com/amplitude/analytics-go/actions/workflows/build.yml)
[![go-doc](https://pkg.go.dev/badge/github.com/amplitude/analytics-go?utm_source=godoc)](https://pkg.go.dev/github.com/amplitude/analytics-go)

</div>

# Announcement 📣

Amplitude is introducing [Go SDK](https://pkg.go.dev/github.com/amplitude/analytics-go). Compared to plain [HTTP V2 API](https://www.docs.developers.amplitude.com/analytics/apis/http-v2-api/), it provides improved developer experience, helps users instrument data more seamlessly and provides more control over data being instrumented using custom plugins. 

To learn more about the new SDK, here are some useful links:

* Go Packages: https://pkg.go.dev/github.com/amplitude/analytics-go
* GitHub: https://github.com/amplitude/analytics-go
* Documentation: https://www.docs.developers.amplitude.com/data/sdks/go/


# Official Amplitude Go SDK

This is Amplitude's latest and official version of Go SDK.

## Installation 

Install analytics-go using go get:

```
go get https://github.com/amplitude/analytics-go
```

## Usage
```go
package main

import (
	"github.com/amplitude/analytics-go/amplitude"
)

func main() {

	config := amplitude.NewConfig("your-api-key")

	client := amplitude.NewClient(config)

	// Track a basic event
	// One of UserID and DeviceID is required
	client.Track(amplitude.Event{
		EventType: "Button Clicked",
		EventOptions: amplitude.EventOptions{
			UserID: "user-id",
		},
	})

	// Flushed queued events and shutdown the client
	client.Shutdown()
}
```


## Need Help?
If you have any issues using our SDK, feel free to [create a GitHub issue](https://github.com/amplitude/analytics-go/issues/new) or submit a request on [Amplitude Help](https://help.amplitude.com/hc/en-us/requests/new).

