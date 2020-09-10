package gpool

import (
	"sync"
	"testing"
)

func BenchmarkSubmit(b *testing.B) {
	b.StopTimer()
	p := New(3)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func() {
		})
	}
	b.StopTimer()
	p.Stop()
}
func BenchmarkWorker(b *testing.B) {
	b.StopTimer()
	p := New(4)
	g := sync.WaitGroup{}
	g.Add(b.N)
	for i := 0; i < b.N; i++ {
		p.Submit(func() {
			g.Done()
		})
	}
	b.StartTimer()
	g.Wait()
}
