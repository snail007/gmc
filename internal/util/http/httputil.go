// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttputil

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

const (
	dieKEY  = "__DIE__"
	stopKEY = "__STOP__"
)

func Die(w io.Writer, data ...interface{}) {
	_, _ = Write(w, data...)
	panic(dieKEY)
}

func Stop(w io.Writer, data ...interface{}) {
	_, _ = Write(w, data...)
	panic(stopKEY)
}

func JustStop() {
	panic(stopKEY)
}

func JustDie() {
	panic(dieKEY)
}

// StopE will exit controller method if error is not nil.
// First argument is an error.
// Secondary argument is fail function, it be called if error is not nil.
// Third argument is success function, it be called if error is nil.
func StopE(err interface{}, fn ...func()) {
	var failFn, okayFn func()
	if len(fn) >= 1 {
		failFn = fn[0]
	}
	if len(fn) >= 2 {
		okayFn = fn[1]
	}
	if err == nil {
		if okayFn != nil {
			okayFn()
		}
		return
	}
	if failFn != nil {
		failFn()
	}
	panic("__STOP__")
}

func Write(w io.Writer, data ...interface{}) (n int, err error) {
	return write(false, w, data...)
}

func WritePretty(w io.Writer, data ...interface{}) (n int, err error) {
	return write(true, w, data...)
}

func write(pretty bool, w io.Writer, data ...interface{}) (n int, err error) {
	for _, d := range data {
		if d == nil {
			continue
		}
		switch v := d.(type) {
		case []byte:
			n, err = w.Write(v)
		case string:
			n, err = w.Write([]byte(v))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			n, err = w.Write([]byte(fmt.Sprintf("%d", v)))
		case bool:
			str := "true"
			if !v {
				str = "false"
			}
			n, err = w.Write([]byte(str))
		case float32:
			n, err = w.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		case float64:
			n, err = w.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
		case error:
			n, err = w.Write([]byte(v.Error()))
		default:
			t := reflect.TypeOf(v)
			//map, slice
			if t.Kind() == reflect.Slice || t.Kind() == reflect.Map || t.Kind() == reflect.Struct {
				var b []byte
				if pretty {
					b, err = json.MarshalIndent(v, "", "    ")
				} else {
					b, err = json.Marshal(v)
				}
				if err == nil {
					n, err = w.Write(b)
				}
			} else {
				n, err = w.Write([]byte(fmt.Sprintf("unsupported type to write: %s", t.String())))
			}
		}
		if err != nil {
			return
		}
	}
	return
}
