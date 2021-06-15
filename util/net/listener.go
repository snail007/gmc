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

var (
	ErrFirstReadTimeout       = fmt.Errorf("ErrFirstReadTimeout")
	ErrConnInitialize         = fmt.Errorf("ErrConnInitialize")
	ErrListenerFilterHijacked = fmt.Errorf("ErrListenerFilterHijacked")
)

type (
	ListenerErrorHandler func(l *EventListener, err error)

	OnCloseHandler func(l *EventListener)

	OnConnCloseHandler func(ConnContext)

	AcceptHandler func(l *EventListener, ctx Context, c net.Conn)

	CodecFactory func() Codec

	ContextFactory func(c net.Conn) Context

	FirstReadTimeoutHandler func(c net.Conn, err error)

	ListenerFilter func(ctx Context, c net.Conn, next NextListenerFilter) (net.Conn, error)
)

type NextListenerFilter interface {
	Call(ctx Context, c net.Conn) (net.Conn, error)
}

type listenerFilters struct {
	filters []ListenerFilter
	idx     int
}

func (s *listenerFilters) Call(ctx Context, c net.Conn) (conn net.Conn, err error) {
	s.idx++
	if s.idx == len(s.filters) {
		// last call, no next filter, just return conn.
		return c, nil
	}
	return s.filters[s.idx](ctx, c, s)
}

func newListenerFilters(filters []ListenerFilter) *listenerFilters {
	f := &listenerFilters{
		idx:     -1,
		filters: filters,
	}
	return f
}

type Listener struct {
	net.Listener
	*sync.Once
	contextFactory                ContextFactory
	codecFactory                  []CodecFactory
	connFilters                   []ConnFilter
	onConnClose                   OnConnCloseHandler
	filters                       []ListenerFilter
	firstReadTimeout              time.Duration
	connCount                     *int64
	autoCloseConnOnReadWriteError bool
}

func (s *Listener) SetAutoCloseConnOnReadWriteError(b bool) {
	s.autoCloseConnOnReadWriteError = b
}

func (s *Listener) SetOnConnClose(onConnClose OnConnCloseHandler) *Listener {
	s.onConnClose = onConnClose
	return s
}

func (s *Listener) AddListenerFilter(f ListenerFilter) *Listener {
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

func (s *Listener) SetConnContextFactory(contextFactory ContextFactory) *Listener {
	s.contextFactory = contextFactory
	return s
}

func (s *Listener) Accept() (c net.Conn, err error) {
	c, err, _ = s.accept()
	return
}

func (s *Listener) accept() (c net.Conn, errTyped, err error) {
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
				return nil, err, err
			}
		}
		break
	}
	tempDelay = 0
	// save raw conn
	rawConn := c

	// root conn context
	ctx := s.contextFactory(c)

	// init listener filters
	fConn, err := newListenerFilters(s.filters).Call(ctx, c)
	if err != nil {
		if err == ErrListenerFilterHijacked {
			// hijacked by filter, just call next accept.
			goto retry
		}
		return nil, err, err
	}
	c = fConn

	// hook to count the connection
	s.onAcceptHook(c)

	// first read timeout check
	if s.firstReadTimeout > 0 {
		bc := NewBufferedConn(c)
		bc.SetReadDeadline(time.Now().Add(s.firstReadTimeout))
		_, err = bc.ReadByte()
		if err != nil {
			bc.Close()
			return nil, ErrFirstReadTimeout, err
		}
		bc.UnreadByte()
		c = bc
	}

	// init Conn
	conn := NewContextConn(ctx.(*defaultContext), c)
	conn.rawConn = rawConn
	// add filters
	for _, f := range s.connFilters {
		conn.AddFilter(f)
	}
	// add codec
	for _, cf := range s.codecFactory {
		conn.AddCodec(cf())
	}
	// init
	err = conn.Initialize()
	if err != nil {
		conn.Close()
		return nil, ErrConnInitialize, err
	}
	// Conn closing hook
	conn.onClose = s.onConnCloseHook

	// sets Conn auto close?
	conn.SetAutoCloseOnReadWriteError(s.autoCloseConnOnReadWriteError)

	c = conn

	return
}

func NewListener(l net.Listener) *Listener {
	if v, ok := l.(*Listener); ok {
		return v
	}
	return &Listener{
		Listener: l,
		contextFactory: func(_ net.Conn) Context {
			return NewContext()
		},
		Once:        &sync.Once{},
		onConnClose: func(ConnContext) {},
		connCount:   new(int64),
	}
}

type EventListener struct {
	l                      *Listener
	onClose                OnCloseHandler
	onAcceptError          ListenerErrorHandler
	onAccept               AcceptHandler
	onFistReadTimeoutError FirstReadTimeoutHandler
	closeOnce              *sync.Once
}

func (s *EventListener) SetAutoCloseConnOnReadWriteError(b bool) *EventListener {
	s.l.SetAutoCloseConnOnReadWriteError(b)
	return s
}

func (s *EventListener) SetOnConnClose(onConnClose OnConnCloseHandler) *EventListener {
	s.l.SetOnConnClose(onConnClose)
	return s
}

func (s *EventListener) AddConnFilter(f ConnFilter) *EventListener {
	s.l.AddConnFilter(f)
	return s
}

func (s *EventListener) AddListenerFilter(f ListenerFilter) *EventListener {
	s.l.AddListenerFilter(f)
	return s
}

func (s *EventListener) SetConnContextFactory(contextFactory ContextFactory) *EventListener {
	s.l.SetConnContextFactory(contextFactory)
	return s
}

func (s *EventListener) OnFistReadTimeout(h FirstReadTimeoutHandler) *EventListener {
	s.onFistReadTimeoutError = h
	return s
}

func (s *EventListener) AddCodecFactory(codecFactory CodecFactory) *EventListener {
	s.l.AddCodecFactory(codecFactory)
	return s
}

func (s *EventListener) OnAcceptError(h ListenerErrorHandler) *EventListener {
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

func (s *EventListener) Start() *EventListener {
	go func() {
		defer func() {
			s.Close()
		}()
		for {
			c, errTyped, err := s.l.accept()
			// root conn context
			switch errTyped {
			case ErrFirstReadTimeout:
				s.onFistReadTimeoutError(c, err)
				return
			case nil:
				// accept
				s.onAccept(s, c.(*Conn).ctx, c)
			case ErrConnInitialize:
				fallthrough
			default:
				s.onAcceptError(s, err)
				return
			}
		}
	}()
	return s
}

func (s *EventListener) Close() *EventListener {
	s.closeOnce.Do(func() {
		s.l.Close()
		s.onClose(s)
	})
	return s
}

func (s *EventListener) ConnCount() int64 {
	return s.l.ConnCount()
}

func NewEventListener(l net.Listener) *EventListener {
	return &EventListener{
		closeOnce: &sync.Once{},
		l:         NewListener(l),
		onAccept: func(l *EventListener, ctx Context, c net.Conn) {
			c.Close()
			fmt.Println("[WARN] you should using OnAccept() to set a accept handler to process the connection")
		},
		onAcceptError:          func(*EventListener, error) {},
		onClose:                func(*EventListener) {},
		onFistReadTimeoutError: func(net.Conn, error) {},
	}
}

func ListenRandom(ip ...string) (l net.Listener, port string, err error) {
	ip0 := ""
	if len(ip) == 1 {
		ip0 = ip[0]
	}
	addr := net.JoinHostPort(ip0, "0")
	l, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	_, port, _ = net.SplitHostPort(l.Addr().String())
	return
}
