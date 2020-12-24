package gcore

import (
	"fmt"
	"github.com/snail007/gmc/module/error"
	"runtime"
)

type StackFrame interface {
	// The path to the file containing this ProgramCounter
	GetFile() string
	// The LineNumber in that file
	GetLineNumber() int
	// The Name of the function that contains this ProgramCounter
	GetName() string
	// The Package that contains this function
	GetPackage() string
	// The underlying ProgramCounter
	GetProgramCounter() uintptr

	Func() *runtime.Func

	String() string

	SourceLine() (string, error)
}

type Error interface {
	Error() string
	Stack() []byte
	StackStr() string
	String() string
	Callers() []uintptr
	ErrorStack() string
	StackFrames() []StackFrame
	TypeName() string
}

func Recover(f ...interface{}) {
	var f0 interface{}
	var printStack bool
	if len(f) == 0 {
		return
	}
	if len(f) == 2 {
		printStack = f[1].(bool)
	}
	if e := recover(); e != nil {
		f0 = f[0]
		switch v := f0.(type) {
		case func(e interface{}):
			v(e)
		case string:
			s := ""
			if printStack {
				s = fmt.Sprintf(",stack: %s", gerror.Wrap(e).ErrorStack())
			}
			fmt.Printf("\nrecover error, %s%s\n", f, s)
		default:
			fmt.Printf("\nrecover error %s\n", gerror.Wrap(e).ErrorStack())
		}
	}
}
