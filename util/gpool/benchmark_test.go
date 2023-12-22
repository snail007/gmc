// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool

import (
	"sync"
	"testing"
)

var pSubmit, pWorker *Pool

func BenchmarkSubmit(b *testing.B) {
	pSubmit = New(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pSubmit.Submit(func() {
		})
	}
	b.StopTimer()
	pSubmit.Stop()
}
func BenchmarkWorker(b *testing.B) {
	pWorker = New(1)
	b.StopTimer()
	g := sync.WaitGroup{}
	g.Add(b.N)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pWorker.Submit(func() {
			g.Done()
		})
	}
	g.Wait()
	b.StopTimer()
	pWorker.Stop()
}
