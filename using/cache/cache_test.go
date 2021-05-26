// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package cache

import (
	"bytes"
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestGMCImplements(t *testing.T) {
	assert := assert2.New(t)
	ctx := gctx.NewCtx()
	for _, v := range []struct {
		factory func() (obj interface{}, err error)
		impl    interface{}
		msg     string
	}{
		{func() (obj interface{}, err error) {
			cfg := gconfig.New()
			ctx.SetConfig(cfg)
			cfg.SetConfigType("toml")
			cfg.ReadConfig(bytes.NewReader([]byte(`
				[cache]
				default="memory"
				[[cache.memory]]
				enable=true
				id="default"
				cleanupinterval=30`)))
			obj, err = gcore.ProviderCache()(ctx)
			return
		}, (*gcore.Cache)(nil), "default ProviderCache server"},
	} {
		obj, err := v.factory()
		assert.Nilf(err, v.msg)
		assert.NotNilf(obj, v.msg)
		assert.Implementsf(v.impl, obj, v.msg)
	}
}
