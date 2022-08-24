package amplitude

import (
	"time"
)

type EventCallback = func(event Event, code int, message ...string)

type Config struct {
	APIKey          string
	FlushInterval   time.Duration
	FlushQueueSize  int
	FlushMaxRetries int
	Logger          Logger
	MinIDLength     int
	Callback        EventCallback
	ServerZone      ServerZone
	UseBatch        bool
	Storage         Storage
	OptOut          bool
	Plan            Plan
	ServerURL       string
}

func NewConfig(apiKey string) Config {
	return Config{
		APIKey:          apiKey,
		FlushInterval:   DefaultFlushInterval,
		FlushQueueSize:  DefaultFlushQueueSize,
		FlushMaxRetries: DefaultFlushMaxRetries,
		Logger:          NewDefaultLogger(),
		MinIDLength:     DefaultMinIDLength,
		Callback:        nil,
		ServerZone:      ServerZoneUS,
		UseBatch:        false,
		Storage:         &InMemoryStorage{},
		OptOut:          false,
		ServerURL:       HTTPV2,
	}
}

func (c Config) IsValid() bool {
	if c.APIKey == "" || c.FlushQueueSize <= 0 || c.FlushInterval <= 0 || !c.IsMinIDLengthValid() {
		return false
	}

	return true
}

func (c Config) IsMinIDLengthValid() bool {
	return c.MinIDLength > 0
}
