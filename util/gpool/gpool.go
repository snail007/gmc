// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gpool

import (
	"context"
	gmccore "github.com/snail007/gmc/core"
	gmcerr "github.com/snail007/gmc/error"
	logutil "github.com/snail007/gmc/util/log"
	"sync"
	"sync/atomic"
)

type GPool struct {
	tasks       []func()
	lock        sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
	runningtCnt *int32
	logger      gmccore.Logger
	workerCnt   int
	workerSig   []chan bool
	isStop      bool
}

//create a gpool object to using
func New(workerCount int) (p *GPool) {
	cnt := int32(0)
	ctx0, cancel0 := context.WithCancel(context.Background())
	p = &GPool{
		tasks:       []func(){},
		lock:        sync.Mutex{},
		ctx:         ctx0,
		cancel:      cancel0,
		runningtCnt: &cnt,
		workerCnt:   workerCount,
		workerSig:   []chan bool{},
		logger:      logutil.New(""),
	}
	for i := 0; i < p.workerCnt; i++ {
		p.workerSig = append(p.workerSig, make(chan bool, 1))
	}
	p.init()
	return
}

//initialize workers to run tasks, a work is a goroutine
func (s *GPool) init() {
	go func() {
		defer gmcerr.Recover(func(e interface{}) {
			s.log("GPool stopped unexpectedly, err: %s", gmcerr.Stack(e))
		})
		//start the workerCnt workers
		for i := 0; i < s.workerCnt; i++ {
			go func(i int, sig chan bool) {
				defer gmcerr.Recover(func(e interface{}) {
					s.log("GPool: a worker stopped unexpectedly, err: %s", gmcerr.Stack(e))
				})
				//s.log("GPool: worker[%d] started ...", i)
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
							if s.isStop {
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
			}(i, s.workerSig[i])
		}
	}()
}

//run a task function, using defer to catch task exception
func (s *GPool) run(fn func()) {
	defer gmcerr.Recover(func(e interface{}) {
		s.log("GPool: a task stopped unexpectedly, err: %s", gmcerr.Stack(e))
	})
	fn()
}

//submit a function as a task ready to run
func (s *GPool) Submit(task func()) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.tasks = append(s.tasks, task)
	s.notifyAll()
}

// notify all workers, only idle workers be awakened
func (s *GPool) notifyAll() {
	for _, v := range s.workerSig {
		select {
		case v <- true:
		default:
		}
	}
}

//shift an element from array head
func (s *GPool) pop() (fn func()) {
	s.lock.Lock()
	defer s.lock.Unlock()
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
func (s *GPool) SetLogger(l gmccore.Logger) {
	s.logger = l
}
