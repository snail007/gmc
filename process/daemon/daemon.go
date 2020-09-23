// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmcdaemon

import (
	"bufio"
	"fmt"
	logger "log"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"
)

var (
	cmd       *exec.Cmd
	log       = logger.New(os.Stderr, "", logger.LstdFlags)
	isDaemon  = false
	isForever = false
	flog      = ""
)

//SetLogger sets the logger for logging
//
//default is logger.New(os.Stderr, "", logger.LstdFlags)
//
//you can SetLogger(nil) to disable logging
func SetLogger(l *logger.Logger) {
	log = l
}

//start daemon or forever or flog
func Start() (err error) {

	args := []string{}
	if len(os.Args) <= 1 {
		return
	}
	a := os.Args[1:]
	preIsFlog := false
	for i, v := range a {
		vv := strings.TrimSpace(v)
		if vv == "--daemon" || vv == "-daemon" {
			isDaemon = true
		} else if vv == "--forever" || vv == "-forever" {
			isForever = true
		} else if vv == "--flog" || vv == "-flog" {
			if len(a) < i+2 {
				err = fmt.Errorf("logging file path required")
				return
			}
			flog = a[i+1]
			preIsFlog = true
		} else if strings.HasPrefix(vv, "--flog=") || strings.HasPrefix(vv, "-flog=") {
			a := strings.Split(vv, "=")
			if len(a) != 2 {
				err = fmt.Errorf("logging file path required")
				return
			}
			a[1] = strings.Trim(a[1], `"'`)
			if a[1] == "" {
				err = fmt.Errorf("logging file path required")
				return
			}
			flog = a[1]
			preIsFlog = true
		} else {
			if preIsFlog {
				preIsFlog = false
			} else {
				args = append(args, v)
			}
		}
	}
	os.Args = append([]string{os.Args[0]}, args...)

	w := os.Stdout
	if flog != "" {
		f, e := os.OpenFile(flog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if e != nil {
			log.Fatal(e)
		}
		log.SetOutput(f)
		w = f
	}
	if isDaemon {
		var cmd *exec.Cmd
		//fmt.Println("$$" + serviceArgsStr + "$$")
		cmd = exec.Command(os.Args[0], trimArgs("daemon", a)...)
		err = cmd.Start()
		if err != nil {
			err = fmt.Errorf("starting forever fail, error: %s", err)
			return
		}
		if cmd.Process == nil {
			err = fmt.Errorf("starting forever fail, process is nil")
			return
		}
		f := ""
		if isForever {
			f = "forever "
		}
		l("%s%s [PID] %d running...\n", f, os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
	if isForever || flog != "" {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Fprintf(w, "crashed, err: %s\nstack:%s", e, string(debug.Stack()))
				}
			}()
			for {
				if cmd != nil {
					cmd.Process.Kill()
					cmd.Process.Release()
					if !isForever {
						break
					}
					time.Sleep(time.Second * 5)
				}
				cmd = exec.Command(os.Args[0], args...)
				cmdReaderStderr, err := cmd.StderrPipe()
				if err != nil {
					l("ERR:%s,restarting...\n", err)
					continue
				}
				cmdReader, err := cmd.StdoutPipe()
				if err != nil {
					l("ERR:%s,restarting...\n", err)
					continue
				}
				scanner := bufio.NewScanner(cmdReader)
				scannerStdErr := bufio.NewScanner(cmdReaderStderr)
				go func() {
					defer func() {
						if e := recover(); e != nil {
							fmt.Fprintf(w, "crashed, err: %s\nstack:%s", e, string(debug.Stack()))
						}
					}()
					for scanner.Scan() {
						fmt.Fprintf(w, scanner.Text()+"\n")
					}
				}()
				go func() {
					defer func() {
						if e := recover(); e != nil {
							fmt.Fprintf(w, "crashed, err: %s\nstack:%s", e, string(debug.Stack()))
						}
					}()
					for scannerStdErr.Scan() {
						fmt.Fprintf(w, scannerStdErr.Text()+"\n")
					}
				}()
				if err := cmd.Start(); err != nil {
					l("ERR:%s,restarting...\n", err)
					continue
				}
				pid := cmd.Process.Pid
				l("worker %s [PID] %d running...\n", os.Args[0], pid)
				if err := cmd.Wait(); err != nil {
					l("ERR:%s,restarting...", err)
					continue
				}
				if isForever {
					l("worker %s [PID] %d unexpected exited, restarting...\n", os.Args[0], pid)
				}
			}
		}()
	}
	return
}

//clean process, should be call before program exit.
func Clean() {
	if cmd != nil && cmd.ProcessState == nil {
		l("clean process %d", cmd.Process.Pid)
		cmd.Process.Kill()
		cmd.Process.Release()
	}
}

//check if program can be run
func CanRun() bool {
	return !isDaemon && !isForever && flog == ""
}
func trimArgs(trim string, _args []string) []string {
	args := []string{}
	for _, arg := range _args {
		if arg != "-"+trim && arg != "--"+trim {
			args = append(args, arg)
		}
	}
	return args
}

func l(format string, val ...interface{}) {
	if log == nil {
		return
	}
	log.Printf(format, val...)
}
