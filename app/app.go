package gmcapp

import (
	"fmt"

	appconfig "github.com/snail007/gmc/config/app"
	"github.com/snail007/gmc/process/hook"
	"github.com/spf13/viper"
)

type GMCApp struct {
	beforeRun      []func() error
	beforeShutdown []func()
	cfg            *appconfig.APPConfig
	isBlock        bool
	isParsed       bool
	config         *viper.Viper
	configfile     string
}

func New() GMCApp {
	return GMCApp{
		isBlock: true,
	}
}
func (s *GMCApp) SetConfigFile(file string) *GMCApp {
	s.configfile = file
	return s
}
func (s *GMCApp) ParseConfig() (err error) {
	if s.isParsed {
		return
	}
	defer func() {
		if err == nil {
			s.isParsed = true
		}
	}()
	//create config file object
	s.config = viper.New()
	if s.configfile == "" {
		s.config.SetConfigName("app")
		s.config.AddConfigPath(".")
		s.config.AddConfigPath("config")
		s.config.AddConfigPath("conf")
		//for testing
		// s.config.AddConfigPath("../app")
	} else {
		s.config.SetConfigFile(s.configfile)
	}
	err = s.config.ReadInConfig()
	if err != nil {
		return
	}
	return
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
