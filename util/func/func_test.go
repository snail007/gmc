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
