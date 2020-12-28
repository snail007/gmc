package main

import (
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
 	gctx "github.com/snail007/gmc/module/ctx"
	gerror "github.com/snail007/gmc/module/error"
	gmap "github.com/snail007/gmc/util/map"
)

func main() {

	api := gmc.New.APIServer(gctx.NewCtx(), ":7082")
	api.API("/", func(c gmc.C) {
		c.Write(gmap.M{
			"code":    0,
			"message": "Hello GMC!",
			"data":    nil,
		})
	})

	app := gmc.New.App()
	app.AddService(gcore.ServiceItem{
		Service: api,
		BeforeInit: func(s gcore.Service, cfg gcore.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})

	if e := gerror.Stack(app.Run()); e != "" {
		app.Logger().Panic(e)
	}
}
