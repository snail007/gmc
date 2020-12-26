package gmc

import (
	gcore "github.com/snail007/gmc/core"
	gcontroller "github.com/snail007/gmc/http/controller"
	ghttpserver "github.com/snail007/gmc/http/server"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
	gview "github.com/snail007/gmc/http/view"
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

func init() {
	providers := gcore.Providers

	providers.RegisterSession("", func(ctx gcore.Ctx) (gcore.Session, error) {
		return gsession.NewSession(), nil
	})

	providers.RegisterSessionStorage("", func(ctx gcore.Ctx) (gcore.SessionStorage, error) {
		return gsession.Init(ctx.Config())
	})

	providers.RegisterView("", func(ctx gcore.Ctx) (gcore.View, error) {
		return gview.New(ctx.Response(), ctx.WebServer().Tpl()), nil
	})

	providers.RegisterTemplate("", func(ctx gcore.Ctx) (gcore.Template, error) {
		return gtemplate.Init(ctx.Config())
	})
}
