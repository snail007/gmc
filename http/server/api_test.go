package gmchttpserver

import (
	"testing"

	gmcrouter "github.com/snail007/gmc/http/router"
	gmchttputil "github.com/snail007/gmc/util/httputil"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIServer(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(":")
	assert.NotNil(api.server)
	assert.Equal(len(api.address), 1)
}
func TestBefore(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(":")
	api.Before(func(c *gmcrouter.Ctx) bool {
		c.Write("okay")
		return false
	})
	api.API("/hello", func(c *gmcrouter.Ctx) {
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
	api := NewAPIServer(":")
	api.API("/hello", func(c *gmcrouter.Ctx) {
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
	api := NewAPIServer(":")
	api.After(func(c *gmcrouter.Ctx, isPanic bool) {
		c.Write("okay")
		return
	})
	api.API("/hello", func(c *gmcrouter.Ctx) {
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
	api := NewAPIServer(":")
	api.API("/hello", func(c *gmcrouter.Ctx) {
		c.Write("a")
		gmchttputil.Stop(c.Response, "c")
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
	api := NewAPIServer(":")
	api.Handle404(func(c *gmcrouter.Ctx) {
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
	api := NewAPIServer(":")
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("Page not found", data)
}
func TestHandle500(t *testing.T) {
	assert := assert.New(t)
	api := NewAPIServer(":")
	api.Handle500(func(c *gmcrouter.Ctx, err interface{}) {
		c.Write("500")
	})
	api.API("/hello", func(c *gmcrouter.Ctx) {
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
	api := NewAPIServer(":").ShowErrorStack(false)
	api.API("/hello", func(c *gmcrouter.Ctx) {
		a := 0
		a /= a
	})
	w, r := mockRequest("/hello")
	api.ServeHTTP(w, r)
	data, _, err := mockResponse(w)
	assert.Nil(err)
	assert.Equal("Internal Server Error", data)
}
