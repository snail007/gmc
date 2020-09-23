package initialize

import (
	gmchttpserver "github.com/snail007/gmc/http/server"
)

func Initialize(s *gmchttpserver.HTTPServer) (err error) {
	// initialize your something here
	s.Logger().Println("using config file : ", s.Config().ConfigFileUsed())
	return
}
