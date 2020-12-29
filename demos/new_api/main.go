package main

import (
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	gctx "github.com/snail007/gmc/module/ctx"
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
	ctx := gctx.NewCtx()
	ctx.SetConfig(cfg)
	api, err := gmc.New.APIServerDefault(ctx)
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
	app.AddService(gcore.ServiceItem{
		Service: api.(gcore.Service),
	})
	// 6. run app
	if e := gmc.Err.Stack(app.Run()); e != "" {
		fmt.Println(e)
	}
}
