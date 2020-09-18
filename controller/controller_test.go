// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/snail007/gmc/config/app"
	"github.com/stretchr/testify/assert"
)

func TestSessionStart(t *testing.T) {
	assert := assert.New(t)
	//init app
	cfg := appconfig.NewAPPConfig()
	err := appconfig.Parse(cfg)
	assert.Nil(err)
	//call controller
	c := Controller{}
	r := httptest.NewRequest("GET", "http://example.com/foo", nil)
	r.AddCookie(&http.Cookie{
		Name:  "gmcsid",
		Value: "442445e36aaf565cbefc65a8bd5675df",
	})
	w := httptest.NewRecorder()
	c.MethodCallPre__(w, r, nil)
	c.SessionStart()
	//response
	resp := w.Result()
	t.Log(resp.Header, c.Session.SessionID())
	assert.Fail("")
}
