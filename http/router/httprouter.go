// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmcrouter

import (
	"fmt"
	gmcerr "github.com/snail007/gmc/error"
	"io"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"
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
		"Write":          true,
		"Tr":             true,
	}
)

type HTTPRouter struct {
	*Router
	handle50x func(val *reflect.Value, err interface{})
	// parent httprouter of current group
	hr *HTTPRouter
	//namespace of current group
	ns string
}

func NewHTTPRouter() *HTTPRouter {
	hr := &HTTPRouter{
		Router: &Router{
			RedirectTrailingSlash:  false,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			HandleOPTIONS:          true,
			SaveMatchedRoutePath:   true, //if true this.Args in controller is always not be nil, if false it maybe nil.
		},
		ns: "/",
	}
	return hr
}

//Group create a group in namespace ns
func (s *HTTPRouter) Group(ns string) *HTTPRouter {
	if !strings.HasSuffix(ns, "/") {
		ns += "/"
	}
	return &HTTPRouter{
		Router: s.Router,
		hr:     s,
		ns:     ns,
	}
}
func (s *HTTPRouter) Namespace() string {
	parentNS := ""
	if s.hr != nil {
		parentNS = s.hr.ns
	}
	return strings.TrimRight(parentNS, "/") + s.ns
}

//SetHandle50x sets handler func to handle exception error
func (s *HTTPRouter) SetHandle50x(fn func(c *reflect.Value, err interface{})) {
	s.handle50x = fn
}

//Controller binds a controller's methods to router
func (s *HTTPRouter) Controller(urlPath string, obj interface{}) {
	s.controller(urlPath, obj, "")
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
	fmt.Fprintf(w, "\n:ROUTE TABLE\n| %-"+fmt.Sprintf("%d", maxplen)+"s | %s\n", "PATH", "METHOD")
	fmt.Fprintf(w, "| %s | %s\n", t1, t2)
	for _, k := range keys {
		fmt.Fprintf(w, "| %-"+fmt.Sprintf("%d", maxplen)+"s | %s\n", k, strings.Join(m[k], ","))
	}
	fmt.Fprintf(w, "| %s | %s\n", t1, t2)
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

//ControllerMethod binds a controller's method to router
func (s *HTTPRouter) ControllerMethod(urlPath string, obj interface{}, method string) {
	s.controller(urlPath, obj, method)
}
func (s *HTTPRouter) controller(urlPath string, obj interface{}, method string) {
	beforeIsFound := false
	afterIsFound := false
	isMehtodFound := false
	for _, objMethod := range methods(obj) {
		if isMehtodFound {
			break
		}
		if skipMethods[objMethod] {
			continue
		}
		switch objMethod {
		case "Before__":
			beforeIsFound = true
		case "After__":
			afterIsFound = true
		}
		path := ""
		if method != "" {
			if objMethod != method {
				continue
			}
			isMehtodFound = true
			path = urlPath
		} else {
			if strings.HasSuffix(objMethod, "__") {
				continue
			}
			p := urlPath
			if !strings.HasSuffix(p, "/") {
				p += "/"
			}
			path = p + strings.ToLower(objMethod)
		}
		objMethod0 := objMethod
		s.HandleAny(path, func(w http.ResponseWriter, r *http.Request, ps Params) {
			objv := reflect.ValueOf(obj)
			isPanic := false
			invoke(objv, "MethodCallPre__", w, r, ps)
			defer invoke(objv, "MethodCallPost__")
			if beforeIsFound && s.call(&objv, "Before__", &isPanic) {
				return
			}
			if !isPanic && s.call(&objv, objMethod0, &isPanic) {
				return
			}
			if !isPanic && afterIsFound && s.call(&objv, "After__", &isPanic) {
				return
			}
		})
	}
	if method != "" && !isMehtodFound {
		panic(gmcerr.New("route [ " + urlPath + " ], method [ " + method + " ] not found"))
	}
}

func (s *HTTPRouter) call(objv *reflect.Value, objMethod string, isPanic *bool) (isDIE bool) {
	func() {
		defer func() {
			e := recover()
			if e != nil {
				switch fmt.Sprintf("%s", e) {
				case "__DIE__":
					isDIE = true
				case "__STOP__":
					//do nothing
				default:
					//exception
					*isPanic = true
					if s.handle50x != nil {
						s.handle50x(objv, e)
					}
				}
			}
		}()
		invoke(*objv, objMethod)
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
func (s *HTTPRouter) Handle(method, path string, handle Handle) {
	p := s.path(path)
	s.Router.Handle(method, p, handle)
}

// HandleAny registers a new request handle with the given path and all http methods,
// GET, POST, PUT, PATCH, DELETE and OPTIONS
func (s *HTTPRouter) HandleAny(path string, handle Handle) {
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
