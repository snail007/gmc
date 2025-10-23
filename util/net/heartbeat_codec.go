// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	heartbeatCodecMsgFlag uint16 = 0x6699
	headerSize            = 6
)

var (
	// A buffer of 64KB for received data.
	readBufferSize = 64 * 1024
	// Pool for write buffers to reduce memory allocation.
	writeBufferPool = sync.Pool{
		New: func() interface{} {
			// The buffer is sized to accommodate the header plus a typical payload.
			// It will grow if necessary.
			return new(bytes.Buffer)
		},
	}
)

// backgroundRead reads from the connection and writes to the internal buffer.
func (s *HeartbeatCodec) backgroundRead() {
	defer s.Close()
	readBuf := make([]byte, readBufferSize)
	var data bytes.Buffer

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		// Read from the underlying connection.
		n, err := s.Conn.Read(readBuf)
		if err != nil {
			if s.isTimeout(err) {
				// A read timeout is not a fatal error, just continue to wait for data.
				// The heartbeat goroutine will check the real connection status.
				continue
			}
			s.setError(err)
			return
		}
		if n == 0 {
			continue
		}

		// Append new data to the buffer.
		data.Write(readBuf[:n])

		// Process all complete messages in the buffer.
		var pendingWrite []byte
		for {
			if data.Len() < 2 {
				// Not enough data to check the flag, wait for more.
				break
			}

			// First, check for the magic flag.
			flag := binary.LittleEndian.Uint16(data.Bytes()[:2])
			if flag != heartbeatCodecMsgFlag {
				s.setError(fmt.Errorf("unrecognized heartbeat msg type: %x, remote: %s", flag, s.Conn.RemoteAddr()))
				return
			}

			if data.Len() < headerSize {
				// We have the flag, but not the full header. Wait for more data.
				break
			}

			header := data.Bytes()[:headerSize]
			msgLen := binary.LittleEndian.Uint32(header[2:6])

			if msgLen == 0 { // It's a heartbeat message.
				data.Next(headerSize)
				continue
			}

			if uint32(data.Len()) < headerSize+msgLen {
				// Not enough data for the full message, wait for more.
				break
			}

			// We have a complete message.
			data.Next(headerSize)
			// Collect payload to write to the read buffer later.
			pendingWrite = append(pendingWrite, data.Next(int(msgLen))...)
		}

		if len(pendingWrite) > 0 {
			s.readBufLock.Lock()
			s.readBuf.Write(pendingWrite)
			s.readBufLock.Unlock()
			// Signal the reader goroutine that data is available.
			s.readCond.Signal()
		}
	}
}

func (s *HeartbeatCodec) isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

// heartbeat sends a heartbeat message periodically.
func (s *HeartbeatCodec) heartbeat() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	defer s.Close()
	for {
		select {
		case <-ticker.C:
			// The write deadline is set only for the heartbeat write.
			s.Conn.SetWriteDeadline(time.Now().Add(s.timeout))
			_, err := s.Write(nil)
			s.Conn.SetWriteDeadline(time.Time{})
			if err != nil {
				if s.isTimeout(err) {
					// A write timeout is not a fatal error, just continue to send heartbeats.
					// The read goroutine will check the real connection status.
					continue
				}
				s.setError(err)
				return
			}
		case <-s.ctx.Done():
			return
		}
	}
}

type HeartbeatCodec struct {
	net.Conn
	sync.Once
	writeLock   sync.Mutex
	readBuf     *bytes.Buffer
	readBufLock sync.Mutex
	readCond    *sync.Cond
	err         error
	errLock     sync.Mutex
	interval    time.Duration
	timeout     time.Duration
	ctx         context.Context
	ctxCancel   context.CancelFunc
}

func NewHeartbeatCodec() *HeartbeatCodec {
	s := &HeartbeatCodec{}
	s.readBuf = &bytes.Buffer{}
	s.readCond = sync.NewCond(&s.readBufLock)
	return s
}

func (s *HeartbeatCodec) Initialize(ctx Context) (err error) {
	if s.timeout == 0 {
		s.timeout = time.Second * 5
	}
	if s.interval == 0 {
		s.interval = time.Second * 5
	}
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	go s.backgroundRead()
	go s.heartbeat()
	return
}

func (s *HeartbeatCodec) SetConn(c net.Conn) Codec {
	s.Conn = c
	return s
}

func (s *HeartbeatCodec) Read(b []byte) (n int, err error) {
	if err = s.getError(); err != nil {
		return 0, err
	}
	s.readBufLock.Lock()
	defer s.readBufLock.Unlock()
	for s.readBuf.Len() == 0 {
		if err = s.getError(); err != nil {
			return 0, err
		}
		s.readCond.Wait()
	}
	return s.readBuf.Read(b)
}

func (s *HeartbeatCodec) Write(b []byte) (n int, err error) {
	if err = s.getError(); err != nil {
		return 0, err
	}

	buf := writeBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer writeBufferPool.Put(buf)

	// Write header directly into the buffer to avoid allocation.
	binary.Write(buf, binary.LittleEndian, heartbeatCodecMsgFlag)
	binary.Write(buf, binary.LittleEndian, uint32(len(b)))

	if len(b) > 0 {
		buf.Write(b)
	}

	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	written, err := s.Conn.Write(buf.Bytes())
	if err != nil {
		s.setError(err)
		return 0, err
	}

	n = written - headerSize
	if n < 0 {
		n = 0
	}
	return n, nil
}

func (s *HeartbeatCodec) Close() (err error) {
	s.Do(func() {
		s.ctxCancel()
		s.readCond.Broadcast() // Wake up all waiting readers
		err = s.Conn.Close()
	})
	return
}

func (s *HeartbeatCodec) Interval() time.Duration {
	return s.interval
}

func (s *HeartbeatCodec) SetInterval(interval time.Duration) *HeartbeatCodec {
	s.interval = interval
	return s
}

func (s *HeartbeatCodec) Timeout() time.Duration {
	return s.timeout
}

func (s *HeartbeatCodec) SetTimeout(timeout time.Duration) *HeartbeatCodec {
	s.timeout = timeout
	return s
}

func (s *HeartbeatCodec) getError() error {
	s.errLock.Lock()
	defer s.errLock.Unlock()
	return s.err
}

func (s *HeartbeatCodec) setError(err error) {
	s.errLock.Lock()
	defer s.errLock.Unlock()
	if s.err == nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		s.err = err
		// Wake up waiting readers to receive the error
		s.readCond.Broadcast()
	}
}