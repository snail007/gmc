package log

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

type bufChnItem struct {
	level gcore.LOG_LEVEL
	msg   string
}
type GMCLog struct {
	l         *log.Logger
	parent    *GMCLog
	ns        string
	level     gcore.LOG_LEVEL
	async     bool
	asyncOnce *sync.Once
	bufChn    chan bufChnItem
	asyncWG   *sync.WaitGroup
}

func NewGMCLog(prefix ...string) gcore.Logger {
	pre := ""
	if len(prefix) == 1 {
		pre = prefix[0]
	}
	l := log.New(os.Stdout, pre, log.LstdFlags|log.Lmicroseconds)
	return &GMCLog{
		l:         l,
		level:     gcore.LDEBUG,
		asyncOnce: &sync.Once{},
	}
}

func (s *GMCLog) WaitAsyncDone() {
	s.asyncWG.Wait()
}

func (s *GMCLog) Async() bool {
	return s.async
}

func (s *GMCLog) asyncWriterInit() {
	s.bufChn = make(chan bufChnItem, 2048)
	s.asyncWG = &sync.WaitGroup{}
	go func() {
		for {
			item := <-s.bufChn
			s.output(item.msg)
			s.asyncWG.Done()
		}
	}()
}

func (s *GMCLog) EnableAsync() {
	s.async = true
	s.asyncOnce.Do(func() {
		s.asyncWriterInit()
	})
}

func (s *GMCLog) SetLevel(i gcore.LOG_LEVEL) {
	s.level = i
}

func (s *GMCLog) With(namespace string) gcore.Logger {
	return &GMCLog{
		l:      s.l,
		parent: s,
		ns:     namespace,
		level:  s.level,
	}
}

func (s *GMCLog) Namespace() string {
	ns := s.ns
	if s.parent != nil {
		ns = s.parent.Namespace() + "/" + s.ns
	}
	return strings.TrimLeft(ns, "/")
}

func (s *GMCLog) namespace() string {
	if s.parent != nil {
		return "[" + s.Namespace() + "] "
	}
	return ""
}

func (s *GMCLog) Panicf(format string, v ...interface{}) {
	if s.level > gcore.LPANIC {
		return
	}
	str := s.caller(fmt.Sprintf(s.namespace()+"PANIC "+format, v...))
	s.Write(str)
	panic(str)
}

func (s *GMCLog) Panic(v ...interface{}) {
	if s.level > gcore.LPANIC {
		return
	}
	v0 := []interface{}{s.namespace() + "PANIC "}
	str := s.caller(fmt.Sprint(append(v0, v...)...))
	s.Write(str)
	panic(str)
}

func (s *GMCLog) Errorf(format string, v ...interface{}) {
	if s.level > gcore.LERROR {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"ERROR "+format, v...)))
	os.Exit(1)
}

func (s *GMCLog) Error(v ...interface{}) {
	if s.level > gcore.LERROR {
		return
	}
	v0 := []interface{}{s.namespace() + "ERROR "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...)))
	os.Exit(1)
}

func (s *GMCLog) Warnf(format string, v ...interface{}) {
	if s.level > gcore.LWARN {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"WARN "+format, v...)))
}

func (s *GMCLog) Warn(v ...interface{}) {
	if s.level > gcore.LWARN {
		return
	}
	v0 := []interface{}{s.namespace() + "WARN "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...)))
}

func (s *GMCLog) Infof(format string, v ...interface{}) {
	if s.level > gcore.LINFO {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"INFO "+format, v...)))

}

func (s *GMCLog) Info(v ...interface{}) {
	if s.level > gcore.LINFO {
		return
	}
	v0 := []interface{}{s.namespace() + "INFO "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...)))
}

func (s *GMCLog) Debugf(format string, v ...interface{}) {
	if s.level > gcore.LDEBUG {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"DEBUG "+format, v...)))
}

func (s *GMCLog) Debug(v ...interface{}) {
	if s.level > gcore.LDEBUG {
		return
	}
	v0 := []interface{}{s.namespace() + "DEBUG "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...)))
}

func (s *GMCLog) Tracef(format string, v ...interface{}) {
	if s.level > gcore.LTRACE {
		return
	}
	s.Write(s.caller(fmt.Sprintf(s.namespace()+"TRACE "+format, v...)))
}

func (s *GMCLog) Trace(v ...interface{}) {
	if s.level > gcore.LTRACE {
		return
	}
	v0 := []interface{}{s.namespace() + "TRACE "}
	s.Write(s.caller(fmt.Sprint(append(v0, v...)...)))
}

func (s *GMCLog) Writer() io.Writer {
	return s.l.Writer()
}

func (s *GMCLog) SetOutput(w io.Writer) {
	s.l.SetOutput(w)
}

func (s *GMCLog) SetFlags(f int) {
	s.l.SetFlags(f)
}

func (s *GMCLog) Write(str string) {
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

func (s *GMCLog) output(str string) {
	s.l.Print(str)
}

func (s *GMCLog) caller(msg string) string {
	file := "unknown"
	line := 0
	if _, file0, line0, ok := runtime.Caller(2); ok {
		file0 = strings.Replace(file0, "\\", "/", -1)
		p := "/github.com/snail007/gmc/"
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
