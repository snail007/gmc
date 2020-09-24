package gmcapp

import (
	"fmt"
	"log"

	gmchook "github.com/snail007/gmc/process/hook"

	gmcservice "github.com/snail007/gmc/service"

	gmcconfig "github.com/snail007/gmc/config/gmc"
	"github.com/snail007/gmc/util/logutil"
)

type GMCApp struct {
	afterParse       []func(*gmcconfig.GMCConfig) error
	beforeRun        []func(*gmcconfig.GMCConfig) error
	beforeShutdown   []func()
	isBlock          bool
	isParsed         bool
	extraConfig      map[string]*gmcconfig.GMCConfig
	extraConfigfiles map[string]string
	services         []ServiceItem
	logger           *log.Logger
	mainConfigFile   string
	mainConfig       *gmcconfig.GMCConfig
	afterServiceInit []func(app *GMCApp, srv interface{}) error
}
type ServiceItem struct {
	BeforeInit   func(srv *gmcconfig.GMCConfig) (err error)
	AfterInit    func(srv *ServiceItem) (err error)
	Service      gmcservice.Service
	ConfigIDname string
}

func New() *GMCApp {
	return &GMCApp{
		isBlock:          true,
		services:         []ServiceItem{},
		logger:           logutil.New("[gmc]"),
		extraConfig:      map[string]*gmcconfig.GMCConfig{},
		extraConfigfiles: map[string]string{},
	}
}

func (s *GMCApp) SetMainConfigFile(file string) *GMCApp {
	s.mainConfigFile = file
	return s
}
func (s *GMCApp) AddExtraConfigFile(idname, file string) *GMCApp {
	s.extraConfigfiles[idname] = file
	return s
}

//Config acquires the main or extra config object.
//if `idanem` is empty , it return main  config object,
//other return extra config object of `idname`.
func (s *GMCApp) Config(idname ...string) *gmcconfig.GMCConfig {
	if len(idname) > 0 {
		return s.extraConfig[idname[0]]
	}
	return s.mainConfig
}
func (s *GMCApp) ParseConfig() (err error) {
	if s.isParsed {
		return
	}
	defer func() {
		if err == nil {
			s.isParsed = true
			err = s.callRunE(s.afterParse)
		}
	}()
	//create config file object
	s.mainConfig = gmcconfig.New()
	if s.mainConfigFile == "" {
		s.mainConfig.SetConfigName("app")
		s.mainConfig.AddConfigPath(".")
		s.mainConfig.AddConfigPath("config")
		s.mainConfig.AddConfigPath("conf")
		//for testing
		// s.config.AddConfigPath("../app")
	} else {
		s.mainConfig.SetConfigFile(s.mainConfigFile)
	}
	err = s.mainConfig.ReadInConfig()
	if err != nil {
		return
	}
	s.mainConfigFile = s.mainConfig.ConfigFileUsed()
	for idname, cfgfile := range s.extraConfigfiles {
		cfg := gmcconfig.New()
		cfg.SetConfigFile(cfgfile)
		err = cfg.ReadInConfig()
		if err != nil {
			return
		}
		s.extraConfig[idname] = cfg
	}
	return
}
func (s *GMCApp) callRunE(fns []func(*gmcconfig.GMCConfig) error) (err error) {
	hasError := false
	for _, fn := range fns {
		func() {
			defer func() {
				if e := recover(); e != nil {
					hasError = true
					err = fmt.Errorf("%s", e)
				}
			}()
			err = fn(s.mainConfig)
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
	err = s.callRunE(s.beforeRun)
	if err != nil {
		return
	}
	err = s.run()
	if err != nil {
		return
	}
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
	for _, fn := range s.beforeShutdown {
		func() {
			defer func() {
				if e := recover(); e != nil {
					s.logger.Printf("run beforeShutdown hook fail, error : %s", e)
				}
			}()
			fn()
		}()
	}
	for _, srv := range s.services {
		srv.Service.Stop()
	}
}
func (s *GMCApp) AfterParse(fn func(*gmcconfig.GMCConfig) (err error)) *GMCApp {
	s.afterParse = append(s.afterParse, fn)
	return s
}
func (s *GMCApp) BeforeRun(fn func(*gmcconfig.GMCConfig) (err error)) *GMCApp {
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
func (s *GMCApp) AddService(item ServiceItem) *GMCApp {
	item.Service.SetLog(s.logger)
	s.services = append(s.services, item)
	return s
}
func (s *GMCApp) Logger() *log.Logger {
	return s.logger
}

// run all services
func (s *GMCApp) run() (err error) {
	for _, srvI := range s.services {
		srv := srvI.Service
		var cfg *gmcconfig.GMCConfig
		if srvI.ConfigIDname != "" {
			cfg = s.extraConfig[srvI.ConfigIDname]
		} else {
			cfg = s.mainConfig
		}
		//BeforeInit
		if srvI.BeforeInit != nil {
			err = srvI.BeforeInit(cfg)
			if err != nil {
				return
			}
		}
		//init service
		err = srv.Init(cfg)
		if err != nil {
			return
		}
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
