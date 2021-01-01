package gset

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSet_Add(t *testing.T) {
	assert := assert.New(t)
	s := NewSet()
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
	s := NewSet()
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
	s := NewSet()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	for i := 0; i < 10; i++ {
		assert.Equal(9-i, s.Pop())
	}
}

func TestSet_Contains(t *testing.T) {
	assert := assert.New(t)
	s := NewSet()
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
	s := NewSet()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	s1 := NewSet()
	s1.Merge(s)
	assert.Equal(10, s1.Len())
	for i := 0; i < 10; i++ {
		assert.Equal(i, s1.Shift())
	}
}

func TestSet_MergeSlice(t *testing.T) {
	assert := assert.New(t)
	s := NewSet()
	for i := 0; i < 10; i++ {
		s.Add(i)
	}
	s1 := NewSet()
	s1.MergeSlice(s.ToSlice())
	assert.Equal(10, s1.Len())
	for i := 0; i < 10; i++ {
		assert.Equal(i, s1.Shift())
	}
}

func TestRange_1(t *testing.T) {
	assert := assert.New(t)
	m := NewSet()
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
	m := NewSet()
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
	m := NewSet()
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
