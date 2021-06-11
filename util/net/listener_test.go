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
	t.Parallel()
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

func TestNewEventListener_OnFistReadTimeout2(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	timeout := false
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		c = NewBufferedConn(c)
		return c, nil
	})
	el.SetFirstReadTimeout(time.Millisecond * 100).
		OnFistReadTimeout(func(l *EventListener, c net.Conn, err error) {
			timeout = true
		}).Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, timeout)
}

func TestNewEventListener_FilterError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	hasErr := false
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		c = NewBufferedConn(c)
		return c, nil
	})
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		return nil, ErrConnFilterSkipped
	})
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		return nil, fmt.Errorf("error")
	})
	el.OnAcceptError(func(l *EventListener, err error) {
		hasErr = true
		assert.Equal(t, "error", err.Error())
	}).Start()
	_, err := net.Dial("tcp", "127.0.0.1:"+p)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	assert.True(t, hasErr)
}

func TestNewEventListener_MissingAccept(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	NewEventListener(l).Start()
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Millisecond * 1000)
	buf := make([]byte, 1)
	_, err := c.Read(buf)
	assert.Error(t, err)
}

func TestNewEventListener_OnAcceptError(t *testing.T) {
	t.Parallel()
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

type initSkippedCodec struct {
	net.Conn
}

func (i *initSkippedCodec) Initialize(context ConnContext) error {
	return ErrCodecSkipped
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
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	hasErr := false
	el := NewEventListener(l)
	el.AddCodecFactory(func() Codec {
		return &initOkayCodec{}
	})
	el.AddCodecFactory(func() Codec {
		return &initSkippedCodec{}
	})
	el.AddCodecFactory(func() Codec {
		return &initErrorCodec{}
	})
	el.OnAcceptError(func(l *EventListener, err error) {
		hasErr = true
	}).Start()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c.Write([]byte("hello"))
	time.Sleep(time.Second)
	assert.True(t, hasErr)
}

func TestNewEventListener(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	var conn net.Conn
	el := NewEventListener(l)
	el.SetConnContextFactory(func(conn net.Conn) Context {
		ctx := NewContext()
		ctx.SetData("cfg", "abc")
		return ctx
	})
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		c = NewBufferedConn(c)
		ctx.SetData("bufConn", c)
		return c, nil
	})
	el.AddCodecFactory(func() Codec {
		return NewAESCodec("abc")
	}).SetFirstReadTimeout(time.Second)
	el.OnAccept(func(_ *EventListener, ctx Context, c net.Conn) {
		conn = ctx.Data("bufConn").(net.Conn)
		c.Write([]byte("hello"))
		assert.Equal(t, "abc", c.(*Conn).ctx.Data("cfg"))
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
	assert.Implements(t, (*BufferedConn)(nil), conn)
}
