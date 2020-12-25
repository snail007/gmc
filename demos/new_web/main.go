package main

import (
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	ghttpserver "github.com/snail007/gmc/http/server"
	gerror "github.com/snail007/gmc/module/error"
	"mygmcweb/initialize"
)

func main() {

	// 1. create an default app to run.
	app := gmc.New.AppDefault()

	// 2. add a http server service to app.
	app.AddService(gcore.ServiceItem{
		Service: ghttpserver.NewHTTPServer(app.Ctx()),
		AfterInit: func(s *gcore.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = initialize.Initialize(s.Service.(*ghttpserver.HTTPServer))
			return
		},
	})

	// 3. run the app
	if e := gerror.Stack(app.Run()); e != "" {
		app.Logger().Panic(e)
	}
}
