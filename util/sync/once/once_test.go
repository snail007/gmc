// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gonce

import (
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadOnce(t *testing.T) {
	assert := assert2.New(t)
	once := LoadOnce("test")
	assert.NotNil(once)
	//RemoveOnce("test")
	//once = LoadOnce("test")
	//assert.Nil(once)
	once1 := LoadOnce("test")
	assert.NotNil(once, once1)
}

func TestLoadOnce_1(t *testing.T) {
	assert := assert2.New(t)
	once := LoadOnce("test")
	once1 := LoadOnce("test")
	assert.Equal(fmt.Sprintf("%p", once), fmt.Sprintf("%p", once1))
}

func TestLoadOnce_2(t *testing.T) {
	assert := assert2.New(t)
	once := LoadOnce("")
	once1 := LoadOnce("")
	assert.NotEqual(fmt.Sprintf("%p", once), fmt.Sprintf("%p", once1))
}
