// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"strings"
	"sync"
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
		OnFistReadTimeout(func(c net.Conn, err error) {
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
		OnFistReadTimeout(func(c net.Conn, err error) {
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
	el.AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		return c, nil
	})
	el.AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		return nil, ErrListenerFilterSkipped
	})
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

func TestNewEventListener_ListenerFilterError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	hasErr := false
	el.AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
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

func TestListener_Accept(t *testing.T) {
	t.Parallel()
	l, p, _ := ListenRandom("")
	l1 := NewListener(NewListener(l))
	addr := ""
	cnt0 := int64(0)
	cnt1 := int64(0)
	var hasErr error
	l1.SetOnConnClose(func(ctx ConnContext) {
		addr = ctx.Conn().LocalAddr().String()
		cnt1 = l1.ConnCount()
	})
	l1.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		cnt0 = l1.ConnCount()
		return c, nil
	})
	l1.SetAutoCloseConnOnReadWriteError(true)
	g := sync.WaitGroup{}
	g.Add(2)

	go func() {
		defer g.Done()
		c, _ := l1.Accept()
		go func() {
			defer g.Done()
			Read(c, 1)
		}()
		c.Write([]byte("hello"))
		time.Sleep(time.Millisecond * 200)
		_, hasErr = c.Write([]byte("_"))
	}()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", ":"+p)
	str, _ := Read(c, 5)
	assert.Equal(t, "hello", str)
	c.Close()
	g.Wait()
	assert.Equal(t, int64(1), cnt0)
	assert.Equal(t, int64(0), cnt1)
	t.Log(addr)
	assert.True(t, strings.HasPrefix(addr, "127.0.0.1"))
	assert.Error(t, hasErr)
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

func TestNewEventListener_OnConnClose(t *testing.T) {
	t.Parallel()
	l, p, _ := ListenRandom("")
	el := NewEventListener(l)
	cnt := int64(0)
	cnt1 := int64(0)
	el.SetAutoCloseConnOnReadWriteError(true)
	el.OnAccept(func(l *EventListener, ctx Context, c net.Conn) {
		cnt = el.ConnCount()
	})
	el.SetOnConnClose(func(ctx ConnContext) {
		cnt1 = el.ConnCount()
	})
	el.Start()
	net.Dial("tcp", ":"+p)
	time.Sleep(time.Millisecond * 300)
	assert.Equal(t, int64(1), cnt)
	assert.Equal(t, int64(0), cnt1)
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
