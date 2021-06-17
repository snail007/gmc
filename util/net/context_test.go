// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gnet

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	_, p, _ := net.SplitHostPort(l.Addr().String())
	c, _ := net.Dial("tcp", "127.0.0.1:"+p)
	conn := NewConn(c)
	conn.SetTimeout(time.Second)
	conn.Initialize()
	conn.Write([]byte("hello"))
	ctx := conn.ctx
	l.Close()
	c.Close()
	l0:=NewListener(l)
	ctx.SetData("test", "abc")
	assert.Exactly(t, c, ctx.Conn())
	assert.Equal(t, time.Second, ctx.ReadTimeout())
	assert.Equal(t, time.Second, ctx.WriteTimeout())
	assert.Equal(t, int64(5), ctx.WriteBytes())
	assert.Equal(t, int64(0), ctx.ReadBytes())
	assert.Contains(t, ctx.LocalAddr().String(), "127.0.0.1:")
	assert.Contains(t, ctx.RemoteAddr().String(), "127.0.0.1:")
	assert.Exactly(t, "abc", ctx.Data("test"))
	assert.IsType(t, (*net.TCPConn)(nil), ctx.RawConn())
	assert.Equal(t, l0, l0.Ctx().Listener())
}
