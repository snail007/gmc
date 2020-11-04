package main

import (
	"github.com/snail007/gmc"
)

func main() {

	api := gmc.New.APIServer(":7082")
	api.API("/", func(c gmc.C) {
 		c.Write(gmc.M{
 			"code":0,
 			"message":"Hello GMC!",
 			"data":nil,
		})
	})

	app := gmc.New.App()
	app.AddService(gmc.ServiceItem{
		Service: api,
		BeforeInit: func(s gmc.Service, cfg *gmc.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})

	if e := gmc.StackE(app.Run());e!=""{
		app.Logger().Panic(e)
	}
}
