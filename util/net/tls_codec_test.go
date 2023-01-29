// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"crypto/x509"
	"fmt"
	gatomic "github.com/snail007/gmc/util/sync/atomic"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	gtest "github.com/snail007/gmc/util/testing"
	"github.com/stretchr/testify/assert"
)

func TestTLSCodec_1(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		// server certificates
		c0.AddCertificate(testCert, testKEY)
		c0.AddCertificate(demoCert, demoKEY)
		// client ca
		c0.RequireClientAuth(true)
		c0.AddClientCa(demoCert)
		c0.AddClientCa(testCert)
		c0.AddClientCa(helloCert)
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, err := c.Write([]byte("hello"))
		assert.NoError(t, err)
	})
	l0.OnAcceptError(func(ctx Context, err error) {
		t.Log(err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)

	c1 := NewTLSClientCodec()
	c1.AddCertificate(helloCert, helloKEY)
	c1.SetServerName("demo.com")
	c1.AddServerCa(demoCert)
	c1.SkipVerify(false)
	conn0.AddCodec(c1)
	d, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "hello", d)
}

func TestTLSCodec_2(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		// server certificates
		c0.AddCertificate(testCert, testKEY)
		// client ca
		c0.RequireClientAuth(true)
		c0.AddClientCa(demoCert)
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, err := c.Write([]byte("hello"))
		assert.NoError(t, err)
	})
	l0.OnAcceptError(func(ctx Context, err error) {
		t.Log(err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)

	c1 := NewTLSClientCodec()
	c1.AddCertificate(demoCert, demoKEY)
	c1.AddCertificate(testCert, testKEY)
	c1.SetServerName("demo.com")
	c1.AddServerCa(helloCert)
	c1.AddServerCa(testCert)
	c1.SkipVerify(false)
	c1.SkipVerifyCommonName(true)
	conn0.AddCodec(c1)
	d, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "hello", d)
}

func TestTLSCodec_3(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		// server certificates
		c0.AddCertificate(testCert, testKEY)
		c0.AddCertificate(demoCert, demoKEY)
		// client ca
		c0.RequireClientAuth(true)
		c0.AddClientCa(demoCert)
		c0.AddClientCa(testCert)
		c0.AddClientCa(helloCert)
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, err := c.Write([]byte("hello"))
		assert.NoError(t, err)
	})
	l0.OnAcceptError(func(ctx Context, err error) {
		t.Log(err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)

	c1 := NewTLSClientCodec()
	c1.AddCertificate(helloCert, helloKEY)
	c1.SetServerName("test.com")
	c1.AddServerCa(testCert)
	c1.SkipVerify(false)
	conn0.AddCodec(c1)

	d, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "hello", d)
}

func TestTLSCodec_4(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		c0.RequireClientAuth(true)
		// server certificates
		c0.AddCertificate(testCert, testKEY)
		c0.AddCertificate(demoCert, demoKEY)
		// client ca
		c0.AddClientCa(demoCert)
		c0.LoadSystemCas()
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, err := c.Write([]byte("hello"))
		assert.Error(t, err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)
	c1 := NewTLSClientCodec()
	c1.AddCertificate(helloCert, helloKEY)
	c1.AddCertificate(testCert, testKEY)
	c1.SetServerName("demo.com")
	c1.AddServerCa(demoCert)
	c1.AddServerCa(testCert)
	c1.SkipVerify(false)
	conn0.AddCodec(c1)

	_, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.Error(t, err)
}

func TestTLSCodec_5(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		// server certificates
		c0.AddCertificate(testCert, testKEY)
		c0.AddCertificate(demoCert, demoKEY)
		// client ca
		c0.RequireClientAuth(true)
		c0.AddClientCa(demoCert)
		c0.AddClientCa(testCert)
		c0.AddClientCa(helloCert)
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, err := c.Write([]byte("hello"))
		assert.NoError(t, err)
	})
	l0.OnAcceptError(func(ctx Context, err error) {
		t.Log(err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)

	c1 := NewTLSClientCodec()
	assert.Error(t, c1.PinServerCert([]byte("aaa")))
	c1.PinServerCert(testCert)
	c1.AddCertificate(helloCert, helloKEY)
	c1.SetServerName("test.com")
	c1.AddServerCa(testCert)
	c1.SkipVerify(false)
	conn0.AddCodec(c1)

	d, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "hello", d)
}

