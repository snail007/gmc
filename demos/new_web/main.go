package main

import (
	"github.com/snail007/gmc"
)

func main() {

	// 1. create an default app to run.
	app := gmc.New.AppDefault()

	// 2. add a http server service to app.
	app.AddService(gmc.ServiceItem{
		Service: gmc.New.HTTPServer(),
		AfterInit: func(s *gmc.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = Initialize(s.Service.(*gmc.HTTPServer))
			return
		},
	})

	// 3. run the app
	app.Logger().Panic(gmc.StackE(app.Run()))
}
