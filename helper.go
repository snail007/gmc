package gmc

import (
	gmcapp "github.com/snail007/gmc/app"
	gmccachehelper "github.com/snail007/gmc/cache/helper"
	gmcconfig "github.com/snail007/gmc/config"
	"github.com/snail007/gmc/core"
	gmcdb "github.com/snail007/gmc/db"
	gmcerr "github.com/snail007/gmc/error"
	gmcrouter "github.com/snail007/gmc/http/router"
	gmchttpserver "github.com/snail007/gmc/http/server"
	gmci18n "github.com/snail007/gmc/i18n"
	gmccaptcha "github.com/snail007/gmc/util/captcha"
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
func (s *New0) Config() *gmcconfig.Config {
	return gmcconfig.New()
}

// ConfigFile creates a new object of gmc.Config from a toml file.
func (s *New0) ConfigFile(file string) (cfg *gmcconfig.Config, err error) {
	return gmcconfig.NewFile(file)
}

// App creates an new object of gmc.APP
func (s *New0) App() *APP {
	return gmcapp.New()
}

// Captcha creates an captcha object
func (s *New0) Captcha() *gmccaptcha.Captcha {
	return gmccaptcha.New()
}

// Captcha creates an captcha object, and sets default configuration
func (s *New0) CaptchaDefault() *gmccaptcha.Captcha {
	return gmccaptcha.NewDefault()
}

// Tr creates an new object of gmci18n.I18nTool
// This only worked after gmc.I18n.Init() called.
// lang is the language translation to.
func (s *New0) Tr(lang string) *gmci18n.I18nTool {
	tool := gmci18n.NewI18nTool(gmci18n.Tr)
	tool.Lang(lang)
	return tool
}

// AppDefault creates a new object of APP and search config file locations:
// ./app.toml or ./conf/app.toml or ./config/app.toml
func (s *New0) AppDefault() *APP {
	return gmcapp.Default()
}

// Router creates a new object of gmc.Router
func (s *New0) Router() *gmcrouter.HTTPRouter {
	return gmcrouter.NewHTTPRouter()
}

// HTTPServer creates a new object of gmc.HTTPServer
func (s *New0) HTTPServer() *gmchttpserver.HTTPServer {
	return gmchttpserver.New()
}

// APIServer creates a new object of gmc.APIServer
func (s *New0) APIServer(address string) *gmchttpserver.APIServer {
	return gmchttpserver.NewAPIServer(address)
}

// APIServer creates a new object of gmc.APIServer and initialized from app.toml [apiserver] section.
// cfg is a gmc.Config object contains section [apiserver] in `app.toml`.
func (s *New0) APIServerDefault(cfg *gmcconfig.Config) (api *gmchttpserver.APIServer, err error) {
	return gmchttpserver.NewDefaultAPIServer(cfg)
}

// Map creates a gmc.Map object, gmc.Map's keys are sequenced.
func (s *New0) Map() *_map.Map {
	return _map.NewMap()
}

// Error creates a gmc.Error from an error or string, the gmc.Error
// keeps the full stack information.
func (s *New0) Error(e interface{}) error {
	return gmcerr.New(e)
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
func (s *DB0) Init(cfg *gmcconfig.Config) error {
	return gmcdb.Init(cfg)
}

// SQLite3DB acquires the default db group object, you must be call Init firstly.
// And you must assert it to the correct type to use.
func (s *DB0) DB(id ...string) gmcdb.Database {
	return gmcdb.DB(id...)
}

// MySQL acquires the mysql db group object, you must be call Init firstly.
func (s *DB0) MySQL(id ...string) *gmcdb.MySQLDB {
	return gmcdb.DBMySQL(id...)
}

// SQLite3 acquires the sqlite3 db group object, you must be call Init firstly.
func (s *DB0) SQLite3(id ...string) *gmcdb.SQLite3DB {
	return gmcdb.DBSQLite3(id...)
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
func (s *Cache0) Init(cfg *gmcconfig.Config) error {
	return gmccachehelper.Init(cfg)
}

// Cache acquires the default cache object, you must be call Init firstly.
func (s *Cache0) Cache(id ...string) gmccore.Cache {
	return gmccachehelper.Cache(id...)
}

// Redis acquires the default redis cache object, you must be call Init firstly.
func (s *Cache0) Redis(id ...string) gmccore.Cache {
	return gmccachehelper.Redis(id...)
}

// File acquires the default file cache object, you must be call Init firstly.
func (s *Cache0) File(id ...string) gmccore.Cache {
	return gmccachehelper.File(id...)
}

// Memory acquires the default memory cache object, you must be call Init firstly.
func (s *Cache0) Memory(id ...string) gmccore.Cache {
	return gmccachehelper.Memory(id...)
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
func (s *I18n0) Init(cfg *gmcconfig.Config) error {
	return gmci18n.Init(cfg)
}

// Tr search the key in the `lang` i18n file, if not found, then search the
// `fallback` (default) lang file, if both fail `defaultMessage` will be returned. You must be call Init firstly.
func (s *I18n0) Tr(lang, key string, defaultMessage ...string) string {
	return gmci18n.Tr(lang, key, defaultMessage...)
}


