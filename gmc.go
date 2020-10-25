package gmc

import (
	"github.com/snail007/gmc/core"
	gmcmysql "github.com/snail007/gmc/db/mysql"
	gmcsqlite3 "github.com/snail007/gmc/db/sqlite3"
	gmci18n "github.com/snail007/gmc/i18n"
	gmchttputil "github.com/snail007/gmc/util/http"
	"net/http"

	gmcconfig "github.com/snail007/gmc/config"

	gmcapp "github.com/snail007/gmc/app"
	gmcdb "github.com/snail007/gmc/db"
	gmcerr "github.com/snail007/gmc/error"
	gmccontroller "github.com/snail007/gmc/http/controller"
	gmccookie "github.com/snail007/gmc/http/cookie"
	gmcrouter "github.com/snail007/gmc/http/router"
	gmchttpserver "github.com/snail007/gmc/http/server"
	gmcsession "github.com/snail007/gmc/http/session"
	gmctemplate "github.com/snail007/gmc/http/template"
	gmcmap "github.com/snail007/gmc/util/map"
)

type (
	// app
	APP         = gmcapp.GMCApp
	ServiceItem = gmcapp.ServiceItem
	Service     = gmccore.Service
	Config      = gmcconfig.Config

	// database
	ResultSet = gmcdb.ResultSet
	DBCache   = gmcdb.Cache
	MySQL     = gmcmysql.DB
	SQLite3   = gmcsqlite3.DB

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
	Map = gmcmap.Map

	// Map
	M   = map[string]interface{}
	Mii = map[interface{}]interface{}
	Mss = map[string]string

	// error
	Error = gmcerr.Error
)

var (
	// Errors
	StackE = gmcerr.Stack
	Errorf = gmcerr.Errorf
	WrapE  = gmcerr.Wrap
	StopE  = gmchttputil.StopE

	// i18n
	Tr = gmci18n.Tr
)
