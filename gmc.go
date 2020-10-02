package gmc

import (
	"net/http"

	gmchttputil "github.com/snail007/gmc/util/httputil"

	gmcapp "github.com/snail007/gmc/app"
	gmccachehelper "github.com/snail007/gmc/cache/helper"
	gmccacheredis "github.com/snail007/gmc/cache/redis"
	gmcconfig "github.com/snail007/gmc/config/gmc"
	gmcdb "github.com/snail007/gmc/db"
	gmcdbhelper "github.com/snail007/gmc/db/helper"
	gmcerr "github.com/snail007/gmc/error"
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
	// app
	APP         = gmcapp.GMCApp
	ServiceItem = gmcapp.ServiceItem
	Service     = gmcservice.Service
	Config      = gmcconfig.GMCConfig
	// database
	ResultSet = gmcdb.ResultSet
	DbCache   = gmcdb.Cache
	// http server
	Controller    = gmccontroller.Controller
	Cookie        = gmccookie.Cookies
	CookieOptions = gmccookie.Options
	Router        = gmcrouter.HTTPRouter
	HTTPServer    = gmchttpserver.HTTPServer
	APIServer     = gmchttpserver.APIServer
	Session       = gmcsession.Session
	SessionStore  = gmcsession.Store
	Template      = gmctemplate.Template
	W             = http.ResponseWriter
	R             = *http.Request
	C             = *gmcrouter.Ctx
	// data
	Map                = maputil.Map
	MapStringString    = maputil.MapStringString
	MapStringInterface = maputil.MapStringInterface
	// error
	Error = gmcerr.Error
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
	InitDB    = gmcdbhelper.Init
	DBMySQL   = gmcdbhelper.DBMySQL
	DBSQLite3 = gmcdbhelper.DBSQLite3

	//Cache
	InitCache = gmccachehelper.Init
	Reids     = gmccachehelper.Redis
	Cache     = gmccachehelper.Cache

	//http util
	Stop       = gmchttputil.Stop
	Die        = gmchttputil.Die
	Write      = gmchttputil.Write
	StatusCode = gmchttputil.StatusCode

	// Errors
	StackE = gmcerr.Stack
	Errorf = gmcerr.Errorf
	WrapE  = gmcerr.Wrap
)
