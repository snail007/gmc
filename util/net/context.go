// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"time"

	gmap "github.com/snail007/gmc/util/map"
)

var (
	errHijackedFail = fmt.Errorf("hijack error, unsupported argument")
)

type defaultContext struct {
	conn          *Conn
	listener      *Listener
	eventConn     *EventConn
	eventListener *EventListener
	connBinder    *ConnBinder
	data          *gmap.Map
	hijacked      bool
}

func (s *defaultContext) Listener() *Listener {
	return s.listener
}

func (s *defaultContext) SetListener(listener *Listener) Context {
	s.listener = listener
	return s
}

func (s *defaultContext) Clone() Context {
	return &defaultContext{
		conn:          s.conn,
		eventConn:     s.eventConn,
		eventListener: s.eventListener,
		connBinder:    s.connBinder,
		data:          s.data,
		hijacked:      s.hijacked,
	}
}

func (s *defaultContext) ConnBinder() *ConnBinder {
	return s.connBinder
}

func (s *defaultContext) SetConnBinder(connBinder *ConnBinder) Context {
	s.connBinder = connBinder
	return s
}

func (s *defaultContext) EventConn() *EventConn {
	return s.eventConn
}

func (s *defaultContext) SetEventConn(eventConn *EventConn) Context {
	s.eventConn = eventConn
	return s
}

func (s *defaultContext) EventListener() *EventListener {
	return s.eventListener
}

func (s *defaultContext) SetEventListener(eventListener *EventListener) Context {
	s.eventListener = eventListener
	return s
}

func (s *defaultContext) Hijacked() bool {
	return s.hijacked
}

func (s *defaultContext) ReadBytes() int64 {
	return s.conn.ReadBytes()
}

func (s *defaultContext) Hijack(c ...Codec) (net.Conn, error) {
	if s.hijacked {
		return nil, nil
	}
	if len(c) == 1 {
		s.conn.Conn = c[0]
	} else if len(c) >= 2 {
		return nil, errHijackedFail
	}
	s.hijacked = true
	return nil, nil
}

func (s *defaultContext) WriteBytes() int64 {
	return s.conn.WriteBytes()
}

func (s *defaultContext) Data(key interface{}) interface{} {
	v, _ := s.data.Load(key)
	return v
}

func (s *defaultContext) SetData(key, value interface{}) Context {
	s.data.Store(key, value)
	return s
}

func (s *defaultContext) SetConn(c net.Conn) Context {
	s.conn.Conn = c
	return s
}

func (s *defaultContext) Conn() net.Conn {
	return s.conn.Conn
}

func (s *defaultContext) RawConn() net.Conn {
	return s.conn.RawConn()
}

func (s *defaultContext) ReadTimeout() time.Duration {
	return s.conn.ReadTimeout()
}

func (s *defaultContext) WriteTimeout() time.Duration {
	return s.conn.WriteTimeout()
}

func (s *defaultContext) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *defaultContext) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func NewContext() Context {
	return &defaultContext{
		data: gmap.New(),
	}
}
