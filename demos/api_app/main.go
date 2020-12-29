package main

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gmchttp "github.com/snail007/gmc/http"
	gctx "github.com/snail007/gmc/module/ctx"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/snail007/gmc"
)

func main() {
	api, _ := gmc.New.APIServer(gctx.NewCtx(), ":8030")
	// add a middleware typed 1 to filter all request registered in router,
	// exclude 404 requests.
	api.AddMiddleware1(func(c gmc.C, s gcore.APIServer) (isStop bool) {
		s.Logger().Infof("before request %s", c.Request().RequestURI)
		return false
	})
	// add a middleware typed 2 to logging every request registered in router,
	// exclude 404 requests.
	api.AddMiddleware2(func(c gmc.C, s gcore.APIServer) (isStop bool) {
		s.Logger().Infof("after request %s %d %d %s", c.Request().Method, c.StatusCode(), c.WriteCount(), c.Request().RequestURI)
		return false
	})
	// sets a function to handle 404 requests.
	api.Handle404(func(c gmc.C) {
		c.Write("404")
	})
	// sets a function to handle panic error.
	api.Handle500(func(c gmc.C, err interface{}) {
		c.WriteHeader(http.StatusInternalServerError)
		c.Write("panic error : ", err)
	})
	// sets an api in url path: /hello
	// add more api , just call api.API() repeatedly
	api.API("/hello", func(c gmc.C) {
		api.Logger().Infof("request %s", c.Request().RequestURI)
		c.Write("hello world!")
	})
	api.API("/hi", func(c gmc.C) {
		c.Write("hi!")
	})
	// trigger a panic error
	api.API("/error", func(c gmc.C) {
		a := 0
		a /= a
	})
	api.API("/version", func(c gmc.C) {
		c.Write(1.1)
	})
	api.API("/sleep", func(c gmc.C) {
		time.Sleep(time.Second * 10)
		c.Write("reload")
	})
	// access ctx
	// http://foo.com/ctxfoo
	api.Router().HandlerFunc("GET", "/ctx:name", func(w http.ResponseWriter, r *http.Request) {
		ctx := gmchttp.GetCtx(w)
		ctx.Write(ctx.Param().ByName("name"), " ", ctx.Conn().LocalAddr().String(), " ", ctx.Conn().RemoteAddr().String())
	})
	// http://foo.com/2ctxfoo
	api.Router().Handle("GET", "/2ctx:name", func(w http.ResponseWriter, r *http.Request, ps gcore.Params) {
		ctx := gmchttp.GetCtx(w)
		ctx.Write(ctx.Param().ByName("name"), " ", ctx.Conn().LocalAddr().String(), " ", ctx.Conn().RemoteAddr().String())
	})
	// routing by group is supported
	group1 := api.Group("/v1")
	group1.API("/time", func(c gmc.C) {
		c.Write(time.Now().In(time.Local).Format("2006-01-02 15:04:05"))
	})

	app := gmc.New.App()

	app.AddService(gcore.ServiceItem{
		Service: api.(gcore.Service),
		BeforeInit: func(s gcore.Service, cfg gcore.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})

	// all path in router
	_, port, _ := net.SplitHostPort(api.Address())
	fmt.Println("please visit:")
	for path := range api.Router().RouteTable() {
		if strings.Contains(path, "*") {
			continue
		}
		if strings.Contains(path, ":type") {
			path = "/hi.json"
		}
		fmt.Println("http://127.0.0.1:" + port + path)
	}
	if e := gmc.Err.Stack(app.Run()); e != "" {
		app.Logger().Panic(e)
	}
}
