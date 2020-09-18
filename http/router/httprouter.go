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
)

type HttpRouter struct {
	Router
}

func NewHttpRouter() *HttpRouter {
	hr := &HttpRouter{
		Router: Router{
			RedirectTrailingSlash:  true,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			HandleOPTIONS:          true,
		},
	}
	return hr
}

//Controller binds a controller's methods to router
func (s *HttpRouter) Controller(urlPath string, obj interface{}) {
	s.controller(urlPath, obj, "")
}

//ControllerMethod binds a controller's method to router
func (s *HttpRouter) ControllerMethod(urlPath string, obj interface{}, method string) {
	s.controller(urlPath, obj, method)
}
func (s *HttpRouter) controller(urlPath string, obj interface{}, method string) {
	beforeIsFound := false
	afterIsFound := false
	for _, m := range methods(obj) {
		switch m {
		case "Before__":
			beforeIsFound = true
		case "After__":
			afterIsFound = true
		}
		if strings.HasSuffix(m, "__") {
			continue
		}
		path := urlPath + strings.ToLower(m)
		if method != "" {
			if m != method {
				continue
			}
			path = urlPath
		}
		for _, vv := range shortcutMethods {
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
			}(vv, m)
		}
	}
}

func (s *HttpRouter) call(objv *reflect.Value, objMethod string) (isDIE bool) {
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
				}
			}
		}()
		invoke(*objv, objMethod)
	}()
	return
}
