// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttppprof

import (
	gcore "github.com/snail007/gmc/core"
	gctx "github.com/snail007/gmc/module/ctx"
	gjson "github.com/snail007/gmc/util/json"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strings"
)

type checkerFunc func(ctx gcore.Ctx) bool

func BindRouter(r gcore.HTTPRouter, prefix string, checker ...checkerFunc) {
	if prefix == "" {
		prefix = "/debug/pprof/"
	}
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	var c checkerFunc
	if len(checker) == 1 {
		c = checker[0]
	}
	root := strings.TrimSuffix(prefix, "/")
	r.HandlerAny(root+"/allocs", newCheckerHandler(pprof.Handler("allocs"), c))
	r.HandlerAny(root+"/block", newCheckerHandler(pprof.Handler("block"), c))
	r.HandlerAny(root+"/goroutine", newCheckerHandler(pprof.Handler("goroutine"), c))
	r.HandlerAny(root+"/mutex", newCheckerHandler(pprof.Handler("mutex"), c))
	r.HandlerAny(root+"/heap", newCheckerHandler(pprof.Handler("heap"), c))
	r.HandlerAny(root+"/threadcreate", newCheckerHandler(pprof.Handler("threadcreate"), c))
	r.HandlerAny(root+"/profile", newCheckerHandler(pprofHandler(pprof.Profile), c))
	r.HandlerAny(root+"/cmdline", newCheckerHandler(pprofHandler(pprof.Cmdline), c))
	r.HandlerAny(root+"/symbol", newCheckerHandler(pprofHandler(pprof.Symbol), c))
	r.HandlerAny(root+"/trace", newCheckerHandler(pprofHandler(pprof.Trace), c))
	r.HandlerAny(root+"/enable_block_mutex", newCheckerHandler(pprofHandler(func(w http.ResponseWriter, r *http.Request) {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
		gjson.NewResult().SetMessage("success").WriteTo(w)
	}), c))
	r.HandlerAny(root+"/disable_block_mutex", newCheckerHandler(pprofHandler(func(w http.ResponseWriter, r *http.Request) {
		runtime.SetMutexProfileFraction(0)
		runtime.SetBlockProfileRate(0)
		gjson.NewResult().SetMessage("success").WriteTo(w)
	}), c))
	r.HandlerAny(root+"/", newCheckerHandler(pprofHandler(pprof.Index), c))
	r.HandlerAny(root, newCheckerHandler(http.RedirectHandler(root+"/", 302), c))
}

func pprofHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}

type checkerHandler struct {
	h       http.Handler
	checker func(ctx gcore.Ctx) bool
}

func (c checkerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := gctx.NewCtxWithHTTP(w, r)
	if c.checker != nil && !c.checker(ctx) {
		return
	}
	c.h.ServeHTTP(w, r)
}

func newCheckerHandler(h http.Handler, checker func(ctx gcore.Ctx) bool) *checkerHandler {
	return &checkerHandler{h: h, checker: checker}
}
