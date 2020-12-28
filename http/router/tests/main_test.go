package grouter_test

import (
	gcore "github.com/snail007/gmc/core"
	grouter "github.com/snail007/gmc/http/router"
	gconfig "github.com/snail007/gmc/module/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	providers := gcore.Providers

	providers.RegisterHTTPRouter("", func(ctx gcore.Ctx) gcore.HTTPRouter {
		return grouter.NewHTTPRouter(ctx)
	})

	providers.RegisterConfig("", func() gcore.Config {
		return gconfig.NewConfig()
	})

	os.Exit(m.Run())
}
