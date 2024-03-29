// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gerror

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	str := ""
	f1 := func() {
		defer Recover(func(err interface{}) {
			str = Wrap(err).Error()
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

func TestStack(t *testing.T) {
	str := ""
	f1 := func() {
		defer Recover(func(err interface{}) {
			str = Stack(err)
		})
		panic("okay")
		str = "okay"
	}
	f1()
	assert.Contains(t, str, "*errors.errorString")
}

func TestNew(t *testing.T) {
	str := ""
	f1 := func() {
		defer Recover(func(err interface{}) {
			str = New(err).ErrorStack()
		})
		panic("okay")
		str = "okay"
	}
	f1()
	assert.Contains(t, str, "*errors.errorString")
}

func TestWrap(t *testing.T) {
	str := ""
	f1 := func() {
		defer Recover(func(err interface{}) {
			str = Wrap(err).ErrorStack()
		})
		panic("okay")
		str = "okay"
	}
	f1()
	assert.Contains(t, str, "*errors.errorString")
}

func TestTry(t *testing.T) {
	e := Try(func() {
		a := 0
		fmt.Print(1 / a)
	})
	assert.Contains(t, e, "integer divide by zero")
}

func TestTryWithStack(t *testing.T) {
	e := TryWithStack(func() {
		a := 0
		fmt.Print(1 / a)
	})
	assert.Contains(t, e.StackStr(), "divideError")
	assert.Contains(t, e.ErrorStack(), "integer divide by zero")
}
