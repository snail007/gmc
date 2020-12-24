package gmc

import (
	gcore "github.com/snail007/gmc/core"
	gcontroller "github.com/snail007/gmc/http/controller"
	ghttpserver "github.com/snail007/gmc/http/server"
	"net/http"
)

type (
	// Alias of type gcontroller.Controller
	Controller = gcontroller.Controller
	// Alias of type ghttpserver.HTTPServer
	HTTPServer = ghttpserver.HTTPServer
	// Alias of type ghttpserver.APIServer
	APIServer = ghttpserver.APIServer
	// Alias of type gcore.Params
	P = gcore.Params
	// Alias of type http.ResponseWriter
	W = http.ResponseWriter
	// Alias of type *http.Request
	R = *http.Request
	// Alias of type gcore.Ctx
	C = gcore.Ctx
)
