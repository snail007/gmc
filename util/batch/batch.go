package gbatch

import (
	"context"
	"fmt"
	gerror "github.com/snail007/gmc/module/error"
	glist "github.com/snail007/gmc/util/list"
	gsync "github.com/snail007/gmc/util/sync"
	"sync"
)

type task func(ctx context.Context) (value interface{}, cancelFunc func(), err error)
type Executor struct {
	workers      int
	tasks        []task
	rootCtx      context.Context
	rootCancel   context.CancelFunc
	panicHandler func(e interface{})
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

func (s *Executor) SetPanicHandler(panicHandler func(e interface{})) {
	s.panicHandler = panicHandler
}

func (s *Executor) WaitAll() (allResults []taskResult) {
	allResult := glist.New()
	g := sync.WaitGroup{}
	g.Add(len(s.tasks))
	for _, t := range s.tasks {
		go func(t task) {
			defer g.Done()
			item := s.run(t)
			allResult.Add(item)
		}(t)
	}
	g.Wait()
	allResult.RangeFast(func(index int, value interface{}) bool {
		allResults = append(allResults, value.(taskResult))
		return true
	})
	return
}

func (s *Executor) run(fn task) (result taskResult) {
	defer func() {
		if e := recover(); e != nil {
			err := gerror.Wrap(e)
			result.err = err
			msg := fmt.Sprintf("[WARN] task panic, err: %s", err.StackStr())
			if s.panicHandler != nil {
				s.panicHandler(e)
			} else {
				fmt.Println(msg)
			}
		}
	}()
	v, f, e := fn(s.rootCtx)
	result.value = v
	result.cancelFunc = f
	result.err = e
	return
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
	value, err = s.waitFirst(true)
	return
}

// WaitFirstDone wait first done, value and err is the task returned.
func (s *Executor) WaitFirstDone() (value interface{}, err error) {
	return s.waitFirst(false)
}

func (s *Executor) waitFirst(checkSuccess bool) (value interface{}, err error) {
	g := sync.WaitGroup{}
	g.Add(len(s.tasks))
	waitChan := make(chan taskResult)
	allResult := glist.New()
	for _, t := range s.tasks {
		go func(t task) {
			defer g.Done()
			item := s.run(t)
			allResult.Add(item)
			if checkSuccess && item.err != nil {
				return
			}
			select {
			case waitChan <- item:
			default:
				if item.cancelFunc != nil {
					item.cancelFunc()
				}
			}
		}(t)
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
