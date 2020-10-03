package gmcapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	gmcconfig "github.com/snail007/gmc/config"

	gmcerr "github.com/snail007/gmc/error"
	gmchook "github.com/snail007/gmc/process/hook"
	gmcservice "github.com/snail007/gmc/service"
	"github.com/snail007/gmc/util/logutil"
)

var (
	onRunOnceFlag = false
)

type GMCApp struct {
	onRunOnce         []func(config *gmcconfig.Config) error
	onRun             []func(*gmcconfig.Config) error
	onShutdown        []func()
	isBlock           bool
	attachConfig      map[string]*gmcconfig.Config
	attachConfigfiles map[string]string
	services          []ServiceItem
	logger            *log.Logger
	configFile        string
	config            *gmcconfig.Config
}
type ServiceItem struct {
	BeforeInit   func(srv *gmcconfig.Config) (err error)
	AfterInit    func(srv *ServiceItem) (err error)
	Service      gmcservice.Service
	ConfigIDname string
}

func New() *GMCApp {
	return &GMCApp{
		isBlock:           true,
		onRunOnce:         []func(*gmcconfig.Config) error{},
		onRun:             []func(*gmcconfig.Config) error{},
		onShutdown:        []func(){},
		services:          []ServiceItem{},
		logger:            logutil.New(""),
		attachConfig:      map[string]*gmcconfig.Config{},
		attachConfigfiles: map[string]string{},
	}
}

func (s *GMCApp) SetConfigFile(file string) *GMCApp {
	s.configFile = file
	return s
}
func (s *GMCApp) SetConfig(cfg *gmcconfig.Config) *GMCApp {
	s.config = cfg
	return s
}
func (s *GMCApp) AttachConfigFile(idname, file string) *GMCApp {
	s.attachConfigfiles[idname] = file
	return s
}
func (s *GMCApp) AttachConfig(idname string, cfg *gmcconfig.Config) *GMCApp {
	s.attachConfig[idname] = cfg
	return s
}

//Config acquires the  or attach config object.
//if `idanem` is empty , it return   config object,
//other return attach config object of `idname`.
func (s *GMCApp) Config(idname ...string) *gmcconfig.Config {
	if len(idname) > 0 {
		return s.attachConfig[idname[0]]
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
		err = s.config.ReadInConfig()
		if err != nil {
			return
		}
		s.configFile = s.config.ConfigFileUsed()
	}
	for idname, cfgfile := range s.attachConfigfiles {
		cfg := gmcconfig.New()
		cfg.SetConfigFile(cfgfile)
		err = cfg.ReadInConfig()
		if err != nil {
			return
		}
		s.attachConfig[idname] = cfg
	}
	return
}
func (s *GMCApp) callRunE(fns []func(*gmcconfig.Config) error) (err error) {
	hasError := false
	for _, fn := range fns {
		func() {
			defer func() {
				if e := recover(); e != nil {
					hasError = true
					err = gmcerr.Wrap(e)
				}
			}()
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
	// on run once
	if !onRunOnceFlag {
		onRunOnceFlag = true
		err = s.callRunE(s.onRunOnce)
		if err != nil {
			return
		}
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
	s.logger.Printf("gmc app start done.")
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
				if e := recover(); e != nil {
					s.logger.Printf("run beforeShutdown hook fail, error : %s", gmcerr.Stack(e))
				}
			}()
			fn()
		}()
	}
	for _, srv := range s.services {
		srv.Service.Stop()
	}
}
func (s *GMCApp) OnRunOnce(fn func(*gmcconfig.Config) (err error)) *GMCApp {
	s.onRunOnce = append(s.onRunOnce, fn)
	return s
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
func (s *GMCApp) Logger() *log.Logger {
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
		if srvI.ConfigIDname != "" {
			cfg = s.attachConfig[srvI.ConfigIDname]
		} else {
			cfg = s.config
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