func TestTLSCodec_6(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		// server certificates
		c0.AddCertificate(testCert, testKEY)
		c0.AddCertificate(demoCert, demoKEY)
		// client ca
		c0.RequireClientAuth(true)
		c0.AddClientCa(demoCert)
		c0.AddClientCa(testCert)
		c0.AddClientCa(helloCert)
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, err := c.Write([]byte("hello"))
		assert.NoError(t, err)
	})
	l0.OnAcceptError(func(ctx Context, err error) {
		t.Log(err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)

	c1 := NewTLSClientCodec()
	c1.AddCertificate(helloCert, helloKEY)
	c1.SetServerName("test.com")
	c1.SkipVerify(true)
	conn0.AddCodec(c1)

	d, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "hello", d)
}

func TestTLSCodec_7(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		c0.RequireClientAuth(false)
		// server certificates
		c0.AddCertificate(demoCert, demoKEY)
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, err := c.Write([]byte("hello"))
		assert.NoError(t, err)
	})
	l0.OnAcceptError(func(ctx Context, err error) {
		t.Log(err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)

	c1 := NewTLSClientCodec()
	c1.SkipVerify(true)
	c1.SkipVerifyCommonName(true)
	conn0.AddCodec(c1)

	d, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "hello", d)
}

func TestTLSCodec_8(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		c0.RequireClientAuth(false)
		// server certificates
		c0.AddCertificate(demoCert, demoKEY)
		return c0
	})
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, err := c.Write([]byte("hello"))
		assert.Error(t, err)
	})
	l0.OnAcceptError(func(ctx Context, err error) {
		t.Log(err)
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)

	c1 := NewTLSClientCodec()
	c1.SkipVerify(false)
	c1.AddServerCa(testCert)
	conn0.AddCodec(c1)

	_, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.Error(t, err)

	c2 := NewTLSClientCodec()
	c2.SkipVerify(false)
	c2.AddServerCa([]byte("aaa"))
	conn0.AddCodec(c2)
	assert.Error(t, conn0.initialize())
}

func TestTLSCodec_9(t *testing.T) {
	t.Parallel()
	l, p, _ := RandomListen("")
	l0 := NewEventListener(l)
	l0.AddCodecFactory(func(ctx Context) Codec {
		c0 := NewTLSServerCodec()
		c0.AddCertificate(testCert, testKEY)
		c0.AddClientCa([]byte("aaa"))
		return c0
	})
	var hasErr = gatomic.NewBool()
	l0.OnAccept(func(ctx Context, c net.Conn) {
		_, e := Read(c, 5)
		if e != nil {
			hasErr.SetTrue()
		}
		c.Close()
	})
	l0.Start()
	time.Sleep(time.Second)
	conn, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn0 := NewConn(conn)

	c1 := NewTLSClientCodec()
	c1.AddCertificate(helloCert, helloKEY)
	c1.SetServerName("demo.com")
	c1.AddServerCa(demoCert)
	c1.SkipVerify(true)
	conn0.AddCodec(c1)
	assert.Error(t, c1.AddCertificate([]byte("aaa"), []byte("aaa")))

	_, err := Read(conn0, 5)
	time.Sleep(time.Second)
	assert.Error(t, err)
	assert.True(t, hasErr.IsTrue())
}

