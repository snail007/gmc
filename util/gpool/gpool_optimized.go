// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
)

var (
	ErrMaxQueuedJobCountReachedOptimized = errors.New("max queued job count reached")
)

// OptimizedPool is an optimized goroutine pool with better performance
type OptimizedPool struct {
	maxWorkCount         int32
	workers              []*optimizedWorker
	workersMutex         sync.RWMutex
	jobQueue             chan *OptimizedJobItem
	g                    *sync.WaitGroup
	opt                  *Option
	idleWorkerCounter    int64
	runningWorkerCounter int64
	stopped              int32
	nextWorkerID         int64
}

type OptimizedJobItem struct {
	Stack string
	Job   func()
}

// NewOptimized creates an optimized pool with better performance
func NewOptimized(workerCount int) *OptimizedPool {
	return NewOptimizedWithOption(workerCount, &Option{
		WithStack: false, // 默认关闭stack trace以提升性能
	})
}

// NewOptimizedWithOption creates an optimized pool with options
func NewOptimizedWithOption(workerCount int, opt *Option) *OptimizedPool {
	if opt.MaxJobCount == 0 {
		opt.MaxJobCount = 10000 // 默认队列大小
	}

	p := &OptimizedPool{
		maxWorkCount: int32(workerCount),
		workers:      make([]*optimizedWorker, 0, workerCount),
		jobQueue:     make(chan *OptimizedJobItem, opt.MaxJobCount),
		g:            &sync.WaitGroup{},
		opt:          opt,
	}

	// 预分配workers
	if opt.PreAlloc || workerCount > 0 {
		for i := 0; i < workerCount; i++ {
			p.addWorker()
		}
	}

	return p
}

// Submit adds a job to the queue (optimized version with no global lock)
func (p *OptimizedPool) Submit(job func()) error {
	if atomic.LoadInt32(&p.stopped) == 1 {
		return errors.New("pool is stopped")
	}

	j := &OptimizedJobItem{
		Job: job,
	}

	// 只在需要时获取stack trace
	if p.opt.WithStack {
		// 使用runtime.Caller替代debug.Stack，性能更好
		if _, file, line, ok := runtime.Caller(1); ok {
			j.Stack = fmt.Sprintf("%s:%d", file, line)
		}
	}

	p.g.Add(1)

	// 使用channel的非阻塞发送，避免全局锁
	select {
	case p.jobQueue <- j:
		return nil
	default:
		// 队列满了
		p.g.Done()
		if p.opt.Blocking {
			// 阻塞模式：等待队列有空位
			p.jobQueue <- j
			p.g.Add(1)
			return nil
		}
		return ErrMaxQueuedJobCountReachedOptimized
	}
}

// WorkerCount returns the count of workers
func (p *OptimizedPool) WorkerCount() int {
	p.workersMutex.RLock()
	defer p.workersMutex.RUnlock()
	return len(p.workers)
}

// RunningWorkerCount returns the count of running workers
func (p *OptimizedPool) RunningWorkerCount() int {
	return int(atomic.LoadInt64(&p.runningWorkerCounter))
}

// IdleWorkerCount returns the count of idle workers
func (p *OptimizedPool) IdleWorkerCount() int {
	return int(atomic.LoadInt64(&p.idleWorkerCounter))
}

// QueuedJobCount returns the count of queued jobs
func (p *OptimizedPool) QueuedJobCount() int {
	return len(p.jobQueue)
}

// WaitDone waits for all jobs to complete
func (p *OptimizedPool) WaitDone() {
	p.g.Wait()
}

// Stop stops all workers
func (p *OptimizedPool) Stop() {
	if !atomic.CompareAndSwapInt32(&p.stopped, 0, 1) {
		return
	}

	close(p.jobQueue)

	p.workersMutex.Lock()
	for _, w := range p.workers {
		w.stop()
	}
	p.workers = nil
	p.workersMutex.Unlock()
}

