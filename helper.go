package gmc

import (
	"github.com/snail007/gmc/core"
	grouter "github.com/snail007/gmc/http/router"
	ghttpserver "github.com/snail007/gmc/http/server"
	gapp "github.com/snail007/gmc/module/app"
	"github.com/snail007/gmc/module/cache"
	gdb "github.com/snail007/gmc/module/db"
	gerr "github.com/snail007/gmc/module/error"
	gi18n "github.com/snail007/gmc/module/i18n"
	gcaptcha "github.com/snail007/gmc/util/captcha"
	gconfig "github.com/snail007/gmc/util/config"
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
)

// #########################################
// # New stuff helper
// #########################################
type NewAssistant struct {
}

// Config creates a new object of gconfig.Config
func (s *NewAssistant) Config() *gconfig.Config {
	return gconfig.NewConfig()
}

// ConfigFile creates a new object of gconfig.Config from a toml file.
func (s *NewAssistant) ConfigFile(file string) (cfg *gconfig.Config, err error) {
	return gconfig.NewConfigFile(file)
}

// App creates an new object of gcore.GMCApp
func (s *NewAssistant) App() gcore.GMCApp {
	return gapp.New()
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
func (s *NewAssistant) Tr(lang string) *gi18n.I18nTool {
	tool := gi18n.NewI18nTool(gi18n.Tr)
	tool.Lang(lang)
	return tool
}

// AppDefault creates a new object of APP and search config file locations:
// ./app.toml or ./conf/app.toml or ./config/app.toml
func (s *NewAssistant) AppDefault() gcore.GMCApp {
	return gapp.Default()
}

// Router creates a new object of gcore.HTTPRouter
func (s *NewAssistant) Router(ctx gcore.Ctx) *grouter.HTTPRouter {
	return grouter.NewHTTPRouter(ctx)
}

// HTTPServer creates a new object of gmc.HTTPServer
func (s *NewAssistant) HTTPServer(ctx gcore.Ctx) *ghttpserver.HTTPServer {
	return ghttpserver.NewHTTPServer(ctx)
}

// APIServer creates a new object of gmc.APIServer
func (s *NewAssistant) APIServer(ctx gcore.Ctx, address string) *ghttpserver.APIServer {
	return ghttpserver.NewAPIServer(ctx, address)
}

// APIServer creates a new object of gmc.APIServer and initialized from app.toml [apiserver] section.
// cfg is a gconfig.Config object contains section [apiserver] in `app.toml`.
func (s *NewAssistant) APIServerDefault(ctx gcore.Ctx, cfg *gconfig.Config) (api *ghttpserver.APIServer, err error) {
	return ghttpserver.NewDefaultAPIServer(ctx, cfg)
}

// Map creates a gmap.Map object, gmap.Map's keys are sequenced.
func (s *NewAssistant) Map() *_map.Map {
	return _map.NewMap()
}

// Error creates a gmc.Error from an error or string, the gmc.Error
// keeps the full stack information.
func (s *NewAssistant) Error(e interface{}) error {
	return gerr.New(e)
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
func (s *DBAssistant) Init(cfg *gconfig.Config) error {
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
func (s *CacheAssistant) Init(cfg *gconfig.Config) error {
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
func (s *I18nAssistant) Init(cfg *gconfig.Config) error {
	return gi18n.Init(cfg)
}

// Tr search the key in the `lang` i18n file, if not found, then search the
// `fallback` (default) lang file, if both fail `defaultMessage` will be returned. You must be call Init firstly.
func (s *I18nAssistant) Tr(lang, key string, defaultMessage ...string) string {
	return gi18n.Tr(lang, key, defaultMessage...)
}
