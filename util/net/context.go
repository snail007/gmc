// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"net"
	"time"

	gmap "github.com/snail007/gmc/util/map"
)

type defaultContext struct {
	conn          *Conn
	listener      *Listener
	eventConn     *EventConn
	eventListener *EventListener
	connBinder    *ConnBinder
	data          *gmap.Map
	isHijacked    bool
	isBreak       bool
	isContinue    bool
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
		data:          s.data.Clone(),
		isHijacked:    s.isHijacked,
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

func (s *defaultContext) ReadBytes() int64 {
	return s.conn.ReadBytes()
}

func (s *defaultContext) IsHijacked() bool {
	return s.isHijacked
}

func (s *defaultContext) Hijack() {
	s.isHijacked = true
	s.isBreak = true
	return
}

func (s *defaultContext) IsBreak() bool {
	return s.isBreak
}

func (s *defaultContext) Break() {
	s.isBreak = true
}

func (s *defaultContext) Continue() {
	s.isContinue = true
}

func (s *defaultContext) IsContinue() bool {
	return s.isContinue
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

func (s *defaultContext) IsTLS() bool {
	v := s.Data(isTLSKey)
	if v == nil {
		return false
	}
	if vv, ok := v.(bool); ok {
		return vv
	}
	return false
}

func NewContext() Context {
	return &defaultContext{
		data: gmap.New(),
	}
}
