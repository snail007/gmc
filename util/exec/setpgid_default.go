//go:build !windows
// +build !windows

package gexec

import "syscall"

func setPgid(attr *syscall.SysProcAttr) {
	attr.Setsid = true
}
