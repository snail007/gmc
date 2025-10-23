// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/snail007/gmc/util/sync/atomic"

	"github.com/stretchr/testify/assert"
)

func TestMultipleCodec(t *testing.T) {
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	password := "abc"
	g := sync.WaitGroup{}
	g.Add(4)
	go func() {
		c, _ := l.Accept()
		conn := NewConn(c)
		conn.AddCodec(NewHeartbeatCodec())
		conn.AddCodec(NewAESCodec(password))
		assert.Equal(t, conn.ctx, conn.Ctx())
		go func() {
			defer g.Done()
			conn.Write([]byte("hello from server"))
		}()
		go func() {
			defer g.Done()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("server read error", err)
				return
			}
			assert.Equal(t, "hello from client", string(buf[:n]))
		}()
	}()
	time.Sleep(time.Second * 2)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn := NewConn(c)
	conn.AddCodec(NewHeartbeatCodec())
	conn.AddCodec(NewAESCodec(password))
	go func() {
		defer g.Done()
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client read error", err)
			return
		}
		assert.Equal(t, "hello from server", string(buf[:n]))
	}()
	go func() {
		defer g.Done()
		conn.Write([]byte("hello from client"))
	}()
	g.Wait()
	assert.IsType(t, (*net.TCPConn)(nil), conn.RawConn())
	assert.IsType(t, (*AESCodec)(nil), conn.Conn)
}

func TestMultipleCodec2(t *testing.T) {
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	password := "abc"
	g := sync.WaitGroup{}
	g.Add(4)
	go func() {
		c, _ := l.Accept()
		conn := NewConn(c)
		conn.AddCodec(NewAESCodec(password))
		conn.AddCodec(NewHeartbeatCodec())
		go func() {
			defer g.Done()
			conn.Write([]byte("hello from server"))
		}()
		go func() {
			defer g.Done()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("server read error", err)
				return
			}
			assert.Equal(t, "hello from client", string(buf[:n]))
		}()
		select {}
	}()
	time.Sleep(time.Second * 2)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn := NewConn(c)
	conn.AddCodec(NewAESCodec(password))
	conn.AddCodec(NewHeartbeatCodec())
	go func() {
		defer g.Done()
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client read error", err)
			return
		}
		assert.Equal(t, "hello from server", string(buf[:n]))
	}()
	go func() {
		defer g.Done()
		conn.Write([]byte("hello from client"))
	}()
	g.Wait()
}

func TestMultipleCodec3(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	el := NewEventListener(l)
	conn := &gatomic.Value{}
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitCodec2(called)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.(*Conn).doInitialize()
		conn.SetVal(c.(*Conn).Conn)
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, called.IsTrue())
	assert.IsType(t, &initCodec{}, conn.Val())
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
	closed := gatomic.NewBool(false)
	conn.OnData(func(ctx Context, data []byte) {
		s := ctx.EventConn()
		assert.Equal(t, "hello", string(data))
		assert.Equal(t, int64(5), s.ReadBytes())
		s.Write([]byte("hi"))
		assert.Equal(t, int64(2), s.WriteBytes())
	}).OnClose(func(ctx Context) {
		closed.SetTrue()
	}).Start()
	conn.Start()
	time.AfterFunc(time.Millisecond*1500, func() {
		conn.Close()
	})
	time.Sleep(time.Second * 3)
	assert.Contains(t, conn.LocalAddr().String(), "127.0.0.1:")
	assert.Contains(t, conn.RemoteAddr().String(), "127.0.0.1:")
	assert.True(t, closed.IsTrue())
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
	closed := gatomic.NewBool(false)
	readErr := gatomic.NewBool(false)
	writeErr := gatomic.NewBool(false)
	conn.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		return c, nil
	})
	conn.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
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
		readErr.SetTrue()
	}).OnWriterError(func(ctx Context, err error) {
		writeErr.SetTrue()
	}).OnClose(func(ctx Context) {
		closed.SetTrue()
	})
	conn.Write([]byte("hi"))
	conn.Start()
	time.AfterFunc(time.Millisecond*1500, func() {
		conn.Close()
	})
	time.Sleep(time.Second * 3)
	assert.Contains(t, conn.LocalAddr().String(), "127.0.0.1:")
	assert.Contains(t, conn.RemoteAddr().String(), "127.0.0.1:")
	assert.True(t, closed.IsTrue())
	assert.True(t, readErr.IsTrue())
	assert.True(t, writeErr.IsTrue())
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
	conn.AddCodec(newInitPassThroughCodec(gatomic.NewBool(false)))
	conn.AddCodec(newInitErrorCodec(gatomic.NewBool(false)))
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
	var str = gatomic.NewString("")
	var n int
	el := NewEventListener(l).AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		bc := NewBufferedConn(NewBufferedConn(c))
		d, err := bc.PeekMax(1024)
		str.SetVal(string(d))
		c = bc
		n = bc.Buffered()
		s, _ := Read(bc, 10)
		assert.Equal(t, "hello", s)
		return c, err
	})
	el.Start()
	time.Sleep(time.Second)
	err = Write("127.0.0.1:"+p, "hello")
	time.Sleep(time.Second)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond * 200)
	assert.NoError(t, err)
	assert.Equal(t, "hello", str.Val())
	assert.Equal(t, n, 5)
	Write(":0", "error")
}

