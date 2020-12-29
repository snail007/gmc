package main

import (
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/demos/website/initialize"
	"runtime/debug"
)

var (
	app gcore.App
)

func main() {
	defer gmc.Err.Recover(func(e interface{}) {
		fmt.Printf("main exited, %s\n", gmc.Err.Stack(e))
		debug.PrintStack()
	})
	defer func() {
		if app != nil && app.Logger() != nil && app.Logger().Async() {
			app.Logger().WaitAsyncDone()
		}
	}()

	// 1. create an app to run.
	app = gmc.New.App()
	app.SetConfigFile("../../module/app/app.toml")

	// 2. add a http server service to app.
	app.AddService(gcore.ServiceItem{
		Service: gmc.New.HTTPServer(app.Ctx()).(gcore.Service),
		AfterInit: func(s *gcore.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = initialize.Initialize(s.Service.(*gmc.HTTPServer))
			return
		},
	})

	// 3. run the app
	if e := gmc.Err.Stack(app.Run()); e != "" {
		panic(e)
	}
}
