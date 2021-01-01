package gmap

import (
	"container/list"
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
// have more useful function, Len(), Shift(), Pop(), Keys(), etc.
type Map struct {
	keys     *list.List
	data     sync.Map
	lock     sync.Mutex
	keyElMap map[interface{}]*list.Element
}

// NewMap creates a Map object.
func NewMap() *Map {
	return &Map{
		keys:     list.New(),
		data:     sync.Map{},
		lock:     sync.Mutex{},
		keyElMap: map[interface{}]*list.Element{},
	}
}

// Pop returns the last element of map s or nil if the map is empty.
func (s *Map) Pop() (k, v interface{}, ok bool) {
	return s.removeElement(s.keys.Back())
}

// Shift returns the first element of map s or nil if the map is empty.
func (s *Map) Shift() (k, v interface{}, ok bool) {
	return s.removeElement(s.keys.Front())
}

// used for shift and pop
func (s *Map) removeElement(el *list.Element) (k, v interface{}, ok bool) {
	if el == nil {
		return
	}
	v, ok = s.data.Load(el.Value)
	if ok {
		k = el.Value
		s.Delete(el.Value)
	}
	return
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (s *Map) Load(key interface{}) (value interface{}, ok bool) {
	value, ok = s.data.Load(key)
	return
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (s *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data.LoadOrStore(key, value)
	if !loaded {
		s.lock.Lock()
		s.keyElMap[key] = s.keys.PushBack(key)
		s.lock.Unlock()
	}
	return
}

// Store sets the value for a key.
func (s *Map) Store(key, value interface{}) *Map {
	s.data.Store(key, value)
	s.lock.Lock()
	if v, ok := s.keyElMap[key]; ok {
		s.keys.Remove(v)
	}
	s.keyElMap[key] = s.keys.PushBack(key)
	s.lock.Unlock()
	return s
}

// Delete deletes the value for a key.
func (s *Map) Delete(key interface{}) *Map {
	s.data.Delete(key)
	s.lock.Lock()
	if el, ok := s.keyElMap[key]; ok {
		s.keys.Remove(el)
		delete(s.keyElMap, key)
	}
	s.lock.Unlock()
	return s
}

// Len returns the length of the map s.
func (s *Map) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.keys.Len()
}

// Clear deletes all data in the map s.
func (s *Map) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data = sync.Map{}
	s.keys = list.New()
	s.keyElMap = map[interface{}]*list.Element{}
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range keep the sequence of store sequence.
func (s *Map) Range(f func(key, value interface{}) bool) *Map {
	for _, k := range s.Keys() {
		v, _ := s.data.Load(k)
		if !f(k, v) {
			break
		}
	}
	return s
}

// Keys returns all keys in map s and keep the sequence of store sequence.
func (s *Map) Keys() (keys []interface{}) {
	p := s.keys.Front()
	for {
		if p == nil {
			break
		}
		keys = append(keys, p.Value)
		p = p.Next()
	}
	return
}

// MapStringString relative to map[string]string
type MapStringString struct {
	data *Map
}

// NewMapStringString create a map MapStringString
func NewMapStringString() *MapStringString {
	return &MapStringString{
		data: NewMap(),
	}
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (s *MapStringString) Load(key string) (value string, ok bool) {
	v, ok := s.data.Load(key)
	if ok {
		value = v.(string)
	}
	return
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (s *MapStringString) LoadOrStore(key, value string) (actual string, loaded bool) {
	v, loaded := s.data.LoadOrStore(key, value)
	if loaded {
		actual = v.(string)
	}
	return
}

// Store sets the value for a key.
func (s *MapStringString) Store(key, value string) *MapStringString {
	s.data.Store(key, value)
	return s
}

// Delete deletes the value for a key.
func (s *MapStringString) Delete(key string) *MapStringString {
	s.data.Delete(key)
	return s
}

// Len returns the length of the map s.
func (s *MapStringString) Len() int {
	return s.data.Len()
}

// Keys returns all keys in map s and keep the sequence of store sequence.
func (s *MapStringString) Keys() (keys []string) {
	for _, k := range s.data.Keys() {
		keys = append(keys, k.(string))
	}
	return
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range keep the sequence of store sequence.
func (s *MapStringString) Range(f func(key, value string) bool) *MapStringString {
	fn := func(key, value interface{}) bool {
		return f(key.(string), value.(string))
	}
	s.data.Range(fn)
	return s
}

// Pop returns the last element of map s or nil if the map is empty.
func (s *MapStringString) Pop() (k, v string, ok bool) {
	kI, vI, ok := s.data.Pop()
	if ok {
		k = kI.(string)
		v = vI.(string)
	}
	return
}

// Shift returns the first element of map s or nil if the map is empty.
func (s *MapStringString) Shift() (k, v string, ok bool) {
	kI, vI, ok := s.data.Shift()
	if ok {
		k = kI.(string)
		v = vI.(string)
	}
	return
}

// MapStringInterface relative to map[string]interface{}
type MapStringInterface struct {
	data *Map
}

// NewMapStringInterface creates a map MapStringInterface
func NewMapStringInterface() *MapStringInterface {
	return &MapStringInterface{
		data: NewMap(),
	}
}

// Pop returns the last element of map s or nil if the map is empty.
func (s *MapStringInterface) Pop(key string) (k string, v interface{}, ok bool) {
	kI, v, ok := s.data.Pop()
	if ok {
		k = kI.(string)
	}
	return
}

// Shift returns the first element of map s or nil if the map is empty.
func (s *MapStringInterface) Shift(key string) (k string, v interface{}, ok bool) {
	kI, v, ok := s.data.Shift()
	if ok {
		k = kI.(string)
	}
	return
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (s *MapStringInterface) Load(key string) (value interface{}, ok bool) {
	value, ok = s.data.Load(key)
	return
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (s *MapStringInterface) LoadOrStore(key string, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data.LoadOrStore(key, value)
	return
}

// Store sets the value for a key.
func (s *MapStringInterface) Store(key string, value interface{}) *MapStringInterface {
	s.data.Store(key, value)
	return s
}

// Delete deletes the value for a key.
func (s *MapStringInterface) Delete(key string) *MapStringInterface {
	s.data.Delete(key)
	return s
}

// Len returns the length of the map s.
func (s *MapStringInterface) Len() int {
	return s.data.Len()
}

// Keys returns all keys in map s and keep the sequence of store sequence.
func (s *MapStringInterface) Keys() (keys []string) {
	for _, k := range s.data.Keys() {
		keys = append(keys, k.(string))
	}
	return
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range keep the sequence of store sequence.
func (s *MapStringInterface) Range(f func(key string, value interface{}) bool) *MapStringInterface {
	fn := func(key, value interface{}) bool {
		return f(key.(string), value)
	}
	s.data.Range(fn)
	return s
}
