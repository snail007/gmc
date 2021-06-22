// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultipleCodec(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	outputCnt := new(int32)
	password := "abc"
	debug := false
	go func() {
		c, _ := l.Accept()
		conn := NewConn(c)
		conn.AddCodec(NewHeartbeatCodec())
		conn.AddCodec(NewAESCodec(password))
		assert.Equal(t, conn.ctx, conn.Ctx())
		go func() {
			for {
				n, err := conn.Write([]byte("hello from server"))
				if err != nil {
					fmt.Printf("server write error %s, %d\n", err, n)
					return
				}
				time.Sleep(time.Millisecond * 100)
			}
		}()
		go func() {
			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println("server read error", err)
					return
				}
				assert.Equal(t, "hello from client", string(buf[:n]))
				atomic.AddInt32(outputCnt, 1)
				if debug {
					fmt.Println(string(buf[:n]))
				}
			}
		}()
		select {}
	}()
	time.Sleep(time.Second * 2)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn := NewConn(c)
	conn.AddCodec(NewHeartbeatCodec())
	conn.AddCodec(NewAESCodec(password))
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("client read error", err)
				return
			}
			atomic.AddInt32(outputCnt, 1)
			assert.Equal(t, "hello from server", string(buf[:n]))
			if debug {
				fmt.Println(string(buf[:n]))
			}
		}
	}()
	go func() {
		for {
			n, err := conn.Write([]byte("hello from client"))
			if err != nil {
				fmt.Printf("client write error %s, %d\n", err, n)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	time.Sleep(time.Second * 3)
	assert.True(t, *outputCnt > 50)
	assert.IsType(t, (*net.TCPConn)(nil), conn.RawConn())
	assert.IsType(t, (*AESCodec)(nil), conn.Conn)
}

func TestMultipleCodec2(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	outputCnt := new(int32)
	password := "abc"
	debug := false
	go func() {
		c, _ := l.Accept()
		conn := NewConn(c)
		conn.AddCodec(NewAESCodec(password))
		conn.AddCodec(NewHeartbeatCodec())
		go func() {
			for {
				_, err := conn.Write([]byte("hello from server"))
				if err != nil {
					fmt.Println("server write error", err)
					return
				}
				time.Sleep(time.Millisecond * 100)
			}
		}()
		go func() {
			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println("server read error", err)
					return
				}
				assert.Equal(t, "hello from client", string(buf[:n]))
				atomic.AddInt32(outputCnt, 1)
				if debug {
					fmt.Println(string(buf[:n]))
				}
			}
		}()
		select {}
	}()
	time.Sleep(time.Second * 2)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn := NewConn(c)
	conn.AddCodec(NewAESCodec(password))
	conn.AddCodec(NewHeartbeatCodec())
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("client read error", err)
				return
			}
			atomic.AddInt32(outputCnt, 1)
			assert.Equal(t, "hello from server", string(buf[:n]))
			if debug {
				fmt.Println(string(buf[:n]))
			}
		}
	}()
	go func() {
		for {
			_, err := conn.Write([]byte("hello from client"))
			if err != nil {
				fmt.Println("client write error", err)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	time.Sleep(time.Second * 3)
	assert.True(t, conn.ReadBytes() > 450)
	assert.True(t, conn.WriteBytes() > 450)
	assert.True(t, *outputCnt > 50)
}

func TestMultipleCodec3(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := new(bool)
	el := NewEventListener(l)
	var conn net.Conn
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(new(bool))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitCodec2(called)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(new(bool))
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.(*Conn).doInitialize()
		conn = c.(*Conn).Conn
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, *called)
	assert.IsType(t, &initCodec{}, conn)
}

func TestEventConn(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	go func() {
		c, _ := l.Accept()
		c.Write([]byte("hello"))
		time.AfterFunc(time.Second, func() {
			c.Close()
		})
	}()
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn := NewEventConn(c)
	conn.SetReadBufferSize(1024)
	conn.SetTimeout(time.Second)
	assert.Equal(t, time.Second, conn.ReadTimeout())
	assert.Equal(t, time.Second, conn.WriteTimeout())
	closed := false
	conn.OnData(func(ctx Context, data []byte) {
		s := ctx.EventConn()
		assert.Equal(t, "hello", string(data))
		assert.Equal(t, int64(5), s.ReadBytes())
		s.Write([]byte("hi"))
		assert.Equal(t, int64(2), s.WriteBytes())
	}).OnClose(func(ctx Context) {
		closed = true
	}).Start()
	time.AfterFunc(time.Millisecond*1500, func() {
		conn.Close()
	})
	time.Sleep(time.Second * 3)
	assert.Contains(t, conn.LocalAddr().String(), "127.0.0.1:")
	assert.Contains(t, conn.RemoteAddr().String(), "127.0.0.1:")
	assert.True(t, closed)
	assert.Equal(t, conn.ctx, conn.Ctx())
}

func TestEventConn2(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	go func() {
		c, _ := l.Accept()
		c.Write([]byte("hello"))
		time.AfterFunc(time.Second, func() {
			c.Close()
		})
	}()
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c.Close()
	conn := NewEventConn(NewConn(NewConn(c)))
	conn.SetReadBufferSize(1024)
	conn.SetTimeout(time.Second)
	assert.Equal(t, time.Second, conn.ReadTimeout())
	assert.Equal(t, time.Second, conn.WriteTimeout())
	closed := false
	readErr := false
	writeErr := false
	conn.AddConnFilter(func(ctx Context, c net.Conn, next NextConnFilter) (net.Conn, error) {
		return next.Call(ctx, c)
	})
	conn.AddConnFilter(func(ctx Context, c net.Conn, next NextConnFilter) (net.Conn, error) {
		// just return c, skip the next
		return c, nil
	})
	conn.OnData(func(ctx Context, data []byte) {
		s := ctx.EventConn()
		assert.Equal(t, "hello", string(data))
		assert.Equal(t, int64(5), s.ReadBytes())
		time.Sleep(time.Millisecond * 300)
		s.Write([]byte("hi"))
		assert.Equal(t, int64(2), s.WriteBytes())
	}).OnReadError(func(ctx Context, err error) {
		readErr = true
	}).OnWriterError(func(ctx Context, err error) {
		writeErr = true
	}).OnClose(func(ctx Context) {
		closed = true
	})
	conn.Write([]byte("hi"))
	conn.Start()
	time.AfterFunc(time.Millisecond*1500, func() {
		conn.Close()
	})
	time.Sleep(time.Second * 3)
	assert.Contains(t, conn.LocalAddr().String(), "127.0.0.1:")
	assert.Contains(t, conn.RemoteAddr().String(), "127.0.0.1:")
	assert.True(t, closed)
	assert.True(t, readErr)
	assert.True(t, writeErr)
}

func TestEventConn3(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	go func() {
		c, _ := l.Accept()
		c.Write([]byte("hello"))
	}()
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn := NewEventConn(c)
	conn.AddCodec(newInitPassThroughCodec(new(bool)))
	conn.AddCodec(newInitErrorCodec(new(bool)))
	conn.OnReadError(func(ctx Context, err error) {
		ctx.SetData("test", "abc")
	})
	conn.Start()
	time.Sleep(time.Second * 1)
	assert.Equal(t, "abc", conn.Ctx().Data("test"))
}

func TestBufferedConn_PeekMax(t *testing.T) {
	t.Parallel()
	l, p, err := RandomListen("")
	assert.NoError(t, err)
	var d []byte
	var n int
	el := NewEventListener(l).AddListenerFilter(func(ctx Context, c net.Conn, next NextConnFilter) (net.Conn, error) {
		bc := NewBufferedConn(NewBufferedConn(c))
		d, err = bc.PeekMax(1024)
		c = bc
		n = bc.Buffered()
		s, _ := Read(bc, 10)
		assert.Equal(t, "hello", s)
		return next.Call(ctx, c)
	})
	el.Start()
	time.Sleep(time.Second)
	err = Write("127.0.0.1:"+p, "hello")
	time.Sleep(time.Second)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 200)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(d))
	assert.Equal(t, n, 5)
	Write(":0", "error")
}

