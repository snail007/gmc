// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmc

import (
	gcore "github.com/snail007/gmc/core"
	gcontroller "github.com/snail007/gmc/http/controller"
	ghttpserver "github.com/snail007/gmc/http/server"
	// using gmc basic: cache, db, web
	_ "github.com/snail007/gmc/using/cache"
	_ "github.com/snail007/gmc/using/db"
	_ "github.com/snail007/gmc/using/web"
	"net/http"
	"sync"
)

type (
	// Controller Alias of type gcontroller.Controller
	Controller = gcontroller.Controller
	// HTTPServer Alias of type ghttpserver.HTTPServer
	HTTPServer = ghttpserver.HTTPServer
	// APIServer Alias of type ghttpserver.APIServer
	APIServer = ghttpserver.APIServer
	// P Alias of type gcore.Params
	P = gcore.Params
	// W Alias of type http.ResponseWriter
	W = http.ResponseWriter
	// R Alias of type *http.Request
	R = *http.Request
	// C Alias of type gcore.Ctx
	C = gcore.Ctx
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		initialize()
	})
}

func initialize() {
	// init others stuff
	initHelper()
}
