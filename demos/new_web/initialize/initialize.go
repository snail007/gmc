package initialize

import (
	"../router"
	"github.com/snail007/gmc"
)

func Initialize(s *gmc.HTTPServer) (err error) {
	// init router
	router.InitRouter(s)

	return
}
