package gnet

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConn(t *testing.T) {
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	go func() {
		c, _ := l.Accept()
		codec := NewHeartbeatCodec()
		conn := NewConn(c).AddCodec(codec)
		conn.SetTimeout(time.Second * 10)
		conn.Initialize()
		assert.Equal(t, time.Second*5, codec.Interval())
		assert.Equal(t, time.Second*5, codec.Timeout())
		go func() {
			for {
				_, err := conn.Write([]byte("hello from server"))
				if err != nil {
					fmt.Println("server write error", err)
					return
				}
				time.Sleep(time.Millisecond * 100)
			}
		}()
		go func() {
			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println("server read error", err)
					return
				}
				fmt.Printf("%s %d\n", string(buf[:n]), n)
			}
		}()
		time.AfterFunc(time.Second*13, func() {
			conn.Close()
		})
		select {}
	}()
	time.Sleep(time.Second * 2)
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	codec := NewHeartbeatCodec().SetTimeout(time.Second).SetInterval(time.Millisecond * 100)
	assert.Equal(t, time.Millisecond*100, codec.Interval())
	assert.Equal(t, time.Second, codec.Timeout())
	conn := NewConn(c).AddCodec(codec)
	conn.SetTimeout(time.Second * 10)
	assert.Equal(t, time.Second*10, conn.ReadTimeout())
	assert.Equal(t, time.Second*10, conn.WriteTimeout())
	conn.Initialize()
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("client read error", err)
				return
			}
			fmt.Printf("%s %d\n", string(buf[:n]), n)
		}
	}()
	go func() {
		for {
			_, err := conn.Write([]byte("hello from client"))
			if err != nil {
				fmt.Println("client write error", err)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	time.Sleep(time.Second * 15)
}
