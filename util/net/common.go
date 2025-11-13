// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"fmt"
	"io"
	"net"
	"strings"
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
		IsHijacked() bool
		Hijack()
		IsBreak() bool
		Break()
		IsContinue() bool
		Continue()
		ReadBytes() int64
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
		SetConn(conn net.Conn) Codec
		Initialize(ctx Context) error
	}
	ErrorHandler            func(ctx Context, err error)
	ConnFilter              func(ctx Context, c net.Conn) (net.Conn, error)
	CloseHandler            func(ctx Context)
	DataHandler             func(ctx Context, data []byte)
	AcceptHandler           func(ctx Context, c net.Conn)
	CodecFactory            func(ctx Context) Codec
	FirstReadTimeoutHandler func(ctx Context, c net.Conn, err error)
	BeforeFirstReadHandler  func(ctx Context, c net.Conn) error
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
func RandomPort(ip ...string) (port string, err error) {
	maxTry := 3
	var l net.Listener
	defer func() {
		if l != nil {
			l.Close()
		}
	}()
	for i := 0; i < maxTry; i++ {
		l, port, err = RandomListen(ip...)
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

func LocalOutgoingIP() (string, error) {
	addr, err := net.ResolveUDPAddr("udp", "1.2.3.4:1")
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	host, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return "", err
	}

	return host, nil
}

func IsPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		// 10.0.0.0/8
		// 172.16.0.0/12
		// 192.168.0.0/16
		return ipv4[0] == 10 || (ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31) || (ipv4[0] == 192 && ipv4[1] == 168)
	}
	if strings.HasPrefix(ip, "fc00:") || strings.HasPrefix(ip, "fd00:") || strings.HasPrefix(ip, "fe80:") {
		return true
	}
	return false
}
