// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gset

import (
	"fmt"
	assert "github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSet_Add(t *testing.T) {
	assert := assert.New(t)
	s := New()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	assert.Equal(10, s.Len())
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	assert.Equal(10, s.Len())
}

func TestSet_Clone(t *testing.T) {
	assert := assert.New(t)
	s := New()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	s1 := s.Clone()
	assert.Equal(10, s1.Len())
	for i := 0; i < 10; i++ {
		assert.Equal(i, s1.Shift())
	}
}

func TestSet_Pop(t *testing.T) {
	assert := assert.New(t)
	s := New()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	for i := 0; i < 10; i++ {
		assert.Equal(9-i, s.Pop())
	}
}

func TestSet_Contains(t *testing.T) {
	assert := assert.New(t)
	s := New()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	for i := 0; i < 11; i++ {
		if i < 10 {
			assert.True(s.Contains(i))
		} else {
			assert.False(s.Contains(i))
		}
	}
}

func TestSet_Merge(t *testing.T) {
	assert := assert.New(t)
	s := New()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	s1 := New()
	s1.Merge(s)
	assert.Equal(10, s1.Len())
	for i := 0; i < 10; i++ {
		assert.Equal(i, s1.Shift())
	}
}

func TestSet_MergeSlice(t *testing.T) {
	assert := assert.New(t)
	s := New()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	s1 := New()
	s1.MergeSlice(s.ToSlice())
	assert.Equal(10, s1.Len())
	for i := 0; i < 10; i++ {
		assert.Equal(i, s1.Shift())
	}
}

func TestSet_MergeStringSlice(t *testing.T) {
	assert := assert.New(t)
	s := New()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	s1 := New()
	s1.MergeStringSlice(s.ToStringSlice())
	assert.Equal(10, s1.Len())
	for i := 0; i < 10; i++ {
		assert.Equal(fmt.Sprintf("%d", i), fmt.Sprintf("%v", s1.Shift()))
	}
}

func TestRange_1(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Add("a")
	m.Add("b")
	m.Add("c")
	a := []interface{}{}
	m.Range(func(k interface{}) bool {
		a = append(a, k)
		return true
	})
	assert.Equal([]interface{}{"a", "b", "c"}, a)
}
func TestRange_2(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Add("a")
	m.Add("b")
	m.Add("c")
	a := []interface{}{}
	m.Range(func(k interface{}) bool {
		if k.(string) == "b" {
			m.Delete(k.(string))
			m.Add("d")
			return false
		}
		a = append(a, k)
		return true
	})
	assert.Equal([]interface{}{"a"}, a)
	assert.Equal([]interface{}{"a", "c", "d"}, m.ToSlice())
}
func TestMap_RangeFast(t *testing.T) {
	assert := assert.New(t)
	m := New()
	m.Add("a")
	m.Add("b")
	m.Add("c")
	a := []interface{}{}
	m.RangeFast(func(k interface{}) bool {
		if k.(string) == "c" {
			return false
		}
		a = append(a, k)
		return true
	})
	assert.Equal([]interface{}{"a", "b"}, a)
}

func TestMap_IndexOf(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		m   *Set
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
			v.m.Add(i)
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

func TestSet_String(t *testing.T) {
	assert := assert.New(t)
	l := New()
	for i := 0; i < 2; i++ {
		l.Add(i)
	}
	assert.Equal("[0 1]", fmt.Sprintf("%s", l))
}

func TestSet_ToStringSlice(t *testing.T) {
	assert := assert.New(t)
	m := New()
	for i := 0; i < 100; i++ {
		m.Add(i)
	}
	m1 := m.ToStringSlice()
	for i := 0; i < 100; i++ {
		assert.Equal(fmt.Sprintf("%d", i), m1[i])
	}
}

func TestSet_AddFront(t *testing.T) {
	assert := assert.New(t)
	l := New()
	for i := 0; i < 2; i++ {
		l.AddFront(i)
	}
	for i := 1; i >= 0; i-- {
		l.AddFront(i)
	}
	assert.Equal("[1 0]", fmt.Sprintf("%s", l))
}

func TestSet_CartesianProduct(t *testing.T) {
	s := New()
	s.Add(1, 2)
	s1 := New()
	s1.Add(3, 4)
	want := [][]interface{}{{1, 3}, {1, 4}, {2, 3}, {2, 4}}
	if got := s.CartesianProduct(s1); !reflect.DeepEqual(got, want) {
		t.Errorf("CartesianProduct() = %v, want %v", got, want)
	}
}
