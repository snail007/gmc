package handlers

import (
	"bytes"
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gmchttp "github.com/snail007/gmc/http"
	"net/http"
	"time"
)

func initHanlder(api gcore.APIServer) {

	// URL: http://foo.com/
	api.API("/", func(c gcore.Ctx) {
		var out bytes.Buffer
		out.WriteString("<title>Hello GMC!</title><h1>This is a GMC API Server!</h1>")
		for k := range api.Router().RouteTable() {
			a := fmt.Sprintf("http://%s%s", c.Request().Host, k)
			out.WriteString(fmt.Sprintf("<p><a href=\"%s\" target=\"_blank\">%s</a></p>", a, a))
		}
		out.WriteString("<p><a href=\"https://github.com/snail007/gmc\" target=\"_blank\">View on GitHub</a></p>")
		c.Write(out.Bytes())
	})
	// URL: http://foo.com/version
	api.API("/version", func(c gmc.C) {
		c.Write(1.1)
	})
	// http://foo.com/sleep
	api.API("/sleep", func(c gmc.C) {
		time.Sleep(time.Second * 10)
		c.Write("reload")
	})

	// routing by group is supported
	// http://foo.com/v1/hello
	group0 := api.Group("/v1")
	group0.API("/hello", func(c gmc.C) {
		api.Logger().Infof("request %s", c.Request().RequestURI)
		c.Write("hello world!")
	})
	// http://foo.com/v1/hi
	group0.API("/hi", func(c gmc.C) {
		c.Write("hi!")
	})
	// http://foo.com/v1/error
	// trigger a panic error
	group0.API("/error", func(c gmc.C) {
		a := 0
		a /= a
	})
	// access ctx
	// http://foo.com/ctxfoo
	group0.Router().HandlerFunc("GET", "/ctx:name", func(w http.ResponseWriter, r *http.Request) {
		ctx := gmchttp.GetCtx(w)
		ctx.Write(ctx.Param().ByName("name"), " ", ctx.Conn().LocalAddr().String(), " ", ctx.Conn().RemoteAddr().String())
	})
	// http://foo.com/2ctxfoo
	group0.Router().Handle("GET", "/2ctx:name", func(w http.ResponseWriter, r *http.Request, ps gcore.Params) {
		ctx := gmchttp.GetCtx(w)
		ctx.Write(ctx.Param().ByName("name"), " ", ctx.Conn().LocalAddr().String(), " ", ctx.Conn().RemoteAddr().String())
	})

	// http://foo.com/v2/time
	group1 := api.Group("/v2")
	group1.Ext(".json")
	group1.API("/time", func(c gmc.C) {
		c.Write(time.Now().In(time.Local).Format("2006-01-02 15:04:05"))
	})

}
