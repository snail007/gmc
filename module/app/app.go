package gapp

import (
	"encoding/json"
	"fmt"
	"github.com/snail007/gmc/core"
	ghook "github.com/snail007/gmc/util/process/hook"
	"net"
	"os"
	"strings"
)

type GMCApp struct {
	onRun             []func(gcore.Config) error
	onShutdown        []func()
	isBlock           bool
	attachConfig      map[string]gcore.Config
	attachConfigfiles map[string]string
	services          []gcore.ServiceItem
	logger            gcore.Logger
	configFile        string
	config            gcore.Config
	ctx               gcore.Ctx
}

func (s *GMCApp) Ctx() gcore.Ctx {
	return s.ctx
}

func (s *GMCApp) SetCtx(ctx gcore.Ctx) {
	s.ctx = ctx
}

func NewApp(isDefault bool) gcore.App {
	if isDefault {
		return Default()
	}
	return New()
}

func New() gcore.App {
	app := &GMCApp{
		isBlock:           true,
		onRun:             []func(gcore.Config) error{},
		onShutdown:        []func(){},
		services:          []gcore.ServiceItem{},
		logger:            nil,
		attachConfig:      map[string]gcore.Config{},
		attachConfigfiles: map[string]string{},
	}
	c := gcore.Providers.Ctx("")()
	c.SetApp(app)
	app.ctx = c
	return app
}

func Default() gcore.App {
	// default config dir and name
	cfg := gcore.Providers.Config("")()
	cfg.SetConfigName("app")
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("conf")
	cfg.AddConfigPath("config")
	a := New()
	a.SetConfig(cfg)
	c := gcore.Providers.Ctx("")()
	c.SetApp(a)
	a.(*GMCApp).ctx = c
	return a
}
func (s *GMCApp) initialize() (err error) {
	defer func() {
		if s.logger == nil {
			s.logger = gcore.Providers.Logger("")(nil, "")
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
		s.logger = gcore.Providers.Logger("")(s.ctx, "")
	}

	// initialize database
	if s.config.Sub("database") != nil {
		_, err = gcore.Providers.Database("")(s.ctx)
		if err != nil {
			return
		}
	}

	// initialize cache
	if s.config.Sub("cache") != nil {
		_, err = gcore.Providers.Cache("")(s.ctx)
		if err != nil {
			return
		}
	}

	// initialize i18n if needed
	if s.config.Sub("i18n") != nil {
		var i18n gcore.I18n
		i18n, err = gcore.Providers.I18n("")(s.ctx)
		if err != nil {
			return err
		}
		s.ctx.SetI18n(i18n)
	}
	return
}

func (s *GMCApp) SetConfigFile(file string) {
	s.configFile = file
}
func (s *GMCApp) SetConfig(cfg gcore.Config) {
	s.config = cfg
}
func (s *GMCApp) AttachConfigFile(id, file string) {
	s.attachConfigfiles[id] = file
}
func (s *GMCApp) AttachConfig(id string, cfg gcore.Config) {
	s.attachConfig[id] = cfg
}

//Config acquires the  or attach config object.
//if `idanem` is empty , it return   config object,
//other return attach config object of `id`.
func (s *GMCApp) Config(id ...string) gcore.Config {
	if len(id) > 0 {
		return s.attachConfig[id[0]]
	}
	return s.config
}
func (s *GMCApp) parseConfigFile() (err error) {
	parse := false
	if s.configFile != "" {
		if s.config == nil {
			s.config = gcore.Providers.Config("")()
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
		cfg := gcore.Providers.Config("")()
		cfg.SetConfigFile(cfgfile)
		err = cfg.ReadInConfig()
		if err != nil {
			return
		}
		s.attachConfig[id] = cfg
	}
	return
}
func (s *GMCApp) callRunE(fns []func(gcore.Config) error) (err error) {
	hasError := false
	for _, fn := range fns {
		func() {
			defer gcore.Providers.Error("")().Recover(func(e interface{}) {
				hasError = true
				err = gcore.Providers.Error("")().Wrap(e)
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
			defer gcore.Providers.Error("")().Recover(func(e interface{}) {
				s.logger.Infof("run beforeShutdown hook fail, error : %s", gcore.Providers.Error("")().StackError(e))
			})
			fn()
		}()
	}
	for _, srv := range s.services {
		srv.Service.Stop()
	}
}

func (s *GMCApp) OnRun(fn func(gcore.Config) (err error)) {
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
		var cfg gcore.Config
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
