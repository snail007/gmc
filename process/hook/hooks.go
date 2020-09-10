package hook

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
)

var (
	shutdown = []func(){}
)

func RegistShutdown(fn func()) {
	shutdown = append(shutdown, fn)
}
func WaitShutdown() {
	signalChan := make(chan os.Signal, 1)
	g := sync.WaitGroup{}
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	g.Add(1)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Printf("shutdown hook manager crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
			}
			g.Done()
		}()
		for range signalChan {
			log.Println("Received an interrupt, stopping service...")
			for _, fn := range shutdown {
				caller := func(fn func()) {
					defer func() {
						if e := recover(); e != nil {
							fmt.Printf("shutdown hook crashed, err: %s\nstack:\n%s", e, string(debug.Stack()))
						}
					}()
					fn()
				}
				caller(fn)
			}
			return
		}
	}()
	g.Wait()
	os.Exit(0)
}
