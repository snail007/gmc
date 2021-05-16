package gonce

import (
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadOnce(t *testing.T) {
	assert := assert2.New(t)
	once := LoadOnce("test")
	assert.NotNil(once)
	RemoveOnce("test")
	once1 := LoadOnce("test")
	assert.NotNil(once, once1)
}

func TestLoadOnce_1(t *testing.T) {
	assert := assert2.New(t)
	once := LoadOnce("test")
	once1 := LoadOnce("test")
	assert.Equal(once, once1)
}
