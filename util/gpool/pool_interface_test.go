// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/snail007/gmc/util/gpool"
	gatomic "github.com/snail007/gmc/util/sync/atomic"
	"github.com/stretchr/testify/assert"
)

// TestPool_BasicOperations 测试接口的基本操作
func TestPool_BasicOperations(t *testing.T) {
	tests := []struct {
		name    string
		factory func(int) gpool.Pool
	}{
		{
			name: "BasicPool",
			factory: func(count int) gpool.Pool {
				return gpool.New(count)
			},
		},
		{
			name: "OptimizedPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewOptimized(count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.factory(3)
			defer pool.Stop()

			// 测试 Submit
			done := make(chan bool)
			err := pool.Submit(func() {
				done <- true
			})
			assert.NoError(t, err)
			assert.True(t, <-done)

			// 测试 WorkerCount
			time.Sleep(time.Millisecond * 50)
			assert.True(t, pool.WorkerCount() > 0)

			// 测试 WaitDone
			counter := int32(0)
			for i := 0; i < 10; i++ {
				pool.Submit(func() {
					atomic.AddInt32(&counter, 1)
				})
			}
			pool.WaitDone()
			assert.Equal(t, int32(10), counter)
		})
	}
}

// TestPool_WorkerManagement 测试接口的工作协程管理
func TestPool_WorkerManagement(t *testing.T) {
	tests := []struct {
		name    string
		factory func(int) gpool.Pool
	}{
		{
			name: "BasicPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewWithPreAlloc(count)
			},
		},
		{
			name: "OptimizedPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewOptimized(count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.factory(3)
			defer pool.Stop()

			time.Sleep(time.Millisecond * 100)
			
			// 测试 Increase
			initialCount := pool.WorkerCount()
			pool.Increase(2)
			time.Sleep(time.Millisecond * 100)
			assert.Equal(t, initialCount+2, pool.WorkerCount())

			// 测试 Decrease
			pool.Decrease(1)
			time.Sleep(time.Millisecond * 100)
			assert.Equal(t, initialCount+1, pool.WorkerCount())
		})
	}
}

// TestPool_WorkerCounters 测试接口的计数器
func TestPool_WorkerCounters(t *testing.T) {
	tests := []struct {
		name    string
		factory func(int) gpool.Pool
	}{
		{
			name: "BasicPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewWithPreAlloc(count)
			},
		},
		{
			name: "OptimizedPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewOptimized(count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.factory(5)
			defer pool.Stop()

			time.Sleep(time.Millisecond * 100)
			
			// 所有worker应该空闲
			assert.Equal(t, 5, pool.IdleWorkerCount())
			assert.Equal(t, 0, pool.RunningWorkerCount())

			// 提交阻塞任务
			done := make(chan bool)
			for i := 0; i < 3; i++ {
				pool.Submit(func() {
					<-done
				})
			}

			time.Sleep(time.Millisecond * 100)
			
			// 3个worker应该在运行
			assert.Equal(t, 3, pool.RunningWorkerCount())
			assert.Equal(t, 2, pool.IdleWorkerCount())

			close(done)
			pool.WaitDone()
			
			time.Sleep(time.Millisecond * 100)
			// 所有worker应该回到空闲状态
			assert.Equal(t, 0, pool.RunningWorkerCount())
		})
	}
}

// TestPool_QueuedJobCount 测试队列任务计数
func TestPool_QueuedJobCount(t *testing.T) {
	tests := []struct {
		name    string
		factory func(int) gpool.Pool
	}{
		{
			name: "BasicPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewWithPreAlloc(count)
			},
		},
		{
			name: "OptimizedPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewOptimized(count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.factory(2)
			defer pool.Stop()

			time.Sleep(time.Millisecond * 50)

			// 占满所有worker
			done := make(chan bool)
			for i := 0; i < 2; i++ {
				pool.Submit(func() {
					<-done
				})
			}

			time.Sleep(time.Millisecond * 100)

			// 提交更多任务到队列
			for i := 0; i < 5; i++ {
				pool.Submit(func() {
					time.Sleep(time.Millisecond * 10)
				})
			}

			time.Sleep(time.Millisecond * 50)
			
			// 应该有任务在队列中
			queuedCount := pool.QueuedJobCount()
			assert.True(t, queuedCount > 0 && queuedCount <= 5)

			close(done)
			pool.WaitDone()
			
			// 所有任务完成后队列应该为空
			assert.Equal(t, 0, pool.QueuedJobCount())
		})
	}
}

// TestPool_ConcurrentSubmit 测试并发提交
func TestPool_ConcurrentSubmit(t *testing.T) {
	tests := []struct {
		name    string
		factory func(int) gpool.Pool
	}{
		{
			name: "BasicPool",
			factory: func(count int) gpool.Pool {
				return gpool.New(count)
			},
		},
		{
			name: "OptimizedPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewOptimized(count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.factory(10)
			defer pool.Stop()

			counter := gatomic.NewInt(0)
			goroutines := 20
			submitsPerGoroutine := 50

			// 并发提交任务
			done := make(chan bool, goroutines)
			for i := 0; i < goroutines; i++ {
				go func() {
					for j := 0; j < submitsPerGoroutine; j++ {
						pool.Submit(func() {
							counter.Increase(1)
						})
					}
					done <- true
				}()
			}

			// 等待所有提交完成
			for i := 0; i < goroutines; i++ {
				<-done
			}

			pool.WaitDone()

			expected := goroutines * submitsPerGoroutine
			assert.Equal(t, expected, counter.Val())
		})
	}
}

