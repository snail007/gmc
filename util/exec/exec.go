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
	"syscall"
	"time"
)

type Command struct {
	cmdStr        string
	args          []string
	env           map[string]string
	async         bool
	execAsync     bool
	detach        bool
	timeout       time.Duration
	log           gcore.Logger
	outputWriter  io.Writer
	workDir       string
	finalCmd      string
	asyncCallback func(cmd *Command, output string, err error)
	cmd           *exec.Cmd
	beforeExec    func(command *Command, cmd *exec.Cmd)
	afterExec     func(command *Command, cmd *exec.Cmd, err error)
	afterExited   func(command *Command, cmd *exec.Cmd, err error)
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

func (s *Command) BeforeExec(f func(command *Command, cmd *exec.Cmd)) *Command {
	s.beforeExec = f
	return s
}

func (s *Command) AfterExec(f func(command *Command, cmd *exec.Cmd, err error)) *Command {
	s.afterExec = f
	return s
}

func (s *Command) AfterExited(f func(command *Command, cmd *exec.Cmd, err error)) *Command {
	s.afterExited = f
	return s
}

func (s *Command) Args(args ...string) *Command {
	s.args = args
	return s
}

func (s *Command) Cmd() *exec.Cmd {
	return s.cmd
}

// Detach make subprocess detach form current process.
// If current process exited, the child process will not be exited.
//
//	Note: this feature only working in go version >=1.20
func (s *Command) Detach(detach bool) *Command {
	s.detach = detach
	return s
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

// Async create a goroutine to wait the command exited, the command be synchronized executed
// if current process exited, the child process will be exited
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
	var run = func() (err error) {
		defer func() {
			if s.afterExited != nil {
				s.afterExited(s, cmd, err)
			}
		}()
		if s.detach {
			s.cmd.SysProcAttr.Setsid = true
		} else {
			s.cmd.SysProcAttr.Setpgid = true
		}
		if s.beforeExec != nil {
			s.beforeExec(s, cmd)
		}
		err = cmd.Start()
		if s.afterExec != nil {
			s.afterExec(s, cmd, err)
		}
		if err != nil {
			return err
		}
		if !s.execAsync {
			err = cmd.Wait()
			if err != nil {
				return err
			}
		}
		return nil
	}
	if s.outputWriter == nil {
		buf := &bytes.Buffer{}
		cmd.Stdout = buf
		cmd.Stderr = buf
		err := run()
		return buf.Bytes(), err
	} else {
		cmd.Stdout = s.outputWriter
		cmd.Stderr = s.outputWriter
		err := run()
		return nil, err
	}
}

// ExecAsync async execute command on linux system.
func (s *Command) ExecAsync() (e error) {
	if s.async {
		panic("ExecAsync can not run with Async is enabled")
	}
	s.execAsync = true
	_, e = s.Exec()
	return e
}

// Exec execute command on linux system.
func (s *Command) Exec() (output string, e error) {
	sid := fmt.Sprintf("/tmp/tmp_%d", grand.New().Int31()) + ".sh"
	s.finalCmd = `
#!/bin/bash
set -e
cleanup_punaelc() {
     rm -rf ` + sid + `
}
trap cleanup_punaelc EXIT
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
	s.cmd.SysProcAttr = &syscall.SysProcAttr{}
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