func TestTLSCodec_10(t *testing.T) {
	if gtest.RunProcess(t, func() {
		x509SystemCertPool = func() (*x509.CertPool, error) {
			return nil, fmt.Errorf("x509SystemCertPool_Error")
		}
		codec1 := NewTLSClientCodec()
		codec1.LoadSystemCas()
		l, _ := net.Listen("tcp", ":0")
		_, p, _ := net.SplitHostPort(l.Addr().String())
		c, _ := Dial("127.0.0.1:"+p, time.Second)
		c.AddCodec(codec1)
		_, err := WriteTo(c, "")
		t.Log(err)
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	assert.True(t, strings.Contains(out, "x509SystemCertPool_Error"))
}

func Test_SetHttpClientTLSCodec(t *testing.T) {
	l, _ := Listen(":")
	l.AddCodecFactory(func(ctx Context) Codec {
		cc := NewTLSServerCodec()
		cc.AddCertificate(testCert, testKEY)
		return cc
	})
	var isTLS interface{}
	go func() {
		http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h, _ := w.(http.Hijacker)
			c, _, _ := h.Hijack()
			isTLS = c.(*Conn).Ctx().IsTLS()
			c.Write([]byte("HTTP/1.0 200 OK \r\n\r\nhello"))
			c.Close()
		}))
	}()
	client := &http.Client{}
	NewTLSClientCodec().AddServerCa(testCert).AddToHTTPClient(client)
	resp, _ := client.Get("https://" + NewAddr(l.Addr()).PortLocalAddr())
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "hello", string(b))
	assert.Equal(t, true, isTLS)
}

func Test_SetHttpClientTLSCodec1(t *testing.T) {
	l, _ := Listen(":")
	l.AddCodecFactory(func(ctx Context) Codec {
		cc := NewTLSServerCodec()
		cc.AddCertificate(testCert, testKEY)
		return cc
	})
	var isTLS interface{}
	go func() {
		http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h, _ := w.(http.Hijacker)
			c, _, _ := h.Hijack()
			isTLS = c.(*Conn).Ctx().IsTLS()
			c.Write([]byte("HTTP/1.0 200 OK \r\n\r\nhello"))
			c.Close()
		}))
	}()
	called := false
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				called = true
				return net.Dial(network, addr)
			},
		},
	}
	NewTLSClientCodec().AddServerCa(testCert).AddToHTTPClient(client)
	resp, _ := client.Get("https://" + NewAddr(l.Addr()).PortLocalAddr())
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "hello", string(b))
	assert.Equal(t, true, isTLS)
	assert.True(t, called)
}

func Test_SetHttpClientTLSCodecError(t *testing.T) {
	//gtest.DebugRunProcess(t)
	if gtest.RunProcess(t, func() {
		l, p, _ := RandomListen()
		l0 := NewListener(l)
		l0.AddCodecFactory(func(ctx Context) Codec {
			c := NewTLSClientCodec()
			c.AddCertificate(testCert, testKEY)
			return c
		})
		server := http.Server{Addr: ":0"}
		go server.Serve(l0)
		time.Sleep(time.Millisecond * 300)
		client := &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return nil, fmt.Errorf("dial_error")
				},
			},
			Timeout: time.Millisecond * 100,
		}
		NewTLSClientCodec().AddServerCa(testCert).AddToHTTPClient(client)
		_, err := client.Get("https://127.0.0.1:" + p)
		t.Log(err.Error())
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	t.Log(out)
	assert.True(t, strings.Contains(out, "dial_error"))
}

func Test_SetHttpClientTLSCodecError1(t *testing.T) {
	if gtest.RunProcess(t, func() {
		netDialTimeout = func(network, address string, timeout time.Duration) (net.Conn, error) {
			return nil, fmt.Errorf("dial_error")
		}
		l, p, _ := RandomListen()
		l0 := NewListener(l)
		l0.AddCodecFactory(func(ctx Context) Codec {
			c := NewTLSClientCodec()
			c.AddCertificate(testCert, testKEY)
			return c
		})
		server := http.Server{Addr: ":0"}
		go server.Serve(l0)
		time.Sleep(time.Millisecond * 300)
		client := &http.Client{
			Timeout: time.Millisecond * 100,
		}
		NewTLSClientCodec().AddServerCa(testCert).AddToHTTPClient(client)
		_, err := client.Get("https://127.0.0.1:" + p)
		t.Log(err.Error())
	}) {
		return
	}
	out, _, _ := gtest.NewProcess(t).Wait()
	t.Log(out)
	assert.True(t, strings.Contains(out, "dial_error"))
}

var testKEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyfZmbS0d60a4RpNFtg9G6gGYvjp7IOLqbnNi8N2dAFuZGShP
gn+X/+8PYaOnai1G1khD7LBF9DqjS2cqJYEQSouK8LVu1mztXMP4X7V2VMIRqVp9
uduVEpARFbap64D7Y2Wv9nj8dX1mcSzGzW6fPJT1STThE3ByqCpSedOxo3CCE/hd
bI2x+GNS3UkPSsibMrvanZCEKiv8a1aJE03PJzvnblmM8pVS3wvr76w5XvQzGir1
idhXU4OvG4ccxZpS/pY3YSgmqwghjeYMdUd0b0jaKpqjHE+PoQ+nwhJ20K+7pVCQ
OAMjQxypelv+aTVD0At3oeeEOMY4qXianM+6BQIDAQABAoIBAA8tWiMoOdBdfymm
lZ2J5l1dg1oAURJ2mwFz4GKTdOH7ADVYxyjaZ9TO5UwEHWeoQWOHCLu3v3oMEgtv
lEY/PbcsZ2ORbuPkSa4n9/lRTLQv5V3htAMMklZTx0TndjuBdOLSWHfgPbCinNky
cTos7wCBfTFkLOnmEGe8znfjRb0va+rH4AmjjGR3olnY/T/pkdgdZeIu3y8m/rtQ
IjWIH5GW4Kpqu8YlA+cAC56xEBRusBndabjRgkSVUUf7VaaDPVKFzR3Lv8s2dcWc
16YpZ47D366nd+zIC3MAtlRvRXWZdM9JCfHbVQBmsuFbiQxciJ7epx+S1ZwDf4FY
9NHqUb0CgYEA/FFuioMbzjr1urECg7czLfDIivja7bTUbHL481Ck0oTUqaO013w5
/u6YWCyLjyKmQjhNwMUDZ32C2LVX4eJEM72Dk8si4GHlQyBDlmmEFbZ5HIS2G8nU
Vu7TytbIdw0oC8q9+7dAlk85t4ng74B+tfXXUnDIrT8y/DxZkSf2MVMCgYEAzOjb
s4qae++I/0aPQn+HZ6kDehSiY1xvOBvFTzB51gdFzh+2ocSze6aknSHZ4WmJ8yQB
7N5yHnT6jhRwvlpFjJRxGtkvc4/eDG0v1RCK4fY6o8ljEa3c+7i5EJRLuExlnnQ7
gIxWSMTplXbuOWG+WhyRSxLEOsioGtr7eWAaREcCgYAgsMA8q+3vU05BCOwFerfj
zN1+u+1JfPNEtcSxaZJhQBp5fB9TB+JPuEP+sI7IVbnqvHa+cggV4XoRb7VaK8Gg
Xn5sqJX1MlnMz6JSG4ukcIbSfhNGGGktdjX0gs1oN0kn9fWVZlG058DXmcKN5T0F
gDuMj9ZAM/78FSmZl+7axwKBgAeTaGQH8NQ6M+90NWG5A1GSzx0ZXDOePEJvzGi0
Gx0NocgQJhlvA0/EBnwEv2B1HXOO1j9irgdwPb85BD4ValLbPh9G/lkgbY46DzWq
aegWyW46yN3jdrMbzkPNp8sFkBA+reB/z8Ta+uPaxM38TiRYwApthDHEL2rmw7tm
ETKLAoGBAI7k+FPWXbn5XIOk2YKanJXG8a431IySiBFIQbFtVjHA4oXx0zVxWx9Y
YENKZH2R+IaxlrnE4IQReB9HXQxyqUWr6Jl9bm/dpOUT4plI9kLz2V4jcnVfezo+
4DKAVBRUIVKlY1BFnWeepY4R8Pg2ClpAL2p9nQLGYYLeN0oz6I3f
-----END RSA PRIVATE KEY-----
`)

// dns name: test.com
var testCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDaDCCAlCgAwIBAgIBATANBgkqhkiG9w0BAQsFADBGMQswCQYDVQQGEwJNVzER
MA8GA1UEChMIdGVzdC5jb20xETAPBgNVBAsTCHRlc3QuY29tMREwDwYDVQQDEwh0
ZXN0LmNvbTAeFw0yMTA2MTgxMjA3NTdaFw0zMTA2MTYxMzA3NTdaMEYxCzAJBgNV
BAYTAk1XMREwDwYDVQQKEwh0ZXN0LmNvbTERMA8GA1UECxMIdGVzdC5jb20xETAP
BgNVBAMTCHRlc3QuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
yfZmbS0d60a4RpNFtg9G6gGYvjp7IOLqbnNi8N2dAFuZGShPgn+X/+8PYaOnai1G
1khD7LBF9DqjS2cqJYEQSouK8LVu1mztXMP4X7V2VMIRqVp9uduVEpARFbap64D7
Y2Wv9nj8dX1mcSzGzW6fPJT1STThE3ByqCpSedOxo3CCE/hdbI2x+GNS3UkPSsib
MrvanZCEKiv8a1aJE03PJzvnblmM8pVS3wvr76w5XvQzGir1idhXU4OvG4ccxZpS
/pY3YSgmqwghjeYMdUd0b0jaKpqjHE+PoQ+nwhJ20K+7pVCQOAMjQxypelv+aTVD
0At3oeeEOMY4qXianM+6BQIDAQABo2EwXzAOBgNVHQ8BAf8EBAMCAqQwHQYDVR0l
BBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0O
BBYEFHw3ZFsv4ZUrcc/JvUIvGwA4sn1QMA0GCSqGSIb3DQEBCwUAA4IBAQB08jtx
zDfr+lzyv9mXDLptcn1XS2F25aGRAQ/XV8nmk2gaGKXEr3pAuoCAoUbQphmLcCvZ
mYIHbXQ/z6KIWonSc/a1ZGIwJfSpj3xt9lyJ72tWBodDkhIhB8l9hG6eHSCXYin9
fx1mvHJo7+jGHOc48//FcZ9vaiwi2e59gBaQRKUpEeUd648AbMuLXFYrshyr32bO
SEDEtk9wRajEds0uakXFxf5ymrABZpvde2VbgOcWqcd7tc76XPnqrxMrRhIgZlPa
vOPqBcUHaGFaWP3WWrQ/i0rq5uh428LD/mYE5El8M8GieVM+3mW8KpLkjkFT5Sc0
r1WcgSLKq+SV/Gip
-----END CERTIFICATE-----
`)

var demoKEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAn3kcOQSXlh5rA/lUBlIcUxNxGgRK1RJDFhsU9yGFL7MYUBe3
b7q10SKJRaDzo381Sz/TD64b/RVN2rDjgy+pA7Rol00xuryXjFFWMvZ5VfXQp5C4
u76R4pLKMRRNBw5jePXQWZWcPsEmObnc83+ta1iYXpVfkYj2Cq9pe/PyOmNpMCB3
T+YK8kkw9PN433IXu+MaOKj993QmBVdqCeKXbkoUrURLMxXi5G7PJiMB2ixVTsrA
btFQgR+XvK4x1S7xkm/y568UMRRRyLRosS2V8dhfO8flkNeyPKza1AuredLZnelU
hs623ZA4nN6UI+Z1XcCjRYnBfg3bPkXrS3mnUQIDAQABAoIBAFXeCGxLJLQYPNcu
8SdWHxo8ZbH0jbac1rKYcnl++w/sBzNZEdR/XFb3maJ8P7PRUwjpnOPchAWJ6xnO
FTMV/pOYGJkfX5+E3LUZNqjKPhsi+O7A5jdxLWwqTeSPYcpi3PzMnxsdi7velI6Q
nYAfR2l9ks2a8JKUhKbMPKgZelwlQVm4UwgJZZb+bn6T4Mx5leh/mdCqsVvNKB6m
ahIRB1giOrMhb220N4y7NptbvrrtZPB0qqxes6qnr1Y3eYQnSryMvMa5sB+yi2/b
CkS9sqNqf7Yts9OB6AG615tggTJtyj20AIOL0d/aiFjl+VbIVgxxRp5h/rVreUah
oz863gECgYEAz7YhL45mjjrpHQA1mliJvX9hYBYDOqhf0xrLv8TTrmuexP0P7c2E
8RNuwOksG+7se+3JIGcMMsU5RtkyI5VHfA9eDSFCHCzxguV9wY0A5NGUZoasthtC
3mbw8vey612oEqHCk0XaakZCV9tGb/JiSEsj7DOPHmaOwy0nQpdW1dkCgYEAxIwX
reXa2qzIj/G04Mmwy0ba+IXbjFvqkSl3z7MwgNN5UjP52GH/32lxUa/rCyFPSbIS
bt2gFqTG2HpFyxMm5ra59lk2MMQaI1BGih/hToN4lo/7QCpEIYUK03HZNec+tUIo
C22Nl8esxDw88w9SDimZgDpeP8uRjySZVcNMGjkCgYBHIVzN91sBfAUWjFrO52EM
BtIm4ILslHp0Ranemx3OjkZJuUu6KPZMxFXaND+JtVFAw1ZsBT31KPsLWxfDfbyE
LJMNtgT4tx9hrwtYu9vBgE/sqFP+7OkCVohO/CpGVcVX1BNY8cPxPuw7P/koHv4v
OaQsoB9zzrU2+4CFWmQ/SQKBgGL5UOtG9kBsBctGohkYN6kFkzrW3Un+904GHck/
qMsWst9MQSJPpzPvuxqxhaDjMzQfMd0WSYldjKxyVjb++/XuShLdtcY02hyyTfM8
Po708YKQGqujHQ/sGRmFGSZlvlQ0bkni7wxhhoSC+QZEzsNG+39w5QknD7OPcI+Z
evcxAoGBAJwfsbSFFLTJjYIi+55wcuD3qcmGSjqhncbi8oqfzJtbFSDnjPcKJIal
OxoU/5keftvaeUpGGiOaDowQLlNdzS+Ua0P8Hd5GcNNKMktbC6q5BN1lFLUTphWJ
jxzcmpeSQ8MB0SjU+yQ3ppDu5BW97wE0BN9joZTmwfR7hbcbg+yq
-----END RSA PRIVATE KEY-----
`)

// dns name: demo.com
var demoCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDaDCCAlCgAwIBAgIBATANBgkqhkiG9w0BAQsFADBGMQswCQYDVQQGEwJPTTER
MA8GA1UEChMIZGVtby5jb20xETAPBgNVBAsTCGRlbW8uY29tMREwDwYDVQQDEwhk
ZW1vLmNvbTAeFw0yMTA2MTgxMjIzMTZaFw0zMTA2MTYxMzIzMTZaMEYxCzAJBgNV
BAYTAk9NMREwDwYDVQQKEwhkZW1vLmNvbTERMA8GA1UECxMIZGVtby5jb20xETAP
BgNVBAMTCGRlbW8uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
n3kcOQSXlh5rA/lUBlIcUxNxGgRK1RJDFhsU9yGFL7MYUBe3b7q10SKJRaDzo381
Sz/TD64b/RVN2rDjgy+pA7Rol00xuryXjFFWMvZ5VfXQp5C4u76R4pLKMRRNBw5j
ePXQWZWcPsEmObnc83+ta1iYXpVfkYj2Cq9pe/PyOmNpMCB3T+YK8kkw9PN433IX
u+MaOKj993QmBVdqCeKXbkoUrURLMxXi5G7PJiMB2ixVTsrAbtFQgR+XvK4x1S7x
km/y568UMRRRyLRosS2V8dhfO8flkNeyPKza1AuredLZnelUhs623ZA4nN6UI+Z1
XcCjRYnBfg3bPkXrS3mnUQIDAQABo2EwXzAOBgNVHQ8BAf8EBAMCAqQwHQYDVR0l
BBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0O
BBYEFCy2GHsINcCVqu3urBKw0QOeJDQ0MA0GCSqGSIb3DQEBCwUAA4IBAQBp0ybq
1b4bPl1tgcn4vMCdyuhxp9Ek15zxMSejelpoU3Z3VBlU6selQoBX2x+TT5Pkx+h1
/Q6OXehkYuXbcNkqUlKMd0f7CWhYtfW4cd3ptDrp9D7wpSPy4zdRM7LO2wOWVZ49
GcS9Tu1qRAKQgjOu4omFpR9q5lG3y50v41TrozwTZGWk4nla8WuugmUzy9J9Amks
3ANUknBVyzK6lqcsp+WxOxVsu8LBzmOJDgMAIM+/4/LFNm2aIbkpkGS5b24YfsLL
Lgh7PneR3vG/ynWrXJpW/Khw8SWYY9OUVQ+ov230H0rtPZ0BYx0JRy9dT4A/LF5G
cb97v8dmkm2Po4dl
-----END CERTIFICATE-----
`)

