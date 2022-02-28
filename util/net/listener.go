// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Listener struct {
	net.Listener
	codecFactory                  []CodecFactory
	connFilters                   []ConnFilter
	onConnClose                   CloseHandler
	filters                       []ConnFilter
	firstReadTimeout              time.Duration
	onFistReadTimeoutError        FirstReadTimeoutHandler
	connCount                     *int64
	autoCloseConnOnReadWriteError bool
	ctx                           Context
}

func (s *Listener) Ctx() Context {
	return s.ctx
}

func (s *Listener) SetAutoCloseConnOnReadWriteError(b bool) {
	s.autoCloseConnOnReadWriteError = b
}

func (s *Listener) SetOnConnClose(onConnClose CloseHandler) *Listener {
	s.onConnClose = onConnClose
	return s
}

func (s *Listener) AddListenerFilter(f ConnFilter) *Listener {
	s.filters = append(s.filters, f)
	return s
}

func (s *Listener) SetFirstReadTimeout(firstReadTimeout time.Duration) *Listener {
	s.firstReadTimeout = firstReadTimeout
	return s
}

func (s *Listener) AddConnFilter(f ConnFilter) *Listener {
	s.connFilters = append(s.connFilters, f)
	return s
}

func (s *Listener) AddCodecFactory(codecFactory CodecFactory) *Listener {
	s.codecFactory = append(s.codecFactory, codecFactory)
	return s
}

func (s *Listener) onConnCloseHook(c *Conn) {
	atomic.AddInt64(s.connCount, -1)
	s.onConnClose(c.ctx)
}

func (s *Listener) onAcceptHook(c net.Conn) {
	atomic.AddInt64(s.connCount, 1)
}

func (s *Listener) ConnCount() int64 {
	return atomic.LoadInt64(s.connCount)
}

func (s *Listener) OnFistReadTimeout(h FirstReadTimeoutHandler) {
	s.onFistReadTimeoutError = h
}

func (s *Listener) Accept() (c net.Conn, err error) {
	c, err = s.accept()
	return
}

func (s *Listener) accept() (c net.Conn, err error) {
	var tempDelay time.Duration // how long to sleep on accept failure
retry:
	for {
		c, err = s.Listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			} else {
				return nil, err
			}
		}
		break
	}
	tempDelay = 0
	// save raw conn
	rawConn := c

	// root conn context
	ctx := s.ctx.Clone()

	// init listener filters and checking hijack
	fConn, err := newConnFilters(s.filters).Call(ctx, c)
	if ctx.Hijacked() {
		// hijacked by filter, just call next accept.
		goto retry
	}
	if err != nil {
		return nil, err
	}
	c = fConn

	// hook to count the connection
	s.onAcceptHook(c)

	// first read timeout check
	if s.firstReadTimeout > 0 {
		bc := NewBufferedConn(c)
		bc.SetReadDeadline(time.Now().Add(s.firstReadTimeout))
		_, err = bc.ReadByte()
		bc.SetReadDeadline(time.Time{})
		if err != nil {
			bc.Close()
			if s.onFistReadTimeoutError != nil {
				s.onFistReadTimeoutError(ctx, c, err)
			}
			goto retry
		}
		bc.UnreadByte()
		c = bc
	}

	// init Conn
	conn := NewContextConn(ctx, c)
	conn.rawConn = rawConn
	// add filters
	for _, f := range s.connFilters {
		conn.AddFilter(f)
	}
	// add codec
	for _, cf := range s.codecFactory {
		conn.AddCodec(cf(ctx))
	}

	// Conn closing hook
	conn.onClose = s.onConnCloseHook

	// sets Conn auto close?
	conn.SetAutoCloseOnReadWriteError(s.autoCloseConnOnReadWriteError)

	c = conn

	return
}

func NewListener(l net.Listener) *Listener {
	return NewContextListener(NewContext(), l)
}

