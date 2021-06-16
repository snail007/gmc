// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	gbytes "github.com/snail007/gmc/util/bytes"

	gmap "github.com/snail007/gmc/util/map"
)

const (
	defaultBufferedConnSize         = 8192
	defaultEventConnReadBufferSize  = 8192
	defaultConnBinderReadBufferSize = 8192
	defaultConnKeepAlivePeriod      = time.Second * 15
)

var (
	ErrHijacked = fmt.Errorf("ErrHijacked")
)

type BufferedConn interface {
	net.Conn
	Peek(n int) ([]byte, error)
	ReadByte() (byte, error)
	UnreadByte() error
	Buffered() int
	PeekMax(n int) (d []byte, err error)
}

type Codec interface {
	net.Conn
	Initialize(CodecContext, NextCodec) (net.Conn, error)
}

type NextCodec interface {
	Call(CodecContext) (net.Conn, error)
}

type codecs struct {
	codecs []Codec
	idx    int
}

func (s *codecs) Call(ctx CodecContext) (conn net.Conn, err error) {
	defer func() {
		if err == nil {
			return
		}
		// if occurs error, try clean.
		for _, codec := range s.codecs {
			func(c Codec) {
				defer func() { _ = recover() }()
				c.Close()
			}(codec)
		}
	}()
	return s.call(ctx)
}
func (s *codecs) call(ctx CodecContext) (conn net.Conn, err error) {
	s.idx++
	if s.idx == len(s.codecs) {
		// last call, no next codec, just return conn.
		return ctx.Conn(), nil
	}
	return s.codecs[s.idx].Initialize(ctx, s)
}

func newCodecs(cs []Codec) *codecs {
	f := &codecs{
		idx:    -1,
		codecs: cs,
	}
	return f
}

type ConnFilter func(ctx Context, c net.Conn, next NextConnFilter) (net.Conn, error)

type NextConnFilter interface {
	Call(ctx Context, c net.Conn) (net.Conn, error)
}

type connFilters struct {
	filters []ConnFilter
	idx     int
}

func (s *connFilters) Call(ctx Context, c net.Conn) (conn net.Conn, err error) {
	s.idx++
	if s.idx == len(s.filters) {
		// last call, no next filter, just return conn.
		return c, nil
	}
	return s.filters[s.idx](ctx, c, s)
}

func newConnFilters(filters []ConnFilter) *connFilters {
	f := &connFilters{
		idx:     -1,
		filters: filters,
	}
	return f
}

type Conn struct {
	ctx ConnContext
	net.Conn
	codec                     []Codec
	filters                   []ConnFilter
	readTimeout               time.Duration
	writeTimeout              time.Duration
	readBytes                 *int64
	writeBytes                *int64
	rawConn                   net.Conn
	onClose                   func(*Conn)
	closeOnce                 *sync.Once
	autoCloseOnReadWriteError bool
}

func (s *Conn) SetAutoCloseOnReadWriteError(b bool) {
	s.autoCloseOnReadWriteError = b
}

func (s *Conn) RawConn() net.Conn {
	return s.rawConn
}

func (s *Conn) AddFilter(f ConnFilter) {
	s.filters = append(s.filters, f)
}

func (s *Conn) ReadBytes() int64 {
	return atomic.LoadInt64(s.readBytes)
}

func (s *Conn) WriteBytes() int64 {
	return atomic.LoadInt64(s.writeBytes)
}

func (s *Conn) ReadTimeout() time.Duration {
	return s.readTimeout
}

func (s *Conn) SetReadTimeout(readTimeout time.Duration) *Conn {
	s.readTimeout = readTimeout
	return s
}

func (s *Conn) SetTimeout(readWriteTimeout time.Duration) *Conn {
	s.SetReadTimeout(readWriteTimeout)
	s.SetWriteTimeout(readWriteTimeout)
	return s
}

func (s *Conn) WriteTimeout() time.Duration {
	return s.writeTimeout
}

