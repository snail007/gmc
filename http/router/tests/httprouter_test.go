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
func TestGroup_1(t *testing.T) {
	assert := assert.New(t)
	r := gmcrouter.NewHTTPRouter()
	r.Controller("/user", new(Controller))
	g1 := r.Group("/v1")
	g1.Controller("/user", new(Controller))
	g11 := g1.Group("/1.0")
	g11.Controller("/user", new(Controller))
	g2 := r.Group("/v2")
	g2.Controller("/user", new(Controller))
	g3 := r.Group("/img/v3")
	g3.Controller("/user", new(Controller))
	want := []string{
		"/user/method1",
		"/user/testmethod",
		"/v1/user/method1",
		"/v1/user/testmethod",
		"/v1/1.0/user/method1",
		"/v1/1.0/user/testmethod",
		"/v2/user/method1",
		"/v2/user/testmethod",
		"/img/v3/user/method1",
		"/img/v3/user/testmethod",
	}
	rt := r.RouteTable()
	for _, v := range want {
		assert.NotNil(rt[v])
	}
}
func TestGroup_2(t *testing.T) {
	assert := assert.New(t)
	r := gmcrouter.NewHTTPRouter()
	r.Controller("/user", new(Controller))
	g1 := r.Group("/v1")
	g1.Controller("/user", new(Controller))
	g11 := g1.Group("/1.0")
	g11.Controller("/user", new(Controller))
	g2 := r.Group("/v2")
	g2.Controller("/user", new(Controller))
	g3 := r.Group("/img/v3")
	g3.Controller("/user", new(Controller))
	want := []string{
		"/user/method1",
		"/user/testmethod",
		"/v1/user/method1",
		"/v1/user/testmethod",
		"/v1/1.0/user/method1",
		"/v1/1.0/user/testmethod",
		"/v2/user/method1",
		"/v2/user/testmethod",
		"/img/v3/user/method1",
		"/img/v3/user/testmethod",
	}
	rt := r.RouteTable()
	for _, v := range want {
		assert.NotNil(rt[v])
	}
	r.PrintRouteTable(nil)
	// t.Fail()
}
