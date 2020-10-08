package gmc

import (
	gmcmysql "github.com/snail007/gmc/db/mysql"
	gmcsqlite3 "github.com/snail007/gmc/db/sqlite3"
	gmci18n "github.com/snail007/gmc/i18n"
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
	gmcservice "github.com/snail007/gmc/service"
	"github.com/snail007/gmc/util/maputil"
)

type (
	// app
	APP         = gmcapp.GMCApp
	ServiceItem = gmcapp.ServiceItem
	Service     = gmcservice.Service
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
	Map = maputil.Map

	// error
	Error = gmcerr.Error
)

var (
	// Errors
	StackE = gmcerr.Stack
	Errorf = gmcerr.Errorf
	WrapE  = gmcerr.Wrap

	// i18n
	Tr = gmci18n.Tr
)
