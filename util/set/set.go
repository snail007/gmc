package gset

import gmap "github.com/snail007/gmc/util/map"

// Set each value in set is unique, and keep the sequence of Add() or Merge()
type Set struct {
	data *gmap.Map
}

// Add adds a value to set, if value not exists in set s.
func (s *Set) Add(value interface{}) {
	s.data.LoadOrStore(value, nil)
}

// Delete removes a value from set.
func (s *Set) Delete(value interface{}) {
	s.data.Delete(value)
}

// Contains returns true if value exists in set s, otherwise false.
func (s *Set) Contains(value interface{}) (exists bool) {
	_, exists = s.data.Load(value)
	return
}

// Range ranges the value in set s, if function f return false, the range loop will break.
//
// Range keeps the Add() sequence.
func (s *Set) Range(f func(value interface{}) bool) {
	s.data.Range(func(key, _ interface{}) bool {
		return f(key)
	})
	return
}

// RangeFast ranges the value in set s, if function f return false, the range loop will break.
//
// RangeFast do not create a snapshot for range, so you can not modify list s in range loop,
// indicate do not call Remove(), Add(), Merge() etc.
func (s *Set) RangeFast(f func(value interface{}) bool) {
	s.data.RangeFast(func(key, _ interface{}) bool {
		return f(key)
	})
	return
}

// ToSlice convert set s to slice []interface{}.
//
// Kept the Add() sequence.
func (s *Set) ToSlice() []interface{} {
	return s.data.Keys()
}

// Merge merges set values to set s.
//
// Kept the Add() sequence.
func (s *Set) Merge(set *Set) {
	s.MergeSlice(set.data.Keys())
}

// MergeSlice merges slice values to set s.
//
// Kept the Add() sequence.
func (s *Set) MergeSlice(data []interface{}) {
	for _, v := range data {
		s.data.LoadOrStore(v, nil)
	}
}

// Pop returns last value of set s, and remove it.
func (s *Set) Pop() (value interface{}) {
	value, _, _ = s.data.Pop()
	return
}

// Shift returns first value of set s, and remove it.
func (s *Set) Shift() (value interface{}) {
	value, _, _ = s.data.Shift()
	return
}

// Len returns length of the set s.
func (s *Set) Len() int {
	return s.data.Len()
}

// Clone duplicates the set s.
func (s *Set) Clone() *Set {
	set := NewSet()
	set.Merge(s)
	return set
}

// NewSet returns a new *Set object
func NewSet() *Set {
	return &Set{data: gmap.NewMap()}
}
