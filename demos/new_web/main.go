package main

import (
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	ghttpserver "github.com/snail007/gmc/http/server"
	"mygmcweb/initialize"
)

func main() {

	// 1. create an default app to run.
	app := gmc.New.AppDefault()

	// 2. add a http server service to app.
	app.AddService(gmc.ServiceItem{
		Service: ghttpserver.New(),
		AfterInit: func(s *gcore.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = initialize.Initialize(s.Service.(*ghttpserver.HTTPServer))
			return
		},
	})

	// 3. run the app
	if e := gmc.StackE(app.Run());e!=""{
		app.Logger().Panic(e)
	}
}
