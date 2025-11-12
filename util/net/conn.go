// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"bufio"
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"

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
}

func (s *codecs) Call(ctx Context) (err error) {
	for _, c := range s.codecs {
		c.SetConn(ctx.Conn())
		err = c.Initialize(ctx)
		if err != nil {
			return
		}
		if ctx.IsContinue() {
			ctx.(*defaultContext).isContinue = false
			continue
		}
		ctx.SetConn(c)
		if ctx.IsBreak() {
			return
		}
	}
	return nil
}

func newCodecs(cs []Codec) *codecs {
	f := &codecs{
		codecs: cs,
	}
	return f
}

type connFilters struct {
	filters []ConnFilter
}

func (s *connFilters) Call(ctx Context, c net.Conn) (net.Conn, error) {
	var err error
	conn := c
	for _, cf := range s.filters {
		conn, err = cf(ctx, conn)
		if err != nil {
			return nil, err
		}
		if ctx.IsContinue() {
			ctx.(*defaultContext).isContinue = false
			continue
		}
		if ctx.IsBreak() {
			return conn, err
		}
	}
	return conn, err
}

func newConnFilters(filters []ConnFilter) *connFilters {
	f := &connFilters{
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
	onCloseHook               func(*Conn)
	onInitializeHook          func(*Conn)
	closeOnce                 *sync.Once
	autoCloseOnReadWriteError bool
	initOnce                  *sync.Once
	maxIdleTimeout            time.Duration
	touchTime                 time.Time
	idleChn                   chan bool
	onIdleTimeout             func(*Conn)
	readLimiter, writeLimiter *rate.Limiter
	rwCtx                     context.Context
	createTime                time.Time
}

func (s *Conn) ReadLimiter() *rate.Limiter {
	return s.readLimiter
}

func (s *Conn) SetReadLimiter(readLimiter *rate.Limiter) {
	s.readLimiter = readLimiter
}

func (s *Conn) WriteLimiter() *rate.Limiter {
	return s.writeLimiter
}

func (s *Conn) SetWriteLimiter(writeLimiter *rate.Limiter) {
	s.writeLimiter = writeLimiter
}

func (s *Conn) SetOnClose(f func(conn *Conn)) {
	s.onCloseHook = f
}

func (s *Conn) SetOnInitialize(f func(conn *Conn)) {
	s.onInitializeHook = f
}

func (s *Conn) SetOnIdleTimeout(f func(conn *Conn)) {
	s.onIdleTimeout = f
}

func (s *Conn) TouchTime() time.Time {
	return s.touchTime
}

func (s *Conn) touch() {
	s.touchTime = time.Now()
}

func (s *Conn) SetMaxIdleTimeout(d time.Duration) {
	s.maxIdleTimeout = d
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
		close(s.idleChn)
		s.onClose(s)
		s.onCloseHook(s)
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
	if err != nil {
		return err
	}
	s.Conn = fConn

	// init codec
	err = newCodecs(s.codec).Call(s.ctx)
	if err != nil {
		return err
	}
	s.Conn = s.ctx.Conn()

	//init idle monitor, this must be after filters and codecs
	s.initIdleTimeoutDaemon()

	//call Initialize hook
	s.onInitializeHook(s)

	return nil
}

func (s *Conn) initIdleTimeoutDaemon() {
	if s.maxIdleTimeout <= 0 {
		return
	}
	go func() {
		t := time.NewTimer(s.maxIdleTimeout)
		defer t.Stop()
		for {
			t.Reset(s.maxIdleTimeout)
			select {
			case <-s.idleChn:
				return
			case <-t.C:
				if s.touchTime.Add(s.maxIdleTimeout).Before(time.Now()) {
					//idle timeout, close the connection
					s.Close()
					if s.onIdleTimeout != nil {
						s.onIdleTimeout(s)
					}
					return
				}
			}
		}
	}()
}

func (s *Conn) Read(b []byte) (n int, err error) {
	if e := s.doInitialize(); e != nil {
		return 0, e
	}
	defer func() {
		if n > 0 {
			atomic.AddInt64(s.readBytes, int64(n))
			s.touch()
		}
		if err != nil && s.autoCloseOnReadWriteError {
			s.Close()
		}
	}()
	if s.readTimeout > 0 {
		s.Conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		defer s.Conn.SetReadDeadline(time.Time{})
	}
	if s.readLimiter == nil {
		return s.Conn.Read(b)
	}
	n, err = s.Conn.Read(b)
	if err != nil {
		return n, err
	}
	err = s.readLimiter.WaitN(s.rwCtx, n)
	return
}

func (s *Conn) Write(b []byte) (n int, err error) {
	if e := s.doInitialize(); e != nil {
		return 0, e
	}
	defer func() {
		if n > 0 {
			atomic.AddInt64(s.writeBytes, int64(n))
			s.touch()
		}
		if err != nil && s.autoCloseOnReadWriteError {
			s.Close()
		}
	}()
	if s.writeTimeout > 0 {
		s.SetWriteDeadline(time.Now().Add(s.readTimeout))
		defer s.SetWriteDeadline(time.Time{})
	}
	if s.writeLimiter == nil {
		return s.Conn.Write(b)
	}
	n, err = s.Conn.Write(b)
	if err != nil {
		return n, err
	}
	err = s.writeLimiter.WaitN(s.rwCtx, n)
	return
}

func (s *Conn) AddCodec(codec Codec) *Conn {
	s.codec = append(s.codec, codec)
	return s
}

func NewConn(conn net.Conn, force ...bool) *Conn {
	return newContextConn(NewContext(), conn, force...)
}

func NewContextConn(ctx Context, conn net.Conn, force ...bool) *Conn {
	return newContextConn(ctx, conn, force...)
}

func newContextConn(ctx Context, conn net.Conn, f ...bool) *Conn {
	force := false
	if len(f) == 1 {
		force = f[0]
	}
	var c *Conn
	if v, ok := conn.(*Conn); ok && !force {
		c = v
	} else {
		c = &Conn{
			Conn:                      conn,
			createTime:                time.Now(),
			writeBytes:                new(int64),
			readBytes:                 new(int64),
			rawConn:                   conn,
			onClose:                   func(conn *Conn) {},
			onCloseHook:               func(conn *Conn) {},
			onInitializeHook:          func(conn *Conn) {},
			closeOnce:                 &sync.Once{},
			initOnce:                  &sync.Once{},
			autoCloseOnReadWriteError: true,
			idleChn:                   make(chan bool),
			touchTime:                 time.Now(),
			rwCtx:                     context.Background(),
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
	started        bool
	createTime     time.Time
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

func (s *EventConn) StartAndWait() {
	if s.started {
		return
	}
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
}

func (s *EventConn) Start() {
	go s.StartAndWait()
}

func NewEventConn(c net.Conn) *EventConn {
	return NewContextEventConn(NewContext(), c)
}

func NewContextEventConn(ctx Context, c net.Conn) *EventConn {
	ec := &EventConn{
		ctx:            ctx,
		createTime:     time.Now(),
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
	src               net.Conn
	dst               net.Conn
	onSrcClose        CloseHandler
	onDstClose        CloseHandler
	onClose           func()
	readBufSize       int
	trafficBytes      *int64
	ctx               Context
	started           bool
	autoClose         bool
	err               error
	afterSrcFirstRead func(b []byte) []byte
	afterDstFirstRead func(b []byte) []byte
}

func (s *ConnBinder) SetAfterSrcFirstRead(afterSrcFirstRead func(b []byte) []byte) {
	s.afterSrcFirstRead = afterSrcFirstRead
}

func (s *ConnBinder) SetAfterDstFirstRead(afterDstFirstRead func(b []byte) []byte) {
	s.afterDstFirstRead = afterDstFirstRead
}

func (s *ConnBinder) setError(err error) {
	if s.err != nil || err == nil {
		return
	}
	s.err = err
}

func (s *ConnBinder) Error() error {
	return s.err
}

func (s *ConnBinder) SetAutoClose(autoClose bool) {
	s.autoClose = autoClose
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

func (s *ConnBinder) copy(a, b net.Conn, aIsSrc bool) error {
	buf := gbytes.GetPool(s.readBufSize).Get().([]byte)
	defer func() {
		if s.autoClose {
			a.Close()
			b.Close()
		}
		gbytes.GetPool(s.readBufSize).Put(buf)
	}()
	isFirst := true
	for {
		n, err := a.Read(buf)
		if err != nil {
			err = errors.Wrap(err, "failed to read from src: "+a.RemoteAddr().String())
		}
		if n > 0 {
			if isFirst {
				isFirst = false
				if aIsSrc {
					if s.afterSrcFirstRead != nil {
						_, err = b.Write(s.afterSrcFirstRead(buf[:n]))
					} else {
						_, err = b.Write(buf[:n])
					}
				} else {
					if s.afterDstFirstRead != nil {
						_, err = b.Write(s.afterDstFirstRead(buf[:n]))
					} else {
						_, err = b.Write(buf[:n])
					}
				}
			} else {
				_, err = b.Write(buf[:n])
			}
			if err != nil {
				err = errors.Wrap(err, "failed to write to dst: "+b.RemoteAddr().String())
			}
			atomic.AddInt64(s.trafficBytes, int64(n))
		}
		if err != nil {
			return err
		}
	}
}

func (s *ConnBinder) StartAndWait() {
	if s.started {
		return
	}
	g := sync.WaitGroup{}
	g.Add(2)
	go func() {
		defer g.Done()
		s.setError(s.copy(s.src, s.dst, true))
		s.onSrcClose(s.ctx)
	}()
	go func() {
		defer g.Done()
		s.setError(s.copy(s.dst, s.src, false))
		s.onDstClose(s.ctx)
	}()
	g.Wait()
	s.onClose()
}

func (s *ConnBinder) Start() {
	go s.StartAndWait()
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
		autoClose:    true,
	}
	cb.ctx.SetConnBinder(cb)
	return cb
}

func NewConnBinder(src net.Conn, dst net.Conn) *ConnBinder {
	return NewContextConnBinder(NewContext(), src, dst)
}
