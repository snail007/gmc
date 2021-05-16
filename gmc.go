// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmc

import (
	gcore "github.com/snail007/gmc/core"
	gcontroller "github.com/snail007/gmc/http/controller"
	gcookie "github.com/snail007/gmc/http/cookie"
	grouter "github.com/snail007/gmc/http/router"
	ghttpserver "github.com/snail007/gmc/http/server"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
	gview "github.com/snail007/gmc/http/view"
	gapp "github.com/snail007/gmc/module/app"
	gcache "github.com/snail007/gmc/module/cache"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	gdb "github.com/snail007/gmc/module/db"
	gerror "github.com/snail007/gmc/module/error"
	gi18n "github.com/snail007/gmc/module/i18n"
	glog "github.com/snail007/gmc/module/log"
	"github.com/snail007/gmc/util/sync/once"
	"io"
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

	gcore.RegisterConfig(gcore.DefaultProviderKey, func() gcore.Config {
		return gconfig.NewConfig()
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

	gcore.RegisterError(gcore.DefaultProviderKey, func() gcore.Error {
		return gerror.New()
	})

	gcore.RegisterLogger(gcore.DefaultProviderKey, func(ctx gcore.Ctx, prefix string) gcore.Logger {
		if ctx == nil {
			return glog.NewLogger(prefix)
		}
		return glog.NewFromConfig(ctx.Config(), prefix)
	})

	gcore.RegisterCache(gcore.DefaultProviderKey, func(ctx gcore.Ctx) (gcore.Cache, error) {
		var err error
		gonce.OnceDo("gmc-cache-init", func() {
			err = gcache.Init(ctx.Config())
		})
		if err != nil {
			return nil, err
		}
		return gcache.Cache(), nil
	})

	gcore.RegisterDatabase(gcore.DefaultProviderKey, func(ctx gcore.Ctx) (gcore.Database, error) {
		var err error
		gonce.OnceDo("gmc-cache-init", func() {
			err = gdb.Init(ctx.Config())
		})
		if err != nil {
			return nil, err
		}
		return gdb.DB(), nil
	})

	gcore.RegisterCtx(gcore.DefaultProviderKey, func() gcore.Ctx {
		return gctx.NewCtx()
	})

	gcore.RegisterController(gcore.DefaultProviderKey, func(ctx gcore.Ctx) gcore.Controller {
		c := &gcontroller.Controller{}
		c.Ctx = ctx
		return c
	})

	gcore.RegisterHTTPServer(gcore.DefaultProviderKey, func(ctx gcore.Ctx) gcore.HTTPServer {
		return ghttpserver.NewHTTPServer(ctx)
	})

	gcore.RegisterApp(gcore.DefaultProviderKey, func(isDefault bool) gcore.App {
		return gapp.NewApp(isDefault)
	})

	gcore.RegisterAPIServer(gcore.DefaultProviderKey, func(ctx gcore.Ctx, address string) (gcore.APIServer, error) {
		return ghttpserver.NewAPIServerForProvider(ctx, address)
	})

	// init others stuff
	initHelper()
}
