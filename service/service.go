package gmcservice

import (
	"log"
	"net"

	gmcconfig "github.com/snail007/gmc/config/gmc"
)

type Service interface {
	// init servcie
	Init(cfg *gmcconfig.GMCConfig) error
	//nonblocking, called After Init -> InjectListeners (when reload) -> Start
	Start() error
	Stop()
	// blocking until all resource are released
	GracefulStop()
	SetLog(*log.Logger)
	// called After Init
	InjectListeners([]net.Listener)
	Listeners() []net.Listener
}
