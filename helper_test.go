// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmc

import (
	"bytes"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	gdb "github.com/snail007/gmc/module/db"
	gcaptcha "github.com/snail007/gmc/util/captcha"
	gmap "github.com/snail007/gmc/util/map"
	gtest "github.com/snail007/gmc/util/testing"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

var (
	dbOnce   = sync.Once{}
	mysqlCfg = bytes.NewReader([]byte(`
				[database]
				default="mysql"
				
				[[database.mysql]]
				enable=true
				id="default"
				host="127.0.0.1"
				port="3306"
				username="user"
				password="user"
				database="test"
				prefix=""
				prefix_sql_holder="__PREFIX__"
				charset="utf8"
				collate="utf8_general_ci"
				maxidle=30
				maxconns=200
				timeout=15000
				readtimeout=15000
				writetimeout=15000
				maxlifetimeseconds=1800`))
	sqlite3Cfg = bytes.NewReader([]byte(`
				[database]
				default="sqlite3"
				
				[[database.sqlite3]]
				enable=true
				id="default"
				database="test.db"
				password=""
				prefix=""
				prefix_sql_holder="__PREFIX__"
				# sync mode: 0:OFF, 1:NORMAL, 2:FULL, 3:EXTRA
				syncmode=0
				# open mode: ro,rw,rwc,memory
				openmode="rw"
				# cache mode: shared,private
				cachemode="shared"`))
	cacheRedisCfg = bytes.NewReader([]byte(`
				[cache]
				default="redis"
				[[cache.redis]]
				debug=true
				enable=true
				id="default"
				address="127.0.0.1:6379"
				prefix=""
				password=""
				timeout=10
				dbnum=0
				maxidle=10
				maxactive=30
				idletimeout=300
				maxconnlifetime=3600
				wait=false`))
	cacheFileCfg = bytes.NewReader([]byte(`
				[cache]
				default="file"
				[[cache.file]]
				enable=true
				id="default"
				dir="{tmp}"
				cleanupinterval=30`))
	cacheMemCfg = bytes.NewReader([]byte(`
				[cache]
				default="memory"
				[[cache.memory]]
				enable=true
				id="default"
				cleanupinterval=30`))
	i18nCfg = bytes.NewReader([]byte(`
				[i18n]
				enable=true
				dir="i18n"
				default="zh-cn"`))
)

func TestCacheAssistant_Cache(t *testing.T) {
	assert := assert.New(t)
	if gtest.RunProcess(t, func() {
		cfg := gconfig.New()
		cfg.SetConfigType("toml")
		cfg.ReadConfig(cacheMemCfg)
		err := Cache.Init(cfg)
		assert.NoError(err)
		assert.Implements((*gcore.Cache)(nil), Cache.Memory())
		assert.Implements((*gcore.Cache)(nil), Cache.Cache())
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
}

func TestCacheAssistant_File(t *testing.T) {
	assert := assert.New(t)
	if gtest.RunProcess(t, func() {
		cfg := gconfig.New()
		cfg.SetConfigType("toml")
		cfg.ReadConfig(cacheFileCfg)
		Cache.Init(cfg)
		assert.Implements((*gcore.Cache)(nil), Cache.File())
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
}

func TestCacheAssistant_Redis(t *testing.T) {

	assert := assert.New(t)
	if gtest.RunProcess(t, func() {
		cfg := gconfig.New()
		cfg.SetConfigType("toml")
		cfg.ReadConfig(cacheRedisCfg)
		Cache.Init(cfg)
		assert.Implements((*gcore.Cache)(nil), Cache.Redis())
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
}

func TestDBAssistant_DB(t *testing.T) {
	assert := assert.New(t)
	if gtest.RunProcess(t, func() {
		cfg := gconfig.New()
		cfg.SetConfigType("toml")
		cfg.ReadConfig(mysqlCfg)
		DB.Init(cfg)
		assert.Implements((*gcore.Database)(nil), DB.DB())
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
}

func TestDBAssistant_InitFromFile(t *testing.T) {
	DB.InitFromFile("module/app/app.toml")
	assert.Implements(t, (*gcore.Database)(nil), DB.DB())
}

func TestDBAssistant_MySQL(t *testing.T) {
	assert := assert.New(t)
	if gtest.RunProcess(t, func() {
		cfg := gconfig.New()
		cfg.SetConfigType("toml")
		cfg.ReadConfig(mysqlCfg)
		DB.Init(cfg)
		assert.Implements((*gcore.Database)(nil), DB.MySQL())
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
}

func TestDBAssistant_SQLite3(t *testing.T) {
	assert := assert.New(t)
	if gtest.RunProcess(t, func() {
		cfg := gconfig.New()
		cfg.SetConfigType("toml")
		cfg.ReadConfig(sqlite3Cfg)
		DB.Init(cfg)
		assert.Implements((*gcore.Database)(nil), DB.SQLite3())
		defer os.Remove("test.db")
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
}

func TestDBAssistant_Table(t *testing.T) {
	assert := assert.New(t)
	if gtest.RunProcess(t, func() {
		cfg := gconfig.New()
		cfg.SetConfigType("toml")
		cfg.ReadConfig(mysqlCfg)
		DB.Init(cfg)
		assert.IsType(&gdb.Model{}, DB.Table("test"))
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
}

func TestErrorAssistant_New(t *testing.T) {
	e := fmt.Errorf("abc")
	assert.Equal(t, "abc", Err.New(e).Error())
}

func TestErrorAssistant_Recover(t *testing.T) {
	s := ""
	g := sync.WaitGroup{}
	g.Add(1)
	a := 1
	b := 0
	go func() {
		defer Err.Recover(func(e interface{}) {
			s = fmt.Sprintf("%v", e)
			g.Done()
		})
		_ = a / b
	}()
	g.Wait()
	assert.Contains(t, s, "zero")
}

func TestErrorAssistant_RecoverString(t *testing.T) {
	g := sync.WaitGroup{}
	g.Add(1)
	a := 1
	b := 0
	go func() {
		defer Err.Recover("abc",true)
		defer g.Done()
		_ = a / b
	}()
	g.Wait()
}

func TestErrorAssistant_RecoverObj(t *testing.T) {
	g := sync.WaitGroup{}
	g.Add(1)
	a := 1
	b := 0
	go func() {
		defer Err.Recover(fmt.Errorf(""),true)
		defer g.Done()
		_ = a / b
	}()
	g.Wait()
}

func TestErrorAssistant_Stack(t *testing.T) {
	e := fmt.Errorf("abc")
	assert.Contains(t, Err.Stack(e), "helper_test.go")
	assert.Contains(t, Err.Stack(e), "abc")
}

func TestErrorAssistant_Wrap(t *testing.T) {
	e := fmt.Errorf("abc")
	assert.Contains(t, Err.Wrap(e).Error(), "abc")
}

func TestI18nAssistant_Init(t *testing.T) {
	assert := assert.New(t)
	if gtest.RunProcess(t, func() {
		cfg := gconfig.New()
		cfg.SetConfigType("toml")
		cfg.ReadConfig(i18nCfg)
		i, err := I18n.Init(cfg)
		ctx := New.Ctx()
		ctx.SetConfig(cfg)
		assert.NoError(err)
		assert.Implements((*gcore.I18n)(nil), i)
		assert.Implements((*gcore.I18n)(nil), New.Tr("", ctx))
	}) {
		return
	}
	_, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
}

func TestNewAssistant_APIServer(t *testing.T) {
	s, err := New.APIServer(New.Ctx(), ":")
	assert.NoError(t, err)
	assert.Implements(t, (*gcore.APIServer)(nil), s)
}

func TestNewAssistant_APIServerDefault(t *testing.T) {
	s, err := New.APIServerDefault(New.Ctx())
	assert.NoError(t, err)
	assert.Implements(t, (*gcore.APIServer)(nil), s)
}

func TestNewAssistant_App(t *testing.T) {
	assert.Implements(t, (*gcore.App)(nil), New.App())
}

func TestNewAssistant_AppDefault(t *testing.T) {
	assert.Implements(t, (*gcore.App)(nil), New.AppDefault())
}

func TestNewAssistant_Captcha(t *testing.T) {
	assert.IsType(t, &gcaptcha.Captcha{}, New.Captcha())
}

func TestNewAssistant_CaptchaDefault(t *testing.T) {
	assert.IsType(t, &gcaptcha.Captcha{}, New.CaptchaDefault())
}

func TestNewAssistant_Config(t *testing.T) {
	assert.Implements(t, (*gcore.Config)(nil), New.Config())
}

func TestNewAssistant_ConfigFile(t *testing.T) {
	s, err := New.ConfigFile("module/app/app.toml")
	assert.NoError(t, err)
	assert.Implements(t, (*gcore.Config)(nil), s)
}

func TestNewAssistant_Ctx(t *testing.T) {
	assert.Implements(t, (*gcore.Ctx)(nil), New.Ctx())
}

func TestNewAssistant_HTTPServer(t *testing.T) {
	assert.Implements(t, (*gcore.HTTPServer)(nil), New.HTTPServer(New.Ctx()))
}

func TestNewAssistant_Logger(t *testing.T) {
	assert.Implements(t, (*gcore.Logger)(nil), New.Logger(New.Ctx(), ""))
}

func TestNewAssistant_Map(t *testing.T) {
	assert.IsType(t, &gmap.Map{}, New.Map())
}

func TestNewAssistant_Router(t *testing.T) {
	assert.Implements(t, (*gcore.HTTPRouter)(nil), New.Router(New.Ctx()))
}

func TestNewErrorAssistant(t *testing.T) {
	assert.IsType(t, &ErrorAssistant{}, NewErrorAssistant())
}

func Test_initHelper(t *testing.T) {
	initHelper()
	assert.NotNil(t, Err.p)
}
