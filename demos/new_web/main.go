// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package main

import (
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	"mygmcweb/initialize"
)

func main() {

	// 1. create an default app to run.
	app := gmc.New.AppDefault()

	// 2. add a http server service to app.
	app.AddService(gcore.ServiceItem{
		Service: gmc.New.HTTPServer(app.Ctx()).(gcore.Service),
		AfterInit: func(s *gcore.ServiceItem) (err error) {
			// do some initialize after http server initialized.
			err = initialize.Initialize(s.Service.(gcore.HTTPServer))
			return
		},
	})

	// 3. run the app
	if e := gmc.Err.Stack(app.Run()); e != "" {
		fmt.Println(e)
	}
}
