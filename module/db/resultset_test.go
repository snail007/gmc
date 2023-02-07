// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type a struct {
	B map[string]interface{} `json:"B"`
}
type c struct {
	B B `json:"B"`
}

type B struct {
	C string `json:"C"`
}

type cc struct {
	B *B `json:"B"`
}

func TestResultSet_mapToStruct(t *testing.T) {
	m := map[string]string{"B": `{"C":{"1":"1"}}`}
	rs := ResultSet{}
	r, e := rs.mapToStruct(m, a{})
	assert.Nil(t, e)
	assert.Equal(t, "1", r.(a).B["C"].(map[string]interface{})["1"])
}

func TestResultSet_mapToStruct1(t *testing.T) {
	m := map[string]string{"B": `{"C":"1"}`}
	rs := ResultSet{}
	r, e := rs.mapToStruct(m, c{})
	assert.Nil(t, e)
	assert.Equal(t, "1", r.(c).B.C)
}

func TestResultSet_mapToStruct2(t *testing.T) {
	m := map[string]string{"B": `{"C":"1"}`}
	rs := ResultSet{}
	r, e := rs.mapToStruct(m, cc{})
	assert.Nil(t, e)
	assert.Equal(t, "1", r.(cc).B.C)
}
