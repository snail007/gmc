// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttpserver_test

import (
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	"net"
	"net/http"
	"strings"
	"time"
)

func apiServer() (addr string, err error) {
	api, _ := gmc.New.APIServer(gmc.New.Ctx(), ":")
	api.Ext(".html")
	// add a middleware typed 1 to filter all request registered in router,
	// exclude 404 requests.
	api.AddMiddleware1(func(c gmc.C) (isStop bool) {
		c.Logger().Infof("before request %s", c.Request().RequestURI)
		return false
	})
	// add a middleware typed 2 to logging every request registered in router,
	// exclude 404 requests.
	api.AddMiddleware2(func(c gmc.C) (isStop bool) {
		c.Logger().Infof("after request %s %d %d %s %s",
			c.Request().Method,
			c.StatusCode(),
			c.WriteCount(),
			c.TimeUsed(),
			c.Request().RequestURI)
		return false
	})
	// sets a function to handle 404 requests.
	api.SetNotFoundHandler(func(c gmc.C) {
		c.Write("404")
	})
	// sets a function to handle panic error.
	api.SetErrorHandler(func(c gmc.C, err interface{}) {
		c.WriteHeader(http.StatusInternalServerError)
		c.Write("panic error : ", gmc.Err.Stack(err))
	})
	// sets an api in url path: /hello
	// add more api , just call api.API() repeatedly
	api.API("/hello", func(c gmc.C) {
		api.Logger().Infof("request %s", c.Request().RequestURI)
		c.Write("hello world!")
	})
	api.API("/hi.:type", func(c gmc.C) {
		c.Write("hi!" + c.Param().ByName("type"))
	}, "")
	// trigger a panic error
	api.API("/error", func(c gmc.C) {
		a := 0
		a /= a
	})
	// access under layer conn
	api.API("/conn", func(c gmc.C) {
		c.Write(c.Conn().RemoteAddr().String(), c.Conn().LocalAddr().String())
	})
	// access ctx
	// http://foo.com/ctxfoo
	api.Router().HandlerFunc("GET", "/ctx:name", func(w http.ResponseWriter, r *http.Request) {
		ctx := gcore.GetCtx(w)
		ctx.Write(ctx.Param().ByName("name"), " ", ctx.Conn().LocalAddr().String(), " ", ctx.Conn().RemoteAddr().String())
	})
	// http://foo.com/2ctxfoo
	api.Router().Handle("GET", "/2ctx:name", func(w http.ResponseWriter, r *http.Request, ps gcore.Params) {
		ctx := gcore.GetCtx(w)
		ctx.Write(ctx.Param().ByName("name"), " ", ctx.Conn().LocalAddr().String(), " ", ctx.Conn().RemoteAddr().String())
	})

	// routing by group is supported
	group1 := api.Group("/v1")
	group1.Ext(".json")
	group1.API("/time", func(c gmc.C) {
		c.Write(time.Now().In(time.Local).Format("2006-01-02 15:04:05"))
	})
	api.PrintRouteTable(nil)
	// start the api server
	err = api.Run()
	if err != nil {
		panic(err)
	}
	addr = api.Address()
	_, p, _ := net.SplitHostPort(addr)
	addr = "http://" + net.JoinHostPort("127.0.0.1", p)
	// all path in router
	_, port, _ := net.SplitHostPort(api.Listener().Addr().String())
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
	client := http.Client{Timeout: time.Millisecond * 100}
	i := 0
	for {
		if i > 10 {
			return "", fmt.Errorf("start testing api server fail")
		}
		_, err = client.Get(addr + "/hi.json")
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
		i++
	}
	return
}
