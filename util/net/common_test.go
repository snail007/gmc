// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	gtest "github.com/snail007/gmc/util/testing"
	"github.com/stretchr/testify/assert"
)

func TestWriteBytes(t *testing.T) {
	err := WriteBytes("", []byte(""))
	assert.Error(t, err)
}

func TestListenRandom(t *testing.T) {
	if gtest.RunProcess(t, func() {
		netListen = func(network, address string) (net.Listener, error) {
			return nil, fmt.Errorf("listen_error")
		}
		_, _, err := RandomListen()
		t.Log(err)
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	assert.True(t, strings.Contains(out, "listen_error"))
}

func TestListen(t *testing.T) {
	l, err := Listen(":0")
	assert.NoError(t, err)
	assert.NotNil(t, l)
	l.Close()
}

func TestListenError(t *testing.T) {
	if gtest.RunProcess(t, func() {
		netListen = func(network, address string) (net.Listener, error) {
			return nil, fmt.Errorf("listen_error")
		}
		_, err := Listen("127.0.0.1:0")
		t.Log(err.Error())
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	assert.True(t, strings.Contains(out, "listen_error"))
}

func TestDial(t *testing.T) {
	l, port, _ := RandomListen()
	NewEventListener(l).Start()
	c, _ := Dial("127.0.0.1:"+port, time.Second)
	assert.NotNil(t, c)
}

func TestDialError(t *testing.T) {
	if gtest.RunProcess(t, func() {
		netDialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
			return nil, fmt.Errorf("dial_error")
		}
		_, err := Dial("", time.Millisecond)
		t.Log(err.Error())
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	assert.True(t, strings.Contains(out, "dial_error"))
}

func TestRandomPort(t *testing.T) {
	p, err := RandomPort()
	assert.NoError(t, err)
	assert.True(t, len(p) > 1)
}

func TestRandomPortError(t *testing.T) {
	if gtest.RunProcess(t, func() {
		netListen = func(network, address string) (net.Listener, error) {
			return nil, fmt.Errorf("listen_error")
		}
		_, err := RandomPort()
		t.Log(err.Error())
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	assert.True(t, strings.Contains(out, "listen_error"))
}

func TestNewTCPAddr(t *testing.T) {
	a := NewTCPAddr("")
	assert.Error(t, a.Err())
	a = NewTCPAddr(":")
	assert.NoError(t, a.Err())
	assert.Empty(t, a.Host())
	assert.Equal(t, "0", a.Port())
	assert.Equal(t, ":", a.String())
	assert.Equal(t, "tcp", a.Network())
	assert.Equal(t, "127.0.0.1:0", a.PortAddr("127.0.0.1"))
}

func TestWriteTo(t *testing.T) {
	l, _ := ListenEvent(":")
	l.OnAccept(func(ctx Context, c net.Conn) {
		WriteTo(c, "hello")
	})
	l.SetAutoCloseConn(true)
	l.Start()
	c, _ := Dial(l.Addr().PortLocalAddr(), time.Second)
	d, _ := Read(c, 5)
	assert.Equal(t, "hello", d)
}

func TestListenEventError(t *testing.T) {
	if gtest.RunProcess(t, func() {
		netListen = func(network, address string) (net.Listener, error) {
			return nil, fmt.Errorf("listen_error")
		}
		_, err := ListenEvent("127.0.0.1:0")
		t.Log(err.Error())
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	assert.True(t, strings.Contains(out, "listen_error"))
}

func TestLocalOutgoingIP(t *testing.T) {
	ip, err := LocalOutgoingIP()
	assert.Nil(t, err)
	assert.NotNil(t, net.ParseIP(ip))
}

func TestIsPrivateIP(t *testing.T) {
	testCases := []struct {
		ip       string
		expected bool
	}{
		{"192.168.1.1", true},
		{"2001:0db8::1", false},
		{"8.8.8.8", false},
		{"2001:4860:4860::8888", false},
		{"invalidIP", false},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"192.168.2.1", true},
		{"fc00::1", false},
		{"fe80::1", false},
		{"2001:0db8::2", false},
	}

	for _, testCase := range testCases {
		result := IsPrivateIP(testCase.ip)
		if result != testCase.expected {
			t.Errorf("For IP %s, expected %t, but got %t", testCase.ip, testCase.expected, result)
		}
	}
}
