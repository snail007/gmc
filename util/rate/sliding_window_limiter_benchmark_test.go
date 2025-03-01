// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package grate

import (
	"sync"
	"testing"
	"time"
)

// Benchmark 测试
func BenchmarkLimiter(b *testing.B) {
	limiter := NewSlidingWindowLimiterLimiter(1000, time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Allow()
	}
}

// 高并发基准测试
func BenchmarkHighConcurrency(b *testing.B) {
	limiter := NewSlidingWindowLimiterLimiter(1000, time.Second)
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			limiter.Allow()
			wg.Done()
		}()
	}
	wg.Wait()
}

// AllowN 基准测试
func BenchmarkAllowN(b *testing.B) {
	limiter := NewSlidingWindowLimiterLimiter(1000, time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.AllowN(5)
	}
}

// 高并发 AllowN 基准测试
func BenchmarkAllowNHighConcurrency(b *testing.B) {
	limiter := NewSlidingWindowLimiterLimiter(1000, time.Second)
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			limiter.AllowN(5)
			wg.Done()
		}()
	}
	wg.Wait()
}
