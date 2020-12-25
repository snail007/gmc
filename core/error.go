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
	Error() string
	Stack() []byte
	StackStr() string
	String() string
	Callers() []uintptr
	ErrorStack() string
	StackFrames() []StackFrame
	TypeName() string
}

