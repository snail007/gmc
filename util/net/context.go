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

type Context interface {
	Data(key interface{}) interface{}
	SetData(key, value interface{})
	Hijack(c ...Codec) (net.Conn, error)
}

type ConnContext interface {
	Context
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	RemoteAddr() net.Addr
	LocalAddr() net.Addr
	Conn() net.Conn
	RawConn() net.Conn
	ReadBytes() int64
	WriteBytes() int64
}

type CodecContext interface {
	ConnContext
	SetConn(net.Conn) CodecContext
}

type defaultContext struct {
	conn *Conn
	data *gmap.Map
}

func (s *defaultContext) ReadBytes() int64 {
	return s.conn.ReadBytes()
}

func (s *defaultContext) Hijack(c ...Codec) (net.Conn, error) {
	if len(c) == 1 {
		s.conn.Conn = c[0]
	} else if len(c) >= 2 {
		return nil, fmt.Errorf("hijack error, unsupported argument")
	}
	return nil, ErrHijacked
}

func (s *defaultContext) WriteBytes() int64 {
	return s.conn.WriteBytes()
}

func (s *defaultContext) Data(key interface{}) interface{} {
	v, _ := s.data.Load(key)
	return v
}

func (s *defaultContext) SetData(key, value interface{}) {
	s.data.Store(key, value)
}

func (s *defaultContext) SetConn(c net.Conn) CodecContext {
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
