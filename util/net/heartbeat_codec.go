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
)

var (
	// a buffer of 64KB for received data
	readBufferSize = 64 * 1024
)

// read from conn and write to buffer
func (s *HeartbeatCodec) backgroundRead() {
	defer s.Close()
	buf := make([]byte, readBufferSize)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}
		// read from conn
		n, err := s.Conn.Read(buf)
		if err != nil {
			if s.isTimeout(err) {
				// read timeout is not a fatal error, just continue to wait for data.
				// the heartbeat goroutine will check the real connection status.
				continue
			}
			s.setError(err)
			return
		}
		if n == 0 {
			continue
		}
		// parse received data stream
		// a stream may contains multi-msg
		readBytes := buf[:n]
		for {
			// read flag
			if len(readBytes) < 2 {
				s.setError(fmt.Errorf("invalid heartbeat msg, remote: %s", s.Conn.RemoteAddr()))
				return
			}
			flag := binary.LittleEndian.Uint16(readBytes[:2])
			if flag != heartbeatCodecMsgFlag {
				s.setError(fmt.Errorf("unrecognized heartbeat msg type: %x, remote: %s", flag, s.Conn.RemoteAddr()))
				return
			}
			// read len
			if len(readBytes) < 6 {
				s.setError(fmt.Errorf("invalid heartbeat msg, remote: %s", s.Conn.RemoteAddr()))
				return
			}
			msgLen := binary.LittleEndian.Uint32(readBytes[2:6])
			// it's a heartbeat msg
			if msgLen == 0 {
				readBytes = readBytes[6:]
				if len(readBytes) == 0 {
					break
				}
				continue
			}
			// it's a data msg
			if int(msgLen) > len(readBytes)-6 {
				s.setError(fmt.Errorf("invalid heartbeat msg data, remote: %s", s.Conn.RemoteAddr()))
				return
			}
			data := readBytes[6 : 6+msgLen]
			s.readBufLock.Lock()
			s.readBuf.Write(data)
			s.readBufLock.Unlock()
			// signal the reader goroutine that data is available
			s.readCond.Signal()
			readBytes = readBytes[6+msgLen:]
			if len(readBytes) == 0 {
				break
			}
		}
	}
}

func (s *HeartbeatCodec) isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}

func (s *HeartbeatCodec) heartbeat() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	defer s.Close()
	for {
		select {
		case <-ticker.C:
			s.Conn.SetWriteDeadline(time.Now().Add(s.timeout))
			_, err := s.Write(nil)
			s.Conn.SetWriteDeadline(time.Time{})
			if err != nil {
				if s.isTimeout(err) {
					// write timeout is not a fatal error, just continue to send heartbeat.
					// the read goroutine will check the real connection status.
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
	s.writeLock.Lock()
	defer s.writeLock.Unlock()

	buf := make([]byte, 6+len(b))
	binary.LittleEndian.PutUint16(buf[0:2], heartbeatCodecMsgFlag)
	binary.LittleEndian.PutUint32(buf[2:6], uint32(len(b)))
	if len(b) > 0 {
		copy(buf[6:], b)
	}
	n, err = s.Conn.Write(buf)
	if err != nil {
		s.setError(err)
		return
	}
	n -= 6
	if n < 0 {
		n = 0
	}
	return
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
