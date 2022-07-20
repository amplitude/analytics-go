package amplitude

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigIsValid(t *testing.T) {
	c := Config{}
	assert.False(t, c.IsValid())

	c.APIKey = "test"
	assert.False(t, c.IsValid())

	c.FlushQueueSize = DefaultFlushQueueSize
	assert.False(t, c.IsValid())

	c.FlushInterval = DefaultFlushInterval
	assert.False(t, c.IsValid())

	c.MinIDLength = DefaultMinIDLength
	assert.True(t, c.IsValid())
}

func TestConfigIsMinIDLengthValid(t *testing.T) {
	c := Config{}
	assert.False(t, c.IsMinIDLengthValid())

	c.MinIDLength = DefaultMinIDLength
	assert.True(t, c.IsMinIDLengthValid())
}
