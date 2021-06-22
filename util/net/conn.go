// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"bufio"
	"net"
	"sync"
	"sync/atomic"
	"time"

	gbytes "github.com/snail007/gmc/util/bytes"
)

type BufferedConn interface {
	net.Conn
	Peek(n int) ([]byte, error)
	ReadByte() (byte, error)
	UnreadByte() error
	Buffered() int
	PeekMax(n int) (d []byte, err error)
}

type codecs struct {
	codecs []Codec
	idx    int
}

func (s *codecs) Call(ctx Context) (conn net.Conn, err error) {
	conn, err = s.call(ctx)
	if err != nil {
		// if occurs error, try clean.
		for _, codec := range s.codecs {
			func() {
				defer func() {
					_ = recover()
				}()
				codec.Close()
			}()
		}
	}
	return
}

func (s *codecs) call(ctx Context) (conn net.Conn, err error) {
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
	net.Conn
	ctx                       Context
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
	initOnce                  *sync.Once
}

func (s *Conn) Ctx() Context {
	return s.ctx
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
		func() {
			defer func() { _ = recover() }()
			err = s.Conn.Close()
		}()
		s.onClose(s)
	})
	return
}
func (s *Conn) doInitialize() (err error) {
	s.initOnce.Do(func() {
		err = s.initialize()
	})
	return
}
func (s *Conn) initialize() (err error) {
	// sets tcp keepalive
	if v, ok := s.rawConn.(*net.TCPConn); ok {
		v.SetKeepAlive(true)
		v.SetKeepAlivePeriod(defaultConnKeepAlivePeriod)
	}

	// init filters
	fConn, err := newConnFilters(s.filters).Call(s.ctx, s.Conn)
	// checking hijack
	if s.ctx.Hijacked() {
		// hijacked by filter, just return nil.
		return nil
	}
	if err != nil {
		return err
	}
	s.Conn = fConn

	// init codec
	codecConn, err := newCodecs(s.codec).Call(s.ctx.(*defaultContext))
	// checking hijack
	if s.ctx.Hijacked() {
		// hijacked by codec, just return nil.
		return nil
	}
	if err != nil {
		return err
	}
	s.Conn = codecConn
	return nil
}
func (s *Conn) Read(b []byte) (n int, err error) {
	if e := s.doInitialize(); e != nil {
		return 0, e
	}
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
	if e := s.doInitialize(); e != nil {
		return 0, e
	}
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
	return NewContextConn(NewContext(), conn)
}

func NewContextConn(ctx Context, conn net.Conn) *Conn {
	var c *Conn
	if v, ok := conn.(*Conn); ok {
		c = v
	} else {
		c = &Conn{
			Conn:                      conn,
			writeBytes:                new(int64),
			readBytes:                 new(int64),
			rawConn:                   conn,
			onClose:                   func(conn *Conn) {},
			closeOnce:                 &sync.Once{},
			initOnce:                  &sync.Once{},
			autoCloseOnReadWriteError: true,
		}
	}
	ctx.(*defaultContext).conn = c
	c.ctx = ctx
	return c
}

type EventConn struct {
	conn           net.Conn
	onReadError    ErrorHandler
	onWriterError  ErrorHandler
	onClose        CloseHandler
	onData         DataHandler
	connFilters    []ConnFilter
	readBufferSize int
	readBytes      *int64
	writeBytes     *int64
	readTimeout    time.Duration
	writeTimeout   time.Duration
	closeOnce      sync.Once
	codec          []Codec
	ctx            Context
}

func (s *EventConn) Ctx() Context {
	return s.ctx
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
		s.onWriterError(s.ctx, err)
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

func (s *EventConn) OnReadError(onReadError ErrorHandler) *EventConn {
	s.onReadError = onReadError
	return s
}

func (s *EventConn) OnWriterError(onWriterError ErrorHandler) *EventConn {
	s.onWriterError = onWriterError
	return s
}

func (s *EventConn) OnClose(onClose CloseHandler) *EventConn {
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
		s.onClose(s.ctx)
	})
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
				s.onReadError(s.ctx, err)
				return
			}
			s.onData(s.ctx, buf[:n])
		}
	}()
}

func NewEventConn(c net.Conn) *EventConn {
	return NewContextEventConn(NewContext(), c)
}

func NewContextEventConn(ctx Context, c net.Conn) *EventConn {
	ec := &EventConn{
		ctx:            ctx,
		writeBytes:     new(int64),
		readBytes:      new(int64),
		readBufferSize: defaultEventConnReadBufferSize,
		conn:           c,
		onReadError:    func(Context, error) {},
		onWriterError:  func(Context, error) {},
		onClose:        func(Context) {},
		onData:         func(Context, []byte) {},
	}
	ec.ctx.SetEventConn(ec)
	return ec
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

type ConnBinder struct {
	src          net.Conn
	dst          net.Conn
	onSrcClose   CloseHandler
	onDstClose   CloseHandler
	onClose      func()
	readBufSize  int
	trafficBytes *int64
	ctx          Context
}

func (s *ConnBinder) Ctx() Context {
	return s.ctx
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
			s.copy(s.src, s.dst)
			s.onSrcClose(s.ctx)
		}()
		go func() {
			defer g.Done()
			s.copy(s.dst, s.src)
			s.onDstClose(s.ctx)
		}()
		g.Wait()
		s.onClose()
	}()
}

func NewContextConnBinder(ctx Context, src net.Conn, dst net.Conn) *ConnBinder {
	cb := &ConnBinder{
		ctx:          ctx,
		src:          src,
		dst:          dst,
		onSrcClose:   func(Context) {},
		onDstClose:   func(Context) {},
		onClose:      func() {},
		trafficBytes: new(int64),
		readBufSize:  defaultConnBinderReadBufferSize,
	}
	cb.ctx.SetConnBinder(cb)
	return cb
}

func NewConnBinder(src net.Conn, dst net.Conn) *ConnBinder {
	return NewContextConnBinder(NewContext(), src, dst)
}
