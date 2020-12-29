package initialize

import (
	gcore "github.com/snail007/gmc/core"
	"mygmcweb/router"
)

func Initialize(s gcore.HTTPServer) (err error) {
	// init router
	router.InitRouter(s)

	return
}
