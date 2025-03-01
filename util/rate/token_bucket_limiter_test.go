// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package grate

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

// TestTokenBucketLimiter 测试在单位时间内允许固定数量的请求
func TestTokenBucketLimiter(t *testing.T) {
	// 创建一个每秒允许 5 个请求的限流器
	limiter := NewTokenBucketLimiter(5, time.Second)

	// 测试连续调用 Allow() 是否允许前5个请求
	for i := 0; i < 5; i++ {
		assert.True(t, limiter.Allow(), "第%d个请求应被允许", i+1)
	}

	// 第6个请求应被拒绝
	assert.False(t, limiter.Allow(), "不应允许超过5个请求")
}

// TestTokenBucketLimiter_WithBurst 测试带突发容量的限流器
func TestTokenBucketLimiter_WithBurst(t *testing.T) {
	// 创建一个每秒允许 5 个请求，但突发容量设置为 10 的限流器
	limiter := NewTokenBucketBurstLimiter(5, time.Second, 10)

	// 测试允许10个请求，因为有突发能力
	for i := 0; i < 10; i++ {
		assert.True(t, limiter.Allow(), "第%d个请求应被允许", i+1)
	}

	// 第11个请求应被拒绝
	assert.False(t, limiter.Allow(), "不应允许超过突发容量的请求")
}

// TestTokenBucketLimiter_Concurrency 高并发情况下的测试
func TestTokenBucketLimiter_Concurrency(t *testing.T) {
	limiter := NewTokenBucketLimiter(5, time.Second)
	const numRequests = 100

	var allowed int32 = 0
	var denied int32 = 0

	// 使用 WaitGroup 确保所有 goroutine 执行完成
	done := make(chan struct{})

	for i := 0; i < numRequests; i++ {
		go func() {
			if limiter.Allow() {
				atomic.AddInt32(&allowed, 1)
			} else {
				atomic.AddInt32(&denied, 1)
			}
			done <- struct{}{}
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < numRequests; i++ {
		<-done
	}

	// 由于限流器每秒只允许 5 个请求，故允许的请求应为5，其余均被拒绝
	assert.Equal(t, int32(5), allowed, "预期每秒允许5个请求")
	assert.Equal(t, int32(numRequests-5), denied, "预期拒绝剩余请求")
}
