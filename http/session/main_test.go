// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gsession

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gcache "github.com/snail007/gmc/module/cache"
	gconfig "github.com/snail007/gmc/module/config"
	gerror "github.com/snail007/gmc/module/error"
	glog "github.com/snail007/gmc/module/log"
	"github.com/snail007/gmc/util/sync/once"
	"os"
	"testing"
)

var (
	fileStore   gcore.SessionStorage
	memoryStore gcore.SessionStorage
	redisStore  gcore.SessionStorage
)

func TestMain(m *testing.M) {

	gcore.RegisterSession(gcore.DefaultProviderKey, func() gcore.Session {
		return NewSession()
	})

	gcore.RegisterSessionStorage(gcore.DefaultProviderKey, func(ctx gcore.Ctx) (gcore.SessionStorage, error) {
		return Init(ctx.Config())
	})

	gcore.RegisterConfig(gcore.DefaultProviderKey, func() gcore.Config {
		return gconfig.New()
	})

	gcore.RegisterCache(gcore.DefaultProviderKey, func(ctx gcore.Ctx) (gcore.Cache, error) {
		var err error
		gonce.OnceDo("gmc-cache-init", func() {
			err = gcache.Init(ctx.Config())
		})
		if err != nil {
			return nil, err
		}
		return gcache.Cache(), nil
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

	var err error
	cfg := NewFileStoreConfig()
	cfg.GCtime = 1
	cfg.TTL = 1
	fileStore, err = NewFileStore(cfg)
	if err != nil {
		fmt.Println(err)
	}

	cfg0 := NewMemoryStoreConfig()
	cfg0.GCtime = 1
	cfg0.TTL = 1
	memoryStore, err = NewMemoryStore(cfg0)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(m.Run())
}

//var onceDoDataMap = sync.Map{}
//
//func OnceDo(uniqueKey string, f func()) {
//	once, _ := onceDoDataMap.LoadOrStore(uniqueKey, &sync.Once{})
//	once.(*sync.Once).Do(f)
//	return
//}
