// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdefer

import (
	"testing"

	gcore "github.com/snail007/gmc/core"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	str := ""
	f1 := func() {
		defer Recover(func(err gcore.Error) {
			str = err.Error()
		})
		panic("okay")
	}
	f1()
	assert.Equal(t, "okay", str)
}

func TestRecoverNop(t *testing.T) {
	str := ""
	f1 := func() {
		defer RecoverNop()
		panic("okay")
		str = "okay"
	}
	f1()
	assert.Equal(t, "", str)
}

func TestRecoverNopFunc(t *testing.T) {
	str := ""
	f1 := func() {
		defer RecoverNopFunc(func() {
			str = "okay"
		})
		panic("foo")
	}
	f1()
	assert.Equal(t, "okay", str)
}
