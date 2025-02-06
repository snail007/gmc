package gfunc

import (
	"errors"
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
	gvalue "github.com/snail007/gmc/util/value"
	"time"
)

type panicError struct{}

var (
	panicErr = &panicError{}

	//SafetyCall call func f, and returned recover() error or nil
	SafetyCall = gerror.Try

	//SafetyCallError call func f, and wrap recover() error to gcore.Error, or nil returned
	SafetyCallError = gerror.TryWithStack

	//RecoverNop call recover() and do nothing
	RecoverNop = gerror.RecoverNop

	//RecoverNopAndFunc call recover() and then call the f
	RecoverNopAndFunc = gerror.RecoverNopFunc
)

func Recover(f func(gcore.Error)) {
	e := recover()
	if e != nil {
		f(gerror.Wrap(e))
	}
}

// CatchCheckError catch the panic error throwing by CheckError.
// If func have another defer to catch recover(), put the defer before defer CatchCheckError.
// For example:
//
//	func myFunc() (err error) {
//			defer func() {
//				e := recover()
//				if e != nil {
//					err = fmt.Errorf("%s", e)
//				}
//			}()
//			defer CatchCheckError()
//			panic("abc")
//		}
func CatchCheckError() {
	if err := recover(); err != nil {
		if _, ok := err.(*panicError); ok {
			return
		} else {
			panic(err)
		}
	}
}

// CheckError if err is not nil, the func call will return immediately at CheckError call.
// this should be worked with defer CatchCheckError()
func CheckError(err error) {
	if err != nil {
		panic(panicErr)
	}
	return
}

// CheckError2 check the err and return the value if err is nil
func CheckError2(v interface{}, err error) *gvalue.AnyValue {
	if err != nil {
		panic(panicErr)
	}
	return gvalue.NewAny(v)
}

// CheckError3 check the err and return the value if err is nil
func CheckError3(v1 interface{}, v2 interface{}, err error) (*gvalue.AnyValue, *gvalue.AnyValue) {
	if err != nil {
		panic(panicErr)
	}
	return gvalue.NewAny(v1), gvalue.NewAny(v2)
}

func Wait(f func()) <-chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- SafetyCall(f)
	}()
	return ch
}

var errWaitTimeout = errors.New("wait timeout")

func IsWaitTimeoutErr(e error) bool {
	return errors.Is(e, errWaitTimeout)
}

func WaitTimeout(f func(), timeout time.Duration) error {
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case err := <-Wait(f):
		return err
	case <-t.C:
		return errWaitTimeout
	}
}
