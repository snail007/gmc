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
	"log"
	"os"
	"testing"
)

func TestGlog(t *testing.T) {
	assert := assert2.New(t)
	var out bytes.Buffer
	glog.SetLevel(gcore.LTRACE)
	glog.SetOutput(&out)
	tests := []struct {
		name    string
		prepare func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool)
		test    func(args []interface{}) (out string, contains []string)
	}{
		{"glog.Trace", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			out.Reset()
			glog.Trace("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"TRACE a"}
		}},
		{"glog.Tracef", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			out.Reset()
			glog.Tracef("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"TRACE a"}
		}},
		{"glog.Debug", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			out.Reset()
			glog.Debug("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"DEBUG a"}
		}},
		{"glog.Debugf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			out.Reset()
			glog.Debugf("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"DEBUG a"}
		}},
		{"glog.Info", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			out.Reset()
			glog.Info("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"INFO a"}
		}},
		{"glog.Infof", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			out.Reset()
			glog.Infof("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"INFO a"}
		}},
		{"glog.Warn", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			out.Reset()
			glog.Warn("a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"WARN a"}
		}},
		{"glog.Warnf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			out.Reset()
			glog.Warnf("%s", "a")
			return []interface{}{&out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return (args[0].(*bytes.Buffer)).String(), []string{"WARN a"}
		}},
		{"glog.Error", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			if os.Getenv("ASSERT_EXISTS_"+t.Name()) == "1" {
				glog.SetExitCode(-1)
				glog.SetOutput(os.Stdout)
				glog.Error("a")
				return nil, true
			}
			os.Setenv("ASSERT_EXISTS_"+t.Name(), "1")
			out, err := gtest.ExecTestFunc(t.Name())
			assert.Nil(err)
			//assert.Contains(err.Error(), "exit status 1")
			return []interface{}{out}, false
		}, func(args []interface{}) (out string, contains []string) {
			return args[0].(string), []string{"ERROR a"}
		}},
		{"glog.Errorf", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
			if os.Getenv("ASSERT_EXISTS_"+t.Name()) == "1" {
				glog.SetExitCode(-1)
				glog.SetOutput(os.Stdout)
				glog.Errorf("%s", "a")
				return nil, true
			}
			os.Setenv("ASSERT_EXISTS_"+t.Name(), "1")
			out, err := gtest.ExecTestFunc(t.Name())
			assert.Nil(err)
			//assert.Contains(err.Error(), "exit status 1")
			return []interface{}{string(out)}, false
		}, func(args []interface{}) (out string, contains []string) {
			return args[0].(string), []string{"ERROR a"}
		}},
		{"glog.Panic", func(t *testing.T, assert *assert2.Assertions) (args []interface{}, stop bool) {
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
	var out bytes.Buffer
	glog.SetFlags(log.LstdFlags | log.Llongfile)
	glog.SetOutput(&out)
	glog.Info("a")
	assert.Contains(out.String(), "gmc/module/log/log.go:")
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
	if os.Getenv("ASSERT_EXISTS_"+t.Name()) == "1" {
		glog.Write("abc")
		return
	}
	assert := assert2.New(t)
	os.Setenv("ASSERT_EXISTS_"+t.Name(), "1")
	//os.Setenv("GMCT_COVER_VERBOSE","true")
	out, err := gtest.ExecTestFunc(t.Name())
	assert.Nil(err)
	assert.Contains(out, "abc")
}

func TestGlog_Async(t *testing.T) {
	assert := assert2.New(t)
	if os.Getenv("ASSERT_EXISTS_"+t.Name()) == "1" {
		glog.SetOutput(os.Stdout)
		glog.EnableAsync()
		assert.True(glog.Async())
		assert.Equal(1, glog.ExitCode())
		glog.Info("a")
		glog.WaitAsyncDone()
		return
	}
	os.Setenv("ASSERT_EXISTS_"+t.Name(), "1")
	out, err := gtest.ExecTestFunc(t.Name())
	assert.Nil(err)
	assert.Contains(string(out), "INFO a")
}
