package gcore

import (
	"net"

	gconfig "github.com/snail007/gmc/config"
)

type Service interface {
	// init servcie
	Init(cfg *gconfig.Config) error
	//nonblocking, called After Init -> InjectListeners (when reload) -> Start
	Start() error
	Stop()
	// blocking until all resource are released
	GracefulStop()
	SetLog(log Logger)
	// called After Init
	InjectListeners([]net.Listener)
	Listeners() []net.Listener
}

type ServiceItem interface {
	AfterInit() func(srv ServiceItem) (err error)
	BeforeInit() func(srv Service, cfg *gconfig.Config) (err error)
	SetBeforeInit(func(srv Service, cfg *gconfig.Config) (err error)) ServiceItem
	SetAfterInit(func(srv ServiceItem) (err error)) ServiceItem
	GetService() Service
	GetConfigID() string
}
