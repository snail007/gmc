package gbatch

import (
	"context"
	"errors"
	"testing"
	"time"

	gloop "github.com/snail007/gmc/util/loop"
	gatomic "github.com/snail007/gmc/util/sync/atomic"
	"github.com/stretchr/testify/assert"
)

func TestNewBatchExecutor(t *testing.T) {
	t.Parallel()
	i := gatomic.NewInt(0)
	be := NewBatchExecutor()
	be.SetWorkers(10)
	start := time.Now()
	gloop.For(10, func(idx int) {
		be.AppendTask(func(_ context.Context) (interface{}, error) {
			i.Increase(idx)
			time.Sleep(time.Second)
			return "okay", nil
		})
	})
	rs := be.WaitAll()
	assert.Equal(t, len(rs), 10)
	assert.Equal(t, "okay", rs[0].Value())
	assert.Nil(t, rs[0].Err())
	diff := time.Now().Sub(start)
	assert.Equal(t, 45, i.Val())
	t.Log(diff)
	assert.True(t, diff >= time.Second)
	assert.True(t, diff < time.Second*10)
}

func TestNewBatchExecutor2(t *testing.T) {
	t.Parallel()
	i := gatomic.NewInt(0)
	be := NewBatchExecutor()
	be.SetWorkers(5)
	start := time.Now()
	gloop.For(10, func(idx int) {
		be.AppendTask(func(_ context.Context) (interface{}, error) {
			i.Increase(idx)
			time.Sleep(time.Second)
			return nil, nil
		})
	})
	be.WaitAll()
	diff := time.Now().Sub(start)
	assert.Equal(t, 45, i.Val())
	t.Log(diff)
	assert.True(t, diff >= time.Second)
	assert.True(t, diff < time.Second*5)
}

func TestWaitFirstSuccess(t *testing.T) {
	// 创建一个包含两个成功任务和一个失败任务的 Executor
	executor := NewBatchExecutor()
	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(100 * time.Millisecond) // Simulate some work
			return "Task 1 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 3 error")
		},
	)

	value, err := executor.WaitFirstSuccess()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedValue := "Task 1 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}
}

func TestWaitFirstSuccess_AllFail(t *testing.T) {
	// 创建一个只包含失败任务的 Executor
	executor := NewBatchExecutor()
	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 1 error")
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 2 error")
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 3 error")
		},
	)

	value, err := executor.WaitFirstSuccess()

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
	if value != nil {
		t.Errorf("Expected value to be nil, but got: %v", value)
	}
}

func TestWaitFirstSuccess_PartialSuccess(t *testing.T) {
	// 创建一个包含一个成功任务和一个失败任务的 Executor
	executor := NewBatchExecutor()
	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(200 * time.Millisecond) // Simulate some work
			return "Task 1 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 2 error")
		},
	)

	value, err := executor.WaitFirstSuccess()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedValue := "Task 1 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}
}

func TestWaitFirstDone(t *testing.T) {
	// 创建一个包含两个成功任务和一个失败任务的 Executor
	executor := NewBatchExecutor()
	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(10 * time.Second) // Simulate some work
			return "Task 1 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(200 * time.Millisecond) // Simulate some work
			return "Task 2 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(20 * time.Second) // Simulate some work
			return nil, errors.New("Task 3 error")
		},
	)

	value, err := executor.WaitFirstDone()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
		return
	}
	expectedValue := "Task 2 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}
}

func TestWaitFirstDone_AllFail(t *testing.T) {
	// 创建一个只包含失败任务的 Executor
	executor := NewBatchExecutor()
	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 1 error")
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 2 error")
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 3 error")
		},
	)

	value, err := executor.WaitFirstDone()

	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
	if value != nil {
		t.Errorf("Expected value to be nil, but got: %v", value)
	}
}

func TestWaitFirstDone_PartialSuccess(t *testing.T) {
	// 创建一个包含一个成功任务和一个失败任务的 Executor
	executor := NewBatchExecutor()
	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(200 * time.Millisecond) // Simulate some work
			return "Task 1 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			return nil, errors.New("Task 2 error")
		},
	)

	value, err := executor.WaitFirstSuccess()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedValue := "Task 1 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}
}

