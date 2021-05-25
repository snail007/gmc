// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog_test

import (
	"bytes"
	"github.com/snail007/gmc/core"
	glog "github.com/snail007/gmc/module/log"
	_ "github.com/snail007/gmc/using/basic"
	gtest "github.com/snail007/gmc/util/testing"
	assert2 "github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
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
			glog.SetLevel(gcore.LDEBUG)
			glog.Trace("b")
			glog.SetLevel(gcore.LTRACE)
			out.Reset()
			glog.Trace("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"TRACE a"}
		}},
		{"glog.Tracef", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LDEBUG)
			glog.Tracef("b")
			glog.SetLevel(gcore.LTRACE)
			out.Reset()
			glog.Tracef("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"TRACE a"}
		}},
		{"glog.Debug", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LINFO)
			glog.Debug("b")
			glog.SetLevel(gcore.LDEBUG)
			out.Reset()
			glog.Debug("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"DEBUG a"}
		}},
		{"glog.Debugf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LINFO)
			glog.Debugf("b")
			glog.SetLevel(gcore.LDEBUG)
			out.Reset()
			glog.Debugf("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"DEBUG a"}
		}},
		{"glog.Info", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LWARN)
			glog.Info("b")
			glog.SetLevel(gcore.LINFO)
			out.Reset()
			glog.Info("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"INFO a"}
		}},
		{"glog.Infof", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LWARN)
			glog.Infof("b")
			glog.SetLevel(gcore.LINFO)
			out.Reset()
			glog.Infof("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"INFO a"}
		}},
		{"glog.Warn", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LERROR)
			glog.Warn("b")
			glog.SetLevel(gcore.LWARN)
			out.Reset()
			glog.Warn("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"WARN a"}
		}},
		{"glog.Warnf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LERROR)
			glog.Warnf("b")
			glog.SetLevel(gcore.LWARN)
			out.Reset()
			glog.Warnf("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"WARN a"}
		}},
		{"glog.Error", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			oldExit := glog.ExitFunc()
			glog.SetLevel(gcore.LPANIC)
			glog.Error("b")
			glog.SetLevel(gcore.LERROR)
			out.Reset()
			gotCode := 0
			glog.SetExitCode(2)
			glog.SetExitFunc(func(code int) {
				gotCode = code
			})
			glog.Error("a")
			assert.Equal(2, gotCode)
			glog.SetExitFunc(oldExit)
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"ERROR a"}
		}},
		{"glog.Errorf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			oldExit := glog.ExitFunc()
			glog.SetLevel(gcore.LPANIC)
			glog.Errorf("b")
			glog.SetLevel(gcore.LERROR)
			out.Reset()
			gotCode := 0
			glog.SetExitCode(3)
			glog.SetExitFunc(func(code int) {
				gotCode = code
			})
			glog.Errorf("%s", "a")
			assert.Equal(3, gotCode)
			glog.SetExitFunc(oldExit)
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"ERROR a"}
		}},
		{"glog.Panic", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			glog.SetLevel(gcore.LNONE)
			glog.Panic("b")
			glog.SetLevel(gcore.LPANIC)
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
			glog.SetLevel(gcore.LNONE)
			glog.Panicf("b")
			glog.SetLevel(gcore.LPANIC)
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
	glog.SetLevel(gcore.LTRACE)
	var out bytes.Buffer
	glog.SetFlag(gcore.LFLAG_SHORT)
	glog.SetOutput(&out)
	glog.Info("a")
	assert.Contains(out.String(), "/module/log/glog_test.go:")
	assert.Contains(out.String(), "INFO a")
}

func testGlog_Level(t *testing.T, assert *assert2.Assertions) {
	glog.SetLevel(gcore.LERROR)
	assert.Equal(glog.Level(), gcore.LERROR)
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
	log.SetFlag(gcore.LFLAG_NORMAL)
	log.Info("abc")
	assert.Contains(buf.String(), "abc")
	assert.NotContains(buf.String(), ".go")
}

func TestGlog_Short(t *testing.T) {
	assert := assert2.New(t)
	os.Setenv("LOG_SKIP_CHECK_GMC","yes")
	log := glog.New()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlag(gcore.LFLAG_SHORT)
	log.Info("abc")
	assert.Contains(buf.String(), "abc")
	assert.Contains(buf.String(), ".go")
}

func TestGlog_Long(t *testing.T) {
	os.Setenv("LOG_SKIP_CHECK_GMC","yes")
	assert := assert2.New(t)
	log := glog.New()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlag(gcore.LFLAG_LONG)
	log.Info("abc")
	assert.Contains(buf.String(), "abc")
	assert.Contains(buf.String(), ".go")
}