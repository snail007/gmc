// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"crypto/cipher"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	gtest "github.com/snail007/gmc/util/testing"
	"github.com/stretchr/testify/assert"
)

func TestAESCodec(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	outputCnt := new(int32)
	password := "abc"
	debug := false
	go func() {
		c, _ := l.Accept()
		conn := NewConn(c).AddCodec(NewAESCodec(password))

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
	conn := NewConn(c).AddCodec(NewAESCodec(password))

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
	assert.True(t, *outputCnt > 50)
}

func TestAESCodec128(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	outputCnt := new(int32)
	options := &AESOptions{
		Type:     "aes-128",
		Password: "abc",
	}
	debug := false
	go func() {
		c, _ := l.Accept()
		conn := NewConn(c).AddCodec(NewAESCodecFromOptions(options))

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
	conn := NewConn(c).AddCodec(NewAESCodecFromOptions(options))

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
	assert.True(t, *outputCnt > 50)
}

func TestAESCodec192(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	outputCnt := new(int32)
	options := &AESOptions{
		Type:     "aes-192",
		Password: "abc",
	}
	debug := false
	go func() {
		c, _ := l.Accept()
		conn := NewConn(c).AddCodec(NewAESCodecFromOptions(options))

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
	conn := NewConn(c).AddCodec(NewAESCodecFromOptions(options))

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
	assert.True(t, *outputCnt > 50)
}

func TestAESCodec256(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	outputCnt := new(int32)
	options := &AESOptions{
		Type:     "aes-256",
		Password: "abc",
	}
	debug := false
	go func() {
		c, _ := l.Accept()
		conn := NewConn(c).AddCodec(NewAESCodecFromOptions(options))

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
	conn := NewConn(c).AddCodec(NewAESCodecFromOptions(options))

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
	assert.True(t, *outputCnt > 50)
}

func TestNewAESCodec_Error(t *testing.T) {
	//gtest.DebugRunProcess(t)
	if gtest.RunProcess(t, func() {
		aesNewCipher = func(key []byte) (cipher.Block, error) {
			return nil, fmt.Errorf("new_aes_error")
		}
		c, _ := net.Dial("tcp", ":")
		c0 := NewConn(c)
		c0.AddCodec(NewAESCodec(""))
		t.Log(c0.doInitialize().Error())
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	assert.True(t, strings.Contains(out, "new_aes_error"))
}
