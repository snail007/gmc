// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package handlers

import (
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
)

func initMiddleware(api gcore.APIServer) {
	// add a middleware typed 1 to filter all request registered in router,
	// exclude 404 requests.
	api.AddMiddleware1(middleware1)
	// add a middleware typed 2 to logging every request registered in router,
	// exclude 404 requests.
	api.AddMiddleware2(middleware2)
}

func middleware1(c gmc.C) (isStop bool) {
	c.APIServer().Logger().Infof("before request %s", c.Request().RequestURI)
	return false
}

func middleware2(c gmc.C) (isStop bool) {
	c.APIServer().Logger().Infof("after request %s %d %d %s", c.Request().Method, c.StatusCode(), c.WriteCount(), c.Request().RequestURI)
	return false
}
