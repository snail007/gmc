package router

import (
	"fmt"
	"net/http"
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
func (s *HttpRouter) Route(urlPath string, obj interface{}) {
	for _, m := range methods(obj) {

		path := urlPath + strings.ToLower(m)
		fmt.Println(m, path)
		for _, vv := range shortcutMethods {
			if strings.HasSuffix(vv, "_") {
				continue
			}
			s.HandlerFunc(vv, path, func(w http.ResponseWriter, r *http.Request) {
				// objv := reflect.ValueOf(obj)
				objv := obj
				fmt.Println(">>>>", m, ">>>>", s, ">>>>", w, ">>>>", r)
				// invoke(objv, "Pre_", s, w, r)
				invoke(objv, "Pre_", w)
				invoke(objv, m)
				invoke(objv, "Post_")
			})
		}
	}
}
func (s *HttpRouter) RouteMethod(urlPath, handler func(http.ResponseWriter, *http.Request)) {

}
