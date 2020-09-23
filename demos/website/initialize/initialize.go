package initialize

import "github.com/snail007/gmc"

func Initialize(s *gmc.HTTPServer) (err error) {
	// initialize your something here
	s.Logger().Println("using config file : ", s.Config().ConfigFileUsed())
	return
}
