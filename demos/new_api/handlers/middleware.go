package handlers

import (
	"github.com/snail007/gmc"
	ghttpserver "github.com/snail007/gmc/http/server"
)

func initMiddleware(api *ghttpserver.APIServer) {
	// add a middleware typed 1 to filter all request registered in router,
	// exclude 404 requests.
	api.AddMiddleware1(middleware1)
	// add a middleware typed 2 to logging every request registered in router,
	// exclude 404 requests.
	api.AddMiddleware2(middleware2)
}

func middleware1(c gmc.C, s *gmc.APIServer) (isStop bool) {
	s.Logger().Infof("before request %s", c.Request.RequestURI)
	return false
}

func middleware2(c gmc.C, s *gmc.APIServer) (isStop bool) {
	s.Logger().Infof("after request %s %d %d %s", c.Request.Method, c.StatusCode(), c.WriteCount(), c.Request.RequestURI)
	return false
}
