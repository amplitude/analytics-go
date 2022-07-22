package amplitude

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigIsValid(t *testing.T) {
	config := NewConfig("test-api-key")
	assert.True(t, config.IsValid())

	config = NewConfig("")
	assert.False(t, config.IsValid())

	config = NewConfig("test-api-key")
	config.FlushQueueSize = 0
	assert.False(t, config.IsValid())

	config = NewConfig("test-api-key")
	config.MinIDLength = 0
	assert.False(t, config.IsValid())
}
