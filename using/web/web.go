// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package web

import (
	gcore "github.com/snail007/gmc/core"
	gcontroller "github.com/snail007/gmc/http/controller"
	gcookie "github.com/snail007/gmc/http/cookie"
	grouter "github.com/snail007/gmc/http/router"
	ghttpserver "github.com/snail007/gmc/http/server"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
	gview "github.com/snail007/gmc/http/view"
	gi18n "github.com/snail007/gmc/module/i18n"
	_ "github.com/snail007/gmc/using/basic"
	gonce "github.com/snail007/gmc/util/sync/once"
	"io"
	"sync"
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

	gcore.RegisterSession(gcore.DefaultProviderKey, func() gcore.Session {
		return gsession.NewSession()
	})

	gcore.RegisterSessionStorage(gcore.DefaultProviderKey, func(ctx gcore.Ctx) (gcore.SessionStorage, error) {
		return gsession.Init(ctx.Config())
	})

	gcore.RegisterView(gcore.DefaultProviderKey, func(w io.Writer, tpl gcore.Template) gcore.View {
		return gview.New(w, tpl)
	})

	gcore.RegisterTemplate(gcore.DefaultProviderKey, func(ctx gcore.Ctx, rootDir string) (gcore.Template, error) {
		if ctx.Config().Sub("template") != nil {
			return gtemplate.Init(ctx)
		}
		return gtemplate.NewTemplate(ctx, rootDir)
	})

	gcore.RegisterHTTPRouter(gcore.DefaultProviderKey, func(ctx gcore.Ctx) gcore.HTTPRouter {
		return grouter.NewHTTPRouter(ctx)
	})

	gcore.RegisterCookies(gcore.DefaultProviderKey, func(ctx gcore.Ctx) gcore.Cookies {
		return gcookie.New(ctx.Response(), ctx.Request())
	})

	gcore.RegisterI18n(gcore.DefaultProviderKey, func(ctx gcore.Ctx) (gcore.I18n, error) {
		var err error
		gonce.OnceDo("gmc-i18n-init", func() {
			err = gi18n.Init(ctx.Config())
		})
		return gi18n.I18N, err
	})

	gcore.RegisterController(gcore.DefaultProviderKey, func(ctx gcore.Ctx) gcore.Controller {
		c := &gcontroller.Controller{}
		c.Ctx = ctx
		return c
	})

	gcore.RegisterHTTPServer(gcore.DefaultProviderKey, func(ctx gcore.Ctx) gcore.HTTPServer {
		return ghttpserver.NewHTTPServer(ctx)
	})

	gcore.RegisterAPIServer(gcore.DefaultProviderKey, func(ctx gcore.Ctx, address string) (gcore.APIServer, error) {
		return ghttpserver.NewAPIServerForProvider(ctx, address)
	})
}