func TestNewConnBinder(t *testing.T) {
	t.Parallel()
	l1, p1, err := RandomListen("")
	assert.NoError(t, err)
	l2, p2, err := RandomListen("")
	assert.NoError(t, err)
	closed := false
	srcClosed := false
	dstClosed := false
	str := ""
	NewEventListener(l2).OnAccept(func(ctx Context, c net.Conn) {
		str, _ = Read(c, 3)
	}).Start()
	NewEventListener(l1).OnAccept(func(ctx Context, c net.Conn) {
		c2, _ := net.Dial("tcp", "127.0.0.1:"+p2)
		b := NewConnBinder(c, c2).OnClose(func() {
			closed = true
		}).OnSrcClose(func(ctx Context) {
			srcClosed = true
		}).OnDstClose(func(ctx Context) {
			dstClosed = true
		}).SetReadBufSize(100)
		b.Start()
		assert.Equal(t, b, b.Ctx().ConnBinder())
	}).Start()
	time.Sleep(time.Second)
	err = Write("127.0.0.1:"+p1, "hello")
	assert.NoError(t, err)
	time.Sleep(time.Second * 2)
	assert.Equal(t, "hel", str)
	assert.True(t, closed)
	assert.True(t, srcClosed)
	assert.True(t, dstClosed)
}

