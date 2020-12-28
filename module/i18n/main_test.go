package gi18n

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {

	providers := gcore.Providers

	providers.RegisterI18n("", func(ctx gcore.Ctx) (gcore.I18n, error) {
		var err error
		OnceDo("gmc-i18n-init", func() {
			err = Init(ctx.Config())
		})
		return I18N, err
	})

	providers.RegisterConfig("", func() gcore.Config {
		return gconfig.NewConfig()
	})

	os.Exit(m.Run())
}

var onceDoDataMap = sync.Map{}

func OnceDo(uniqueKey string, f func()) {
	once, _ := onceDoDataMap.LoadOrStore(uniqueKey, &sync.Once{})
	once.(*sync.Once).Do(f)
	return
}
