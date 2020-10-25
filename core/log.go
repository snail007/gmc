package gmccore

import "io"

const (
	LTRACE=iota
	LDEBUG
	LINFO
	LWARN
	LERROR
	LPANIC
)

type LOG_LEVEL int

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

	SetLevel(LOG_LEVEL)

	With(name string) Logger
	Namespace() string

	Writer() io.Writer
	SetOutput(w io.Writer)

}
