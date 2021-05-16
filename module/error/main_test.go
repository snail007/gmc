// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gerror

import (
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gcore.RegisterConfig(gcore.DefaultProviderKey, func() gcore.Config {
		return gconfig.New()
	})

	os.Exit(m.Run())
}
