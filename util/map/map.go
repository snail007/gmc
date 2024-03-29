// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmap

import (
	"fmt"
	"sync"

	"github.com/snail007/gmc/util/linklist"
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
// have more useful function, Len(), Shift(), Pop(),Keys(), etc.
type Map struct {
	keys     *glinklist.LinkList
	data     map[interface{}]interface{}
	keyElMap map[interface{}]interface{}
	l        sync.RWMutex
}

// Clone duplicates the map s.
func (s *Map) Clone() *Map {
	s.l.RLock()
	defer s.l.RUnlock()
	m := New()
	s.keys.Range(func(v interface{}) bool {
		m.store(v, s.data[v])
		return true
	})
	return m
}

// GC rebuild the internal map to release memory used by the map.
func (s *Map) GC() {
	s.l.Lock()
	defer s.l.Unlock()
	m := New()
	s.keys.Range(func(v interface{}) bool {
		m.store(v, s.data[v])
		return true
	})
	s.clear()
	s.merge(m)
}

// CloneAndClear duplicates the map s.
func (s *Map) CloneAndClear() *Map {
	s.l.RLock()
	defer s.l.RUnlock()
	m := New()
	s.keys.Range(func(v interface{}) bool {
		m.store(v, s.data[v])
		return true
	})
	s.clear()
	return m
}

// ToMap duplicates the map s.
func (s *Map) ToMap() map[interface{}]interface{} {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.toMap()
}
func (s *Map) toMap() map[interface{}]interface{} {
	m := map[interface{}]interface{}{}
	for k, v := range s.data {
		m[k] = v
	}
	return m
}

// ToStringMap duplicates the map s.
func (s *Map) ToStringMap() map[string]interface{} {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.toStringMap()
}

func (s *Map) toStringMap() map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range s.data {
		m[fmt.Sprintf("%v", k)] = v
	}
	return m
}

// Merge merges a Map to Map s.
func (s *Map) Merge(m *Map) *Map {
	s.l.Lock()
	defer s.l.Unlock()
	s.merge(m)
	return s
}

func (s *Map) merge(m *Map) {
	for k, v := range m.toMap() {
		s.store(k, v)
	}
}

// MergeMap merges a map to Map s.
func (s *Map) MergeMap(m Mii) *Map {
	s.l.Lock()
	defer s.l.Unlock()
	s.mergeMap(m)
	return s
}

func (s *Map) mergeMap(m Mii) {
	for key, value := range m {
		s.store(key, value)
	}
}

// MergeStrMap merges a map to Map s.
func (s *Map) MergeStrMap(m M) *Map {
	s.l.Lock()
	defer s.l.Unlock()
	s.mergeStrMap(m)
	return s
}

func (s *Map) mergeStrMap(m M) {
	for key, value := range m {
		s.store(key, value)
	}
}

// MergeStrStrMap merges a map to Map s.
func (s *Map) MergeStrStrMap(m Mss) *Map {
	s.l.Lock()
	defer s.l.Unlock()
	s.mergeStrStrMap(m)
	return s
}

func (s *Map) mergeStrStrMap(m Mss) {
	for key, value := range m {
		s.store(key, value)
	}
}

// MergeSyncMap merges a sync.Map to Map s.
func (s *Map) MergeSyncMap(m *sync.Map) *Map {
	s.l.Lock()
	defer s.l.Unlock()
	s.mergeSyncMap(m)
	return s
}

func (s *Map) mergeSyncMap(m *sync.Map) {
	m.Range(func(key, value interface{}) bool {
		s.store(key, value)
		return true
	})
}

// Pop returns the last element of map s or nil if the map is empty.
func (s *Map) Pop() (k, v interface{}, ok bool) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.removeElement(s.keys.Back())
}

// Shift returns the first element of map s or nil if the map is empty.
func (s *Map) Shift() (k, v interface{}, ok bool) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.removeElement(s.keys.Front())
}

// used for shift and pop
func (s *Map) removeElement(el *glinklist.Element) (k, v interface{}, ok bool) {
	if el == nil {
		return
	}
	v, ok = s.data[el.Value]
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
	s.l.RLock()
	defer s.l.RUnlock()
	return s.load(key)
}

func (s *Map) load(key interface{}) (value interface{}, ok bool) {
	value, ok = s.data[key]
	return
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (s *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.loadOrStore(key, value)
}

func (s *Map) loadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data[key]
	if !loaded {
		actual = value
		s.keyElMap[key] = s.keys.PushBack(key)
		s.store(key, actual)
	}
	return
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (s *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
	s.l.Lock()
	defer s.l.Unlock()
	value, loaded = s.load(key)
	if loaded {
		s.delete(key)
	}
	return
}

// LoadAndStoreFunc call the given func and stores it returns value.
// The loaded result is true if the keys was exists, false if not exists.
// If loaded, the given func firstly parameter is the loaded value, otherwise is nil.
func (s *Map) LoadAndStoreFunc(key interface{}, f func(oldValue interface{}, loaded bool) (newValue interface{})) (newValue interface{}, loaded bool) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.loadAndStoreFunc(key, f)
}

func (s *Map) loadAndStoreFunc(key interface{}, f func(oldValue interface{}, loaded bool) (newValue interface{})) (newValue interface{}, loaded bool) {
	oldValue, loaded := s.data[key]
	newValue = f(oldValue, loaded)
	s.store(key, newValue)
	return
}

// LoadAndStoreFuncErr call the given func and stores it returns value,
// if func returns an error, the error be returned.
// The loaded result is true if the keys was exists, false if not exists.
// If loaded, the given func firstly parameter is the loaded value, otherwise is nil.
func (s *Map) LoadAndStoreFuncErr(key interface{}, f func(oldValue interface{}, loaded bool) (newValue interface{}, err error)) (newValue interface{}, loaded bool, err error) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.loadAndStoreFuncErr(key, f)
}

