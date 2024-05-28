package internal_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/amplitude/analytics-go/amplitude/plugins/destination/internal"
	"github.com/amplitude/analytics-go/amplitude/types"
)

func TestAmplitudeHTTPClient(t *testing.T) {
	suite.Run(t, new(AmplitudeHTTPClientSuiteSuite))
}

type AmplitudeHTTPClientSuiteSuite struct {
	suite.Suite
}

func (t *AmplitudeHTTPClientSuiteSuite) TestSend_Success() {
	server := t.createTestServer(0, 200, `{"code": 234, "error": "some server error"}`)

	client := internal.NewAmplitudeHTTPClient(
		server.URL,
		internal.AmplitudePayloadOptions{MinIDLength: 7},
		noopLogger{},
		time.Millisecond*1000,
	)

	response := client.Send(internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{
			t.createEvent(3),
			t.createEvent(1),
			t.createEvent(2),
		},
	})

	t.Require().Equal(internal.AmplitudeResponse{
		Status: 200,
		Code:   234,
		Error:  "some server error",
	}, response)

	server.Close()

	server.mu.Lock()
	defer server.mu.Unlock()

	t.Require().Equal(1, len(server.payloads))
	t.Require().JSONEq(`
{
  "api_key": "my-api-key",
  "events": [
    {
      "event_type": "event-3",
      "user_id": "user-3",
      "time": 3,
      "insert_id": "insert-3",
      "event_properties": {
        "prop-3": 3
      }
    },
    {
      "event_type": "event-1",
      "user_id": "user-1",
      "time": 1,
      "insert_id": "insert-1",
      "event_properties": {
        "prop-1": 1
      }
    },
    {
      "event_type": "event-2",
      "user_id": "user-2",
      "time": 2,
      "insert_id": "insert-2",
      "event_properties": {
        "prop-2": 2
      }
    }
  ],
  "options": {
    "min_id_length": 7
  }
}`, server.payloads[0])
}

func (t *AmplitudeHTTPClientSuiteSuite) TestSend_Empty() {
	server := t.createTestServer(0, 200, `{"code": 234, "error": "some server error"}`)

	client := internal.NewAmplitudeHTTPClient(
		server.URL,
		internal.AmplitudePayloadOptions{MinIDLength: 7},
		noopLogger{},
		time.Millisecond*1000,
	)

	response := client.Send(internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: nil,
	})

	t.Require().Equal(internal.AmplitudeResponse{}, response)

	server.Close()

	server.mu.Lock()
	defer server.mu.Unlock()

	t.Require().Empty(server.payloads)
}

func (t *AmplitudeHTTPClientSuiteSuite) TestSend_Timeout() {
	timeout := time.Millisecond * 100
	server := t.createTestServer(timeout * 2, 200, `{"code": 234, "error": "some server error"}`)

	client := internal.NewAmplitudeHTTPClient(
		server.URL,
		internal.AmplitudePayloadOptions{MinIDLength: 7},
		noopLogger{},
		timeout,
	)

	response := client.Send(internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{
			t.createEvent(1),
		},
	})

	t.Require().NotNil(response.Err)

	var urlErr *url.Error
	isURLErr := errors.As(response.Err, &urlErr)
	t.Require().True(isURLErr && urlErr.Timeout())

	server.Close()
}

func (t *AmplitudeHTTPClientSuiteSuite) TestSend_NonJsonResponse() {
	server := t.createTestServer(0, 413, `<html>
	<head><title>413 Request Entity Too Large</title></head>
	<body>
	<center><h1>413 Request Entity Too Large</h1></center>
	<hr><center>nginx</center>
	</body>
	</html>`)

	client := internal.NewAmplitudeHTTPClient(
		server.URL,
		internal.AmplitudePayloadOptions{MinIDLength: 7},
		noopLogger{},
		time.Millisecond*1000,
	)

	response := client.Send(internal.AmplitudePayload{
		APIKey: "my-api-key",
		Events: []*types.Event{
			t.createEvent(1),
		},
	})

	t.Require().Equal(internal.AmplitudeResponse{
		Status: 413,
		Code:   413,
	}, response)

	server.Close()

	server.mu.Lock()
	defer server.mu.Unlock()

	t.Require().Equal(1, len(server.payloads))
	t.Require().JSONEq(`
{
  "api_key": "my-api-key",
  "events": [
    {
      "event_type": "event-1",
      "user_id": "user-1",
      "time": 1,
      "insert_id": "insert-1",
      "event_properties": {
        "prop-1": 1
      }
    }
  ],
  "options": {
    "min_id_length": 7
  }
}`, server.payloads[0])
}

type testServer struct {
	*httptest.Server
	mu       sync.Mutex
	payloads []string
}

func (t *AmplitudeHTTPClientSuiteSuite) createTestServer(delay time.Duration, statusCode int, responseBody string) *testServer {
	server := &testServer{}

	server.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			time.Sleep(delay)
		}

		body, err := io.ReadAll(r.Body)
		t.Require().Nil(err)

		t.Require().True(json.Valid(body))

		server.mu.Lock()
		defer server.mu.Unlock()
		server.payloads = append(server.payloads, string(body))

		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(responseBody))
	}))

	return server
}

func (t *AmplitudeHTTPClientSuiteSuite) createEvent(index int) *types.Event {
	postfix := fmt.Sprintf("-%d", index)

	return &types.Event{
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

type noopLogger struct{}

func (l noopLogger) Debugf(string, ...interface{}) {
}

func (l noopLogger) Infof(string, ...interface{}) {
}

func (l noopLogger) Warnf(string, ...interface{}) {
}

func (l noopLogger) Errorf(string, ...interface{}) {
}
