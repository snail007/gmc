package gmc

import (
	"net/http"

	gmcapp "github.com/snail007/gmc/app"
	gmccachehelper "github.com/snail007/gmc/cache/helper"
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
	NewAPP              = gmcapp.New
	NewConfig           = gmcconfig.New
	NewRouter           = gmcrouter.NewHTTPRouter()
	NewHTTPServer       = gmchttpserver.New
	NewAPIServer        = gmchttpserver.NewAPIServer

	// Map
	NewMap                = maputil.NewMap

	//Database
	InitDB    = gmcdbhelper.Init
	DBMySQL   = gmcdbhelper.DBMySQL
	DBSQLite = gmcdbhelper.DBSQLite3

	//Cache
	InitCache = gmccachehelper.Init
	Reids     = gmccachehelper.Redis
	Cache     = gmccachehelper.Cache

	// Errors
	StackE = gmcerr.Stack
	Errorf = gmcerr.Errorf
	WrapE  = gmcerr.Wrap
)
