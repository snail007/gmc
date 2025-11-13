// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/snail007/gmc/util/sync/atomic"

	"github.com/stretchr/testify/assert"
)

func TestNewEventListener_OnFirstReadTimeout(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	timeout := gatomic.NewBool(false)
	el.SetFirstReadTimeout(time.Millisecond * 100).
		OnFirstReadTimeout(func(ctx Context, c net.Conn, err error) {
			timeout.SetTrue()
		}).Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, timeout.IsTrue())
}

func TestNewEventListener_OnFirstReadTimeout2(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	timeout := gatomic.NewBool(false)
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		c = NewBufferedConn(c)
		return c, nil
	})
	el.SetFirstReadTimeout(time.Millisecond * 100).
		OnFirstReadTimeout(func(ctx Context, c net.Conn, err error) {
			timeout.SetTrue()
		}).Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, timeout.IsTrue())
}

func TestNewEventListener_BeforeFirstRead_Success(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	timeout := gatomic.NewBool(false)
	beforeRead := gatomic.NewBool(false)
	el.BeforeFirstRead(func(ctx Context, c net.Conn) error {
		beforeRead.SetTrue()
		return nil
	})
	el.SetFirstReadTimeout(time.Millisecond * 100).
		OnFirstReadTimeout(func(ctx Context, c net.Conn, err error) {
			timeout.SetTrue()
		}).Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, timeout.IsTrue())
	assert.True(t, beforeRead.IsTrue())
}

func TestNewEventListener_BeforeFirstRead_Error(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	timeout := gatomic.NewBool(false)
	beforeRead := gatomic.NewBool(false)
	el.BeforeFirstRead(func(ctx Context, c net.Conn) error {
		beforeRead.SetTrue()
		return fmt.Errorf("before read error")
	})
	el.SetFirstReadTimeout(time.Millisecond * 100).
		OnFirstReadTimeout(func(ctx Context, c net.Conn, err error) {
			timeout.SetTrue()
		}).Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.False(t, timeout.IsTrue())
	assert.True(t, beforeRead.IsTrue())
}

func TestNewEventListener_Hijacked(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	called := gatomic.NewBool(false)
	isHijacked := gatomic.NewBool(false)
	el.AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		// isHijacked the conn, do anything with c
		isHijacked.SetTrue()
		ctx.Hijack()
		return nil, nil
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		called.SetTrue()
	})
	el.Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, called.IsFalse())
	assert.True(t, isHijacked.IsTrue())
}

func TestNewEventListener_FilterError(t *testing.T) {
	t.Parallel()
	el, err := NewEventListenerAddr(":0")
	assert.Nil(t, err)
	p := NewAddr(el.Addr()).Port()
	var hasErr = gatomic.NewBool(false)
	el.AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		return c, nil
	})
	el.AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		//skip the next filer
		return c, nil
	})
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		c = NewBufferedConn(c)
		return c, nil
	})
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		return nil, fmt.Errorf("error")
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		_, e := Read(c, 5)
		if e != nil && e.Error() == "error" {
			hasErr.SetTrue()
		}
	}).Start()
	_, e := net.Dial("tcp", "127.0.0.1:"+p)
	assert.NoError(t, e)
	time.Sleep(time.Second)
	assert.True(t, hasErr.IsTrue())
}

