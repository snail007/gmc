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

func NewAny(initValue interface{}) *Any {
	return &Any{val: initValue}
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
		s.val = f(s.val)
	})
}

type Value struct {
	baseValue
	val atomic.Value
}

func NewValue(initValue interface{}) *Value {
	var val atomic.Value
	val.Store(initValue)
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

func NewString(initValue string) *String {
	return &String{val: initValue}
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

func NewBool(initValue bool) *Bool {
	a := new(int32)
	*a = 0
	if initValue {
		*a = 1
	}
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

func NewBytes(initValue []byte) *Bytes {
	var v []byte
	if len(initValue) > 0 {
		v = make([]byte, len(initValue))
		copy(v, initValue)
	}
	return &Bytes{bytes: v}
}

func (s *Bytes) Bytes() (x []byte) {
	s.get(func() {
		x = s.bytes
	})
	return
}

func (s *Bytes) CloneBytes() (x []byte) {
	s.get(func() {
		x = make([]byte, len(s.bytes))
		copy(x, s.bytes)
	})
	return
}

func (s *Bytes) Len() int {
	return len(s.Bytes())
}

func (s *Bytes) SetBytes(data []byte) {
	s.set(func() {
		s.bytes = make([]byte, len(data))
		copy(s.bytes, data)
	})
}

func (s *Bytes) Append(bytes []byte) {
	s.set(func() {
		s.bytes = append(s.bytes, bytes...)
	})
}

type Int8 struct {
	baseValue
	val int8
}

func NewInt8(initValue int8) *Int8 {
	return &Int8{val: initValue}
}

func (s *Int8) Increase(cnt int8) (newValue int8) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Int8) Decrease(cnt int8) (newValue int8) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Int8) SetVal(cnt int8) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Int8) Val() (value int8) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Int8) GetAndSet(f func(oldVal int8) (newVal int8)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Int struct {
	baseValue
	val int
}

func NewInt(initValue int) *Int {
	return &Int{val: initValue}
}

func (s *Int) Increase(cnt int) (newValue int) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Int) Decrease(cnt int) (newValue int) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Int) SetVal(cnt int) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Int) Val() (value int) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Int) GetAndSet(f func(oldVal int) (newVal int)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Int32 struct {
	baseValue
	val int32
}

func NewInt32(initValue int32) *Int32 {
	return &Int32{val: initValue}
}

func (s *Int32) Increase(cnt int32) (newValue int32) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Int32) Decrease(cnt int32) (newValue int32) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Int32) SetVal(cnt int32) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Int32) Val() (value int32) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Int32) GetAndSet(f func(oldVal int32) (newVal int32)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Int64 struct {
	baseValue
	val int64
}

func NewInt64(initValue int64) *Int64 {
	return &Int64{val: initValue}
}

func (s *Int64) Increase(cnt int64) (newValue int64) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Int64) Decrease(cnt int64) (newValue int64) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Int64) SetVal(cnt int64) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Int64) Val() (value int64) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Int64) GetAndSet(f func(oldVal int64) (newVal int64)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Uint8 struct {
	baseValue
	val uint8
}

func NewUint8(initValue uint8) *Uint8 {
	return &Uint8{val: initValue}
}

func (s *Uint8) Increase(cnt uint8) (newValue uint8) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Uint8) Decrease(cnt uint8) (newValue uint8) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Uint8) SetVal(cnt uint8) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Uint8) Val() (value uint8) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Uint8) GetAndSet(f func(oldVal uint8) (newVal uint8)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Uint struct {
	baseValue
	val uint
}

func NewUint(initValue uint) *Uint {
	return &Uint{val: initValue}
}

func (s *Uint) Increase(cnt uint) (newValue uint) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Uint) Decrease(cnt uint) (newValue uint) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Uint) SetVal(cnt uint) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Uint) Val() (value uint) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Uint) GetAndSet(f func(oldVal uint) (newVal uint)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Uint32 struct {
	baseValue
	val uint32
}

func NewUint32(initValue uint32) *Uint32 {
	return &Uint32{val: initValue}
}

func (s *Uint32) Increase(cnt uint32) (newValue uint32) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Uint32) Decrease(cnt uint32) (newValue uint32) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Uint32) SetVal(cnt uint32) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Uint32) Val() (value uint32) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Uint32) GetAndSet(f func(oldVal uint32) (newVal uint32)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Uint64 struct {
	baseValue
	val uint64
}

func NewUint64(initValue uint64) *Uint64 {
	return &Uint64{val: initValue}
}

func (s *Uint64) Increase(cnt uint64) (newValue uint64) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Uint64) Decrease(cnt uint64) (newValue uint64) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Uint64) SetVal(cnt uint64) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Uint64) Val() (value uint64) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Uint64) GetAndSet(f func(oldVal uint64) (newVal uint64)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Float32 struct {
	baseValue
	val float32
}

func NewFloat32(initValue float32) *Float32 {
	return &Float32{val: initValue}
}

func (s *Float32) Increase(cnt float32) (newValue float32) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Float32) Decrease(cnt float32) (newValue float32) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Float32) SetVal(cnt float32) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Float32) Val() (value float32) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Float32) GetAndSet(f func(oldVal float32) (newVal float32)) {
	s.set(func() {
		s.val = f(s.val)
	})
}

type Float64 struct {
	baseValue
	val float64
}

func NewFloat64(initValue float64) *Float64 {
	return &Float64{val: initValue}
}

func (s *Float64) Increase(cnt float64) (newValue float64) {
	s.set(func() {
		s.val += cnt
		newValue = s.val
	})
	return
}

func (s *Float64) Decrease(cnt float64) (newValue float64) {
	s.set(func() {
		s.val -= cnt
		newValue = s.val
	})
	return
}

func (s *Float64) SetVal(cnt float64) {
	s.set(func() {
		s.val = cnt
	})
}

func (s *Float64) Val() (value float64) {
	s.get(func() {
		value = s.val
	})
	return
}

func (s *Float64) GetAndSet(f func(oldVal float64) (newVal float64)) {
	s.set(func() {
		s.val = f(s.val)
	})
}
