package gmcmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	assert := assert.New(t)
	m := NewMap()
	m.Store("a", "111").
		Store("b", "111")
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
	m := NewMap()
	m.Store("a", "111")
	v, ok := m.Load("a")
	assert.Equal("111", v)
	v, ok = m.LoadOrStore("b", "222")
	assert.False(ok)
	assert.Equal("222", v)
	v, ok = m.LoadOrStore("a", "333")
	assert.True(ok)
	assert.Equal("111", v)
}
func TestRange_1(t *testing.T) {
	assert := assert.New(t)
	m := NewMap()
	m.Store("a", "111").
		Store("b", "111").
		Store("c", "111")
	for i := 0; i < 10; i++ {
		a := []interface{}{}
		m.Range(func(k, v interface{}) bool {
			a = append(a, k)
			return true
		})
		assert.Equal([]interface{}{"a", "b", "c"}, a)
	}
}
func TestRange_2(t *testing.T) {
	assert := assert.New(t)
	m := NewMap()
	m.Store("a", "111").
		Store("b", "111").
		Store("c", "111")
	for i := 0; i < 10; i++ {
		a := []interface{}{}
		m.Range(func(k, v interface{}) bool {
			if k.(string) == "b" {
				return false
			}
			a = append(a, k)
			return true
		})
		assert.Equal([]interface{}{"a"}, a)
	}
}
func TestKeys(t *testing.T) {
	assert := assert.New(t)
	m := NewMap()
	m.Store("a", "111").
		Store("b", "111").
		Store("c", "111")
	// fmt.Println(m.Keys())
	for i := 0; i < 10; i++ {
		assert.Equal([]interface{}{"a", "b", "c"}, m.Keys())
	}
}
func TestKeys_1(t *testing.T) {
	assert := assert.New(t)
	m := NewMap()
	m.Store("a", "111").
		Store("b", "111").
		Store("c", "111").
		Store("c", "111").
		Store("a", "111").
		Store("b", "111")

	assert.Equal(3, m.Len())
	m.Delete("a")
	m.Delete("b")
	m.Delete("c")
	assert.Equal(0, m.Len())

}
func TestKeys_2(t *testing.T) {
	assert := assert.New(t)
	m := NewMap()
	m.Store("a", "111").
		Store("b", "111").
		Store("c", "111").
		Store("d", "111")
	m.Delete("b")
	m.Delete("d")
	assert.Equal([]interface{}{"a", "c"}, m.Keys())

}