func TestNewConnBinder1(t *testing.T) {
	t.Parallel()
	l1, p1, err := RandomListen("")
	assert.NoError(t, err)
	l2, p2, err := RandomListen("")
	assert.NoError(t, err)
	closed := gatomic.NewBool(false)
	srcClosed := gatomic.NewBool(false)
	dstClosed := gatomic.NewBool(false)
	str := gatomic.NewString("")
	NewEventListener(l2).OnAccept(func(ctx Context, c net.Conn) {
		s, _ := Read(c, 3)
		str.SetVal(s)
	}).Start()
	NewEventListener(l1).OnAccept(func(ctx Context, c net.Conn) {
		c2, _ := net.Dial("tcp", "127.0.0.1:"+p2)
		b := NewConnBinder(c, c2).OnClose(func() {
			closed.SetTrue()
		}).OnSrcClose(func(ctx Context) {
			srcClosed.SetTrue()
		}).OnDstClose(func(ctx Context) {
			dstClosed.SetTrue()
		}).SetReadBufSize(100)
		b.Start()
		assert.Equal(t, b, b.Ctx().ConnBinder())
	}).Start()
	time.Sleep(time.Second)
	err = Write("127.0.0.1:"+p1, "hello")
	assert.NoError(t, err)
	time.Sleep(time.Second * 2)
	assert.Equal(t, "hel", str.Val())
	assert.True(t, closed.IsTrue())
	assert.True(t, srcClosed.IsTrue())
	assert.True(t, dstClosed.IsTrue())
}

func TestNewConnBinder2(t *testing.T) {
	t.Parallel()
	l1, p1, err := RandomListen("")
	assert.NoError(t, err)
	l2, p2, err := RandomListen("")
	assert.NoError(t, err)
	str := gatomic.NewString("")
	NewEventListener(l2).OnAccept(func(ctx Context, c net.Conn) {
		c.Close()
	}).Start()
	NewEventListener(l1).OnAccept(func(ctx Context, c net.Conn) {
		c2, _ := net.Dial("tcp", "127.0.0.1:"+p2)
		b := NewConnBinder(c, c2)
		b.StartAndWait()
		str.SetVal(b.Error().Error())
	}).Start()
	time.Sleep(time.Second)
	Write("127.0.0.1:"+p1, "hello")
	time.Sleep(time.Second)
	assert.NotEmpty(t, str.Val())
}

