package grate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 单线程测试
func TestSingleThread(t *testing.T) {
	limiter := NewSlidingWindowLimiter(5, time.Second)
	assert.Equal(t, limiter.Capacity(), 5)
	assert.Equal(t, limiter.Duration(), time.Second)
	var success, fail int32

	for i := 0; i < 10; i++ {
		if limiter.Allow() {
			success++
		} else {
			fail++
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("[Single Thread] Allowed:", success, "Denied:", fail)
	assert.Equal(t, int32(5), success)
	assert.Equal(t, int32(5), fail)
	time.Sleep(time.Second)
	assert.True(t, limiter.AllowN(5))
	assert.False(t, limiter.AllowN(1))
}

// 高并发测试
func TestHighConcurrency(t *testing.T) {
	limiter := NewSlidingWindowLimiter(5, time.Second)
	var success, fail int32
	var wg sync.WaitGroup

	workerCount := 20
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				atomic.AddInt32(&success, 1)
			} else {
				atomic.AddInt32(&fail, 1)
			}
		}()
	}
	wg.Wait()

	fmt.Println("[High Concurrency] Allowed:", success, "Denied:", fail)
	assert.Equal(t, int32(5), success)
	assert.Equal(t, int32(15), fail)

}

// 测试恢复速率
func TestRecoveryRate(t *testing.T) {
	limiter := NewSlidingWindowLimiter(5, time.Second)
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Fatalf("Expected request %d to be allowed", i)
		}
	}
	if limiter.Allow() {
		t.Fatalf("Expected request to be denied after reaching the limit")
	}
	time.Sleep(time.Second)
	if !limiter.Allow() {
		t.Fatalf("Expected request to be allowed after window reset")
	}
}

// 并发+恢复测试
func TestConcurrencyWithRecovery(t *testing.T) {
	limiter := NewSlidingWindowLimiter(5, time.Second)
	var success, fail int32
	var wg sync.WaitGroup

	workerCount := 10
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				atomic.AddInt32(&success, 1)
			} else {
				atomic.AddInt32(&fail, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("[Concurrency with Recovery] Allowed:", success, "Denied:", fail)
	assert.Equal(t, fail, success)
	time.Sleep(time.Second)
	if !limiter.Allow() {
		t.Fatalf("Expected request to be allowed after window reset")
	}
}