func (s *Map) loadAndStoreFuncErr(key interface{}, f func(oldValue interface{}, loaded bool) (newValue interface{}, err error)) (newValue interface{}, loaded bool, err error) {
	oldValue, loaded := s.data[key]
	newValue, err = f(oldValue, loaded)
	if err != nil {
		return
	}
	s.store(key, newValue)
	return
}

// LoadOrStoreFunc returns the existing value for the key if present.
// Otherwise, it call the given func and stores it returns value.
// The loaded result is true if the value was loaded, false if stored.
func (s *Map) LoadOrStoreFunc(key interface{}, f func() interface{}) (actual interface{}, loaded bool) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.loadOrStoreFunc(key, f)
}

func (s *Map) loadOrStoreFunc(key interface{}, f func() interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data[key]
	if !loaded {
		actual = f()
		s.store(key, actual)
	}
	return
}

// LoadOrStoreFuncErr returns the existing value for the key if present.
// Otherwise, it call the given func and stores it returns value, if the func
// returns an error, nothing will be stored, the function's error be returned.
// The loaded result is true if the value was loaded, false if stored.
func (s *Map) LoadOrStoreFuncErr(key interface{}, f func() (x interface{}, err error)) (actual interface{}, loaded bool, err error) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.loadOrStoreFuncErr(key, f)
}

func (s *Map) loadOrStoreFuncErr(key interface{}, f func() (x interface{}, err error)) (actual interface{}, loaded bool, err error) {
	actual, loaded = s.data[key]
	if !loaded {
		actual, err = f()
		if err != nil {
			return
		}
		s.store(key, actual)
	}
	return
}

// LoadOrStoreFront returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
// The key will be stored the first in keys queue if key not exists.
func (s *Map) LoadOrStoreFront(key, value interface{}) (actual interface{}, loaded bool) {
	s.l.Lock()
	defer s.l.Unlock()
	return s.loadOrStoreFront(key, value)
}

func (s *Map) loadOrStoreFront(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data[key]
	if !loaded {
		s.storeFront(key, value)
	}
	return
}

// StoreFront sets the value for a key.
// The key will be stored the first in keys queue.
func (s *Map) StoreFront(key, value interface{}) {
	s.l.Lock()
	defer s.l.Unlock()
	s.storeFront(key, value)
}

func (s *Map) storeFront(key, value interface{}) {
	s.data[key] = value
	if v, ok := s.keyElMap[key]; ok {
		s.keys.Remove(v.(*glinklist.Element))
	}
	s.keyElMap[key] = s.keys.PushFront(key)
}

// Store sets the value for a key.
func (s *Map) Store(key, value interface{}) {
	s.l.Lock()
	defer s.l.Unlock()
	s.store(key, value)
}

func (s *Map) store(key, value interface{}) {
	s.data[key] = value
	if v, ok := s.keyElMap[key]; ok {
		s.keys.Remove(v.(*glinklist.Element))
	}
	s.keyElMap[key] = s.keys.PushBack(key)
}

// Delete deletes the value for a key.
func (s *Map) Delete(key interface{}) {
	s.l.Lock()
	defer s.l.Unlock()
	s.delete(key)
}

func (s *Map) delete(key interface{}) {
	delete(s.data, key)
	if el, ok := s.keyElMap[key]; ok {
		s.keys.Remove(el.(*glinklist.Element))
		delete(s.keyElMap, key)
	}
}

// Len returns the length of the map s.
func (s *Map) Len() int {
	return int(s.keys.Len())
}

// Clear deletes all data in the map s.
func (s *Map) Clear() {
	s.l.Lock()
	defer s.l.Unlock()
	s.clear()
}

func (s *Map) clear() {
	s.data = map[interface{}]interface{}{}
	s.keys = glinklist.New()
	s.keyElMap = map[interface{}]interface{}{}
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range keep the sequence of store sequence.
func (s *Map) Range(f func(key, value interface{}) bool) {
	s.l.RLock()
	keys := s.keysArr()
	s.l.RUnlock()
	for _, k := range keys {
		v, ok := s.Load(k)
		if ok && !f(k, v) {
			break
		}
	}
}

// RangeFast calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// RangeFast keep the sequence of store sequence.
//
// RangeFast do not create a snapshot for range, so you can not
// modify map s in range loop, indicate do not call Delete(), Store(), LoadOrStore(), Merge(), etc.
func (s *Map) RangeFast(f func(key, value interface{}) bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	s.keys.Range(func(k interface{}) bool {
		if v, ok := s.load(k); ok {
			return f(k, v)
		}
		return true
	})
}

// Keys returns all keys in map s and keep the sequence of store sequence.
func (s *Map) Keys() (keys []interface{}) {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.keysArr()
}

func (s *Map) keysArr() (keyArr []interface{}) {
	s.keys.Range(func(v interface{}) bool {
		keyArr = append(keyArr, v)
		return true
	})
	return
}

// StringKeys returns all keys in map s and keep the sequence of store sequence.
func (s *Map) StringKeys() (keys []string) {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.stringKeys()
}

func (s *Map) stringKeys() (keys []string) {
	s.keys.Range(func(v interface{}) bool {
		keys = append(keys, fmt.Sprintf("%v", v))
		return true
	})
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

// New creates a Map object.
func New() *Map {
	return &Map{
		keys:     glinklist.New(),
		data:     map[interface{}]interface{}{},
		keyElMap: map[interface{}]interface{}{},
	}
}
