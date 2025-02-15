// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
	glist "github.com/snail007/gmc/util/list"
	gmap "github.com/snail007/gmc/util/map"
	"io"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const (
	statusRunning = iota + 1
	statusIdle
	statusStopped
)

var (
	ErrMaxQueuedJobCountReached = errors.New("max queued job count reached")
)

// Pool is a goroutine pool, you can increase or decrease pool size in runtime.
type Pool struct {
	maxWorkCount         int
	jobs                 *glist.List
	workers              *gmap.Map
	g                    *sync.WaitGroup
	submitBlockChanList  *glist.List
	submitLock           *sync.Mutex
	opt                  *Option
	idleWorkerCounter    *int64
	runningWorkerCounter *int64
}

// Option sets the pool
type Option struct {
	//limits the max queued job count, 0 no limit
	MaxJobCount int
	// block the Submit call after the count of queued job to run reach the max, only worked on MaxJobCount is greater 0
	Blocking bool
	// output the debug logging, only worked on the pool Logger is not nil
	Debug bool
	// the logger to output debug logging
	Logger gcore.Logger
	// if IdleDuration nonzero, the worker will exited after idle duration when complete the job
	IdleDuration time.Duration
	// start the worker when the pool created
	PreAlloc bool
	// PanicHandler is used to handle panics from each job function.
	PanicHandler func(e interface{})
	//WithStack sets if fill stack info with submitted job
	WithStack bool
}

// Blocking  the count of queued job to run reach the max, if blocking Submit call
func (s *Pool) Blocking() bool {
	return s.opt.Blocking
}

// SetBlocking sets the count of queued job to run reach the max, if blocking Submit call
func (s *Pool) SetBlocking(blocking bool) {
	s.opt.Blocking = blocking
}

// IdleDuration is the idle time duration before the worker exit,
// duration 0 means the work will not exit.
func (s *Pool) IdleDuration() time.Duration {
	return s.opt.IdleDuration
}

// SetIdleDuration set the idle time duration before the worker exit,
// duration 0 means the work will not exit.
//
// Notice: if idle duration changed from zero, only the new worker will support the idle.
func (s *Pool) SetIdleDuration(idleDuration time.Duration) {
	s.opt.IdleDuration = idleDuration
}

// MaxJobCount returns the max queued job count.
func (s *Pool) MaxJobCount() int {
	return s.opt.MaxJobCount
}

// SetMaxJobCount sets the max queued job count.
func (s *Pool) SetMaxJobCount(maxJobCount int) {
	s.opt.MaxJobCount = maxJobCount
}

// IsDebug returns the pool in debug mode or not.
func (s *Pool) IsDebug() bool {
	return s.opt.Debug
}

// SetDebug sets the pool in debug mode, the pool will output more logging.
func (s *Pool) SetDebug(debug bool) {
	s.opt.Debug = debug
}

// New create a gpool object to using
func New(workerCount int) (p *Pool) {
	return NewWithOption(workerCount, &Option{
		WithStack: true,
	})
}

func NewWithLogger(workerCount int, logger gcore.Logger) (p *Pool) {
	return NewWithOption(workerCount, &Option{Logger: logger})
}

func NewWithPreAlloc(workerCount int) (p *Pool) {
	return NewWithOption(workerCount, &Option{PreAlloc: true})
}

func NewWithOption(workerCount int, opt *Option) (p *Pool) {
	p = &Pool{
		submitBlockChanList:  glist.New(),
		jobs:                 glist.New(),
		workers:              gmap.New(),
		maxWorkCount:         workerCount,
		g:                    &sync.WaitGroup{},
		submitLock:           &sync.Mutex{},
		opt:                  opt,
		idleWorkerCounter:    new(int64),
		runningWorkerCounter: new(int64),
	}
	if opt.PreAlloc {
		p.maxWorkCount = 0
		p.Increase(workerCount)
	}
	return p
}

// Increase add the count of `workerCount` workers
func (s *Pool) Increase(workerCount int) {
	s.increase(workerCount, true)
}

func (s *Pool) increase(workerCount int, modifyCounter bool) {
	if modifyCounter {
		s.maxWorkCount += workerCount
	}
	s.addWorker(workerCount)
}

