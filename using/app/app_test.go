// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package app

import (
	gcore "github.com/snail007/gmc/core"
	_ "github.com/snail007/gmc/using/basic"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestGMCImplements(t *testing.T) {
	assert := assert2.New(t)
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
	} {
		obj, err := v.factory()
		assert.Nilf(err, v.msg)
		assert.NotNilf(obj, v.msg)
		assert.Implementsf(v.impl, obj, v.msg)
	}
}
