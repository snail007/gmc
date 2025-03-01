// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package grate

import (
	"testing"
	"time"
)

// BenchmarkTokenBucketLimiter 在高并发下测量限流器的性能
func BenchmarkTokenBucketLimiter(b *testing.B) {
	limiter := NewTokenBucketLimiter(5, time.Second)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow()
		}
	})
}

// BenchmarkTokenBucketLimiter_Burst 带突发能力的限流器性能测试
func BenchmarkTokenBucketLimiter_Burst(b *testing.B) {
	limiter := NewTokenBucketBurstLimiter(5, time.Second, 10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow()
		}
	})
}