func (s *Pool) removeWorker(w *worker) {
	w.Stop()
	s.workers.Delete(w.id)
}
func (s *Pool) Decrease(workerCount int) {
	s.decrease(workerCount, true)
}

// Decrease stop the count of `workerCount` workers
func (s *Pool) decrease(workerCount int, modifyCounter bool) {
	if modifyCounter {
		s.maxWorkCount -= workerCount
		if s.maxWorkCount < 0 {
			s.maxWorkCount = 0
		}
	}

	// find idle workers
	s.workers.Range(func(_, v interface{}) bool {
		w := v.(*worker)
		if w.Status() == statusIdle {
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
func (s *Pool) ResetTo(workerCount int) {
	length := s.workers.Len()
	if length == workerCount {
		return
	}
	s.maxWorkCount = workerCount
	if workerCount > length {
		s.increase(workerCount-length, false)
	} else {
		s.decrease(length-workerCount, false)
	}
}

// WorkerCount returns the count of workers
func (s *Pool) WorkerCount() int {
	return s.workers.Len()
}

// WaitDone wait all the jobs submitted executed done, if no job, return immediately.
func (s *Pool) WaitDone() {
	s.g.Wait()
}

func (s *Pool) addWorker(cnt int) {
	if s.WorkerCount() >= s.maxWorkCount {
		return
	}
	for i := 0; i < cnt; i++ {
		w := newWorker(s)
		s.workers.Store(w.id, w)
	}
}

func (s *Pool) newWorkerID() string {
	k := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return ""
	}
	return hex.EncodeToString(k)
}

// run a job function, using defer to catch job exception
func (s *Pool) run(j *JobItem) {
	defer func() {
		s.g.Done()
		if e := recover(); e != nil {
			msg := fmt.Sprintf("Pool: a job stopped unexpectedly, err: %s", gcore.ProviderError()().StackError(e))
			if s.opt.Logger != nil {
				s.opt.Logger.Error(msg)
			} else {
				fmt.Println(msg)
			}
			if s.opt.PanicHandler != nil {
				s.opt.PanicHandler(e)
			}
		}
	}()
	j.Job()
}

type JobItem struct {
	Stack string
	Job   func()
}

// Submit adds a function as a job ready to run
func (s *Pool) Submit(job func()) error {
	s.submitLock.Lock()
	defer s.submitLock.Unlock()
	if s.WorkerCount() < s.maxWorkCount && !s.hasIdleWorker() {
		s.addWorker(1)
	}
	if s.opt.MaxJobCount > 0 && s.jobs.Len() >= s.opt.MaxJobCount {
		if s.opt.Blocking {
			ch := make(chan bool)
			s.submitBlockChanList.Add(ch)
			<-ch
		} else {
			return ErrMaxQueuedJobCountReached
		}
	}
	j := &JobItem{
		Job: job,
	}
	if s.opt.WithStack {
		a := bytes.SplitN(debug.Stack(), []byte("\n"), 4)
		stackStr := string(a[0])
		if len(a) > 3 {
			stackStr += "\n" + string(a[3])
		}
		j.Stack = stackStr
	}
	s.g.Add(1)
	s.jobs.Add(j)
	s.notifyAll()
	return nil
}

// notify all workers, only idle workers be awakened
func (s *Pool) notifyAll() {
	s.workers.RangeFast(func(_, v interface{}) bool {
		if v.(*worker).Status() == statusIdle {
			v.(*worker).Wakeup()
		}
		return true
	})
}

// shift an element from array head
func (s *Pool) pop() (fn *JobItem) {
	f := s.jobs.Pop()
	if f != nil {
		fn = f.(*JobItem)
	}
	return
}

// Stop and remove all workers in the pool
func (s *Pool) Stop() {
	s.workers.RangeFast(func(_, v interface{}) bool {
		v.(*worker).Stop()
		return true
	})
	s.workers.Clear()
}

// RunningWorkerCount returns the count of running workers
func (s *Pool) RunningWorkerCount() (workerCount int) {
	return int(atomic.LoadInt64(s.runningWorkerCounter))
}

// IdleWorkerCount returns the count of idle workers
func (s *Pool) IdleWorkerCount() (workerCount int) {
	return int(atomic.LoadInt64(s.idleWorkerCounter))
}

func (s *Pool) hasIdleWorker() (has bool) {
	return atomic.LoadInt64(s.idleWorkerCounter) > 0
}

// QueuedJobCount returns the count of queued job
func (s *Pool) QueuedJobCount() (jobCount int) {
	return s.jobs.Len()
}

func (s *Pool) debugLog(fmt string, v ...interface{}) {
	if !s.opt.Debug {
		return
	}
	s.log(fmt, v...)
}

func (s *Pool) log(fmt string, v ...interface{}) {
	if s.opt.Logger == nil {
		return
	}
	s.opt.Logger.Warnf(fmt, v...)
}

// SetLogger set the logger to logging, you can SetLogger(nil) to disable logging
//
// default is log.New(os.Stdout, "", log.LstdFlags),
func (s *Pool) SetLogger(l gcore.Logger) {
	s.opt.Logger = l
}

type worker struct {
	status    int
	pool      *Pool
	wakeupSig chan bool
	breakSig  chan bool
	id        string
}

func (w *worker) Status() int {
	return w.status
}

func (w *worker) SetStatus(status int) {
	w.status = status
	if status == statusIdle && w.pool.submitBlockChanList.Len() > 0 {
		if ch := w.pool.submitBlockChanList.Pop(); ch != nil {
			ch.(chan bool) <- true
		}
	}
}

func (w *worker) addWorkerCounter(cnt *int64, val int64) {
	if val < 0 {
		if atomic.LoadInt64(cnt)-val >= 0 {
			atomic.AddInt64(cnt, val)
		} else {
			atomic.AddInt64(cnt, 0)
		}
	} else {
		atomic.AddInt64(cnt, val)
	}
}

func (w *worker) addIdleWorkerCounter(val int64) {
	w.addWorkerCounter(w.pool.idleWorkerCounter, val)
}

func (w *worker) addRunningWorkerCounter(val int64) {
	w.addWorkerCounter(w.pool.runningWorkerCounter, val)
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
			w.pool.debugLog("Pool: worker[%s] stopped", w.id)
		}()
		w.pool.debugLog("Pool: worker[%s] started ...", w.id)
		j := new(JobItem)
		var doJob = func() (isReturn bool) {
			w.SetStatus(statusRunning)
			w.addRunningWorkerCounter(1)
			defer w.addRunningWorkerCounter(-1)
			w.pool.debugLog("Pool: worker[%s] running ...", w.id)
			for {
				//w.pool.debugLog("Pool: worker[%s] read break", w.id)
				if w.isBreak() {
					w.pool.debugLog("Pool: worker[%s] break", w.id)
					return true
				}
				if j = w.pool.pop(); j != nil {
					//w.pool.debugLog("Pool: worker[%s] called", w.id)
					w.pool.run(j)
				} else {
					w.pool.debugLog("Pool: worker[%s] no job, break", w.id)
					break
				}
			}
			return
		}
		for {
			w.SetStatus(statusIdle)
			w.addIdleWorkerCounter(1)
			w.pool.debugLog("Pool: worker[%s] waiting ...", w.id)
			if w.pool.opt.IdleDuration > 0 {
				if t == nil {
					t = time.NewTimer(w.pool.opt.IdleDuration)
				} else {
					t.Reset(w.pool.opt.IdleDuration)
				}
				select {
				case <-t.C:
					w.pool.debugLog("Pool: worker[%s] idle timeout, exited", w.id)
					w.addIdleWorkerCounter(-1)
					return
				case _, ok := <-w.wakeupSig:
					if !ok {
						w.addIdleWorkerCounter(-1)
						return
					}
					w.addIdleWorkerCounter(-1)
					if doJob() {
						return
					}
				}
			} else {
				select {
				case _, ok := <-w.wakeupSig:
					if !ok {
						w.addIdleWorkerCounter(-1)
						return
					}
					w.addIdleWorkerCounter(-1)
					if doJob() {
						return
					}
				}
			}
		}
	}()
}

func newWorker(pool *Pool) *worker {
	w := &worker{
		pool:      pool,
		id:        pool.newWorkerID(),
		wakeupSig: make(chan bool, 1),
		breakSig:  make(chan bool, 1),
		status:    statusIdle,
	}
	w.start()
	return w
}
