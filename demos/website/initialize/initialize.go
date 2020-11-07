package initialize

import (
	"bytes"
	"fmt"
	"github.com/snail007/gmc"
	"github.com/snail007/gmc/demos/website/router"
	"net"
	"strings"
)

func Initialize(s *gmc.HTTPServer) (err error) {
	s.Logger().Infof("using config file : %s", s.Config().ConfigFileUsed())

	// initialize database if needed
	err = gmc.DB.Init(s.Config())
	if err != nil {
		return
	}

	// initialize cache if needed
	err = gmc.Cache.Init(s.Config())
	if err != nil {
		return
	}

	// initialize i18n if needed
	// for testing
	s.Config().Set("i18n.enable",true)
	gmc.I18n.Init(s.Config())

	// initialize router
	router.InitRouter(s)

	// add template helper functions here
	funMap := map[string]interface{}{
		"test": func(str string) string {
			return fmt.Sprintf("%d", len(str))
		},
	}
	s.AddFuncMap(funMap)

	// all path in router
	_, port, _ := net.SplitHostPort(s.Config().GetString("httpserver.listen"))
	var buf bytes.Buffer
	buf.WriteString("please visit:\n")
	for path, _ := range s.Router().RouteTable() {
		if strings.Contains(path, "*") {
			continue
		}
		buf.WriteString("http://127.0.0.1:" + port + path+"\n")
	}
	s.Logger().Writer().Write(buf.Bytes())
	return
}
