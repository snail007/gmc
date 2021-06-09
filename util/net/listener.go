// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"net"
	"sync"
	"time"
)

type EventListenerErrorHandler func(l *EventListener, err error)

type EventListenerFirstReadTimeoutHandler func(l *EventListener, c net.Conn, err error)

type EventListenerOnCloseHandler func(l *EventListener)

type EventListenerOnAcceptHandler func(l *EventListener, c net.Conn)

type CodecFactory func() Codec

type EventListener struct {
	l                      net.Listener
	codecFactory           []CodecFactory
	onClose                EventListenerOnCloseHandler
	onAcceptError          EventListenerErrorHandler
	onAccept               EventListenerOnAcceptHandler
	onCodecInitializeError EventListenerErrorHandler
	onFistReadTimeoutError EventListenerFirstReadTimeoutHandler
	firstReadTimeout       time.Duration
	closeOnce              *sync.Once
}

func (s *EventListener) OnFistReadTimeout(h EventListenerFirstReadTimeoutHandler) *EventListener {
	s.onFistReadTimeoutError = h
	return s
}

func (s *EventListener) AddCodecFactory(codecFactory CodecFactory) *EventListener {
	s.codecFactory = append(s.codecFactory, codecFactory)
	return s
}

func (s *EventListener) OnAcceptError(h EventListenerErrorHandler) *EventListener {
	s.onAcceptError = h
	return s
}

func (s *EventListener) OnAccept(h EventListenerOnAcceptHandler) *EventListener {
	s.onAccept = h
	return s
}

func (s *EventListener) OnCodecInitializeError(h EventListenerErrorHandler) *EventListener {
	s.onCodecInitializeError = h
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
			// timeout check
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
			// codec check
			if len(s.codecFactory) > 0 {
				conn := NewConn(c)
				for _, cf := range s.codecFactory {
					conn.AddCodec(cf())
				}
				err = conn.Initialize()
				if err != nil {
					c.Close()
					s.onCodecInitializeError(s, err)
					continue
				}
				c=conn
			}
			// accept
			s.onAccept(s, c)
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
		closeOnce:              &sync.Once{},
		l:                      l,
		onAcceptError:          func(*EventListener, error) {},
		onCodecInitializeError: func(*EventListener, error) {},
		onClose:                func(*EventListener) {},
		onFistReadTimeoutError: func(*EventListener, net.Conn, error) {},
	}
}
