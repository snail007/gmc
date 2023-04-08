// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"errors"
	"fmt"
	"github.com/snail007/gmc/util/gpool"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/snail007/gmc/core"
	"golang.org/x/time/rate"
)

var (
	logger = New("")
	pool   = gpool.NewWithLogger(10, nil)
)

func init() {
	logger.SetCallerSkip(logger.CallerSkip() + 1)
}

func Panic(v ...interface{}) {
	logger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
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
	l := logger.With(name)
	l.SetCallerSkip(l.CallerSkip() - 1)
	return l
}

func WithRate(duration time.Duration) gcore.Logger {
	return logger.WithRate(duration)
}

func AddWriter(writer io.Writer) gcore.Logger {
	logger.AddWriter(writer)
	return logger
}

func AddLevelWriter(writer io.Writer, level gcore.LogLevel) gcore.Logger {
	logger.AddLevelWriter(writer, level)
	return logger
}

func SetRateCallback(cb func(msg string)) gcore.Logger {
	logger.SetRateCallback(cb)
	return logger
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

func SetFlag(f gcore.LogFlag) {
	logger.SetFlag(f)
}

func Write(s string) {
	logger.Write(s)
}

func Level() gcore.LogLevel {
	return logger.Level()
}

func SetExitCode(code int) {
	logger.SetExitCode(code)
}

func ExitCode() int {
	return logger.ExitCode()
}

func ExitFunc() func(int) {
	return logger.ExitFunc()
}

func SetExitFunc(exitFunc func(int)) {
	logger.SetExitFunc(exitFunc)
}

type bufChnItem struct {
	writer *levelWriter
	level  gcore.LogLevel
	msg    string
}
type Logger struct {
	l           *log.Logger
	parent      *Logger
	ns          string
	level       gcore.LogLevel
	async       bool
	asyncOnce   *sync.Once
	bufChn      chan bufChnItem
	asyncWG     *sync.WaitGroup
	callerSkip  int
	exitCode    int
	exitFunc    func(int)
	flag        gcore.LogFlag
	lim         *rate.Limiter
	limCallback func(msg string)
	// for testing purpose
	skipCheckGMC bool
	levelWriters []*levelWriter
}

func (s *Logger) clone() *Logger {
	return &Logger{
		l:            s.l,
		parent:       s.parent,
		ns:           s.ns,
		level:        s.level,
		async:        s.async,
		asyncOnce:    s.asyncOnce,
		bufChn:       s.bufChn,
		asyncWG:      s.asyncWG,
		callerSkip:   s.callerSkip,
		exitCode:     s.exitCode,
		exitFunc:     s.exitFunc,
		flag:         s.flag,
		lim:          s.lim,
		skipCheckGMC: s.skipCheckGMC,
		levelWriters: s.levelWriters,
	}
}

func (s *Logger) exit() {
	if s.exitCode >= 0 {
		s.exitFunc(s.exitCode)
	}
}

func (s *Logger) ExitFunc() func(int) {
	return s.exitFunc
}

func (s *Logger) SetRateCallback(cb func(msg string)) gcore.Logger {
	s.limCallback = cb
	return s
}

func (s *Logger) AddLevelWriter(w io.Writer, level gcore.LogLevel) gcore.Logger {
	s.levelWriters = append(s.levelWriters, newLogWriter(w, level))
	return s
}

func (s *Logger) AddWriter(w io.Writer) gcore.Logger {
	s.l.SetOutput(newMultiWriter(s.l.Writer(), w))
	return s
}

func (s *Logger) SetExitFunc(exitFunc func(int)) {
	s.exitFunc = exitFunc
}

func (s *Logger) ExitCode() int {
	return s.exitCode
}

func (s *Logger) SetExitCode(exitCode int) {
	s.exitCode = exitCode
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
		l:            l,
		level:        gcore.LogLeveDebug,
		asyncOnce:    &sync.Once{},
		callerSkip:   2,
		exitCode:     1,
		exitFunc:     os.Exit,
		flag:         gcore.LogFlagShort,
		skipCheckGMC: os.Getenv("LOG_SKIP_CHECK_GMC") == "yes",
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
			s.output(item.msg, item.writer)
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
	l := s.clone()
	l.ns = namespace
	l.parent = s
	return l
}

func (s *Logger) WithRate(duration time.Duration) gcore.Logger {
	r := rate.NewLimiter(rate.Every(duration), 1)
	l := s.clone()
	l.lim = r
	return l
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

func (s *Logger) levelWrite(str string, level gcore.LogLevel) {
	if len(s.levelWriters) == 0 {
		return
	}
	g := sync.WaitGroup{}
	g.Add(len(s.levelWriters))
	for _, w := range s.levelWriters {
		if w.level <= level {
			w0 := w
			pool.Submit(func() {
				defer g.Done()
				s.write(str, w0)
			})
		}
	}
	g.Wait()
}

func (s *Logger) canLevelWrite(level gcore.LogLevel) bool {
	for _, w := range s.levelWriters {
		if w.level <= level {
			return true
		}
	}
	return false
}

func (s *Logger) Panicf(format string, v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLevePanic)
	if s.level > gcore.LogLevePanic && !levelWrite {
		return
	}
	str := s.caller(fmt.Sprintf(s.namespace()+"PANIC "+format, v...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLevePanic)
	}

	if s.level <= gcore.LogLevePanic {
		s.write(str, nil)
	}
	s.WaitAsyncDone()
	panic(str)
}

func (s *Logger) Panic(v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLevePanic)
	if s.level > gcore.LogLevePanic && !levelWrite {
		return
	}
	v0 := []interface{}{s.namespace() + "PANIC "}
	str := s.caller(fmt.Sprint(append(v0, v...)...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLevePanic)
	}

	if s.level <= gcore.LogLevePanic {
		s.write(str, nil)
	}
	s.WaitAsyncDone()
	panic(str)
}

