package gapp

import (
	"encoding/json"
	"fmt"
	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/module/cache"
	gdb "github.com/snail007/gmc/module/db"
	gerr "github.com/snail007/gmc/module/error"
	gi18n "github.com/snail007/gmc/module/i18n"
	"github.com/snail007/gmc/module/log"
	gconfig "github.com/snail007/gmc/util/config"
	ghook "github.com/snail007/gmc/util/process/hook"
	"net"
	"os"
	"strings"
)

type GMCApp struct {
	onRun             []func(*gconfig.Config) error
	onShutdown        []func()
	isBlock           bool
	attachConfig      map[string]*gconfig.Config
	attachConfigfiles map[string]string
	services          []gcore.ServiceItem
	logger            gcore.Logger
	configFile        string
	config            *gconfig.Config
}

func New() gcore.GMCApp {
	return &GMCApp{
		isBlock:           true,
		onRun:             []func(*gconfig.Config) error{},
		onShutdown:        []func(){},
		services:          []gcore.ServiceItem{},
		logger:            nil,
		attachConfig:      map[string]*gconfig.Config{},
		attachConfigfiles: map[string]string{},
	}
}

func Default() gcore.GMCApp {
	// default config dir and name
	cfg := gconfig.NewConfig()
	cfg.SetConfigName("app")
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("conf")
	cfg.AddConfigPath("config")
	a := New()
	a.SetConfig(cfg)
	return a
}
func (s *GMCApp) initialize() (err error) {
	defer func() {
		if s.logger == nil {
			s.logger = glog.NewGMCLog()
		}
		if s.logger.Async() {
			s.OnShutdown(func() {
				s.logger.WaitAsyncDone()
			})
		}
	}()

	if s.config == nil {
		return
	}

	// initialize logging
	if s.config.Sub("log") != nil && s.logger == nil {
		s.logger = glog.NewFromConfig(s.config)
	}

	// initialize database
	if s.config.Sub("database") != nil {
		err = gdb.Init(s.config)
		if err != nil {
			return
		}
	}

	// initialize cache
	if s.config.Sub("cache") != nil {
		err = gcache.Init(s.config)
		if err != nil {
			return
		}
	}
	// initialize i18n if needed
	if s.config.Sub("i18n") != nil {
		err = gi18n.Init(s.config)
	}
	return
}

func (s *GMCApp) SetConfigFile(file string) {
	s.configFile = file
}
func (s *GMCApp) SetConfig(cfg *gconfig.Config) {
	s.config = cfg
}
func (s *GMCApp) AttachConfigFile(id, file string) {
	s.attachConfigfiles[id] = file
}
func (s *GMCApp) AttachConfig(id string, cfg *gconfig.Config) {
	s.attachConfig[id] = cfg
}

//Config acquires the  or attach config object.
//if `idanem` is empty , it return   config object,
//other return attach config object of `id`.
func (s *GMCApp) Config(id ...string) *gconfig.Config {
	if len(id) > 0 {
		return s.attachConfig[id[0]]
	}
	return s.config
}
func (s *GMCApp) parseConfigFile() (err error) {
	parse := false
	if s.configFile != "" {
		if s.config == nil {
			s.config = gconfig.NewConfig()
		}
		s.config.SetConfigFile(s.configFile)
		parse = true
	} else if s.config != nil && s.config.ConfigFileUsed() == "" {
		parse = true
	}
	if parse {
		// env binding
		s.config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		s.config.SetEnvPrefix("GMC")
		s.config.AutomaticEnv()
		err = s.config.ReadInConfig()
		if err != nil {
			return
		}
		s.configFile = s.config.ConfigFileUsed()
	}
	for id, cfgfile := range s.attachConfigfiles {
		cfg := gconfig.NewConfig()
		cfg.SetConfigFile(cfgfile)
		err = cfg.ReadInConfig()
		if err != nil {
			return
		}
		s.attachConfig[id] = cfg
	}
	return
}
func (s *GMCApp) callRunE(fns []func(*gconfig.Config) error) (err error) {
	hasError := false
	for _, fn := range fns {
		func() {
			defer gcore.Recover(func(e interface{}) {
				hasError = true
				err = gerr.Wrap(e)
			})
			err = fn(s.config)
			if err != nil {
				hasError = true
			}
		}()
		if hasError {
			break
		}
	}
	return
}
func (s *GMCApp) Run() (err error) {
	err = s.parseConfigFile()
	if err != nil {
		return
	}
	err = s.initialize()
	if err != nil {
		return
	}
	// on run
	err = s.callRunE(s.onRun)
	if err != nil {
		return
	}
	err = s.run()
	if err != nil {
		return
	}
	s.reloadSignalMonitor()
	s.logger.Infof("gmc app started done.")
	ghook.RegistShutdown(func() {
		s.Stop()
	})
	if s.isBlock {
		ghook.WaitShutdown()
	} else {
		go ghook.WaitShutdown()
	}
	return
}
func (s *GMCApp) Stop() {
	for _, fn := range s.onShutdown {
		func() {
			defer gcore.Recover(func(e interface{}) {
				s.logger.Infof("run beforeShutdown hook fail, error : %s", gerr.Stack(e))
			})
			fn()
		}()
	}
	for _, srv := range s.services {
		srv.Service.Stop()
	}
}

func (s *GMCApp) OnRun(fn func(*gconfig.Config) (err error)) {
	s.onRun = append(s.onRun, fn)
}
func (s *GMCApp) OnShutdown(fn func()) {
	s.onShutdown = append(s.onShutdown, fn)
}
func (s *GMCApp) SetBlock(isBlockRun bool) {
	s.isBlock = isBlockRun
}
func (s *GMCApp) AddService(item gcore.ServiceItem) {
	item.Service.SetLog(s.logger)
	s.services = append(s.services, item)
}

func (s *GMCApp) SetLogger(logger gcore.Logger) {
	s.logger = logger
}

func (s *GMCApp) Logger() gcore.Logger {
	return s.logger
}

// run all services
func (s *GMCApp) run() (err error) {
	isReload := os.Getenv("GMC_REALOD") == "yes"
	data := os.Getenv("GMC_REALOD_DATA")
	fdMap := map[int]map[int]bool{}
	json.Unmarshal([]byte(data), &fdMap)
	for i, srvI := range s.services {
		srv := srvI.Service
		var cfg *gconfig.Config
		if srvI.ConfigID != "" {
			cfg = s.attachConfig[srvI.ConfigID]
		} else {
			cfg = s.config
		}
		//BeforeInit
		if srvI.BeforeInit != nil {
			err = srvI.BeforeInit(srv, cfg)
			if err != nil {
				return
			}
		}

		//reload checking
		if isReload && len(fdMap[i]) > 0 {
			// fmt.Println(fdMap)
			listeners := []net.Listener{}
			for k := range fdMap[i] {
				listener, e := net.FileListener(os.NewFile(uintptr(k)+3, ""))
				if e != nil {
					err = fmt.Errorf("reload fail, %s", e)
					return
				}
				listeners = append(listeners, listener)
			}
			srv.InjectListeners(listeners)
		}

		//init service
		err = srv.Init(cfg)
		if err != nil {
			return
		}
		srv.SetLog(s.logger)

		//AfterInit
		if srvI.AfterInit != nil {
			err = srvI.AfterInit(&srvI)
			if err != nil {
				return
			}
		}
		//run service
		err = srv.Start()
		if err != nil {
			return
		}
	}
	return
}
