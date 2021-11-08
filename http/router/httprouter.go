// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package grouter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"

	gcore "github.com/snail007/gmc/core"
)

var (
	anyMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}
	//helper functions in Controller in the list will exclude in router
	skipMethods = map[string]bool{
		"Die":            true,
		"SessionDestroy": true,
		"SessionStart":   true,
		"Stop":           true,
		"StopE":          true,
		"Write":          true,
		"WriteE":         true,
		"Tr":             true,
		"Before":         true,
		"After":          true,
		"MethodCallPost": true,
		"MethodCallPre":  true,
		"GetCtx":         true,
	}
)

type HTTPRouter struct {
	*Router
	// parent httprouter of current group
	hr *HTTPRouter
	//namespace of current group
	ns  string
	ext string
	ctx gcore.Ctx
}

func NewHTTPRouter(ctx gcore.Ctx) *HTTPRouter {
	hr := &HTTPRouter{
		Router: &Router{
			RedirectTrailingSlash:  false,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			HandleOPTIONS:          true,
			SaveMatchedRoutePath:   true, //if true this.Args in controller is always not be nil, if false it maybe nil.
		},
		ns:  "/",
		ctx: ctx,
	}
	return hr
}

//Group create a group in namespace ns
func (s *HTTPRouter) Group(ns string) gcore.HTTPRouter {
	if !strings.HasSuffix(ns, "/") {
		ns += "/"
	}
	return &HTTPRouter{
		Router: s.Router,
		hr:     s,
		ns:     ns,
		ctx:    s.ctx,
	}
}
func (s *HTTPRouter) Namespace() string {
	parentNS := ""
	if s.hr != nil {
		parentNS = s.hr.ns
	}
	return strings.TrimRight(parentNS, "/") + s.ns
}

// Ext sets Controller()'s default ext
func (this *HTTPRouter) Ext(ext string) {
	this.ext = ext
}

// Controller binds a controller's methods to router
// ext is method's extension in url.
func (s *HTTPRouter) Controller(urlPath string, obj gcore.Controller, ext ...string) {
	s.controller(urlPath, obj, "", ext...)
}

// PrintRouteTable dump all routes into `w`, if `w` is nil, os.Stdout will be used.
func (s *HTTPRouter) PrintRouteTable(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	m := s.RouteTable()
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return
	}
	sort.Strings(keys)
	maxmlen := 0
	maxplen := 0
	for _, k := range keys {
		l := len(strings.Join(m[k], ","))
		if maxmlen < l {
			maxmlen = l
		}
		l = len(k)
		if maxplen < l {
			maxplen = l
		}
	}
	t1 := strings.Repeat("-", maxplen)
	t2 := strings.Repeat("-", maxmlen)
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\n:ROUTE TABLE\n| %-"+fmt.Sprintf("%d", maxplen)+"s | %s\n", "PATH", "METHOD"))
	buf.WriteString(fmt.Sprintf("| %s | %s\n", t1, t2))
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("| %-"+fmt.Sprintf("%d", maxplen)+"s | %s\n", k, strings.Join(m[k], ",")))
	}
	buf.WriteString(fmt.Sprintf("| %s | %s\n\n", t1, t2))
	fmt.Fprint(w, buf.String())
}

//RouteTable returns all routes in router. KEY is url path, VALUE is http methods.
func (s *HTTPRouter) RouteTable() (table map[string][]string) {
	t := &map[string][]string{}
	for k, v := range s.trees {
		s.visit(v, "", k, t)
	}
	return *t
}
func (s *HTTPRouter) visit(n *node, prefix, m string, p *map[string][]string) {
	path := prefix + n.path
	if n.handle != nil {
		(*p)[path] = append((*p)[path], m)
	}
	for _, v := range n.children {
		s.visit(v, path, m, p)
	}
}

