package gmccache_test

import (
	gmccachehelper "github.com/snail007/gmc/cache/helper"
	gmcconfig "github.com/snail007/gmc/config"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func Test_Cache(t *testing.T) {
	assert := assert2.New(t)
	cfg, err := gmcconfig.NewFile("../app/app.toml")
	assert.Nil(err)
	//fmt.Println(cfg.Get("cache"))

	var c0 map[string]interface{}
	cfg.UnmarshalKey("cache", &c0)
	c0["memory"].([]interface{})[0].(map[string]interface{})["enable"] = true
	cfg.Set("cache", c0)

	//fmt.Println(cfg.Get("cache"))
	gmccachehelper.Init(cfg)
	assert.NotNil(gmccachehelper.Memory())
	assert.NotNil(gmccachehelper.Redis())
	assert.Same(gmccachehelper.Cache(), gmccachehelper.Redis())
}
