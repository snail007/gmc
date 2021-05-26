// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtemplate

import (
	gcore "github.com/snail007/gmc/core"
	grouter "github.com/snail007/gmc/http/router"
	gsession "github.com/snail007/gmc/http/session"
	gview "github.com/snail007/gmc/http/view"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	gi18n "github.com/snail007/gmc/module/i18n"
	glog "github.com/snail007/gmc/module/log"
	"io"
	"os"
	"sync"
	"testing"
)

var (
	tpl *Template
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
			return Init(ctx)
		}
		return NewTemplate(ctx, rootDir)
	})

	gcore.RegisterHTTPRouter(gcore.DefaultProviderKey, func(ctx gcore.Ctx) gcore.HTTPRouter {
		return grouter.NewHTTPRouter(ctx)
	})

	gcore.RegisterConfig(gcore.DefaultProviderKey, func() gcore.Config {
		return gconfig.New()
	})

	gcore.RegisterI18n(gcore.DefaultProviderKey, func(ctx gcore.Ctx) (gcore.I18n, error) {
		var err error
		OnceDo("gmc-i18n-init", func() {
			err = gi18n.Init(ctx.Config())
		})
		return gi18n.I18N, err
	})

	gcore.RegisterLogger(gcore.DefaultProviderKey, func(ctx gcore.Ctx, prefix string) gcore.Logger {
		if ctx == nil {
			return glog.New(prefix)
		}
		return glog.NewFromConfig(ctx.Config(), prefix)
	})

	gcore.RegisterCtx(gcore.DefaultProviderKey, func() gcore.Ctx {
		return gctx.NewCtx()
	})

	ctx := gcore.ProviderCtx()()
	ctx.SetConfig(gcore.ProviderConfig()())
	tpl, _ = NewTemplate(ctx, "tests/views")
	tpl.Delims("{{", "}}")
	tpl.Parse()
	os.Exit(m.Run())
}

var onceDoDataMap = sync.Map{}

func OnceDo(uniqueKey string, f func()) {
	once, _ := onceDoDataMap.LoadOrStore(uniqueKey, &sync.Once{})
	once.(*sync.Once).Do(f)
	return
}
