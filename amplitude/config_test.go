package amplitude

import "testing"

func TestConfigIsValid(t *testing.T) {
	c := Config{}
	if c.isValid() {
		t.Error("Config is not valid")
	}
}

func TestConfigIsMinIdLengthValid(t *testing.T) {
	c := Config{}
	if c.isMinIdLengthValid() {
		t.Error("isMinIdLengthValid is not valid")
	}
}
