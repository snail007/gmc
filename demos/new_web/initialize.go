package main

import (
	"github.com/snail007/gmc"
)

func Initialize(s *gmc.HTTPServer) (err error) {
	// init router
	InitRouter(s)

	return
}
