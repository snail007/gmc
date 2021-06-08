package gnet

import (
	"net"
	"time"
)

type Codec interface {
	net.Conn
	Initialize(net.Conn) error
}

type Conn struct {
	net.Conn
	codec        []Codec
	readTimeout  time.Duration
	writeTimeout time.Duration
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
	for _, c := range s.codec {
		err = c.Initialize(s.Conn)
		if err != nil {
			return
		}
		s.Conn = c
	}
	return nil
}
func (s *Conn) Read(b []byte) (n int, err error) {
	if s.readTimeout > 0 {
		s.Conn.SetReadDeadline(time.Now().Add(s.readTimeout))
		defer s.Conn.SetReadDeadline(time.Time{})
	}
	return s.Conn.Read(b)
}

func (s *Conn) Write(b []byte) (n int, err error) {
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
	return &Conn{Conn: conn}
}
