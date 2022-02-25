// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"context"
	"net"
	"sync"

	"golang.org/x/time/rate"
)

var (
	burstCnt = 1024 * 1024 * 1024
)

type RateCodec struct {
	net.Conn
	readLimiter  *rate.Limiter
	writeLimiter *rate.Limiter
	readCtx      context.Context
	writeCtx     context.Context
	readCancel   context.CancelFunc
	writCancel   context.CancelFunc
	closeOnce    *sync.Once
	burst        int
}

func NewRateCodec(bytesPerSecond int64) *RateCodec {
	r := &RateCodec{
		readLimiter:  rate.NewLimiter(rate.Limit(bytesPerSecond), burstCnt),
		writeLimiter: rate.NewLimiter(rate.Limit(bytesPerSecond), burstCnt),
		closeOnce:    &sync.Once{},
		burst:        burstCnt,
	}
	r.readCtx, r.readCancel = context.WithCancel(context.Background())
	r.writeCtx, r.writCancel = context.WithCancel(context.Background())
	// clear the bucket
	r.readLimiter.WaitN(r.readCtx, burstCnt)
	r.writeLimiter.WaitN(r.writeCtx, burstCnt)
	return r
}

func (s *RateCodec) Read(p []byte) (int, error) {
	n, err := s.Conn.Read(p)
	if err != nil {
		return n, err
	}
	s.waitN(true, n)
	return n, nil
}

func (s *RateCodec) Write(p []byte) (int, error) {
	n, err := s.Conn.Write(p)
	if err != nil {
		return n, err
	}
	s.waitN(false, n)
	return n, err
}

func (s *RateCodec) Close() (err error) {
	s.closeOnce.Do(func() {
		go s.readCancel()
		s.writCancel()
		if s.Conn!=nil{
			err = s.Conn.Close()
		}
	})
	return
}

func (s *RateCodec) waitN(isRead bool, n int) (err error) {
	if n == 0 {
		return
	}
	var limiter *rate.Limiter
	var ctx context.Context
	if isRead {
		ctx = s.readCtx
		limiter = s.readLimiter
	} else {
		ctx = s.writeCtx
		limiter = s.writeLimiter
	}
	if n <= s.burst {
		return limiter.WaitN(ctx, n)
	}
	for i := 0; i < n/s.burst; i++ {
		err = limiter.WaitN(ctx, s.burst)
		if err != nil {
			return
		}
	}
	if v := n % s.burst; v > 0 {
		err = limiter.WaitN(ctx, v)
	}
	return
}
func (s *RateCodec) Initialize(ctx Context, next NextCodec) (net.Conn, error) {
	s.Conn = ctx.Conn()
	return next.Call(ctx.SetConn(s))
}
