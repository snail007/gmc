// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package router

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

var (
	shortcutMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}
	skipMethods     = map[string]bool{
		"Die":            true,
		"SessionDestory": true,
		"SessionStart":   true,
		"Stop":           true,
		"Write":          true,
	}
)

type HTTPRouter struct {
	Router
	handle50x func(val *reflect.Value, err interface{})
}

func NewHTTPRouter() *HTTPRouter {
	hr := &HTTPRouter{
		Router: Router{
			RedirectTrailingSlash:  false,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			HandleOPTIONS:          true,
			SaveMatchedRoutePath:   true, //if true this.Args in controller is always not be nil, if false it maybe nil.
		},
	}
	return hr
}

//SetHandle50x sets handler func to handle exception error
func (s *HTTPRouter) SetHandle50x(fn func(c *reflect.Value, err interface{})) {
	s.handle50x = fn
}

//Controller binds a controller's methods to router
func (s *HTTPRouter) Controller(urlPath string, obj interface{}) {
	s.controller(urlPath, obj, "")
}

//RouteTable returns all routes in router
// func (s *HTTPRouter) RouteTable() (table map[string]string) {
// 	table = map[string]string{}
// 	for k, v := range s.trees {
// 		table[k] = v.path + v.indices
// 	}
// 	return
// }

//ControllerMethod binds a controller's method to router
func (s *HTTPRouter) ControllerMethod(urlPath string, obj interface{}, method string) {
	s.controller(urlPath, obj, method)
}
func (s *HTTPRouter) controller(urlPath string, obj interface{}, method string) {
	beforeIsFound := false
	afterIsFound := false
	for _, objMethod := range methods(obj) {
		if skipMethods[objMethod] {
			continue
		}
		switch objMethod {
		case "Before__":
			beforeIsFound = true
		case "After__":
			afterIsFound = true
		}
		if strings.HasSuffix(objMethod, "__") {
			continue
		}
		path := ""
		if method != "" {
			if objMethod != method {
				continue
			}
			path = urlPath
		} else {
			p := urlPath
			if !strings.HasSuffix(p, "/") {
				p += "/"
			}
			path = p + strings.ToLower(objMethod)
		}
		for _, httpMethod := range shortcutMethods {
			func(httpMethod, objMethod string) {
				s.Handle(httpMethod, path, func(w http.ResponseWriter, r *http.Request, ps Params) {
					objv := reflect.ValueOf(obj)
					invoke(objv, "MethodCallPre__", w, r, ps)
					if beforeIsFound && s.call(&objv, "Before__") {
						return
					}
					if s.call(&objv, objMethod) {
						return
					}
					if afterIsFound && s.call(&objv, "After__") {
						return
					}
					invoke(objv, "MethodCallPost__")
				})
			}(httpMethod, objMethod)
		}
	}
}

func (s *HTTPRouter) call(objv *reflect.Value, objMethod string) (isDIE bool) {
	func() {
		defer func() {
			e := recover()
			if e != nil {
				switch fmt.Sprintf("%s", e) {
				case "__DIE__":
					//ignore PostCall__
					isDIE = true
				case "__STOP__":
					//do nothing
				default:
					//@todo
					//exception
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
