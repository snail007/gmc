package gcore

import gconfig "github.com/snail007/gmc/util/config"

type GMCApp interface {
	SetConfigFile(file string)
	SetConfig(cfg *gconfig.Config)
	AttachConfigFile(id, file string)
	AttachConfig(id string, cfg *gconfig.Config)
	Config(id ...string) *gconfig.Config
	Run() (err error)
	OnRun(fn func(*gconfig.Config) (err error))
	OnShutdown(fn func())
	SetBlock(isBlockRun bool)
	AddService(item ServiceItem)
	SetLogger(logger Logger)
	Logger() Logger
	Stop()
	Ctx() Ctx
}
