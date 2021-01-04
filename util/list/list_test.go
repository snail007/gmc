package glist

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestList_Add(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	for i := 0; i < 100; i++ {
		assert.Equal(i, l.Get(i))
	}
}
func TestList_Set(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	for i := 0; i < 101; i++ {
		l.Set(i, i+1)
	}
	for i := 0; i < 100; i++ {
		assert.Equal(i+1, l.Get(i))
	}
}
func TestList_Get(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	for i := 0; i < 101; i++ {
		if i < 100 {
			assert.Equal(i, l.Get(i))
		} else {
			assert.Equal(nil, l.Get(i))
		}
	}
}
func TestList_AddFirst(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i <= 100; i++ {
		l.AddFront(i)
	}
	for i := 0; i <= 100; i++ {
		assert.Equal(100-i, l.Get(i))
	}
}

func TestList_Clear(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i <= 100; i++ {
		l.Add(i)
	}
	l.Clear()
	assert.Equal(0, l.Len())
}

func TestList_Clone(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i <= 100; i++ {
		l.Add(i)
	}
	l1 := l.Clone()
	for i := 0; i < 100; i++ {
		assert.Equal(i, l1.Get(i))
	}
}

func TestList_Contains(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	for i := 0; i < 101; i++ {
		if i >= 100 {
			assert.False(l.Contains(i))
		} else {
			assert.True(l.Contains(i))
		}
	}
}

func TestList_Merge(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l1 := NewList()
	for i := 100; i < 200; i++ {
		l1.Add(i)
	}
	l.Merge(l1)
	for i := 0; i < 200; i++ {
		assert.Equal(i, l.Get(i))
	}
	assert.Equal(200, l.Len())
}

func TestList_MergeSlice(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l1 := []interface{}{}
	for i := 100; i < 200; i++ {
		l1 = append(l1, i)
	}
	l.MergeSlice(l1)
	for i := 0; i < 200; i++ {
		assert.Equal(i, l.Get(i))
	}
	assert.Equal(200, l.Len())
}



func TestList_MergeStringSlice(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(fmt.Sprintf("%v",i))
	}
	l1 := []string{}
	for i := 100; i < 200; i++ {
		l1 = append(l1, fmt.Sprintf("%v",i))
	}
	l.MergeStringSlice(l1)
	for i := 0; i < 200; i++ {
		assert.Equal(fmt.Sprintf("%d",i), l.Get(i))
	}
	assert.Equal(200, l.Len())
}

func TestList_Pop(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	for i := 0; i < 101; i++ {
		if i < 100 {
			assert.Equal(99-i, l.Pop())
		} else {
			assert.Equal(nil, l.Pop())
		}
	}
	assert.Equal(0, l.Len())
}

func TestList_Shift(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	for i := 0; i < 101; i++ {
		if i < 100 {
			assert.Equal(i, l.Shift())
		} else {
			assert.Equal(nil, l.Shift())
		}
	}
	assert.Equal(0, l.Len())
}

func TestList_Sub(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	data := []struct {
		list  *List
		len   int
		isNil bool
	}{
		{l.Sub(0, 0), 0, true},
		{l.Sub(99, 101), 0, true},
		{l.Sub(-1, 2), 0, true},
		{l.Sub(0, 101), 0, true},
		{l.Sub(99, 100), 1, false},
		{l.Sub(0, 10), 10, false},
	}
	for _, v := range data {
		if v.isNil {
			assert.Nil(v.list)
		} else {
			assert.Equal(v.len, v.list.Len())
		}
	}
}

func TestList_Remove(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	l.Add(1)
	l.Add(2)
	l.Remove(0)
	assert.Equal(1, l.Len())
	l.Remove(0)
	assert.Equal(0, l.Len())
	l = NewList()
	l.Add(1)
	l.Add(2)
	l.Remove(1)
	l.Remove(1)
	assert.Equal(1, l.Len())
	l.Remove(0)
	assert.Equal(0, l.Len())
}

func TestList_ToSlice(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 2; i++ {
		l.Add(i)
	}
	assert.Equal([]interface{}{0, 1}, l.ToSlice())
}

func TestList_Range(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	k := 0
	j := 0
	l.Range(func(_ int, v interface{}) bool {
		if v.(int) < 90 {
			k = v.(int)
		} else {
			return false
		}
		j++
		return true
	})
	assert.Equal(90, j)
	assert.Equal(89, k)
}
func TestList_Range2(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l.Range(func(idx int, v interface{}) bool {
		if idx < 90 {
			l.Remove(0)
		} else {
			l.Add(idx)
		}
		return true
	})
	assert.Equal(20, l.Len())
}
func TestList_RangeFast(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	k := 0
	j := 0
	l.RangeFast(func(_ int, v interface{}) bool {
		if v.(int) < 90 {
			k = v.(int)
		} else {
			return false
		}
		j++
		return true
	})
	assert.Equal(90, j)
	assert.Equal(89, k)
}

func TestList_IsEmpty(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 100; i++ {
		l.Add(i)
	}
	l.Clear()
	assert.True(l.IsEmpty())
}

func TestMap_IndexOf(t *testing.T) {
	assert := assert.New(t)
	data := []struct {
		m   *List
		max int
	}{
		{NewList(), 0},
		{NewList(), 1},
		{NewList(), 2},
		{NewList(), 3},
		{NewList(), 100},
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

func TestList_String(t *testing.T) {
	assert := assert.New(t)
	l := NewList()
	for i := 0; i < 2; i++ {
		l.Add(i)
	}
	assert.Equal("[0 1]",fmt.Sprintf("%s",l))
}

func TestList_ToStringSlice(t *testing.T) {
	assert := assert.New(t)
	m := NewList()
	for i := 0; i < 100; i++ {
		m.Add(i)
	}
	m1 := m.ToStringSlice()
	for i := 0; i < 100; i++ {
		assert.Equal(fmt.Sprintf("%d", i), m1[i])
	}
}