func TestConn_FilterBreak(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	isBreak := gatomic.NewBool(false)
	el := NewEventListener(l)
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		isBreak.SetTrue()
		ctx.Break()
		return c, nil
	})
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		called.SetTrue()
		return c, nil
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		// trigger lazy initialization
		c.Write([]byte(""))
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Millisecond * 500)
	assert.True(t, isBreak.IsTrue())
	assert.False(t, called.IsTrue())
}
func TestConn_FilterContinue(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	isContinue := gatomic.NewBool(false)
	el := NewEventListener(l)
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		isContinue.SetTrue()
		ctx.Continue()
		return c, nil
	})
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		called.SetTrue()
		return c, nil
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		// trigger lazy initialization
		c.Write([]byte(""))
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Millisecond * 500)
	assert.True(t, isContinue.IsTrue())
	assert.True(t, called.IsTrue())
}

type breakCodec struct {
	net.Conn
	called *gatomic.Bool
}

func (i *breakCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}

func (i *breakCodec) Initialize(ctx Context) error {
	i.called.SetTrue()
	ctx.Break()
	return nil
}

func newBreakCodec(called *gatomic.Bool) *breakCodec {
	return &breakCodec{
		called: called,
	}
}

type continueCodec struct {
	net.Conn
	called *gatomic.Bool
}

func (i *continueCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}

func (i *continueCodec) Initialize(ctx Context) error {
	i.called.SetTrue()
	ctx.Continue()
	return nil
}

func newContinueCodec(called *gatomic.Bool) *continueCodec {
	return &continueCodec{
		called: called,
	}
}

type breakFailCodec struct {
	net.Conn
	called    *gatomic.Bool
	modifyCtx func(ctx Context)
}

func (i *breakFailCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}

func (i *breakFailCodec) Initialize(ctx Context) error {
	i.called.SetTrue()
	if i.modifyCtx != nil {
		i.modifyCtx(ctx)
	}
	ctx.Break()
	return fmt.Errorf("break fail")
}

func newBreakFailCodec(called *gatomic.Bool, m func(ctx Context)) *breakFailCodec {
	return &breakFailCodec{
		called:    called,
		modifyCtx: m,
	}
}

func TestConn_CodecBreak(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	isBreak := gatomic.NewBool(false)
	el := NewEventListener(l)
	var conn net.Conn
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newBreakCodec(isBreak)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(called)
	})
	accepted := gatomic.NewBool(false)
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.(*Conn).doInitialize()
		assert.False(t, ctx.IsTLS())
		conn = c.(*Conn).Conn
		accepted.SetTrue()
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, isBreak.IsTrue())
	assert.False(t, called.IsTrue())
	assert.True(t, accepted.IsTrue())
	assert.IsType(t, &breakCodec{}, conn)
}
func TestConn_CodecContinue(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	isContinue := gatomic.NewBool(false)
	el := NewEventListener(l)
	var conn = gatomic.NewAny("")
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		called.SetTrue()
		return newContinueCodec(isContinue)
	})
	accepted := gatomic.NewBool(false)
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.(*Conn).doInitialize()
		assert.False(t, ctx.IsTLS())
		conn.SetVal(c.(*Conn).Conn)
		accepted.SetTrue()
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, isContinue.IsTrue())
	assert.True(t, called.IsTrue())
	assert.True(t, accepted.IsTrue())
	assert.IsType(t, &initPassThroughCodec{}, conn.Val())
}

func TestConn_CodecBreakFail(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	isBreak := gatomic.NewBool(false)
	var hasError = gatomic.NewBool(false)
	el := NewEventListener(l)

	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newBreakFailCodec(isBreak, nil)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(called)
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		_, e := Read(c, 1)
		if e != nil {
			hasError.SetTrue()
		}
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, isBreak.IsTrue())
	assert.False(t, called.IsTrue())
	assert.True(t, hasError.IsTrue())
}

func TestConn_CodecBreakFail1(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	isBreak := gatomic.NewBool(false)
	var hasError = gatomic.NewBool(false)
	el := NewEventListener(l)

	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newBreakFailCodec(isBreak, nil)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(called)
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		//if codec init error, Conn will close the conn, so here read EOF
		_, e := Read(c, 1)
		if e != nil {
			hasError.SetTrue()
		}
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, isBreak.IsTrue())
	assert.False(t, called.IsTrue())
	assert.True(t, hasError.IsTrue())
}

