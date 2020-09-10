## Daemon 

daemon package does add daemon, logging and forever function to your program.

the below command line arguments you can passed to your program after you using daemon package.

--forver or -forver

the argument will fork a worker process and master process monit worker process , and restart it when it crashed.

--daemon or -daemon

the argument will put your program running in background , only working in linux uinx etc.

--flog or -flog <filename.log>

the argument will logging your program stdout to the log file.

notice:

before you maybe execute your program like this:

./foobar -u root

after using daemon package, execute your program can be like this:

./foobar -u root -forever -daemon -flog foobar.log

## Demo

```golang
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

```