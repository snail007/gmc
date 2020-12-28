package gapp

import (
	gcore "github.com/snail007/gmc/core"
	grouter "github.com/snail007/gmc/http/router"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
	gview "github.com/snail007/gmc/http/view"
	gcache "github.com/snail007/gmc/module/cache"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	gdb "github.com/snail007/gmc/module/db"
	gerror "github.com/snail007/gmc/module/error"
	gi18n "github.com/snail007/gmc/module/i18n"
	glog "github.com/snail007/gmc/module/log"
	gutil "github.com/snail007/gmc/util"
	"io"
	"os"
	"testing"
 )

func TestMain(m *testing.M) {

	providers := gcore.Providers

	providers.RegisterSession("", func() gcore.Session {
		return gsession.NewSession()
	})

	providers.RegisterSessionStorage("", func(ctx gcore.Ctx) (gcore.SessionStorage, error) {
		return gsession.Init(ctx.Config())
	})

	providers.RegisterView("", func(w io.Writer, tpl gcore.Template) gcore.View {
		return gview.New(w, tpl)
	})

	providers.RegisterTemplate("", func(ctx gcore.Ctx) (gcore.Template, error) {
		return gtemplate.Init(ctx)
	})

	providers.RegisterHTTPRouter("", func(ctx gcore.Ctx) gcore.HTTPRouter {
		return grouter.NewHTTPRouter(ctx)
	})

	providers.RegisterConfig("", func() gcore.Config {
		return gconfig.NewConfig()
	})

	providers.RegisterCtx("", func() gcore.Ctx {
		return gctx.NewCtx()
	})


	providers.RegisterI18n("", func(ctx gcore.Ctx) (gcore.I18n, error) {
		var err error
		gutil.OnceDo("gmc-i18n-init", func() {
			err = gi18n.Init(ctx.Config())
		})
		return gi18n.I18N, err
	})

	providers.RegisterError("", func() gcore.Error {
		return gerror.New()
	})

	providers.RegisterLogger("", func(ctx gcore.Ctx,prefix string) gcore.Logger {
		if ctx == nil {
			return glog.NewGMCLog(prefix)
		}
		return glog.NewFromConfig(ctx.Config(), prefix)
	})

	providers.RegisterCache("", func(ctx gcore.Ctx) (gcore.Cache, error) {
		var err error
		gutil.OnceDo("gmc-cache-init", func() {
			err = gcache.Init(ctx.Config())
		})
		if err != nil {
			return nil, err
		}
		return gcache.Cache(), nil
	})

	providers.RegisterDatabase("", func(ctx gcore.Ctx) (gcore.Database, error) {
		var err error
		gutil.OnceDo("gmc-cache-init", func() {
			err = gdb.Init(ctx.Config())
		})
		if err != nil {
			return nil, err
		}
		return gdb.DB(), nil
	})

	os.Exit(m.Run())
}
