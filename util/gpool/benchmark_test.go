// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gpool

import (
	gcast "github.com/snail007/gmc/util/cast"
	gloop "github.com/snail007/gmc/util/loop"
	"sync"
	"testing"
)

var pSubmit, pWorker *BasicPool

func BenchmarkSubmit(b *testing.B) {
	var do = func(max, step int) {
		gloop.ForBy(max, step, func(loopIndex, loopValue int) {
			size := loopValue + step
			b.Run("pool size:"+gcast.ToString(size), func(b *testing.B) {
				pSubmit = New(size)
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					pSubmit.Submit(func() {
					})
				}
				b.StopTimer()
				pSubmit.Stop()
			})
		})
	}
	do(100, 20)
	do(1000, 200)
	do(50000, 10000)

}
func BenchmarkJob(b *testing.B) {
	var do = func(max, step int) {
		gloop.ForBy(max, step, func(loopIndex, loopValue int) {
			size := loopValue + step
			b.Run("pool size:"+gcast.ToString(size), func(b *testing.B) {
				pWorker = New(size)
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
			})
		})
	}
	do(100, 20)
	do(1000, 200)
	do(50000, 10000)
}
