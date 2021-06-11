// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type ListenerErrorHandler func(l *EventListener, err error)

type FirstReadTimeoutHandler func(l *EventListener, c net.Conn, err error)

type OnCloseHandler func(l *EventListener)

type AcceptHandler func(l *EventListener, ctx Context, c net.Conn)

type CodecFactory func() Codec

type ContextFactory func(c net.Conn) Context

type EventListener struct {
	l                      net.Listener
	contextFactory         ContextFactory
	codecFactory           []CodecFactory
	connFilters            []ConnFilter
	onClose                OnCloseHandler
	onAcceptError          ListenerErrorHandler
	onAccept               AcceptHandler
	onFistReadTimeoutError FirstReadTimeoutHandler
	firstReadTimeout       time.Duration
	closeOnce              *sync.Once
}

func (s *EventListener) AddConnFilter(f ConnFilter) *EventListener {
	s.connFilters = append(s.connFilters, f)
	return s
}

func (s *EventListener) SetConnContextFactory(contextFactory ContextFactory) {
	s.contextFactory = contextFactory
}

func (s *EventListener) OnFistReadTimeout(h FirstReadTimeoutHandler) *EventListener {
	s.onFistReadTimeoutError = h
	return s
}

func (s *EventListener) AddCodecFactory(codecFactory CodecFactory) *EventListener {
	s.codecFactory = append(s.codecFactory, codecFactory)
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
	s.firstReadTimeout = firstReadTimeout
	return s
}

func (s *EventListener) Start() *EventListener {
	go func() {
		defer func() {
			s.Close()
		}()
		for {
			c, err := s.l.Accept()
			if err != nil {
				if v, ok := err.(*net.OpError); ok && !v.Timeout() && v.Temporary() {
					continue
				} else {
					s.onAcceptError(s, err)
					return
				}
			}
			ctx := s.contextFactory(c)

			// first read timeout check
			if s.firstReadTimeout > 0 {
				bc := NewBufferedConn(c)
				bc.SetReadDeadline(time.Now().Add(s.firstReadTimeout))
				_, err = bc.ReadByte()
				if err != nil {
					c.Close()
					s.onFistReadTimeoutError(s, c, err)
					continue
				}
				bc.UnreadByte()
				c = bc
			}
			// init Conn
			if len(s.codecFactory) > 0 || len(s.connFilters) > 0 {
				conn := NewContextConn(ctx.(*defaultContext), c)
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
					c.Close()
					s.onAcceptError(s, err)
					continue
				}
				c = conn
			}
			// accept
			s.onAccept(s, ctx, c)
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

func NewEventListener(l net.Listener) *EventListener {
	return &EventListener{
		contextFactory: func(net.Conn) Context {
			return NewContext()
		},
		closeOnce: &sync.Once{},
		l:         l,
		onAccept: func(l *EventListener, ctx Context, c net.Conn) {
			c.Close()
			fmt.Println("[WARN] you should using OnAccept() to set a accept handler to process the connection")
		},
		onAcceptError:          func(*EventListener, error) {},
		onClose:                func(*EventListener) {},
		onFistReadTimeoutError: func(*EventListener, net.Conn, error) {},
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
