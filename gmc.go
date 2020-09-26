package gmc

import (
	"net/http"

	gmchttputil "github.com/snail007/gmc/util/httputil"

	gmcapp "github.com/snail007/gmc/app"
	gmccache "github.com/snail007/gmc/cache"
	gmccacheredis "github.com/snail007/gmc/cache/redis"
	gmcconfig "github.com/snail007/gmc/config/gmc"
	gmcdb "github.com/snail007/gmc/db"
	gmcdbhelper "github.com/snail007/gmc/db/helper"
	gmccontroller "github.com/snail007/gmc/http/controller"
	gmccookie "github.com/snail007/gmc/http/cookie"
	gmcrouter "github.com/snail007/gmc/http/router"
	gmchttpserver "github.com/snail007/gmc/http/server"
	gmcsession "github.com/snail007/gmc/http/session"
	gmctemplate "github.com/snail007/gmc/http/template"
	gmcservice "github.com/snail007/gmc/service"
	"github.com/snail007/gmc/util/maputil"
)

type (
	APP                = gmcapp.GMCApp
	ServiceItem        = gmcapp.ServiceItem
	Cache              = gmccache.Cache
	CacheRedis         = gmccacheredis.RedisCache
	CacheRedisConfig   = gmccacheredis.RedisCacheConfig
	Config             = gmcconfig.GMCConfig
	ResultSet          = gmcdb.ResultSet
	DbCache            = gmcdb.Cache
	Controller         = gmccontroller.Controller
	Cookie             = gmccookie.Cookies
	CookieOptions      = gmccookie.Options
	Router             = gmcrouter.HTTPRouter
	P                  = gmcrouter.Params
	HTTPServer         = gmchttpserver.HTTPServer
	Session            = gmcsession.Session
	SessionStore       = gmcsession.Store
	Template           = gmctemplate.Template
	Service            = gmcservice.Service
	Map                = maputil.Map
	MapStringString    = maputil.MapStringString
	MapStringInterface = maputil.MapStringInterface
	W                  = http.ResponseWriter
	R                  = *http.Request
)

var (
	NewCacheRedisConfig = gmccacheredis.NewRedisCacheConfig
	NewAPP              = gmcapp.New
	NewCacheRedis       = gmccacheredis.New
	NewConfig           = gmcconfig.New
	NewResultSet        = gmcdb.NewResultSet
	NewRouter           = gmcrouter.NewHTTPRouter()
	NewHTTPServer       = gmchttpserver.New
	NewAPIServer        = gmchttpserver.NewAPIServer
	NewTemplate         = gmctemplate.New
	NewCookies          = gmccookie.New

	// Map
	NewMap                = maputil.NewMap
	NewMapStringString    = maputil.NewMapStringString
	NewMapStringInterface = maputil.NewMapStringInterface

	//Database
	InitDB    = gmcdbhelper.RegistGroup
	DBMySQL   = gmcdbhelper.DBMySQL
	DBSQLite3 = gmcdbhelper.DBSQLite3

	//http util
	Stop       = gmchttputil.Stop
	Die        = gmchttputil.Die
	Write      = gmchttputil.Write
	StatusCode = gmchttputil.StatusCode
)
