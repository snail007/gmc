// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcache

import (
	gcore "github.com/snail007/gmc/core"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func Test_Cache(t *testing.T) {
	assert := assert2.New(t)
	ctx := gcore.ProviderCtx()()
	c := gcore.ProviderConfig()()
	c.SetConfigFile("../app/app.toml")
	err := c.ReadInConfig()
	assert.Nil(err)
	ctx.SetConfig(c)
	_, err = gcore.ProviderCache()(ctx)
	assert.Nil(err)
	assert.NotNil(Memory())
	assert.NotNil(Redis())
	assert.NotNil(File())
	assert.Same(Cache(), Redis())
}
