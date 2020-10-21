/*
Package maputil is a order map implements
*/
package _map

import (
	"container/list"
	"sync"
)

type Map struct {
	keys     *list.List
	data     sync.Map
	lock     sync.Mutex
	keyElMap map[interface{}]*list.Element
}

func NewMap() *Map {
	return &Map{
		keys:     list.New(),
		data:     sync.Map{},
		lock:     sync.Mutex{},
		keyElMap: map[interface{}]*list.Element{},
	}
}
func (s *Map) Load(key interface{}) (value interface{}, ok bool) {
	value, ok = s.data.Load(key)
	return
}

func (s *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data.LoadOrStore(key, value)
	if !loaded {
		s.lock.Lock()
		s.keyElMap[key] = s.keys.PushBack(key)
		s.lock.Unlock()
	}
	return
}

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
func (s *Map) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.keys.Len()
}
func (s *Map) Range(f func(key, value interface{}) bool) *Map {
	for _, k := range s.Keys() {
		v, _ := s.data.Load(k)
		if !f(k, v) {
			break
		}
	}
	return s
}
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

//MapStringString relative to map[string]string
type MapStringString struct {
	data *Map
}

func NewMapStringString() *MapStringString {
	return &MapStringString{
		data: NewMap(),
	}
}
func (s *MapStringString) Load(key string) (value string, ok bool) {
	v, ok := s.data.Load(key)
	if ok {
		value = v.(string)
	}
	return
}
func (s *MapStringString) LoadOrStore(key, value string) (actual string, loaded bool) {
	v, loaded := s.data.LoadOrStore(key, value)
	if loaded {
		actual = v.(string)
	}
	return
}
func (s *MapStringString) Store(key, value string) *MapStringString {
	s.data.Store(key, value)
	return s
}

func (s *MapStringString) Delete(key string) *MapStringString {
	s.data.Delete(key)
	return s
}

func (s *MapStringString) Len() int {
	return s.data.Len()
}
func (s *MapStringString) Keys() (keys []string) {
	for _, k := range s.data.Keys() {
		keys = append(keys, k.(string))
	}
	return
}
func (s *MapStringString) Range(f func(key, value string) bool) *MapStringString {
	fn := func(key, value interface{}) bool {
		return f(key.(string), value.(string))
	}
	s.data.Range(fn)
	return s
}

//MapStringInterface relative to map[string]interface{}
type MapStringInterface struct {
	data *Map
}

func NewMapStringInterface() *MapStringInterface {
	return &MapStringInterface{
		data: NewMap(),
	}
}
func (s *MapStringInterface) Load(key string) (value interface{}, ok bool) {
	value, ok = s.data.Load(key)
	return
}
func (s *MapStringInterface) LoadOrStore(key string, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = s.data.LoadOrStore(key, value)
	return
}
func (s *MapStringInterface) Store(key string, value interface{}) *MapStringInterface {
	s.data.Store(key, value)
	return s
}

func (s *MapStringInterface) Delete(key string) *MapStringInterface {
	s.data.Delete(key)
	return s
}

func (s *MapStringInterface) Len() int {
	return s.data.Len()
}
func (s *MapStringInterface) Keys() (keys []string) {
	for _, k := range s.data.Keys() {
		keys = append(keys, k.(string))
	}
	return
}
func (s *MapStringInterface) Range(f func(key string, value interface{}) bool) *MapStringInterface {
	fn := func(key, value interface{}) bool {
		return f(key.(string), value)
	}
	s.data.Range(fn)
	return s
}
