package gmcapp

import (
	"encoding/json"
	"fmt"
	gmclog "github.com/snail007/gmc/base/log"
	gmccachehelper "github.com/snail007/gmc/cache/helper"
	gmcconfig "github.com/snail007/gmc/config"
	"github.com/snail007/gmc/core"
	gmcdbhelper "github.com/snail007/gmc/db/helper"
	gmcerr "github.com/snail007/gmc/error"
	gmci18n "github.com/snail007/gmc/i18n"
	gmchook "github.com/snail007/gmc/process/hook"
	"net"
	"os"
	"strings"
)

type GMCApp struct {
	onRun             []func(*gmcconfig.Config) error
	onShutdown        []func()
	isBlock           bool
	attachConfig      map[string]*gmcconfig.Config
	attachConfigfiles map[string]string
	services          []ServiceItem
	logger            gmccore.Logger
	configFile        string
	config            *gmcconfig.Config
}

type ServiceItem struct {
	BeforeInit func(srv gmccore.Service, cfg *gmcconfig.Config) (err error)
	AfterInit  func(srv *ServiceItem) (err error)
	Service    gmccore.Service
	ConfigID   string
}

func New() *GMCApp {
	return &GMCApp{
		isBlock:           true,
		onRun:             []func(*gmcconfig.Config) error{},
		onShutdown:        []func(){},
		services:          []ServiceItem{},
		logger:            nil,
		attachConfig:      map[string]*gmcconfig.Config{},
		attachConfigfiles: map[string]string{},
	}
}

func Default() *GMCApp {
	// default config dir and name
	cfg := gmcconfig.New()
	cfg.SetConfigName("app")
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("conf")
	cfg.AddConfigPath("config")
	return New().SetConfig(cfg)
}
func (s *GMCApp) initialize() (err error) {
	defer func() {
		if s.logger == nil {
			s.logger = gmclog.NewGMCLog()
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
		s.logger = gmclog.NewFromConfig(s.config)
	}

	// initialize database
	if s.config.Sub("database") != nil {
		err = gmcdbhelper.Init(s.config)
		if err != nil {
			return
		}
	}

	// initialize cache
	if s.config.Sub("cache") != nil {
		err = gmccachehelper.Init(s.config)
		if err != nil {
			return
		}
	}
	// initialize i18n if needed
	if s.config.Sub("i18n") != nil {
		err = gmci18n.Init(s.config)
	}
	return
}
func (s *GMCApp) SetConfigFile(file string) *GMCApp {
	s.configFile = file
	return s
}
func (s *GMCApp) SetConfig(cfg *gmcconfig.Config) *GMCApp {
	s.config = cfg
	return s
}
func (s *GMCApp) AttachConfigFile(id, file string) *GMCApp {
	s.attachConfigfiles[id] = file
	return s
}
func (s *GMCApp) AttachConfig(id string, cfg *gmcconfig.Config) *GMCApp {
	s.attachConfig[id] = cfg
	return s
}

//Config acquires the  or attach config object.
//if `idanem` is empty , it return   config object,
//other return attach config object of `id`.
func (s *GMCApp) Config(id ...string) *gmcconfig.Config {
	if len(id) > 0 {
		return s.attachConfig[id[0]]
	}
	return s.config
}
func (s *GMCApp) parseConfigFile() (err error) {
	parse := false
	if s.configFile != "" {
		if s.config == nil {
			s.config = gmcconfig.New()
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
		cfg := gmcconfig.New()
		cfg.SetConfigFile(cfgfile)
		err = cfg.ReadInConfig()
		if err != nil {
			return
		}
		s.attachConfig[id] = cfg
	}
	return
}
func (s *GMCApp) callRunE(fns []func(*gmcconfig.Config) error) (err error) {
	hasError := false
	for _, fn := range fns {
		func() {
			defer gmcerr.Recover(func(e interface{}){
				hasError = true
				err = gmcerr.Wrap(e)
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
	gmchook.RegistShutdown(func() {
		s.Stop()
	})
	if s.isBlock {
		gmchook.WaitShutdown()
	} else {
		go gmchook.WaitShutdown()
	}
	return
}
func (s *GMCApp) Stop() {
	for _, fn := range s.onShutdown {
		func() {
			defer func() {
				gmcerr.Recover(func(e interface{}){
					s.logger.Infof("run beforeShutdown hook fail, error : %s", gmcerr.Stack(e))
				})
			}()
			fn()
		}()
	}
	for _, srv := range s.services {
		srv.Service.Stop()
	}
}

func (s *GMCApp) OnRun(fn func(*gmcconfig.Config) (err error)) *GMCApp {
	s.onRun = append(s.onRun, fn)
	return s
}
func (s *GMCApp) OnShutdown(fn func()) *GMCApp {
	s.onShutdown = append(s.onShutdown, fn)
	return s
}
func (s *GMCApp) Block(isBlockRun bool) *GMCApp {
	s.isBlock = isBlockRun
	return s
}
func (s *GMCApp) AddService(item ServiceItem) *GMCApp {
	item.Service.SetLog(s.logger)
	s.services = append(s.services, item)
	return s
}

func (s *GMCApp) SetLogger(logger gmccore.Logger) {
	s.logger = logger
}

func (s *GMCApp) Logger() gmccore.Logger {
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
		var cfg *gmcconfig.Config
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
