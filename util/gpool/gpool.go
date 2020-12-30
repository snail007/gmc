// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gpool

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	gcore "github.com/snail007/gmc/core"
	"io"
	"sync"
	"sync/atomic"
)

type GPool struct {
	tasks        []func()
	taskLock     *sync.Mutex
	cntLock      *sync.Mutex
	ctx          context.Context
	cancel       context.CancelFunc
	runningtCnt  *int32
	logger       gcore.Logger
	maxWorkerCnt *int32
	workerCnt    *int32
	workerSig    *sync.Map
	isStop       bool
}

//create a gpool object to using
func NewGPool(workerCount int) (p *GPool) {
	ctx0, cancel0 := context.WithCancel(context.Background())
	p = &GPool{
		tasks:        []func(){},
		taskLock:     &sync.Mutex{},
		cntLock:      &sync.Mutex{},
		ctx:          ctx0,
		cancel:       cancel0,
		runningtCnt:  new(int32),
		workerCnt:    new(int32),
		workerSig:    &sync.Map{},
		maxWorkerCnt: new(int32),
		logger:       gcore.Providers.Logger("")(nil, ""),
	}
	p.addWorker(workerCount)
	return
}

func (s *GPool) Increase(workerCount int) {
	s.cntLock.Lock()
	defer s.cntLock.Unlock()
	s.addWorker(workerCount)
	s.notifyAll()
}

func (s *GPool) Decrease(workerCount int) {
	s.cntLock.Lock()
	defer s.cntLock.Unlock()
	if atomic.LoadInt32(s.maxWorkerCnt) == 0 {
		return
	}
	if atomic.LoadInt32(s.maxWorkerCnt)-int32(workerCount) < 0 {
		atomic.AddInt32(s.maxWorkerCnt, 0)
	} else {
		atomic.AddInt32(s.maxWorkerCnt, -int32(workerCount))
	}
	s.notifyAll()
}

func (s *GPool) Reset(workerCount int) {
	s.cntLock.Lock()
	defer s.cntLock.Unlock()
	atomic.StoreInt32(s.maxWorkerCnt, int32(workerCount))
	s.notifyAll()
}

func (s *GPool) WorkerCount() int {
	return int(atomic.LoadInt32(s.workerCnt))
}

func (s *GPool) newWorkerID() string {
	k := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return ""
	}
	return hex.EncodeToString(k)
}

func (s *GPool) addWorker(cnt int) {
	atomic.AddInt32(s.workerCnt, int32(cnt))
	atomic.AddInt32(s.maxWorkerCnt, int32(cnt))
	for i := 0; i < cnt; i++ {
		idx := s.newWorkerID()
		sig := make(chan bool, 1)
		s.workerSig.Store(idx, sig)
		go func(idx string, sig chan bool) {
			defer gcore.Providers.Error("")().Recover(func(e interface{}) {
				s.log("GPool: a worker stopped unexpectedly, err: %s", gcore.Providers.Error("")().StackError(e))
			})
			s.log("GPool: worker[%s] started ...", idx)
			ctx, cancel := context.WithCancel(s.ctx)
			defer cancel()
			var fn func()
			for {
				select {
				//checking stop is called
				case <-ctx.Done():
					// s.log("GPool: worker[%d] stopped.", i)
					return
				case <-sig:
					atomic.AddInt32(s.runningtCnt, 1)
					for {
						if s.isStop || s.needExit() {
							atomic.AddInt32(s.runningtCnt, -1)
							s.deleteWorker(idx)
							return
						}
						if fn = s.pop(); fn != nil {
							s.run(fn)
						} else {
							break
						}
					}
					atomic.AddInt32(s.runningtCnt, -1)
				}
			}
		}(idx, sig)
	}
	atomic.AddInt32(s.workerCnt, int32(cnt))
}

func (s *GPool) needExit() bool {
	if atomic.LoadInt32(s.maxWorkerCnt) == 0 {
		return false
	}
	return atomic.LoadInt32(s.workerCnt) > atomic.LoadInt32(s.maxWorkerCnt)
}

func (s *GPool) deleteWorker(idx string) {
	if _, ok := s.workerSig.Load(idx); ok {
		s.workerSig.Delete(idx)
		atomic.AddInt32(s.workerCnt, -1)
	}
}

//run a task function, using defer to catch task exception
func (s *GPool) run(fn func()) {
	defer gcore.Providers.Error("")().Recover(func(e interface{}) {
		s.log("GPool: a task stopped unexpectedly, err: %s", gcore.Providers.Error("")().StackError(e))
	})
	fn()
}

//submit a function as a task ready to run
func (s *GPool) Submit(task func()) {
	s.taskLock.Lock()
	defer s.taskLock.Unlock()
	s.tasks = append(s.tasks, task)
	s.notifyAll()
}

// notify all workers, only idle workers be awakened
func (s *GPool) notifyAll() {
	s.workerSig.Range(func(_, v interface{}) bool {
		select {
		case v.(chan bool) <- true:
		default:
		}
		return true
	})
}

//shift an element from array head
func (s *GPool) pop() (fn func()) {
	s.taskLock.Lock()
	defer s.taskLock.Unlock()
	l := len(s.tasks)
	if l > 0 {
		fn = s.tasks[0]
		s.tasks[0] = nil
		if l == 1 {
			s.tasks = []func(){}
		} else {
			s.tasks = s.tasks[1:]
		}
	}
	return
}

//stop all workers in the pool
func (s *GPool) Stop() {
	s.isStop = true
	s.cancel()
}

//return the number of running workers
func (s *GPool) Running() (cnt int32) {
	return atomic.LoadInt32(s.runningtCnt)
}

//return the number of task ready to run
func (s *GPool) Awaiting() (cnt int32) {
	return int32(len(s.tasks))
}
func (s *GPool) log(fmt string, v ...interface{}) {
	if s.logger == nil {
		return
	}
	s.logger.Infof(fmt, v...)
}

//SetLogger set the logger to logging, you can SetLogger(nil) to disable logging
//
//default is log.New(os.Stdout, "", log.LstdFlags),
func (s *GPool) SetLogger(l gcore.Logger) {
	s.logger = l
}
