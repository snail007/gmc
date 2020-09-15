package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Object0 struct {
}

func (s *Object0) Method3() {

}
func (s *Object0) Method4() {

}

type Object struct {
	Object0
}

func (s *Object) Method1(arg string) string {
	return arg
}
func (s *Object) Method2() {

}
func (s *Object) method5() {

}
func TestParseObject(t *testing.T) {
	assert := assert.New(t)
	obj := &Object{}
	var obj1 *Object
	var obj2 = new(Object)
	m := methods(obj)
	assert.Len(m, 4)
	v := invoke(obj, m[0], "called")
	assert.Equal("called", v[0].String())
	v = invoke(obj1, m[0], "called")
	assert.Equal("called", v[0].String())
	v = invoke(obj2, m[0], "called")
	assert.Equal("called", v[0].String())
}
