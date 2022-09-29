package destination

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/amplitude/analytics-go/amplitude/types"
)

type AmplitudeResponseProcessor struct {
	EventStorage           types.EventStorage
	MaxRetries             int
	RetryBaseInterval      time.Duration
	RetryThrottledInterval time.Duration
	Now                    func() time.Time
	Logger                 types.Logger
}

func (p *AmplitudeResponseProcessor) Process(events []*types.Event, response AmplitudeResponse) AmplitudeResult {
	responseStatus := response.normalizedStatus()

	var urlErr *url.Error
	isURLErr := errors.As(response.Err, &urlErr)

	switch {
	case response.Err == nil && responseStatus == http.StatusOK:
		return p.processSuccess(events, response)
	case (isURLErr && urlErr.Timeout()) || responseStatus == http.StatusRequestTimeout || responseStatus == http.StatusInternalServerError:
		return p.processTimeout(events, response)
	case responseStatus == http.StatusRequestEntityTooLarge:
		return p.processTooLargeRequest(events, response)
	case responseStatus == http.StatusBadRequest:
		return p.processBadRequest(events, response)
	case responseStatus == http.StatusTooManyRequests:
		return p.processTooManyRequests(events, response)
	}

	return p.processUnknownError(events, response)
}

func (p *AmplitudeResponseProcessor) processSuccess(events []*types.Event, response AmplitudeResponse) AmplitudeResult {
	return AmplitudeResult{
		Events:  events,
		Code:    response.Code,
		Message: "Event sent successfully.",
	}
}

func (p *AmplitudeResponseProcessor) processTimeout(events []*types.Event, response AmplitudeResponse) AmplitudeResult {
	eventsForCallback := make([]*types.Event, 0, len(events))
	eventsForRetry := make([]*types.Event, 0, len(events))
	now := p.Now()
	for _, event := range events {
		if event.RetryCount >= p.MaxRetries {
			eventsForCallback = append(eventsForCallback, event)
		} else {
			event.RetryCount++
			event.RetryAt = now.Add(p.retryInterval(event.RetryCount))
			eventsForRetry = append(eventsForRetry, event)
		}
	}

	p.EventStorage.ReturnBack(eventsForRetry...)

	result := AmplitudeResult{
		Events:  eventsForCallback,
		Code:    response.Code,
		Message: fmt.Sprintf("Event reached max retry times %d", p.MaxRetries),
	}
	p.logResult(p.Logger.Errorf, result)

	return result
}

func (p *AmplitudeResponseProcessor) processTooLargeRequest(events []*types.Event, response AmplitudeResponse) AmplitudeResult {
	if len(events) == 1 {
		result := AmplitudeResult{
			Events:  events,
			Code:    response.Code,
			Message: response.Error,
		}
		p.logResult(p.Logger.Errorf, result)

		return result
	}

	p.Logger.Warnf("RequestEntityTooLarge: chunk size is reduced")
	p.EventStorage.ReduceChunkSize()
	p.EventStorage.ReturnBack(events...)

	return AmplitudeResult{
		Code:    response.Code,
		Message: response.Error,
	}
}

func (p *AmplitudeResponseProcessor) processBadRequest(events []*types.Event, response AmplitudeResponse) AmplitudeResult {
	switch {
	case strings.HasPrefix(response.Error, "Invalid API key:"):
		result := AmplitudeResult{
			Events:  events,
			Code:    response.Code,
			Message: "Invalid API key",
		}
		p.logResult(p.Logger.Errorf, result)

		return result
	case response.MissingField != "":
		result := AmplitudeResult{
			Events:  events,
			Code:    response.Code,
			Message: fmt.Sprintf("Request missing required field %s", response.MissingField),
		}
		p.logResult(p.Logger.Errorf, result)

		return result
	}

	invalidIndexes := response.invalidOrSilencedEventIndexes()
	eventsForCallback := make([]*types.Event, 0, len(events))
	eventsForRetry := make([]*types.Event, 0, len(events))
	for i, event := range events {
		if _, ok := invalidIndexes[i]; ok {
			eventsForCallback = append(eventsForCallback, event)
		} else {
			eventsForRetry = append(eventsForRetry, event)
		}
	}

	p.EventStorage.ReturnBack(eventsForRetry...)

	result := AmplitudeResult{
		Events:  eventsForCallback,
		Code:    response.Code,
		Message: response.Error,
	}
	p.logResult(p.Logger.Errorf, result)

	return result
}

func (p *AmplitudeResponseProcessor) processTooManyRequests(events []*types.Event, response AmplitudeResponse) AmplitudeResult {
	eventsForCallback := make([]*types.Event, 0, len(events))
	eventsForRetry := make([]*types.Event, 0, len(events))
	eventsForRetryDelay := make([]*types.Event, 0, len(events))
	now := p.Now()

	for i, event := range events {
		if response.throttledEventIndex(i) {
			if response.exceedDailyQuota(event) {
				eventsForCallback = append(eventsForCallback, event)
			} else {
				event.RetryAt = now.Add(p.RetryThrottledInterval)
				eventsForRetryDelay = append(eventsForRetryDelay, event)
			}
		} else {
			eventsForRetry = append(eventsForRetry, event)
		}
	}

	p.EventStorage.ReturnBack(eventsForRetryDelay...)
	p.EventStorage.ReturnBack(eventsForRetry...)

	result := AmplitudeResult{
		Events:  eventsForCallback,
		Code:    response.Code,
		Message: "Exceeded daily quota",
	}
	p.logResult(p.Logger.Errorf, result)

	return result
}

func (p *AmplitudeResponseProcessor) processUnknownError(events []*types.Event, response AmplitudeResponse) AmplitudeResult {
	errMessage := response.Error
	if response.Err != nil {
		errMessage = response.Err.Error()
	}
	if errMessage == "" {
		errMessage = "Unknown error"
	}

	result := AmplitudeResult{
		Events:  events,
		Code:    response.Code,
		Message: errMessage,
	}
	p.logResult(p.Logger.Errorf, result)

	return result
}

func (p *AmplitudeResponseProcessor) retryInterval(retries int) time.Duration {
	return p.RetryBaseInterval * (1 << ((retries - 1) / 2))
}

func (p *AmplitudeResponseProcessor) logResult(logFunc func(message string, args ...interface{}), result AmplitudeResult) {
	if len(result.Events) == 0 {
		return
	}

	eventsJSON, _ := json.Marshal(result.Events)
	logFunc("%s: code=%d, events=%s", result.Message, result.Code, eventsJSON)
}
