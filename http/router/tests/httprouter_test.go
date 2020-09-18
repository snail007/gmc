// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package router

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snail007/gmc/http/router"

	"github.com/snail007/gmc/controller"
)

type Controller struct {
	controller.Controller
}

func (this *Controller) Before__() {
	this.Response.Write([]byte("OKAY"))
}
func (this *Controller) Method1() {
	this.Response.Write([]byte("OKAY"))
}
func (this *Controller) After__() {
	this.Response.Write([]byte("OKAY"))
}
func (this *Controller) TestMethod() {
	this.Response.Write([]byte("OKAY" + this.Args.ByName("name")))
}
func TestController(t *testing.T) {
	assert := assert.New(t)
	r := router.NewHttpRouter()
	r.Controller("/user/", new(Controller))
	h, args, _ := r.Lookup("GET", "/user/method1")
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	h(w, req, args)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("OKAYOKAYOKAY", string(body))
}
func TestControllerMethod(t *testing.T) {
	assert := assert.New(t)
	r := router.NewHttpRouter()
	r.ControllerMethod("/method/:name", new(Controller), "TestMethod")
	//test Controller
	h, args, _ := r.Lookup("GET", "/method/hello")
	req := httptest.NewRequest("GET", "http://example.com/foo/", nil)
	w := httptest.NewRecorder()
	h(w, req, args)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("OKAYOKAYhelloOKAY", string(body))
}
