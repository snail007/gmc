package gbatch

import (
	"context"
	"fmt"
	gerror "github.com/snail007/gmc/module/error"
	"github.com/snail007/gmc/util/gpool"
	glist "github.com/snail007/gmc/util/list"
	gsync "github.com/snail007/gmc/util/sync"
	"sync"
)

type task func(ctx context.Context) (value interface{}, cancelFunc func(), err error)
type Executor struct {
	workers    int
	tasks      []task
	rootCtx    context.Context
	rootCancel context.CancelFunc
}

func NewBatchExecutor() *Executor {
	r, c := context.WithCancel(context.Background())
	return &Executor{
		rootCtx:    r,
		rootCancel: c,
		workers:    10,
	}
}

func (s *Executor) SetWorkers(workersCnt int) {
	s.workers = workersCnt
}

func (s *Executor) AppendTask(tasks ...task) {
	s.tasks = append(s.tasks, tasks...)
}

func (s *Executor) WaitAll() (allResults []taskResult) {
	p := s.getPool()
	defer p.Stop()
	allResult := glist.New()
	g := sync.WaitGroup{}
	g.Add(len(s.tasks))
	for _, t := range s.tasks {
		task := t
		p.Submit(func() {
			defer g.Done()
			v, f, e := task(s.rootCtx)
			allResult.Add(taskResult{value: v, err: e, cancelFunc: f})
		})
	}
	g.Wait()
	allResult.RangeFast(func(index int, value interface{}) bool {
		allResults = append(allResults, value.(taskResult))
		return true
	})
	return
}
func (s *Executor) getPool() *gpool.Pool {
	workers := len(s.tasks)
	if s.workers > 0 {
		workers = s.workers
	}
	return gpool.NewWithPreAlloc(workers)
}

type taskResult struct {
	value      interface{}
	err        error
	cancelFunc func()
}

func (t taskResult) Value() interface{} {
	return t.value
}

func (t taskResult) Err() error {
	return t.err
}

// WaitFirstSuccess wait first success, value is the first success task's returned value, if all task fail, err is the last error.
func (s *Executor) WaitFirstSuccess() (value interface{}, err error) {
	return s.waitFirst(true)
}

// WaitFirstDone wait first done, value and err is the task returned.
func (s *Executor) WaitFirstDone() (value interface{}, err error) {
	return s.waitFirst(false)
}

func (s *Executor) waitFirst(checkSuccess bool) (value interface{}, err error) {
	p := s.getPool()
	defer p.Stop()
	g := sync.WaitGroup{}
	g.Add(len(s.tasks))
	waitChan := make(chan taskResult)
	allResult := glist.New()
	for _, t := range s.tasks {
		task := t
		p.Submit(func() {
			defer g.Done()
			var v interface{}
			var f func()
			var e error
			if err := gerror.TryWithStack(func() {
				v, f, e = task(s.rootCtx)
			}); err != nil {
				e = err
				fmt.Println("[WARN] run task panic, error:" + err.String())
			}
			item := taskResult{value: v, err: e, cancelFunc: f}
			allResult.Add(item)
			if checkSuccess && e != nil {
				return
			}
			select {
			case waitChan <- item:
			default:
				if f != nil {
					f()
				}
			}
		})
	}
	select {
	case v := <-waitChan:
		//a task returned, call rootCancel to cancel others task.
		go s.rootCancel()
		return v.value, v.err
	case <-gsync.WaitGroupToChan(&g):
		// all task done, return the last err
		return nil, allResult.Pop().(taskResult).err
	}
}
