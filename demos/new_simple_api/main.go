package main

import (
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
	gconfig "github.com/snail007/gmc/util/config"
	gmap "github.com/snail007/gmc/util/map"
)

func main() {

	api := gmc.New.APIServer(":7082")
	api.API("/", func(c gmc.C) {
 		c.Write(gmap.M{
 			"code":0,
 			"message":"Hello GMC!",
 			"data":nil,
		})
	})

	app := gmc.New.App()
	app.AddService(gcore.ServiceItem{
		Service: api,
		BeforeInit: func(s gcore.Service, cfg *gconfig.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})

	if e := gerror.Stack(app.Run());e!=""{
		app.Logger().Panic(e)
	}
}
