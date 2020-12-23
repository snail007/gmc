package gcore

import gconfig "github.com/snail007/gmc/config"

type GMCApp interface {
	SetConfigFile(file string) GMCApp
	SetConfig(cfg *gconfig.Config) GMCApp
	AttachConfigFile(id, file string) *GMCApp
	AttachConfig(id string, cfg *gconfig.Config) *GMCApp
	Config(id ...string) *gconfig.Config
	Run() (err error)
	OnRun(fn func(*gconfig.Config) (err error)) *GMCApp
	OnShutdown(fn func()) *GMCApp
	SetBlock(isBlockRun bool) *GMCApp
	AddService(item ServiceItem) *GMCApp
	SetLogger(logger Logger)
	Logger() Logger
}
