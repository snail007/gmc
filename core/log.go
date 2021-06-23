// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcore

import (
	"io"
	"time"
)

const (
	LogLevelTrace LogLevel = iota + 1
	LogLeveDebug
	LogLeveInfo
	LogLeveWarn
	LogLeveError
	LogLevePanic
	LogLeveNone
)
const (
	LogFlagNormal LogFlag = iota + 1
	LogFlagShort
	LogFlagLong
)

type LogLevel int
type LogFlag int

type Logger interface {
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Trace(v ...interface{})
	Tracef(format string, v ...interface{})

	Level() LogLevel
	SetLevel(LogLevel)

	With(name string) Logger
	Namespace() string

	Writer() io.Writer
	AddWriter(writer io.Writer) Logger
	SetOutput(w io.Writer)
	SetFlag(f LogFlag)

	Async() bool
	WaitAsyncDone()
	EnableAsync()

	CallerSkip() int
	SetCallerSkip(callerSkip int)
	Write(string)

	ExitCode() int
	SetExitCode(exitCode int)
	ExitFunc() func(int)
	SetExitFunc(exitFunc func(int))

	WithRate(duration time.Duration) Logger
	SetRateCallback(cb func(msg string)) Logger
}
