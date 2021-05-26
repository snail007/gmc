// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package basic_test

import (
	gcore "github.com/snail007/gmc/core"
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
			obj = gcore.ProviderApp()(false)
			return
		}, (*gcore.App)(nil), "default app"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderApp()(true)
			return
		}, (*gcore.App)(nil), "default appDefault"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderConfig()()
			return
		}, (*gcore.Config)(nil), "default config server"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderLogger()(ctx, "")
			return
		}, (*gcore.Logger)(nil), "default logger"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderLogger()(nil, "")
			return
		}, (*gcore.Logger)(nil), "default logger, ctx is nil"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderError()()
			return
		}, (*gcore.Error)(nil), "default error"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderCtx()()
			return
		}, (*gcore.Ctx)(nil), "default ctx"},
	} {
		obj, err := v.factory()
		assert.Nil(err)
		assert.NotNil(obj)
		assert.Implements(v.impl, obj)
	}
}
