// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmc

import (
	"fmt"
	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/module/cache"
	gdb "github.com/snail007/gmc/module/db"
	gi18n "github.com/snail007/gmc/module/i18n"
	gcaptcha "github.com/snail007/gmc/util/captcha"
	gmap "github.com/snail007/gmc/util/map"
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
	// Err shortcut of gerror
	Err = &ErrorAssistant{}
)

// NewAssistant helper to new different gmc objects.
type NewAssistant struct {
}

// Config creates a new object of gconfig.Config
func (s *NewAssistant) Config() gcore.Config {
	return gcore.ProviderConfig()()
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
	return gcore.ProviderCtx()()
}

// Logger creates an new object of gcore.Logger
func (s *NewAssistant) Logger(ctx gcore.Ctx, prefix string) gcore.Logger {
	return gcore.ProviderLogger()(ctx, prefix)
}

// App creates an new object of gcore.App
func (s *NewAssistant) App() gcore.App {
	return gcore.ProviderApp()(false)
}

// Captcha creates an captcha object
func (s *NewAssistant) Captcha() *gcaptcha.Captcha {
	return gcaptcha.New()
}

// CaptchaDefault creates an captcha object, and sets default configuration
func (s *NewAssistant) CaptchaDefault() *gcaptcha.Captcha {
	return gcaptcha.NewDefault()
}

// Tr creates an new object of gi18n.I18nTool
// This only worked after gmc.I18n.Init() called.
// lang is the language translation to.
func (s *NewAssistant) Tr(lang string, ctx gcore.Ctx) (tr gcore.I18n) {
	tr, _ = gcore.ProviderI18n()(ctx)
	tr.Lang(lang)
	return
}

// AppDefault creates a new object of APP and search config file locations:
// ./app.toml or ./conf/app.toml or ./config/app.toml
func (s *NewAssistant) AppDefault() gcore.App {
	return gcore.ProviderApp()(true)
}

// Router creates a new object of gcore.HTTPRouter
func (s *NewAssistant) Router(ctx gcore.Ctx) gcore.HTTPRouter {
	return gcore.ProviderHTTPRouter()(ctx)
}

// HTTPServer creates a new object of gmc.HTTPServer
func (s *NewAssistant) HTTPServer(ctx gcore.Ctx) gcore.HTTPServer {
	return gcore.ProviderHTTPServer()(ctx)
}

// APIServer creates a new object of gmc.APIServer
func (s *NewAssistant) APIServer(ctx gcore.Ctx, address string) (gcore.APIServer, error) {
	return gcore.ProviderAPIServer()(ctx, address)
}

// APIServerDefault creates a new object of gmc.APIServer and initialized from app.toml [apiserver] section.
// cfg is a gconfig.Config object contains section [apiserver] in `app.toml`.
func (s *NewAssistant) APIServerDefault(ctx gcore.Ctx) (gcore.APIServer, error) {
	return gcore.ProviderAPIServer()(ctx, "")
}

// Map creates a gmap.Map object, gmap.Map's keys are sequenced.
func (s *NewAssistant) Map() *gmap.Map {
	return gmap.New()
}

// DBAssistant helper to access database, DB.Init Must be called firstly with config object of app.toml.
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

// DB acquires the default db group object, you must be call Init firstly.
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

// CacheAssistant helper to access cache, Init Must be called firstly with config object of app.toml.
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

// I18nAssistant helper to using i18n, Init Must be called firstly with config object of app.toml.
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

// ErrorAssistant helper to using error or add stack info to error.
type ErrorAssistant struct {
	p gcore.ErrorProvider
}

// NewErrorAssistant returns new a ErrorAssistant object.
func NewErrorAssistant() *ErrorAssistant {
	return &ErrorAssistant{p: gcore.ProviderError()}
}

// Stack acquires an error full stack info string
func (e *ErrorAssistant) Stack(err interface{}) string {
	return e.p().StackError(err)
}

// Recover catch fatal error in defer,
// f can be func(err interface{}) or string or object as string
func (e *ErrorAssistant) Recover(f ...interface{}) {
	if err := recover(); err != nil {
		eObj := e.p().New(nil)
		var f0 interface{}
		var printStack bool
		if len(f) == 0 {
			return
		}
		if len(f) == 2 {
			printStack = f[1].(bool)
		}
		f0 = f[0]
		switch v := f0.(type) {
		case func(interface{}):
			v(err)
		case string:
			s := ""
			if printStack {
				s = fmt.Sprintf(",stack: %s", eObj.StackError(err))
			}
			fmt.Printf("\nrecover error, %v%s\n", f, s)
		default:
			fmt.Printf("\nrecover error %s\n", eObj.Wrap(err).ErrorStack())
		}
	}
}

// New creates a gcore.Error from an error or string, the gcore.Error
// keeps the full stack information.
func (e *ErrorAssistant) New(err interface{}) error {
	return e.p().New(err)
}

// Wrap wraps an error to gcore.Error, which keeps
// the full stack information.
func (e *ErrorAssistant) Wrap(err interface{}) error {
	return e.p().WrapN(err, 2)
}

func initHelper() {
	Err = NewErrorAssistant()
}
