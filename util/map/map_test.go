// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmap

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", "111")
	assert.Equal(2, m.Len())
	m.Delete("a")
	assert.Equal(1, m.Len())
	m.Range(func(k, v interface{}) bool {
		m.Delete(k)
		return true
	})
	assert.Equal(0, m.Len())
}

func TestLoad(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	v, _ := m.Load("a")
	assert.Equal("111", v)
	v, ok := m.LoadOrStore("b", "222")
	assert.False(ok)
	assert.Equal("222", v)
	v, ok = m.LoadOrStore("a", "333")
	assert.True(ok)
	assert.Equal("111", v)
}
func TestRange_1(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", "111")
	m.Store("c", "111")
	a := []interface{}{}
	m.Range(func(k, v interface{}) bool {
		a = append(a, k)
		return true
	})
	assert.Equal([]interface{}{"a", "b", "c"}, a)
}
func TestRange_2(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", "111")
	m.Store("c", "111")
	a := []interface{}{}
	m.Range(func(k, v interface{}) bool {
		if k.(string) == "b" {
			m.Delete(k.(string))
			m.Store("d", "111")
			return false
		}
		a = append(a, k)
		return true
	})
	assert.Equal([]interface{}{"a"}, a)
	assert.Equal([]interface{}{"a", "c", "d"}, m.Keys())
}
func TestMap_RangeFast(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", "111")
	m.Store("c", "111")
	a := []interface{}{}
	m.RangeFast(func(k, v interface{}) bool {
		if k.(string) == "c" {
			return false
		}
		a = append(a, k)
		return true
	})
	assert.Equal([]interface{}{"a", "b"}, a)
}
func TestKeys(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", "111")
	m.Store("c", "111")
	// fmt.Println(m.Keys())
	for i := 0; i < 10; i++ {
		assert.Equal([]interface{}{"a", "b", "c"}, m.Keys())
	}
}
func TestKeys_1(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", "111")
	m.Store("c", "111")
	m.Store("c", "111")
	m.Store("a", "111")
	m.Store("b", "111")

	assert.Equal(3, m.Len())
	m.Delete("a")
	m.Delete("b")
	m.Delete("c")
	assert.Equal(0, m.Len())

}
func TestKeys_2(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", " 11")
	m.Store("c", "111")
	m.Store("d", "111")
	m.Delete("b")
	m.Delete("d")
	assert.Equal([]interface{}{"a", "c"}, m.Keys())

}

func TestShift(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", "222")
	m.Store("c", "333")
	m.Store("d", "444")
	data := []struct {
		key    interface{}
		except interface{}
		len    int
		ok     bool
	}{
		{"a", "111", 3, true},
		{"b", "222", 2, true},
		{"c", "333", 1, true},
		{"d", "444", 0, true},
		{nil, nil, 0, false},
	}
	for _, v := range data {
		key, value, ok := m.Shift()
		assert.Equal(v.ok, ok)
		assert.Equal(v.except, value)
		assert.Equal(v.len, m.Len())
		assert.Equal(v.key, key)
	}
}

func TestPop(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Store("a", "111")
	m.Store("b", "222")
	m.Store("c", "333")
	m.Store("d", "444")
	data := []struct {
		key    interface{}
		except interface{}
		len    int
		ok     bool
	}{
		{nil, nil, 0, false},
		{"a", "111", 0, true},
		{"b", "222", 1, true},
		{"c", "333", 2, true},
		{"d", "444", 3, true},
	}
	for i := len(data) - 1; i >= 0; i-- {
		v := data[i]
		key, value, ok := m.Pop()
		assert.Equal(v.ok, ok)
		assert.Equal(v.except, value)
		assert.Equal(v.len, m.Len())
		assert.Equal(v.key, key)
	}
}

func TestMap_Clone(t *testing.T) {
	assert := assert.New(t)
	m := New()
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	m1 := m.Clone()
	for i := 0; i < 101; i++ {
		v, ok := m1.Load(i)
		if i < 100 {
			assert.True(ok)
			assert.Equal(i, v)
		} else {
			assert.False(false)
		}
	}
}

func TestMap_ToMap(t *testing.T) {
	assert := assert.New(t)
	m := New()
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	m1 := m.ToMap()
	for i := 0; i < 101; i++ {
		v, ok := m1[i]
		if i < 100 {
			assert.True(ok)
			assert.Equal(i, v)
		} else {
			assert.False(false)
		}
	}
}

func TestMap_ToStringMap(t *testing.T) {
	assert := assert.New(t)
	m := New()
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	m1 := m.ToStringMap()
	for i := 0; i < 101; i++ {
		v, ok := m1[fmt.Sprintf("%d", i)]
		if i < 100 {
			assert.True(ok)
			assert.Equal(i, v)
		} else {
			assert.False(false)
		}
	}
}

func TestMap_Merge(t *testing.T) {
	assert := assert.New(t)
	m := New()
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	assert.Equal(100, m.Len())
	m1 := New()
	m1.Merge(m)
	assert.Equal(100, m1.Len())
}

