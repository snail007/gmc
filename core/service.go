package gmccore

import (
	"net"

	gmcconfig "github.com/snail007/gmc/config"
)

type Service interface {
	// init servcie
	Init(cfg *gmcconfig.Config) error
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
