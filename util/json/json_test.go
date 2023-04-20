// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gjson

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONResult_Parse(t *testing.T) {
	j := NewResult(`{"code":1,"data":"1","message":"2"}`)
	assert.NotNil(t, j)
	assert.Equal(t, 1, j.Code())
	assert.Equal(t, "1", j.Data())
	assert.Equal(t, "2", j.Message())
	j = NewResult([]byte(`{"code":1,"data":"1","message":"2"}`))
	assert.NotNil(t, j)
	assert.Equal(t, 1, j.Code())
	assert.Equal(t, "1", j.Data())
	assert.Equal(t, "2", j.Message())
	j = NewResult(1, 2, 3)
	assert.Equal(t, 1, j.Code())
	assert.Equal(t, "2", j.Message())
	assert.Equal(t, 3, j.Data())
	j = NewResult([]byte(`{`))
	assert.Nil(t, j)
}
