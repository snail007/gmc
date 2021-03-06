// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package db

import (
	gcore "github.com/snail007/gmc/core"
	gdb "github.com/snail007/gmc/module/db"
	// basic requirements
	_ "github.com/snail007/gmc/using/basic"
	gonce "github.com/snail007/gmc/util/sync/once"
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

	gcore.RegisterDatabase(gcore.DefaultProviderKey, func(ctx gcore.Ctx) (gcore.Database, error) {
		var err error
		gonce.OnceDo("gmc-cache-init", func() {
			err = gdb.Init(ctx.Config())
		})
		if err != nil {
			return nil, err
		}
		return gdb.DB(), nil
	})
}
