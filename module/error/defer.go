// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gerror

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	"reflect"
)

func Recover(f func(err interface{})) {
	e := recover()
	if e != nil {
		f(New().Wrap(e))
	}
}

func RecoverNop() {
	_ = recover()
}

func RecoverNopFunc(f func()) {
	_ = recover()
	f()
}

func Try(f func()) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = ParseRecover(err)
		}
	}()
	f()
	return
}

func TryWithStack(f func()) (e gcore.Error) {
	defer func() {
		if err := recover(); err != nil {
			e = Wrap(err)
		}
	}()
	f()
	return
}

func ParseRecover(e interface{}) error {
	if e != nil {
		var err error
		switch ev := e.(type) {
		case error:
			err = ev
		default:
			vp := reflect.ValueOf(ev)
			if vp.Kind() == reflect.Ptr {
				err = fmt.Errorf("%s->%v", reflect.TypeOf(ev).String(), vp.Interface())
			} else {
				err = fmt.Errorf("%v", ev)
			}
		}
		return err
	}
	return nil
}
