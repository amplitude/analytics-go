package types

import (
	"time"
)

type EventCallback = func(event Event, code int, message ...string)

type Config struct {
	APIKey             string
	FlushInterval      time.Duration
	FlushQueueSize     int
	FlushMaxRetries    int
	Logger             Logger
	MinIDLength        int
	Callback           EventCallback
	ServerZone         ServerZone
	UseBatch           bool
	StorageFactory     func() EventStorage
	OptOut             bool
	Plan               *Plan
	ServerURL          string
	ConnectionTimeout  time.Duration
	MaxStorageCapacity int
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
