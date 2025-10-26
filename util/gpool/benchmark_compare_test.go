// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool

import (
	"sync"
	"testing"
)

// 对比原版和优化版的性能

func BenchmarkOriginal_Submit(b *testing.B) {
	p := New(100)
	defer p.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func() {})
	}
}

func BenchmarkOptimized_Submit(b *testing.B) {
	p := NewOptimized(100)
	defer p.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func() {})
	}
}

func BenchmarkOriginal_SubmitAndWait(b *testing.B) {
	p := New(100)
	defer p.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func() {})
	}
	p.WaitDone()
}

func BenchmarkOptimized_SubmitAndWait(b *testing.B) {
	p := NewOptimized(100)
	defer p.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func() {})
	}
	p.WaitDone()
}

// 并发提交测试
func BenchmarkOriginal_ConcurrentSubmit(b *testing.B) {
	p := New(100)
	defer p.Stop()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Submit(func() {})
		}
	})
}

func BenchmarkOptimized_ConcurrentSubmit(b *testing.B) {
	p := NewOptimized(100)
	defer p.Stop()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Submit(func() {})
		}
	})
}

// 不同worker数量的对比
func BenchmarkOriginal_Workers10(b *testing.B) {
	benchmarkPoolSubmit(b, New(10))
}

func BenchmarkOptimized_Workers10(b *testing.B) {
	benchmarkPoolSubmit(b, NewOptimized(10))
}

func BenchmarkOriginal_Workers100(b *testing.B) {
	benchmarkPoolSubmit(b, New(100))
}

func BenchmarkOptimized_Workers100(b *testing.B) {
	benchmarkPoolSubmit(b, NewOptimized(100))
}

func BenchmarkOriginal_Workers1000(b *testing.B) {
	benchmarkPoolSubmit(b, New(1000))
}

func BenchmarkOptimized_Workers1000(b *testing.B) {
	benchmarkPoolSubmit(b, NewOptimized(1000))
}

// 通用基准测试辅助函数
func benchmarkPoolSubmit(b *testing.B, p Pool) {
	defer p.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func() {})
	}
}

// 实际工作负载测试
func BenchmarkOriginal_WithWork(b *testing.B) {
	p := New(100)
	defer p.Stop()
	
	b.ResetTimer()
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		p.Submit(func() {
			// 模拟一些实际工作
			sum := 0
			for j := 0; j < 100; j++ {
				sum += j
			}
			wg.Done()
		})
	}
	wg.Wait()
}

func BenchmarkOptimized_WithWork(b *testing.B) {
	p := NewOptimized(100)
	defer p.Stop()
	
	b.ResetTimer()
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		p.Submit(func() {
			// 模拟一些实际工作
			sum := 0
			for j := 0; j < 100; j++ {
				sum += j
			}
			wg.Done()
		})
	}
	wg.Wait()
}

// 高并发场景测试
func BenchmarkOriginal_HighConcurrency(b *testing.B) {
	p := New(100)
	defer p.Stop()
	
	b.ResetTimer()
	b.SetParallelism(16) // 16个goroutine并发提交
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg := sync.WaitGroup{}
			wg.Add(1)
			p.Submit(func() {
				sum := 0
				for j := 0; j < 50; j++ {
					sum += j
				}
				wg.Done()
			})
			wg.Wait()
		}
	})
}

func BenchmarkOptimized_HighConcurrency(b *testing.B) {
	p := NewOptimized(100)
	defer p.Stop()
	
	b.ResetTimer()
	b.SetParallelism(16)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg := sync.WaitGroup{}
			wg.Add(1)
			p.Submit(func() {
				sum := 0
				for j := 0; j < 50; j++ {
					sum += j
				}
				wg.Done()
			})
			wg.Wait()
		}
	})
}
