package router

import (
	"strings"

	"github.com/snail007/gmc"

	"github.com/snail007/gmc/demos/website/controller"
)

func InitRouter(s *gmc.HTTPServer) {
	// sets pre routing handler, it be called with any request.
	s.BeforeRouting(filterAll)

	// sets post routing handler, it be called only when url's path be found in router.
	s.RoutingFiliter(filter)

	// acquire router object
	r := s.Router()

	// bind a controller, /demo is path of controller, after this you can visit http://127.0.0.1:7080/demo/hello
	// "hello" is full lower case name of controller method.
	r.Controller("/demo", new(controller.Demo))

	// indicates router initialized
	s.Logger().Printf("router inited.")
}
func filterAll(w gmc.W, r gmc.R, server *gmc.HTTPServer) bool {
	server.Logger().Printf(r.RequestURI)
	return true
}
func filter(w gmc.W, r gmc.R, ps gmc.P, server *gmc.HTTPServer) bool {
	path := strings.TrimRight(r.URL.Path, "/\\")

	// we want to prevent user to access method `controller.Demo.Protected`
	if strings.HasSuffix(path, "protected") {
		w.Write([]byte("404"))
		return false
	}
	// server.Logger().Printf(r.RequestURI)
	return true
}
