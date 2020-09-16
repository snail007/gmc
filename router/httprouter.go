package router

import (
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
func (s *HttpRouter) Controller(urlPath string, obj interface{}) {
	for _, m := range methods(obj) {
		path := urlPath + strings.ToLower(m)
		for _, vv := range shortcutMethods {
			if strings.HasSuffix(vv, "__") {
				continue
			}
			func(httpMethod, objMethod string) {
				s.Handle(httpMethod, path, func(w http.ResponseWriter, r *http.Request, ps Params) {
					objv := reflect.ValueOf(obj)
					invoke(objv, "PreCall__", w, r, ps)
					invoke(objv, objMethod)
					invoke(objv, "PostCall__")
				})
			}(vv, m)
		}
	}
}
func (s *HttpRouter) ControllerMethod(urlPath, handler func(http.ResponseWriter, *http.Request)) {

}