func (s *Logger) Fatalf(format string, v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveFatal)
	if s.level > gcore.LogLeveFatal && !levelWrite {
		return
	}
	str := s.caller(fmt.Sprintf(s.namespace()+"FATAL "+format, v...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveFatal)
	}

	if s.level <= gcore.LogLeveFatal {
		s.write(str, nil)
	}
	s.WaitAsyncDone()
	s.exit()
}

func (s *Logger) Fatal(v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveFatal)
	if s.level > gcore.LogLeveFatal && !levelWrite {
		return
	}
	v0 := []interface{}{s.namespace() + "FATAL "}
	str := s.caller(fmt.Sprint(append(v0, v...)...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveFatal)
	}

	if s.level <= gcore.LogLeveFatal {
		s.write(str, nil)
	}
	s.WaitAsyncDone()
	s.exit()
}

func (s *Logger) Errorf(format string, v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveError)
	if s.level > gcore.LogLeveError && !levelWrite {
		return
	}

	str := s.caller(fmt.Sprintf(s.namespace()+"ERROR "+format, v...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveError)
	}

	if s.level <= gcore.LogLeveError {
		s.write(str, nil)
	}
}

func (s *Logger) Error(v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveError)
	if s.level > gcore.LogLeveError && !levelWrite {
		return
	}

	v0 := []interface{}{s.namespace() + "ERROR "}
	str := s.caller(fmt.Sprint(append(v0, v...)...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveError)
	}

	if s.level <= gcore.LogLeveError {
		s.write(str, nil)
	}
}

func (s *Logger) Warnf(format string, v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveWarn)
	if s.level > gcore.LogLeveWarn && !levelWrite {
		return
	}
	str := s.caller(fmt.Sprintf(s.namespace()+"WARN "+format, v...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveWarn)
	}

	if s.level <= gcore.LogLeveWarn {
		s.write(str, nil)
	}
}

func (s *Logger) Warn(v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveWarn)
	if s.level > gcore.LogLeveWarn && !levelWrite {
		return
	}
	v0 := []interface{}{s.namespace() + "WARN "}
	str := s.caller(fmt.Sprint(append(v0, v...)...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveWarn)
	}

	if s.level <= gcore.LogLeveWarn {
		s.write(str, nil)
	}
}

func (s *Logger) Infof(format string, v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveInfo)
	if s.level > gcore.LogLeveInfo && !levelWrite {
		return
	}
	str := s.caller(fmt.Sprintf(s.namespace()+"INFO "+format, v...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveInfo)
	}

	if s.level <= gcore.LogLeveInfo {
		s.write(str, nil)
	}
}

func (s *Logger) Info(v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveInfo)
	if s.level > gcore.LogLeveInfo && !levelWrite {
		return
	}

	v0 := []interface{}{s.namespace() + "INFO "}
	str := s.caller(fmt.Sprint(append(v0, v...)...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveInfo)
	}

	if s.level <= gcore.LogLeveInfo {
		s.write(str, nil)
	}
}

func (s *Logger) Debugf(format string, v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveDebug)
	if s.level > gcore.LogLeveDebug && !levelWrite {
		return
	}
	str := s.caller(fmt.Sprintf(s.namespace()+"DEBUG "+format, v...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveDebug)
	}

	if s.level <= gcore.LogLeveDebug {
		s.write(str, nil)
	}
}

func (s *Logger) Debug(v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLeveDebug)
	if s.level > gcore.LogLeveDebug && !levelWrite {
		return
	}
	v0 := []interface{}{s.namespace() + "DEBUG "}
	str := s.caller(fmt.Sprint(append(v0, v...)...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLeveDebug)
	}

	if s.level <= gcore.LogLeveDebug {
		s.write(str, nil)
	}
}

