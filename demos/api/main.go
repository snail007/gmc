package main

import (
	"fmt"
	"net/http"

	"github.com/snail007/gmc"
	gmchook "github.com/snail007/gmc/process/hook"
)

func main() {
	api := gmc.NewAPIServer(":8030")
	// sets a filter to filter all request registered in router,
	// exclude 404 requests.
	api.Before(func(w gmc.W, r gmc.R) bool {
		api.Logger().Printf("before request %s", r.RequestURI)
		return true
	})
	// sets a function called after every request registered in router,
	// exclude 404 requests.
	api.After(func(w gmc.W, r gmc.R, ps gmc.P, isPanic bool) {
		status := "200"
		if isPanic {
			status = "500"
		}
		api.Logger().Printf("after request %s %s", status, r.RequestURI)
		return
	})
	// sets a function to handle 404 requests.
	api.Handle404(func(w gmc.W, r gmc.R) {
		gmc.Write(w, "404")
	})
	// sets a function to handle panic error.
	api.Handle500(func(w gmc.W, r gmc.R, ps gmc.P, err interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		gmc.Write(w, "panic error : ", err)
	})
	// sets an api in url path: /hello
	// add more api , just call api.API() repeatedly
	api.API("/hello", func(w gmc.W, r gmc.R, ps gmc.P) {
		api.Logger().Printf("request %s", r.RequestURI)
		gmc.Write(w, "hello world!")
	}).API("/hi", func(w gmc.W, r gmc.R, ps gmc.P) {
		gmc.Write(w, "hi!")
	})
	// trigger a panic error
	api.API("/error", func(w gmc.W, r gmc.R, ps gmc.P) {
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
