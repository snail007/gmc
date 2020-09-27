package router

import (
	"strings"

	"github.com/snail007/gmc"

	"github.com/snail007/gmc/demos/website/controller"
)

func InitRouter(s *gmc.HTTPServer) {
	// sets pre routing handler, it be called with any request.
	s.AddMiddleware0(filterAll)

	// sets post routing handler, it be called only when url's path be found in router.
	s.AddMiddleware1(filter)

	// acquire router object
	r := s.Router()

	// bind a controller, /demo is path of controller, after this you can visit http://127.0.0.1:7080/demo/hello
	// "hello" is full lower case name of controller method.
	r.Controller("/demo", new(controller.Demo))

	// indicates router initialized
	s.Logger().Printf("router inited.")
}
func filterAll(c gmc.C, server *gmc.HTTPServer) bool {
	server.Logger().Printf(c.Request.RequestURI)
	return false
}
func filter(c gmc.C, server *gmc.HTTPServer) bool {
	path := strings.TrimRight(c.Request.URL.Path, "/\\")

	// we want to prevent user to access method `controller.Demo.Protected`
	if strings.HasSuffix(path, "protected") {
		c.Write([]byte("404"))
		return true
	}
	// server.Logger().Printf(r.RequestURI)
	return false
}
