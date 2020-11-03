package main

import (
	"fmt"
	"github.com/snail007/gmc"
	"github.com/snail007/gmc/demos/website/initialize"
)

var (
	app *gmc.APP
)

func main() {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("main exited, %s\n", e)
		}
		if app != nil && app.Logger() != nil && app.Logger().Async() {
			app.Logger().WaitAsyncDone()
		}
	}()

	// 1. create an app to run.
	app = gmc.New.App().SetConfigFile("../../app/app.toml")

	// 2. add a http server service to app.
	app.AddService(gmc.ServiceItem{
		Service: gmc.New.HTTPServer(),
		AfterInit: func(s *gmc.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = initialize.Initialize(s.Service.(*gmc.HTTPServer))
			return
		},
	})

	// 3. run the app
	err := gmc.StackE(app.Run())
	app.Logger().Panic(err)
}
