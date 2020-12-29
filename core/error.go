package gcore

import (
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
	New(e interface{}) Error
	StackError(e interface{}) string
	Wrap(e interface{}) Error
	WrapN(e interface{}, skip int) Error
	WrapPrefix(e interface{}, prefix string, skip int) Error
	WrapPrefixN(e interface{}, prefix string, skip int) Error
	Errorf(format string, a ...interface{}) Error
	Error() string
	Stack() []byte
	StackStr() string
	String() string
	Callers() []uintptr
	ErrorStack() string
	StackFrames() []StackFrame
	TypeName() string
	Recover(f ...interface{})
}
