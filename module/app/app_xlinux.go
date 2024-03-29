// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

//go:build !windows
// +build !windows

package gapp

import (
	"encoding/json"
	gcore "github.com/snail007/gmc/core"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

func (s *GMCApp) reloadSignalMonitor() {
	go func() {
		// s.logger.Printf("monitor USR2 signal ...")
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGUSR2)
		<-ch
		s.logger.Infof("Received USR2 signal , now reloading ...")
		s.reload()
	}()
}

func (s *GMCApp) reload() {
	files := []*os.File{}
	fdMap := map[int]map[int]bool{}
	k := 0
	for i, srvI := range s.services {
		srv := srvI.Service
		if _, ok := fdMap[i]; !ok {
			fdMap[i] = map[int]bool{}
		}
		for _, l := range srv.Listeners() {
			f, e := l.(*net.TCPListener).File()
			if e != nil {
				s.logger.Warnf("reload fail, %s", e)
				return
			}
			files = append(files, f)
			fdMap[i][k] = true
			k++
		}
	}
	// fmt.Println(fdMap, len(files))
	data, _ := json.Marshal(fdMap)
	cmd := exec.Cmd{}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GMC_REALOD=yes", "GMC_REALOD_DATA="+string(data))
	if len(os.Args) > 1 {
		cmd.Args = os.Args[1:]
	}
	cmd.Path = os.Args[0]
	cmd.ExtraFiles = files
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		s.logger.Warnf("reload fail, fork error : %s", err)
		return
	}

	g := sync.WaitGroup{}
	g.Add(len(s.services))
	for _, srvI := range s.services {
		go func(s gcore.ServiceItem) {
			defer g.Add(-1)
			s.Service.GracefulStop()
		}(srvI)
	}
	g.Wait()
	s.logger.Infof("gmc app reload done.")
	os.Exit(0)
	return
}
