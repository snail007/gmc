// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEventListener_OnFistReadTimeout(t *testing.T) {
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	timeout := false
	el.SetFirstReadTimeout(time.Millisecond * 100).
		OnFistReadTimeout(func(l *EventListener, c net.Conn, err error) {
			timeout = true
		}).Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, timeout)
}

func TestNewEventListener_OnAcceptError(t *testing.T) {
	l, _ := net.Listen("tcp", ":0")
	el := NewEventListener(l)
	hasErr := false
	el.OnAcceptError(func(l *EventListener, err error) {
		hasErr = true
	}).Start()
	el.Close()
	time.Sleep(time.Second)
	assert.True(t, hasErr)
}

type initErrorCodec struct {
	net.Conn
}

func (i *initErrorCodec) Initialize(context ConnContext) error {
	return fmt.Errorf("init error")
}

type initOkayCodec struct {
	net.Conn
}

func (i *initOkayCodec) Initialize(context ConnContext) error {
	return nil
}

func (i *initOkayCodec) Close() error {
	return nil
}

func TestNewEventListener_OnCodecError(t *testing.T) {
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	hasErr := false
	el := NewEventListener(l)
	el.AddCodecFactory(func() Codec {
		return &initOkayCodec{}
	})
	el.AddCodecFactory(func() Codec {
		return &initErrorCodec{}
	})
	el.OnCodecInitializeError(func(l *EventListener, err error) {
		hasErr = true
	}).Start()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c.Write([]byte("hello"))
	time.Sleep(time.Second)
	assert.True(t, hasErr)
}

func TestNewEventListener(t *testing.T) {
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	el.SetConnContextFactory(func(conn net.Conn) Context {
		ctx:= NewContext()
		ctx.SetData("cfg","abc")
		return ctx
	})
	el.AddCodecFactory(func() Codec {
		return NewAESCodec("abc")
	}).SetFirstReadTimeout(time.Second)
	el.OnAccept(func(_ *EventListener,ctx Context, c net.Conn) {
		c.Write([]byte("hello"))
		assert.Equal(t,"abc",c.(*Conn).ctx.Data("cfg"))
	}).Start()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	hasData := false
	ec := NewEventConn(c)
	ec.AddCodec(NewAESCodec("abc"))
	ec.OnData(func(ec *EventConn, data []byte) {
		hasData = true
		assert.Equal(t, "hello", string(data))
	}).Start()
	ec.Write([]byte("hello"))
	time.Sleep(time.Second)
	assert.True(t, hasData)
}
