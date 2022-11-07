package destination

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/amplitude/analytics-go/amplitude/types"
)

func newAmplitudeHTTPClient(
	serverURL string, options clientPayloadOptions, logger types.Logger, connectionTimeout time.Duration,
) *amplitudeHTTPClient {
	var payloadOptions *clientPayloadOptions
	if options != (clientPayloadOptions{}) {
		payloadOptions = &options
	}

	return &amplitudeHTTPClient{
		serverURL:      serverURL,
		logger:         logger,
		payloadOptions: payloadOptions,
		httpClient: &http.Client{
			Timeout: connectionTimeout,
		},
	}
}

type clientPayloadOptions struct {
	MinIDLength int `json:"min_id_length,omitempty"`
}

type clientPayload struct {
	APIKey  string                `json:"api_key"`
	Events  []*types.Event        `json:"events"`
	Options *clientPayloadOptions `json:"options,omitempty"`
}

type sendResult struct {
	Code    int
	Message string
}

type amplitudeResponse struct {
	Error string `json:"error"`
}

type amplitudeHTTPClient struct {
	serverURL      string
	logger         types.Logger
	payloadOptions *clientPayloadOptions
	httpClient     *http.Client
}

func (h *amplitudeHTTPClient) send(payload clientPayload) sendResult {
	if len(payload.Events) == 0 {
		return sendResult{}
	}

	payload.Options = h.payloadOptions
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		h.logger.Errorf("payload encoding failed: \n\tError: %w\n\tpayload: %+v", err, payload)

		return sendResult{
			Message: fmt.Sprintf("Payload encoding failed: %s", err),
		}
	}

	h.logger.Debugf("payloadBytes:\n\t%s", string(payloadBytes))

	request, err := http.NewRequest(http.MethodPost, h.serverURL, bytes.NewReader(payloadBytes))
	if err != nil {
		h.logger.Errorf("Building new request failed: \n\t%w", err)

		return sendResult{
			Message: fmt.Sprintf("Building new request failed: %s", err),
		}
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")

	response, err := h.httpClient.Do(request)
	if err != nil {
		h.logger.Errorf("HTTP request failed: %s", err)

		if response != nil {
			return sendResult{
				Code:    response.StatusCode,
				Message: fmt.Sprintf("HTTP request failed: %s", err),
			}
		} else {
			return sendResult{
				Message: fmt.Sprintf("HTTP request failed: %s", err),
			}
		}
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			h.logger.Warnf("HTTP response, close body: %s", err)
		}
	}()

	h.logger.Infof("HTTP response code: %s", response.Status)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		h.logger.Warnf("HTTP response, can't read body: %s", err)

		return sendResult{
			Code:    response.StatusCode,
			Message: fmt.Sprintf("HTTP response, can't read body: %s", err),
		}
	}

	h.logger.Infof("HTTP response body: %s", string(body))

	var message string
	var amplitudeResponse amplitudeResponse
	if err := json.Unmarshal(body, &amplitudeResponse); err == nil {
		message = amplitudeResponse.Error
	}

	return sendResult{
		Code:    response.StatusCode,
		Message: message,
	}
}
