// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcache

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	glog "github.com/snail007/gmc/module/log"
	"os"
	"sync"
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
		OnceDo("gmc-cache-init", func() {
			err = Init(ctx.Config())
		})
		if err != nil {
			return nil, err
		}
		return Cache(), nil
	})

	providers.RegisterLogger("", func(ctx gcore.Ctx, prefix string) gcore.Logger {
		if ctx == nil {
			return glog.NewLogger(prefix)
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

	logger = gcore.Providers.Logger("")(c, "")

	cFile, err = NewFileCache(NewFileCacheConfig())
	if err != nil {
		panic(err)
	}
	cMem = NewMemCache(NewMemCacheConfig())
	os.Exit(m.Run())
}

var onceDoDataMap = sync.Map{}

func OnceDo(uniqueKey string, f func()) {
	once, _ := onceDoDataMap.LoadOrStore(uniqueKey, &sync.Once{})
	once.(*sync.Once).Do(f)
	return
}
