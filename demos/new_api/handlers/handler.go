package handlers

import (
	"github.com/snail007/gmc"
	gmcrouter "github.com/snail007/gmc/http/router"
	gmchttpserver "github.com/snail007/gmc/http/server"
	"github.com/snail007/gmc/util/timeutil"
	"time"
)

func  initHanlder(api *gmchttpserver.APIServer)  {

	// /
	api.API("/", func(c *gmcrouter.Ctx) {
		//var out bytes.Buffer
		//for k:=range api.Router().RouteTable(){
		//
		//	out.WriteString("http://127.0.0.1")
		//}
	})

	api.API("/hello", func(c gmc.C) {
		api.Logger().Printf("request %s", c.Request.RequestURI)
		c.Write("hello world!")
	})

	api.API("/hi", func(c gmc.C) {
		c.Write("hi!")
	})

	// trigger a panic error
	api.API("/error", func(c gmc.C) {
		a := 0
		a /= a
	})

	api.API("/version", func(c gmc.C) {
		c.Write(1.1)
	})

	api.API("/sleep", func(c gmc.C) {
		time.Sleep(time.Second * 10)
		c.Write("reload")
	})

	// routing by group is supported
	group1 := api.Group("/v1")
	group1.API("/time", func(c gmc.C) {
		c.Write(timeutil.TimeToStr(time.Now()))
	})



}
