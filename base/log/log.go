package gmclog

import (
	"github.com/snail007/gmc/core"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type GMCLog struct {
	l      *log.Logger
	parent *GMCLog
	ns     string
	level  gmccore.LOG_LEVEL
	mu  sync.RWMutex
}

func NewGMCLog() gmccore.Logger {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	return &GMCLog{
		l:     l,
		level: gmccore.LDEBUG,
	}
}

func (G *GMCLog) SetLevel(i gmccore.LOG_LEVEL) {
	G.level = i
}

func (G *GMCLog) With(namespace string) gmccore.Logger {
	return &GMCLog{
		l:      G.l,
		parent: G,
		ns:     namespace,
		level:  G.level,
	}
}

func (G *GMCLog) Namespace() string {
	ns := G.ns
	if G.parent != nil {
		ns = G.parent.ns + "/" + G.ns
	}
	return strings.TrimLeft(ns, "/")
}

func (G *GMCLog) namespace() string {
	if G.parent != nil {
		return "[" + G.Namespace() + "] "
	}
	return ""
}

func (G *GMCLog) Panic(v ...interface{}) {
	if G.level > gmccore.LPANIC {
		return
	}
	v0 := []interface{}{G.namespace() + "PANIC "}
	G.l.Panic(append(v0, v...)...)
}

func (G *GMCLog) Panicf(format string, v ...interface{}) {
	if G.level > gmccore.LPANIC {
		return
	}
	G.l.Panicf(G.namespace()+"PANIC "+format, v...)
}

func (G *GMCLog) Errorf(format string, v ...interface{}) {
	if G.level > gmccore.LERROR {
		return
	}
	G.l.Fatalf(G.namespace()+"ERROR "+format, v...)
}

func (G *GMCLog) Error(v ...interface{}) {
	if G.level > gmccore.LERROR {
		return
	}
	v0 := []interface{}{G.namespace() + "ERROR "}
	G.l.Fatal(append(v0, v...)...)
}

func (G *GMCLog) Warnf(format string, v ...interface{}) {
	if G.level > gmccore.LWARN {
		return
	}
	G.l.Printf(G.namespace()+"WARN "+format, v...)
}

func (G *GMCLog) Warn(v ...interface{}) {
	if G.level > gmccore.LWARN {
		return
	}
	v0 := []interface{}{G.namespace() + "WARN "}
	G.l.Print(append(v0, v...)...)
}


func (G *GMCLog) Infof(format string, v ...interface{}) {
	if G.level > gmccore.LINFO {
		return
	}
	G.l.Printf(G.namespace()+"INFO "+format, v...)
}

func (G *GMCLog) Info(v ...interface{}) {
	if G.level > gmccore.LINFO {
		return
	}
	v0 := []interface{}{G.namespace() + "INFO "}
	G.l.Print(append(v0, v...)...)
}

func (G *GMCLog) Debugf(format string, v ...interface{}) {
	if G.level > gmccore.LDEBUG {
		return
	}
	G.l.Printf(G.namespace()+"DEBUG "+format, v...)
}

func (G *GMCLog) Debug(v ...interface{}) {
	if G.level > gmccore.LDEBUG {
		return
	}
	v0 := []interface{}{G.namespace() + "DEBUG "}
	G.l.Print(append(v0, v...)...)
}


func (G *GMCLog) Tracef(format string, v ...interface{}) {
	if G.level > gmccore.LTRACE {
		return
	}
	G.l.Printf(G.namespace()+"TRACE "+format, v...)
}

func (G *GMCLog) Trace(v ...interface{}) {
	if G.level > gmccore.LTRACE {
		return
	}
	v0 := []interface{}{G.namespace() + "TRACE "}
	G.l.Print(append(v0, v...)...)
}


func (G *GMCLog) Writer() io.Writer {
	return G.l.Writer()
}

func (G *GMCLog) SetOutput(w io.Writer) {
	G.mu.Lock()
	defer G.mu.Unlock()
	G.l.SetOutput(w)
}
