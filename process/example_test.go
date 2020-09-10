package process

import (
	"fmt"

	"github.com/snail007/gmc/process/daemon"
	"github.com/snail007/gmc/process/hook"
)

func Example() {
	//the below code block should be replace into your main function code, and your actual main() change to doMain()
	if err := daemon.Start(); err != nil {
		fmt.Println(err)
		return
	}
	if daemon.CanRun() {
		go doMain()
	}
	hook.RegistShutdown(func() {
		daemon.Clean()
	})
	hook.WaitShutdown()
}

//do your main function
func doMain() {
	//flag parse should be called here if you have
}