func TestNewEventListener_ListenerFilterError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	hasErr := gatomic.NewBool(false)
	el.AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		return nil, fmt.Errorf("error")
	})
	el.OnAcceptError(func(ctx Context, err error) {
		hasErr.SetTrue()
		assert.Equal(t, "error", err.Error())
	}).Start()
	_, err := net.Dial("tcp", "127.0.0.1:"+p)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	assert.True(t, hasErr.IsTrue())
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
	l, p, _ := RandomListen("")
	l1 := NewListener(NewListener(l))
	addr := ""
	cnt0 := int64(0)
	cnt1 := int64(0)
	var hasErr error
	l1.SetOnConnClose(func(ctx Context) {
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
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
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

type errListener struct {
	net.Listener
	cnt int
	t   *testing.T
}

func newErrListener(listener net.Listener) *errListener {
	return &errListener{
		Listener: listener,
	}
}

func (s *errListener) Accept() (conn net.Conn, err error) {
	defer func() { s.cnt++ }()
	if s.cnt == 10 || s.cnt == 0 {
		return s.Listener.Accept()
	} else if s.cnt > 10 {
		return nil, fmt.Errorf("error")
	}
	s.Listener.Accept()
	return nil, newTempNetError()
}

func TestListener_AcceptError(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := newErrListener(l)
	l0.t = t
	l1 := NewEventListener(NewListener(l0))
	l1.SetAutoCloseConnOnReadWriteError(true)
	var acceptCnt = new(int32)
	l1.OnAccept(func(ctx Context, c net.Conn) {
		atomic.AddInt32(acceptCnt, 1)
	})
	l1.Start()
	time.Sleep(time.Second * 1)
	net.Dial("tcp", "127.0.0.1:"+p)
	for i := 0; i <= 13; i++ {
		net.Dial("tcp", "127.0.0.1:"+p)
		time.Sleep(time.Millisecond * 300)
	}
	assert.Equal(t, int32(2), atomic.LoadInt32(acceptCnt))
	assert.Equal(t, l1, l1.Ctx().EventListener())
}

func TestNewEventListener_OnAcceptError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	el := NewEventListener(l)
	hasErr := gatomic.NewBool(false)
	el.OnAcceptError(func(ctx Context, err error) {
		hasErr.SetTrue()
	}).Start()
	el.Close()
	time.Sleep(time.Second)
	assert.True(t, hasErr.IsTrue())
}

func TestNewEventListener_OnConnClose(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	el := NewEventListener(l)
	cnt := new(int64)
	cnt1 := new(int64)
	el.SetAutoCloseConnOnReadWriteError(true)
	el.OnAccept(func(ctx Context, c net.Conn) {
		atomic.StoreInt64(cnt, el.ConnCount())
	})
	el.SetOnConnClose(func(ctx Context) {
		atomic.StoreInt64(cnt1, el.ConnCount())
	})
	el.Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Millisecond * 300)
	assert.Equal(t, int64(1), atomic.LoadInt64(cnt))
	assert.Equal(t, int64(0), atomic.LoadInt64(cnt1))
}

type initErrorCodec struct {
	net.Conn
	called *gatomic.Bool
}

func newInitErrorCodec(called *gatomic.Bool) *initErrorCodec {
	return &initErrorCodec{called: called}
}
func (i *initErrorCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}
func (i *initErrorCodec) Initialize(ctx Context) error {
	i.called.SetTrue()
	return fmt.Errorf("init error")
}

type tempNetError struct {
}

func (t tempNetError) Error() string {
	return ""
}

func (t tempNetError) Timeout() bool {
	return true
}

func (t tempNetError) Temporary() bool {
	return true
}

func newTempNetError() *tempNetError {
	return &tempNetError{}
}

type initHijackedCodec struct {
	net.Conn
	called *bool
}

func (i *initHijackedCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}

func newInitHijackedCodec(called *bool) *initHijackedCodec {
	return &initHijackedCodec{
		called: called,
	}
}

func (i *initHijackedCodec) Initialize(ctx Context) error {
	// isHijacked the conn, do anything with conn, and just call ctx.Hijack at return.
	*i.called = true
	ctx.Hijack()
	return nil
}

type initHijackedFailCodec struct {
	net.Conn
	called    *bool
	modifyCtx func(ctx Context)
}

func (i *initHijackedFailCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}

func newInitHijackedFailCodec(called *bool, m func(ctx Context)) *initHijackedFailCodec {
	return &initHijackedFailCodec{
		called:    called,
		modifyCtx: m,
	}
}

func (i *initHijackedFailCodec) Initialize(ctx Context) error {
	// isHijacked the conn, do anything with conn, and just return ErrCodecHijacked at return.
	*i.called = true
	if i.modifyCtx != nil {
		i.modifyCtx(ctx)
	}
	ctx.Hijack()
	return fmt.Errorf("isHijacked fail")
}

type initPassThroughCodec struct {
	net.Conn
	called *gatomic.Bool
}

func newInitPassThroughCodec(called *gatomic.Bool) *initPassThroughCodec {
	return &initPassThroughCodec{
		called: called,
	}
}
func (i *initPassThroughCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}
func (i *initPassThroughCodec) Initialize(ctx Context) error {
	// pass through
	i.called.SetTrue()
	return nil
}

