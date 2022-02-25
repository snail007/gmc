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

type heartbeatCodecMsg struct {
	flag uint16
	len  uint32
	data []byte
}

func (h *heartbeatCodecMsg) Data() []byte {
	if h.len == 0 {
		return nil
	}
	d := make([]byte, h.len)
	copy(d, h.data)
	return d
}

func (h *heartbeatCodecMsg) DataLength() uint32 {
	return h.len
}

func (h *heartbeatCodecMsg) SetData(data []byte) {
	h.len = uint32(len(data))
	h.data = data
}

func (h *heartbeatCodecMsg) ReadFrom(r net.Conn) (err error) {
	var flag uint16
	err = binary.Read(r, binary.LittleEndian, &flag)
	if err != nil {
		return
	}
	if flag != heartbeatCodecMsgFlag {
		return fmt.Errorf("unrecognized msg type: %x, remote: %s", flag, r.RemoteAddr().String())
	}
	err = binary.Read(r, binary.LittleEndian, &h.len)
	if err != nil {
		return
	}
	if h.len == 0 {
		return
	}
	h.data = make([]byte, h.len)
	_, err = io.ReadFull(r, h.data)
	if err != nil {
		h.data = nil
		return
	}
	return
}

func (h *heartbeatCodecMsg) Bytes() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, h.flag)
	binary.Write(buf, binary.LittleEndian, uint32(len(h.data)))
	if h.len > 0 {
		binary.Write(buf, binary.LittleEndian, h.data)
	}
	return buf.Bytes()
}

func newHeartbeatCodecMsg() *heartbeatCodecMsg {
	return &heartbeatCodecMsg{
		flag: heartbeatCodecMsgFlag,
	}
}

type HeartbeatCodec struct {
	net.Conn
	sync.Mutex
	sync.Once
	interval    time.Duration
	timeout     time.Duration
	bufReader   *io.PipeReader
	bufWriter   *io.PipeWriter
	readErrChn  chan error
	writeErrChn chan error
	ctx         context.Context
	ctxCancel   context.CancelFunc
}

func (s *HeartbeatCodec) sendReadErr(err error) {
	select {
	case s.readErrChn <- err:
	default:
	}
}

func (s *HeartbeatCodec) sendWriteErr(err error) {
	select {
	case s.writeErrChn <- err:
	default:
	}
}
func (s *HeartbeatCodec) tempDelay(tempDelay time.Duration) time.Duration {
	if tempDelay == 0 {
		tempDelay = time.Millisecond * 5
	}
	if tempDelay > time.Second {
		tempDelay = time.Second
	} else {
		tempDelay *= 2
	}
	return tempDelay
}
func (s *HeartbeatCodec) heartbeat() {
	done := make(chan error)
	tempDelay := time.Duration(0)
retry:
	for {
		go func() {
			s.SetWriteDeadline(time.Now().Add(s.timeout))
			_, err := s.Write(nil)
			s.SetWriteDeadline(time.Time{})
			done <- err
		}()
		select {
		case err := <-done:
			if err == nil {
				break
			}
			if e, ok := err.(net.Error); ok && e.Temporary() {
				tempDelay = s.tempDelay(tempDelay)
				time.Sleep(tempDelay)
				continue retry
			}
			s.sendWriteErr(err)
			return
		case <-s.ctx.Done():
			return
		}
		tempDelay = 0
		time.Sleep(s.interval)
	}
}
func (s *HeartbeatCodec) backgroundRead() {
	tempDelay := time.Duration(0)
	msg := newHeartbeatCodecMsg()
	out := make(chan error)
retry:
	for {
		go func() {
			out <- msg.ReadFrom(s.Conn)
		}()
		select {
		case err := <-out:
			if err == nil {
				break
			}
			if e, ok := err.(net.Error); ok && e.Temporary() {
				tempDelay = s.tempDelay(tempDelay)
				time.Sleep(tempDelay)
				continue retry
			}
			s.sendReadErr(err)
			return
		case <-s.ctx.Done():
			return
		}
		if msg.DataLength() == 0 {
			// it's a heartbeat msg.
			continue
		}
		// it's a data msg.
		s.bufWriter.Write(msg.Data())
	}
}

func (s *HeartbeatCodec) Read(b []byte) (n int, err error) {
	done := make(chan bool)
	go func() {
		defer close(done)
		n, err = s.bufReader.Read(b)
		_ = err
	}()
	select {
	case <-done:
	case err = <-s.readErrChn:
	}
	return
}

func (s *HeartbeatCodec) Write(b []byte) (n int, err error) {
	s.Lock()
	defer s.Unlock()
	msg := newHeartbeatCodecMsg()
	msg.SetData(b)
	done := make(chan bool)
	go func() {
		defer close(done)
		n, err = s.Conn.Write(msg.Bytes())
		if err == nil {
			// decrease the n with msg header length 6 bytes
			n = n - 6
		}
	}()
	select {
	case <-done:
	case err = <-s.writeErrChn:
	}
	return
}

func (s *HeartbeatCodec) Close() (err error) {
	s.Do(func() {
		s.ctxCancel()
		s.bufReader.Close()
		s.bufWriter.Close()
		if s.Conn != nil {
			err = s.Conn.Close()
		}
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

func (s *HeartbeatCodec) Initialize(ctx Context, next NextCodec) (conn net.Conn, err error) {
	s.Conn = ctx.Conn()
	s.bufReader, s.bufWriter = io.Pipe()
	if s.timeout == 0 {
		s.timeout = time.Second * 5
	}
	if s.interval == 0 {
		s.interval = time.Second * 5
	}
	s.readErrChn = make(chan error, 1)
	s.writeErrChn = make(chan error, 1)
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	go s.backgroundRead()
	go s.heartbeat()
	return next.Call(ctx.SetConn(s))
}

func NewHeartbeatCodec() *HeartbeatCodec {
	return &HeartbeatCodec{}
}
