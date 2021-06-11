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
	defaultBufferedConnSize        = 4096
	defaultEventConnReadBufferSize = 8192
)

var (
	ErrCodecSkipped      = fmt.Errorf("")
	ErrConnFilterSkipped = fmt.Errorf("")
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
	Initialize(ConnContext) error
}

type ConnFilter func(ctx Context, c net.Conn) (net.Conn, error)

type Conn struct {
	ctx ConnContext
	net.Conn
	codec        []Codec
	filters      []ConnFilter
	readTimeout  time.Duration
	writeTimeout time.Duration
	readBytes    *int64
	writeBytes   *int64
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

func (s *Conn) Initialize() (err error) {
	var okayIdx []int
	defer func() {
		if err != nil {
			for i := range okayIdx {
				s.codec[i].Close()
			}
		}
	}()

	// init filters
	for _, f := range s.filters {
		conn, err := f(s.ctx, s.Conn)
		if err != nil {
			if err == ErrConnFilterSkipped {
				continue
			}
			return err
		}
		s.Conn = conn
	}

	// init codec
	for i, c := range s.codec {
		err = c.Initialize(s.ctx)
		if err != nil {
			if err == ErrCodecSkipped {
				continue
			}
			return
		}
		okayIdx = append(okayIdx, i)
		s.Conn = c
	}
	return nil
}
func (s *Conn) Read(b []byte) (n int, err error) {
	defer func() {
		if n > 0 {
			atomic.AddInt64(s.readBytes, int64(n))
		}
	}()
	if s.readTimeout > 0 {
		s.Conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		defer s.Conn.SetReadDeadline(time.Time{})
	}
	return s.Conn.Read(b)
}

func (s *Conn) Write(b []byte) (n int, err error) {
	defer func() {
		if n > 0 {
			atomic.AddInt64(s.writeBytes, int64(n))
		}
	}()
	if s.writeTimeout > 0 {
		s.SetWriteDeadline(time.Now().Add(s.readTimeout))
		defer s.SetWriteDeadline(time.Time{})
	}
	return s.Conn.Write(b)
}

func (s *Conn) AddCodec(codec Codec) *Conn {
	s.codec = append(s.codec, codec)
	return s
}

func NewConn(conn net.Conn) *Conn {
	return NewContextConn(nil, conn)
}

func NewContextConn(ctx Context, conn net.Conn) *Conn {
	c := &Conn{
		Conn:       conn,
		writeBytes: new(int64),
		readBytes:  new(int64),
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

type ConnCloseHandler func(ec *EventConn)

type EventConn struct {
	conn                   net.Conn
	onReadError            ConnErrorHandler
	onWriterError          ConnErrorHandler
	onClose                ConnCloseHandler
	onData                 DataHandler
	onCodecInitializeError ConnErrorHandler
	readBufferSize         int
	readBytes              *int64
	writeBytes             *int64
	readTimeout            time.Duration
	writeTimeout           time.Duration
	closeOnce              sync.Once
	codec                  []Codec
	data                   *gmap.Map
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

func (s *EventConn) OnCodecInitializeError(h ConnErrorHandler) *EventConn {
	s.onCodecInitializeError = h
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

func (s *EventConn) OnClose(onClose ConnCloseHandler) *EventConn {
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

		// codec check
		if len(s.codec) > 0 {
			conn := NewConn(s.conn)
			for _, c := range s.codec {
				conn.AddCodec(c)
			}
			err := conn.Initialize()
			if err != nil {
				s.onCodecInitializeError(s, err)
				return
			}
			s.conn = conn
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