type initCodec struct {
	net.Conn
	called *gatomic.Bool
}

func newInitCodec2(called *gatomic.Bool) *initCodec {
	return &initCodec{
		called: called,
	}
}
func (i *initCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}
func (i *initCodec) Initialize(ctx Context) error {
	// do something
	i.called.SetTrue()
	ctx.Hijack()
	return nil
}

func TestNewEventListener_OnCodecError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	hasErr := gatomic.NewBool(false)
	el := NewEventListener(l)
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(gatomic.NewBool(false))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitErrorCodec(gatomic.NewBool(false))
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		_, e := Read(c, 5)
		if e != nil {
			hasErr.SetTrue()
		}
	}).Start()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c.Write([]byte("hello"))
	time.Sleep(time.Second)
	assert.True(t, hasErr.IsTrue())
}

func TestNewEventListener(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	var conn = gatomic.NewAny(nil)
	ctx := NewContext()
	ctx.SetData("cfg1", "abc")
	el := NewContextEventListener(ctx, l)
	el.AddConnFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		ctx.SetData("cfg2", "abc")
		c = NewBufferedConn(c)
		ctx.SetData("bufConn", c)
		return c, nil
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return NewAESCodec("abc")
	}).SetFirstReadTimeout(time.Second)
	el.OnAccept(func(ctx Context, c net.Conn) {
		c.Write([]byte("hello"))
		conn.SetVal(ctx.Data("bufConn").(net.Conn))
		assert.Equal(t, "abc", c.(*Conn).ctx.Data("cfg1"))
		assert.Equal(t, "abc", c.(*Conn).ctx.Data("cfg2"))
	}).Start()
	time.Sleep(time.Second)
	c, err := net.Dial("tcp", "127.0.0.1:"+p)
	assert.NoError(t, err)
	hasData := gatomic.NewBool(false)
	ec := NewEventConn(c)
	ec.AddCodec(NewAESCodec("abc"))
	ec.OnData(func(ctx Context, data []byte) {
		hasData.SetTrue()
		assert.Equal(t, "hello", string(data))
	}).Start()
	ec.Write([]byte("hello"))
	time.Sleep(time.Second)
	assert.True(t, hasData.IsTrue())
	assert.Implements(t, (*BufferedConn)(nil), conn.Val())
}

func TestEventListener_AutoCloseConn(t *testing.T) {
	l, _ := ListenEvent(":")
	l.SetAutoCloseConn(true)
	var hasErr = gatomic.NewBool(false)
	l.OnAccept(func(ctx Context, c net.Conn) {
		WriteTo(c, "hello")
		time.AfterFunc(time.Millisecond*300, func() {
			_, e := WriteTo(c, "1")
			if e != nil {
				hasErr.SetTrue()
			}
		})
	}).Start()
	c, _ := Dial(l.Addr().PortLocalAddr(), time.Second)
	d, _ := Read(c, 5)
	_, err := Read(c, 1)
	assert.Equal(t, "hello", d)
	assert.Error(t, err)
	time.Sleep(time.Second)
	assert.True(t, hasErr.IsTrue())
}
func TestProtocolListener_1(t *testing.T) {
	t.Parallel()
	l, err := NewListenerAddr(":0")
	assert.Nil(t, err)
	httpListener := l.NewProtocolListener(&ProtocolListenerOption{
		Name: "http",
		Checker: func(listener *Listener, conn BufferedConn) bool {
			h, err := conn.PeekMax(7)
			if err != nil {
				return false
			}
			return isHTTP(h)
		},
	})
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("okay"))
	})
	server := http.Server{
		Handler: mux,
	}
	go server.Serve(httpListener)
	time.Sleep(time.Second)
	addr := fmt.Sprintf("http://%s/http", NewAddr(httpListener.Addr()).PortLocalAddr())
	resp, err := http.Get(addr)
	assert.Nil(t, err)
	buf := make([]byte, 4)
	n, err := io.ReadFull(resp.Body, buf)
	assert.Nil(t, err)
	assert.Equal(t, "okay", string(buf[:n]))
	l.Close()
}