// Increase adds workers (optimized version)
func (p *OptimizedPool) Increase(count int) {
	atomic.AddInt32(&p.maxWorkCount, int32(count))
	for i := 0; i < count; i++ {
		p.addWorker()
	}
}

// Decrease removes workers (optimized version)
func (p *OptimizedPool) Decrease(count int) {
	atomic.AddInt32(&p.maxWorkCount, -int32(count))

	p.workersMutex.Lock()
	defer p.workersMutex.Unlock()

	for i := 0; i < count && i < len(p.workers); i++ {
		w := p.workers[len(p.workers)-1]
		p.workers = p.workers[:len(p.workers)-1]
		w.stop()
	}
}

func (p *OptimizedPool) addWorker() {
	w := newOptimizedWorker(p)
	p.workersMutex.Lock()
	p.workers = append(p.workers, w)
	p.workersMutex.Unlock()
	w.start()
}

func (p *OptimizedPool) run(j *OptimizedJobItem) {
	defer func() {
		p.g.Done()
		if e := recover(); e != nil {
			msg := fmt.Sprintf("Pool: a job stopped unexpectedly, err: %s", gcore.ProviderError()().StackError(e))
			if p.opt.Logger != nil {
				p.opt.Logger.Error(msg)
			} else {
				fmt.Println(msg)
			}
			if p.opt.PanicHandler != nil {
				p.opt.PanicHandler(e)
			}
		}
	}()
	j.Job()
}

type optimizedWorker struct {
	id       int64
	pool     *OptimizedPool
	stopChan chan struct{}
	stopped  int32
}

func newOptimizedWorker(pool *OptimizedPool) *optimizedWorker {
	// 使用atomic counter生成ID，比crypto/rand快得多
	id := atomic.AddInt64(&pool.nextWorkerID, 1)
	return &optimizedWorker{
		id:       id,
		pool:     pool,
		stopChan: make(chan struct{}),
	}
}

func (w *optimizedWorker) start() {
	go func() {
		defer func() {
			gerror.RecoverNop()
		}()

		var idleTimer *time.Timer
		if w.pool.opt.IdleDuration > 0 {
			idleTimer = time.NewTimer(w.pool.opt.IdleDuration)
			defer idleTimer.Stop()
		}

		for {
			// 使用atomic操作保证并发安全
			atomic.AddInt64(&w.pool.idleWorkerCounter, 1)

			var job *OptimizedJobItem
			var ok bool

			if idleTimer != nil {
				idleTimer.Reset(w.pool.opt.IdleDuration)
				select {
				case job, ok = <-w.pool.jobQueue:
					if !ok {
						atomic.AddInt64(&w.pool.idleWorkerCounter, -1)
						return
					}
				case <-idleTimer.C:
					atomic.AddInt64(&w.pool.idleWorkerCounter, -1)
					// 空闲超时，退出
					w.pool.workersMutex.Lock()
					for i, worker := range w.pool.workers {
						if worker == w {
							w.pool.workers = append(w.pool.workers[:i], w.pool.workers[i+1:]...)
							break
						}
					}
					w.pool.workersMutex.Unlock()
					return
				case <-w.stopChan:
					atomic.AddInt64(&w.pool.idleWorkerCounter, -1)
					return
				}
			} else {
				select {
				case job, ok = <-w.pool.jobQueue:
					if !ok {
						atomic.AddInt64(&w.pool.idleWorkerCounter, -1)
						return
					}
				case <-w.stopChan:
					atomic.AddInt64(&w.pool.idleWorkerCounter, -1)
					return
				}
			}

			atomic.AddInt64(&w.pool.idleWorkerCounter, -1)
			atomic.AddInt64(&w.pool.runningWorkerCounter, 1)

			w.pool.run(job)

			atomic.AddInt64(&w.pool.runningWorkerCounter, -1)
		}
	}()
}

func (w *optimizedWorker) stop() {
	if !atomic.CompareAndSwapInt32(&w.stopped, 0, 1) {
		return
	}
	close(w.stopChan)
}
