// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gbytes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPool(t *testing.T) {
	p1 := GetPool(1024)
	p2 := GetPool(1024)
	assert.Exactly(t, p1, p2)
	assert.Equal(t, 1024, len(p1.Get().([]byte)))
}

func TestPool_Put(t *testing.T) {
	p := GetPool(1024)
	buf1 := p.Get()
	p.Put(buf1)
	buf2 := p.Get()
	assert.Exactly(t, buf1, buf2)
	assert.Equal(t, 1024, len(buf1.([]byte)))
}
