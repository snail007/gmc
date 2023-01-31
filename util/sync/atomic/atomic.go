// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gatomic

import (
	"sync"
	"sync/atomic"
)

type baseValue struct {
	l sync.RWMutex
}

func (s *baseValue) get(f func()) {
	s.l.RLock()
	defer s.l.RUnlock()
	f()
}

func (s *baseValue) set(f func()) {
	s.l.Lock()
	defer s.l.Unlock()
	f()
}

type Any struct {
	baseValue
	val interface{}
}

func NewAny(x ...interface{}) *Any {
	var val interface{}
	if len(x) == 1 && x[0] != nil {
		val = x[0]
	}
	return &Any{val: val}
}

func (s *Any) SetVal(x interface{}) {
	s.set(func() {
		s.val = x
	})
}

func (s *Any) Val() (x interface{}) {
	s.get(func() {
		x = s.val
	})
	return
}

func (s *Any) GetAndSet(f func(oldVal interface{}) (newVal interface{})) {
	s.set(func() {
		f(s.val)
	})
}

type Value struct {
	baseValue
	val atomic.Value
}

func NewValue(x ...interface{}) *Value {
	var val atomic.Value
	if len(x) == 1 && x[0] != nil {
		val.Store(x[0])
	}
	return &Value{val: val}
}

func (s *Value) SetVal(x interface{}) {
	s.set(func() {
		s.val.Store(x)
	})
}

func (s *Value) Val() (x interface{}) {
	s.get(func() {
		x = s.val.Load()
	})
	return
}

func (s *Value) GetAndSet(f func(oldVal interface{}) (newVal interface{})) {
	s.set(func() {
		s.val.Store(f(s.val.Load()))
	})
}

type String struct {
	baseValue
	val string
}

func NewString(val ...string) *String {
	v := ""
	if len(val) == 1 {
		v = val[0]
	}
	return &String{val: v}
}

func (s *String) IsEmpty() bool {
	return len(s.Val()) == 0
}

func (s *String) Val() (x string) {
	s.get(func() {
		x = s.val
	})
	return
}

func (s *String) String() string {
	return s.Val()
}

func (s *String) SetVal(val string) {
	s.set(func() {
		s.val = val
	})
}

type Bool struct {
	val *int32
}

func NewBool() *Bool {
	a := new(int32)
	return &Bool{val: a}
}

func (s *Bool) IsTrue() bool {
	return atomic.LoadInt32(s.val) == 1
}

func (s *Bool) IsFalse() bool {
	return atomic.LoadInt32(s.val) == 0
}

func (s *Bool) SetTrue() {
	atomic.StoreInt32(s.val, 1)
}

func (s *Bool) SetFalse() {
	atomic.StoreInt32(s.val, 0)
}

type Bytes struct {
	baseValue
	bytes []byte
}

func NewBytes(val ...[]byte) *Bytes {
	var v []byte
	if len(val) == 1 {
		v = val[0]
	}
	return &Bytes{bytes: v}
}

func (s *Bytes) Bytes() (x []byte) {
	s.get(func() {
		x = s.bytes
	})
	return
}

func (s *Bytes) Len() int {
	return len(s.Bytes())
}

func (s *Bytes) SetBytes(data []byte) {
	s.set(func() {
		s.bytes = data
	})
}

func (s *Bytes) Append(bytes []byte) {
	s.set(func() {
		s.bytes = append(s.bytes, bytes...)
	})
}