func TestConn_FilterHijacked(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := false
	hijacked := false
	el := NewEventListener(l)
	el.AddConnFilter(func(ctx Context, c net.Conn, next NextConnFilter) (net.Conn, error) {
		hijacked = true
		return ctx.Hijack()
	})
	el.AddConnFilter(func(ctx Context, c net.Conn, next NextConnFilter) (net.Conn, error) {
		called = true
		return next.Call(ctx, c)
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		// trigger lazy initialization
		c.Write([]byte(""))
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Millisecond * 500)
	assert.True(t, hijacked)
	assert.False(t, called)
}

func TestConn_CodecHijacked(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := new(bool)
	hijacked := new(bool)
	el := NewEventListener(l)
	var conn net.Conn
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(new(bool))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitHijackedCodec(hijacked)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(called)
	})
	accepted := false
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.(*Conn).doInitialize()
		a, b := ctx.Hijack()
		assert.Nil(t, a)
		assert.Equal(t, errHijackedAlready, b)
		assert.False(t, ctx.IsTLS())
		conn = c.(*Conn).Conn
		accepted = true
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, *hijacked)
	assert.False(t, *called)
	assert.True(t, accepted)
	assert.IsType(t, &initHijackedCodec{}, conn)
}

func TestConn_CodecHijackedFail(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := new(bool)
	hijacked := new(bool)
	var hasError error
	el := NewEventListener(l)

	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(new(bool))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitHijackedFailCodec(hijacked, nil)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(called)
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		_, hasError = Read(c, 1)
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, *hijacked)
	assert.False(t, *called)
	assert.Equal(t, errHijackedFail, hasError)
}

func TestConn_CodecHijackedFail1(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := new(bool)
	hijacked := new(bool)
	var hasError error
	el := NewEventListener(l)

	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(new(bool))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitHijackedFailCodec(hijacked, nil)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(called)
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		//if codec init error, Conn will close the conn, so here read EOF
		_, hasError = Read(c, 1)
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, *hijacked)
	assert.False(t, *called)
	assert.Error(t, hasError)
}

func TestConn_CodecHijackedFail2(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := new(bool)
	hijacked := new(bool)
	var hasError error
	el := NewEventListener(l)

	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(new(bool))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitHijackedFailCodec(hijacked, func(ctx Context) {
			ctx.Hijack(nil, nil)
		})
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(called)
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		_, hasError = Read(c, 1)
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, *hijacked)
	assert.False(t, *called)
	assert.Equal(t, errHijackedFail, hasError)
}

func TestConn_CodecError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := false
	hasErr := new(bool)
	el := NewEventListener(l)
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(new(bool))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitErrorCodec(hasErr)
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.(*Conn).doInitialize()
	})
	el.Start()
	time.Sleep(time.Second)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, *hasErr)
	assert.False(t, called)
}

func TestConn_CodecError1(t *testing.T) {
	t.Parallel()
	c, _ := net.Dial("tcp", ":")
	c0 := NewConn(c)
	c0.AddCodec(newInitErrorCodec(new(bool)))
	_, err := c0.Write([]byte(""))
	assert.Error(t, err)
}
