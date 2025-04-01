// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"fmt"
	"github.com/goccy/go-json"
	gcore "github.com/snail007/gmc/core"
	"os"
	"time"
)

type JSONLogger struct {
	logger *Logger
	data   map[string]interface{}
}

func (l *JSONLogger) setValue(key string, value interface{}) *JSONLogger {
	l.data[key] = value
	return l
}

func (l *JSONLogger) Uint8(key string, value uint8) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Uint32(key string, value uint32) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Uint(key string, value uint) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Uint64(key string, value uint64) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Int(key string, value int) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Int32(key string, value int32) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Int64(key string, value int64) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Float32(key string, value float32) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Float64(key string, value float64) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) String(key string, value string) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Bool(key string, value bool) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Time(key string, value time.Time) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Duration(key string, value time.Duration) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Any(key string, value interface{}) *JSONLogger {
	return l.setValue(key, value)
}

func (l *JSONLogger) Trace(msg string) string {
	return l.write(gcore.LogLevelTrace, msg)
}

func (l *JSONLogger) Tracef(format string, values ...interface{}) string {
	return l.write(gcore.LogLevelTrace, format, values...)
}

func (l *JSONLogger) Debug(msg string) string {
	return l.write(gcore.LogLeveDebug, msg)
}

func (l *JSONLogger) Debugf(format string, values ...interface{}) string {
	return l.write(gcore.LogLeveDebug, format, values...)
}

func (l *JSONLogger) Info(msg string) string {
	return l.write(gcore.LogLeveInfo, msg)
}

func (l *JSONLogger) Infof(format string, values ...interface{}) string {
	return l.write(gcore.LogLeveInfo, format, values...)
}

func (l *JSONLogger) Warn(msg string) string {
	return l.write(gcore.LogLeveWarn, msg)
}

func (l *JSONLogger) Warnf(format string, values ...interface{}) string {
	return l.write(gcore.LogLeveWarn, format, values...)
}

func (l *JSONLogger) Error(msg string) string {
	return l.write(gcore.LogLeveError, msg)
}

func (l *JSONLogger) Errorf(format string, values ...interface{}) string {
	return l.write(gcore.LogLeveError, format, values...)
}

func (l *JSONLogger) Panic(msg string) {
	panic(l.write(gcore.LogLevePanic, msg))
}

func (l *JSONLogger) Panicf(format string, values ...interface{}) {
	panic(l.write(gcore.LogLevePanic, format, values...))
}

func (l *JSONLogger) Fatal(msg string) {
	l.write(gcore.LogLeveFatal, msg)
	os.Exit(-1)
}

func (l *JSONLogger) Fatalf(format string, values ...interface{}) {
	l.write(gcore.LogLeveFatal, format, values...)
	os.Exit(-1)
}

func (l *JSONLogger) write(level gcore.LogLevel, format string, values ...interface{}) (jsonStr string) {
	msg := format
	if len(values) > 0 {
		msg = fmt.Sprintf(format, values...)
	}
	l.setValue("log_msg", msg)
	l.setValue("log_level", level.String())
	l.setValue("log_time", time.Now().Format("2006-01-02 15:04:05.000 -07"))
	b, _ := json.Marshal(l.data)
	jsonStr = string(append(b, '\n'))
	l.logger.WriteRaw(jsonStr, level)
	return
}
