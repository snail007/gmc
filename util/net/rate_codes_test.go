// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"io"
	"io/ioutil"
	"net"
	"strings"
	"testing"
	"time"

	gtest "github.com/snail007/gmc/util/testing"
	"github.com/stretchr/testify/assert"
)

func TestRateCodec1(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen()
	l0 := NewEventListener(l)
	l0.SetAutoCloseConn(true)
	l0.AddCodecFactory(func(ctx Context) Codec {
		return NewRateCodec(1)
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		c.Write([]byte{0, 0, 0})
	}).Start()
	time.Sleep(time.Millisecond * 300)
	c, _ := Dial("127.0.0.1:"+p, time.Second)
	start := time.Now()
	io.Copy(ioutil.Discard, c)
	assert.True(t, time.Now().Sub(start).Seconds() > 2)
}

func TestRateCodec2(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen()
	l0 := NewEventListener(l)
	l0.SetAutoCloseConn(true)
	l0.OnAccept(func(ctx Context, c net.Conn) {
		c.Write(make([]byte, 3))
	}).Start()
	time.Sleep(time.Millisecond * 300)
	c, _ := Dial("127.0.0.1:"+p, time.Second)
	c.AddCodec(NewRateCodec(1))
	start := time.Now()
	io.Copy(ioutil.Discard, c)
	assert.True(t, time.Now().Sub(start).Seconds() > 2)
}

func TestRateCodec3(t *testing.T) {
	l, p, _ := RandomListen()
	l0 := NewEventListener(l)
	l0.SetAutoCloseConn(true)
	data := make([]byte, 1024*1024) //1MB
	l0.OnAccept(func(ctx Context, c net.Conn) {
		for i := 0; i < 100; i++ {
			c.Write(data)
		}
	}).Start()
	time.Sleep(time.Millisecond * 400)
	c, err := Dial("127.0.0.1:"+p, time.Second)
	assert.NoError(t, err)
	c.AddCodec(NewRateCodec(1024 * 1024 * 1024))
	start := time.Now()
	_, err = io.Copy(ioutil.Discard, c)
	assert.NoError(t, err)
	used := time.Now().Sub(start).Seconds()
	assert.True(t, used < 0.3)
}

func TestRateCodec4(t *testing.T) {
	if gtest.RunProcess(t, func() {
		t.Parallel()
		l, p, _ := RandomListen()
		l0 := NewEventListener(l)
		l0.SetAutoCloseConn(true)
		l0.OnAccept(func(ctx Context, c net.Conn) {
			c.Write(make([]byte, 10))
		}).Start()
		l0.Start()
		c, _ := Dial("127.0.0.1:"+p, time.Second)
		burstCnt = 3
		c.AddCodec(NewRateCodec(3))
		start := time.Now()
		io.Copy(ioutil.Discard, c)
		if time.Now().Sub(start).Seconds() > 2 {
			t.Log("burstCnt_1_okay")
		} else {
			t.Log("burstCnt_1_fail")
		}
	}) {
		return
	}
	t.Parallel()
	out, _, _ := gtest.NewProcess(t).Wait()
	assert.True(t, strings.Contains(out, "burstCnt_1_okay"))
}

func TestRateCodec_Close(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen()
	l0 := NewEventListener(l)
	l0.OnAccept(func(ctx Context, c net.Conn) {
		time.Sleep(time.Second * 2)
		c.Close()
	}).Start()
	
	time.Sleep(time.Millisecond * 300)
	c, _ := Dial("127.0.0.1:"+p, time.Second)
	c.AddCodec(NewRateCodec(1024))
	
	err := c.Close()
	assert.NoError(t, err)
}

func TestRateCodec_ReadWriteError(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen()
	done := make(chan bool)
	l0 := NewEventListener(l)
	l0.OnAccept(func(ctx Context, c net.Conn) {
		<-done
		c.Close()
	}).Start()
	
	time.Sleep(time.Millisecond * 300)
	c, _ := Dial("127.0.0.1:"+p, time.Second)
	c.AddCodec(NewRateCodec(1024))
	
	c.RawConn().Close()
	time.Sleep(time.Millisecond * 100)
	
	_, err := c.Write([]byte("test"))
	assert.Error(t, err)
	
	buf := make([]byte, 100)
	_, err = c.Read(buf)
	assert.Error(t, err)
	
	done <- true
}

func TestRateCodec_LargeData(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen()
	largeData := make([]byte, 1024*1024*5) // 5MB
	
	l0 := NewEventListener(l)
	l0.SetAutoCloseConn(true)
	l0.OnAccept(func(ctx Context, c net.Conn) {
		c.Write(largeData)
	}).Start()
	
	time.Sleep(time.Millisecond * 300)
	c, _ := Dial("127.0.0.1:"+p, time.Second)
	c.AddCodec(NewRateCodec(1024 * 1024 * 100)) // 100MB/s
	
	received, err := io.ReadAll(c)
	assert.NoError(t, err)
	assert.Equal(t, len(largeData), len(received))
}