var helloCert = []byte(`-----BEGIN CERTIFICATE-----
MIIDbjCCAlagAwIBAgIBATANBgkqhkiG9w0BAQsFADBJMQswCQYDVQQGEwJHRDES
MBAGA1UEChMJaGVsbG8uY29tMRIwEAYDVQQLEwloZWxsby5jb20xEjAQBgNVBAMT
CWhlbGxvLmNvbTAeFw0yMTA2MTkwMDE1MDFaFw0zMTA2MTcwMTE1MDFaMEkxCzAJ
BgNVBAYTAkdEMRIwEAYDVQQKEwloZWxsby5jb20xEjAQBgNVBAsTCWhlbGxvLmNv
bTESMBAGA1UEAxMJaGVsbG8uY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAqsEyNZzzb8Zupqcs2LCrZBnmOaKl1MdoVyxhnitKHR+y1IT3R1Fv3P/c
Tk6DYABnUeaiDC7GpmQRNl7XeNIwZZ9frZOEJQklJTCYGlIVtRNoqP3MS+ttalMD
nHjSACa/9VNFrNvzZY0OtoaYLcTnzAdiLCG7R3O4XTR77vzVYqKTT3oScjCyq4aP
cCgSc7MLG2CmWYvbzQERMzi8RgyTx/Ol+JziNJXdOtr0scW8xZ2jw4gWoVVWxAcT
ATTVqRz8hSqpceiv7ZeHjAeqn4uVAFSaPpiWaDaKD/V30/ty5iggZxmuXh/sF3wW
T4Vw+aHBkcI9WDeGapC0YANuQoGN1wIDAQABo2EwXzAOBgNVHQ8BAf8EBAMCAqQw
HQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8w
HQYDVR0OBBYEFBi1j87c0nVtdsPGx9VqN0eSHSqaMA0GCSqGSIb3DQEBCwUAA4IB
AQBJRHZQTlo3sUTY8R1d+zdANVPG/EKiQ/Qzg07eFf2/LiwVDuUzGXgw6q4SDUEI
VrSjD819E0Gp0/56okooD2qlv7NkQHOpGAhF89X/gdRUdtz4/CJwQiJ96Jlgc7nK
vPYjLgBeuizEEqOVZyTOp6mSEL8jTyC/Td0+yvJ5odMaMpA0Hq2i1HYYRcXBW+kq
1IHU0XVJTuVLsMkavU4HyhtHvokcG2l9bzZaGOfbKfisajUxXp7naN/+0+cFROtb
uMakvtCL1Jp9oy7Pmj+4yi2gNne8dHd2z4iX6pz6Np0I0dqctpJRSoXa9zKVUOVF
GUqLq7P1fOHyQsDzhWrD+43X
-----END CERTIFICATE-----
`)

