//go:build !windows
// +build !windows

package gexec

import "syscall"

func setSid(attr *syscall.SysProcAttr) {
	attr.Setsid = true
}
