// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package main

import (
	gcore "github.com/snail007/gmc/core"
	"net"

	"github.com/snail007/gmc"
)

type MyService struct {
	gcore.Service
	l       net.Listener
	log     gcore.Logger
	address string
}

func NewMyService() gcore.Service {
	return &MyService{}
}
func (s *MyService) Init(cfg gcore.Config) error {
	s.address = cfg.GetString("listen")
	s.log = gcore.ProviderLogger()(nil, "")
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
	app.AddService(gcore.ServiceItem{
		Service:  NewMyService(),
		ConfigID: "mycfg",
	})
	app.Logger().Panic(app.Run())
}
