package gexec

import (
	"bytes"
	"context"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gerror "github.com/snail007/gmc/module/error"
	gfile "github.com/snail007/gmc/util/file"
	grand "github.com/snail007/gmc/util/rand"
	"io"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type Command struct {
	cmdStr        string
	args          []string
	env           map[string]string
	async         bool
	timeout       time.Duration
	log           gcore.Logger
	outputWriter  io.Writer
	workDir       string
	finalCmd      string
	asyncCallback func(cmd *Command, output string, err error)
	cmd           *exec.Cmd
	beforeExec    func(command *Command, cmd *exec.Cmd)
}

func NewCommand(cmd string) *Command {
	if runtime.GOOS == "windows" {
		panic("only worked in linux*")
	}
	return &Command{
		cmdStr:  cmd,
		env:     map[string]string{},
		workDir: "./",
	}
}

func (s *Command) BeforeExec(f func(command *Command, cmd *exec.Cmd)) {
	s.beforeExec = f
}

func (s *Command) Args(args ...string) *Command {
	s.args = args
	return s
}

func (s *Command) Cmd() *exec.Cmd {
	return s.cmd
}

func (s *Command) Kill() {
	if s.cmd == nil {
		return
	}
	if s.cmd.ProcessState != nil {
		return
	}
	if s.cmd.Process != nil {
		s.cmd.Process.Kill()
	}
	return
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

func (s *Command) Output(w io.Writer) *Command {
	if (strings.Contains(reflect.TypeOf(w).Kind().String(), "Pointer") ||
		strings.Contains(reflect.TypeOf(w).Kind().String(), "ptr")) &&
		reflect.ValueOf(w).IsNil() {
		w = nil
	}
	s.outputWriter = w
	return s
}

func (s *Command) errorLog(msg string) {
	if s.log == nil {
		return
	}
	s.log.Error(msg)
}

func (s *Command) combinedOutput(cmd *exec.Cmd) ([]byte, error) {
	if s.beforeExec != nil {
		s.beforeExec(s, cmd)
	}
	if s.outputWriter == nil {
		buf := &bytes.Buffer{}
		cmd.Stdout = buf
		cmd.Stderr = buf
		err := cmd.Run()
		return buf.Bytes(), err
	} else {
		cmd.Stdout = s.outputWriter
		cmd.Stderr = s.outputWriter
		err := cmd.Run()
		return nil, err
	}
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
` + s.cmdStr
	gfile.WriteString(sid, s.finalCmd, false)
	var cancel context.CancelFunc
	var ctx context.Context
	if s.timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), s.timeout)
		if !s.async {
			defer cancel()
		}
		s.cmd = exec.CommandContext(ctx, "bash", append([]string{sid}, s.args...)...)
	} else {
		s.cmd = exec.Command("bash", append([]string{sid}, s.args...)...)
	}
	s.cmd.Dir = s.workDir
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
		s.cmd.Env = append(s.cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	if !s.async {
		b, err := s.combinedOutput(s.cmd)
		if err != nil {
			e = fmt.Errorf("exec fail, exit code: %d, command: %s, error: %v, output: %s",
				s.cmd.ProcessState.ExitCode(), s.finalCmd, err, string(b))
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
			out, err := s.combinedOutput(s.cmd)
			if err != nil && s.log != nil {
				s.errorLog(fmt.Sprintf("async exec fail, exit code: %d, command: %s, error: %v, output: %s",
					s.cmd.ProcessState.ExitCode(), s.finalCmd, err, string(out)))
			}
			if s.asyncCallback != nil {
				s.asyncCallback(s, string(out), err)
			}
		}()
	}
	return
}
