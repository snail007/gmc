// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool

import (
	"crypto/rand"
	"encoding/hex"
	gcore "github.com/snail007/gmc/core"
	glist "github.com/snail007/gmc/util/list"
	gmap "github.com/snail007/gmc/util/map"
	"io"
	"sync"
)

const (
	statusRunning = iota + 1
	statusWaiting
	statusStopped
)

// GPool is a goroutine pool, you can increase or decrease pool size in runtime.
type GPool struct {
	tasks         *glist.List
	logger        gcore.Logger
	workers       *gmap.Map
	debug         bool
	maxWaitCount  int
	lazy          sync.Once
	initWorkCount int
}

//MaxTaskAwaitCount returns the max waiting task count.
func (s *GPool) MaxTaskAwaitCount() int {
	return s.maxWaitCount
}

//SetMaxTaskAwaitCount sets the max waiting task count.
func (s *GPool) SetMaxTaskAwaitCount(maxWaitCount int) {
	s.maxWaitCount = maxWaitCount
}

// IsDebug returns the pool in debug mode or not.
func (s *GPool) IsDebug() bool {
	return s.debug
}

// SetDebug sets the pool in debug mode, the pool will output more logging.
func (s *GPool) SetDebug(debug bool) {
	s.debug = debug
}

// New create a gpool object to using
func New(workerCount int) (p *GPool) {
	p = &GPool{
		tasks:         glist.New(),
		logger:        gcore.ProviderLogger()(nil, ""),
		workers:       gmap.New(),
		initWorkCount: workerCount,
	}
	return
}

// Increase add the count of `workerCount` workers
func (s *GPool) Increase(workerCount int) {
	s.addWorker(workerCount)
}

// Decrease stop the count of `workerCount` workers
func (s *GPool) Decrease(workerCount int) {
	// find idle workers
	s.workers.Range(func(_, v interface{}) bool {
		w := v.(*worker)
		if w.Status() == statusWaiting {
			w.Stop()
			s.workers.Delete(w.id)
			workerCount--
			if workerCount == 0 {
				return false
			}
		}
		return true
	})
	// workerCount still great 0, stop some running workers
	if workerCount > 0 {
		s.workers.Range(func(_, v interface{}) bool {
			w := v.(*worker)
			if w.Status() == statusRunning {
				v.(*worker).Stop()
				s.workers.Delete(w.id)
				workerCount--
				if workerCount == 0 {
					return false
				}
			}
			return true
		})
	}
}

// ResetTo set the count of workers
func (s *GPool) ResetTo(workerCount int) {
	length := s.workers.Len()
	if length == workerCount {
		return
	}
	if workerCount > length {
		s.Increase(workerCount - length)
	} else {
		s.Decrease(length - workerCount)
	}
}

// WorkerCount returns the count of workers
func (s *GPool) WorkerCount() int {
	return s.workers.Len()
}

func (s *GPool) addWorker(cnt int) {
	for i := 0; i < cnt; i++ {
		w := newWorker(s)
		s.workers.Store(w.id, w)
	}
}

func (s *GPool) newWorkerID() string {
	k := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return ""
	}
	return hex.EncodeToString(k)
}

//run a task function, using defer to catch task exception
func (s *GPool) run(fn func()) {
	defer gcore.ProviderError()().Recover(func(e interface{}) {
		s.log("GPool: a task stopped unexpectedly, err: %s", gcore.ProviderError()().StackError(e))
	})
	fn()
}

//Submit adds a function as a task ready to run
func (s *GPool) Submit(task func()) bool {
	s.lazy.Do(func() {
		s.addWorker(s.initWorkCount)
	})
	if s.maxWaitCount > 0 && s.tasks.Len() > s.maxWaitCount {
		return false
	}
	s.tasks.Add(task)
	s.notifyAll()
	return true
}

// notify all workers, only idle workers be awakened
func (s *GPool) notifyAll() {
	s.workers.Range(func(_, v interface{}) bool {
		v.(*worker).Wakeup()
		return true
	})
}

//shift an element from array head
func (s *GPool) pop() (fn func()) {
	f := s.tasks.Pop()
	if f != nil {
		fn = f.(func())
	}
	return
}

// Stop stop and remove all workers in the pool
func (s *GPool) Stop() {
	s.workers.Range(func(_, v interface{}) bool {
		v.(*worker).Stop()
		return true
	})
	s.workers.Clear()
}

// Running returns the count of running workers
func (s *GPool) Running() (workerCount int) {
	s.workers.Range(func(_, v interface{}) bool {
		if v.(*worker).Status() == statusRunning {
			workerCount++
		}
		return true
	})
	return
}

// Awaiting returns the count of task ready to run
func (s *GPool) Awaiting() (taskCount int) {
	return s.tasks.Len()
}
func (s *GPool) debugLog(fmt string, v ...interface{}) {
	if !s.debug {
		return
	}
	s.log(fmt, v...)
}
func (s *GPool) log(fmt string, v ...interface{}) {
	if s.logger == nil {
		return
	}
	s.logger.Warnf(fmt, v...)
}

//SetLogger set the logger to logging, you can SetLogger(nil) to disable logging
//
//default is log.New(os.Stdout, "", log.LstdFlags),
func (s *GPool) SetLogger(l gcore.Logger) {
	s.logger = l
}

type worker struct {
	status    int
	pool      *GPool
	wakeupSig chan bool
	breakSig  chan bool
	id        string
}

func (w *worker) Status() int {
	return w.status
}

func (w *worker) SetStatus(status int) {
	w.status = status
}

func (w *worker) Wakeup() bool {
	defer gcore.ProviderError()().Recover()
	select {
	case w.wakeupSig <- true:
	default:
		return false
	}
	return true
}

func (w *worker) isBreak() bool {
	select {
	case <-w.breakSig:
		return true
	default:
		return false
	}
}

func (w *worker) breakLoop() bool {
	defer gcore.ProviderError()().Recover()
	select {
	case w.breakSig <- true:
	default:
		return false
	}
	return true
}

func (w *worker) Stop() {
	defer gcore.ProviderError()().Recover()
	w.breakLoop()
	close(w.wakeupSig)
}

func (w *worker) start() {
	go func() {
		w.Wakeup()
		defer func() {
			w.SetStatus(statusStopped)
			w.pool.debugLog("GPool: worker[%s] stopped", w.id)
		}()
		w.pool.debugLog("GPool: worker[%s] started ...", w.id)
		var fn func()
		for {
			w.SetStatus(statusWaiting)
			w.pool.debugLog("GPool: worker[%s] waiting ...", w.id)
			select {
			case _, ok := <-w.wakeupSig:
				if !ok {
					return
				}
				w.SetStatus(statusRunning)
				w.pool.debugLog("GPool: worker[%s] running ...", w.id)
				for {
					//w.pool.debugLog("GPool: worker[%s] read break", w.id)
					if w.isBreak() {
						w.pool.debugLog("GPool: worker[%s] break", w.id)
						return
					}
					if fn = w.pool.pop(); fn != nil {
						//w.pool.debugLog("GPool: worker[%s] called", w.id)
						w.pool.run(fn)
					} else {
						w.pool.debugLog("GPool: worker[%s] no task, break", w.id)
						break
					}
				}
			}
		}
	}()
}

func newWorker(pool *GPool) *worker {
	w := &worker{
		pool:      pool,
		id:        pool.newWorkerID(),
		wakeupSig: make(chan bool,1),
		breakSig:  make(chan bool,1),
		status:    statusWaiting,
	}
	w.start()
	return w
}
