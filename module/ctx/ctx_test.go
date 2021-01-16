package gctx

import (
	"bytes"
	gcore "github.com/snail007/gmc/core"
	ghttputil "github.com/snail007/gmc/internal/util/http"
	assert2 "github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
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

func TestCtx_GetXXX(t *testing.T) {
	assert := assert2.New(t)
	for _, v := range []struct {
		ctx      *Ctx
		setter   func(*Ctx)
		getter   func(*Ctx) interface{}
		excepted func(*Ctx) interface{}
		testTag  string
	}{
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetAPIServer((gcore.APIServer)(nil))
		}, func(ctx *Ctx) interface{} {
			return ctx.APIServer()
		}, func(ctx *Ctx) interface{} {
			return ctx.apiServer
		}, "APIServer"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetApp((gcore.App)(nil))
		}, func(ctx *Ctx) interface{} {
			return ctx.App()
		}, func(ctx *Ctx) interface{} {
			return ctx.app
		}, "App"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetI18n((gcore.I18n)(nil))
		}, func(ctx *Ctx) interface{} {
			return ctx.I18n()
		}, func(ctx *Ctx) interface{} {
			return ctx.i18n
		}, "I18n"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetConn((net.Conn)(nil))
		}, func(ctx *Ctx) interface{} {
			return ctx.Conn()
		}, func(ctx *Ctx) interface{} {
			return ctx.conn
		}, "Conn"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetRemoteAddr("a")
		}, func(ctx *Ctx) interface{} {
			return ctx.RemoteAddr()
		}, func(ctx *Ctx) interface{} {
			return ctx.remoteAddr
		}, "RemoteAddr"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetTemplate((gcore.Template)(nil))
		}, func(ctx *Ctx) interface{} {
			return ctx.Template()
		}, func(ctx *Ctx) interface{} {
			return ctx.template
		}, "Template"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetConfig((gcore.Config)(nil))
		}, func(ctx *Ctx) interface{} {
			return ctx.Config() != nil
		}, func(ctx *Ctx) interface{} {
			return true
		}, "Config"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetLogger((gcore.Logger)(nil))
		}, func(ctx *Ctx) interface{} {
			return ctx.Logger()
		}, func(ctx *Ctx) interface{} {
			return ctx.logger
		}, "Logger"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetWebServer((gcore.HTTPServer)(nil))
		}, func(ctx *Ctx) interface{} {
			return ctx.WebServer()
		}, func(ctx *Ctx) interface{} {
			return ctx.webServer
		}, "WebServer"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetLocalAddr("a")
		}, func(ctx *Ctx) interface{} {
			return ctx.LocalAddr()
		}, func(ctx *Ctx) interface{} {
			return ctx.localAddr
		}, "LocalAddr"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetParam(gcore.Params{})
		}, func(ctx *Ctx) interface{} {
			return ctx.Param()
		}, func(ctx *Ctx) interface{} {
			return ctx.param
		}, "Param"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetTimeUsed(111)
		}, func(ctx *Ctx) interface{} {
			return ctx.TimeUsed()
		}, func(ctx *Ctx) interface{} {
			return ctx.timeUsed
		}, "TimeUsed"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.WriteE("a")
		}, func(ctx *Ctx) interface{} {
			return ctx.response.(*httptest.ResponseRecorder).Code == http.StatusInternalServerError &&
				ctx.response.(*httptest.ResponseRecorder).Body.String() == "a"
		}, func(ctx *Ctx) interface{} {
			return true
		}, "WriteE"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.WriteHeader(500)
		}, func(ctx *Ctx) interface{} {
			return ctx.response.(*httptest.ResponseRecorder).Code == http.StatusInternalServerError
		}, func(ctx *Ctx) interface{} {
			return true
		}, "WriteHeader"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.response = ghttputil.NewResponseWriter(ctx.response)
			ctx.Status(http.StatusInternalServerError)
		}, func(ctx *Ctx) interface{} {
			return ctx.StatusCode()
		}, func(ctx *Ctx) interface{} {
			return http.StatusInternalServerError
		}, "StatusCode"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.response = ghttputil.NewResponseWriter(ctx.response)
			ctx.Write("aaa")
		}, func(ctx *Ctx) interface{} {
			return ctx.WriteCount()
		}, func(ctx *Ctx) interface{} {
			return int64(3)
		}, "WriteCount"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsGET() }, func(ctx *Ctx) interface{} { return true }, "IsGET"},
		{mockCtx("POST", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsPOST() }, func(ctx *Ctx) interface{} { return true }, "IsPOST"},
		{mockCtx("DELETE", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsDELETE() }, func(ctx *Ctx) interface{} { return true }, "IsDELETE"},
		{mockCtx("PUT", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsPUT() }, func(ctx *Ctx) interface{} { return true }, "IsPUT"},
		{mockCtx("PATCH", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsPATCH() }, func(ctx *Ctx) interface{} { return true }, "IsPATCH"},
		{mockCtx("OPTIONS", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsOPTIONS() }, func(ctx *Ctx) interface{} { return true }, "IsOPTIONS"},
		{mockCtx("HEAD", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsHEAD() }, func(ctx *Ctx) interface{} { return true }, "IsHEAD"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.request.Header.Set("X-Requested-With", "XMLHttpRequest")
		}, func(ctx *Ctx) interface{} { return ctx.IsAJAX() }, func(ctx *Ctx) interface{} { return true }, "IsAJAX"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.request.Header.Set("Connection", "upgrade")
			ctx.request.Header.Set("Upgrade", "websocket")
		}, func(ctx *Ctx) interface{} { return ctx.IsWebsocket() }, func(ctx *Ctx) interface{} { return true }, "IsWebsocket"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsAJAX() }, func(ctx *Ctx) interface{} { return false }, "IsAJAX"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {}, func(ctx *Ctx) interface{} { return ctx.IsWebsocket() }, func(ctx *Ctx) interface{} { return false }, "IsWebsocket"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.request.Header.Set("X-Forwarded-For", "1.1.1.1")
		}, func(ctx *Ctx) interface{} {
			return ctx.ClientIP()
		}, func(ctx *Ctx) interface{} {
			return "1.1.1.1"
		}, "ClientIP-X-Forwarded-For"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.request.Header.Set("X-Real-IP", "1.1.1.1")
		}, func(ctx *Ctx) interface{} {
			return ctx.ClientIP()
		}, func(ctx *Ctx) interface{} {
			a := ctx.Clone().(*Ctx)
			c := a.CloneWithHTTP(a.response, a.request, a.param).(*Ctx)
			*ctx = *c
			return "1.1.1.1"
		}, "ClientIP-X-Real-IP"},
		{mockCtx("GET", "/", ""), func(ctx *Ctx) {
			ctx.SetCookie("test", "a", 0, "", "", true, true)
		}, func(ctx *Ctx) interface{} {
			return ctx.response.(*httptest.ResponseRecorder).Header().Get("Set-Cookie")
		}, func(ctx *Ctx) interface{} {
			return "test=a; Path=/; HttpOnly; Secure"
		}, "SetCookie"},
	} {
		v.setter(v.ctx)
		assert.Equal(v.excepted(v.ctx), v.getter(v.ctx), v.testTag)
	}
}

func TestNewCtx(t *testing.T) {
	assert := assert2.New(t)
	assert.Implements((*gcore.Ctx)(nil), NewCtx())
	assert.Implements((*gcore.Ctx)(nil), NewCtxFromConfig((gcore.Config)(nil)))
	c, e := NewCtxFromConfigFile("../app/app.toml")
	assert.Nil(e)
	assert.Implements((*gcore.Ctx)(nil), c)
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
