// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package basic

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	gerror "github.com/snail007/gmc/module/error"
	glog "github.com/snail007/gmc/module/log"
	"sync"
)

var (
	once sync.Once
)

func init() {
	once.Do(func() {
		initialize()
	})
}

func initialize() {
	gcore.RegisterConfig(gcore.DefaultProviderKey, func() gcore.Config {
		return gconfig.New()
	})

	gcore.RegisterError(gcore.DefaultProviderKey, func() gcore.Error {
		return gerror.New()
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
}
