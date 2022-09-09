package destination

import (
	"bytes"
	"encoding/json"
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

type amplitudeHTTPClient struct {
	serverURL      string
	logger         types.Logger
	payloadOptions *clientPayloadOptions
	httpClient     *http.Client
}

func (h *amplitudeHTTPClient) send(payload clientPayload) {
	if len(payload.Events) == 0 {
		return
	}

	payload.Options = h.payloadOptions
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		h.logger.Errorf("payload encoding failed: \n\tError: %w\n\tpayload: %+v", err, payload)

		return
	}

	h.logger.Debugf("payloadBytes:\n\t%s", string(payloadBytes))

	request, err := http.NewRequest(http.MethodPost, h.serverURL, bytes.NewReader(payloadBytes))
	if err != nil {
		h.logger.Errorf("Building new request failed: \n\t%w", err)

		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")

	response, err := h.httpClient.Do(request)
	if err != nil {
		h.logger.Errorf("HTTP request failed: %s", err)

		return
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
	} else {
		h.logger.Infof("HTTP response body: %s", string(body))
	}
}
