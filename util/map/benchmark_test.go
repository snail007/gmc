// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmap

import (
	"sync"
	"testing"
)

func BenchmarkGMapLoad(b *testing.B) {
	gMap := NewMap()
	b.StopTimer()
	for i := 0; i < b.N/2; i++ {
		gMap.Store(i, Map{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gMap.Load(i)
	}
}

func BenchmarkGoMapLoad(b *testing.B) {
	goMap := sync.Map{}
	b.StopTimer()
	for i := 0; i < b.N/2; i++ {
		goMap.Store(i, Map{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		goMap.Load(i)
	}
}

func BenchmarkGMapLoadOrStore(b *testing.B) {
	gMap := NewMap()
	b.StopTimer()
	for i := 0; i < b.N/2; i++ {
		gMap.Store(i, Map{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gMap.LoadOrStore(i, Map{})
	}
}

func BenchmarkGoMapLoadOrStore(b *testing.B) {
	goMap := sync.Map{}
	b.StopTimer()
	for i := 0; i < b.N/2; i++ {
		goMap.Store(i, Map{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		goMap.LoadOrStore(i, Map{})
	}
}

func BenchmarkGMapRange(b *testing.B) {
	gMap := NewMap()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		gMap.Store(i, Map{})
	}
	b.StartTimer()
	gMap.Range(func(_, _ interface{}) bool {
		return true
	})
}

func BenchmarkGoMapRange(b *testing.B) {
	goMap := sync.Map{}
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		goMap.Store(i, Map{})
	}
	b.StartTimer()
	goMap.Range(func(_, _ interface{}) bool {
		return true
	})
}

func BenchmarkGMapPop(b *testing.B) {
	gMap := NewMap()
	b.StopTimer()
	for i := 0; i < b.N/2; i++ {
		gMap.Store(i, Map{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gMap.Pop()
	}
}

func BenchmarkGMapShift(b *testing.B) {
	gMap := NewMap()
	b.StopTimer()
	for i := 0; i < b.N/2; i++ {
		gMap.Store(i, Map{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gMap.Shift()
	}
}