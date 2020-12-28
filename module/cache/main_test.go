package gcache

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	glog "github.com/snail007/gmc/module/log"
	gutil "github.com/snail007/gmc/util"
	"os"
	"testing"
)

var (
	cFile  gcore.Cache
	cMem   gcore.Cache
	cRedis gcore.Cache
)

func TestMain(m *testing.M) {
	var err error

	providers := gcore.Providers

	providers.RegisterConfig("", func() gcore.Config {
		return gconfig.NewConfig()
	})

	providers.RegisterCache("", func(ctx gcore.Ctx) (gcore.Cache, error) {
		var err error
		gutil.OnceDo("gmc-cache-init", func() {
			err = Init(ctx.Config())
		})
		if err != nil {
			return nil, err
		}
		return Cache(), nil
	})

	providers.RegisterLogger("", func(ctx gcore.Ctx,prefix string) gcore.Logger {
		if ctx == nil {
			return glog.NewGMCLog(prefix)
		}
		return glog.NewFromConfig(ctx.Config(), prefix)
	})

	providers.RegisterCtx("", func() gcore.Ctx {
		return gctx.NewCtx()
	})

	c, e := gctx.NewCtxFromConfigFile("../app/app.toml")
	if e != nil {
		panic(e)
	}

	logger = gcore.Providers.Logger("")(c,"")


	cFile, err = NewFileCache(NewFileCacheConfig())
	if err != nil {
		panic(err)
	}
	cMem = NewMemCache(NewMemCacheConfig())
	os.Exit(m.Run())
}
