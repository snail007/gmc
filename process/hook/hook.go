// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmchook

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
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
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("shutdown hook manager crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
		}
	}()
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
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("shutdown hook crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
			}
		}()
		fn()
	}
	for _, fn := range shutdown {
		caller(fn)
	}
}
