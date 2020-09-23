package router

import (
	"net/http"
	"strings"

	gmcrouter "github.com/snail007/gmc/http/router"

	"github.com/snail007/gmc/demos/website/controller"
	gmchttpserver "github.com/snail007/gmc/http/server"
)

func InitRouter(s *gmchttpserver.HTTPServer) {
	// sets pre routing handler, it be called with any request.
	s.BeforeRouting(filiterAll)

	// sets post routing handler, it be called only when url's path be found in router.
	s.RoutingFiliter(filiter)

	// acquire router object
	r := s.Router()

	// bind a controller, /demo is path of controller, after this you can visit http://127.0.0.1:7080/demo/hello
	// "hello" is full lower case name of controller method.
	r.Controller("/demo", new(controller.Demo))

	// indicates router initialized
	s.Logger().Printf("router inited.")
}
func filiterAll(w http.ResponseWriter, r *http.Request, server *gmchttpserver.HTTPServer) bool {
	server.Logger().Printf(r.RequestURI)
	return true
}
func filiter(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params, server *gmchttpserver.HTTPServer) bool {
	path := strings.TrimRight(r.URL.Path, "/\\")

	// we want to prevent user to access method `controller.Demo.Protected`
	if strings.HasSuffix(path, "protected") {
		w.Write([]byte("404"))
		return false
	}
	server.Logger().Printf(r.RequestURI)
	return true
}
