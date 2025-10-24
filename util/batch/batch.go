package gbatch

import (
	"context"
	"errors"
	"fmt"

	gerror "github.com/snail007/gmc/module/error"
	"github.com/snail007/gmc/util/gpool"
	glist "github.com/snail007/gmc/util/list"
)

type task func(ctx context.Context) (value interface{}, err error)
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
	p := gpool.New(s.workers)
	defer p.Stop()
	allResult := glist.New()
	for _, t := range s.tasks {
		t0 := t
		p.Submit(func() {
			item := s.run(t0)
			allResult.Add(item)
		})
	}
	p.WaitDone()
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
			msg := fmt.Sprintf("[WARN] task panic, err: %s", err.ErrorStack())
			if s.panicHandler != nil {
				s.panicHandler(e)
			} else {
				fmt.Println(msg)
			}
		}
	}()
	v, e := fn(s.rootCtx)
	result.value = v
	result.err = e
	return
}

type taskResult struct {
	value interface{}
	err   error
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
	if len(s.tasks) == 0 {
		return nil, errors.New("tasks is empty")
	}
	allResultChan := make(chan taskResult, len(s.tasks))
	for _, t := range s.tasks {
		t0 := t
		go func() {
			if s.rootCtx.Err() != nil {
				return
			}
			allResultChan <- s.run(t0)
		}()
	}
	cnt := 0
	for item := range allResultChan {
		cnt++
		if checkSuccess {
			if item.err == nil {
				go s.rootCancel()
				return item.value, nil
			} else if cnt == len(s.tasks) {
				return item.value, item.err
			}
		} else {
			go s.rootCancel()
			return item.value, item.err
		}
	}
	return nil, errors.New("failed task not found, this should not happen")
}
