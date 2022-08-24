package amplitude

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
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
	if len(p.Events) == 0 {
		return
	}

	payloadBytes, err := json.Marshal(p)
	if err != nil {
		h.logger.Errorf("payload encoding failed: \n\tError: %w\n\tpayload: %+v", err, p)
	}

	h.logger.Debugf("payloadBytes:\n\t%s", string(payloadBytes))

	request, err := http.NewRequest("POST", h.serverURL, bytes.NewReader(payloadBytes))
	if err != nil {
		h.logger.Errorf("Building new request failed: \n\t%w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	response, err := httpClient.Do(request)
	if err != nil {
		h.logger.Errorf("HTTP request failed", err)
	}
	defer response.Body.Close()

	h.logger.Infof("HTTP response:\n\t%+v", response)
}
