// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gset

import (
	"fmt"
	"github.com/snail007/gmc/util/collection"
	gmap "github.com/snail007/gmc/util/map"
)

// Set each value in set is unique, and keep the sequence of Add() or Merge()
type Set struct {
	data *gmap.Map
}

// Add adds a value to set, if value not exists in set s.
func (s *Set) Add(value ...interface{}) {
	for _, v := range value {
		s.data.LoadOrStore(v, nil)
	}
}

// AddFront adds a value to set, if value not exists in set s.
// The value will be stored the first in set s if value not exists.
func (s *Set) AddFront(value ...interface{}) {
	for _, v := range value {
		s.data.LoadOrStoreFront(v, nil)
	}
}

// Delete removes a value from set.
func (s *Set) Delete(value ...interface{}) {
	for _, k := range value {
		s.data.Delete(k)
	}
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

// ToStringSlice convert set s to slice []interface{}.
//
// Kept the Add() sequence.
func (s *Set) ToStringSlice() []string {
	return s.data.StringKeys()
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

// MergeStringSlice merges slice values to set s.
//
// Kept the Add() sequence.
func (s *Set) MergeStringSlice(data []string) {
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
	set := New()
	set.Merge(s)
	return set
}

// CloneAndClear duplicates the set s.
func (s *Set) CloneAndClear() *Set {
	set := New()
	set.MergeSlice(s.data.CloneAndClear().Keys())
	return set
}

// Clear duplicates the set s.
func (s *Set) Clear() *Set {
	s.data.Clear()
	return s
}

// IndexOf indicates the index of value in Map s, if not found returns -1.
//
// idx start with 0.
func (s *Set) IndexOf(k interface{}) int {
	return s.data.IndexOf(k)
}

// String returns string format of the Set.
func (s *Set) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

// CartesianProduct → {(x, y) | ∀ x ∈ s1, ∀ y ∈ s2}
// For example:
//		CartesianProduct(A, B), where A = {1, 2} and B = {7, 8}
//        => {(1, 7), (1, 8), (2, 7), (2, 8)}
func (s *Set) CartesianProduct(sets ...*Set) [][]interface{} {
	var dstSets [][]interface{}
	dstSets = append(dstSets, s.ToSlice())
	for _, set := range sets {
		if set.Len() == 0 {
			continue
		}
		dstSets = append(dstSets, set.ToSlice())
	}
	return collection.CartesianProduct(dstSets...)
}

// New returns a new *Set object
func New() *Set {
	return &Set{data: gmap.New()}
}
