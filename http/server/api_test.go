// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttpserver

import (
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/http/template/testdata"
	"net"
	"testing"

	ghttputil "github.com/snail007/gmc/internal/util/http"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIServer(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	assert.NotNil(api.server)
	assert.Equal(len(api.address), 1)
}
func TestBefore(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.AddMiddleware0(func(c gcore.Ctx) bool {
		c.Write("okay")
		return true
	})
	api.API("/hello", func(c gcore.Ctx) {
		c.Write("a")
	})
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("okay", data)
}

func TestAPI(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.API("/hello", func(c gcore.Ctx) {
		c.Write("a")
	})
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("a", data)
}

func TestAfter(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.AddMiddleware2(func(c gcore.Ctx) bool {
		c.Write("okay")
		return false
	})
	api.API("/hello", func(c gcore.Ctx) {
		c.Write("a")
	})
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("aokay", data)
}
func TestStop(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.API("/hello", func(c gcore.Ctx) {
		c.Write("a")
		ghttputil.Stop(c.Response(), "c")
		c.Write("b")
	})
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("ac", data)
}
func TestHandle404(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.SetNotFoundHandler(func(c gcore.Ctx) {
		c.Write("404")
	})
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("404", data)
}
func TestHandle404_1(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("Page not found", data)
}
func TestHandle500(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.SetErrorHandler(func(c gcore.Ctx, err interface{}) {
		c.Write("500")
	})
	api.API("/hello", func(c gcore.Ctx) {
		a := 0
		a /= a
	})
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("500", data)
}

func TestHandle500_1(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.ShowErrorStack(false)
	api.API("/hello", func(c gcore.Ctx) {
		a := 0
		a /= a
	})
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("Internal Server Error", data)
}

type MyListener struct {
	net.Listener
}

func (s *MyListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 123,
		Zone: "",
	}
}

type MyListener1 struct {
	net.Listener
}

func (s *MyListener1) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 123,
		Zone: "",
	}
}

func TestAPIServer_createListener(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.InjectListeners([]net.Listener{&MyListener{}})
	e := api.createListener()
	assert.Nil(e)
	assert.IsType((*MyListener)(nil), api.listener)
	api.InjectListeners([]net.Listener{nil})

	api.SetListenerFactory(func(_ string) (net.Listener, error) {
		return &MyListener1{}, nil
	})
	assert.NotNil(api.ListenerFactory())
	e = api.createListener()
	assert.Nil(e)
	assert.IsType((*MyListener1)(nil), api.listener)

	api.InjectListeners([]net.Listener{&MyListener{}})
	e = api.createListener()
	assert.Nil(e)
	assert.IsType((*MyListener)(nil), api.listener)
}

func TestAPIServer_ServeFiles(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(gcore.ProviderCtx()(), ":")
	api.ServeEmbedFS(testdata.TplFS, "/tpls/")
	api.ServeFiles("tests", "/files/")

	w, r := mockRequest("/tpls/f/f.txt")
	api.ServeHTTP(w, r)
	str, _ := result(w)
	assert.Equal("f", str)

	w, r = mockRequest("/files/d.txt")
	api.ServeHTTP(w, r)
	str, _ = result(w)
	assert.Equal("d", str)
}
