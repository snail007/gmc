// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	glog "github.com/snail007/gmc/module/log"
)

var _ net.Listener = &Listener{}
var _ net.Listener = &ProtocolListener{}
var (
	defaultConnQueueSize = 2048
	ErrClosedListener    = errors.New("listener is closed")
)

type ProtocolChecker func(listener *Listener, conn BufferedConn) bool

type ProtocolListenerOption struct {
	Name              string
	Checker           ProtocolChecker
	ConnQueueSize     int
	OnQueueOverflow   func(l net.Listener, opt *ProtocolListenerOption, conn BufferedConn)
	OverflowAutoClose bool
}

// ProtocolListener 在构建多协议服务时，必须使用 `Listener.SetFirstReadTimeout` 设置一个合理的时间，来防止恶意或行为不当的客户端阻塞整个服务。
type ProtocolListener struct {
	opt       *ProtocolListenerOption
	connChn   chan net.Conn
	closed    bool
	closeOnce sync.Once
	l         *Listener
}

func (s *ProtocolListener) Accept() (net.Conn, error) {
	conn, ok := <-s.connChn
	if !ok {
		return nil, ErrClosedListener
	}
	return conn, nil
}

func (s *ProtocolListener) Close() error {
	s.closeOnce.Do(func() {
		s.closed = true
		close(s.connChn)
	})
	return nil
}

func (s *ProtocolListener) Addr() net.Addr {
	return s.l.listener.Addr()
}

type Listener struct {
	listener                      net.Listener
	protocolListeners             []*ProtocolListener
	codecFactory                  []CodecFactory
	connFilters                   []ConnFilter
	onConnClose                   CloseHandler
	filters                       []ConnFilter
	firstReadTimeout              time.Duration
	onFistReadTimeoutError        FirstReadTimeoutHandler
	beforeFirstRead               BeforeFirstReadHandler
	connCount                     *int64
	autoCloseConnOnReadWriteError bool
	ctx                           Context
}

func (s *Listener) NewProtocolListener(opt *ProtocolListenerOption) *Listener {
	if opt.ConnQueueSize <= 0 {
		opt.ConnQueueSize = defaultConnQueueSize
	}
	l := &ProtocolListener{
		connChn: make(chan net.Conn, opt.ConnQueueSize),
		l:       s,
		opt:     opt,
	}
	s.protocolListeners = append(s.protocolListeners, l)
	return NewListener(l)
}

func (s *Listener) Close() error {
	for _, v := range s.protocolListeners {
		if !v.closed {
			v.Close()
		}
	}
	return s.listener.Close()
}

func (s *Listener) Addr() net.Addr {
	return s.listener.Addr()
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

func (s *Listener) SetBeforeFirstRead(h BeforeFirstReadHandler) *Listener {
	s.beforeFirstRead = h
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
		c, err = s.listener.Accept()
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
	if ctx.IsHijacked() {
		// isHijacked by filter, just call next accept.
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
		if s.beforeFirstRead != nil {
			err = s.beforeFirstRead(ctx, c)
			if err != nil {
				c.Close()
				goto retry
			}
		}
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

	//protocol check
	if len(s.protocolListeners) > 0 {
		bc := NewBufferedConn(c)
		for _, v := range s.protocolListeners {
			if !v.closed && v.opt.Checker(s, bc) {
				select {
				case v.connChn <- bc:
				default:
					if v.opt.OnQueueOverflow != nil {
						v.opt.OnQueueOverflow(v, v.opt, bc)
					} else {
						bc.Close()
						glog.Warnf("protocol listener's conn queue is overflow, name: %s, size: %d", v.opt.Name, v.opt.ConnQueueSize)
					}
					if v.opt.OverflowAutoClose {
						bc.Close()
					}
				}
				goto retry
			}
		}
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
func NewListenerAddr(addr string) (*Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewListener(l), nil
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
			listener:    l,
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
func NewEventListenerAddr(addr string) (*EventListener, error) {
	l, err := NewListenerAddr(addr)
	if err != nil {
		return nil, err
	}
	return NewEventListener(l), nil
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
