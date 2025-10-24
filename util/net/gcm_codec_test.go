// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGCMCodec_Basic(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	password := "test-password-123"
	g := sync.WaitGroup{}
	g.Add(4)
	
	go func() {
		c, _ := l.Accept()
		codec, err := NewGCMCodec(password)
		assert.NoError(t, err)
		conn := NewConn(c).AddCodec(codec)
		
		go func() {
			defer g.Done()
			conn.Write([]byte("hello from server"))
		}()
		go func() {
			defer g.Done()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			assert.Equal(t, "hello from client", string(buf[:n]))
		}()
	}()
	
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	codec, err := NewGCMCodec(password)
	assert.NoError(t, err)
	conn := NewConn(c).AddCodec(codec)
	
	go func() {
		defer g.Done()
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
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

func TestGCMCodec_InvalidPassword(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	
	go func() {
		c, _ := l.Accept()
		codec, _ := NewGCMCodec("password1")
		conn := NewConn(c).AddCodec(codec)
		conn.Write([]byte("hello"))
		time.Sleep(time.Second * 2)
		c.Close()
	}()
	
	time.Sleep(time.Millisecond * 500)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	codec, _ := NewGCMCodec("password2")
	conn := NewConn(c).AddCodec(codec)
	
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	assert.Error(t, err)
}

func TestGCMCodec_LargeData(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	password := "test-password"
	
	largeData := make([]byte, 1024*100) // 100KB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	
	received := make([]byte, len(largeData))
	done := make(chan bool)
	
	go func() {
		c, _ := l.Accept()
		codec, _ := NewGCMCodec(password)
		conn := NewConn(c).AddCodec(codec)
		
		total := 0
		for total < len(received) {
			n, err := conn.Read(received[total:])
			if err != nil {
				break
			}
			total += n
		}
		done <- true
	}()
	
	time.Sleep(time.Second)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	codec, _ := NewGCMCodec(password)
	conn := NewConn(c).AddCodec(codec)
	
	conn.Write(largeData)
	<-done
	
	assert.Equal(t, largeData, received)
}

func TestGCMCodec_EmptyData(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	password := "test-password"
	
	done := make(chan bool)
	
	go func() {
		c, _ := l.Accept()
		codec, _ := NewGCMCodec(password)
		conn := NewConn(c).AddCodec(codec)
		
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		assert.Equal(t, 0, n)
		done <- true
	}()
	
	time.Sleep(time.Millisecond * 500)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	codec, _ := NewGCMCodec(password)
	conn := NewConn(c).AddCodec(codec)
	
	conn.Write([]byte{})
	<-done
}