func TestCancelFunc(t *testing.T) {
	t.Parallel()
	executor := NewBatchExecutor()
	task1Canceled := false
	task2Canceled := false
	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			go func() {
				select {
				case <-ctx.Done():
					// 检查是否是第一个成功的
					if !IsFirstSuccess(ctx) {
						task1Canceled = true
					}
				}
			}()
			time.Sleep(time.Second * 2)
			return "Task 1 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			go func() {
				select {
				case <-ctx.Done():
					// 检查是否是第一个成功的
					if !IsFirstSuccess(ctx) {
						task2Canceled = true
					}
				}
			}()
			time.Sleep(time.Second)
			return "Task 2 result", nil
		})

	value, err := executor.WaitFirstSuccess()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedValue := "Task 2 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}
	time.Sleep(time.Second * 3)
	// Task 2 是第一个完成的，不应该被标记为 canceled
	assert.False(t, task2Canceled, "First completed task should not be canceled")
	// Task 1 不是第一个完成的，应该被取消
	assert.True(t, task1Canceled, "Non-first task should be canceled")
}

func TestPanic(t *testing.T) {
	t.Parallel()
	called := false
	executor := NewBatchExecutor()
	executor.SetPanicHandler(func(e interface{}) {
		called = true
	})
	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			time.Sleep(time.Second * 2)
			return "Task 2 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			panic("abc")
			return "Task 1 result", nil
		})

	value, err := executor.WaitFirstSuccess()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedValue := "Task 2 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}
	assert.True(t, called)
	time.Sleep(time.Second * 3)
}

func TestCancelFuncForFirstCompletedTask(t *testing.T) {
	t.Parallel()
	executor := NewBatchExecutor()
	firstTaskCanceledCount := 0
	secondTaskCanceledCount := 0

	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			go func() {
				select {
				case <-ctx.Done():
					// 检查是否是第一个完成的
					if !IsFirstSuccess(ctx) {
						firstTaskCanceledCount++
					}
				}
			}()
			time.Sleep(time.Second)
			return "Task 1 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			go func() {
				select {
				case <-ctx.Done():
					// 检查是否是第一个完成的
					if !IsFirstSuccess(ctx) {
						secondTaskCanceledCount++
					}
				}
			}()
			time.Sleep(time.Second * 3)
			return "Task 2 result", nil
		})

	value, err := executor.WaitFirstSuccess()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedValue := "Task 1 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}

	time.Sleep(time.Second * 2)

	// 第一个完成的 task 不应该执行取消逻辑
	assert.Equal(t, 0, firstTaskCanceledCount, "First completed task should not execute cancel logic")
	// 第二个 task 应该执行取消逻辑
	assert.Equal(t, 1, secondTaskCanceledCount, "Second task should execute cancel logic")
}

func TestIsFirstDone(t *testing.T) {
	t.Parallel()
	executor := NewBatchExecutor()
	task1Cancelled := false
	task2Cancelled := false

	executor.AppendTask(
		func(ctx context.Context) (interface{}, error) {
			go func() {
				<-ctx.Done()
				if !IsFirstDone(ctx) {
					task1Cancelled = true
				}
			}()
			time.Sleep(time.Second * 2)
			return "Task 1 result", nil
		},
		func(ctx context.Context) (interface{}, error) {
			go func() {
				<-ctx.Done()
				if !IsFirstDone(ctx) {
					task2Cancelled = true
				}
			}()
			time.Sleep(100 * time.Millisecond)
			return "Task 2 result", nil
		})

	value, err := executor.WaitFirstDone()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedValue := "Task 2 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}

	time.Sleep(time.Second * 3)

	// Task 2 是第一个完成的，不应该被标记为 cancelled
	assert.False(t, task2Cancelled, "First done task should not be marked as cancelled")
	// Task 1 不是第一个完成的，应该被取消
	assert.True(t, task1Cancelled, "Non-first task should be cancelled")
}

func TestWaitAllResultsOrder(t *testing.T) {
	t.Parallel()
	executor := NewBatchExecutor()
	executor.SetWorkers(5)
	
	// 添加10个任务，每个任务返回它的索引
	// 用不同的延迟来确保完成顺序是随机的
	for i := 0; i < 10; i++ {
		idx := i
		executor.AppendTask(func(_ context.Context) (interface{}, error) {
			// 让后面的任务更快完成
			time.Sleep(time.Millisecond * time.Duration(100 - idx*10))
			return idx, nil
		})
	}
	
	results := executor.WaitAll()
	
	// 验证结果数量
	assert.Equal(t, 10, len(results), "Results length should be 10")
	
	// 验证每个结果的值与索引对应
	for i, r := range results {
		assert.Nil(t, r.Err(), "Task %d should have no error", i)
		assert.Equal(t, i, r.Value().(int), "Task %d result should be %d", i, i)
	}
}
