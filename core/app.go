package gcore

type App interface {
	SetConfigFile(file string)
	SetConfig(cfg Config)
	AttachConfigFile(id, file string)
	AttachConfig(id string, cfg Config)
	Config(id ...string) Config
	Run() (err error)
	OnRun(fn func(Config) (err error))
	OnShutdown(fn func())
	SetBlock(isBlockRun bool)
	AddService(item ServiceItem)
	SetLogger(logger Logger)
	Logger() Logger
	Stop()
	Ctx() Ctx
}
