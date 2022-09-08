package amplitude

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func newAmplitudeClient(
	serverURL string, options clientPayloadOptions, logger Logger, connectionTimeout time.Duration,
) *amplitudeClient {
	var payloadOptions *clientPayloadOptions
	if options != (clientPayloadOptions{}) {
		payloadOptions = &options
	}

	return &amplitudeClient{
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
	Events  []*Event              `json:"events"`
	Options *clientPayloadOptions `json:"options,omitempty"`
}

type amplitudeClient struct {
	serverURL      string
	logger         Logger
	payloadOptions *clientPayloadOptions
	httpClient     *http.Client
}

func (h *amplitudeClient) send(payload clientPayload) {
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
