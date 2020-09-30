package main

import (
	"github.com/snail007/gmc"
	"github.com/snail007/gmc/demos/website/initialize"
	"github.com/snail007/gmc/demos/website/router"
)

func main() {
	//1. create an app to run
	app := gmc.NewAPP()

	//2. set app main config file
	app.SetMainConfigFile("../../app/app.toml")

	//3. when config file parsed, then initialize global database and cache objects in main config file app.toml, if have.
	app.AfterParse(func(_ *gmc.Config) (err error) {
		err = gmc.InitDB(app.Config())
		return
	})
	app.AfterParse(func(_ *gmc.Config) (err error) {
		err = gmc.InitCache(app.Config())
		return
	})

	//4. parse app config file
	err := app.ParseConfig()
	if err != nil {
		panic(err)
	}

	//5.add service to app
	app.AddService(gmc.ServiceItem{
		Service: gmc.NewHTTPServer(), //create a http server
		AfterInit: func(s *gmc.ServiceItem) (err error) {
			server := s.Service.(*gmc.HTTPServer)
			//1.do something after http server inited
			initialize.Initialize(server)
			//2.configuration your routers
			router.InitRouter(server)
			return
		},
	})
	//and add more service to app
	app.AddService(gmc.ServiceItem{
		Service: gmc.NewHTTPServer(), //create a http server
		BeforeInit: func(cfg *gmc.Config) (err error) {
			//change port listen on
			cfg.Set("httpserver.listen", ":")
			return
		},
		AfterInit: func(s *gmc.ServiceItem) (err error) {
			server := s.Service.(*gmc.HTTPServer)
			//1.do something after http server inited
			initialize.Initialize(server)
			//2.configuration your routers
			router.InitRouter(server)
			return
		},
	})

	//6.run the app
	app.Logger().Panic(app.Run())
}