func (s *Conn) SetWriteTimeout(writeTimeout time.Duration) *Conn {
	s.writeTimeout = writeTimeout
	return s
}

func (s *Conn) Close() (err error) {
	s.closeOnce.Do(func() {
		err = s.Conn.Close()
		s.onClose(s)
	})
	return
}

func (s *Conn) Initialize() (err error) {
	// sets tcp keepalive
	if v, ok := s.rawConn.(*net.TCPConn); ok {
		v.SetKeepAlive(true)
		v.SetKeepAlivePeriod(defaultConnKeepAlivePeriod)
	}

	// init filters
	fConn, err := newConnFilters(s.filters).Call(s.ctx, s.Conn)
	if err != nil {
		if err == ErrHijacked {
			// hijacked by filter, just return nil.
			return nil
		}
		return err
	}
	s.Conn = fConn

	// init codec
	codecConn, err := newCodecs(s.codec).Call(s.ctx.(*defaultContext))
	if err != nil {
		if err == ErrHijacked {
			// hijacked by codec, just return nil.
			return nil
		}
		return err
	}
	s.Conn = codecConn
	return nil
}
func (s *Conn) Read(b []byte) (n int, err error) {
	defer func() {
		if n > 0 {
			atomic.AddInt64(s.readBytes, int64(n))
		}
		if err != nil && s.autoCloseOnReadWriteError {
			s.Close()
		}
	}()
	if s.readTimeout > 0 {
		s.Conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		defer s.Conn.SetReadDeadline(time.Time{})
	}
	n, err = s.Conn.Read(b)
	return
}

func (s *Conn) Write(b []byte) (n int, err error) {
	defer func() {
		if n > 0 {
			atomic.AddInt64(s.writeBytes, int64(n))
		}
		if err != nil && s.autoCloseOnReadWriteError {
			s.Close()
		}
	}()
	if s.writeTimeout > 0 {
		s.SetWriteDeadline(time.Now().Add(s.readTimeout))
		defer s.SetWriteDeadline(time.Time{})
	}
	n, err = s.Conn.Write(b)
	return
}

func (s *Conn) AddCodec(codec Codec) *Conn {
	s.codec = append(s.codec, codec)
	return s
}

func NewConn(conn net.Conn) *Conn {
	return NewContextConn(nil, conn)
}

func NewContextConn(ctx Context, conn net.Conn) *Conn {
	var c *Conn
	if v, ok := conn.(*Conn); ok {
		c = v
	} else {
		c = &Conn{
			Conn:       conn,
			writeBytes: new(int64),
			readBytes:  new(int64),
			rawConn:    conn,
			onClose:    func(conn *Conn) {},
			closeOnce:  &sync.Once{},
		}
	}
	if ctx == nil {
		c.ctx = NewContext().(*defaultContext)
	} else {
		c.ctx = ctx.(*defaultContext)
	}
	c.ctx.(*defaultContext).conn = c
	return c
}

type ConnErrorHandler func(ec *EventConn, err error)

type DataHandler func(ec *EventConn, data []byte)

type EventConnCloseHandler func(ec *EventConn)

type EventConn struct {
	conn                  net.Conn
	onReadError           ConnErrorHandler
	onWriterError         ConnErrorHandler
	onClose               EventConnCloseHandler
	onData                DataHandler
	onConnInitializeError ConnErrorHandler
	connFilters           []ConnFilter
	readBufferSize        int
	readBytes             *int64
	writeBytes            *int64
	readTimeout           time.Duration
	writeTimeout          time.Duration
	closeOnce             sync.Once
	codec                 []Codec
	data                  *gmap.Map
}

func (s *EventConn) AddConnFilter(f ConnFilter) *EventConn {
	s.connFilters = append(s.connFilters, f)
	return s
}

func (s *EventConn) ReadBytes() int64 {
	return atomic.LoadInt64(s.readBytes)
}

func (s *EventConn) WriteBytes() int64 {
	return atomic.LoadInt64(s.writeBytes)
}

