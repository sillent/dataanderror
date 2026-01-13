package dataanderror

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type CustomKey string

const (
	One CustomKey = "one"
)

func (c *CustomKey) String() string {
	return string(*c)
}

func Test_DataAndErrorFromCustomTypeKey(t *testing.T) {
	de := New[CustomKey, string, string]()
	de.Store(One, "one")
	value, exist := de.Load(One)
	assert.Equal(t, true, exist)
	assert.Equal(t, "one", value)
}
