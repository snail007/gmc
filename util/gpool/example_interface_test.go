package gpool_test

import (
	"fmt"
	"time"

	"github.com/snail007/gmc/util/gpool"
)

// ExamplePoolInterface 演示接口的灵活性
func ExamplePoolInterface() {
	// 定义一个接受接口的函数
	processJobs := func(pool gpool.Pool, jobCount int) {
		for i := 0; i < jobCount; i++ {
			pool.Submit(func() {
				time.Sleep(10 * time.Millisecond)
			})
		}
		
		// 接口提供的监控方法
		fmt.Printf("Total workers: %d\n", pool.WorkerCount())
		pool.WaitDone()
	}

	// 场景1：使用标准 Pool
	pool1 := gpool.New(5)
	processJobs(pool1, 10)
	pool1.Stop()

	// 场景2：使用优化版 OptimizedPool
	pool2 := gpool.NewOptimized(5)
	processJobs(pool2, 10)
	pool2.Stop()

	fmt.Println("Both implementations work!")
	// Output:
	// Total workers: 5
	// Total workers: 5
	// Both implementations work!
}

// ExamplePoolInterface_switchImplementation 演示如何轻松切换实现
func ExamplePoolInterface_switchImplementation() {
	// 配置决定使用哪个实现
	useOptimized := true

	var pool gpool.Pool
	if useOptimized {
		pool = gpool.NewOptimized(10)
	} else {
		pool = gpool.New(10)
	}
	defer pool.Stop()

	// 使用池处理任务，无需关心具体实现
	for i := 0; i < 5; i++ {
		pool.Submit(func() {
			fmt.Println("Task executed")
		})
	}

	pool.WaitDone()
	fmt.Println("Done")
}

// ExamplePoolInterface_dynamicScaling 演示接口支持的动态扩缩容
func ExamplePoolInterface_dynamicScaling() {
	pool := gpool.NewOptimized(2)
	defer pool.Stop()

	fmt.Printf("Initial workers: %d\n", pool.WorkerCount())

	// 增加工作协程
	pool.Increase(3)
	fmt.Printf("After increase: %d\n", pool.WorkerCount())

	// 减少工作协程
	pool.Decrease(2)
	fmt.Printf("After decrease: %d\n", pool.WorkerCount())

	// Output:
	// Initial workers: 2
	// After increase: 5
	// After decrease: 3
}
