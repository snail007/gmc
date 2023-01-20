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

	"github.com/stretchr/testify/assert"
)

func TestNewEventListener_OnFistReadTimeout(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	timeout := false
	el.SetFirstReadTimeout(time.Millisecond * 100).
		OnFistReadTimeout(func(ctx Context, c net.Conn, err error) {
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
		OnFistReadTimeout(func(ctx Context, c net.Conn, err error) {
			timeout = true
		}).Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.True(t, timeout)
}

func TestNewEventListener_Hijacked(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	el := NewEventListener(l)
	called := false
	isHijacked := false
	el.AddListenerFilter(func(ctx Context, c net.Conn) (net.Conn, error) {
		// isHijacked the conn, do anything with c
		isHijacked = true
		ctx.Hijack()
		return nil, nil
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		called = true
	})
	el.Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Second)
	assert.False(t, called)
	assert.True(t, isHijacked)
}

func TestNewEventListener_FilterError(t *testing.T) {
	t.Parallel()
	el, err := NewEventListenerAddr(":0")
	assert.Nil(t, err)
	p := NewAddr(el.Addr()).Port()
	var hasErr error
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
		_, hasErr = Read(c, 5)
	}).Start()
	_, err = net.Dial("tcp", "127.0.0.1:"+p)
	assert.NoError(t, err)
	time.Sleep(time.Second)
	assert.Equal(t, "error", hasErr.Error())
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
	el.OnAcceptError(func(ctx Context, err error) {
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
	var acceptCnt int
	l1.OnAccept(func(ctx Context, c net.Conn) {
		acceptCnt++
	})
	l1.Start()
	time.Sleep(time.Second * 1)
	net.Dial("tcp", "127.0.0.1:"+p)
	for i := 0; i <= 13; i++ {
		net.Dial("tcp", "127.0.0.1:"+p)
		time.Sleep(time.Millisecond * 300)
	}
	assert.Equal(t, 2, acceptCnt)
	assert.Equal(t, l1, l1.Ctx().EventListener())
}

func TestNewEventListener_OnAcceptError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	el := NewEventListener(l)
	hasErr := false
	el.OnAcceptError(func(ctx Context, err error) {
		hasErr = true
	}).Start()
	el.Close()
	time.Sleep(time.Second)
	assert.True(t, hasErr)
}

func TestNewEventListener_OnConnClose(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	el := NewEventListener(l)
	cnt := int64(0)
	cnt1 := int64(0)
	el.SetAutoCloseConnOnReadWriteError(true)
	el.OnAccept(func(ctx Context, c net.Conn) {
		cnt = el.ConnCount()
	})
	el.SetOnConnClose(func(ctx Context) {
		cnt1 = el.ConnCount()
	})
	el.Start()
	net.Dial("tcp", "127.0.0.1:"+p)
	time.Sleep(time.Millisecond * 300)
	assert.Equal(t, int64(1), cnt)
	assert.Equal(t, int64(0), cnt1)
}

type initErrorCodec struct {
	net.Conn
	called *bool
}

func newInitErrorCodec(called *bool) *initErrorCodec {
	return &initErrorCodec{called: called}
}
func (i *initErrorCodec) SetConn(conn net.Conn) Codec {
	i.Conn = conn
	return i
}
func (i *initErrorCodec) Initialize(ctx Context) error {
	*i.called = true
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
	called *bool
}

func newInitPassThroughCodec(called *bool) *initPassThroughCodec {
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
	*i.called = true
	return nil
}

type initCodec struct {
	net.Conn
	called *bool
}

func newInitCodec2(called *bool) *initCodec {
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
	*i.called = true
	ctx.Hijack()
	return nil
}

func TestNewEventListener_OnCodecError(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	var hasErr error
	el := NewEventListener(l)
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitPassThroughCodec(new(bool))
	})
	el.AddCodecFactory(func(ctx Context) Codec {
		return newInitErrorCodec(new(bool))
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		_, hasErr = Read(c, 5)
	}).Start()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c.Write([]byte("hello"))
	time.Sleep(time.Second)
	assert.Error(t, hasErr)
}

func TestNewEventListener(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	var conn net.Conn
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
		conn = ctx.Data("bufConn").(net.Conn)
		assert.Equal(t, "abc", c.(*Conn).ctx.Data("cfg1"))
		assert.Equal(t, "abc", c.(*Conn).ctx.Data("cfg2"))
	}).Start()
	time.Sleep(time.Second)
	c, err := net.Dial("tcp", "127.0.0.1:"+p)
	assert.NoError(t, err)
	hasData := false
	ec := NewEventConn(c)
	ec.AddCodec(NewAESCodec("abc"))
	ec.OnData(func(ctx Context, data []byte) {
		hasData = true
		assert.Equal(t, "hello", string(data))
	}).Start()
	ec.Write([]byte("hello"))
	time.Sleep(time.Second)
	assert.True(t, hasData)
	assert.Implements(t, (*BufferedConn)(nil), conn)
}

func TestEventListener_AutoCloseConn(t *testing.T) {
	l, _ := ListenEvent(":")
	l.SetAutoCloseConn(true)
	var hasErr error
	l.OnAccept(func(ctx Context, c net.Conn) {
		WriteTo(c, "hello")
		time.AfterFunc(time.Millisecond*300, func() {
			_, hasErr = WriteTo(c, "1")
		})
	}).Start()
	c, _ := Dial(l.Addr().PortLocalAddr(), time.Second)
	d, _ := Read(c, 5)
	_, err := Read(c, 1)
	assert.Equal(t, "hello", d)
	assert.Error(t, err)
	time.Sleep(time.Second)
	assert.Error(t, hasErr)
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
	d, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "okay", string(d))
	okay := false
	go func() {
		jsonLister.Accept()
		okay = true
	}()
	time.Sleep(time.Second)
	c, _ := Dial(NewAddr(jsonLister.Addr()).PortLocalAddr(), time.Second)
	c.Write([]byte("{}"))
	time.Sleep(time.Second)
	assert.True(t, okay)
	l.Close()
}
