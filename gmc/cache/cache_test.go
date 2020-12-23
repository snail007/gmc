package gcache_test

import (
	gconfig "github.com/snail007/gmc/config"
	gcachehelper "github.com/snail007/gmc/gmc/cache/helper"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func Test_Cache(t *testing.T) {
	assert := assert2.New(t)
	cfg, err := gconfig.NewConfigFile("../app/app.toml")
	assert.Nil(err)
	//fmt.Println(cfg.Get("cache"))

	var c0 map[string]interface{}
	cfg.UnmarshalKey("cache", &c0)
	c0["memory"].([]interface{})[0].(map[string]interface{})["enable"] = true
	c0["file"].([]interface{})[0].(map[string]interface{})["enable"] = true
	cfg.Set("cache", c0)

	//fmt.Println(cfg.Get("cache"))
	gcachehelper.Init(cfg)
	assert.NotNil(gcachehelper.Memory())
	assert.NotNil(gcachehelper.Redis())
	assert.NotNil(gcachehelper.File())
	assert.Same(gcachehelper.Cache(), gcachehelper.Redis())
}
