// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package router

import (
	"testing"

	gmccontroller "github.com/snail007/gmc/http/controller"

	"github.com/stretchr/testify/assert"

	gmcrouter "github.com/snail007/gmc/http/router"
)

type Controller struct {
	gmccontroller.Controller
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
	r := gmcrouter.NewHTTPRouter()
	r.Controller("/user/", new(Controller))
	h, _, _ := r.Lookup("GET", "/user/method1")
	assert.NotNil(h)
}
func TestControllerMethod(t *testing.T) {
	assert := assert.New(t)
	r := gmcrouter.NewHTTPRouter()
	r.ControllerMethod("/method/:name", new(Controller), "TestMethod")
	//test Controller
	h, _, _ := r.Lookup("GET", "/method/hello")
	assert.NotNil(h)
}