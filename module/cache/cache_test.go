package gcache_test

import (
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/module/cache"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func Test_Cache(t *testing.T) {
	assert := assert2.New(t)
	ctx := gcore.Providers.Ctx("")()
	c := gcore.Providers.Config("")()
	c.SetConfigFile("../app/app.toml")
	err := c.ReadInConfig()
	assert.Nil(err)
	ctx.SetConfig(c)
	_, err = gcore.Providers.Cache("")(ctx)
	assert.Nil(err)
	assert.NotNil(gcache.Memory())
	assert.NotNil(gcache.Redis())
	assert.NotNil(gcache.File())
	assert.Same(gcache.Cache(), gcache.Redis())
}
