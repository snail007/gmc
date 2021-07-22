// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdefer

import (
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
)

func Recover(f func(err gcore.Error)) {
	e := recover()
	if e != nil {
		f(gerror.New().Wrap(e))
	}
}

func RecoverNop() {
	_ = recover()
}

func RecoverNopFunc(f func()) {
	_ = recover()
	f()
}
