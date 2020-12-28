package gi18n

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	providers := gcore.Providers

	providers.RegisterI18n("", func(ctx gcore.Ctx) (gcore.I18n,error) {
		return Init(ctx.Config())
	})

	providers.RegisterConfig("", func() gcore.Config {
		return gconfig.NewConfig()
	})

	os.Exit(m.Run())
}
