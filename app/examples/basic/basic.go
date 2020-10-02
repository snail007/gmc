package main

import (
	"log"
	"net"

	"github.com/snail007/gmc/util/logutil"

	"github.com/snail007/gmc"
)

type MyService struct {
	gmc.Service
	l       net.Listener
	log     *log.Logger
	address string
}

func NewMyService() gmc.Service {
	return &MyService{}
}
func (s *MyService) Init(cfg *gmc.Config) error {
	s.address = cfg.GetString("listen")
	s.log = logutil.New("")
	return nil
}
func (s *MyService) Start() (err error) {
	if s.l == nil {
		s.l, err = net.Listen("tcp", s.address)
	}
	// do something
	s.log.Printf("server listen on %s", s.l.Addr().String())
	return
}
func (s *MyService) Stop() {
	s.log.Printf("server stoped on %s", s.l.Addr().String())
	// 1. close active connections ...
	// 2. stop accept
	s.l.Close()
}
func (s *MyService) GracefulStop() {
	s.log.Printf("server graceful stop on %s", s.l.Addr().String())
	// 1. stop accept
	s.l.Close()
}
func (s *MyService) SetLog(l *log.Logger) {
	s.log = l
}
func (s *MyService) InjectListeners(ls []net.Listener) {
	s.l = ls[0]
}
func (s *MyService) Listeners() []net.Listener {
	return []net.Listener{s.l}
}

func main() {
	cfg := gmc.NewConfig()
	cfg.Set("listen", ":")
	app := gmc.NewAPP()
	app.AttachConfig("mycfg", cfg)
	app.AddService(gmc.ServiceItem{
		Service:      NewMyService(),
		ConfigIDname: "mycfg",
	})
	app.Logger().Panic(app.Run())
}
