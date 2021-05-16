// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	gcore "github.com/snail007/gmc/core"
	grouter "github.com/snail007/gmc/http/router"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
	gview "github.com/snail007/gmc/http/view"
	gconfig "github.com/snail007/gmc/module/config"
	gerror "github.com/snail007/gmc/module/error"
	"io"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	
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
		return gconfig.New()
	})

	gcore.RegisterLogger(gcore.DefaultProviderKey, func(ctx gcore.Ctx, prefix string) gcore.Logger {
		if ctx == nil {
			return New(prefix)
		}
		return NewFromConfig(ctx.Config(), prefix)
	})

	gcore.RegisterError(gcore.DefaultProviderKey, func() gcore.Error {
		return gerror.New()
	})

	os.Exit(m.Run())
}
