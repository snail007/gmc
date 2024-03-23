package grate

import (
	gatomic "github.com/snail007/gmc/util/sync/atomic"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	t.Parallel()
	rate := 20
	interval := time.Minute
	requests := 30

	limiter := NewRateLimiter(rate, interval)
	assert.Equal(t, rate, limiter.Rate())
	assert.Equal(t, interval, limiter.Interval())

	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()

			// 模拟请求处理时间
			time.Sleep(time.Millisecond * 100)

			if limiter.Allow() {
				t.Logf("Request %d allowed\n", requestNum)
			} else {
				t.Logf("Request %d denied\n", requestNum)
			}
		}(i + 1)
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	t.Logf("Total elapsed time: %s\n", elapsedTime)

	if elapsedTime.Seconds() > 10 {
		t.Errorf("Total elapsed time is more than 10 seconds")
	}
}

func TestRateLimiterAllow(t *testing.T) {
	t.Parallel()
	rate := 5
	interval := time.Second
	limiter := NewRateLimiter(rate, interval)

	// 设置等待组
	var wg sync.WaitGroup

	// 记录允许的请求数
	allowedCount := gatomic.NewInt32(0)

	// 模拟并发请求
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				allowedCount.Increase(1)
			}
		}()
	}

	// 等待所有请求完成
	wg.Wait()

	// 预期允许的请求数应小于等于限制的速率
	expectedAllowed := int32(rate)
	if allowedCount.Val() > expectedAllowed {
		t.Errorf("Allowed requests exceeded expected rate")
	}
}

func TestRateLimiterWait(t *testing.T) {
	t.Parallel()
	rate := 5
	interval := time.Second
	requests := 10

	limiter := NewRateLimiter(rate, interval)

	var wg sync.WaitGroup

	startTime := time.Now()

	var totalWaitTime time.Duration

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()

			// 模拟请求处理时间
			time.Sleep(time.Millisecond * 100)

			waitStart := time.Now()
			limiter.Wait()
			waitTime := time.Since(waitStart)
			t.Logf("Request %d allowed after waiting %s\n", requestNum, waitTime)
			totalWaitTime += waitTime
		}(i + 1)
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	t.Logf("Total elapsed time: %s\n", elapsedTime)

	// 最大可能的等待时间
	maxPossibleWaitTime := 16 * time.Second
	if totalWaitTime > maxPossibleWaitTime {
		t.Errorf("Total elapsed time is more than expected")
	}
}

func TestRateLimiter_WaitN(t *testing.T) {
	t.Parallel()
	rl := NewRateLimiter(5, time.Second)

	// Create a ticker that simulates time passing
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		time.Sleep(2 * time.Second)
		rl.AllowN(3)
	}()

	start := time.Now()

	// Simulate the passage of time with the ticker
	for range ticker.C {
		if elapsed := time.Since(start); elapsed >= 2*time.Second {
			break
		}
	}

	// Wait for the remaining tokens
	rl.WaitN(3)
	elapsed := time.Since(start)

	if elapsed < 2*time.Second {
		t.Errorf("WaitN didn't wait long enough")
	}
}

func TestRateLimiter_AllowN(t *testing.T) {
	t.Parallel()
	rl := NewRateLimiter(5, time.Second)

	// Test initial allowance
	if !rl.AllowN(3) {
		t.Errorf("Expected initial allowance for 3 tokens")
	}

	// Test exceeding allowance
	if rl.AllowN(10) {
		t.Errorf("Expected allowance to fail for 10 tokens")
	}

	// Wait for a short period of time to refill tokens
	time.Sleep(2 * time.Second)

	// Test allowance after refill
	if !rl.AllowN(4) {
		t.Errorf("Expected allowance after refill for 4 tokens")
	}
}

func TestRateLimiter_AllowN_Refill(t *testing.T) {
	t.Parallel()
	rl := NewRateLimiter(5, time.Second)

	// Consume initial tokens
	rl.AllowN(3)

	// Wait for a short period of time to refill tokens
	time.Sleep(2 * time.Second)

	// Test allowance after refill
	if !rl.AllowN(4) {
		t.Errorf("Expected allowance after refill for 4 tokens")
	}
}
