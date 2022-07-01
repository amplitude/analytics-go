package analytics

import "github.com/amplitude/Amplitude-Go/internal"

type Config struct {
	apiKey              string
	flushQueueSize      int64
	flushIntervalMillis int64
	flushMaxRetries     int64
	//logger
	midIdLength int64
	//callback
	serverZone      string
	useBatch        bool
	serverUrl       string
	storageProvider internal.StorageProvider
	plan            internal.Plan
}

func getStorage(c *Config) {
}

func (c Config) isValid() {

}

func (c Config) isMinIdLengthValid() {

}

func (c Config) increaseFlushDivider() {

}

func (c Config) resetFlushDivider() {

}