func NewContextListener(ctx Context, l net.Listener) *Listener {
	var lis *Listener
	if v, ok := l.(*Listener); ok {
		lis = v
	} else {
		lis = &Listener{
			Listener:    l,
			onConnClose: func(Context) {},
			connCount:   new(int64),
		}
	}
	ctx.SetListener(lis)
	lis.ctx = ctx
	return lis
}

type EventListener struct {
	l             *Listener
	onClose       CloseHandler
	onAcceptError ErrorHandler
	onAccept      AcceptHandler
	closeOnce     *sync.Once
	addr          *Addr
	ctx           Context
	autoCloseConn bool
	started       bool
}

// SetAutoCloseConn if true, EventListener will close Conn after accept handler return.
func (s *EventListener) SetAutoCloseConn(autoCloseConn bool) {
	s.autoCloseConn = autoCloseConn
}

func (s *EventListener) Addr() *Addr {
	return s.addr
}

func (s *EventListener) Ctx() Context {
	return s.ctx
}

func (s *EventListener) SetAutoCloseConnOnReadWriteError(b bool) *EventListener {
	s.l.SetAutoCloseConnOnReadWriteError(b)
	return s
}

func (s *EventListener) SetOnConnClose(onConnClose CloseHandler) *EventListener {
	s.l.SetOnConnClose(onConnClose)
	return s
}

func (s *EventListener) AddConnFilter(f ConnFilter) *EventListener {
	s.l.AddConnFilter(f)
	return s
}

func (s *EventListener) AddListenerFilter(f ConnFilter) *EventListener {
	s.l.AddListenerFilter(f)
	return s
}

func (s *EventListener) OnFistReadTimeout(h FirstReadTimeoutHandler) *EventListener {
	s.l.onFistReadTimeoutError = h
	return s
}

func (s *EventListener) AddCodecFactory(codecFactory CodecFactory) *EventListener {
	s.l.AddCodecFactory(codecFactory)
	return s
}

func (s *EventListener) OnAcceptError(h ErrorHandler) *EventListener {
	s.onAcceptError = h
	return s
}

func (s *EventListener) OnAccept(h AcceptHandler) *EventListener {
	s.onAccept = h
	return s
}

func (s *EventListener) SetFirstReadTimeout(firstReadTimeout time.Duration) *EventListener {
	s.l.SetFirstReadTimeout(firstReadTimeout)
	return s
}

func (s *EventListener) StartAndWait() {
	if s.started {
		return
	}
	defer func() {
		s.Close()
	}()
	for {
		c, err := s.l.accept()
		// root conn context
		switch err {
		case nil:
			go func() {
				// accept
				s.onAccept(c.(*Conn).ctx, c)
				if s.autoCloseConn {
					c.Close()
				}
			}()
		default:
			s.onAcceptError(s.ctx, err)
			return
		}
	}
}

func (s *EventListener) Start() *EventListener {
	go s.StartAndWait()
	return s
}

func (s *EventListener) Close() *EventListener {
	if s == nil {
		return nil
	}
	s.closeOnce.Do(func() {
		s.l.Close()
		s.onClose(s.ctx)
	})
	return s
}

func (s *EventListener) ConnCount() int64 {
	return s.l.ConnCount()
}

func NewEventListener(l net.Listener) *EventListener {
	return NewContextEventListener(NewContext(), l)
}

func NewContextEventListener(ctx Context, l net.Listener) *EventListener {
	el := &EventListener{
		closeOnce: &sync.Once{},
		onAccept: func(ctx Context, c net.Conn) {
			c.Close()
			fmt.Println("[WARN] you should using OnAccept() to set a accept handler to process the connection")
		},
		onAcceptError: func(Context, error) {},
		onClose:       func(Context) {},
	}
	ctx.SetEventListener(el)
	el.ctx = ctx
	el.l = NewContextListener(el.ctx, l)
	el.addr = NewAddr(l.Addr())
	return el
}
