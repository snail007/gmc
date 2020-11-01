package gmc

import (
	gmcapp "github.com/snail007/gmc/app"
	gmcconfig "github.com/snail007/gmc/config"
	gmccore "github.com/snail007/gmc/core"
	gmcdb "github.com/snail007/gmc/db"
	gmcmysql "github.com/snail007/gmc/db/mysql"
	gmcsqlite3 "github.com/snail007/gmc/db/sqlite3"
	gmcerr "github.com/snail007/gmc/error"
	gmccontroller "github.com/snail007/gmc/http/controller"
	gmccookie "github.com/snail007/gmc/http/cookie"
	gmcrouter "github.com/snail007/gmc/http/router"
	gmchttpserver "github.com/snail007/gmc/http/server"
	gmcsession "github.com/snail007/gmc/http/session"
	gmctemplate "github.com/snail007/gmc/http/template"
	gmci18n "github.com/snail007/gmc/i18n"
	gmchttputil "github.com/snail007/gmc/util/http"
	gmcmap "github.com/snail007/gmc/util/map"
	"net/http"
)

type (
	// Alias of type gmcapp.GMCApp
	APP = gmcapp.GMCApp
	// Alias of type gmcapp.ServiceItem
	ServiceItem = gmcapp.ServiceItem
	// Alias of type gmccore.Service
	Service = gmccore.Service
	// Alias of type gmcconfig.Config
	Config = gmcconfig.Config
	// Alias of type gmcdb.ResultSet
	ResultSet = gmcdb.ResultSet
	// Alias of type gmcdb.Cache
	DBCache = gmcdb.Cache
	// Alias of type gmcmysql.DB
	MySQL = gmcmysql.DB
	// Alias of type gmcsqlite3.DB
	SQLite3 = gmcsqlite3.DB
	// Alias of type gmccontroller.Controller
	Controller = gmccontroller.Controller
	// Alias of type gmccookie.Cookies
	Cookie = gmccookie.Cookies
	// Alias of type gmccookie.Options
	CookieOptions = gmccookie.Options
	// Alias of type gmcrouter.HTTPRouter
	Router = gmcrouter.HTTPRouter
	// Alias of type gmcrouter.Handle
	Handle = gmcrouter.Handle
	// Alias of type gmcrouter.Params
	P = gmcrouter.Params
	// Alias of type gmchttpserver.HTTPServer
	HTTPServer = gmchttpserver.HTTPServer
	// Alias of type gmchttpserver.APIServer
	APIServer = gmchttpserver.APIServer
	// Alias of type gmcsession.Session
	Session = gmcsession.Session
	// Alias of type gmcsession.Store
	SessionStore = gmcsession.Store
	// Alias of type gmctemplate.Template
	Template = gmctemplate.Template
	// Alias of type http.ResponseWriter
	W = http.ResponseWriter
	// Alias of type *http.Request
	R = *http.Request
	// Alias of type *gmcrouter.Ctx
	C = *gmcrouter.Ctx
	// Alias of type gmcmap.Map
	Map = gmcmap.Map
	// Alias of type map[string]interface{}
	M = map[string]interface{}
	// Alias of type map[interface{}]interface{}
	Mii = map[interface{}]interface{}
	// Alias of type map[string]string
	Mss = map[string]string
	// Alias of type gmcerr.Error
	Error = gmcerr.Error
)

var (

	// Shortcut of gmcerr.Stack.
	// Acquires the full stack, returns string.
	StackE = gmcerr.Stack

	// Shortcut of gmcerr.Errorf
	// Create a gmcerr.Error with format string.
	Errorf = gmcerr.Errorf

	// Shortcut of gmcerr.Wrap
	// Wrap an error or string to a gmc.Error.
	WrapE = gmcerr.Wrap

	// Shortcut of gmchttputil.StopE
	// StopE will exit controller method if error is not nil.
	// First argument is an error.
	// Secondary argument is fail function, it be called if error is not nil.
	// Third argument is success function, it be called if error is nil.
	StopE = gmchttputil.StopE

	// Shortcut of gmci18n.Tr
	// This only worked after gmci18n.Init called.
	Tr = gmci18n.Tr
)
