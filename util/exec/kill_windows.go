package gexec

import "os/exec"

func (s *Command) killCmd(c *exec.Cmd) {
	c.Process.Kill()
}
