package amplitude

import (
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	c := Config{}

	if c.IsValid() {
		t.Errorf("Expected: false, , but got %t", c.IsValid())
	}

	c.APIKey = "test"
	if c.IsValid() {
		t.Errorf("Expected: false, , but got %t", c.IsValid())
	}

	c.FlushQueueSize = DefaultFlushQueueSize
	if c.IsValid() {
		t.Errorf("Expected: false, , but got %t", c.IsValid())
	}

	c.FlushInterval = DefaultFlushInterval
	if c.IsValid() {
		t.Errorf("Expected: false, , but got %t", c.IsValid())
	}

	c.MinIDLength = DefaultMinIDLength
	if !c.IsValid() {
		t.Errorf("Expected: true, , but got %t", c.IsValid())
	}
}

func TestConfigIsMinIDLengthValid(t *testing.T) {
	c := Config{}

	if c.IsMinIDLengthValid() {
		t.Errorf("Expected: false, , but got %t", c.IsMinIDLengthValid())
	}

	c.MinIDLength = DefaultMinIDLength
	if !c.IsMinIDLengthValid() {
		t.Errorf("Expected: true, , but got %t", c.IsMinIDLengthValid())
	}
}