func (s *EventConn) SetReadBufferSize(readBufferSize int) *EventConn {
	s.readBufferSize = readBufferSize
	return s
}

func (s *EventConn) OnConnInitializeError(h ConnErrorHandler) *EventConn {
	s.onConnInitializeError = h
	return s
}

func (s *EventConn) Write(b []byte) (n int, err error) {
	defer func() {
		if n > 0 {
			atomic.AddInt64(s.writeBytes, int64(n))
		}
	}()
	if s.writeTimeout > 0 {
		s.conn.SetWriteDeadline(time.Now().Add(s.readTimeout))
	}
	n, err = s.conn.Write(b)
	if s.writeTimeout > 0 {
		s.conn.SetWriteDeadline(time.Time{})
	}
	if err != nil {
		s.onWriterError(s, err)
		s.Close()
	}
	return
}

func (s *EventConn) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *EventConn) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *EventConn) ReadTimeout() time.Duration {
	return s.readTimeout
}

func (s *EventConn) SetReadTimeout(readTimeout time.Duration) *EventConn {
	s.readTimeout = readTimeout
	return s
}

func (s *EventConn) SetTimeout(readWriteTimeout time.Duration) *EventConn {
	s.SetReadTimeout(readWriteTimeout)
	s.SetWriteTimeout(readWriteTimeout)
	return s
}

func (s *EventConn) WriteTimeout() time.Duration {
	return s.writeTimeout
}

func (s *EventConn) SetWriteTimeout(writeTimeout time.Duration) *EventConn {
	s.writeTimeout = writeTimeout
	return s
}

func (s *EventConn) OnReadError(onReadError ConnErrorHandler) *EventConn {
	s.onReadError = onReadError
	return s
}

func (s *EventConn) OnWriterError(onWriterError ConnErrorHandler) *EventConn {
	s.onWriterError = onWriterError
	return s
}

func (s *EventConn) OnClose(onClose EventConnCloseHandler) *EventConn {
	s.onClose = onClose
	return s
}

func (s *EventConn) OnData(onData DataHandler) *EventConn {
	s.onData = onData
	return s
}

func (s *EventConn) Close() {
	s.closeOnce.Do(func() {
		s.conn.Close()
		s.onClose(s)
	})
}

func (s *EventConn) Data(key interface{}) interface{} {
	v, _ := s.data.Load(key)
	return v
}

func (s *EventConn) SetData(key, value interface{}) {
	s.data.Store(key, value)
}

func (s *EventConn) AddCodec(codec Codec) *EventConn {
	s.codec = append(s.codec, codec)
	return s
}

func (s *EventConn) Start() {
	go func() {
		defer func() {
			s.Close()
		}()

		// init Conn
		if len(s.codec) > 0 || len(s.connFilters) > 0 {
			conn := NewConn(s.conn)
			// add filters
			for _, f := range s.connFilters {
				conn.AddFilter(f)
			}
			// add codec
			for _, c := range s.codec {
				conn.AddCodec(c)
			}
			// init
			err := conn.Initialize()
			if err != nil {
				s.onConnInitializeError(s, err)
				return
			}
			s.conn = conn
		} else {
			if v, ok := s.conn.(*net.TCPConn); ok {
				v.SetKeepAlive(true)
				v.SetKeepAlivePeriod(defaultConnKeepAlivePeriod)
			}
		}

		buf := gbytes.GetPool(s.readBufferSize).Get().([]byte)
		defer gbytes.GetPool(s.readBufferSize).Put(buf)
		for {
			if s.readTimeout > 0 {
				s.conn.SetReadDeadline(time.Now().Add(s.readTimeout))
			}
			n, err := s.conn.Read(buf)
			if s.readTimeout > 0 {
				s.conn.SetReadDeadline(time.Time{})
			}
			if n > 0 {
				atomic.AddInt64(s.readBytes, int64(n))
			}
			if err != nil {
				s.onReadError(s, err)
				return
			}
			s.onData(s, buf[:n])
		}
	}()
}

