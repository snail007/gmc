// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmcdaemon

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Example() {
	//the code block should be in your main function,and the actual main() change to doMain()
	if err := Start(); err != nil {
		fmt.Println(err)
		return
	}
	if CanRun() {
		//call actual main()
		go doMain()
	}
	//waiting kill signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	//do clean
	Clean()
}

//your main function
func doMain() {
	//flag.Parse() should be called here if you have
}
