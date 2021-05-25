// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gatomic

import (
	assert2 "github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	assert := assert2.New(t)
	type data struct {
		cnt int
	}
	g:=sync.WaitGroup{}
	g.Add(6)
	value := NewValue(data{1})
	for i := 0; i < 2; i++ {
		go func() {
			defer g.Done()
			value.Store(data{1})
		}()
	}

	for i := 0; i < 2; i++ {
		go func() {
			defer g.Done()
			assert.Equal(1, value.Load().(data).cnt)
		}()
	}

	for i := 0; i < 2; i++ {
		go func() {
			defer g.Done()
			value.LoadAndStore(func(x interface{}) interface{} {
				assert.Equal(1, x.(data).cnt)
				d := x.(data)
				d.cnt =1
				return d
			})
		}()
	}
	g.Wait()
}
