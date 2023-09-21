package gsync

import "sync"

func Wait(g *sync.WaitGroup) <-chan bool {
	ch := make(chan bool)
	go func() {
		g.Wait()
		ch <- true
	}()
	return ch
}
