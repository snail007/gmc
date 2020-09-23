package gmc

import (
	gmcapp "github.com/snail007/gmc/app"
	gmccache "github.com/snail007/gmc/cache"
	gmccacheredis "github.com/snail007/gmc/cache/redis"
	gmcconfig "github.com/snail007/gmc/config/gmc"
	gmcdb "github.com/snail007/gmc/db"
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
	RouterParams       = gmcrouter.Params
	HTTPServer         = gmchttpserver.HTTPServer
	Session            = gmcsession.Session
	SessionStore       = gmcsession.Store
	Template           = gmctemplate.Template
	Service            = gmcservice.Service
	Map                = maputil.Map
	MapStringString    = maputil.MapStringString
	MapStringInterface = maputil.MapStringInterface
)

var (
	NewCacheRedisConfig   = gmccacheredis.NewRedisCacheConfig
	NewAPP                = gmcapp.New
	NewCacheRedis         = gmccacheredis.New
	NewConfig             = gmcconfig.New
	NewResultSet          = gmcdb.NewResultSet
	NewRouter             = gmcrouter.NewHTTPRouter()
	NewHTTPServer         = gmchttpserver.New
	NewTemplate           = gmctemplate.New
	NewMap                = maputil.NewMap
	NewMapStringString    = maputil.NewMapStringString
	NewMapStringInterface = maputil.NewMapStringInterface
)
