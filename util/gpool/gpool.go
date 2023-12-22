// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
	glist "github.com/snail007/gmc/util/list"
	gmap "github.com/snail007/gmc/util/map"
	"io"
	"sync"
	"time"
)

const (
	statusRunning = iota + 1
	statusWaiting
	statusStopped
)

var (
	MaxWaitCountReached = errors.New("max await job count reached")
)

// GPool is a goroutine pool, you can increase or decrease pool size in runtime.
type GPool struct {
	tasks               *glist.List
	logger              gcore.Logger
	workers             *gmap.Map
	debug               bool
	maxWaitCount        int
	maxWorkCount        int
	g                   *sync.WaitGroup
	idleDuration        time.Duration
	submitBlockChanList *glist.List
	blockingOnMaxWait   bool
}

// BlockingOnMaxWait if the count of await job to run reach the max,
// if blocking Submit call
func (s *GPool) BlockingOnMaxWait() bool {
	return s.blockingOnMaxWait
}

func (s *GPool) SetBlockingOnMaxWait(blockingOnMaxWait bool) {
	s.blockingOnMaxWait = blockingOnMaxWait
}

// IdleDuration is the idle time duration before the worker exit,
// duration 0 means the work will not exit.
func (s *GPool) IdleDuration() time.Duration {
	return s.idleDuration
}

// SetIdleDuration set the idle time duration before the worker exit,
// duration 0 means the work will not exit.
func (s *GPool) SetIdleDuration(idleDuration time.Duration) {
	s.idleDuration = idleDuration
}

// MaxTaskAwaitCount returns the max waiting task count.
func (s *GPool) MaxTaskAwaitCount() int {
	return s.maxWaitCount
}

// SetMaxTaskAwaitCount sets the max waiting task count.
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
	return NewWithLogger(workerCount, nil)
}

func NewWithLogger(workerCount int, logger gcore.Logger) (p *GPool) {
	p = &GPool{
		submitBlockChanList: glist.New(),
		tasks:               glist.New(),
		logger:              logger,
		workers:             gmap.New(),
		maxWorkCount:        workerCount,
		g:                   &sync.WaitGroup{},
	}
	return
}

// Increase add the count of `workerCount` workers
func (s *GPool) Increase(workerCount int) {
	s.maxWorkCount += workerCount
	s.addWorker(workerCount)
}
func (s *GPool) removeWorker(w *worker) {
	w.Stop()
	s.workers.Delete(w.id)
}

