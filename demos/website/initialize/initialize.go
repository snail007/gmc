package initialize

import (
	httpserver "github.com/snail007/gmc/http/server"
)

func Initialize(s *httpserver.HTTPServer) (err error) {
	// initialize your something here
	s.Logger().Println("using config file : ", s.Config().ConfigFileUsed())
	return
}
