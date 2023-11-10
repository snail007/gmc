//go:build !windows
// +build !windows

package gexec

import (
	"os/exec"
	"syscall"
)

func (s *Command) killCmd(c *exec.Cmd) {
	if s.detach {
		c.Process.Kill()
	} else {
		syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
	}
}
