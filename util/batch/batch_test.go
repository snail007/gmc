package gbatch

import (
	"context"
	"errors"
	gloop "github.com/snail007/gmc/util/loop"
	gatomic "github.com/snail007/gmc/util/sync/atomic"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewBatchExecutor(t *testing.T) {
	t.Parallel()
	i := gatomic.NewInt(0)
	be := NewBatchExecutor()
	be.SetWorkers(10)
	start := time.Now()
	gloop.For(10, func(idx int) {
		be.AppendTask(func(_ context.Context) (interface{}, func(), error) {
			i.Increase(idx)
			time.Sleep(time.Second)
			return "okay", nil, nil
		})
	})
	rs := be.WaitAll()
	assert.Equal(t, len(rs), 10)
	assert.Equal(t, "okay", rs[0].Value())
	assert.Nil(t, rs[0].Err())
	diff := time.Now().Sub(start)
	assert.Equal(t, 45, i.Val())
	assert.True(t, diff >= time.Second)
	assert.True(t, diff < time.Millisecond*1500)
}

func TestNewBatchExecutor2(t *testing.T) {
	t.Parallel()
	i := gatomic.NewInt(0)
	be := NewBatchExecutor()
	be.SetWorkers(5)
	start := time.Now()
	gloop.For(10, func(idx int) {
		be.AppendTask(func(_ context.Context) (interface{}, func(), error) {
			i.Increase(idx)
			time.Sleep(time.Second)
			return nil, nil, nil
		})
	})
	be.WaitAll()
	diff := time.Now().Sub(start)
	assert.Equal(t, 45, i.Val())
	assert.True(t, diff >= time.Second*2)
	assert.True(t, diff < time.Millisecond*2500)
}

func TestWaitFirstSuccess(t *testing.T) {
	// 创建一个包含两个成功任务和一个失败任务的 Executor
	executor := NewBatchExecutor()
	executor.AppendTask(
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(100 * time.Millisecond) // Simulate some work
			return "Task 1 result", nil, nil
		},
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(200 * time.Millisecond) // Simulate some work
			return "Task 2 result", nil, nil
		},
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 3 error")
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
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 1 error")
		},
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 2 error")
		},
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 3 error")
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
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(200 * time.Millisecond) // Simulate some work
			return "Task 1 result", nil, nil
		},
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 2 error")
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
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(100 * time.Second) // Simulate some work
			return "Task 1 result", nil, nil
		},
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(200 * time.Millisecond) // Simulate some work
			return "Task 2 result", nil, nil
		},
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(200 * time.Second) // Simulate some work
			return nil, nil, errors.New("Task 3 error")
		},
	)

	value, err := executor.WaitFirstDone()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
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
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 1 error")
		},
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 2 error")
		},
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 3 error")
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
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(200 * time.Millisecond) // Simulate some work
			return "Task 1 result", func() {
			}, nil
		},
		func(ctx context.Context) (interface{}, func(), error) {
			return nil, nil, errors.New("Task 2 error")
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
	canceled := false
	executor.AppendTask(
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(time.Second * 2)
			return "Task 2 result", func() {
				canceled = true
			}, nil
		},
		func(ctx context.Context) (interface{}, func(), error) {
			time.Sleep(time.Second)
			return "Task 1 result", nil, nil
		})

	value, err := executor.WaitFirstSuccess()

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	expectedValue := "Task 1 result"
	if value != expectedValue {
		t.Errorf("Expected value %q, but got: %q", expectedValue, value)
	}
	time.Sleep(time.Second * 3)
	assert.True(t, canceled)
}
