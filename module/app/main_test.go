package gapp

import (
	gcore "github.com/snail007/gmc/core"
	grouter "github.com/snail007/gmc/http/router"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
	gview "github.com/snail007/gmc/http/view"
	"os"
	"testing"
)

func TestMain(m *testing.M)  {
	providers:=gcore.Providers
	providers.RegisterSession("", func(ctx gcore.Ctx) gcore.Session {
		return gsession.NewSession()
	})

	providers.RegisterSessionStorage("", func(ctx gcore.Ctx) (gcore.SessionStorage, error) {
		return gsession.Init(ctx.Config())
	})

	providers.RegisterView("", func(ctx gcore.Ctx) gcore.View {
		return gview.New(ctx.Response(), ctx.WebServer().Tpl())
	})

	providers.RegisterTemplate("", func(ctx gcore.Ctx) (gcore.Template, error) {
		return gtemplate.Init(ctx.Config())
	})

	providers.RegisterHTTPRouter("", func(ctx gcore.Ctx) gcore.HTTPRouter {
		return grouter.NewHTTPRouter(ctx)
	})
	os.Exit(m.Run())
}
