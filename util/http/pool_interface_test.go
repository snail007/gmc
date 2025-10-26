package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/snail007/gmc/util/gpool"
)

// TestPoolInterfaceCompatibility verifies both Pool implementations work with BatchRequest
func TestPoolInterfaceCompatibility(t *testing.T) {
	// 编译时检查类型兼容性
	var _ gpool.Pool = (*gpool.BasicPool)(nil)
	var _ gpool.Pool = (*gpool.OptimizedPool)(nil)

	t.Run("OriginalPool", func(t *testing.T) {
		pool := gpool.New(2)
		defer pool.Stop()

		// 验证可以传入 Pool
		testWithPool(t, pool)
	})

	t.Run("OptimizedPool", func(t *testing.T) {
		pool := gpool.NewOptimized(2)
		defer pool.Stop()

		// 验证可以传入 OptimizedPool
		testWithPool(t, pool)
	})
}

// testWithPool 是一个接受 gpool.Pool 的通用测试函数
func testWithPool(t *testing.T, pool gpool.Pool) {
	// 提交一些任务到池中
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		err := pool.Submit(func() {
			time.Sleep(10 * time.Millisecond)
			done <- true
		})
		if err != nil {
			t.Errorf("Submit failed: %v", err)
		}
	}

	// 等待完成
	pool.WaitDone()

	// 验证所有任务都完成了
	count := 0
	timeout := time.After(1 * time.Second)
loop:
	for {
		select {
		case <-done:
			count++
			if count == 10 {
				break loop
			}
		case <-timeout:
			t.Errorf("Timeout waiting for tasks, completed: %d/10", count)
			break loop
		}
	}

	if count != 10 {
		t.Errorf("Expected 10 tasks completed, got %d", count)
	}
}

// ExamplePoolInterface demonstrates using gpool.Pool for flexibility
func ExamplePoolInterface() {
	// 函数接受接口类型，可以使用任意实现
	processJobs := func(pool gpool.Pool) {
		for i := 0; i < 5; i++ {
			pool.Submit(func() {
				// 执行任务
			})
		}
		pool.WaitDone()
	}

	// 使用标准 Pool
	pool1 := gpool.New(2)
	processJobs(pool1)
	pool1.Stop()

	// 使用优化版 Pool
	pool2 := gpool.NewOptimized(2)
	processJobs(pool2)
	pool2.Stop()

	fmt.Println("Both implementations work!")
	// Output: Both implementations work!
}

