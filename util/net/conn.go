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
)

type Codec interface {
	net.Conn
	Initialize(Context) error
}

type Conn struct {
	ctx Context
	net.Conn
	codec        []Codec
	readTimeout  time.Duration
	writeTimeout time.Duration
	readBytes    *int64
	writeBytes   *int64
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
	errIdx := -1
	defer func() {
		if err != nil {
			for i := 0; i < errIdx; i++ {
				s.codec[i].Close()
			}
		}
	}()
	for i, c := range s.codec {
		err = c.Initialize(s.ctx)
		if err != nil {
			errIdx = i
			return
		}
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

func (c *Conn) AddCodec(codec Codec) *Conn {
	c.codec = append(c.codec, codec)
	return c
}

func NewConn(conn net.Conn) *Conn {
	c := &Conn{
		Conn:       conn,
		writeBytes: new(int64),
		readBytes:  new(int64),
	}
	c.ctx = NewContext(c)
	return c
}

type EventConnErrorHandler func(ec *EventConn, err error)

type EventConnDataHandler func(ec *EventConn, data []byte)

type EventConnCloseHandler func(ec *EventConn)

type EventConn struct {
	conn                   net.Conn
	onReadError            EventConnErrorHandler
	onWriterError          EventConnErrorHandler
	onClose                EventConnCloseHandler
	onData                 EventConnDataHandler
	onCodecInitializeError EventConnErrorHandler
	readBufferSize         int
	readBytes              *int64
	writeBytes             *int64
	readTimeout            time.Duration
	writeTimeout           time.Duration
	closeOnce              sync.Once
	codec                  []Codec
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

func (s *EventConn) OnCodecInitializeError(h EventConnErrorHandler) *EventConn {
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

func (s *EventConn) OnReadError(onReadError EventConnErrorHandler) *EventConn {
	s.onReadError = onReadError
	return s
}

func (s *EventConn) OnWriterError(onWriterError EventConnErrorHandler) *EventConn {
	s.onWriterError = onWriterError
	return s
}

func (s *EventConn) OnClose(onClose EventConnCloseHandler) *EventConn {
	s.onClose = onClose
	return s
}

func (s *EventConn) OnData(onData EventConnDataHandler) *EventConn {
	s.onData = onData
	return s
}

func (s *EventConn) Close() {
	s.closeOnce.Do(func() {
		s.conn.Close()
		s.onClose(s)
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

		buf := make([]byte, s.readBufferSize)
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
		writeBytes:     new(int64),
		readBytes:      new(int64),
		readBufferSize: 1024 * 8,
		conn:           c,
		onReadError:    func(*EventConn, error) {},
		onWriterError:  func(*EventConn, error) {},
		onClose:        func(*EventConn) {},
		onData:         func(*EventConn, []byte) {},
	}
}

type BufferedConn struct {
	r *bufio.Reader
	net.Conn
}

func NewBufferedConn(c net.Conn) BufferedConn {
	return NewBufferedConnSize(c, 4096)
}

func NewBufferedConnSize(c net.Conn, n int) BufferedConn {
	return BufferedConn{bufio.NewReaderSize(c, n), c}
}

func (b BufferedConn) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func (b BufferedConn) Read(p []byte) (int, error) {
	return b.r.Read(p)
}
func (b BufferedConn) ReadByte() (byte, error) {
	return b.r.ReadByte()
}
func (b BufferedConn) UnreadByte() error {
	return b.r.UnreadByte()
}
func (b BufferedConn) Buffered() int {
	return b.r.Buffered()
}
