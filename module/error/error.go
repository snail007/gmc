// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gerror

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"

	gcore "github.com/snail007/gmc/core"
)

// MaxStackDepth the maximum number of stack frames on any error.
var MaxStackDepth = 50

// Error is an error with an attached stacktrace. It can be used
// wherever the builtin error interface is expected.
type Error struct {
	Err    error
	stack  []uintptr
	frames []gcore.StackFrame
	prefix string
}

func New(e ...interface{}) gcore.Error {
	if len(e) > 0 {
		return (&Error{}).New(e)
	}
	return &Error{}
}

func Stack(e interface{}) string {
	return New().Wrap(e).ErrorStack()
}

func Wrap(e interface{}) gcore.Error {
	return New().Wrap(e)
}

// New makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The stacktrace will point to the line of code that
// called New.
func (this *Error) New(e interface{}) gcore.Error {
	var err error

	switch e := e.(type) {
	case error:
		err = e
	default:
		err = fmt.Errorf("%v", e)
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2, stack[:])
	return &Error{
		Err:   err,
		stack: stack[:length],
	}
}

func (this *Error) StackError(e interface{}) string {
	v := fmt.Sprintf("%v", e)
	if e == nil || v == "" {
		return ""
	}
	var err gcore.Error
	if v, ok := e.(*Error); ok {
		err = v
	} else {
		err = this.Wrap(e)
	}
	return err.Error() + "\n" + string(err.Stack())
}

// Wrap makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The skip parameter indicates how far up the stack
// to start the stacktrace. 0 is from the current call, 1 from its caller, etc.
func (this *Error) Wrap(e interface{}) gcore.Error {
	return this.WrapN(e, 2)
}

func (this *Error) WrapN(e interface{}, skip int) gcore.Error {
	if e == nil {
		return nil
	}

	var err error

	switch ev := e.(type) {
	case *Error:
		return ev
	case error:
		err = ev
	default:
		vp := reflect.ValueOf(ev)
		if vp.Kind() == reflect.Ptr {
			err = fmt.Errorf("%s->%v", reflect.TypeOf(ev).String(), vp.Interface())
		} else {
			err = fmt.Errorf("%v", ev)
		}
	}

	stack := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(2+skip, stack[:])
	return &Error{
		Err:   err,
		stack: stack[:length],
	}
}

// WrapPrefix makes an Error from the given value. If that value is already an
// error then it will be used directly, if not, it will be passed to
// fmt.Errorf("%v"). The prefix parameter is used to add a prefix to the
// error message when calling Error(). The skip parameter indicates how far
// up the stack to start the stacktrace. 0 is from the current call,
// 1 from its caller, etc.
func (this *Error) WrapPrefix(e interface{}, prefix string, skip int) gcore.Error {
	return this.WrapPrefixN(e, prefix, skip)
}
func (this *Error) WrapPrefixN(e interface{}, prefix string, skip int) gcore.Error {
	if e == nil {
		return nil
	}

	err := this.WrapN(e, 1+skip).(*Error)

	if err.prefix != "" {
		prefix = fmt.Sprintf("%s: %s", prefix, err.prefix)
	}

	return &Error{
		Err:    err.Err,
		stack:  err.stack,
		prefix: prefix,
	}

}

// Errorf creates a new error with the given message. You can use it
// as a drop-in replacement for fmt.Errorf() to provide descriptive
// errors in return values.
func (this *Error) Errorf(format string, a ...interface{}) gcore.Error {
	return this.WrapN(fmt.Errorf(format, a...), 1)
}

// Error returns the underlying error's message.
func (err *Error) Error() string {

	msg := err.Err.Error()
	if err.prefix != "" {
		msg = fmt.Sprintf("%s: %s", err.prefix, msg)
	}

	return msg
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (err *Error) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range err.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}
func (err *Error) StackStr() string {
	return string(err.Stack())
}

func (err *Error) String() string {
	return string(err.Stack())
}

// Callers satisfies the bugsnag ErrorWithCallerS() interface
// so that the stack can be read out.
func (err *Error) Callers() []uintptr {
	return err.stack
}

// ErrorStack returns a string that contains both the
// error message and the callstack.
func (err *Error) ErrorStack() string {
	return err.TypeName() + " " + err.Error() + "\n" + string(err.Stack())
}

// StackFrames returns an array of frames containing information about the
// stack.
func (err *Error) StackFrames() []gcore.StackFrame {
	if err.frames == nil {
		err.frames = make([]gcore.StackFrame, len(err.stack))

		for i, pc := range err.stack {
			err.frames[i] = NewStackFrame(pc)
		}
	}

	return err.frames
}

// TypeName returns the type this error. e.g. *errors.stringError.
func (err *Error) TypeName() string {
	if _, ok := err.Err.(uncaughtPanic); ok {
		return "panic"
	}
	return reflect.TypeOf(err.Err).String()
}

// Recover usage Recover(func(e interface{})) or Recover(func(e interface{}),bool printStack)
func (this *Error) Recover(f ...interface{}) {
	if e := recover(); e != nil {
		var f0 interface{}
		var printStack bool
		if len(f) == 0 {
			return
		}
		if len(f) == 2 {
			printStack = f[1].(bool)
		}
		f0 = f[0]
		switch v := f0.(type) {
		case func(e interface{}):
			v(e)
		case string:
			s := v
			if printStack {
				s += fmt.Sprintf(",stack: %s", this.StackError(e))
			}
			fmt.Printf("\nrecover error, %v%s\n", f, s)
		default:
			fmt.Printf("\nrecover error %s\n", this.Wrap(e).ErrorStack())
		}
	}
}
