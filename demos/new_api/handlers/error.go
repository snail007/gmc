// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package handlers

import (
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	"net/http"
)

func initError(api gcore.APIServer) {
	// sets a function to handle 404 requests.
	api.SetNotFoundHandler(error404)
	// sets a function to handle panic error.
	api.SetErrorHandler(error500)
}

func error404(c gmc.C) {
	c.Write("404")
}

func error500(c gmc.C, err interface{}) {
	c.WriteHeader(http.StatusInternalServerError)
	c.Write(gmc.Err.Stack(err))
}