var helloKEY = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAqsEyNZzzb8Zupqcs2LCrZBnmOaKl1MdoVyxhnitKHR+y1IT3
R1Fv3P/cTk6DYABnUeaiDC7GpmQRNl7XeNIwZZ9frZOEJQklJTCYGlIVtRNoqP3M
S+ttalMDnHjSACa/9VNFrNvzZY0OtoaYLcTnzAdiLCG7R3O4XTR77vzVYqKTT3oS
cjCyq4aPcCgSc7MLG2CmWYvbzQERMzi8RgyTx/Ol+JziNJXdOtr0scW8xZ2jw4gW
oVVWxAcTATTVqRz8hSqpceiv7ZeHjAeqn4uVAFSaPpiWaDaKD/V30/ty5iggZxmu
Xh/sF3wWT4Vw+aHBkcI9WDeGapC0YANuQoGN1wIDAQABAoIBAQCQNSFmTer51yfT
7xPc3TeiDo1013wdu1rPZFf88Kpi9kZdXP5JaOmER0GTkJM7HJwlexYYG9kA5Tn0
JRzsmPbunC59tTvA23xXcDbE49YZWw7kyZMj+uwpA3rlRtRz9EXhtjX9yrRAa2Sl
mf4jiUwJ76JliwdTTNPDQ3P3XegIp8r5LFTk//6g0MQ6qjn6rssAJ6FN6XgMaHMW
QAreMAfEPYbAm3fZijs+Rb6+y/0uPOfwxNfut7CdQLTu491OKs93ljxzKVXOe/Q6
MBF6uECzvi2R4IGzlP8mZia4aq+nczbavudqezcyYLwS6lpRz0mCY4a3RHkbH5+L
LX1Hv+vBAoGBAMXY4Tp0vgEeCiqnSaEQX4M5u+r5Q1ZESUYU3yBjMwBbVh556Ug6
EKTfoR0V2VIu2I4cse8lIFce+SUscldkkcsuypC4MXmplcQ2exXA5BUn68VWs8tg
iyF4qCNc+40IRShv4/qUzp3n04QsGR3rIJ8p8AvutPlnidslAyx8KJSpAoGBANzx
uae25i3lHCoextncD5+3Ogw4B9Bq+/2rbNOQbrkN5vJGBbnQeemOyLQCqoP2UYj2
IPN2SdV+rU5acfT64+XlshfnqgjiG/wtn1MfU4HJuYvFSXpm06xAnkkhWfYG3y8I
0JCu8PsXG5l3aDFnAaC5M+RvP8KTDfpHbdGF1B5/AoGAERxkvk2CcU5LysyVDZ0A
5bSEkBnmvPtC6xC7C24I5yr/E7uvdVOwRNIieQV+uiDbEc9hhDFNzrsbCSAC85P7
F/uAAWwsuzzzevjLRGJeV4YQWgzZl+lNnyN0Rzqvds8UTB8BNJbSF84I+RFnSrMf
KyTRYfbPKBLQVWeqEpraV6ECgYEAl65bhohJ/bgMXd5DJc2t7Dgd4cWVl7/av4uw
ao39dY3Vvv3TcH1vNKiRoQMzjOTNlPlkJcBPcAJHeEMfeM/FJU9LtJ2WXgLcs4Oe
nbIj4jZa61nF2AI/z0GNaSc8W2rcTa3/gVSYm8iBahpPrZrJw01iErFNVIcgUXI4
Ml9uAIECgYEAqt8i55jfMkQObmZh1r4sw7x/DRVEOzSJS26DjGvK/CXcIhso3OXs
+tTPCwmrHfnXWOmlIzY1qAzv5nUqdGYY92qGsmbwgbWHHgGDfn3zrZ706xNecZRD
YUsY7gAW9zuM2dW8g6j4pc3PGeA7bxzZdcXOVoY41dot8UTYIs+mqBo=
-----END RSA PRIVATE KEY-----
`)
