// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package controller

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/snail007/gmc/appconfig"
	"github.com/stretchr/testify/assert"
)

func TestSessionStart(t *testing.T) {
	assert := assert.New(t)
	c := Controller{}
	r := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	c.MethodCallPre__(w, r, nil)
	c.SessionStart()
	err := appconfig.Parse()
	assert.Nil(err)
	resp := w.Result()
	t.Log(w.HeaderMap, w.Code)
	b, _ := ioutil.ReadAll(resp.Body)
	t.Log(string(b))

	assert.Fail("")
}
