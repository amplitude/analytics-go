package amplitude

import (
	"log"
)

type Config struct {
	ApiKey              string
	FlushIntervalMillis int
	FlushMaxRetries     int
	Logger              log.Logger
	MinIdLength         int
	Callback            interface{}
	ServerZone          string
	UseBatch            bool
	StorageProvider     StorageProvider
	OptOut              bool
	plan                Plan

	Url              string
	FlushQueueSize   int
	FlushSizeDivider int
}

func getStorage(c *Config) Storage {
	return c.StorageProvider.GetStorage()
}

func (c Config) isValid() bool {
	if c.ApiKey != "" || c.FlushQueueSize <= 0 || c.FlushIntervalMillis <= 0 || !c.isMinIdLengthValid() {
		return false
	}
	return true
}

func (c Config) isMinIdLengthValid() bool {
	if c.MinIdLength > 0 {
		return true
	}
	return false
}

func (c Config) increaseFlushDivider() {
	c.FlushSizeDivider += 1
}

func (c Config) resetFlushDivider() {
	c.FlushSizeDivider = 1
}
