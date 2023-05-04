//go:build !windows
// +build !windows

package gexec

import (
	"context"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
	grand "github.com/snail007/gmc/util/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	cmd           string
	env           map[string]string
	async         bool
	timeout       time.Duration
	log           gcore.Logger
	workDir       string
	finalCmd      string
	asyncCallback func(cmd *Command, output string, err error)
}

func NewCommand(cmd string) *Command {
	return &Command{
		cmd:     cmd,
		env:     map[string]string{},
		workDir: "./",
	}
}

func (s *Command) Log(log gcore.Logger) *Command {
	s.log = log
	return s
}

func (s *Command) Timeout(timeout time.Duration) *Command {
	s.timeout = timeout
	return s
}

func (s *Command) Async(async bool) *Command {
	s.async = async
	return s
}

func (s *Command) AsyncCallback(asyncCallback func(cmd *Command, output string, err error)) *Command {
	s.asyncCallback = asyncCallback
	return s
}

func (s *Command) Env(env map[string]string) *Command {
	for k, v := range env {
		s.env[k] = v
	}
	return s
}

func (s *Command) WorkDir(workDir string) *Command {
	s.workDir = workDir
	return s
}

func (s *Command) errorLog(msg string) {
	if s.log == nil {
		return
	}
	s.log.Error(msg)
}

// Exec execute command on linux system.
func (s *Command) Exec() (output string, e error) {
	sid := fmt.Sprintf("/tmp/tmp_%d", grand.New().Int31()) + ".sh"
	defer func() {
		if !s.async {
			os.Remove(sid)
		}
	}()
	s.finalCmd = `
#!/bin/bash
set -e
` + s.cmd
	os.WriteFile(sid, []byte(s.finalCmd), 0755)
	var cmd *exec.Cmd
	var cancel context.CancelFunc
	var ctx context.Context
	if s.timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), s.timeout)
		if !s.async {
			defer cancel()
		}
		cmd = exec.CommandContext(ctx, "bash", sid)
	} else {
		cmd = exec.Command("bash", sid)
	}
	cmd.Dir = s.workDir
	env := map[string]string{}
	for _, v := range os.Environ() {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) != 2 {
			continue
		}
		env[kv[0]] = kv[1]
	}
	if len(s.env) > 0 {
		for k, v := range s.env {
			env[k] = v
		}
	}
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	if !s.async {
		b, err := cmd.CombinedOutput()
		if err != nil {
			e = fmt.Errorf("exec fail, exit code: %d, command: %s, error: %v, output: %s",
				cmd.ProcessState.ExitCode(), s.finalCmd, err, string(b))
			s.errorLog(e.Error())
			return
		}
		output = string(b)
	} else {
		go func() {
			defer func() {
				if cancel != nil {
					cancel()
				}
				os.Remove(sid)
			}()
			defer gerror.Recover(func(e interface{}) {
				s.errorLog(fmt.Sprintf("async exec crashed, command: %s, error: %v", s.finalCmd, e))
			})
			out, err := cmd.CombinedOutput()
			if err != nil && s.log != nil {
				s.errorLog(fmt.Sprintf("async exec fail, exit code: %d, command: %s, error: %v, output: %s",
					cmd.ProcessState.ExitCode(), s.finalCmd, err, string(out)))
			}
			if s.asyncCallback != nil {
				s.asyncCallback(s, string(out), err)
			}
		}()
	}
	return
}
