package main

import (
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/gmc/log"
	"net"

	"github.com/snail007/gmc"
)

type MyService struct {
	gmc.Service
	l       net.Listener
	log     gcore.Logger
	address string
}

func NewMyService() gmc.Service {
	return &MyService{}
}
func (s *MyService) Init(cfg *gmc.Config) error {
	s.address = cfg.GetString("listen")
	s.log = log.NewGMCLog()
	return nil
}
func (s *MyService) Start() (err error) {
	if s.l == nil {
		s.l, err = net.Listen("tcp", s.address)
	}
	// do something
	s.log.Infof("server listen on %s", s.l.Addr().String())
	return
}
func (s *MyService) Stop() {
	s.log.Infof("server stoped on %s", s.l.Addr().String())
	// 1. close active connections ...
	// 2. stop accept
	s.l.Close()
}
func (s *MyService) GracefulStop() {
	s.log.Infof("server graceful stop on %s", s.l.Addr().String())
	// 1. stop accept
	s.l.Close()
}
func (s *MyService) SetLog(l gcore.Logger) {
	s.log = l
}
func (s *MyService) InjectListeners(ls []net.Listener) {
	s.l = ls[0]
}
func (s *MyService) Listeners() []net.Listener {
	return []net.Listener{s.l}
}

func main() {
	cfg := gmc.New.Config()
	cfg.Set("listen", ":")
	app := gmc.New.App()
	app.AttachConfig("mycfg", cfg)
	app.AddService(gmc.ServiceItem{
		Service:      NewMyService(),
		ConfigID: "mycfg",
	})
	app.Logger().Panic(app.Run())
}
