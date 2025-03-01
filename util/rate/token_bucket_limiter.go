// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package grate

import (
	"golang.org/x/time/rate"
	"time"
)

type TokenBucketLimiter struct {
	*rate.Limiter
	count    int
	duration time.Duration
}

func NewTokenBucketLimiter(count int, duration time.Duration) *TokenBucketLimiter {
	return NewTokenBucketBurstLimiter(count, duration, count)
}

func NewTokenBucketBurstLimiter(count int, duration time.Duration, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		Limiter:  rate.NewLimiter(rate.Every(duration/time.Duration(count)), burst),
		count:    count,
		duration: duration,
	}
}
