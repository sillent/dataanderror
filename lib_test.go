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

func Test_Remove(t *testing.T) {
	de := New[string, string, string]()
	de.Store("first", "firstdata")
	value, exist := de.Load("first")
	assert.Equal(t, true, exist)
	assert.Equal(t, "firstdata", value)
	de.Remove("first")
	value, exist = de.Load("first")
	assert.Equal(t, false, exist)
	assert.Equal(t, "", value)
}

func Test_RemoveError(t *testing.T) {
	de := New[string, string, string]()
	de.StoreError("first", "firsterror")
	value := de.CopiedError()
	assert.Equal(t, "firsterror", value["first"])
	de.RemoveError("first")
	value = de.CopiedError()
	assert.Equal(t, "", value["first"])
}

func Test_LoadError(t *testing.T) {
	de := New[string, string, string]()
	de.StoreError("first", "error")
	value, exist := de.LoadError("first")
	assert.Equal(t, true, exist)
	assert.Equal(t, "error", value)
}
