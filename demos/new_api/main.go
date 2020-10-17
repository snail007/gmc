package main

import (
	"github.com/snail007/gmc"
	"mygmcapi/handlers"
)

func main() {
	app := gmc.New.App()
	//cfg,err:=gmc.New.ConfigFile("conf/app.toml")
	//if err!=nil{
	//	app.Logger().Fatal(err)
	//}
	//
	//addres := cfg.GetString("apiserver.listen")
	//api := gmc.New.APIServer(addres)
	//isTLS := cfg.GetBool("apiserver.tlsenable")
	//if isTLS {
	//	cert, key := cfg.GetString("apiserver.tlscert"), cfg.GetString("apiserver.tlskey")
	//	api.SetTLSFile(cert, key)
	//}
	api := gmc.New.APIServer(":8030")
	handlers.Init(api)

	app.AddService(gmc.ServiceItem{
		Service: api,
		BeforeInit: func(s gmc.Service, cfg *gmc.Config) (err error) {
			api.PrintRouteTable(nil)
			return
		},
	})
	app.Logger().Panic(gmc.StackE(app.Run()))
}
