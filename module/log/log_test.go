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
	"sync/atomic"
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
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Info("a")
	assert.True(strings.HasSuffix(out.String(), "INFO a\n"))
}

func TestLogger_Writer(t *testing.T) {
	assert := assert2.New(t)
	out := glog.NewLoggerWriter(bytes.NewBuffer(nil))
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(out)
	assert.Equal(out, l.Writer())
}

func TestLogger_SetLevel(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.SetLevel(gcore.LogLeveWarn)
	l.Info("a")
	assert.Empty(out.String())
}

func TestLogger_With_1(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := glog.New()
	l.SetOutput(glog.NewLoggerWriter(&out))
	l0 := l.With("api")
	l0.Info("a", "b")
	t.Log(out.String())
	assert.True(strings.Contains(out.String(), "[api] INFO ab\n"))
	assert.True(strings.Contains(out.String(), "log_test.go"))
	assert.Equal(l0.Namespace(), "api")
}

func TestLogger_With_2(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l0 := l.With("api").With("user").With("list")
	l0.Info("a")
	t.Log(out.String())
	assert.True(strings.HasSuffix(out.String(), "[api/user/list] INFO a\n"))
	assert.True(strings.Contains(out.String(), "log_test.go"))
}

func TestLogger_Infof(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Infof("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "INFO a10\n"))
}

func TestLogger_Trace(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetLevel(gcore.LogLevelTrace)
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Trace("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "TRACE a\n"))
}

func TestLogger_Tracef(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetLevel(gcore.LogLevelTrace)
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Tracef("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "TRACE a10\n"))
}

func TestLogger_Debug(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Debug("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "DEBUG a\n"))
}

func TestLogger_Debugf(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Debugf("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "DEBUG a10\n"))
}

func TestLogger_Warn(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Warn("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "WARN a\n"))
}

func TestLogger_Warnf(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Warnf("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "WARN a10\n"))
}

func TestLogger_Error(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Error("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "ERROR a\n"))
}

func TestLogger_Errorf(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Errorf("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "ERROR a10\n"))
}

func TestLogger_Panic(t *testing.T) {
	os.Setenv("DISABLE_CONSOLE_COLOR", "true")
	defer os.Unsetenv("DISABLE_CONSOLE_COLOR")
	assert := assert2.New(t)
	l := gcore.ProviderLogger()(nil, "")
	var err interface{}
	defer func() {
		assert.Contains(err, "PANIC a")
	}()
	defer gcore.ProviderError()().Recover(func(e interface{}) {
		err = e
	})
	l.Panic("a")
}

func TestLogger_Panicf(t *testing.T) {
	os.Setenv("DISABLE_CONSOLE_COLOR", "true")
	defer os.Unsetenv("DISABLE_CONSOLE_COLOR")
	assert := assert2.New(t)
	l := gcore.ProviderLogger()(nil, "")
	defer gcore.ProviderError()().Recover(func(e interface{}) {
		assert.Contains(e, "PANIC a10")
	})
	l.Panicf("a%d", 10)
}

func TestLogger_Fatal(t *testing.T) {
	assert := assert2.New(t)
	l := gcore.ProviderLogger()(nil, "")
	if os.Getenv("ASSERT_EXISTS_"+t.Name()) == "1" {
		l.Fatal("a")
	}
	var out bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), "ASSERT_EXISTS_"+t.Name()+"=1")
	cmd.Stdout = &out
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.True(strings.Contains(out.String(), "FATAL a\n"))
		return
	} else {
		assert.Fail("expecting unsuccessful exit")
	}
}

func TestLogger_Fatalf(t *testing.T) {
	assert := assert2.New(t)
	l := gcore.ProviderLogger()(nil, "")
	if os.Getenv("ASSERT_EXISTS_"+t.Name()) == "1" {
		l.Fatalf("a%d", 10)
	}
	var out bytes.Buffer
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), "ASSERT_EXISTS_"+t.Name()+"=1")
	cmd.Stdout = &out
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.True(strings.Contains(out.String(), "FATAL a10\n"))
		return
	} else {
		assert.Fail("expecting unsuccessful exit")
	}
}

func TestLogger_WithRate(t *testing.T) {
	t.Parallel()
	cnt := new(int32)
	l := glog.New()
	l0 := l.WithRate(time.Second)
	l0.SetOutput(glog.NewLoggerWriter(ioutil.Discard))
	l0.SetRateCallback(func(msg string) {
		atomic.AddInt32(cnt, 1)
	})
	for i := 0; i < 35; i++ {
		l0.Write("hello", gcore.LogLeveInfo)
		time.Sleep(time.Millisecond * 100)
	}
	assert2.True(t, atomic.LoadInt32(cnt) >= 3)
}

func TestLogger_Write1(t *testing.T) {
	t.Parallel()
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.Write("abc", gcore.LogLeveInfo)
	assert.Equal(strings.Contains(out.String(), "log/log_test.go:"),
		strings.HasSuffix(out.String(), "abc\n"))
}

func TestLogger_Write2(t *testing.T) {
	t.Parallel()
	assert := assert2.New(t)
	var out bytes.Buffer
	l := gcore.ProviderLogger()(nil, "")
	l.SetOutput(glog.NewLoggerWriter(&out))
	l.WriteRaw("abc", gcore.LogLeveInfo)
	assert.Equal(!strings.Contains(out.String(), "log/log_test.go:"),
		strings.HasSuffix(out.String(), "abc\n"))
}

func TestLogger_Write3(t *testing.T) {
	t.Parallel()
	assert := assert2.New(t)
	var out bytes.Buffer
	l := glog.New()
	l.SetLevel(gcore.LogLeveNone)
	l.AddLevelWriter(&out, gcore.LogLevelTrace)
	l.Trace("foo1")
	l.Debug("foo2")
	l.Info("foo3")
	l.Warn("foo4")
	assert.Contains(out.String(), "foo1\n")
	assert.Contains(out.String(), "foo2\n")
	assert.Contains(out.String(), "foo3\n")
	assert.Contains(out.String(), "foo4\n")
}

func TestLogger_Write4(t *testing.T) {
	t.Parallel()
	assert := assert2.New(t)
	var out bytes.Buffer
	l := glog.New()
	l.SetLevel(gcore.LogLeveNone)
	l.AddLevelWriter(&out, gcore.LogLevelTrace)
	l.Tracef("foo1")
	l.Debugf("foo2")
	l.Infof("foo3")
	l.Warnf("foo4")
	assert.Contains(out.String(), "foo1\n")
	assert.Contains(out.String(), "foo2\n")
	assert.Contains(out.String(), "foo3\n")
	assert.Contains(out.String(), "foo4\n")
	t.Log(out.String())
}

func TestLogger_Write5(t *testing.T) {
	t.Parallel()
	assert := assert2.New(t)
	var out bytes.Buffer
	l := glog.New()
	l.AddLevelWriter(&out, gcore.LogLeveInfo)
	l.SetLevel(gcore.LogLeveNone)
	l.Tracef("foo1")
	l.Debugf("foo2")
	l.Infof("foo3")
	l.Warnf("foo4")
	assert.NotContains(out.String(), "foo1")
	assert.NotContains(out.String(), "foo2")
	assert.Contains(out.String(), "foo3\n")
	assert.Contains(out.String(), "foo4\n")
}

func TestLogger_Write6(t *testing.T) {
	t.Parallel()
	assert := assert2.New(t)
	var out bytes.Buffer
	l := glog.New()
	l.AddLevelsWriter(&out, gcore.LogLeveDebug, gcore.LogLeveWarn)
	l.SetLevel(gcore.LogLeveNone)
	l.Tracef("foo1")
	l.Debugf("foo2")
	l.Infof("foo3")
	l.Warnf("foo4")
	assert.NotContains(out.String(), "foo1")
	assert.NotContains(out.String(), "foo3")
	assert.Contains(out.String(), "foo2\n")
	assert.Contains(out.String(), "foo4\n")
}

func TestLogger_Write7(t *testing.T) {
	t.Parallel()
	assert := assert2.New(t)
	var out1 = bytes.NewBuffer(nil)
	var out2 = bytes.NewBuffer(nil)
	l := glog.New()
	l.SetOutput(glog.NewLoggerWriter(ioutil.Discard))
	l.AddWriter(glog.NewLoggerWriter(out1))
	l.AddWriter(glog.NewLoggerWriter(out2))
	l.SetLevel(gcore.LogLevelTrace)
	l.Tracef("foo1")
	l.Debugf("foo2")
	l.Infof("foo3")
	l.Warnf("foo4")
	assert.Contains(out1.String(), "foo1")
	assert.Contains(out1.String(), "foo2")
	assert.Contains(out1.String(), "foo2\n")
	assert.Contains(out1.String(), "foo4\n")

	assert.Contains(out2.String(), "foo1")
	assert.Contains(out2.String(), "foo2")
	assert.Contains(out2.String(), "foo2\n")
	assert.Contains(out2.String(), "foo4\n")
}
