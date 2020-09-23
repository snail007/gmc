// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmcprofiling

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"strings"
	"sync"
)

var (
	profilingLock = &sync.Mutex{}
	isProfiling   = false
	cpuProfilingFile, memProfilingFile, blockProfilingFile,
	goroutineProfilingFile, threadcreateProfilingFile *os.File
	found   bool
	bindArg = ""
	argv    = ""
	inited  = false
)

func _init() {
	os.Args, argv, found = trimArgs(bindArg, os.Args)
	//fmt.Println(os.Args, argv, found)
}

//bind command line argument, no - or -- prefix, such as: profiling
//
//./foobar -profiling=debug
//
//all profiling files will be stored in named debug folder, folder must be exists
func StartArg(argName string) {
	start(argName)
}

//Start's arguments is store floder path,
//
//all profiling files will be stored in "storePath" folder, folder must be exists
func Start(storePath_ string) {
	start("", storePath_)
}
func start(cliName string, storePath_ ...string) {
	storePath := ""
	if cliName != "" {
		if !inited {
			bindArg = cliName
			_init()
		}
		storePath = argv
	} else if len(storePath_) > 0 {
		storePath = storePath_[0]
	} else {
		return
	}
	profilingLock.Lock()
	defer profilingLock.Unlock()
	if !isProfiling {
		isProfiling = true
		if storePath == "" {
			storePath = "."
		}
		cpuProfilingFile, _ = os.Create(filepath.Join(storePath, "cpu.prof"))
		memProfilingFile, _ = os.Create(filepath.Join(storePath, "memory.prof"))
		blockProfilingFile, _ = os.Create(filepath.Join(storePath, "block.prof"))
		goroutineProfilingFile, _ = os.Create(filepath.Join(storePath, "goroutine.prof"))
		threadcreateProfilingFile, _ = os.Create(filepath.Join(storePath, "threadcreate.prof"))
		pprof.StartCPUProfile(cpuProfilingFile)
	}
}

//Stop profiling, and write profiling data to file, must be call when you needed, if it not be called, nothing profiling data will be write.
func Stop() {
	profilingLock.Lock()
	defer profilingLock.Unlock()
	if isProfiling {
		isProfiling = false
		pprof.StopCPUProfile()
		goroutine := pprof.Lookup("goroutine")
		goroutine.WriteTo(goroutineProfilingFile, 0)
		heap := pprof.Lookup("heap")
		heap.WriteTo(memProfilingFile, 0)
		block := pprof.Lookup("block")
		block.WriteTo(blockProfilingFile, 0)
		threadcreate := pprof.Lookup("threadcreate")
		threadcreate.WriteTo(threadcreateProfilingFile, 0)
		//close
		goroutineProfilingFile.Close()
		memProfilingFile.Close()
		blockProfilingFile.Close()
		threadcreateProfilingFile.Close()
	}
}
func trimArgs(trim string, _args []string) ([]string, string, bool) {
	if trim == "" {
		return nil, "", false
	}
	args := []string{}
	found := false
	value := ""
	for i, arg := range _args {
		if match, _ := regexp.MatchString(`--?`+bindArg+`(=.+)?`, arg); match {
			found = true
			if strings.Contains(arg, "=") {
				a := strings.Split(arg, "=")
				value = strings.Trim(a[1], `"'`)
			} else if len(_args) >= i+2 {
				value = _args[i+1]
				i++
			} else {
				found = false
			}
		} else {
			args = append(args, arg)
		}
	}
	return args, value, found
}
