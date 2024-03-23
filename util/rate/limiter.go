package grate

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rate       int           // 最大请求数
	interval   time.Duration // 时间窗口
	tokens     int           // 当前令牌数
	lastUpdate time.Time     // 上次更新时间
	mu         sync.Mutex    // 互斥锁
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		interval:   interval,
		tokens:     rate,
		lastUpdate: time.Now(),
	}
}

func (rl *RateLimiter) Rate() int {
	return rl.rate
}

func (rl *RateLimiter) Interval() time.Duration {
	return rl.interval
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate)

	// 根据时间间隔计算应该恢复的令牌数量
	refillTokens := int(elapsed.Seconds() / rl.interval.Seconds())

	if refillTokens > 0 {
		rl.tokens = rl.tokens + refillTokens
		if rl.tokens > rl.rate {
			rl.tokens = rl.rate
		}
		rl.lastUpdate = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) Wait() {
	for !rl.Allow() {
		time.Sleep(rl.interval)
	}
}

func (rl *RateLimiter) AllowN(n int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate)

	// 根据时间间隔计算应该恢复的令牌数量
	refillTokens := int(elapsed.Seconds() / rl.interval.Seconds())

	if refillTokens > 0 {
		rl.tokens = rl.tokens + refillTokens
		if rl.tokens > rl.rate {
			rl.tokens = rl.rate
		}
		rl.lastUpdate = now
	}

	if rl.tokens >= n {
		rl.tokens -= n
		return true
	}

	return false
}

func (rl *RateLimiter) WaitN(n int) {
	for !rl.AllowN(n) {
		time.Sleep(rl.interval)
	}
}
