// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

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
