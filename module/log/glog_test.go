// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog_test

import (
	"bytes"
	"fmt"
	"github.com/snail007/gmc"
	gctx "github.com/snail007/gmc/module/ctx"
	ghttppprof "github.com/snail007/gmc/util/pprof"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/snail007/gmc/core"
	glog "github.com/snail007/gmc/module/log"
	_ "github.com/snail007/gmc/using/basic"
	gtest "github.com/snail007/gmc/util/testing"
	assert2 "github.com/stretchr/testify/assert"
)

func TestGlog(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	glog.SetOutput(&out)
	tests := []struct {
		name    string
		prepare func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool)
		test    func(args []interface{}) (out string, contains []string)
	}{
		{"glog.Trace", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveDebug)
			glog.Trace("b")
			glog.SetLevel(gcore.LogLevelTrace)
			out.Reset()
			glog.Trace("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"TRACE a"}
		}},
		{"glog.Tracef", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveDebug)
			glog.Tracef("b")
			glog.SetLevel(gcore.LogLevelTrace)
			out.Reset()
			glog.Tracef("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"TRACE a"}
		}},
		{"glog.Debug", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveInfo)
			glog.Debug("b")
			glog.SetLevel(gcore.LogLeveDebug)
			out.Reset()
			glog.Debug("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"DEBUG a"}
		}},
		{"glog.Debugf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveInfo)
			glog.Debugf("b")
			glog.SetLevel(gcore.LogLeveDebug)
			out.Reset()
			glog.Debugf("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"DEBUG a"}
		}},
		{"glog.Info", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveWarn)
			glog.Info("b")
			glog.SetLevel(gcore.LogLeveInfo)
			out.Reset()
			glog.Info("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"INFO a"}
		}},
		{"glog.Infof", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveWarn)
			glog.Infof("b")
			glog.SetLevel(gcore.LogLeveInfo)
			out.Reset()
			glog.Infof("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"INFO a"}
		}},
		{"glog.Warn", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveError)
			glog.Warn("b")
			glog.SetLevel(gcore.LogLeveWarn)
			out.Reset()
			glog.Warn("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"WARN a"}
		}},
		{"glog.Warnf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveError)
			glog.Warnf("b")
			glog.SetLevel(gcore.LogLeveWarn)
			out.Reset()
			glog.Warnf("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"WARN a"}
		}},
		{"glog.Error", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLevePanic)
			glog.Warnf("b")
			glog.SetLevel(gcore.LogLeveError)
			out.Reset()
			glog.Error("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"ERROR a"}
		}},
		{"glog.Errorf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLevePanic)
			glog.Warnf("b")
			glog.SetLevel(gcore.LogLeveError)
			out.Reset()
			glog.Errorf("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"ERROR a"}
		}},
		{"glog.Fatal", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveFatal)
			glog.SetExitFunc(func(i int) {})
			out.Reset()
			args = []interface{}{&out}
			glog.Fatal("a")
			return
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"FATAL a"}
		}},
		{"glog.Fatalf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveFatal)
			glog.SetExitFunc(func(i int) {})
			out.Reset()
			args = []interface{}{&out}
			glog.Fatalf("%s", "a")
			return
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"FATAL a"}
		}},
		{"glog.Panic", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveNone)
			glog.Panic("b")
			glog.SetLevel(gcore.LogLevePanic)
			out.Reset()
			args = []interface{}{&out}
			defer func() {
				assert.NotNil(recover())
			}()
			glog.Panic("a")
			return
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"PANIC a"}
		}},
		{"glog.Panicf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LogLeveNone)
			glog.Panicf("b")
			glog.SetLevel(gcore.LogLevePanic)
			out.Reset()
			args = []interface{}{&out}
			defer func() {
				assert.NotNil(recover())
			}()
			glog.Panicf("%s", "a")
			return
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"PANIC a"}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert = assert2.New(t)
			args, stop := tt.prepare(t, assert)
			if stop {
				return
			}
			out, want := tt.test(args)
			for _, w := range want {
				assert.Contains(out, w)
				assert.Contains(out, "glog_test.go:")
			}
		})
	}
	testGlog_SetFlags(t, assert)
	testGlog_Level(t, assert)
}

func testGlog_SetFlags(t *testing.T, assert *assert2.Assertions) {
	glog.SetLevel(gcore.LogLevelTrace)
	var out bytes.Buffer
	glog.SetFlag(gcore.LogFlagShort)
	glog.SetOutput(&out)
	glog.Info("a")
	assert.Contains(out.String(), "/module/log/glog_test.go:")
	assert.Contains(out.String(), "INFO a")
}

