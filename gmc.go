package gmc

import (
	gconfig "github.com/snail007/gmc/config"
	gcore "github.com/snail007/gmc/core"
	gapp "github.com/snail007/gmc/gmc/app"
	gdb "github.com/snail007/gmc/gmc/db"
	gerr "github.com/snail007/gmc/gmc/error"
	gi18n "github.com/snail007/gmc/gmc/i18n"
	gcontroller "github.com/snail007/gmc/http/controller"
	gcookie "github.com/snail007/gmc/http/cookie"
	grouter "github.com/snail007/gmc/http/router"
	ghttpserver "github.com/snail007/gmc/http/server"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
	ghttputil "github.com/snail007/gmc/util/http"
	gmap "github.com/snail007/gmc/util/map"
	"net/http"
)

type (
	// Alias of type gapp.GMCApp
	APP = gapp.GMCApp
	// Alias of type gapp.ServiceItem
	ServiceItem = gcore.ServiceItem
	// Alias of type gcore.Service
	Service = gcore.Service
	// Alias of type gconfig.Config
	Config = gconfig.Config
	// Alias of type gdb.ResultSet
	ResultSet = gdb.ResultSet
	// Alias of type gdb.DBCache
	DBCache = gcore.DBCache
	// Alias of type gdb.DatabaseGroup
	DBGroup = gcore.DatabaseGroup
	// Alias of type gdb.Database
	Database = gcore.Database
	// Alias of type gmysql.SQLite3DB
	MySQL = gdb.MySQLDB
	// Alias of type gsqlite3.SQLite3DB
	SQLite3 = gdb.MySQLDB
	// Alias of type gmysql.SQLite3DB
	MySQLG = gdb.MySQLDBGroup
	// Alias of type gsqlite3.SQLite3DB
	SQLite3G = gdb.MySQLDBGroup
	// Alias of type gcontroller.Controller
	Controller = gcontroller.Controller
	// Alias of type gcookie.Cookies
	Cookie = gcookie.Cookies
	// Alias of type gcookie.Options
	CookieOptions = gcore.CookieOptions
	// Alias of type grouter.HTTPRouter
	Router = gcore.HTTPRouter
	// Alias of type grouter.Handle
	Handle = gcore.Handle
	// Alias of type grouter.Params
	P = gcore.Params
	// Alias of type ghttpserver.HTTPServer
	HTTPServer = ghttpserver.HTTPServer
	// Alias of type ghttpserver.APIServer
	APIServer = ghttpserver.APIServer
	// Alias of type gsession.Session
	Session = gsession.Session
	// Alias of type gsession.SessionStorage
	SessionStore = gcore.SessionStorage
	// Alias of type gtemplate.Template
	Template = gtemplate.Template
	// Alias of type http.ResponseWriter
	W = http.ResponseWriter
	// Alias of type *http.Request
	R = *http.Request
	// Alias of type gcore.Ctx
	C = gcore.Ctx
	// Alias of type gmap.Map
	Map = gmap.Map
	// Alias of type map[string]interface{}
	M = map[string]interface{}
	// Alias of type map[interface{}]interface{}
	Mii = map[interface{}]interface{}
	// Alias of type map[string]string
	Mss = map[string]string
	// Alias of type gerr.Error
	Error = gerr.Error
)

var (

	// Alias of gerr.Stack.
	// Acquires the full stack, returns string.
	StackE = gerr.Stack

	// Alias of gerr.Errorf
	// Create a gerr.Error with format string.
	Errorf = gerr.Errorf

	// Alias of gerr.Wrap
	// Wrap an error or string to a gmc.Error.
	WrapE = gerr.Wrap

	// Alias of ghttputil.StopE
	// StopE will exit controller method if error is not nil.
	// First argument is an error.
	// Secondary argument is fail function, it be called if error is not nil.
	// Third argument is success function, it be called if error is nil.
	StopE = ghttputil.StopE

	// Alias of gi18n.Tr
	// This only worked after gi18n.Init called.
	Tr = gi18n.Tr
	// Alias of gi18n.Tr grouter.ParamsFromContext
	ParamsFromContext = grouter.ParamsFromContext

	// Alias of gdb.Model
	Table = gdb.Table
)
