package amplitude_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude"
	"github.com/amplitude/analytics-go/amplitude/types"
)

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

type ClientSuite struct {
	suite.Suite
}

func (t *ClientSuite) TestOneEvent() {
	var payloads []string

	server := t.createTestServer(func(payload string) {
		payloads = append(payloads, payload)
	})
	defer server.Close()

	config := amplitude.NewConfig("your_api_key")
	config.ServerURL = server.URL
	config.FlushQueueSize = 3

	client := amplitude.NewClient(config)
	client.Add(&testBeforePlugin{})
	destPlugin := &testDestinationPlugin{}
	client.Add(destPlugin)

	client.Track(t.createEvent(1))

	client.Flush()
	client.Shutdown()

	t.Require().Equal(1, len(payloads))
	t.Require().JSONEq(`
{
  "api_key": "your_api_key",
  "events": [
    {
      "event_type": "event-1",
      "user_id": "user-1",
      "time": 1,
      "insert_id": "insert-1",
      "library": "amplitude-go/0.0.7",
      "ip": "IP 1",
      "event_properties": {
        "prop-1": 1
      }
    }
  ]
}`, payloads[0])

	t.Require().Equal(1, len(destPlugin.events))
}

func (t *ClientSuite) TestFlushQueueSize() {
	var payloads []string

	server := t.createTestServer(func(payload string) {
		payloads = append(payloads, payload)
	})
	defer server.Close()

	config := amplitude.NewConfig("your_api_key")
	config.ServerURL = server.URL
	config.FlushQueueSize = 2

	client := amplitude.NewClient(config)
	client.Add(&testBeforePlugin{})
	destPlugin := &testDestinationPlugin{}
	client.Add(destPlugin)

	client.Track(t.createEvent(2))
	client.Track(t.createEvent(1))
	client.Track(t.createEvent(3))

	client.Flush()
	client.Shutdown()

	t.Require().Equal(2, len(payloads))
	t.Require().JSONEq(`
{
  "api_key": "your_api_key",
  "events": [
    {
      "event_type": "event-2",
      "user_id": "user-2",
      "time": 2,
      "insert_id": "insert-2",
      "library": "amplitude-go/0.0.7",
      "ip": "IP 1",
      "event_properties": {
        "prop-2": 2
      }
    },
    {
      "event_type": "event-1",
      "user_id": "user-1",
      "time": 1,
      "insert_id": "insert-1",
      "library": "amplitude-go/0.0.7",
      "ip": "IP 2",
      "event_properties": {
        "prop-1": 1
      }
    }
  ]
}`, payloads[0])

	t.Require().JSONEq(`
{
  "api_key": "your_api_key",
  "events": [
    {
      "event_type": "event-3",
      "user_id": "user-3",
      "time": 3,
      "insert_id": "insert-3",
      "library": "amplitude-go/0.0.7",
      "ip": "IP 3",
      "event_properties": {
        "prop-3": 3
      }
    }
  ]
}`, payloads[1])

	t.Require().Equal(3, len(destPlugin.events))
}

func (t *ClientSuite) TestFlushInterval() {
	var payloads []string

	server := t.createTestServer(func(payload string) {
		payloads = append(payloads, payload)
	})
	defer server.Close()

	config := amplitude.NewConfig("your_api_key")
	config.ServerURL = server.URL
	config.FlushQueueSize = 999
	config.FlushInterval = time.Millisecond * 300

	client := amplitude.NewClient(config)
	client.Add(&testBeforePlugin{})
	destPlugin := &testDestinationPlugin{}
	client.Add(destPlugin)

	client.Track(t.createEvent(2))
	client.Track(t.createEvent(1))

	time.Sleep(time.Millisecond * 500)

	client.Track(t.createEvent(3))

	client.Flush()
	client.Shutdown()

	t.Require().Equal(2, len(payloads))
	t.Require().JSONEq(`
{
  "api_key": "your_api_key",
  "events": [
    {
      "event_type": "event-2",
      "user_id": "user-2",
      "time": 2,
      "insert_id": "insert-2",
      "library": "amplitude-go/0.0.7",
      "ip": "IP 1",
      "event_properties": {
        "prop-2": 2
      }
    },
    {
      "event_type": "event-1",
      "user_id": "user-1",
      "time": 1,
      "insert_id": "insert-1",
      "library": "amplitude-go/0.0.7",
      "ip": "IP 2",
      "event_properties": {
        "prop-1": 1
      }
    }
  ]
}`, payloads[0])

	t.Require().JSONEq(`
{
  "api_key": "your_api_key",
  "events": [
    {
      "event_type": "event-3",
      "user_id": "user-3",
      "time": 3,
      "insert_id": "insert-3",
      "library": "amplitude-go/0.0.7",
      "ip": "IP 3",
      "event_properties": {
        "prop-3": 3
      }
    }
  ]
}`, payloads[1])

	t.Require().Equal(3, len(destPlugin.events))
}

