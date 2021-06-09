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
		conn.Initialize()
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
	conn.Initialize()
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
}

func TestMultipleCodec2(t *testing.T) {
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
		conn.Initialize()
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
	conn.Initialize()
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
	assert.Equal(t, int64(510), conn.ReadBytes())
	assert.Equal(t, int64(510), conn.WriteBytes())
	assert.True(t, *outputCnt > 50)
}

func TestEventConn(t *testing.T) {
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
	conn.OnData(func(s *EventConn, data []byte) {
		assert.Equal(t, "hello", string(data))
		assert.Equal(t, int64(5), s.ReadBytes())
		s.Write([]byte("hi"))
		assert.Equal(t, int64(2), s.WriteBytes())
	}).OnClose(func(s *EventConn) {
		closed = true
	}).Start()
	time.AfterFunc(time.Millisecond*1500, func() {
		conn.Close()
	})
	time.Sleep(time.Second * 3)
	assert.Contains(t, conn.LocalAddr().String(), "127.0.0.1:")
	assert.Contains(t, conn.RemoteAddr().String(), "127.0.0.1:")
	assert.True(t, closed)
}

func TestEventConn2(t *testing.T) {
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
	conn := NewEventConn(c)
	conn.SetReadBufferSize(1024)
	conn.SetTimeout(time.Second)
	assert.Equal(t, time.Second, conn.ReadTimeout())
	assert.Equal(t, time.Second, conn.WriteTimeout())
	closed := false
	readErr := false
	writeErr := false
	conn.OnData(func(s *EventConn, data []byte) {
		assert.Equal(t, "hello", string(data))
		assert.Equal(t, int64(5), s.ReadBytes())
		time.Sleep(time.Millisecond * 300)
		s.Write([]byte("hi"))
		assert.Equal(t, int64(2), s.WriteBytes())
	}).OnReadError(func(s *EventConn, err error) {
		readErr = true
	}).OnWriterError(func(s *EventConn, err error) {
		writeErr = true
	}).OnClose(func(s *EventConn) {
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
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	go func() {
		c, _ := l.Accept()
		c.Write([]byte("hello"))
	}()
	hasErr:=false
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn := NewEventConn(c)
	conn.AddCodec(&initOkayCodec{})
	conn.AddCodec(&initErrorCodec{})
 	conn.OnCodecInitializeError(func(_ *EventConn, err error) {
		hasErr=true
	}).Start()
	time.Sleep(time.Second * 1)
 	assert.True(t, hasErr)
}