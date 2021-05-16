// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package app

import (
	gcore "github.com/snail007/gmc/core"
	gapp "github.com/snail007/gmc/module/app"
	_ "github.com/snail007/gmc/using/basic"
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

	gcore.RegisterApp(gcore.DefaultProviderKey, func(isDefault bool) gcore.App {
		return gapp.NewApp(isDefault)
	})
}
