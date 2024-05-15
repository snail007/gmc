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
	LogLeveFatal
	LogLeveNone
)

type LogLevel int

func (s LogLevel)String() string {
	switch s {
	case LogLevelTrace:
		return "TRACE"
	case LogLeveDebug:
		return "DEBUG"
	case LogLeveInfo:
		return "INFO"
	case LogLeveWarn:
		return "WARN"
	case LogLeveError:
		return "ERROR"
	case LogLevePanic:
		return "PANIC"
	case LogLeveFatal:
		return "FATAL"
	case LogLeveNone:
		return "NONE"
	}
	return ""
}

const (
	LogFlagNormal LogFlag = iota + 1
	LogFlagShort
	LogFlagLong
)

type LogFlag int

type Logger interface {
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

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

	Writer() LoggerWriter
	AddWriter(LoggerWriter) Logger
	AddLevelWriter(io.Writer, LogLevel) Logger
	AddLevelsWriter(io.Writer, ...LogLevel) Logger
	SetOutput(LoggerWriter)
	SetFlag(f LogFlag)

	Async() bool
	WaitAsyncDone()
	EnableAsync()

	CallerSkip() int
	SetCallerSkip(callerSkip int)
	Write(string, LogLevel)
	WriteRaw(string, LogLevel)

	ExitCode() int
	SetExitCode(exitCode int)
	ExitFunc() func(int)
	SetExitFunc(exitFunc func(int))

	WithRate(duration time.Duration) Logger
	SetRateCallback(cb func(msg string)) Logger
	SetTimeLayout(layout string)
	SetAsyncBufferSize(asyncBufferSize int)
	SetErrHandler(errHandler func(error))
}

type LoggerWriter interface {
	Write(p []byte, level LogLevel) (n int, err error)
}
