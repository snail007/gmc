// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package main

import (
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gmchttp "github.com/snail007/gmc/http"
	gmap "github.com/snail007/gmc/util/map"
	"net/http"
)

func main() {

	api, _ := gmc.New.APIServer(gmc.New.Ctx(), ":7082")

	api.API("/", func(c gmc.C) {
		c.Write(gmap.M{
			"code":    0,
			"message": "Hello GMC!",
			"data":    nil,
		})
	})

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

	app := gmc.New.App()
	app.AddService(gcore.ServiceItem{
		Service: api.(gcore.Service),
		BeforeInit: func(s gcore.Service, cfg gcore.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})

	if e := gmc.Err.Stack(app.Run()); e != "" {
		fmt.Println(e)
	}
}
