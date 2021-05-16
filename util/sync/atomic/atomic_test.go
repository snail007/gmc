// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gatomic

import (
	assert2 "github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)
	type vt struct {
		cnt int
	}
	value := NewValue(vt{1})
	for i := 0; i < 2; i++ {
		go func() {
			value.Store(vt{1})
		}()
	}
	for i := 0; i < 2; i++ {
		go func() {
			assert.Equal(1, value.Load().(vt).cnt)
		}()
	}
	for i := 0; i < 2; i++ {
		go func() {
			value.LoadAndStore(func(x interface{}) interface{} {
				assert.Equal(1, x.(vt).cnt)
				return 1
			})
		}()
	}
	time.Sleep(time.Millisecond *200 * time.Duration(runtime.NumCPU()))
}
