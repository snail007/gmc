package gmc

import (
	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/module/cache"
	gdb "github.com/snail007/gmc/module/db"
	gi18n "github.com/snail007/gmc/module/i18n"
	gcaptcha "github.com/snail007/gmc/util/captcha"
	_map "github.com/snail007/gmc/util/map"
)

var (
	// New shortcut of New gmc stuff
	New = &NewAssistant{}
	// DB shortcut of gmc database stuff
	DB = &DBAssistant{}
	// Cache shortcut of gmc cache stuff
	Cache = &CacheAssistant{}
	// I18n shortcut of gmc i18n
	I18n = &I18nAssistant{}
	Err  *ErrorAssistant
)

// #########################################
// # New stuff helper
// #########################################
type NewAssistant struct {
}

// Config creates a new object of gconfig.Config
func (s *NewAssistant) Config() gcore.Config {
	return gcore.Providers.Config("")()
}

// ConfigFile creates a new object of gconfig.Config from a toml file.
func (s *NewAssistant) ConfigFile(file string) (cfg gcore.Config, err error) {
	v := s.Config()
	v.SetConfigFile(file)
	err = v.ReadInConfig()
	if err != nil {
		return
	}
	cfg = v
	return
}

// Ctx creates an new object of gcore.Ctx
func (s *NewAssistant) Ctx() gcore.Ctx {
	return gcore.Providers.Ctx("")()
}

// App creates an new object of gcore.App
func (s *NewAssistant) App() gcore.App {
	return gcore.Providers.App("")(false)
}

// Captcha creates an captcha object
func (s *NewAssistant) Captcha() *gcaptcha.Captcha {
	return gcaptcha.New()
}

// Captcha creates an captcha object, and sets default configuration
func (s *NewAssistant) CaptchaDefault() *gcaptcha.Captcha {
	return gcaptcha.NewDefault()
}

// Tr creates an new object of gi18n.I18nTool
// This only worked after gmc.I18n.Init() called.
// lang is the language translation to.
func (s *NewAssistant) Tr(lang string) (tr gcore.I18n) {
	tr, _ = gcore.Providers.I18n("")(nil)
	tr.Lang(lang)
	return
}

// AppDefault creates a new object of APP and search config file locations:
// ./app.toml or ./conf/app.toml or ./config/app.toml
func (s *NewAssistant) AppDefault() gcore.App {
	return gcore.Providers.App("")(true)
}

// Router creates a new object of gcore.HTTPRouter
func (s *NewAssistant) Router(ctx gcore.Ctx) gcore.HTTPRouter {
	return gcore.Providers.HTTPRouter("")(ctx)
}

// HTTPServer creates a new object of gmc.HTTPServer
func (s *NewAssistant) HTTPServer(ctx gcore.Ctx) gcore.HTTPServer {
	return gcore.Providers.HTTPServer("")(ctx)
}

// APIServer creates a new object of gmc.APIServer
func (s *NewAssistant) APIServer(ctx gcore.Ctx, address string) (gcore.APIServer, error) {
	return gcore.Providers.APIServer("")(ctx, address)
}

// APIServer creates a new object of gmc.APIServer and initialized from app.toml [apiserver] section.
// cfg is a gconfig.Config object contains section [apiserver] in `app.toml`.
func (s *NewAssistant) APIServerDefault(ctx gcore.Ctx) (gcore.APIServer, error) {
	return gcore.Providers.APIServer("")(ctx, "")
}

// Map creates a gmap.Map object, gmap.Map's keys are sequenced.
func (s *NewAssistant) Map() *_map.Map {
	return _map.NewMap()
}

// ##################################################
// # Database helper
// # Init Must be called firstly with config object
// # of app.toml.
// ##################################################
type DBAssistant struct {
}

// Init initialize the db group objects from a config object
// contains app.toml section [database].
func (s *DBAssistant) Init(cfg gcore.Config) error {
	return gdb.Init(cfg)
}

// InitFromFile initialize the db group objects from a foo.toml config file.
// foo.toml must be contains section [database].
func (s *DBAssistant) InitFromFile(cfgFile string) error {
	return gdb.InitFromFile(cfgFile)
}

// SQLite3DB acquires the default db group object, you must be call Init firstly.
// And you must assert it to the correct type to use.
func (s *DBAssistant) DB(id ...string) gcore.Database {
	return gdb.DB(id...)
}

// MySQL acquires the mysql db group object, you must be call Init firstly.
func (s *DBAssistant) MySQL(id ...string) *gdb.MySQLDB {
	return gdb.DBMySQL(id...)
}

// SQLite3 acquires the sqlite3 db group object, you must be call Init firstly.
func (s *DBAssistant) SQLite3(id ...string) *gdb.SQLite3DB {
	return gdb.DBSQLite3(id...)
}

// Table acquires the table model object, you must be call Init firstly.
func (s *DBAssistant) Table(tableName string) *gdb.Model {
	return gdb.Table(tableName)
}

// #################################################
// # Cache helper
// # Init Must be called firstly with config object
// # of app.toml.
// ##################################################
type CacheAssistant struct {
}

// Init initialize the cache group objects from a config object
// contains app.toml section [cache].
func (s *CacheAssistant) Init(cfg gcore.Config) error {
	return gcache.Init(cfg)
}

// Cache acquires the default cache object, you must be call Init firstly.
func (s *CacheAssistant) Cache(id ...string) gcore.Cache {
	return gcache.Cache(id...)
}

// Redis acquires the default redis cache object, you must be call Init firstly.
func (s *CacheAssistant) Redis(id ...string) *gcache.RedisCache {
	return gcache.Redis(id...)
}

// File acquires the default file cache object, you must be call Init firstly.
func (s *CacheAssistant) File(id ...string) *gcache.FileCache {
	return gcache.File(id...)
}

// Memory acquires the default memory cache object, you must be call Init firstly.
func (s *CacheAssistant) Memory(id ...string) *gcache.MemCache {
	return gcache.Memory(id...)
}

// ##################################################
// # I18n helper
// # Init Must be called firstly with config object
// # of app.toml.
// ##################################################
type I18nAssistant struct {
}

// Init initialize the i18n object from a config object
// contains app.toml section [i18n].
func (s *I18nAssistant) Init(cfg gcore.Config) (i18n gcore.I18n, err error) {
	err = gi18n.Init(cfg)
	if err != nil {
		return
	}
	i18n = gi18n.I18N
	return
}

// ##################################################
// # Error helper
// ##################################################
type ErrorAssistant struct {
	p gcore.ErrorProvider
}

func NewErrorAssistant() *ErrorAssistant {
	return &ErrorAssistant{p: gcore.Providers.Error("")}
}

// ErrStack acquires an error full stack info string
func (e *ErrorAssistant) Stack(err interface{}) string {
	return e.p().StackError(err)
}

// ErrRecover catch fatal error in defer,
// f can be func(err interface{}) or string or object as string
func (e *ErrorAssistant) Recover(f ...interface{}) {
	e.p().Recover(f)
}

// New creates a gcore.Error from an error or string, the gcore.Error
// keeps the full stack information.
func (e *ErrorAssistant) New(err interface{}) error {
	return e.p().New(e)
}

// Wrap wraps an error to gcore.Error, which keeps
// the full stack information.
func (e *ErrorAssistant) Wrap(err interface{}) error {
	return e.p().New(e)
}

func initHelper() {
	Err = NewErrorAssistant()
}
