// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/snail007/gmc/core"
	glog "github.com/snail007/gmc/module/log"
	_ "github.com/snail007/gmc/using/basic"
	assert2 "github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	assert := assert2.New(t)
	assert.Implements(new(gcore.Logger), gcore.ProviderLogger()(nil, ""))
}

func TestLogger_SetOutput(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l.Info("a")
	assert.True(strings.HasSuffix(out.String(), "INFO a\n"))
}

func TestLogger_Writer(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	assert.Equal(&out, l.Writer())
}

func TestLogger_SetLevel(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l.SetLevel(gcore.LogLeveWarn)
	l.Info("a")
	assert.Empty(out.String())
}

func TestLogger_With_1(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l0 := l.With("api")
	l0.Info("a", "b")
	t.Log(out.String())
	assert.True(strings.HasSuffix(out.String(), "[api] INFO ab\n"))
	assert.Equal(l0.Namespace(), "api")
}

func TestLogger_With_2(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l0 := l.With("api").With("user").With("list")
	l0.Info("a")
	t.Log(out.String())
	assert.True(strings.HasSuffix(out.String(), "[api/user/list] INFO a\n"))
}

func TestLogger_Infof(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l.Infof("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "INFO a10\n"))
}

func TestLogger_Trace(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetLevel(gcore.LogLevelTrace)
	l.SetOutput(&out)
	l.Trace("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "TRACE a\n"))
}

func TestLogger_Tracef(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetLevel(gcore.LogLevelTrace)
	l.SetOutput(&out)
	l.Tracef("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "TRACE a10\n"))
}

func TestLogger_Debug(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l.Debug("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "DEBUG a\n"))
}

func TestLogger_Debugf(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l.Debugf("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "DEBUG a10\n"))
}

func TestLogger_Warn(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l.Warn("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "WARN a\n"))
}

func TestLogger_Warnf(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(&out)
	l.Warnf("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "WARN a10\n"))
}

func TestLogger_Panic(t *testing.T) {
	assert := assert2.New(t)
	l := gcore.ProviderLogger()(nil, "")
	defer gcore.ProviderError()().Recover(func(e interface{}) {
		assert.Contains(e, "PANIC a")
	})
	l.Panic("a")
}

func TestLogger_Panicf(t *testing.T) {
	assert := assert2.New(t)
	l := gcore.ProviderLogger()(nil, "")
	defer gcore.ProviderError()().Recover(func(e interface{}) {
		assert.Contains(e, "PANIC a10")
	})
	l.Panicf("a%d", 10)
}

func TestLogger_Error(t *testing.T) {
	assert := assert2.New(t)
	l := gcore.ProviderLogger()(nil, "")
	if os.Getenv("ASSERT_EXISTS_"+t.Name()) == "1" {
		l.Error("a")
	}
	var out bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), "ASSERT_EXISTS_"+t.Name()+"=1")
	cmd.Stdout = &out
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.True(strings.HasSuffix(out.String(), "ERROR a\n"))
		return
	} else {
		assert.Fail("expecting unsuccessful exit")
	}
}

func TestLogger_Errorf(t *testing.T) {
	assert := assert2.New(t)
	l := gcore.ProviderLogger()(nil, "")
	if os.Getenv("ASSERT_EXISTS_"+t.Name()) == "1" {
		l.Errorf("a%d", 10)
	}
	var out bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), "ASSERT_EXISTS_"+t.Name()+"=1")
	cmd.Stdout = &out
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.True(strings.HasSuffix(out.String(), "ERROR a10\n"))
		return
	} else {
		assert.Fail("expecting unsuccessful exit")
	}
}

func TestLogger_WithRate(t *testing.T) {
	t.Parallel()
	cnt := 0
	l := glog.New()
	l0 := l.WithRate(time.Second)
	l0.SetOutput(ioutil.Discard)
	l0.SetRateCallback(func(msg string) {
		cnt++
	})
	for i := 0; i < 35; i++ {
		l0.Write("hello")
		time.Sleep(time.Millisecond * 100)
	}
	assert2.True(t, cnt >= 3)
}
