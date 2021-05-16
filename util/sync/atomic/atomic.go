package gatomic

import (
	"sync"
	"sync/atomic"
)

type Value struct {
	atomic.Value
	sync.RWMutex
}

func NewValue(x interface{}) *Value {
	var val atomic.Value
	val.Store(x)
	return &Value{Value:val}
}

func (s *Value) Store(x interface{}) {
	s.RLock()
	defer s.RUnlock()
	s.Value.Store(x)
}

func (s *Value) Load() interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.Value.Load()
}

func (s *Value) LoadAndStore(f func(x interface{}) interface{}) {
	s.Lock()
	defer s.Unlock()
	s.Store(f(s.Load()))
}
