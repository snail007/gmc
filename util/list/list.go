// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glist

import (
	"fmt"
	"sync"
)

// List then list container
type List struct {
	lock *sync.RWMutex
	data []interface{}
}

// Len returns length of the list s.
func (s *List) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.data)
}

// Get gets a value on idx in list s.
func (s *List) Get(idx int) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if idx < 0 || idx >= len(s.data) {
		return nil
	}
	return s.data[idx]
}

// Set sets a value on idx in list s.
func (s *List) Set(idx int, v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if idx < 0 || idx >= len(s.data) {
		return
	}
	s.data[idx] = v
}

// Remove removes a value on idx in list s.
func (s *List) Remove(idx int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	length := len(s.data)
	if idx < 0 || idx > length-1 {
		return
	}
	data := s.data[0:idx]
	if idx+1 <= length-1 {
		data = append(data, s.data[idx+1:]...)
	}
	s.data = data
}

// Add adds a value to end of list s.
func (s *List) Add(v ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = append(s.data, v...)
}

// AddFront adds a value to first of list s.
func (s *List) AddFront(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	arr := []interface{}{v}
	s.data = append(arr, s.data...)
}

// MergeSlice merge a array slice to end of list s.
func (s *List) MergeSlice(arr []interface{}) *List {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = append(s.data, arr...)
	return s
}

// MergeStringSlice merge a array slice to end of list s.
func (s *List) MergeStringSlice(arr []string) *List {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, v := range arr {
		s.data = append(s.data, v)
	}
	return s
}

// Merge merges a List to end of list s.
func (s *List) Merge(l *List) *List {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = append(s.data, l.data...)
	return s
}

// Pop returns last value of list s, and remove it.
func (s *List) Pop() (v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	length := len(s.data)
	if length == 0 {
		return nil
	}
	v = s.data[length-1]
	if length == 1 {
		s.data = []interface{}{}
	} else {
		s.data = s.data[:length-1]
	}
	return
}

// Shift returns first value of list s, and remove it.
func (s *List) Shift() (v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.data) == 0 {
		return nil
	}
	v = s.data[0]
	if len(s.data) == 1 {
		s.data = []interface{}{}
	} else {
		s.data = s.data[1:]
	}
	return
}

// Range ranges the value in list s, if function f return false, the range loop will break.
func (s *List) Range(f func(index int, value interface{}) bool) {
	snapshot := s.Clone()
	for idx, v := range snapshot.data {
		if !f(idx, v) {
			break
		}
	}
	return
}

// RangeFast ranges the value in list s, if function f return false, the range loop will break.
//
// RangeFast do not create a snapshot for range, so you can not modify list s in range loop,
// indicate do not call Delete(), Add(), Merge() etc.
func (s *List) RangeFast(f func(index int, value interface{}) bool) {
	for idx, v := range s.data {
		if !f(idx, v) {
			break
		}
	}
	return
}

// Clone duplicates the list s.
func (s *List) Clone() *List {
	s.lock.RLock()
	defer s.lock.RUnlock()
	newList := New()
	newList.data = append(newList.data, s.data...)
	return newList
}

// ToSlice returns the list as an array.
func (s *List) ToSlice() []interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	var newList []interface{}
	newList = append(newList, s.data...)
	return newList
}

// ToStringSlice returns the list as an array.
func (s *List) ToStringSlice() []string {
	var newList []string
	for _, v := range s.ToSlice() {
		newList = append(newList, fmt.Sprintf("%v", v))
	}
	return newList
}

// Clear empty the list s.
func (s *List) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = []interface{}{}
}

// Sub returns the sub list between two ranges
// With x the start and y the end, we are sending [x,y) sublist.
func (s *List) Sub(start, end int) *List {
	s.lock.RLock()
	defer s.lock.RUnlock()
	length := len(s.data)
	if start >= end || start < 0 || end > length {
		return nil
	}
	l := New()
	l.MergeSlice(s.data[start:end])
	return l
}

// Contains  indicates if the list contains an element.
func (s *List) Contains(v interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, val := range s.data {
		if val == v {
			return true
		}
	}
	return false
}

// IsEmpty indicates if the list is empty.
func (s *List) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.data) == 0
}

// IndexOf indicates the first index of value in list s, if not found returns -1.
//
// idx start with 0.
func (s *List) IndexOf(v interface{}) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for i := 0; i < len(s.data); i++ {
		if s.data[i] == v {
			return i
		}
	}
	//j := len(s.data) - 1
	//length := len(s.data)
	//for i := 0; i <= length/2; i++ {
	//	if i < length && s.data[i] == v {
	//		return i
	//	}
	//	if j >= 0 && s.data[j] == v {
	//		return j
	//	}
	//	j--
	//}
	return -1
}

// String returns string format of the list.
func (s *List) String() string {
	return fmt.Sprintf("%v", s.ToSlice())
}

// New returns a new *List object
func New() *List {
	return &List{
		data: []interface{}{},
		lock: &sync.RWMutex{},
	}
}
