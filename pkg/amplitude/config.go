package amplitude

import (
	"github.com/amplitude/Amplitude-Go/internal"
	"log"
)

type Config struct {
	ApiKey              string
	FlushIntervalMillis int64
	FlushMaxRetries     int64
	Logger              log.Logger
	MidIdLength         int64
	Callback            interface{}
	ServerZone          string
	UseBatch            bool
	StorageProvider     internal.StorageProvider
	OptOut              bool
	plan                Plan

	Url              string
	FlushQueueSize   int64
	FlushSizeDivider int
}

func getStorage(c *Config) internal.Storage {
	return c.StorageProvider.GetStorage()
}

func (c Config) isValid() bool {
	if c.ApiKey != "" || c.FlushQueueSize <= 0 || c.FlushIntervalMillis <= 0 || !c.isMinIdLengthValid() {
		return false
	}
	return true
}

func (c Config) isMinIdLengthValid() bool {
	if c.MidIdLength > 0 {
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
