package initialize

import (
	"fmt"
	"github.com/snail007/gmc"
	"github.com/snail007/gmc/demos/website/router"
	"net"
	"strings"
)

func Initialize(s *gmc.HTTPServer) (err error) {
	s.Logger().Println("using config file : ", s.Config().ConfigFileUsed())
	// initialize your HTTPServer something here
	err = gmc.DB.Init(s.Config())
	if err!=nil{
		return
	}
	err = gmc.Cache.Init(s.Config())
	if err!=nil{
		return
	}
	router.InitRouter(s)
	// all path in router
	_,port,_:=net.SplitHostPort(s.Config().GetString("httpserver.listen"))
	fmt.Println("please visit:")
	for path,_:=range s.Router().RouteTable(){
		if strings.Contains(path,"*"){
			continue
		}
		fmt.Println("http://127.0.0.1:"+port+path)
	}

	return
}
