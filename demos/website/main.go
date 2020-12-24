package main

import (
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/demos/website/initialize"
	gerr "github.com/snail007/gmc/module/error"
	"runtime/debug"
)

var (
	app gcore.GMCApp
)

func main() {
	defer gcore.Recover(func(e interface{}) {
		fmt.Printf("main exited, %s\n", gerr.Wrap(e))
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
	app.AddService(gmc.ServiceItem{
		Service: gmc.New.HTTPServer(),
		AfterInit: func(s *gcore.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = initialize.Initialize(s.Service.(*gmc.HTTPServer))
			return
		},
	})

	// 3. run the app
	if e := gmc.StackE(app.Run());e!=""{
		panic(e)
	}
}