func TestMap_MergeMap(t *testing.T) {
	assert := assert.New(t)
	m := map[interface{}]interface{}{}
	for i := 0; i < 100; i++ {
		m[i] = i
	}
	assert.Equal(100, len(m))
	m1 := New()
	m1.MergeMap(m)
	assert.Equal(100, m1.Len())
}

func TestMap_MergeStrMap(t *testing.T) {
	assert := assert.New(t)
	m := map[string]interface{}{}
	for i := 0; i < 100; i++ {
		m[fmt.Sprintf("%v", i)] = i
	}
	assert.Equal(100, len(m))
	m1 := New()
	m1.MergeStrMap(m)
	assert.Equal(100, m1.Len())
}

func TestMap_MergeStrStrMap(t *testing.T) {
	assert := assert.New(t)
	m := map[string]string{}
	for i := 0; i < 100; i++ {
		m[fmt.Sprintf("%v", i)] = fmt.Sprintf("%v", i)
	}
	assert.Equal(100, len(m))
	m1 := New()
	m1.MergeStrStrMap(m)
	assert.Equal(100, m1.Len())
}

func TestMap_MergeSyncMap(t *testing.T) {
	assert := assert.New(t)
	m := &sync.Map{}
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	m1 := New()
	m1.MergeSyncMap(m)
	assert.Equal(100, m1.Len())
}

func TestMap_Clear(t *testing.T) {
	assert := assert.New(t)
	m := New()
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	assert.Equal(100, m.Len())
	m.Clear()
	assert.Equal(0, m.Len())
}

func TestMap_IsEmpty(t *testing.T) {
	assert := assert.New(t)
	m := New()
	for i := 0; i < 100; i++ {
		m.Store(i, i)
	}
	assert.Equal(100, m.Len())
	assert.False(m.IsEmpty())
	m.Clear()
	assert.True(m.IsEmpty())
}

func TestMap_IndexOf(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		m   *Map
		max int
	}{
		{New(), 0},
		{New(), 1},
		{New(), 2},
		{New(), 3},
		{New(), 100},
	}
	for _, v := range data {
		for i := 0; i < v.max; i++ {
			v.m.Store(i, i)
		}
		for i := 0; i < v.max+1; i++ {
			if i < v.max {
				assert.Equal(i, v.m.IndexOf(i))
			} else {
				assert.Equal(-1, v.m.IndexOf(i))
			}
		}
	}
}

func TestMap_String(t *testing.T) {
	assert := assert.New(t)
	l := New()
	for i := 0; i < 1; i++ {
		l.Store(i, i)
	}
	assert.Equal("map[0:0]", l.String())
}

func TestMap_StoreFront(t *testing.T) {
	assert := assert.New(t)
	l := New()
	for i := 0; i < 3; i++ {
		l.StoreFront(i, i)
	}
	assert.Equal("[2 1 0]", fmt.Sprintf("%v", l.StringKeys()))
}

func TestMap_StoreFront2(t *testing.T) {
	assert := assert.New(t)
	l := New()
	for i := 0; i < 3; i++ {
		l.StoreFront(i, i)
	}
	for i := 0; i < 3; i++ {
		l.StoreFront(i, i)
	}
	assert.Equal("[2 1 0]", fmt.Sprintf("%v", l.StringKeys()))
}

func TestMap_LoadOrStoreFront(t *testing.T) {
	assert := assert.New(t)
	l := New()
	for i := 0; i < 3; i++ {
		l.LoadOrStoreFront(i, i)
	}
	for i := 2; i >= 0; i-- {
		l.LoadOrStoreFront(i, i)
	}
	assert.Equal("[2 1 0]", fmt.Sprintf("%v", l.StringKeys()))
}

func TestMap_LoadOrStoreFunc(t *testing.T) {
	assert := assert.New(t)
	l := New()
	for i := 0; i < 3; i++ {
		j := i
		l.LoadOrStoreFunc(i, func() interface{} {
			return j
		})
	}
	k := 0
	for i := 0; i < 3; i++ {
		j := i
		l.LoadOrStoreFunc(i, func() interface{} {
			k++
			return j
		})
	}
	assert.Equal(0, k)
	assert.Equal("[0 1 2]", fmt.Sprintf("%v", l.StringKeys()))
}

func TestMap_LoadAndStoreFunc(t *testing.T) {
	assert := assert.New(t)
	l := New()
	for i := 0; i < 3; i++ {
		l.LoadAndStoreFunc(i, func(_ interface{}, loaded bool) (newValue interface{}) {
			assert.False(loaded)
			return i + 3
		})
	}
	k := 0
	for i := 0; i < 3; i++ {
		l.LoadAndStoreFunc(i, func(oldValue interface{}, loaded bool) (newValue interface{}) {
			k++
			assert.Equal(i+3, oldValue)
			return i
		})
	}
	assert.Equal(3, k)
	assert.Equal("[0 1 2]", fmt.Sprintf("%v", l.StringKeys()))
}