// ControllerMethod binds a controller's method to router
func (s *HTTPRouter) ControllerMethod(urlPath string, obj gcore.Controller, method string) {
	s.controller(urlPath, obj, method)
}
func (s *HTTPRouter) controller(urlPath string, obj interface{}, method string, ext ...string) {
	bindMethods := []string{}
	allMethods := map[string]bool{}
	if method != "" {
		bindMethods = append(bindMethods, method)
	}
	for _, objMethod := range methods(obj) {
		allMethods[objMethod] = true
		if method != "" ||
			skipMethods[objMethod] ||
			strings.HasSuffix(objMethod, "__") ||
			strings.HasSuffix(objMethod, "_") {
			continue
		}
		bindMethods = append(bindMethods, objMethod)
	}

	beforeIsFound := allMethods["Before"]
	afterIsFound := allMethods["After"]

	for _, objMethod := range bindMethods {
		path := ""
		if objMethod == method {
			path = urlPath
		} else {
			p := urlPath
			if !strings.HasSuffix(p, "/") {
				p += "/"
			}
			// default
			ext1 := s.ext
			if len(ext) > 0 {
				// cover
				ext1 = ext[0]
			}
			path = p + strings.ToLower(objMethod) + ext1
		}
		objMethod0 := objMethod
		s.HandleAny(path, func(w http.ResponseWriter, _ *http.Request, ps gcore.Params) {
			reqCtx := gcore.GetCtx(w)
			// fix param not contains matched route path
			reqCtx.SetParam(ps)
			reqCtx.SetControllerMethod(objMethod0)
			obj0 := reflect.ValueOf(obj)
			var val reflect.Value
			if obj0.Kind() == reflect.Ptr {
				val = obj0.Elem()
			} else {
				val = obj0
			}
			vp := reflect.New(val.Type())
			vp.Elem().Set(val)
			objv := vp.Interface()
			reqCtx.SetController(objv.(gcore.Controller))
			defer invoke(objv, "MethodCallPost")
			invoke(objv, "MethodCallPre", reqCtx)
			if beforeIsFound {
				invoke(objv, "Before")
			}
			s.call(func() { invoke(objv, objMethod0) })
			if afterIsFound {
				invoke(objv, "After")
			}
		})
	}
	if method != "" && !allMethods[method] {
		panic(gcore.ProviderError()().New("route [ " + urlPath + " ], method [ " + method + " ] not found"))
	}
}

func (s *HTTPRouter) call(fn func()) {
	func() {
		defer gcore.ProviderError()().Recover(func(e interface{}) {
			if fmt.Sprintf("%s", e) == "__STOP__" {
				return
			}
			panic(gcore.ProviderError()().Wrap(e))
		})
		fn()
	}()
	return
}

func (s *HTTPRouter) path(path string) string {
	ns := strings.TrimRight(s.Namespace(), "/")
	return ns + path
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (s *HTTPRouter) Handle(method, path string, handle gcore.Handle) {
	p := s.path(path)
	s.Router.Handle(method, p, handle)
}

// HandleAny registers a new request handle with the given path and all http methods,
// GET, POST, PUT, PATCH, DELETE and OPTIONS
func (s *HTTPRouter) HandleAny(path string, handle gcore.Handle) {
	for _, method := range anyMethods {
		s.Handle(method, path, handle)
	}
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
// The Params are available in the request context under ParamsKey.
func (s *HTTPRouter) Handler(method, path string, handler http.Handler) {
	s.Router.Handler(method, s.path(path), handler)
}

// HandlerAny is an adapter which allows the usage of an http.Handler as a
// request handle match all http methods,
// GET, POST, PUT, PATCH, DELETE and OPTIONS
func (s *HTTPRouter) HandlerAny(path string, handler http.Handler) {
	for _, method := range anyMethods {
		s.Handler(method, path, handler)
	}
}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
func (s *HTTPRouter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	s.Router.HandlerFunc(method, s.path(path), handler)
}

// HandlerFuncAny is an adapter which allows the usage of an http.HandlerFunc as a
// request handle match all http methods,
// GET, POST, PUT, PATCH, DELETE and OPTIONS
func (s *HTTPRouter) HandlerFuncAny(path string, handler http.HandlerFunc) {
	for _, method := range anyMethods {
		s.HandlerFunc(method, path, handler)
	}
}