func NewEventConn(c net.Conn) *EventConn {
	return &EventConn{
		data:           gmap.New(),
		writeBytes:     new(int64),
		readBytes:      new(int64),
		readBufferSize: defaultEventConnReadBufferSize,
		conn:           c,
		onReadError:    func(*EventConn, error) {},
		onWriterError:  func(*EventConn, error) {},
		onClose:        func(*EventConn) {},
		onData:         func(*EventConn, []byte) {},
	}
}

type defaultBufferedConn struct {
	r *bufio.Reader
	net.Conn
}

func NewBufferedConn(c net.Conn) BufferedConn {
	return NewBufferedConnSize(c, defaultBufferedConnSize)
}

func NewBufferedConnSize(c net.Conn, n int) BufferedConn {
	if v, ok := c.(BufferedConn); ok {
		return v
	}
	return &defaultBufferedConn{
		r:    bufio.NewReaderSize(c, n),
		Conn: c,
	}
}

func (b *defaultBufferedConn) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func (b *defaultBufferedConn) Read(p []byte) (int, error) {
	return b.r.Read(p)
}

func (b *defaultBufferedConn) ReadByte() (byte, error) {
	return b.r.ReadByte()
}

func (b *defaultBufferedConn) UnreadByte() error {
	return b.r.UnreadByte()
}

func (b *defaultBufferedConn) Buffered() int {
	return b.r.Buffered()
}

func (b *defaultBufferedConn) PeekMax(n int) (d []byte, err error) {
	_, err = b.ReadByte()
	if err != nil {
		return
	}
	err = b.UnreadByte()
	if err != nil {
		return
	}
	if n > b.r.Buffered() {
		n = b.r.Buffered()
	}
	return b.Peek(n)
}

type CloseHandler func(err error)

type ConnBinder struct {
	src          net.Conn
	dst          net.Conn
	onSrcClose   CloseHandler
	onDstClose   CloseHandler
	onClose      func()
	readBufSize  int
	trafficBytes *int64
}

func (s *ConnBinder) SetReadBufSize(readBufSize int) *ConnBinder {
	s.readBufSize = readBufSize
	return s
}

func (s *ConnBinder) OnSrcClose(onSrcClose CloseHandler) *ConnBinder {
	s.onSrcClose = onSrcClose
	return s
}

func (s *ConnBinder) OnDstClose(onDstClose CloseHandler) *ConnBinder {
	s.onDstClose = onDstClose
	return s
}

func (s *ConnBinder) OnClose(onClose func()) *ConnBinder {
	s.onClose = onClose
	return s
}

func (s *ConnBinder) copy(src, dst net.Conn) error {
	buf := gbytes.GetPool(s.readBufSize).Get().([]byte)
	defer func() {
		src.Close()
		dst.Close()
		gbytes.GetPool(s.readBufSize).Put(buf)
	}()
	for {
		n, err := src.Read(buf)
		if n > 0 {
			_, err = dst.Write(buf[:n])
			atomic.AddInt64(s.trafficBytes, int64(n))
		}
		if err != nil {
			return err
		}
	}
}

func (s *ConnBinder) Start() {
	g := sync.WaitGroup{}
	g.Add(2)
	go func() {
		go func() {
			defer g.Done()
			s.onSrcClose(s.copy(s.src, s.dst))
		}()
		go func() {
			defer g.Done()
			s.onDstClose(s.copy(s.dst, s.src))
		}()
		g.Wait()
		s.onClose()
	}()
}

func NewConnBinder(src net.Conn, dst net.Conn) *ConnBinder {
	return &ConnBinder{
		src:          src,
		dst:          dst,
		onSrcClose:   func(error) {},
		onDstClose:   func(error) {},
		onClose:      func() {},
		trafficBytes: new(int64),
		readBufSize:  defaultConnBinderReadBufferSize,
	}
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
