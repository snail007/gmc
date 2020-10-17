package main

import (
	"github.com/snail007/gmc"
	"mygmcapi/handlers"
)

func main() {
	app := gmc.New.App()
	_, err := gmc.New.ConfigFile("conf/app.toml")
	if err != nil {
		app.Logger().Fatal(err)
	}
api:=gmc.New.APIServer(":")

	// init api
	handlers.Init(api)

	app.AddService(gmc.ServiceItem{
		Service: api,
		BeforeInit: func(s gmc.Service, cfg *gmc.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})
	app.Logger().Panic(gmc.StackE(app.Run()))
}
