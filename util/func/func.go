package gfunc

import (
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
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
