// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gtest

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

const (
	execFlagPrefix = "GMCT_COVER_EXEC_"
)

func newCmdFromEnv(runName string) (cmd, binName string) {
	if strings.Contains(runName, "/") {
		rs := strings.Split(runName, "/")
		runName = strings.Join(rs, "$/^")
	}
	runName = "^" + runName + "$"
	rb := make([]byte, 16)
	io.ReadFull(rand.Reader, rb)
	d, _ := os.Getwd()
	coverfile := filepath.Join(d, fmt.Sprintf("%x", rb)) + ".gocc.tmp"
	race := os.Getenv("GMCT_COVER_RACE")
	packages := os.Getenv("GMCT_COVER_PACKAGES")
	pkg := strings.TrimPrefix(d, filepath.Join(os.Getenv("GOPATH"), "src"))[1:]
	if race == "true" {
		race = "-race"
	} else {
		race = ""
	}
	if packages != "" {
		// packages is not empty, means run with gmct.
		cover := ""
		rb := make([]byte, 16)
		io.ReadFull(rand.Reader, rb)
		binName = fmt.Sprintf("gmct_testing_%x.bin", rb)
		cover = fmt.Sprintf("-covermode=atomic -coverpkg=%s", packages)
		testCompileCmd := fmt.Sprintf(`go test -c -o %s -run=%s %s %s %s`,
			binName, runName, race, cover, pkg)
		c := exec.Command("bash", "-c", testCompileCmd)
		c.Env = append(c.Env, os.Environ()...)
		out, err := c.CombinedOutput()
		if err != nil {
			fmt.Println(out)
			os.Exit(0)
		}
		binName = filepath.Join(d, binName)
		os.Chmod(binName, 0755)
		cmd = fmt.Sprintf("%s -test.v=true -test.run=%s -test.coverprofile=%s",
			binName, runName, coverfile)
		return
	}
	cmd = fmt.Sprintf(`go test -v -run=%s %s %s`,
		runName, race, pkg)
	return
}

func InGMCT() bool {
	return os.Getenv("GMCT_COVER") == "true"
}

// RunProcess checking if testing is called in NewProcess,
// if true, then call the function f and returns true,
// you should return current testing t function after call CanExec.
func RunProcess(t *testing.T, f func()) bool {
	if os.Getenv(execFlagPrefix+t.Name()) == "true" {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGUSR2)
		go func() {
			f()
			signalChan <- syscall.SIGUSR2
		}()
		<-signalChan
		if msg := os.Getenv("GMCT_COVER_SHOW_KILLED"); msg != "" {
			fmt.Println(msg)
		}
		return true
	}
	return false
}

// NewProcess fork a subprocess runs current testing t function.
func NewProcess(t *testing.T) *Process {
	testFuncName := t.Name()
	isVerbose := os.Getenv("GMCT_COVER_VERBOSE") == "true"
	cmdStr, binName := newCmdFromEnv(testFuncName)
	cleanFile := ""
	if binName!="" {
		cleanFile = strings.SplitN(cmdStr, " ", 2)[0]
	}
	c := exec.Command("bash", "-c", cmdStr)
	c.Env = append(c.Env, os.Environ()...)
	c.Env = append(c.Env, execFlagPrefix+t.Name()+"=true")
	return &Process{
		c:            c,
		testFuncName: testFuncName,
		isVerbose:    isVerbose,
		cmdStr:       cmdStr,
		cleanFile:    cleanFile,
	}
}

type Process struct {
	c            *exec.Cmd
	cleanFile    string
	testFuncName string
	isVerbose    bool
	cmdStr       string
	buf          bytes.Buffer
}

// Wait starts testing subprocess and wait for it exited.
func (s *Process) Wait() (out string, exitCode int, err error) {
	defer func() {
		os.Remove(s.cleanFile)
		if s.isVerbose {
			fmt.Printf("OUTPUT:\n %s", out)
			fmt.Printf(">>> end child testing process %s\n", s.testFuncName)
		}
	}()
	if s.isVerbose {
		fmt.Printf(">>> start child testing process %s\n", s.testFuncName)
		fmt.Println(s.cmdStr)
	}
	b, err := s.c.CombinedOutput()
	out = string(b)
	if s.c.ProcessState != nil {
		exitCode = s.c.ProcessState.ExitCode()
	}
	return
}

// Start starts testing subprocess and return immediately with no wait.
func (s *Process) Start() (err error) {
	if s.isVerbose {
		fmt.Printf(">>> start child testing process %s\n", s.testFuncName)
		fmt.Println(s.cmdStr)
	}
	s.c.Stdout = &s.buf
	s.c.Stderr = &s.buf
	return s.c.Start()
}

// Kill killing testing subprocess.
func (s *Process) Kill() (err error) {
	// s.cleanFile is empty, means not run with gmct.
	if s.cleanFile == "" {
		if s.c.Process != nil {
			return s.c.Process.Kill()
		}
		return
	}
	defer func() {
		os.Remove(s.cleanFile)
		if s.isVerbose {
			fmt.Printf("OUTPUT:\n %s", s.Output())
			fmt.Printf(">>> end child testing process %s\n", s.testFuncName)
		}
	}()
	if s.c.Process != nil {
		err = s.c.Process.Signal(syscall.SIGUSR2)
		if err != nil {
			return err
		}
		_, err = s.c.Process.Wait()
	}
	return nil
}

// IsRunning returns true if testing subprocess is running, otherwise returns false.
func (s *Process) IsRunning() bool {
	return s.c.Process != nil && s.c.ProcessState != nil && !s.c.ProcessState.Exited()
}

// Output acquires the subprocess stdout and stderr output after Start called.
func (s *Process) Output() string {
	return s.buf.String()
}