// TestPool_Stop 测试停止功能
func TestPool_Stop(t *testing.T) {
	tests := []struct {
		name    string
		factory func(int) gpool.Pool
	}{
		{
			name: "BasicPool",
			factory: func(count int) gpool.Pool {
				return gpool.New(count)
			},
		},
		{
			name: "OptimizedPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewOptimized(count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.factory(3)

			executed := int32(0)
			pool.Submit(func() {
				time.Sleep(time.Millisecond * 50)
				atomic.AddInt32(&executed, 1)
			})

			time.Sleep(time.Millisecond * 20)
			pool.Stop()

			// 多次Stop应该安全
			pool.Stop()
			pool.Stop()

			time.Sleep(time.Millisecond * 100)
			// 验证pool已停止
			assert.Equal(t, 0, pool.WorkerCount())
		})
	}
}

// TestPool_PerformanceComparison 性能对比测试
func TestPool_PerformanceComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能对比测试")
	}

	tests := []struct {
		name    string
		factory func(int) gpool.Pool
	}{
		{
			name: "BasicPool",
			factory: func(count int) gpool.Pool {
				return gpool.New(count)
			},
		},
		{
			name: "OptimizedPool",
			factory: func(count int) gpool.Pool {
				return gpool.NewOptimized(count)
			},
		},
	}

	jobCount := 10000
	workerCount := 50

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.factory(workerCount)
			defer pool.Stop()

			counter := int32(0)
			start := time.Now()

			for i := 0; i < jobCount; i++ {
				pool.Submit(func() {
					atomic.AddInt32(&counter, 1)
				})
			}

			pool.WaitDone()
			elapsed := time.Since(start)

			assert.Equal(t, int32(jobCount), counter)
			t.Logf("%s: 处理 %d 个任务耗时 %v", tt.name, jobCount, elapsed)
		})
	}
}

// TestPool_MaxJobCount 测试最大任务队列限制
func TestPool_MaxJobCount(t *testing.T) {
	// BasicPool 使用 New() 创建，需要手动设置 MaxJobCount
	t.Run("BasicPool", func(t *testing.T) {
		pool := gpool.New(1)
		pool.SetMaxJobCount(5)
		pool.SetBlocking(false)
		defer pool.Stop()

		// 占用唯一的worker
		done := make(chan bool)
		pool.Submit(func() {
			<-done
		})

		time.Sleep(time.Millisecond * 100)

		// 填满队列
		for i := 0; i < 5; i++ {
			err := pool.Submit(func() {
				time.Sleep(time.Millisecond)
			})
			assert.NoError(t, err)
		}

		// 队列满后应该返回错误
		err := pool.Submit(func() {})
		assert.Error(t, err)

		close(done)
		pool.WaitDone()
	})

	// OptimizedPool 通过 Option 设置
	t.Run("OptimizedPool", func(t *testing.T) {
		opt := &gpool.Option{
			MaxJobCount: 5,
			Blocking:    false,
		}
		pool := gpool.NewOptimizedWithOption(1, opt)
		defer pool.Stop()

		// 占用唯一的worker
		done := make(chan bool)
		pool.Submit(func() {
			<-done
		})

		time.Sleep(time.Millisecond * 100)

		// 填满队列
		for i := 0; i < 5; i++ {
			err := pool.Submit(func() {
				time.Sleep(time.Millisecond)
			})
			assert.NoError(t, err)
		}

		// 队列满后应该返回错误
		err := pool.Submit(func() {})
		assert.Error(t, err)

		close(done)
		pool.WaitDone()
	})
}

// TestPool_Polymorphism 测试多态性
func TestPool_Polymorphism(t *testing.T) {
	// 这个函数接受 Pool 接口，可以处理任何实现
	processWithPool := func(pool gpool.Pool, jobCount int) int32 {
		counter := int32(0)
		for i := 0; i < jobCount; i++ {
			pool.Submit(func() {
				atomic.AddInt32(&counter, 1)
			})
		}
		pool.WaitDone()
		return counter
	}

	// 使用 BasicPool
	basicPool := gpool.New(5)
	result1 := processWithPool(basicPool, 100)
	basicPool.Stop()
	assert.Equal(t, int32(100), result1)

	// 使用 OptimizedPool
	optimizedPool := gpool.NewOptimized(5)
	result2 := processWithPool(optimizedPool, 100)
	optimizedPool.Stop()
	assert.Equal(t, int32(100), result2)
}

// TestPool_DynamicSelection 测试动态选择实现
func TestPool_DynamicSelection(t *testing.T) {
	tests := []struct {
		name         string
		useOptimized bool
		jobCount     int
	}{
		{"BasicPool_SmallLoad", false, 10},
		{"OptimizedPool_SmallLoad", true, 10},
		{"BasicPool_MediumLoad", false, 100},
		{"OptimizedPool_MediumLoad", true, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 根据条件动态选择实现
			var pool gpool.Pool
			if tt.useOptimized {
				pool = gpool.NewOptimized(10)
			} else {
				pool = gpool.New(10)
			}
			defer pool.Stop()

			counter := int32(0)
			for i := 0; i < tt.jobCount; i++ {
				pool.Submit(func() {
					atomic.AddInt32(&counter, 1)
				})
			}

			pool.WaitDone()
			assert.Equal(t, int32(tt.jobCount), counter)
			t.Logf("完成 %d 个任务", counter)
		})
	}
}
