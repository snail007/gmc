// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package grand

import (
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	assert := assert2.New(t)
	assert.Len(String(1), 1)
	assert.Len(String(0), 0)
	assert.Len(String(3), 3)
	for i := 0; i < 10; i++ {
		assert.NotEqual(String(3), String(3))
	}
}
func TestNew(t *testing.T) {
	assert := assert2.New(t)
	for i := 0; i < 1000; i++ {
		assert.NotEqual(New().Int(), New().Int())
	}
}

func TestIntString(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(IntString(0), "")
	assert.Len(IntString(1), 1)
	for i := 0; i < 1000; i++ {
		assert.NotEqual(IntString(2)[:1], "0")
	}
}