func (s *Logger) Tracef(format string, v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLevelTrace)
	if s.level > gcore.LogLevelTrace && !levelWrite {
		return
	}
	str := s.caller(fmt.Sprintf(s.namespace()+"TRACE "+format, v...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLevelTrace)
	}

	if s.level <= gcore.LogLevelTrace {
		s.write(str, nil)
	}
}

func (s *Logger) Trace(v ...interface{}) {
	levelWrite := s.canLevelWrite(gcore.LogLevelTrace)
	if s.level > gcore.LogLevelTrace && !levelWrite {
		return
	}
	v0 := []interface{}{s.namespace() + "TRACE "}
	str := s.caller(fmt.Sprint(append(v0, v...)...), s.skip())

	if levelWrite {
		s.levelWrite(str, gcore.LogLevelTrace)
	}

	if s.level <= gcore.LogLevelTrace {
		s.write(str, nil)
	}
}

func (s *Logger) Writer() io.Writer {
	return s.l.Writer()
}

func (s *Logger) SetOutput(w io.Writer) {
	s.l.SetOutput(w)
}

func (s *Logger) SetFlag(f gcore.LogFlag) {
	s.flag = f
}

func (s *Logger) write(str string, writer *levelWriter) {
	if s.async {
		select {
		case s.bufChn <- bufChnItem{
			msg:    str,
			writer: writer,
		}:
			s.asyncWG.Add(1)
		default:
			s.output("WARN gmclog buf chan overflow", writer)
		}
		return
	}
	s.output(str, writer)
}

func (s *Logger) Write(msg string) {
	s.write(s.caller(msg, s.skip()), nil)
}

func (s *Logger) WriteRaw(msg string) {
	s.write(msg, nil)
}

func (s *Logger) output(str string, writer *levelWriter) {
	if s.lim != nil && !s.lim.Allow() {
		return
	}
	if s.lim != nil && s.limCallback != nil {
		go s.limCallback(str)
	}
	if writer != nil {
		writer.lock.Lock()
		defer writer.lock.Unlock()
		ln := ""
		if len(str) == 0 || str[len(str)-1] != '\n' {
			ln = "\n"
		}
		writer.Write([]byte(time.Now().Format("2006/01/02 15:04:05.000000") + " " + str + ln))
	} else {
		s.l.Print(str)
	}
}

func (s *Logger) skip() int {
	skip := s.callerSkip
	return skip
}

func (s *Logger) caller(msg string, skip int) string {
	if s.flag == gcore.LogFlagNormal {
		return msg
	}
	file := "unknown"
	line := 0
	if _, file0, line0, ok := runtime.Caller(skip); ok {
		file0 = strings.Replace(file0, "\\", "/", -1)
		p := "github.com/snail007/gmc"
		if !s.skipCheckGMC && strings.Contains(file0, p) && !strings.Contains(file0, p+"t") &&
			!strings.Contains(file0, p+"demos") {
			// gmc
			file = "[gmc]" + file0[strings.Index(file0, p)+len(p):]
		} else if s.flag == gcore.LogFlagShort {
			//short
			file = filepath.Base(filepath.Dir(file0)) + "/" + filepath.Base(file0)
		} else {
			//long
			file = file0 + "/" + filepath.Base(file0)
		}
		line = line0
	}
	msg = fmt.Sprintf("%s:%d ", file, line) + msg
	return msg
}

type levelWriter struct {
	io.Writer
	level gcore.LogLevel
	lock  *sync.Mutex
}

func newLogWriter(writer io.Writer, level gcore.LogLevel) *levelWriter {
	return &levelWriter{Writer: writer, level: level, lock: &sync.Mutex{}}
}

type multiWriter struct {
	writers []io.Writer
}

var ErrShortWrite = errors.New("short write")

func (t *multiWriter) Write(p []byte) (n int, err error) {
	if len(t.writers) == 0 {
		return
	}
	g := sync.WaitGroup{}
	g.Add(len(t.writers))
	for _, w := range t.writers {
		w0 := w
		pool.Submit(func() {
			defer g.Done()
			n, e := w0.Write(p)
			if e != nil {
				err = e
				return
			}
			if n != len(p) {
				err = ErrShortWrite
			}
		})
	}
	g.Wait()
	if err != nil {
		return
	}
	return len(p), nil
}

func newMultiWriter(writers ...io.Writer) io.Writer {
	allWriters := make([]io.Writer, 0, len(writers))
	for _, w := range writers {
		if mw, ok := w.(*multiWriter); ok {
			allWriters = append(allWriters, mw.writers...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &multiWriter{allWriters}
}
