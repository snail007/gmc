package glist

import (
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
func (s *List) Add(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = append(s.data, v)
}

// AddFirst adds a value to first of list s.
func (s *List) AddFirst(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	arr := []interface{}{v}
	s.data = append(arr, s.data...)
}

// MergeSlice merge a array slice to end of list s.
func (s *List) MergeSlice(arr []interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = append(s.data, arr...)
}

// Merge merges a List to end of list s.
func (s *List) Merge(l *List) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = append(s.data, l.data...)
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

// Range ranges the value in list s, if function f return false, the range loop will break.
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
	newList := NewList()
	newList.data = append(newList.data, s.data...)
	return newList
}

// ToSlice returns the list as an array.
func (s *List) ToSlice() []interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	newList := []interface{}{}
	newList = append(newList, s.data...)
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
	l := NewList()
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

// NewList returns a new *List object
func NewList() *List {
	return &List{
		data: []interface{}{},
		lock: &sync.RWMutex{},
	}
}
