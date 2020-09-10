package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/snail007/gmc/process/daemon"
)

func Example() {
	//the code block should be in your main function,and the actual main() change to doMain()
	if err := daemon.Start(); err != nil {
		fmt.Println(err)
		return
	}
	if daemon.CanRun() {
		//call actual main()
		go doMain()
	}
	//waiting kill signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	//do clean
	daemon.Clean()
}

//your main function
func doMain() {
	//flag.Parse() should be called here if you have
}
