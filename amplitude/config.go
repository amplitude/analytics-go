package amplitude

import (
	"log"
	"time"
)

type Callback = func(event BaseEvent, code int, message ...string)

type Config struct {
	APIKey              string
	FlushIntervalMillis time.Duration
	FlushQueueSize      int
	FlushMaxRetries     int
	Logger              log.Logger
	MinIDLength         int
	Callback            Callback
	ServerZone          string
	UseBatch            bool
	StorageProvider     StorageProvider
	OptOut              bool
	Plan                Plan
	ServerURL           string
}

func getStorage(c *Config) Storage {
	return c.StorageProvider.GetStorage()
}

func (c Config) IsValid() bool {
	if c.APIKey == "" || c.FlushQueueSize <= 0 || c.FlushIntervalMillis <= 0 || !c.IsMinIDLengthValid() {
		return false
	}
	return true
}

func (c Config) IsMinIDLengthValid() bool {
	return c.MinIDLength > 0
}
