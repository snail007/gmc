// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"io"
	"net"
	"time"

	gbytes "github.com/snail007/gmc/util/bytes"
)

const (
	defaultBufferedConnSize         = 8192
	defaultEventConnReadBufferSize  = 8192
	defaultConnBinderReadBufferSize = 8192
	defaultConnKeepAlivePeriod      = time.Second * 15
	isTLSKey                        = "Gnet-TLSCodecIsTLSKey"
)

var (
	ErrFirstReadTimeout = fmt.Errorf("ErrFirstReadTimeout")
	netListen           = net.Listen
	netDialTimeout      = net.DialTimeout
)

type (
	Context interface {
		Clone() Context
		EventConn() *EventConn
		SetEventConn(eventConn *EventConn) Context
		EventListener() *EventListener
		SetEventListener(eventListener *EventListener) Context
		ConnBinder() *ConnBinder
		SetConnBinder(connBinder *ConnBinder) Context
		Hijacked() bool
		ReadBytes() int64
		Hijack(c ...Codec) (net.Conn, error)
		WriteBytes() int64
		Data(key interface{}) interface{}
		SetData(key, value interface{}) Context
		SetConn(c net.Conn) Context
		Conn() net.Conn
		RawConn() net.Conn
		Listener() *Listener
		SetListener(listener *Listener) Context
		ReadTimeout() time.Duration
		WriteTimeout() time.Duration
		RemoteAddr() net.Addr
		LocalAddr() net.Addr
		IsTLS() bool
	}
	Codec interface {
		net.Conn
		Initialize(ctx Context, next NextCodec) (net.Conn, error)
	}
	NextCodec interface {
		Call(ctx Context) (net.Conn, error)
	}
	NextConnFilter interface {
		Call(ctx Context, c net.Conn) (net.Conn, error)
	}
	ErrorHandler            func(ctx Context, err error)
	ConnFilter              func(ctx Context, c net.Conn, next NextConnFilter) (net.Conn, error)
	CloseHandler            func(ctx Context)
	DataHandler             func(ctx Context, data []byte)
	AcceptHandler           func(ctx Context, c net.Conn)
	CodecFactory            func(ctx Context) Codec
	FirstReadTimeoutHandler func(ctx Context, c net.Conn, err error)
)

// make sure all codec implements Codec
var _ Codec = &AESCodec{}
var _ Codec = &HeartbeatCodec{}
var _ Codec = &TLSServerCodec{}
var _ Codec = &TLSClientCodec{}

func WriteTo(writer io.Writer, data string) (n int, err error) {
	return writer.Write([]byte(data))
}

func Write(addr, data string) (err error) {
	return WriteBytes(addr, []byte(data))
}

func WriteBytes(addr string, data []byte) (err error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer func() {
		time.AfterFunc(time.Second, func() {
			c.Close()
		})
	}()
	_, err = c.Write(data)
	return
}

func Read(c net.Conn, bufSize int) (d string, err error) {
	data, err := ReadBytes(c, bufSize)
	d = string(data)
	return
}

func ReadBytes(c net.Conn, bufSize int) (d []byte, err error) {
	buf := gbytes.GetPool(bufSize).Get().([]byte)
	defer gbytes.GetPool(bufSize).Put(buf)
	n, err := c.Read(buf)
	if n > 0 {
		d = append(d, buf[:n]...)
	}
	return
}

func RandomListen(ip ...string) (l net.Listener, port string, err error) {
	ip0 := ""
	if len(ip) == 1 {
		ip0 = ip[0]
	}
	addr := net.JoinHostPort(ip0, "0")
	l, err = netListen("tcp", addr)
	if err != nil {
		return
	}
	_, port, _ = net.SplitHostPort(l.Addr().String())
	return
}
func RandomPort() (port string, err error) {
	maxTry := 3
	var l net.Listener
	defer func() {
		if l != nil {
			l.Close()
		}
	}()
	for i := 0; i < maxTry; i++ {
		l, port, err = RandomListen()
		if err == nil {
			port = NewAddr(l.Addr()).Port()
			return
		}
		time.Sleep(time.Millisecond * 10)
	}
	return
}

func Listen(addr string) (l *Listener, err error) {
	lis, err := netListen("tcp", addr)
	if err != nil {
		return
	}
	l = NewListener(lis)
	return
}

func ListenEvent(addr string) (l *EventListener, err error) {
	l0, err := Listen(addr)
	if err != nil {
		return
	}
	l = NewEventListener(l0)
	return
}

func Dial(addr string, timeout time.Duration) (c *Conn, err error) {
	c0, err := netDialTimeout("tcp", addr, timeout)
	if err != nil {
		return
	}
	c = NewConn(c0)
	return
}

type Addr struct {
	network string
	addr    string
	port    string
	host    string
	err     error
}

func (s *Addr) Err() error {
	return s.err
}

func (s *Addr) Network() string {
	return s.network
}

func (s *Addr) String() string {
	return s.addr
}

func (s *Addr) Port() string {
	if s.port == "" {
		return "0"
	}
	return s.port
}

// PortAddr returns address string with 'withIP:Addr.Port()',
// 'withIP' can be a valid ip or domain or empty.
func (s *Addr) PortAddr(withIP ...string) string {
	ip0 := ""
	if len(withIP) == 1 {
		ip0 = withIP[0]
	}
	return net.JoinHostPort(ip0, s.Port())
}

// PortLocalAddr returns address string with '127.0.0.1:Addr.Port()',
func (s *Addr) PortLocalAddr() string {
	return s.PortAddr("127.0.0.1")
}

func (s *Addr) Host() string {
	return s.host
}

func NewAddr(addr net.Addr) *Addr {
	h, p, err := net.SplitHostPort(addr.String())
	return &Addr{
		addr:    addr.String(),
		network: addr.Network(),
		host:    h,
		port:    p,
		err:     err,
	}
}

func NewTCPAddr(address string) *Addr {
	h, p, err := net.SplitHostPort(address)
	return &Addr{
		addr:    address,
		network: "tcp",
		host:    h,
		port:    p,
		err:     err,
	}
}
