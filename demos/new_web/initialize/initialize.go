package initialize

import (
	"github.com/snail007/gmc"
	"mygmcweb/router"
)

func Initialize(s *gmc.HTTPServer) (err error) {
	// init router
	router.InitRouter(s)

	return
}
