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
	app.SetConfigFile("../../app/app.toml")

	//3. when config file parsed success, we initialize database and cache only once.
	app.OnRunOnce(func(cfg *gmc.Config) (err error) {
		err = gmc.InitDB(cfg)
		return
	})
	app.OnRunOnce(func(cfg *gmc.Config) (err error) {
		err = gmc.InitCache(cfg)
		return
	})

	//4.add a service to app
	// and you can call AddService more to add service to app
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

	//5.run the app
	app.Logger().Panic(gmc.StackE(app.Run()))

}
