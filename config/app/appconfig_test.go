package appconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)
	cfg := NewAPPConfig()
	err := Parse(cfg)
	assert.Nil(err)
	FileConfig.Set("session.store", "memory")
	err = initSessionStore()
	assert.Nil(err)
	FileConfig.Set("session.store", "redis")
	err = initSessionStore()
	assert.Nil(err)
	assert.NotNil(SessionStore)
	// assert.Fail("")
}
