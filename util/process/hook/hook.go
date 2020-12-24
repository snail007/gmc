// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package ghook

import (
	"fmt"
	"github.com/snail007/gmc/core"
	gerr "github.com/snail007/gmc/module/error"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	shutdown = []func(){}
	IsMock   = false
)

func RegistShutdown(fn func()) {
	shutdown = append(shutdown, fn)
}
func WaitShutdown() {
	defer gcore.Recover(func(e interface{}) {
		fmt.Printf("shutdown hook manager crashed, err: %s", gerr.Stack(e))
	})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-signalChan
	log.Println("Received an interrupt, stopping service...")
	runHooks()
	os.Exit(0)
}

//just for testing
func MockShutdown() {
	log.Println("Received an interrupt, mock stopping service...")
	runHooks()
}
func runHooks() {
	caller := func(fn func()) {
		defer gcore.Recover(func(e interface{}) {
			fmt.Printf("shutdown hook crashed, err: %s", gerr.Stack(e))
		})
		fn()
	}
	for _, fn := range shutdown {
		caller(fn)
	}
}
