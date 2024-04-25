package gfunc

import (
	"errors"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckError(t *testing.T) {
	var f1 = func() (err error) {
		defer func() {
			e := recover()
			if e != nil {
				err = fmt.Errorf("%s", e)
			}
		}()
		defer CatchCheckError()
		panic("abc")
	}
	assert.Equal(t, "abc", f1().Error())
	var f2 = func() (err error) {
		defer CatchCheckError()
		err = errors.New("def")
		CheckError(err)
		err = errors.New("abc")
		return
	}
	assert.Equal(t, "def", f2().Error())

	var f3 = func() (err error) {
		defer CatchCheckError()
		CheckError(nil)
		err = errors.New("abc")
		return
	}
	assert.Equal(t, "abc", f3().Error())

	var f4 = func() int {
		defer CatchCheckError()
		v := CheckError2(123, nil)
		return v.Int()
	}
	assert.Equal(t, 123, f4())

	var f5 = func() interface{} {
		defer CatchCheckError()
		v := CheckError2(123, errors.New("abc"))
		return v.Val()
	}
	assert.Equal(t, nil, f5())

	var f6 = func() (int, string) {
		defer CatchCheckError()
		v1, v2 := CheckError3(123, "abc", nil)
		return v1.Int(), v2.String()
	}
	a, b := f6()
	assert.Equal(t, 123, a)
	assert.Equal(t, "abc", b)

	var f7 = func() (interface{}, interface{}) {
		defer CatchCheckError()
		v1, v2 := CheckError3(123, "abc", errors.New("abc"))
		return v1.Val(), v2.Val()
	}
	c, d := f7()
	assert.Equal(t, nil, c)
	assert.Equal(t, nil, d)
}

func TestRecover(t *testing.T) {
	var e gcore.Error
	var f3 = func() {
		defer Recover(func(err gcore.Error) {
			e = err
		})
		panic("abc")
		return
	}
	f3()
	assert.Equal(t, "abc", e.Error())
}
