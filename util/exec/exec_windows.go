//go:build windows
// +build windows

package gexec

type Command struct {
}

func NewCommand(cmd string) *Command {
	panic("only worked in linux*")
}
