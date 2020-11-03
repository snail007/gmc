package main

import (
	"github.com/snail007/gmc"
	"mygmcapi/handlers"
)

func main() {
	// 1. create app
	app := gmc.New.App()
	// 2. parse config file
	cfg, err := gmc.New.ConfigFile("conf/app.toml")
	if err != nil {
		app.Logger().Error(err)
	}
	// 3. create api server
	api, err := gmc.New.APIServerDefault(cfg)
	if err != nil {
		app.Logger().Error(err)
	}
	//4. init db, cache, handlers
	// int db
	gmc.DB.Init(cfg)
	// init cache
	gmc.Cache.Init(cfg)
	// init api handlers
	handlers.Init(api)
	// 5. add service
	app.AddService(gmc.ServiceItem{
		Service: api,
	})
	// 6. run app
	e := gmc.StackE(app.Run())
	app.Logger().Panic(e)
}