func (t *ClientSuite) TestPanicInPlugins() {
	var payloads []string

	server := t.createTestServer(func(payload string) {
		payloads = append(payloads, payload)
	})
	defer server.Close()

	config := amplitude.NewConfig("your_api_key")
	config.ServerURL = server.URL
	config.FlushQueueSize = 3

	client := amplitude.NewClient(config)
	client.Add(&testBeforePlugin{raisePanic: true})
	destPlugin := &testDestinationPlugin{raisePanic: true}
	client.Add(destPlugin)

	client.Track(t.createEvent(1))

	client.Flush()
	client.Shutdown()

	t.Require().Equal(1, len(payloads))
	t.Require().JSONEq(`
{
  "api_key": "your_api_key",
  "events": [
    {
      "event_type": "event-1",
      "user_id": "user-1",
      "time": 1,
      "insert_id": "insert-1",
      "library": "amplitude-go/0.0.7",
      "event_properties": {
        "prop-1": 1
      }
    }
  ]
}`, payloads[0])

	t.Require().Equal(0, len(destPlugin.events))
}

func (t *ClientSuite) TestConcurrentTrack() {
	var payloads []string
	var mu sync.Mutex

	server := t.createTestServer(func(payload string) {
		mu.Lock()
		defer mu.Unlock()
		payloads = append(payloads, payload)
	})
	defer server.Close()

	config := amplitude.NewConfig("your_api_key")
	config.ServerURL = server.URL
	config.FlushQueueSize = 33
	config.Logger = &noopLogger{}

	events := make(chan types.Event)

	clientCount := 3
	clients := make([]amplitude.Client, clientCount)
	var wg sync.WaitGroup
	for i := range clients {
		client := amplitude.NewClient(config)
		client.Add(&testBeforePlugin{})
		destPlugin := &testDestinationPlugin{}
		client.Add(destPlugin)
		clients[i] = client

		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range events {
				client.Track(event)
			}
		}()
	}

	eventCount := 13579
	for i := 1; i <= eventCount; i++ {
		events <- t.createEvent(1)
	}
	close(events)

	wg.Wait()

	for _, client := range clients {
		client.Flush()
		client.Shutdown()
	}
	time.Sleep(time.Millisecond * 500)
	server.Close()

	mu.Lock()
	defer mu.Unlock()

	serverEventCount := 0
	for _, payload := range payloads {
		var data struct {
			Events []interface{} `json:"events"`
		}
		err := json.Unmarshal([]byte(payload), &data)
		t.Require().Nil(err)
		serverEventCount += len(data.Events)
	}

	t.Require().Equal(eventCount, serverEventCount)
}

func (t *ClientSuite) createEvent(index int) amplitude.Event {
	postfix := fmt.Sprintf("-%d", index)

	return amplitude.Event{
		EventType: "event" + postfix,
		EventOptions: types.EventOptions{
			UserID:   "user" + postfix,
			InsertID: "insert" + postfix,
			Time:     int64(index),
		},
		EventProperties: map[string]interface{}{
			"prop" + postfix: index,
		},
	}
}

func (t *ClientSuite) createTestServer(onPayload func(payload string)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		t.Require().Nil(err)

		t.Require().True(json.Valid(body))

		onPayload(string(body))
	}))
}

type testBeforePlugin struct {
	currentIP  int
	raisePanic bool
}

func (p *testBeforePlugin) Name() string {
	return "test-before-plugin"
}

func (p *testBeforePlugin) Type() amplitude.PluginType {
	return amplitude.PluginTypeBefore
}

func (p *testBeforePlugin) Setup(types.Config) {
}

func (p *testBeforePlugin) Execute(event *amplitude.Event) *amplitude.Event {
	p.currentIP++

	if p.raisePanic {
		panic("panic in test-before-plugin")
	}

	event.IP = "IP " + strconv.Itoa(p.currentIP)

	return event
}

type testDestinationPlugin struct {
	raisePanic bool
	events     []*amplitude.Event
}

func (p *testDestinationPlugin) Name() string {
	return "test-destination-plugin"
}

func (p *testDestinationPlugin) Type() amplitude.PluginType {
	return amplitude.PluginTypeDestination
}

func (p *testDestinationPlugin) Setup(types.Config) {
	if p.raisePanic {
		panic("panic in test-destination-plugin")
	}
}

func (p *testDestinationPlugin) Execute(event *amplitude.Event) {
	if p.raisePanic {
		panic("panic in test-destination-plugin")
	}

	p.events = append(p.events, event)
}

type noopLogger struct{}

func (l noopLogger) Debugf(string, ...interface{}) {
}

func (l noopLogger) Infof(string, ...interface{}) {
}

func (l noopLogger) Warnf(string, ...interface{}) {
}

func (l noopLogger) Errorf(string, ...interface{}) {
}