func TestConn_CodecBreakFail2(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	isBreak := gatomic.NewBool(false)
	var hasError = gatomic.NewBool(false)
	el := NewEventListener(l)

	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newBreakFailCodec(isBreak, func(ctx Context) {
			ctx.Break()
		})
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(called)
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		_, e := Read(c, 1)
		if e != nil {
			hasError.SetTrue()
		}
	})
	el.Start()
	time.Sleep(time.Millisecond * 500)
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, isBreak.IsTrue())
	assert.False(t, called.IsTrue())
	assert.True(t, hasError.IsTrue())
}

func TestConn_CodecError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	hasErr := gatomic.NewBool(false)
	el := NewEventListener(l)
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
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
	assert.True(t, hasErr.IsTrue())
	assert.False(t, called.IsTrue())
}

func TestConn_CodecError1(t *testing.T) {
	t.Parallel()
	c, _ := net.Dial("tcp", ":")
	c0 := NewConn(c)
	c0.AddCodec(newInitErrorCodec(gatomic.NewBool(false)))
	_, err := c0.Write([]byte(""))
	assert.Error(t, err)
}

func TestConn_MaxIdleTimeout1(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	el := NewEventListener(l)
	el.AddCodecFactory(func(ctx Context) Codec {
		return NewHeartbeatCodec().SetInterval(time.Millisecond * 100)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return NewAESCodec("123")
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.Write([]byte("hello"))
	})
	el.Start()
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c0 := NewConn(c)
	c0.SetMaxIdleTimeout(time.Second)
	c0.SetOnIdleTimeout(func(conn *Conn) {
		called.SetTrue()
	})
	c0.AddCodec(NewHeartbeatCodec().SetInterval(time.Millisecond * 100))
	c0.AddCodec(NewAESCodec("123"))
	resp := ""
	go func() {
		buf := make([]byte, 1024)
		n, _ := c0.Read(buf)
		resp = string(buf[:n])
	}()
	time.Sleep(time.Second * 3)
	assert.True(t, called.IsTrue())
	assert.Equal(t, resp, "hello")
	assert.True(t, c0.TouchTime().Add(time.Second).Before(time.Now()))
}

func TestConn_MaxIdleTimeout2(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called := gatomic.NewBool(false)
	el := NewEventListener(l)
	el.AddCodecFactory(func(ctx Context) Codec {
		return NewHeartbeatCodec().SetInterval(time.Millisecond * 100)
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return NewAESCodec("123")
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.Write([]byte("hello"))
	})
	el.Start()
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c0 := NewConn(c)
	c0.SetMaxIdleTimeout(time.Second)
	c0.SetOnIdleTimeout(func(conn *Conn) {
		called.SetTrue()
	})
	c0.AddCodec(NewHeartbeatCodec().SetInterval(time.Millisecond * 100))
	c0.AddCodec(NewAESCodec("123"))
	resp := ""
	go func() {
		buf := make([]byte, 1024)
		n, _ := c0.Read(buf)
		resp = string(buf[:n])
	}()
	time.Sleep(time.Millisecond * 400)
	c0.Close()
	assert.True(t, called.IsFalse())
	assert.Equal(t, resp, "hello")
	assert.False(t, c0.TouchTime().Add(time.Second).Before(time.Now()))
}

func TestConn_Hook(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	called1 := gatomic.NewBool(false)
	called2 := gatomic.NewBool(false)
	el := NewEventListener(l)
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.Write([]byte("hello"))
	})
	el.Start()
	time.Sleep(time.Second)
	c, _ := Dial("127.0.0.1:"+p, time.Second)
	c.SetOnClose(func(conn *Conn) {
		called1.SetTrue()
	})
	c.SetOnInitialize(func(conn *Conn) {
		called2.SetTrue()
	})
	buf := make([]byte, 1024)
	c.Read(buf)
	c.Close()
	assert.True(t, called1.IsTrue())
	assert.True(t, called2.IsTrue())
}
