// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"time"

	gbytes "github.com/snail007/gmc/util/bytes"
)

const (
	defaultBufferedConnSize         = 8192
	defaultEventConnReadBufferSize  = 8192
	defaultConnBinderReadBufferSize = 8192
	defaultConnKeepAlivePeriod      = time.Second * 15
)

var (
	ErrFirstReadTimeout = fmt.Errorf("ErrFirstReadTimeout")
	ErrConnInitialize   = fmt.Errorf("ErrConnInitialize")
)

type (
	Context interface {
		Clone() Context
		EventConn() *EventConn
		SetEventConn(eventConn *EventConn) Context
		EventListener() *EventListener
		SetEventListener(eventListener *EventListener) Context
		ConnBinder() *ConnBinder
		SetConnBinder(connBinder *ConnBinder) Context
		Hijacked() bool
		ReadBytes() int64
		Hijack(c ...Codec) (net.Conn, error)
		WriteBytes() int64
		Data(key interface{}) interface{}
		SetData(key, value interface{}) Context
		SetConn(c net.Conn) Context
		Conn() net.Conn
		RawConn() net.Conn
		Listener() *Listener
		SetListener(listener *Listener) Context
		ReadTimeout() time.Duration
		WriteTimeout() time.Duration
		RemoteAddr() net.Addr
		LocalAddr() net.Addr
	}
	Codec interface {
		net.Conn
		Initialize(ctx Context, next NextCodec) (net.Conn, error)
	}
	NextCodec interface {
		Call(ctx Context) (net.Conn, error)
	}
	NextConnFilter interface {
		Call(ctx Context, c net.Conn) (net.Conn, error)
	}
	ErrorHandler            func(ctx Context, err error)
	ConnFilter              func(ctx Context, c net.Conn, next NextConnFilter) (net.Conn, error)
	CloseHandler            func(ctx Context)
	DataHandler             func(ctx Context, data []byte)
	AcceptHandler           func(ctx Context, c net.Conn)
	CodecFactory            func(ctx Context) Codec
	FirstReadTimeoutHandler func(ctx Context, c net.Conn, err error)
)

func Write(addr, data string) (err error) {
	return WriteBytes(addr, []byte(data))
}

func WriteBytes(addr string, data []byte) (err error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer func() {
		time.AfterFunc(time.Second, func() {
			c.Close()
		})
	}()
	_, err = c.Write(data)
	return
}

func Read(c net.Conn, bufSize int) (d string, err error) {
	data, err := ReadBytes(c, bufSize)
	d = string(data)
	return
}

func ReadBytes(c net.Conn, bufSize int) (d []byte, err error) {
	buf := gbytes.GetPool(bufSize).Get().([]byte)
	defer gbytes.GetPool(bufSize).Put(buf)
	n, err := c.Read(buf)
	if n > 0 {
		d = append(d, buf[:n]...)
	}
	return
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
