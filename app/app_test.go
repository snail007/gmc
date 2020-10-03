package gmcapp

import (
	"fmt"
	"testing"
	"time"

	gmcconfig "github.com/snail007/gmc/config"

	httpserver "github.com/snail007/gmc/http/server"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	assert := assert.New(t)
	app := New()
	assert.Nil(app.ParseConfig())
}
func TestRun(t *testing.T) {
	assert := assert.New(t)
	app := New().Block(false)
	assert.Nil(app.ParseConfig())
	app.Config().Set("httpserver.listen", ":")
	server := httpserver.New()
	app.AddService(ServiceItem{
		Service: server,
	})
	err := app.Run()
	assert.Nil(err)
	app.BeforeShutdown(func() {})
	app.Stop()
	time.Sleep(time.Second)
}

func TestRun_1(t *testing.T) {
	assert := assert.New(t)
	app := New().Block(false)
	app.AddExtraConfigFile("007", "app.toml")
	assert.Nil(app.ParseConfig())
	app.AddService(ServiceItem{
		Service: httpserver.New(),
	})
	app.AddService(ServiceItem{
		Service:      httpserver.New(),
		ConfigIDname: "007",
	})
	err := app.Run()
	assert.NotNil(err)
	assert.NotEqual(app.Config(), app.Config("007"))
	app.Stop()
}
func TestRun_2(t *testing.T) {
	assert := assert.New(t)
	app := New().Block(false)
	assert.Nil(app.ParseConfig())
	app.Config().Set("session.store", "none")
	server := httpserver.New()
	app.AddService(ServiceItem{
		Service: server,
	})
	err := app.Run()
	assert.NotNil(err)
	app.Stop()
}

func TestRun_3(t *testing.T) {
	assert := assert.New(t)
	app := New().Block(false)
	assert.Nil(app.ParseConfig())
	app.AddService(ServiceItem{
		Service: httpserver.New(),
		AfterInit: func(srv *ServiceItem) (err error) {
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
	app := New().Block(false)
	app.SetMainConfigFile("app.toml")
	app.ParseConfig()
	app.ParseConfig()
	app.SetMainConfigFile("none.toml")
	app.ParseConfig()
}
func TestSetConfig_2(t *testing.T) {
	// assert := assert.New(t)
	app := New().Block(false)
	app.SetMainConfigFile("none.toml")
	app.ParseConfig()
}
func TestSetExtraConfig(t *testing.T) {
	assert := assert.New(t)
	app := New().Block(false)
	app.AddExtraConfigFile("extra01", "extra.toml")
	err := app.ParseConfig()
	assert.Nil(err)
	assert.NotEmpty(app.Config().GetString("httpserver.listen"))
	assert.NotNil(app.Config("extra01"))
}
func TestSetExtraConfig_1(t *testing.T) {
	assert := assert.New(t)
	app := New().Block(false)
	app.AddExtraConfigFile("001", "none.toml")
	err := app.ParseConfig()
	assert.NotNil(err)
}
func TestBeforeRun(t *testing.T) {
	assert := assert.New(t)
	//run gmc app
	app := New()
	assert.Nil(app.ParseConfig())
	app.BeforeRun(func(*gmcconfig.Config) error {
		return fmt.Errorf("stop")
	})
	err := app.Run()
	assert.Equal(err.Error(), "stop")
}
func TestBeforeRun_1(t *testing.T) {
	// assert := assert.New(t)
	app := New().Block(false)
	app.BeforeRun(func(*gmcconfig.Config) (err error) {
		a := 0
		_ = a / a
		return
	})
	app.Run()
}
func TestBeforeRun_2(t *testing.T) {
	// assert := assert.New(t)
	app := New().Block(false)
	app.BeforeRun(func(*gmcconfig.Config) (err error) {
		return fmt.Errorf(".")
	})
	app.Run()
}
func TestBeforeShutdown(t *testing.T) {
	assert := assert.New(t)
	//run gmc app
	app := New().Block(false)
	assert.Nil(app.ParseConfig())
	app.BeforeShutdown(func() {
		a := 0
		_ = a / a
		return
	})
	app.Run()
	app.Stop()
}
