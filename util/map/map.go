// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmap

import (
	"fmt"
	"github.com/snail007/gmc/util/linklist"
	"sync"
)

type (
	// M alias of type map[string]interface{}
	M = map[string]interface{}
	// Mii alias of type map[interface{}]interface{}
	Mii = map[interface{}]interface{}
	// Mss alias of type map[string]string
	Mss = map[string]string
)

// Map a map can kept the sequence of keys when range the map
// have more useful function, Len(), Shift(), Pop(),gKeys(), etc.
type Map struct {
	keys     *glinklist.LinkList
	data     sync.Map
	keyElMap sync.Map
	sync.RWMutex
}

// Clone duplicates the map s.
func (s *Map) Clone() *Map {
	s.RLock()
	defer s.RUnlock()
	m := NewMap()
	p := s.keys.Front()
	for {
		if p == nil {
			break
		}
		v, _ := s.data.Load(p.Value)
		m.Store(p.Value, v)
		p = p.Next()
	}
	return m
}

// ToMap duplicates the map s.
func (s *Map) ToMap() map[interface{}]interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.toMap()
}
func (s *Map) toMap() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	s.data.Range(func(key, value interface{}) bool {
		m[key] = value
		return true
	})
	return m
}

// ToStringMap duplicates the map s.
func (s *Map) ToStringMap() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.toStringMap()
}

func (s *Map) toStringMap() map[string]interface{} {
	m := map[string]interface{}{}
	s.data.Range(func(key, value interface{}) bool {
		m[fmt.Sprintf("%v", key)] = value
		return true
	})
	return m
}

// Merge merges a Map to Map s.
func (s *Map) Merge(m *Map) {
	s.Lock()
	defer s.Unlock()
	s.merge(m)
}

func (s *Map) merge(m *Map) {
	m.data.Range(func(key, value interface{}) bool {
		s.store(key, value)
		return true
	})
}

// MergeMap merges a map to Map s.
func (s *Map) MergeMap(m Mii) {
	s.Lock()
	defer s.Unlock()
	s.mergeMap(m)
}

func (s *Map) mergeMap(m Mii) {
	for key, value := range m {
		s.store(key, value)
	}
}

// MergeStrMap merges a map to Map s.
func (s *Map) MergeStrMap(m M) {
	s.Lock()
	defer s.Unlock()
	s.mergeStrMap(m)
}

func (s *Map) mergeStrMap(m M) {
	for key, value := range m {
		s.store(key, value)
	}
}

// MergeStrStrMap merges a map to Map s.
func (s *Map) MergeStrStrMap(m Mss) {
	s.Lock()
	defer s.Unlock()
	s.mergeStrStrMap(m)
}

func (s *Map) mergeStrStrMap(m Mss) {
	for key, value := range m {
		s.store(key, value)
	}
}

// MergeSyncMap merges a sync.Map to Map s.
func (s *Map) MergeSyncMap(m *sync.Map) {
	s.Lock()
	defer s.Unlock()
	s.mergeSyncMap(m)
}

func (s *Map) mergeSyncMap(m *sync.Map) {
	m.Range(func(key, value interface{}) bool {
		s.store(key, value)
		return true
	})
}

// Pop returns the last element of map s or nil if the map is empty.
func (s *Map) Pop() (k, v interface{}, ok bool) {
	s.Lock()
	defer s.Unlock()
	return s.removeElement(s.keys.Back())
}

// Shift returns the first element of map s or nil if the map is empty.
func (s *Map) Shift() (k, v interface{}, ok bool) {
	s.Lock()
	defer s.Unlock()
	return s.removeElement(s.keys.Front())
}

