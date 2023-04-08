// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"bytes"
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNewFromConfig(t *testing.T) {
	cfg := gconfig.New()
	cfg.SetConfigType("toml")
	cfg.ReadConfig(bytes.NewReader([]byte(`
				[log]
				# 0,1,2,3,4,5 => TRACE, DEBUG, INFO, WARN, ERROR, PANIC
				level=0
				# 0,1 => console, file
				output=[0,1]
				# only worked when output contains 1
				logsDir="./logs"
				# filename in logs logsDir.
				# available placeholders are:
				# %Y:Year 2020, %m:Month 10, %d:Day 10, %H:24Hours 21
				filename="web_%s.log"
				gzip=true
				async=true`)))
	l := NewFromConfig(cfg, "")
	l.Info("test")
	time.Sleep(time.Second)
	l.Info("test")
	assert.Implements(t, (*gcore.Logger)(nil), l)
	time.Sleep(time.Second)
	assert.DirExists(t, "logs")
	assert.Nil(t, os.RemoveAll("logs"))
}

func TestNewFromConfig1(t *testing.T) {
	cfg := gconfig.New()
	cfg.SetConfigType("toml")
	cfg.ReadConfig(bytes.NewReader([]byte(`
				[log]
				# 0,1,2,3,4,5 => TRACE, DEBUG, INFO, WARN, ERROR, PANIC
				level=0
				# 0,1 => console, file
				output=[0]
				# only worked when output contains 1
				logsDir="./logs"
				# filename in logs logsDir.
				# available placeholders are:
				# %Y:Year 2020, %m:Month 10, %d:Day 10, %H:24Hours 21
				filename="web_%s.log"
				gzip=true
				async=true`)))
	l := NewFromConfig(cfg, "")
	assert.Implements(t, (*gcore.Logger)(nil), l)
}

func TestNewFromConfig2(t *testing.T) {
	cfg := gconfig.New()
	cfg.SetConfigType("toml")
	cfg.ReadConfig(bytes.NewReader([]byte(` `)))
	l := NewFromConfig(cfg, "")
	assert.Implements(t, (*gcore.Logger)(nil), l)
}

func Test_existsDir(t *testing.T) {
	f := "t.mp.tmp"
	assert.Nil(t, ioutil.WriteFile(f, []byte("\n"), 0755))
	assert.False(t, existsDir(f))
	os.RemoveAll(f)
	assert.Nil(t, os.Mkdir(f, 0755))
	assert.True(t, existsDir(f))
	os.RemoveAll(f)
}