// Decrease stop the count of `workerCount` workers
func (s *GPool) Decrease(workerCount int) {
	s.maxWorkCount -= workerCount
	if s.maxWorkCount < 0 {
		s.maxWorkCount = 0
	}
	// find awaiting workers
	s.workers.Range(func(_, v interface{}) bool {
		w := v.(*worker)
		if w.Status() == statusWaiting {
			s.removeWorker(w)
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
				s.removeWorker(w)
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
	s.maxWorkCount = workerCount
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

// WaitDone wait all the tasks submit task is executed done, if no task, return immediately.
func (s *GPool) WaitDone() {
	s.g.Wait()
}

func (s *GPool) addWorker(cnt int) {
	if s.WorkerCount() >= s.maxWorkCount {
		return
	}
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

// run a task function, using defer to catch task exception
func (s *GPool) run(fn func()) {
	defer func() {
		s.g.Done()
		if e := recover(); e != nil {
			s.log("GPool: a task stopped unexpectedly, err: %s", gcore.ProviderError()().StackError(e))
		}
	}()
	fn()
}

// Submit adds a function as a task ready to run
func (s *GPool) Submit(task func()) error {
	if s.WorkerCount() < s.maxWorkCount {
		s.addWorker(1)
	}
	if s.maxWaitCount > 0 && s.tasks.Len() >= s.maxWaitCount {
		if s.blockingOnMaxWait {
			ch := make(chan bool)
			s.submitBlockChanList.Add(ch)
			<-ch
		} else {
			return MaxWaitCountReached
		}
	}
	s.g.Add(1)
	s.tasks.Add(task)
	s.notifyAll()
	return nil
}

// notify all workers, only idle workers be awakened
func (s *GPool) notifyAll() {
	s.workers.Range(func(_, v interface{}) bool {
		v.(*worker).Wakeup()
		return true
	})
}

// shift an element from array head
func (s *GPool) pop() (fn func()) {
	f := s.tasks.Pop()
	if f != nil {
		fn = f.(func())
	}
	return
}

// Stop and remove all workers in the pool
func (s *GPool) Stop() {
	s.workers.Range(func(_, v interface{}) bool {
		v.(*worker).Stop()
		return true
	})
	s.workers.Clear()
}

// RunningWork returns the count of running workers
func (s *GPool) RunningWork() (workerCount int) {
	s.workers.Range(func(_, v interface{}) bool {
		if v.(*worker).Status() == statusRunning {
			workerCount++
		}
		return true
	})
	return
}

// AwaitingWorker returns the count of awaiting workers
func (s *GPool) AwaitingWorker() (workerCount int) {
	s.workers.Range(func(_, v interface{}) bool {
		if v.(*worker).Status() == statusWaiting {
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

// SetLogger set the logger to logging, you can SetLogger(nil) to disable logging
//
// default is log.New(os.Stdout, "", log.LstdFlags),
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
	if status == statusWaiting && w.pool.submitBlockChanList.Len() > 0 {
		if ch := w.pool.submitBlockChanList.Pop(); ch != nil {
			ch.(chan bool) <- true
		}
	}
}

func (w *worker) Wakeup() bool {
	defer gerror.RecoverNop()
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
	defer gerror.RecoverNop()
	select {
	case w.breakSig <- true:
	default:
		return false
	}
	return true
}

func (w *worker) Stop() {
	defer gerror.RecoverNop()
	w.breakLoop()
	close(w.wakeupSig)
}

func (w *worker) start() {
	go func() {
		w.Wakeup()
		var t *time.Timer
		defer func() {
			if t != nil {
				t.Stop()
			}
			w.SetStatus(statusStopped)
			w.pool.removeWorker(w)
			w.pool.debugLog("GPool: worker[%s] stopped", w.id)
		}()
		w.pool.debugLog("GPool: worker[%s] started ...", w.id)
		var fn func()
		var doJob = func() (isReturn bool) {
			w.SetStatus(statusRunning)
			w.pool.debugLog("GPool: worker[%s] running ...", w.id)
			for {
				//w.pool.debugLog("GPool: worker[%s] read break", w.id)
				if w.isBreak() {
					w.pool.debugLog("GPool: worker[%s] break", w.id)
					return true
				}
				if fn = w.pool.pop(); fn != nil {
					//w.pool.debugLog("GPool: worker[%s] called", w.id)
					w.pool.run(fn)
				} else {
					w.pool.debugLog("GPool: worker[%s] no task, break", w.id)
					break
				}
			}
			return
		}
		for {
			w.SetStatus(statusWaiting)
			w.pool.debugLog("GPool: worker[%s] waiting ...", w.id)
			if w.pool.idleDuration > 0 {
				if t == nil {
					t = time.NewTimer(w.pool.idleDuration)
				} else {
					t.Reset(w.pool.idleDuration)
				}
				select {
				case <-t.C:
					w.pool.debugLog("GPool: worker[%s] idle timeout, exited", w.id)
					return
				case _, ok := <-w.wakeupSig:
					if !ok {
						return
					}
					if doJob() {
						return
					}
				}
			} else {
				select {
				case _, ok := <-w.wakeupSig:
					if !ok {
						return
					}
					if doJob() {
						return
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
		wakeupSig: make(chan bool, 1),
		breakSig:  make(chan bool, 1),
		status:    statusWaiting,
	}
	w.start()
	return w
}
