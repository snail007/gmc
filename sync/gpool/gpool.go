// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gpool

import (
	"context"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type GPool struct {
	tasks       []func()
	lock        sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
	runningtCnt *int32
	logger      *log.Logger
	workerCnt   int
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
		logger:      log.New(os.Stdout, "", log.LstdFlags),
	}
	p.init()
	return
}

//initialize workers to run tasks, a work is a goroutine
func (s *GPool) init() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				s.log("GPool stopped unexceptedly, err: %s", e)
			}
		}()
		//start the workerCnt workers
		for i := 0; i < s.workerCnt; i++ {
			go func(i int) {
				defer func() {
					if e := recover(); e != nil {
						s.log("GPool: a worker stopped unexceptedly, err: %s", e)
					}
				}()
				// s.log("GPool: worker[%d] started ...", i)
				ctx, cancle := context.WithCancel(s.ctx)
				defer cancle()
				var fn func()
				for {
					select {
					//checking stop is called
					case <-ctx.Done():
						// s.log("GPool: worker[%d] stopped.", i)
						return
					default:
						if fn = s.pop(); fn != nil {
							atomic.AddInt32(s.runningtCnt, 1)
							s.run(fn)
							atomic.AddInt32(s.runningtCnt, -1)
							fn = nil
						} else {
							//no task to run, we sleep a while
							time.Sleep(time.Second * 3)
						}
					}
				}
			}(i)
		}
	}()
}

//run a task function, using defer to catch task exception
func (s *GPool) run(fn func()) {
	defer func() {
		if e := recover(); e != nil {
			s.log("GPool: a task stopped unexceptedly, err: %s", e)
		}
	}()
	fn()
}

//submit a function as a task ready to run
func (s *GPool) Submit(task func()) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.tasks = append(s.tasks, task)
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
	s.logger.Printf(fmt, v...)
}

//SetLogger set the logger to logging, you can SetLogger(nil) to disable logging
//
//default is log.New(os.Stdout, "", log.LstdFlags),
func (s *GPool) SetLogger(l *log.Logger) {
	s.logger = l
}
