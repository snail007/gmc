// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gbytes

import (
	"fmt"
	"sync"
)

var (
	poolMap     = map[int]*Pool{}
	poolLock    = &sync.Mutex{}
	poolCapMap  = map[string]*Pool{}
	poolCapLock = &sync.Mutex{}
)

func GetPool(bufSize int) *Pool {
	poolLock.Lock()
	defer poolLock.Unlock()
	if v, ok := poolMap[bufSize]; ok {
		return v
	}
	p := NewPool(bufSize)
	poolMap[bufSize] = p
	return p
}

func GetPoolCap(bufSize, capSize int) *Pool {
	poolLock.Lock()
	defer poolLock.Unlock()
	k := fmt.Sprintf("%d-%d", bufSize, capSize)
	if v, ok := poolCapMap[k]; ok {
		return v
	}
	p := NewPoolCap(bufSize, capSize)
	poolCapMap[k] = p
	return p
}

type Pool struct {
	*sync.Pool
}

func (s *Pool) Get() interface{} {
	return s.Pool.Get()
}

func (s *Pool) Put(x interface{}) {
	s.Pool.Put(x)
}

func NewPool(bufSize int) *Pool {
	p := &sync.Pool{}
	p.New = func() interface{} {
		return make([]byte, bufSize)
	}
	return &Pool{
		Pool: p,
	}
}

func NewPoolCap(bufSize, capSize int) *Pool {
	p := &sync.Pool{}
	p.New = func() interface{} {
		return make([]byte, bufSize, capSize)
	}
	return &Pool{
		Pool: p,
	}
}
