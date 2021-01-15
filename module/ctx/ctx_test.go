package gctx

import (
	"bytes"
	gcore "github.com/snail007/gmc/core"
	assert2 "github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestCtx_GET(t *testing.T) {
	assert := assert2.New(t)
	for _, v := range []struct {
		ctx    *Ctx
		key    string
		defalt string
		except string
	}{
		{mockCtx("GET", "/foo?a=1&b=3", ""), "a", "", "1"},
		{mockCtx("GET", "/foo?a=3&a=1", ""), "a", "", "3"},
		{mockCtx("GET", "/foo", ""), "a", "", ""},
		{mockCtx("GET", "/foo", ""), "a", "4", "4"},
	} {
		assert.Equal(v.except, v.ctx.GET(v.key, v.defalt))
	}
}

func TestCtx_POST(t *testing.T) {
	assert := assert2.New(t)
	for _, v := range []struct {
		ctx    *Ctx
		key    string
		defalt string
		except string
	}{
		{mockCtx("POST", "/foo", "a=1&b=3"), "a", "", "1"},
		{mockCtx("POST", "/foo", "a=3&b=1"), "a", "", "3"},
		{mockCtx("POST", "/foo", ""), "a", "", ""},
		{mockCtx("POST", "/foo", ""), "a", "4", "4"},
	} {
		assert.Equal(v.except, v.ctx.POST(v.key, v.defalt))
	}
}

func TestCtx_Cookie(t *testing.T) {
	assert := assert2.New(t)
	for _, v := range []struct {
		ctx    *Ctx
		cookie string
		key    string
		except string
	}{
		{mockCtx("GET", "/foo", ""), "a=1;b=3", "a", "1"},
		{mockCtx("GET", "/foo", ""), "a=3;b=1", "a", "3"},
		{mockCtx("GET", "/foo", ""), "", "a", ""},
	} {
		v.ctx.request.Header.Set("Cookie", v.cookie)
		assert.Equal(v.except, v.ctx.Cookie(v.key))
	}
}

func TestCtx_GETArray(t *testing.T) {
	assert := assert2.New(t)
	for _, v := range []struct {
		ctx    *Ctx
		key    string
		defalt string
		except interface{}
	}{
		{mockCtx("GET", "/foo", ""), "a", "", *new([]string)},
		{mockCtx("GET", "/foo?a=1", ""), "a", "", []string{"1"}},
		{mockCtx("GET", "/foo?a=1&a=3", ""), "a", "", []string{"1", "3"}},
	} {
		if v.defalt != "" {
			assert.Equal(v.except, v.ctx.GETArray(v.key, v.defalt))
		} else {
			assert.Equal(v.except, v.ctx.GETArray(v.key))
		}
	}
}

func TestCtx_JSON(t *testing.T) {
	assert := assert2.New(t)
	for _, v := range []struct {
		ctx        *Ctx
		code       int
		data       interface{}
		exceptCode int
		exceptData string
	}{
		{mockCtx("GET", "/", ""), 200, map[string]string{"a": "b"}, 200, `{"a":"b"}`},
		{mockCtx("GET", "/", ""), 500, []string{"a"}, 500, `["a"]`},
		{mockCtx("GET", "/", ""), 404, "", 404, ""},
	} {
		v.ctx.JSON(v.code, v.data)
		assert.Equal(v.exceptCode, v.ctx.response.(*httptest.ResponseRecorder).Code)
		data, _ := ioutil.ReadAll(v.ctx.response.(*httptest.ResponseRecorder).Body)
		assert.Equal(v.exceptData, string(data))
	}
}

func TestCtx_JSONP(t *testing.T) {
	assert := assert2.New(t)
	for _, v := range []struct {
		ctx        *Ctx
		code       int
		data       interface{}
		exceptCode int
		exceptData string
	}{
		{mockCtx("GET", "/", ""), 200, map[string]string{"a": "b"}, 200, `_jsonp({"a":"b"})`},
		{mockCtx("GET", "/?callback=jsonp", ""), 500, []string{"a"}, 500, `jsonp(["a"])`},
		{mockCtx("GET", "/", ""), 404, "", 404, "_jsonp()"},
	} {
		v.ctx.JSONP(v.code, v.data)
		assert.Equal(v.exceptCode, v.ctx.response.(*httptest.ResponseRecorder).Code)
		data, _ := ioutil.ReadAll(v.ctx.response.(*httptest.ResponseRecorder).Body)
		assert.Equal(v.exceptData, string(data))
	}
}

func TestNewCtx(t *testing.T) {
	assert := assert2.New(t)
	assert.Implements((*gcore.Ctx)(nil), NewCtx())
}

func mockCtx(method, path string, body string) *Ctx {
	r := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	if body != "" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	w := httptest.NewRecorder()
	c := NewCtx()
	c.SetRequest(r)
	c.SetResponse(w)
	return c
}
