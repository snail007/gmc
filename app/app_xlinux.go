// +build !windows

package gmcapp

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func (s *GMCApp) reloadSignalMonitor() {
	go func() {
		// s.logger.Printf("monitor USR2 signal ...")
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGUSR2)
		<-ch
		s.logger.Printf("Recieved USR2 signal , now reloading ...")
		s.reload()
	}()
}

func (s *GMCApp) reload() {
	files := []*os.File{}
	skip := []string{"-1"}
	for i, srvI := range s.services {
		srv := srvI.Service
		l := srv.Listener()
		if l != nil {
			f, e := l.(*net.TCPListener).File()
			if e != nil {
				s.logger.Printf("reload fail, %s", e)
				return
			}
			files = append(files, f)
		} else {
			skip = append(skip, fmt.Sprintf("%d", i))
			files = append(files, nil)
		}
	}

	cmd := exec.Cmd{}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GMC_REALOD=yes", "GMC_REALOD_SKIP="+strings.Join(skip, ","))
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
		s.logger.Printf("reload fail, fork error : %s", err)
		return
	}

	g := sync.WaitGroup{}
	g.Add(len(s.services))
	for _, srvI := range s.services {
		go func(s ServiceItem) {
			defer g.Add(-1)
			s.Service.GracefulStop()
		}(srvI)
	}
	g.Wait()
	s.logger.Printf("gmc app reload done.")
	os.Exit(0)
	return
}
