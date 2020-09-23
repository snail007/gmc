package main

import (
	gmcapp "github.com/snail007/gmc/app"
	"github.com/snail007/gmc/demos/base/initialize"
	"github.com/snail007/gmc/demos/base/router"
	httpserver "github.com/snail007/gmc/http/server"
)

func main() {
	//1.create an app to run
	app := gmcapp.New()

	//2.set app main config file
	app.SetMainConfigFile("../../app/app.toml")

	//3.parse app config file
	err := app.ParseConfig()
	if err != nil {
		panic(err)
	}

	//4.add service to app
	app.AddService(gmcapp.ServiceItem{
		Service: httpserver.New(), //create a http server
		AfterInit: func(s *gmcapp.ServiceItem) (err error) {
			server := s.Service.(*httpserver.HTTPServer)
			//1.do something after http server inited
			initialize.Initialize(server)
			//2.configuration your routers
			router.InitRouter(server)
			return
		},
	})

	//5.run the app
	app.Run()
}
