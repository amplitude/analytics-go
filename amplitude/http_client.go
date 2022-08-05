package amplitude

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type payload struct {
	APIKey string   `json:"api_key"`
	Events []*Event `json:"Events"`
}

type httpClient struct {
	logger    Logger
	serverURL string
}

func (h *httpClient) send(p payload) {
	payloadBytes, err := json.Marshal(p)
	if err != nil {
		h.logger.Error("payload encoding failed", err)
	}

	h.logger.Debug("payloadBytes: ", string(payloadBytes))

	request, err := http.NewRequest("POST", h.serverURL, bytes.NewReader(payloadBytes))
	if err != nil {
		h.logger.Error("Building new request failed", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)
	if err != nil {
		h.logger.Error("HTTP request failed", err)
	}
	defer response.Body.Close()

	h.logger.Info("HTTP request response", response)
}
