// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttppprof

import (
	gcore "github.com/snail007/gmc/core"
	"net/http"
	"net/http/pprof"
	"strings"
)

func BindRouter(r gcore.HTTPRouter, prefix string) {
	if prefix == "" {
		prefix = "/debug/pprof/"
	}
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}
	root := strings.TrimSuffix(prefix, "/")
	r.HandlerAny(root+"/allocs", pprof.Handler("allocs"))
	r.HandlerAny(root+"/block", pprof.Handler("block"))
	r.HandlerAny(root+"/goroutine", pprof.Handler("goroutine"))
	r.HandlerAny(root+"/mutex", pprof.Handler("mutex"))
	r.HandlerAny(root+"/heap", pprof.Handler("heap"))
	r.HandlerAny(root+"/threadcreate", pprof.Handler("threadcreate"))
	r.HandlerAny(root+"/profile", pprofHandler(pprof.Profile))
	r.HandlerAny(root+"/cmdline", pprofHandler(pprof.Cmdline))
	r.HandlerAny(root+"/symbol", pprofHandler(pprof.Symbol))
	r.HandlerAny(root+"/trace", pprofHandler(pprof.Trace))
	r.HandlerAny(root+"/", pprofHandler(pprof.Index))
	r.HandlerAny(root, http.RedirectHandler(root+"/", 302))
}

func pprofHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}
