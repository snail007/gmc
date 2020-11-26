package gmclog

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
	level gmccore.LOG_LEVEL
	msg   string
}
type GMCLog struct {
	l         *log.Logger
	parent    *GMCLog
	ns        string
	level     gmccore.LOG_LEVEL
	async     bool
	asyncOnce *sync.Once
	bufChn    chan bufChnItem
	asyncWG   *sync.WaitGroup
}

func NewGMCLog() gmccore.Logger {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	return &GMCLog{
		l:         l,
		level:     gmccore.LDEBUG,
		asyncOnce: &sync.Once{},
	}
}

func (g *GMCLog) WaitAsyncDone() {
	g.asyncWG.Wait()
}

func (g *GMCLog) Async() bool {
	return g.async
}

func (g *GMCLog) asyncWriterInit() {
	g.bufChn = make(chan bufChnItem, 2048)
	g.asyncWG = &sync.WaitGroup{}
	go func() {
		for {
			item := <-g.bufChn
			g.output(item.msg)
			g.asyncWG.Done()
		}
	}()
}

func (g *GMCLog) EnableAsync() {
	g.async = true
	g.asyncOnce.Do(func() {
		g.asyncWriterInit()
	})
}

func (g *GMCLog) SetLevel(i gmccore.LOG_LEVEL) {
	g.level = i
}

func (g *GMCLog) With(namespace string) gmccore.Logger {
	return &GMCLog{
		l:      g.l,
		parent: g,
		ns:     namespace,
		level:  g.level,
	}
}

func (g *GMCLog) Namespace() string {
	ns := g.ns
	if g.parent != nil {
		ns = g.parent.Namespace() + "/" + g.ns
	}
	return strings.TrimLeft(ns, "/")
}

func (g *GMCLog) namespace() string {
	if g.parent != nil {
		return "[" + g.Namespace() + "] "
	}
	return ""
}

func (g *GMCLog) Panicf(format string, v ...interface{}) {
	if g.level > gmccore.LPANIC {
		return
	}
	s := g.caller(fmt.Sprintf(g.namespace()+"PANIC "+format, v...))
	g.Write(s)
	panic(s)
}

func (g *GMCLog) Panic(v ...interface{}) {
	if g.level > gmccore.LPANIC {
		return
	}
	v0 := []interface{}{g.namespace() + "PANIC "}
	s := g.caller(fmt.Sprint(append(v0, v...)...))
	g.Write(s)
	panic(s)
}

func (g *GMCLog) Errorf(format string, v ...interface{}) {
	if g.level > gmccore.LERROR {
		return
	}
	g.Write(g.caller(fmt.Sprintf(g.namespace()+"ERROR "+format, v...)))
	os.Exit(1)
}

func (g *GMCLog) Error(v ...interface{}) {
	if g.level > gmccore.LERROR {
		return
	}
	v0 := []interface{}{g.namespace() + "ERROR "}
	g.Write(g.caller(fmt.Sprint(append(v0, v...)...)))
	os.Exit(1)
}

func (g *GMCLog) Warnf(format string, v ...interface{}) {
	if g.level > gmccore.LWARN {
		return
	}
	g.Write(g.caller(fmt.Sprintf(g.namespace()+"WARN "+format, v...)))
}

func (g *GMCLog) Warn(v ...interface{}) {
	if g.level > gmccore.LWARN {
		return
	}
	v0 := []interface{}{g.namespace() + "WARN "}
	g.Write(g.caller(fmt.Sprint(append(v0, v...)...)))
}

func (g *GMCLog) Infof(format string, v ...interface{}) {
	if g.level > gmccore.LINFO {
		return
	}
	g.Write(g.caller(fmt.Sprintf(g.namespace()+"INFO "+format, v...)))

}

func (g *GMCLog) Info(v ...interface{}) {
	if g.level > gmccore.LINFO {
		return
	}
	v0 := []interface{}{g.namespace() + "INFO "}
	g.Write(g.caller(fmt.Sprint(append(v0, v...)...)))
}

func (g *GMCLog) Debugf(format string, v ...interface{}) {
	if g.level > gmccore.LDEBUG {
		return
	}
	g.Write(g.caller(fmt.Sprintf(g.namespace()+"DEBUG "+format, v...)))
}

func (g *GMCLog) Debug(v ...interface{}) {
	if g.level > gmccore.LDEBUG {
		return
	}
	v0 := []interface{}{g.namespace() + "DEBUG "}
	g.Write(g.caller(fmt.Sprint(append(v0, v...)...)))
}

func (g *GMCLog) Tracef(format string, v ...interface{}) {
	if g.level > gmccore.LTRACE {
		return
	}
	g.Write(g.caller(fmt.Sprintf(g.namespace()+"TRACE "+format, v...)))
}

func (g *GMCLog) Trace(v ...interface{}) {
	if g.level > gmccore.LTRACE {
		return
	}
	v0 := []interface{}{g.namespace() + "TRACE "}
	g.Write(g.caller(fmt.Sprint(append(v0, v...)...)))
}

func (g *GMCLog) Writer() io.Writer {
	return g.l.Writer()
}

func (g *GMCLog) SetOutput(w io.Writer) {
	g.l.SetOutput(w)
}

func (g *GMCLog) SetFlags(f int) {
	g.l.SetFlags(f)
}

func (g *GMCLog) Write(s string) {
	if g.async {
		select {
		case g.bufChn <- bufChnItem{
			msg: s,
		}:
			g.asyncWG.Add(1)
		default:
			g.output("WARN gmclog buf chan overflow")
		}
		return
	}
	g.output(s)
}

func (g *GMCLog) output(s string) {
	g.l.Print(s)
}

func (g *GMCLog) caller(msg string) string {
	file := "unknown"
	line := 0
	if _, file0, line0, ok := runtime.Caller(2); ok {
		file0=strings.Replace(file0,"\\","/",-1)
  		p:="/github.com/snail007/gmc/"
		if strings.Contains(file0, p) &&
			!strings.Contains(file0, p+"demos") {
			file = "[gmc]" + file0[strings.Index(file0,p)+len(p):]
		} else {
			file = filepath.Base(filepath.Dir(file0)) + "/" + filepath.Base(file0)
		}
		line = line0
	}
	msg = fmt.Sprintf("%s:%d: ", file, line) + msg
	return msg
}
