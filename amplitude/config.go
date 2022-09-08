package amplitude

import (
	"time"
)

type EventCallback = func(event Event, code int, message ...string)

type Config struct {
	APIKey            string
	FlushInterval     time.Duration
	FlushQueueSize    int
	FlushMaxRetries   int
	Logger            Logger
	MinIDLength       int
	Callback          EventCallback
	ServerZone        ServerZone
	UseBatch          bool
	StorageFactory    func() EventStorage
	OptOut            bool
	Plan              Plan
	ServerURL         string
	ConnectionTimeout time.Duration
}

func NewConfig(apiKey string) Config {
	return Config{
		APIKey: apiKey,
	}
}

func (c Config) IsValid() bool {
	if c.APIKey == "" || c.FlushQueueSize <= 0 || c.FlushInterval <= 0 || c.MinIDLength < 0 {
		return false
	}

	return true
}

func (c Config) setDefaultValues() Config {
	clone := c

	if clone.FlushInterval == 0 {
		clone.FlushInterval = DefaultConfig.FlushInterval
	}
	if clone.FlushQueueSize == 0 {
		clone.FlushQueueSize = DefaultConfig.FlushQueueSize
	}
	if clone.FlushMaxRetries == 0 {
		clone.FlushMaxRetries = DefaultConfig.FlushMaxRetries
	}
	if clone.Logger == nil {
		clone.Logger = NewDefaultLogger()
	}
	if clone.StorageFactory == nil {
		clone.StorageFactory = func() EventStorage {
			return NewInMemoryEventStorage(clone.FlushQueueSize)
		}
	}
	if clone.ServerZone == "" {
		clone.ServerZone = DefaultConfig.ServerZone
	}
	if clone.ServerURL == "" {
		if clone.UseBatch {
			clone.ServerURL = ServerBatchURLs[clone.ServerZone]
		} else {
			clone.ServerURL = ServerURLs[clone.ServerZone]
		}
	}

	return clone
}
