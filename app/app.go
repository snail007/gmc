package gmcapp

import (
	"fmt"

	"github.com/snail007/gmc/config/app"
	"github.com/snail007/gmc/process/hook"
)

type GMCApp struct {
	beforeRun      []func() error
	beforeShutdown []func()
	cfg            *appconfig.APPConfig
	isBlock        bool
}

func NewConfig() *appconfig.APPConfig {
	return appconfig.NewAPPConfig()
}

func New(cfg *appconfig.APPConfig) GMCApp {
	return GMCApp{
		cfg:     cfg,
		isBlock: true,
	}
}
func (s *GMCApp) Run() (err error) {
	//before run hooks
	hasError := false
	for _, fn := range s.beforeRun {
		func() {
			defer func() {
				if e := recover(); e != nil {
					hasError = true
					err = fmt.Errorf("%s", e)
				}
			}()
			err = fn()
			if err != nil {
				hasError = true
			}
		}()
		if hasError {
			break
		}
	}
	if hasError {
		return
	}
	err = s.run()
	if err != nil {
		return
	}
	hook.RegistShutdown(func() {
		for _, fn := range s.beforeShutdown {
			func() {
				defer func() {
					if e := recover(); e != nil {

					}
				}()
				fn()
			}()
		}
	})
	if s.isBlock {
		hook.WaitShutdown()
	} else {
		go hook.WaitShutdown()
	}
	return
}

func (s *GMCApp) BeforeRun(fn func() (err error)) *GMCApp {
	s.beforeRun = append(s.beforeRun, fn)
	return s
}
func (s *GMCApp) BeforeShutdown(fn func()) *GMCApp {
	s.beforeShutdown = append(s.beforeShutdown, fn)
	return s
}
func (s *GMCApp) Block(isBlockRun bool) *GMCApp {
	s.isBlock = isBlockRun
	return s
}

//init all gmc modules
func (s *GMCApp) run() (err error) {
	err = appconfig.Parse(s.cfg)
	return
}
