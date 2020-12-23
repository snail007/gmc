package glog_test

import (
	"bytes"
	"github.com/snail007/gmc/core"
	log "github.com/snail007/gmc/gmc/log"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestNewGMCLog(t *testing.T) {
	assert := assert2.New(t)
	assert.Implements(new(gcore.Logger), log.NewGMCLog())
}

func TestGMCLog_SetOutput(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l.Info("a")
	assert.True(strings.HasSuffix(out.String(), "INFO a\n"))
}

func TestGMCLog_Writer(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	assert.Equal(&out, l.Writer())
}

func TestGMCLog_SetLevel(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l.SetLevel(gcore.LWARN)
	l.Info("a")
	assert.Empty(out.String())
}

func TestGMCLog_With_1(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l0 := l.With("api")
	l0.Info("a", "b")
	t.Log(out.String())
	assert.True(strings.HasSuffix(out.String(), "[api] INFO ab\n"))
	assert.Equal(l0.Namespace(), "api")
}

func TestGMCLog_With_2(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l0 := l.With("api").With("user").With("list")
	l0.Info("a")
	t.Log(out.String())
	assert.True(strings.HasSuffix(out.String(), "[api/user/list] INFO a\n"))
}

func TestGMCLog_Infof(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l.Infof("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "INFO a10\n"))
}

func TestGMCLog_Trace(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetLevel(gcore.LTRACE)
	l.SetOutput(&out)
	l.Trace("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "TRACE a\n"))
}

func TestGMCLog_Tracef(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetLevel(gcore.LTRACE)
	l.SetOutput(&out)
	l.Tracef("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "TRACE a10\n"))
}

func TestGMCLog_Debug(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l.Debug("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "DEBUG a\n"))
}

func TestGMCLog_Debugf(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l.Debugf("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "DEBUG a10\n"))
}

func TestGMCLog_Warn(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l.Warn("a")
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "WARN a\n"))
}

func TestGMCLog_Warnf(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	l := log.NewGMCLog()
	l.SetOutput(&out)
	l.Warnf("a%d", 10)
	t.Log(out.String(), len(out.String()))
	assert.True(strings.HasSuffix(out.String(), "WARN a10\n"))
}

func TestGMCLog_Panic(t *testing.T) {
	assert := assert2.New(t)
	l := log.NewGMCLog()
	assert.PanicsWithValue("[gmc]gmc/log/log_test.go:145: PANIC a", func() {
		l.Panic("a")
	})
}

func TestGMCLog_Panicf(t *testing.T) {
	assert := assert2.New(t)
	l := log.NewGMCLog()
	assert.PanicsWithValue("[gmc]gmc/log/log_test.go:153: PANIC a10", func() {
		l.Panicf("a%d", 10)
	})
}

func TestGMCLog_Error(t *testing.T) {
	assert := assert2.New(t)
	l := log.NewGMCLog()
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

func TestGMCLog_Errorf(t *testing.T) {
	assert := assert2.New(t)
	l := log.NewGMCLog()
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
