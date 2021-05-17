// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"fmt"
	"github.com/snail007/gmc/core"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	logger = New("")
)

func init() {
	logger.SetCallerSkip(logger.CallerSkip()+1)
}

func Panic(v ...interface{}) {
	logger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

func Error(v ...interface{}) {
	logger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Warn(v ...interface{}) {
	logger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func Trace(v ...interface{}) {
	logger.Trace(v...)
}

func Tracef(format string, v ...interface{}) {
	logger.Tracef(format, v...)
}

func SetLevel(level gcore.LogLevel) {
	logger.SetLevel(level)
}

func With(name string) gcore.Logger {
	return logger.With(name)
}

func Namespace() string {
	return logger.Namespace()
}

func Writer() io.Writer {
	return logger.Writer()
}

func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

func Async() bool {
	return logger.Async()
}

func WaitAsyncDone() {
	logger.WaitAsyncDone()
}

func EnableAsync() {
	logger.EnableAsync()
}

type bufChnItem struct {
	level gcore.LogLevel
	msg   string
}
type Logger struct {
	l          *log.Logger
	parent     *Logger
	ns         string
	level      gcore.LogLevel
	async      bool
	asyncOnce  *sync.Once
	bufChn     chan bufChnItem
	asyncWG    *sync.WaitGroup
	callerSkip int
}

func (s *Logger) CallerSkip() int {
	return s.callerSkip
}

func (s *Logger) SetCallerSkip(callerSkip int) {
	s.callerSkip = callerSkip
}

func New(prefix ...string) gcore.Logger {
	pre := ""
	if len(prefix) == 1 {
		pre = prefix[0]
	}
	l := log.New(os.Stdout, pre, log.LstdFlags|log.Lmicroseconds)
	return &Logger{
		l:          l,
		level:      gcore.LDEBUG,
		asyncOnce:  &sync.Once{},
		callerSkip: 2,
	}
}

func (s *Logger) WaitAsyncDone() {
	if s.async {
		s.asyncWG.Wait()
	}
}

func (s *Logger) Async() bool {
	return s.async
}

func (s *Logger) asyncWriterInit() {
	s.bufChn = make(chan bufChnItem, 2048)
	s.asyncWG = &sync.WaitGroup{}
	go func() {
		for {
			item := <-s.bufChn
			s.output(item.msg)
			func() {
				defer func() {
					_ = recover()
				}()
				s.asyncWG.Done()
			}()
		}
	}()
}

func (s *Logger) EnableAsync() {
	s.async = true
	s.asyncOnce.Do(func() {
		s.asyncWriterInit()
	})
}

func (s *Logger) Level() gcore.LogLevel {
	return s.level
}

func (s *Logger) SetLevel(i gcore.LogLevel) {
	s.level = i
}

func (s *Logger) With(namespace string) gcore.Logger {
	return &Logger{
		l:      s.l,
		parent: s,
		ns:     namespace,
		level:  s.level,
	}
}

func (s *Logger) Namespace() string {
	ns := s.ns
	if s.parent != nil {
		ns = s.parent.Namespace() + "/" + s.ns
	}
	return strings.TrimLeft(ns, "/")
}

func (s *Logger) namespace() string {
	if s.parent != nil {
		return "[" + s.Namespace() + "] "
	}
	return ""
}

func (s *Logger) Panicf(format string, v ...interface{}) {
	if s.level > gcore.LPANIC {
		return
	}
	str := s.caller(fmt.Sprintf(s.namespace()+"PANIC "+format, v...), s.skip())
	s.Write(str)
	s.WaitAsyncDone()
	panic(str)
}

func (s *Logger) Panic(v ...interface{}) {
	if s.level > gcore.LPANIC {
		return
	}
	v0 := []interface{}{s.namespace() + "PANIC "}
	str := s.caller(fmt.Sprint(append(v0, v...)...), s.skip())
	s.Write(str)
	s.WaitAsyncDone()
	panic(str)
}

func (s *Logger) Errorf(format string, v ...interface{}) {
	if s.level > gcore.LERROR {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"ERROR "+format, v...), s.skip()))
	s.WaitAsyncDone()
	os.Exit(1)
}

func (s *Logger) Error(v ...interface{}) {
	if s.level > gcore.LERROR {
		return
	}
	v0 := []interface{}{s.namespace() + "ERROR "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...), s.skip()))
	s.WaitAsyncDone()
	os.Exit(1)
}

func (s *Logger) Warnf(format string, v ...interface{}) {
	if s.level > gcore.LWARN {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"WARN "+format, v...), s.skip()))
}

func (s *Logger) Warn(v ...interface{}) {
	if s.level > gcore.LWARN {
		return
	}
	v0 := []interface{}{s.namespace() + "WARN "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...), s.skip()))
}

func (s *Logger) Infof(format string, v ...interface{}) {
	if s.level > gcore.LINFO {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"INFO "+format, v...), s.skip()))

}

func (s *Logger) Info(v ...interface{}) {
	if s.level > gcore.LINFO {
		return
	}
	v0 := []interface{}{s.namespace() + "INFO "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...), s.skip()))
}

func (s *Logger) Debugf(format string, v ...interface{}) {
	if s.level > gcore.LDEBUG {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"DEBUG "+format, v...), s.skip()))
}

func (s *Logger) Debug(v ...interface{}) {
	if s.level > gcore.LDEBUG {
		return
	}
	v0 := []interface{}{s.namespace() + "DEBUG "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...), s.skip()))
}

func (s *Logger) Tracef(format string, v ...interface{}) {
	if s.level > gcore.LTRACE {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"TRACE "+format, v...), s.skip()))
}

func (s *Logger) Trace(v ...interface{}) {
	if s.level > gcore.LTRACE {
		return
	}
	v0 := []interface{}{s.namespace() + "TRACE "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...), s.skip()))
}

func (s *Logger) Writer() io.Writer {
	return s.l.Writer()
}

func (s *Logger) SetOutput(w io.Writer) {
	s.l.SetOutput(w)
}

func (s *Logger) SetFlags(f int) {
	s.l.SetFlags(f)
}

func (s *Logger) Write(str string) {
	if s.async {
		select {
		case s.bufChn <- bufChnItem{
			msg: str,
		}:
			s.asyncWG.Add(1)
		default:
			s.output("WARN gmclog buf chan overflow")
		}
		return
	}
	s.output(str)
}

func (s *Logger) output(str string) {
	s.l.Print(str)
}

func (s *Logger) skip() int {
	skip := s.callerSkip
	return skip
}

func (s *Logger) caller(msg string, skip int) string {
	file := "unknown"
	line := 0
	if _, file0, line0, ok := runtime.Caller(skip); ok {
		file0 = strings.Replace(file0, "\\", "/", -1)
		p := "github.com/snail007/gmc"
		if strings.Contains(file0, p) &&
			!strings.Contains(file0, p+"demos") {
			file = "[gmc]" + file0[strings.Index(file0, p)+len(p):]
		} else {
			file = filepath.Base(filepath.Dir(file0)) + "/" + filepath.Base(file0)
		}
		line = line0
	}
	msg = fmt.Sprintf("%s:%d: ", file, line) + msg
	return msg
}
