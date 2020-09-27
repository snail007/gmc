package main

import (
	"fmt"
	"net/http"

	"github.com/snail007/gmc"
	gmchook "github.com/snail007/gmc/process/hook"
)

func main() {
	api := gmc.NewAPIServer(":8030")
	// add a middleware typed 1 to filter all request registered in router,
	// exclude 404 requests.
	api.AddMiddleware1(func(c gmc.C, s *gmc.APIServer) (isStop bool) {
		s.Logger().Printf("before request %s", c.Request.RequestURI)
		return false
	})
	// add a middleware typed 2 to logging every request registered in router,
	// exclude 404 requests.
	api.AddMiddleware2(func(c gmc.C, s *gmc.APIServer) (isStop bool) {
		s.Logger().Printf("after request %d %s", c.StatusCode(), c.Request.RequestURI)
		return false
	})
	// sets a function to handle 404 requests.
	api.Handle404(func(c gmc.C) {
		c.Write("404")
	})
	// sets a function to handle panic error.
	api.Handle500(func(c gmc.C, err interface{}) {
		c.WriteHeader(http.StatusInternalServerError)
		c.Write("panic error : ", err)
	})
	// sets an api in url path: /hello
	// add more api , just call api.API() repeatedly
	api.API("/hello", func(c gmc.C) {
		api.Logger().Printf("request %s", c.Request.RequestURI)
		c.Write("hello world!")
	}).API("/hi", func(c gmc.C) {
		c.Write("hi!")
	})
	// trigger a panic error
	api.API("/error", func(c gmc.C) {
		a := 0
		a /= a
	})
	// start the api server
	err := api.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(`
please visit:
http://127.0.0.1:8030/none
http://127.0.0.1:8030/hello
http://127.0.0.1:8030/error
http://127.0.0.1:8030/hi
	`)
	// block main routine
	gmchook.WaitShutdown()
}
