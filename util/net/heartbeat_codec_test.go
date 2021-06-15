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

func TestHeartbeatCodec(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	outputCnt := new(int32)
	errCnt := new(int32)
	go func() {
		c, _ := l.Accept()
		codec := NewHeartbeatCodec()
		conn := NewConn(c).AddCodec(codec)
		conn.SetTimeout(time.Second * 10)
		conn.Initialize()
		assert.Equal(t, time.Second*5, codec.Interval())
		assert.Equal(t, time.Second*5, codec.Timeout())
		go func() {
			for {
				_, err := conn.Write([]byte("hello from server"))
				if err != nil {
					atomic.AddInt32(errCnt, 1)
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
					atomic.AddInt32(errCnt, 1)
					fmt.Println("server read error", err)
					return
				}
				assert.Equal(t, "hello from client", string(buf[:n]))
				atomic.AddInt32(outputCnt, 1)
				//fmt.Printf("%s %d\n", string(buf[:n]), n)
			}
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
	conn.Initialize()
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				atomic.AddInt32(errCnt, 1)
				fmt.Println("client read error", err)
				return
			}
			atomic.AddInt32(outputCnt, 1)
			assert.Equal(t, "hello from server", string(buf[:n]))
		}
	}()
	go func() {
		for {
			_, err := conn.Write([]byte("hello from client"))
			if err != nil {
				atomic.AddInt32(errCnt, 1)
				fmt.Println("client write error", err)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	time.Sleep(time.Second * 15)
	assert.Equal(t, *errCnt, int32(4))
	t.Log("outputCnt:", *outputCnt)
	assert.True(t, *outputCnt > 250)
}

func TestHeartbeatCodec_UnknownMsg(t *testing.T) {
	t.Parallel()
	l, p, err := ListenRandom("")
	assert.NoError(t, err)
	el := NewEventListener(l)
	el.AddCodecFactory(func() Codec {
		return NewHeartbeatCodec()
	})
	el.OnAccept(func(l *EventListener, ctx Context, c net.Conn) {
		buf := make([]byte, 1024)
		_, err = c.Read(buf)
	})
	el.Start()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", ":"+p)
	c.Write([]byte("hello"))
	time.Sleep(time.Millisecond * 200)
	assert.Contains(t, err.Error(), "unrecognized msg type:")
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
	l, p, err := ListenRandom("")
	assert.NoError(t, err)
	el := NewEventListener(l)
	el.AddCodecFactory(func() Codec {
		return NewHeartbeatCodec()
	})
	el.OnAccept(func(l *EventListener, ctx Context, c net.Conn) {
		buf := make([]byte, 1024)
		_, err = c.Read(buf)
	})
	el.Start()
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", ":"+p)
	c = newTempErrorConn(c)
	c1 := NewConn(c)
	err = c1.AddCodec(NewHeartbeatCodec()).Initialize()
	assert.NoError(t, err)
	time.Sleep(time.Second * 5)
}
