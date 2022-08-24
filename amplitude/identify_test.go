package amplitude

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyIdentifyInstanceNotValid(t *testing.T) {
	identify := Identify{}

	assert.False(t, identify.IsValid())
}

func TestSetKeySuccessAndValid(t *testing.T) {
	identify := Identify{}
	identify.Set("test-set", 15)

	expectUserProperty := map[IdentityOp]map[string]interface{}{}
	expectUserProperty[IdentityOpSet] = map[string]interface{}{
		"test-set": 15,
	}

	assert.Equal(t, expectUserProperty, identify.Properties)
	assert.True(t, identify.IsValid())
}

func TestSetOnceKeySuccessAndValid(t *testing.T) {
	identify := Identify{}
	identify.SetOnce("test-set-once", 10)

	expectUserProperty := map[IdentityOp]map[string]interface{}{}
	expectUserProperty[IdentityOpSetOnce] = map[string]interface{}{
		"test-set-once": 10,
	}

	assert.Equal(t, expectUserProperty, identify.Properties)
	assert.True(t, identify.IsValid())
}

func TestAppendKeySuccessAndValid(t *testing.T) {
	identify := Identify{}
	identify.Append("test-append", 10)

	expectUserProperty := map[IdentityOp]map[string]interface{}{}
	expectUserProperty[IdentityOpAppend] = map[string]interface{}{
		"test-append": 10,
	}

	assert.Equal(t, expectUserProperty, identify.Properties)
	assert.True(t, identify.IsValid())
}