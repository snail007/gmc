// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gapp

import (
	"fmt"
	"github.com/snail007/gmc/core"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_parseConfigFile(t *testing.T) {
	assert := assert.New(t)
	app := New().(*GMCApp)
	assert.Nil(app.parseConfigFile())
}
func TestRun(t *testing.T) {
	assert := assert.New(t)
	app := New()
	assert.NotNil(app.Ctx())
	app.SetBlock(false)
	app.SetConfigFile("app.toml")
	server := gcore.Providers.HTTPServer("")(app.Ctx())
	assert.NotNil(server)
	app.AddService(gcore.ServiceItem{
		Service: server.(gcore.Service),
		BeforeInit: func(srv gcore.Service, cfg gcore.Config) (err error) {
			cfg.Set("template.dir", "../../http/template/tests/views")
			cfg.Set("httpserver.listen", ":")
			return
		},
	})
	err := app.Run()
	assert.Nil(err)
	app.OnShutdown(func() {})
	app.Stop()
	time.Sleep(time.Second)
}

func TestRun_1(t *testing.T) {
	assert := assert.New(t)
	app := New()
	app.SetBlock(false)
	app.SetConfigFile("app.toml")
	app.AttachConfigFile("007", "app.toml")
	app.AddService(gcore.ServiceItem{
		Service: gcore.Providers.HTTPServer("")(app.Ctx()).(gcore.Service),
	})
	app.AddService(gcore.ServiceItem{
		Service:  gcore.Providers.HTTPServer("")(app.Ctx()).(gcore.Service),
		ConfigID: "007",
	})
	err := app.Run()
	assert.NotNil(err)
	assert.NotSame(app.Config(), app.Config("007"))
}
func TestRun_2(t *testing.T) {
	assert := assert.New(t)
	app := New()
	app.SetBlock(false)
	app.SetConfigFile("app.toml")
	server := gcore.Providers.HTTPServer("")(app.Ctx()).(gcore.Service)
	app.AddService(gcore.ServiceItem{
		Service: server,
		BeforeInit: func(srv gcore.Service, config gcore.Config) (err error) {
			config.Set("session.store", "none")
			return
		},
	})
	err := app.Run()
	assert.NotNil(err)
	app.Stop()
}

func TestRun_3(t *testing.T) {
	assert := assert.New(t)
	app := New()
	app.SetBlock(false)
	app.SetConfigFile("app.toml")
	app.AddService(gcore.ServiceItem{
		Service: gcore.Providers.HTTPServer("")(app.Ctx()).(gcore.Service),
		AfterInit: func(srv *gcore.ServiceItem) (err error) {
			err = fmt.Errorf("error")
			return
		},
	})
	err := app.Run()
	assert.Equal(err.Error(), "error")
	app.Stop()
}

func TestSetConfig(t *testing.T) {
	// assert := assert.New(t)
	app := New().(*GMCApp)
	app.SetBlock(false)
	app.SetConfigFile("app.toml")
	app.parseConfigFile()
	app.parseConfigFile()
	app.SetConfigFile("none.toml")
	app.parseConfigFile()
}
func TestSetConfig_2(t *testing.T) {
	// assert := assert.New(t)
	app := New().(*GMCApp)
	app.SetBlock(false)
	app.SetConfigFile("none.toml")
	app.parseConfigFile()
}
func TestSetExtraConfig(t *testing.T) {
	assert := assert.New(t)
	app := New().(*GMCApp)
	app.SetBlock(false)
	app.SetConfigFile("app.toml")
	app.AttachConfigFile("extra01", "extra.toml")
	err := app.parseConfigFile()
	assert.Nil(err)
	assert.NotEmpty(app.Config().GetString("httpserver.listen"))
	assert.NotNil(app.Config("extra01"))
}
func TestSetExtraConfig_1(t *testing.T) {
	assert := assert.New(t)
	app := New().(*GMCApp)
	app.SetBlock(false)
	app.AttachConfigFile("001", "none.toml")
	err := app.parseConfigFile()
	assert.NotNil(err)
}
func TestBeforeRun(t *testing.T) {
	assert := assert.New(t)
	//run gmc app
	app := New().(*GMCApp)
	assert.Nil(app.parseConfigFile())
	app.OnRun(func(gcore.Config) error {
		return fmt.Errorf("stop")
	})
	err := app.Run()
	assert.Equal(err.Error(), "stop")
}
func TestBeforeRun_1(t *testing.T) {
	// assert := assert.New(t)
	app := New()
	app.SetBlock(false)
	app.OnRun(func(gcore.Config) (err error) {
		a := 0
		_ = a / a
		return
	})
	app.Run()
}
func TestBeforeRun_2(t *testing.T) {
	// assert := assert.New(t)
	app := New()
	app.SetBlock(false)
	app.OnRun(func(gcore.Config) (err error) {
		return fmt.Errorf("_")
	})
	app.Run()
}
func TestBeforeShutdown(t *testing.T) {
	assert := assert.New(t)
	//run gmc app
	app := New().(*GMCApp)
	app.SetBlock(false)
	assert.Nil(app.parseConfigFile())
	app.OnShutdown(func() {
		a := 0
		_ = a / a
		return
	})
	app.Run()
	app.Stop()
}
