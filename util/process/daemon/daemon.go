// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdaemon

import (
	"bufio"
	"fmt"
	glog "github.com/snail007/gmc/module/log"
	gcast "github.com/snail007/gmc/util/cast"
	gfile "github.com/snail007/gmc/util/file"
	"io"
	"io/ioutil"
	logger "log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	gcore "github.com/snail007/gmc/core"
)

var (
	cmd            *exec.Cmd
	log            = logger.New(os.Stdout, "", logger.LstdFlags)
	isDaemon       = false
	isForever      = false
	flog           = ""
	initCalled     bool
	initError      error
	initOSArgs     []string
	initArgs       []string
	flogNameExp    string
	rollFileWriter *glog.FileWriter
)

func initFileLogger() {
	flogNameExp = os.Getenv("FLOG_NAME")
	if flogNameExp == "" && gcast.ToBool(os.Getenv("FLOG_BY_MONTH")) {
		flogNameExp = "logs/gateway_%Y%m.log"
	}
	if flogNameExp == "" && gcast.ToBool(os.Getenv("FLOG_BY_DAY")) {
		flogNameExp = "logs/gateway_%Y%m%d.log"
	}
	if flogNameExp == "" && gcast.ToBool(os.Getenv("FLOG_BY_HOUR")) {
		flogNameExp = "logs/gateway_%Y%m%d-%H.log"
	}
	if flogNameExp == "" && gcast.ToBool(os.Getenv("FLOG_BY_MINUTE")) {
		flogNameExp = "logs/gateway_%Y%m%d-%H%i.log"
	}
	if flogNameExp != "" {
		var err error
		rollFileWriter, err = glog.NewFileWriterE(&glog.FileWriterOption{
			Filename:      gfile.BaseName(flogNameExp),
			LogsDir:       filepath.Dir(flogNameExp),
			AliasFilename: "gateway.log",
		})
		if err != nil {
			fmt.Println("init log file error:", err)
			os.Exit(-1)
		}
	}
}

// SetLogger sets the logger for logging
//
// default is logger.New(os.Stderr, "", logger.LstdFlags)
//
// you can SetLogger(nil) to disable logging
func SetLogger(l *logger.Logger) {
	log = l
}

func InitFlags() {
	if initCalled {
		return
	}
	initCalled = true
	if len(os.Args) <= 1 {
		return
	}
	l1 := len(os.Args)
	initOSArgs = make([]string, l1)
	a := make([]string, l1-1)
	copy(initOSArgs, os.Args[:])
	copy(a, os.Args[1:])
	preIsFlog := false
	for i, v := range a {
		vv := strings.TrimSpace(v)
		if vv == "--daemon" || vv == "-daemon" {
			isDaemon = true
		} else if vv == "--forever" || vv == "-forever" {
			isForever = true
		} else if vv == "--flog" || vv == "-flog" {
			if len(a) < i+2 {
				initError = fmt.Errorf("logging file path required")
				return
			}
			flog = a[i+1]
			preIsFlog = true
		} else if strings.HasPrefix(vv, "--flog=") || strings.HasPrefix(vv, "-flog=") {
			a := strings.Split(vv, "=")
			if len(a) != 2 {
				initError = fmt.Errorf("logging file path required")
				return
			}
			a[1] = strings.Trim(a[1], `"'`)
			if a[1] == "" {
				initError = fmt.Errorf("logging file path required")
				return
			}
			flog = a[1]
			preIsFlog = true
		} else {
			if preIsFlog {
				preIsFlog = false
			} else {
				initArgs = append(initArgs, v)
			}
		}
	}
	os.Args = append([]string{os.Args[0]}, initArgs...)
	return
}

// Start daemon or forever or flog
func Start() (err error) {
	if !initCalled {
		InitFlags()
	}
	if len(initOSArgs) <= 1 {
		return
	}
	if initError != nil {
		return initError
	}
	if isDaemon {
		args := trimArgs("daemon", initOSArgs[1:])
		if flog == "" {
			args = append(args, "-flog", "null")
		}
		var cmd *exec.Cmd
		//fmt.Println("$$" + serviceArgsStr + "$$")
		cmd = exec.Command(os.Args[0], args...)
		cmd.Env = os.Environ()
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
		initFileLogger()
		l("%s%s [PID] %d running...\n", f, os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}
	if isForever || flog != "" {
		var w io.Writer
		w = os.Stdout
		if flog == "null" {
			w = ioutil.Discard
		} else if flog != "" {
			f, e := os.OpenFile(flog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			if e != nil {
				log.Fatal(e)
			}
			log.SetOutput(f)
			w = f
		}
		initFileLogger()
		go func() {
			defer gcore.ProviderError()().Recover(func(e interface{}) {
				fprintf(w, "crashed, err: %s", gcore.ProviderError()().StackError(e))
			})
			for {
				if cmd != nil {
					clean(false)
					if !isForever {
						break
					}
					time.Sleep(time.Second * 5)
				}
				cmd = exec.Command(os.Args[0], initArgs...)
				cmd.Env = os.Environ()
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
					defer gcore.ProviderError()().Recover(func(e interface{}) {
						fprintf(w, "crashed, err: %s", gcore.ProviderError()().StackError(e))
					})
					for scanner.Scan() {
						fprintf(w, scanner.Text()+"\n")
					}
				}()
				go func() {
					defer gcore.ProviderError()().Recover(func(e interface{}) {
						fprintf(w, "crashed, err: %s", gcore.ProviderError()().StackError(e))
					})
					for scannerStdErr.Scan() {
						fprintf(w, scannerStdErr.Text()+"\n")
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

// Clean process, should be call before program exit.
func Clean() {
	clean(true)
}

func clean(showlog bool) {
	if cmd != nil && cmd.ProcessState == nil {
		if showlog {
			l("clean process %d", cmd.Process.Pid)
		}
		cmd.Process.Signal(syscall.SIGHUP)
		time.Sleep(time.Second)
		cmd.Process.Kill()
		cmd.Process.Release()
	}
}

// CanRun check if program can be run
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

// control message
func l(format string, val ...interface{}) {
	msg := fmt.Sprintf(format, val...)
	fileLoggerWrite(msg)
	if log == nil {
		return
	}
	log.Print(msg)
}

// child process std & stderr output
func fprintf(w io.Writer, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprint(w, msg)
	fileLoggerWrite(msg)
}

func fileLoggerWrite(msg string) {
	if rollFileWriter==nil{
		return
	}
	rollFileWriter.Write([]byte(msg), gcore.LogLeveNone)
}
