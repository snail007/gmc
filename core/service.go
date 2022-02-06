// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcore

import (
	"net"
)

type Service interface {
	// init servcie
	Init(cfg Config) error
	//nonblocking, called After Init -> InjectListeners (when reload) -> Start
	Start() error
	Stop()
	// blocking until all resource are released
	GracefulStop()
	SetLog(log Logger)
	// called After Init
	InjectListeners([]net.Listener)
	ListenerFactory() func() (net.Listener, error)
	SetListenerFactory(listenerFactory func() (net.Listener, error))
	Listeners() []net.Listener
}

type ServiceItem struct {
	BeforeInit func(srv Service, cfg Config) (err error)
	AfterInit  func(srv *ServiceItem) (err error)
	Service    Service
	ConfigID   string
}
