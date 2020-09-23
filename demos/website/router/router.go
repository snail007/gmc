package router

import (
	"net/http"
	"strings"

	"github.com/snail007/gmc/http/router"

	"github.com/snail007/gmc/demos/base/controller"
	httpserver "github.com/snail007/gmc/http/server"
)

func InitRouter(s *httpserver.HTTPServer) {
	s.BeforeRouting(filiterAll)
	s.RoutingFiliter(filiter)
	r := s.Router()
	r.Controller("/demo", new(controller.Demo))
	s.Logger().Printf("router inited.")
}
func filiterAll(w http.ResponseWriter, r *http.Request, server *httpserver.HTTPServer) bool {
	server.Logger().Printf(r.RequestURI)
	return true
}
func filiter(w http.ResponseWriter, r *http.Request, ps router.Params, server *httpserver.HTTPServer) bool {
	path := strings.TrimRight(r.URL.Path, "/\\")
	if strings.HasSuffix(path, "protected") {
		w.Write([]byte("404"))
		return false
	}
	server.Logger().Printf(r.RequestURI)
	return true
}
