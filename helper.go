package gmc

import (
	gconfig "github.com/snail007/gmc/config"
	"github.com/snail007/gmc/core"
	gapp "github.com/snail007/gmc/gmc/app"
	gcachehelper "github.com/snail007/gmc/gmc/cache/helper"
	gdb "github.com/snail007/gmc/gmc/db"
	gerr "github.com/snail007/gmc/gmc/error"
	gi18n "github.com/snail007/gmc/gmc/i18n"
	grouter "github.com/snail007/gmc/http/router"
	ghttpserver "github.com/snail007/gmc/http/server"
	gcaptcha "github.com/snail007/gmc/util/captcha"
	_map "github.com/snail007/gmc/util/map"
)

var (
	// New shortcut of New gmc stuff
	New = &New0{}
	// SQLite3DB shortcut of gmc database stuff
	DB = &DB0{}
	// Cache shortcut of gmc cache stuff
	Cache = &Cache0{}
	// I18n shortcut of gmc i18n
	I18n = &I18n0{}
)

// #########################################
// # New stuff helper
// #########################################
type New0 struct {
}

// Config creates a new object of gmc.Config
func (s *New0) Config() *gconfig.Config {
	return gconfig.NewConfig()
}

// ConfigFile creates a new object of gmc.Config from a toml file.
func (s *New0) ConfigFile(file string) (cfg *gconfig.Config, err error) {
	return gconfig.NewConfigFile(file)
}

// App creates an new object of gmc.APP
func (s *New0) App() gcore.GMCApp {
	return gapp.New()
}

// Captcha creates an captcha object
func (s *New0) Captcha() *gcaptcha.Captcha {
	return gcaptcha.New()
}

// Captcha creates an captcha object, and sets default configuration
func (s *New0) CaptchaDefault() *gcaptcha.Captcha {
	return gcaptcha.NewDefault()
}

// Tr creates an new object of gi18n.I18nTool
// This only worked after gmc.I18n.Init() called.
// lang is the language translation to.
func (s *New0) Tr(lang string) *gi18n.I18nTool {
	tool := gi18n.NewI18nTool(gi18n.Tr)
	tool.Lang(lang)
	return tool
}

// AppDefault creates a new object of APP and search config file locations:
// ./app.toml or ./conf/app.toml or ./config/app.toml
func (s *New0) AppDefault() gcore.GMCApp {
	return gapp.Default()
}

// Router creates a new object of gmc.Router
func (s *New0) Router() *grouter.HTTPRouter {
	return grouter.NewHTTPRouter()
}

// HTTPServer creates a new object of gmc.HTTPServer
func (s *New0) HTTPServer() *ghttpserver.HTTPServer {
	return ghttpserver.New()
}

// APIServer creates a new object of gmc.APIServer
func (s *New0) APIServer(address string) *ghttpserver.APIServer {
	return ghttpserver.NewAPIServer(address)
}

// APIServer creates a new object of gmc.APIServer and initialized from app.toml [apiserver] section.
// cfg is a gmc.Config object contains section [apiserver] in `app.toml`.
func (s *New0) APIServerDefault(cfg *gconfig.Config) (api *ghttpserver.APIServer, err error) {
	return ghttpserver.NewDefaultAPIServer(cfg)
}

// Map creates a gmc.Map object, gmc.Map's keys are sequenced.
func (s *New0) Map() *_map.Map {
	return _map.NewMap()
}

// Error creates a gmc.Error from an error or string, the gmc.Error
// keeps the full stack information.
func (s *New0) Error(e interface{}) error {
	return gerr.New(e)
}

// ##################################################
// # Database helper
// # Init Must be called firstly with config object
// # of app.toml.
// ##################################################
type DB0 struct {
}

// Init initialize the db group objects from a config object
// contains app.toml section [database].
func (s *DB0) Init(cfg *gconfig.Config) error {
	return gdb.Init(cfg)
}

// InitFromFile initialize the db group objects from a foo.toml config file.
// foo.toml must be contains section [database].
func (s *DB0) InitFromFile(cfgFile string) error {
	return gdb.InitFromFile(cfgFile)
}

// SQLite3DB acquires the default db group object, you must be call Init firstly.
// And you must assert it to the correct type to use.
func (s *DB0) DB(id ...string) gcore.Database {
	return gdb.DB(id...)
}

// MySQL acquires the mysql db group object, you must be call Init firstly.
func (s *DB0) MySQL(id ...string) *gdb.MySQLDB {
	return gdb.DBMySQL(id...)
}

// SQLite3 acquires the sqlite3 db group object, you must be call Init firstly.
func (s *DB0) SQLite3(id ...string) *gdb.SQLite3DB {
	return gdb.DBSQLite3(id...)
}

// #################################################
// # Cache helper
// # Init Must be called firstly with config object
// # of app.toml.
// ##################################################
type Cache0 struct {
}

// Init initialize the cache group objects from a config object
// contains app.toml section [cache].
func (s *Cache0) Init(cfg *gconfig.Config) error {
	return gcachehelper.Init(cfg)
}

// Cache acquires the default cache object, you must be call Init firstly.
func (s *Cache0) Cache(id ...string) gcore.Cache {
	return gcachehelper.Cache(id...)
}

// Redis acquires the default redis cache object, you must be call Init firstly.
func (s *Cache0) Redis(id ...string) gcore.Cache {
	return gcachehelper.Redis(id...)
}

// File acquires the default file cache object, you must be call Init firstly.
func (s *Cache0) File(id ...string) gcore.Cache {
	return gcachehelper.File(id...)
}

// Memory acquires the default memory cache object, you must be call Init firstly.
func (s *Cache0) Memory(id ...string) gcore.Cache {
	return gcachehelper.Memory(id...)
}

// ##################################################
// # I18n helper
// # Init Must be called firstly with config object
// # of app.toml.
// ##################################################
type I18n0 struct {
}

// Init initialize the i18n object from a config object
// contains app.toml section [i18n].
func (s *I18n0) Init(cfg *gconfig.Config) error {
	return gi18n.Init(cfg)
}

// Tr search the key in the `lang` i18n file, if not found, then search the
// `fallback` (default) lang file, if both fail `defaultMessage` will be returned. You must be call Init firstly.
func (s *I18n0) Tr(lang, key string, defaultMessage ...string) string {
	return gi18n.Tr(lang, key, defaultMessage...)
}

