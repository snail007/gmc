// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package grouter_test

import (
	gcore "github.com/snail007/gmc/core"
	grouter "github.com/snail007/gmc/http/router"
	gconfig "github.com/snail007/gmc/module/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {


	gcore.RegisterHTTPRouter(gcore.DefaultProviderKey, func(ctx gcore.Ctx) gcore.HTTPRouter {
		return grouter.NewHTTPRouter(ctx)
	})

	gcore.RegisterConfig(gcore.DefaultProviderKey, func() gcore.Config {
		return gconfig.New()
	})

	os.Exit(m.Run())
}