// used for shift and pop
func (s *Map) removeElement(el *glinklist.Element) (k, v interface{}, ok bool) {
	if el == nil {
		return
	}
	v, ok = s.data.Load(el.Value)
	if ok {
		k = el.Value
		s.delete(el.Value)
	}
	return
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (s *Map) Load(key interface{}) (value interface{}, ok bool) {
	s.RLock()
	defer s.RUnlock()
	return s.load(key)
}

func (s *Map) load(key interface{}) (value interface{}, ok bool) {
	value, ok = s.data.Load(key)
	return
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (s *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	s.Lock()
	defer s.Unlock()
	return s.loadOrStore(key, value)
}

func (s *Map) loadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data.LoadOrStore(key, value)
	if !loaded {
		s.keyElMap.Store(key, s.keys.PushBack(key))
	}
	return
}

// LoadOrStoreFront returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
// The key will be stored the first in keys queue if key not exists.
func (s *Map) LoadOrStoreFront(key, value interface{}) (actual interface{}, loaded bool) {
	s.Lock()
	defer s.Unlock()
	return s.loadOrStoreFront(key, value)
}

func (s *Map) loadOrStoreFront(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data.LoadOrStore(key, value)
	if !loaded {
		s.keyElMap.Store(key, s.keys.PushFront(key))
	}
	return
}

// StoreFront sets the value for a key.
// The key will be stored the first in keys queue.
func (s *Map) StoreFront(key, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.storeFront(key, value)
}

func (s *Map) storeFront(key, value interface{}) {
	s.data.Store(key, value)
	if v, ok := s.keyElMap.Load(key); ok {
		s.keys.Remove(v.(*glinklist.Element))
	}
	s.keyElMap.Store(key, s.keys.PushFront(key))
}

// Store sets the value for a key.
func (s *Map) Store(key, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.store(key, value)
}

func (s *Map) store(key, value interface{}) {
	s.data.Store(key, value)
	if v, ok := s.keyElMap.Load(key); ok {
		s.keys.Remove(v.(*glinklist.Element))
	}
	s.keyElMap.Store(key, s.keys.PushBack(key))
}

// Delete deletes the value for a key.
func (s *Map) Delete(key interface{}) {
	s.Lock()
	defer s.Unlock()
	s.delete(key)
}

func (s *Map) delete(key interface{}) {
	s.data.Delete(key)
	if el, ok := s.keyElMap.Load(key); ok {
		s.keys.Remove(el.(*glinklist.Element))
		s.keyElMap.Delete(key)
	}
}

// Len returns the length of the map s.
func (s *Map) Len() int {
	return int(s.keys.Len())
}

// Clear deletes all data in the map s.
func (s *Map) Clear() {
	s.Lock()
	defer s.Unlock()
	s.clear()
}

func (s *Map) clear() {
	s.data = sync.Map{}
	s.keys = glinklist.New()
	s.keyElMap = sync.Map{}
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range keep the sequence of store sequence.
func (s *Map) Range(f func(key, value interface{}) bool) {
	s.keys.Clone().Range(func(k interface{}) bool {
		if v, ok := s.data.Load(k); ok {
			return f(k, v)
		}
		return true
	})
}

// RangeFast calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// RangeFast keep the sequence of store sequence.
//
// RangeFast do not create a snapshot for range, so you can not
// modify map s in range loop, indicate do not call Delete(), Store(), LoadOrStore(), Merge(), etc.
func (s *Map) RangeFast(f func(key, value interface{}) bool) {
	s.keys.Range(func(k interface{}) bool {
		_, v := s.data.Load(k)
		return f(k, v)
	})
}

// Keys returns all keys in map s and keep the sequence of store sequence.
func (s *Map) Keys() (keys []interface{}) {
	s.RLock()
	defer s.RUnlock()
	return s.keysArr()
}

func (s *Map) keysArr() (keyArr []interface{}) {
	p := s.keys.Front()
	for {
		if p == nil {
			break
		}
		keyArr = append(keyArr, p.Value)
		p = p.Next()
	}
	return
}

// StringKeys returns all keys in map s and keep the sequence of store sequence.
func (s *Map) StringKeys() (keys []string) {
	s.RLock()
	defer s.RUnlock()
	return s.stringKeys()
}

func (s *Map) stringKeys() (keys []string) {
	p := s.keys.Front()
	for {
		if p == nil {
			break
		}
		keys = append(keys, fmt.Sprintf("%v", p.Value))
		p = p.Next()
	}
	return
}

// IsEmpty indicates if the map is empty.
func (s *Map) IsEmpty() bool {
	return s.keys.Len() == 0
}

// IndexOf indicates the index of value in Map s, if not found returns -1.
//
// idx start with 0.
func (s *Map) IndexOf(k interface{}) int {
	return s.keys.IndexOf(k)
}

// String returns string format of the Set.
func (s *Map) String() string {
	return fmt.Sprintf("%v", s.ToMap())
}

// NewMap creates a Map object.
func NewMap() *Map {
	return &Map{
		keys:     glinklist.New(),
		data:     sync.Map{},
		keyElMap: sync.Map{},
	}
}
