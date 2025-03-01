package grate

import (
	"sync/atomic"
	"time"
)

type SlidingWindowLimiter struct {
	buckets   []int32
	index     int32
	count     int32
	capacity  int32
	duration  time.Duration
	interval  time.Duration
	lastCheck int64
}

// NewSlidingWindowLimiter 创建一个新的滑动窗口限流器
func NewSlidingWindowLimiter(count int, duration time.Duration) *SlidingWindowLimiter {
	bucketSize := 10 // 分成10个时间窗口
	interval := duration / time.Duration(bucketSize)
	return &SlidingWindowLimiter{
		buckets:   make([]int32, bucketSize),
		capacity:  int32(count),
		duration:  duration,
		interval:  interval,
		lastCheck: time.Now().UnixNano(),
	}
}

// AllowN 允许 n 个请求
func (l *SlidingWindowLimiter) AllowN(n int32) bool {
	now := time.Now().UnixNano()
	last := atomic.LoadInt64(&l.lastCheck)
	elapsed := now - last

	bucketsToShift := int(elapsed / l.interval.Nanoseconds())
	if bucketsToShift > 0 {
		if atomic.CompareAndSwapInt64(&l.lastCheck, last, now) {
			if bucketsToShift > len(l.buckets) {
				bucketsToShift = len(l.buckets)
			}
			for i := 0; i < bucketsToShift; i++ {
				oldIndex := (int(atomic.LoadInt32(&l.index)) + i + 1) % len(l.buckets)
				oldVal := atomic.SwapInt32(&l.buckets[oldIndex], 0)
				atomic.AddInt32(&l.count, -oldVal)
			}

			oldIndex := atomic.LoadInt32(&l.index)
			newIndex := (int(oldIndex) + bucketsToShift) % len(l.buckets)
			atomic.StoreInt32(&l.index, int32(newIndex))
		}
	}

	// 防止计数变成负数
	currentCount := atomic.LoadInt32(&l.count)
	if currentCount < 0 {
		atomic.StoreInt32(&l.count, 0)
	}

	// 检查是否超过限流
	if atomic.LoadInt32(&l.count)+n > l.capacity {
		return false
	}

	// 增加请求计数
	atomic.AddInt32(&l.buckets[atomic.LoadInt32(&l.index)], n)
	atomic.AddInt32(&l.count, n)
	return true
}

// Allow 允许 1 个请求
func (l *SlidingWindowLimiter) Allow() bool {
	return l.AllowN(1)
}
