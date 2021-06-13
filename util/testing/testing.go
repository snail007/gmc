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
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	execFlagPrefix = "GMCT_COVER_EXEC_"
)

func encodeTestName(n string) string {
	r := strings.NewReplacer("/", "_", ".", "_")
	return r.Replace(n)
}

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
		binName = filepath.Join(os.TempDir(), fmt.Sprintf("gmct_testing_%x.bin", rb))
		cover = fmt.Sprintf("-covermode=atomic -coverpkg=%s", packages)
		testCompileCmd := fmt.Sprintf(`go test -c -o %s -run=%s %s %s %s`,
			binName, runName, race, cover, pkg)
		c := exec.Command("bash", "-c", testCompileCmd)
		c.Env = append(c.Env, os.Environ()...)
		c.Run()
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

func DebugRunProcess(t *testing.T) {
	os.Setenv(execFlagPrefix+encodeTestName(t.Name()), fmt.Sprintf("%v", true))
}

func panicErr(e interface{}) {
	if e != nil {
		panic(e)
	}
}

// RunProcess checking if testing is called in NewProcess,
// if true, then call the function f and returns true,
// you should return current testing t function after call CanExec.
func RunProcess(t *testing.T, f func()) bool {
	if os.Getenv(execFlagPrefix+encodeTestName(t.Name())) == "true" {
		addrFile := os.Getenv("GMCT_COVER_ADDR_FILE")
		defer os.Remove(addrFile)
		exitChan := make(chan bool, 1)
		go func() {
			l, err := net.Listen("tcp", "127.0.0.1:0")
			panicErr(err)
			if err == nil {
				ioutil.WriteFile(addrFile, []byte(l.Addr().String()), 0755)
			}
			_, err = l.Accept()
			panicErr(err)
			exitChan <- true
		}()
		go func() {
			f()
			exitChan <- true
		}()
		<-exitChan
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
	rb := make([]byte, 16)
	io.ReadFull(rand.Reader, rb)
	addrFile := filepath.Join(os.TempDir(), fmt.Sprintf("%x", rb)) + ".addr.tmp"
	os.Setenv("GMCT_COVER_ADDR_FILE", addrFile)
	export := []string{}
	for _, v := range os.Environ() {
		if strings.ContainsAny(v, `./\ &^*()`) {
			continue
		}
		ev := strings.Split(v, "=")
		if len(ev) != 2 {
			continue
		}
		if ev[0] == "" || ev[1] == "" {
			continue
		}
		export = append(export, "export "+v)
	}
	exportStr := strings.Join(export, ";") + ";"
	c := exec.Command("bash", "-c", exportStr+cmdStr)
	c.Env = append(c.Env, os.Environ()...)
	c.Env = append(c.Env, execFlagPrefix+encodeTestName(testFuncName)+"=true")
	return &Process{
		c:            c,
		testFuncName: testFuncName,
		isVerbose:    isVerbose,
		cmdStr:       cmdStr,
		cleanFile:    binName,
		exitChn:      make(chan bool, 1),
		addrFile:     addrFile,
	}
}

type Process struct {
	c            *exec.Cmd
	cleanFile    string
	testFuncName string
	isVerbose    bool
	cmdStr       string
	buf          bytes.Buffer
	exitChn      chan bool
	exited       bool
	started      bool
	startCalled  bool
	addrFile     string
}

// Verbose sets verbose output of testing process.
func (s *Process) Verbose(isVerbose bool) *Process {
	s.isVerbose = isVerbose
	return s
}

// Wait starts testing subprocess and wait for it exited.
func (s *Process) Wait() (out string, exitCode int, err error) {
	exitCode = -1
	if s.started {
		return "", 0, fmt.Errorf("already started")
	}
	s.started = true
	defer func() {
		s.exited = true
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
	if exitCode != 0 {
		output := ""
		if err != nil {
			output = fmt.Sprintf(", output: %s", err.Error())
		}
		err = fmt.Errorf("testing process FAIL, expect exit 0, but got %d"+output, exitCode)
	}
	return
}

// Start starts testing subprocess and return immediately with no wait.
func (s *Process) Start() (err error) {
	if s.started {
		return fmt.Errorf("already started")
	}
	s.started = true
	s.startCalled = true
	if s.isVerbose {
		fmt.Printf(">>> start child testing process %s\n", s.testFuncName)
		fmt.Println(s.cmdStr)
	}
	s.c.Stdout = &s.buf
	s.c.Stderr = &s.buf
	err = s.c.Start()
	if err == nil {
		go func() {
			s.c.Process.Wait()
			s.exitChn <- true
			s.exited = true
			os.Remove(s.cleanFile)
		}()
	}
	return err
}

// Kill killing testing subprocess.
func (s *Process) Kill() (err error) {
	if !s.startCalled {
		return fmt.Errorf("not run")
	}
	defer func() {
		if s.isVerbose {
			fmt.Printf("OUTPUT:\n %s", s.Output())
			fmt.Printf(">>> end child testing process %s\n", s.testFuncName)
		}
	}()
	i := 0
	for i < 10 {
		time.Sleep(time.Second)
		i++
		contents, _ := ioutil.ReadFile(s.addrFile)
		if len(contents) == 0 {
			continue
		}
		os.Remove(s.addrFile)
		_, err = net.Dial("tcp", string(contents))
		if err == nil {
			break
		}
	}
	select {
	case <-s.exitChn:
	case <-time.After(time.Second * 3):
		fmt.Println("[WARN] kill timeout, force kill pid", s.c.Process.Pid)
		if s.c.Process != nil {
			s.c.Process.Kill()
		}
	}
	return nil
}

// IsRunning returns true if testing subprocess is running, otherwise returns false.
func (s *Process) IsRunning() bool {
	return !s.exited
}

// Output acquires the subprocess stdout and stderr output after Start called.
func (s *Process) Output() string {
	var b bytes.Buffer
	for _, line := range strings.Split(s.buf.String(), "\n") {
		if strings.Contains(line, "warning: no packages being tested depend on matches for pattern") {
			continue
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}
