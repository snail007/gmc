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

	gatomic "github.com/snail007/gmc/util/sync/atomic"

	"github.com/stretchr/testify/assert"
)

func TestHeartbeatCodec(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	g := sync.WaitGroup{}
	g.Add(4)
	go func() {
		c, _ := l.Accept()
		codec := NewHeartbeatCodec()
		conn := NewConn(c).AddCodec(codec)
		conn.SetTimeout(time.Second * 10)
		go func() {
			defer g.Done()
			_, err := conn.Write([]byte("hello from server"))
			assert.Equal(t, time.Second*5, codec.Interval())
			assert.Equal(t, time.Second*5, codec.Timeout())
			if err != nil {
				fmt.Println("server write error", err)
				return
			}
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
		time.AfterFunc(time.Second*13, func() {
			conn.Close()
		})
		select {}
	}()
	time.Sleep(time.Second * 2)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	codec := NewHeartbeatCodec().SetTimeout(time.Second).SetInterval(time.Millisecond * 100)
	assert.Equal(t, time.Millisecond*100, codec.Interval())
	assert.Equal(t, time.Second, codec.Timeout())
	conn := NewConn(c).AddCodec(codec)
	conn.SetTimeout(time.Second * 10)
	assert.Equal(t, time.Second*10, conn.ReadTimeout())
	assert.Equal(t, time.Second*10, conn.WriteTimeout())
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

func TestHeartbeatCodec_UnknownMsg(t *testing.T) {
	t.Parallel()
	l, p, err := RandomListen("")
	assert.NoError(t, err)
	el := NewEventListener(l)
	el.AddCodecFactory(func(ctx Context) Codec {
		return NewHeartbeatCodec()
	})
	str := gatomic.NewString("")
	el.OnAccept(func(ctx Context, c net.Conn) {
		buf := make([]byte, 1024)
		_, err := c.Read(buf)
		if err != nil {
			str.SetVal(err.Error())
		}
	})
	el.Start()
	time.Sleep(time.Second * 3)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c.Write([]byte("hello"))
	time.Sleep(time.Millisecond * 200)
	assert.Contains(t, str.Val(), "unrecognized heartbeat msg type:")
}

type tempErrorConn struct {
	net.Conn
	writeCnt int
	readCnt  int
}

func newTempErrorConn(conn net.Conn) *tempErrorConn {
	return &tempErrorConn{Conn: conn}
}

func (s *tempErrorConn) Read(b []byte) (n int, err error) {
	defer func() { s.readCnt++ }()
	return 0, newTempNetError()
}

func (s *tempErrorConn) Write(b []byte) (n int, err error) {
	defer func() { s.writeCnt++ }()
	return 0, newTempNetError()
}

func TestHeartbeatCodec_TempError(t *testing.T) {
	t.Parallel()
	l, p, err := RandomListen("")
	assert.NoError(t, err)
	el := NewEventListener(l)
	el.AddCodecFactory(func(ctx Context) Codec {
		return NewHeartbeatCodec()
	})
	el.OnAccept(func(ctx Context, c net.Conn) {
		buf := make([]byte, 1024)
		_, err = c.Read(buf)
	})
	el.Start()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	c = newTempErrorConn(c)
	c1 := NewConn(c)
	err = c1.AddCodec(NewHeartbeatCodec()).doInitialize()
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)
}