func testGlog_Level(t *testing.T, assert *assert2.Assertions) {
	glog.SetLevel(gcore.LogLeveError)
	assert.Equal(glog.Level(), gcore.LogLeveError)
}

func TestGlog_Namespace(t *testing.T) {
	assert := assert2.New(t)
	log := glog.With("user")
	assert.Contains(log.Namespace(), "user")
	assert.Empty(glog.Namespace())
}

func TestGlog_Writer(t *testing.T) {
	assert := assert2.New(t)
	assert.Implements((*io.Writer)(nil), glog.Writer())
}

func TestGlog_Write(t *testing.T) {
	if gtest.RunProcess(t, func() {
		glog.Write("abc")
	}) {
		return
	}
	assert := assert2.New(t)
	out, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
	assert.Contains(out, "abc")
}

func TestGlog_Async(t *testing.T) {
	assert := assert2.New(t)
	if gtest.RunProcess(t, func() {
		glog.SetOutput(os.Stdout)
		glog.EnableAsync()
		assert.True(glog.Async())
		assert.Equal(1, glog.ExitCode())
		glog.Info("a")
		glog.WaitAsyncDone()
	}) {
		return
	}
	out, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
	assert.Contains(string(out), "INFO a")
}

func TestGlog_Normal(t *testing.T) {
	assert := assert2.New(t)
	log := glog.New()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlag(gcore.LogFlagNormal)
	log.Info("abc")
	assert.Contains(buf.String(), "abc")
	assert.NotContains(buf.String(), ".go")
}

func TestGlog_Short(t *testing.T) {
	assert := assert2.New(t)
	os.Setenv("LOG_SKIP_CHECK_GMC", "yes")
	log := glog.New()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlag(gcore.LogFlagShort)
	log.Info("abc")
	assert.Contains(buf.String(), "abc")
	assert.Contains(buf.String(), ".go")
}

func TestGlog_Long(t *testing.T) {
	os.Setenv("LOG_SKIP_CHECK_GMC", "yes")
	assert := assert2.New(t)
	log := glog.New()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlag(gcore.LogFlagLong)
	log.Info("abc")
	assert.Contains(buf.String(), "abc")
	assert.Contains(buf.String(), ".go")
}

func TestGlog_WithRate(t *testing.T) {
	t.Parallel()
	glog.SetRateCallback(func(msg string) {})
	l0 := glog.WithRate(time.Second)
	glog.AddWriter(ioutil.Discard)
	cnt := new(int32)
	l0.SetRateCallback(func(msg string) {
		atomic.AddInt32(cnt, 1)
	})
	for i := 0; i < 35; i++ {
		l0.Write("hello")
		time.Sleep(time.Millisecond * 100)
	}
	assert2.True(t, atomic.LoadInt32(cnt) >= 3)
}
func StartWebPProf(log gcore.Logger) {
	s := "pprof server"
	api, err := gmc.New.APIServer(gctx.NewCtx(), ":8809")
	if err != nil {
		log.Panicf("start %s error: %s", s, err)
	}
	ghttppprof.BindRouter(api.Router(), "/proxy/.pprof/")
	err = api.Run()
	if err != nil {
		log.Panicf("%s run error: %s", s, err)
	}
	//log.Infof("%s on: %s", s, api.Listener().Addr())
}
func TestGLog_With(t *testing.T) {
	if gtest.RunProcess(t, func() {
		l := glog.With("api")
		l.Info("abc")
	}) {
		return
	}
	assert := assert2.New(t)
	out, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
	t.Log(out)
	assert.True(strings.Contains(out, "[api] INFO abc\n"))
	assert.Contains(out, "glog_test.go")
}

func TestGLog_AddLevelWriter(t *testing.T) {
	if gtest.RunProcess(t, func() {
		var out bytes.Buffer
		l := glog.New()
		l.AddLevelWriter(&out, gcore.LogLeveInfo)
		l.SetLevel(gcore.LogLeveNone)
		l.Tracef("foo1")
		l.Info("foo2")
		fmt.Println(out.String())
	}) {
		return
	}
	assert := assert2.New(t)
	out, _, err := gtest.NewProcess(t).Wait()
	assert.Nil(err)
	t.Log(">>>>", out)
	assert.False(strings.Contains(out, "TRACE foo1\n"))
	assert.True(strings.Contains(out, "INFO foo2\n"))
	assert.Contains(out, "glog_test.go")
}
