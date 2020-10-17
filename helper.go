package gmc

import (
	gmcapp "github.com/snail007/gmc/app"
	gmccache "github.com/snail007/gmc/cache"
	gmccachehelper "github.com/snail007/gmc/cache/helper"
	gmcconfig "github.com/snail007/gmc/config"
	gmcdbhelper "github.com/snail007/gmc/db/helper"
	gmcmysql "github.com/snail007/gmc/db/mysql"
	gmcsqlite3 "github.com/snail007/gmc/db/sqlite3"
	gmcerr "github.com/snail007/gmc/error"
	gmcrouter "github.com/snail007/gmc/http/router"
	gmchttpserver "github.com/snail007/gmc/http/server"
	gmci18n "github.com/snail007/gmc/i18n"
	"github.com/snail007/gmc/util/maputil"
)

var (
	// New shortcut of New gmc stuff
	New = &New0{}
	// DB shortcut of gmc database stuff
	DB = &DB0{}
	// Cache shortcut of gmc cache stuff
	Cache = &Cache0{}
	// I18n shortcut of gmc i18n
	I18n = &I18n0{}
)

// New stuff helper
type New0 struct {
}

func (s *New0) Config() *gmcconfig.Config {
	return gmcconfig.New()
}
func (s *New0) ConfigFile(file string) (cfg *gmcconfig.Config,err error) {
	return gmcconfig.NewFile(file)
}
func (s *New0) App() *APP {
	return gmcapp.New()
}

func (s *New0) Tr(lang string) *gmci18n.I18nTool {
	tool:=gmci18n.NewI18nTool(gmci18n.Tr)
	tool.Lang(lang)
	return tool
}

func (s *New0) AppDefault() *APP {
	return gmcapp.Default()
}

func (s *New0) Router() *gmcrouter.HTTPRouter {
	return gmcrouter.NewHTTPRouter()
}

func (s *New0) HTTPServer() *gmchttpserver.HTTPServer {
	return gmchttpserver.New()
}

func (s *New0) APIServer(address string) *gmchttpserver.APIServer {
	return gmchttpserver.NewAPIServer(address)
}

func (s *New0) Map() *maputil.Map {
	return maputil.NewMap()
}

func (s *New0) Error(e interface{}) error {
	return gmcerr.New(e)
}

// Database helper
type DB0 struct {
}

func (s *DB0) Init(cfg *gmcconfig.Config) error {
	return gmcdbhelper.Init(cfg)
}

func (s *DB0) DB(id ...string) interface{} {
	return gmcdbhelper.DB(id...)
}

func (s *DB0) MySQL(id ...string) *gmcmysql.DB {
	return gmcdbhelper.DBMySQL(id...)
}

func (s *DB0) SQLite3(id ...string) *gmcsqlite3.DB {
	return gmcdbhelper.DBSQLite3(id...)
}

// Cache helper
type Cache0 struct {
}

func (s *Cache0) Init(cfg *gmcconfig.Config) error {
	return gmccachehelper.Init(cfg)
}

func (s *Cache0) Cache(id ...string) gmccache.Cache {
	return gmccachehelper.Cache(id...)
}

func (s *Cache0) Redis(id ...string) gmccache.Cache {
	return gmccachehelper.Redis(id...)
}

func (s *Cache0) File(id ...string) gmccache.Cache {
	return gmccachehelper.File(id...)
}

func (s *Cache0) Memory(id ...string) gmccache.Cache {
	return gmccachehelper.Memory(id...)
}

// I18n helper
type I18n0 struct {
}

func (s *I18n0) Init(cfg *gmcconfig.Config) error {
	return gmci18n.Init(cfg)
}

func (s *I18n0) Tr(lang, key string, defaultMessage ...string) string{
	return gmci18n.Tr(lang,key,defaultMessage...)
}