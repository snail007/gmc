package gmclog

import (
	"github.com/snail007/gmc/core"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type bufChnItem struct {
	level    gmccore.LOG_LEVEL
	isFormat bool
	format   string
	msg      []interface{}
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
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
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
			if item.isFormat {
				g.l.Printf(item.format, item.msg...)
			} else {
				g.l.Print(item.msg...)
			}
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

func (g *GMCLog) Panic(v ...interface{}) {
	if g.level > gmccore.LPANIC {
		return
	}
	v0 := []interface{}{g.namespace() + "PANIC "}
	g.l.Panic(append(v0, v...)...)
}

func (g *GMCLog) Panicf(format string, v ...interface{}) {
	if g.level > gmccore.LPANIC {
		return
	}
	g.l.Panicf(g.namespace()+"PANIC "+format, v...)
}

func (g *GMCLog) Errorf(format string, v ...interface{}) {
	if g.level > gmccore.LERROR {
		return
	}
	g.l.Fatalf(g.namespace()+"ERROR "+format, v...)
}

func (g *GMCLog) Error(v ...interface{}) {
	if g.level > gmccore.LERROR {
		return
	}
	v0 := []interface{}{g.namespace() + "ERROR "}
	g.l.Fatal(append(v0, v...)...)
}

func (g *GMCLog) Warnf(format string, v ...interface{}) {
	if g.level > gmccore.LWARN {
		return
	}
	g.Writef(g.namespace()+"WARN "+format, v...)
}

func (g *GMCLog) Warn(v ...interface{}) {
	if g.level > gmccore.LWARN {
		return
	}
	v0 := []interface{}{g.namespace() + "WARN "}
	g.Write(append(v0, v...)...)
}

func (g *GMCLog) Infof(format string, v ...interface{}) {
	if g.level > gmccore.LINFO {
		return
	}
	g.Writef(g.namespace()+"INFO "+format, v...)
}

func (g *GMCLog) Info(v ...interface{}) {
	if g.level > gmccore.LINFO {
		return
	}
	v0 := []interface{}{g.namespace() + "INFO "}
	g.Write(append(v0, v...)...)
}

func (g *GMCLog) Debugf(format string, v ...interface{}) {
	if g.level > gmccore.LDEBUG {
		return
	}
	g.Writef(g.namespace()+"DEBUG "+format, v...)
}

func (g *GMCLog) Debug(v ...interface{}) {
	if g.level > gmccore.LDEBUG {
		return
	}
	v0 := []interface{}{g.namespace() + "DEBUG "}
	g.Write(append(v0, v...)...)
}

func (g *GMCLog) Tracef(format string, v ...interface{}) {
	if g.level > gmccore.LTRACE {
		return
	}
	g.Writef(g.namespace()+"TRACE "+format, v...)
}

func (g *GMCLog) Trace(v ...interface{}) {
	if g.level > gmccore.LTRACE {
		return
	}
	v0 := []interface{}{g.namespace() + "TRACE "}
	g.Write(append(v0, v...)...)
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

func (g *GMCLog) Write(v ...interface{}) {
	if g.async {
		select {
		case g.bufChn <- bufChnItem{
			isFormat: false,
			msg:      v,
		}:
			g.asyncWG.Add(1)
		default:
			g.l.Print("WARN gmclog buf chan overflow")
		}
		return
	}
	g.l.Print(v...)
}

func (g *GMCLog) Writef(format string, v ...interface{}) {
	if g.async {
		select {
		case g.bufChn <- bufChnItem{
			isFormat: true,
			format:   format,
			msg:      v,
		}:
			g.asyncWG.Add(1)
		default:
			g.l.Print("WARN gmclog buf chan overflow")
		}
		return
	}
	g.l.Printf(format, v...)
}