func TestProtocolListener_2(t *testing.T) {
	t.Parallel()
	l, err := NewListenerAddr(":0")
	assert.Nil(t, err)
	counter := new(int64)
	httpListener := l.NewProtocolListener(&ProtocolListenerOption{
		Name: "http",
		Checker: func(listener *Listener, conn BufferedConn) bool {
			h, err := conn.PeekMax(7)
			if err != nil {
				return false
			}
			return isHTTP(h)
		},
		OverflowAutoClose: true,
		OnQueueOverflow: func(l net.Listener, opt *ProtocolListenerOption, conn BufferedConn) {
			atomic.AddInt64(counter, 1)
			conn.Close()
		},
		ConnQueueSize: 9,
	})
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	time.Sleep(time.Second)
	addr := fmt.Sprintf("http://%s/http", NewAddr(httpListener.Addr()).PortLocalAddr())
	for i := 0; i < 10; i++ {
		go http.Get(addr)
	}
	time.Sleep(time.Second)
	l.Close()
	assert.Equal(t, atomic.LoadInt64(counter), int64(1))
}

func TestProtocolListener_3(t *testing.T) {
	t.Parallel()
	l, err := NewListenerAddr(":0")
	assert.Nil(t, err)
	httpListener := l.NewProtocolListener(&ProtocolListenerOption{
		Name: "http",
		Checker: func(listener *Listener, conn BufferedConn) bool {
			h, err := conn.PeekMax(7)
			if err != nil {
				return false
			}
			return isHTTP(h)
		},
		ConnQueueSize: 1,
	})
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	time.Sleep(time.Second)
	addr := fmt.Sprintf("http://%s/http", NewAddr(httpListener.Addr()).PortLocalAddr())
	for i := 0; i < 2; i++ {
		go http.Get(addr)
	}
	time.Sleep(time.Second)
	httpListener.Close()
	http.Get(addr)
	time.Sleep(time.Second)
	l.Close()
}

func isHTTP(head []byte) bool {
	if len(head) < 3 {
		return false
	}
	keys := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	for _, key := range keys {
		if bytes.HasPrefix(head, []byte(key)) || bytes.HasPrefix(head, []byte(strings.ToLower(key))) {
			return true
		}
	}
	return false
}

func TestNewListenerAddr_Error(t *testing.T) {
	t.Parallel()
	_, err := NewListenerAddr("invalid:address:format")
	assert.Error(t, err)
}

func TestNewEventListenerAddr_Error(t *testing.T) {
	t.Parallel()
	_, err := NewEventListenerAddr("invalid:address:format")
	assert.Error(t, err)
}

func TestListener_Close(t *testing.T) {
	t.Parallel()
	l, _ := NewListenerAddr(":0")
	err := l.Close()
	assert.NoError(t, err)
}

func TestProtocolListener_4(t *testing.T) {
	t.Parallel()
	l, err := NewListenerAddr(":0")
	assert.Nil(t, err)
	httpListener := l.NewProtocolListener(&ProtocolListenerOption{
		Name: "http",
		Checker: func(listener *Listener, conn BufferedConn) bool {
			h, err := conn.PeekMax(7)
			if err != nil {
				return false
			}
			return isHTTP(h)
		},
	})
	jsonLister := l.NewProtocolListener(&ProtocolListenerOption{
		Name: "json",
		Checker: func(listener *Listener, conn BufferedConn) bool {
			h, err := conn.PeekMax(1)
			if err != nil {
				return false
			}
			return string(h) == "{"
		},
	})
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("okay"))
	})
	server := http.Server{
		Handler: mux,
	}
	go server.Serve(httpListener)
	time.Sleep(time.Second)
	addr := fmt.Sprintf("http://%s/http", NewAddr(httpListener.Addr()).PortLocalAddr())
	resp, err := http.Get(addr)
	assert.Nil(t, err)
	buf := make([]byte, 4)
	n, err := io.ReadFull(resp.Body, buf)
	assert.Nil(t, err)
	assert.Equal(t, "okay", string(buf[:n]))
	okay := gatomic.NewBool(false)
	go func() {
		jsonLister.Accept()
		okay.SetTrue()
	}()
	time.Sleep(time.Second)
	c, _ := Dial(NewAddr(jsonLister.Addr()).PortLocalAddr(), time.Second)
	c.Write([]byte("{}"))
	time.Sleep(time.Second)
	assert.True(t, okay.IsTrue())
	l.Close()
}
