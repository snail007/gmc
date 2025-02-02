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
	gbytes "github.com/snail007/gmc/util/bytes"
	"github.com/snail007/gmc/util/gpool"
	"io"
	"net"
	"sync"
	"time"
)

const (
	heartbeatCodecMsgFlag uint16 = 0x6699
)

var (
	HeartbeatCodecMsgBufferSize = 32 * 1024
	heartbeatMsgBufferSize      = 2 + 4 //uint16 + uint32
	msgBufPool                  = gbytes.GetPoolCap(0, HeartbeatCodecMsgBufferSize)
	hbBufPool                   = gbytes.GetPoolCap(0, heartbeatMsgBufferSize)
)

type heartbeatCodecMsg struct {
	flag   uint16
	len    uint32
	data   []byte
	msgBuf []byte
	hbBuf  []byte
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
	var buf *bytes.Buffer
	if len(h.data) == 0 {
		h.hbBuf = hbBufPool.Get().([]byte)
		buf = bytes.NewBuffer(h.hbBuf)
	} else {
		h.msgBuf = msgBufPool.Get().([]byte)
		buf = bytes.NewBuffer(h.msgBuf)
	}
	//buf := bytes.NewBuffer(make([]byte, 0, 32*1024))
	//buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, h.flag)
	binary.Write(buf, binary.LittleEndian, uint32(len(h.data)))
	if h.len > 0 {
		binary.Write(buf, binary.LittleEndian, h.data)
	}
	return buf.Bytes()
}

// PutBackBytesBuf should be called after Bytes()
func (h *heartbeatCodecMsg) PutBackBytesBuf() {
	if h.hbBuf != nil {
		//reset slice
		h.hbBuf = h.hbBuf[:0]
		hbBufPool.Put(h.hbBuf)
	}
	if h.msgBuf != nil {
		//reset slice
		h.msgBuf = h.msgBuf[:0]
		msgBufPool.Put(h.msgBuf)
	}
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
	readPool    *gpool.Pool
	writePool   *gpool.Pool
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
	p := gpool.New(1)
retry:
	for {
		p.Submit(func() {
			s.SetWriteDeadline(time.Now().Add(s.timeout))
			_, err := s.Write(nil)
			s.SetWriteDeadline(time.Time{})
			done <- err
		})
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
			err = fmt.Errorf("heartbeat error: %s", err)
			s.sendWriteErr(err)
			s.sendReadErr(err)
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
	p := gpool.New(1)
retry:
	for {
		p.Submit(func() {
			out <- msg.ReadFrom(s.Conn)
		})
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
			err = fmt.Errorf("heartbeat error: %s", err)
			s.sendReadErr(err)
			s.sendWriteErr(err)
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

func (s *HeartbeatCodec) SetConn(c net.Conn) Codec {
	s.Conn = c
	return s
}

func (s *HeartbeatCodec) Read(b []byte) (n int, err error) {
	done := make(chan bool)
	s.readPool.Submit(func() {
		defer close(done)
		n, err = s.bufReader.Read(b)
		_ = err
	})
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
	s.writePool.Submit(func() {
		defer func() {
			close(done)
			//put msg buffer back to pool
			msg.PutBackBytesBuf()
		}()
		n, err = s.Conn.Write(msg.Bytes())
		if err == nil {
			// decrease the n with msg header length 6 bytes
			n = n - 6
		}
	})
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

func (s *HeartbeatCodec) Initialize(ctx Context) (err error) {
	s.bufReader, s.bufWriter = io.Pipe()
	if s.timeout == 0 {
		s.timeout = time.Second * 5
	}
	if s.interval == 0 {
		s.interval = time.Second * 5
	}
	s.readErrChn = make(chan error, 1)
	s.writeErrChn = make(chan error, 1)
	s.readPool = gpool.New(2)
	s.writePool = gpool.New(2)
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	go s.backgroundRead()
	go s.heartbeat()
	return
}

func NewHeartbeatCodec() *HeartbeatCodec {
	return &HeartbeatCodec{}
}
