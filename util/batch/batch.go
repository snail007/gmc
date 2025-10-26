package gbatch

import (
	"context"
	"errors"
	"fmt"

	gerror "github.com/snail007/gmc/module/error"
	"github.com/snail007/gmc/util/gpool"
)

type taskStatusKey struct{}

// IsFirstSuccess 检查当前 task 是否是第一个成功的
// 仅在 WaitFirstSuccess 中有效
// 返回 true 表示当前 task 是第一个成功的，不应该执行失败清理逻辑
// 返回 false 表示当前 task 不是第一个成功的，应该执行清理逻辑
func IsFirstSuccess(ctx context.Context) bool {
	if v := ctx.Value(taskStatusKey{}); v != nil {
		if isFirst, ok := v.(*bool); ok {
			return *isFirst
		}
	}
	return false
}

// IsFirstDone 检查当前 task 是否是第一个完成的
// 仅在 WaitFirstDone 中有效
// 返回 true 表示当前 task 是第一个完成的，不应该执行失败清理逻辑
// 返回 false 表示当前 task 不是第一个完成的，应该执行清理逻辑
func IsFirstDone(ctx context.Context) bool {
	// 实现与 IsFirstSuccess 相同，直接调用
	return IsFirstSuccess(ctx)
}

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

// WaitAll 等待所有任务完成并返回结果
// 返回的结果数组与任务添加顺序严格对应，allResults[i] 对应 tasks[i]
func (s *Executor) WaitAll() (allResults []taskResult) {
	p := gpool.New(s.workers)
	defer p.Stop()
	allResults = make([]taskResult, len(s.tasks))
	for idx, t := range s.tasks {
		t0 := t
		idx0 := idx
		p.Submit(func() {
			allResults[idx0] = s.run(t0)
		})
	}
	p.WaitDone()
	return
}

func (s *Executor) run(fn task) (result taskResult) {
	return s.runWithContext(s.rootCtx, fn)
}

func (s *Executor) runWithContext(ctx context.Context, fn task) (result taskResult) {
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
	v, e := fn(ctx)
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
	type taskResultWithIndex struct {
		taskResult
		index int
	}
	allResultChan := make(chan taskResultWithIndex, len(s.tasks))
	isFirstFlags := make([]*bool, len(s.tasks))

	for i, t := range s.tasks {
		t0 := t
		idx := i
		isFirst := false
		isFirstFlags[idx] = &isFirst
		taskCtx := context.WithValue(s.rootCtx, taskStatusKey{}, isFirstFlags[idx])
		taskCtx, _ = context.WithCancel(taskCtx)
		go func() {
			if taskCtx.Err() != nil {
				return
			}
			result := s.runWithContext(taskCtx, t0)
			allResultChan <- taskResultWithIndex{taskResult: result, index: idx}
		}()
	}
	cnt := 0
	for item := range allResultChan {
		cnt++
		if checkSuccess {
			if item.err == nil {
				// 标记第一个成功的 task
				*isFirstFlags[item.index] = true
				// 取消 rootCtx 会自动取消所有派生的 task context
				s.rootCancel()
				return item.value, nil
			} else if cnt == len(s.tasks) {
				// 所有任务都失败了，取消所有 task
				s.rootCancel()
				return item.value, item.err
			}
		} else {
			// 标记第一个完成的 task
			*isFirstFlags[item.index] = true
			// 取消 rootCtx 会自动取消所有派生的 task context
			s.rootCancel()
			return item.value, item.err
		}
	}
	return nil, errors.New("failed task not found, this should not happen")
}
