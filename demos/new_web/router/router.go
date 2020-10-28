package router

import (
	"github.com/snail007/gmc"
	"mygmcweb/controller"
)

func InitRouter(s *gmc.HTTPServer) {
	// acquire router object
	r := s.Router()

	// bind a controller, /demo is path of controller, after this you can visit http://127.0.0.1:7080/demo/hello
	// "hello" is full lower case name of controller method.
	r.Controller("/demo", new(controller.Demo))
	r.ControllerMethod("/", new(controller.Demo), "Index__")
	r.ControllerMethod("/index.html", new(controller.Demo), "Index__")

	// indicates router initialized
	s.Logger().Infof("router inited.")
}
