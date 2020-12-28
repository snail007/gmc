package module

import (
	gcore "github.com/snail007/gmc/core"
	gcontroller "github.com/snail007/gmc/http/controller"
	gcookie "github.com/snail007/gmc/http/cookie"
	grouter "github.com/snail007/gmc/http/router"
	ghttpserver "github.com/snail007/gmc/http/server"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
	gview "github.com/snail007/gmc/http/view"
	gapp "github.com/snail007/gmc/module/app"
	gcache "github.com/snail007/gmc/module/cache"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	gdb "github.com/snail007/gmc/module/db"
	gerror "github.com/snail007/gmc/module/error"
)

// Make sure all modules conform with the correct interface

var _ gcore.SessionStorage = &gsession.MemoryStore{}
var _ gcore.SessionStorage = &gsession.FileStore{}
var _ gcore.SessionStorage = &gsession.RedisStore{}
var _ gcore.Session = &gsession.Session{}
var _ gcore.Controller = &gcontroller.Controller{}
var _ gcore.Cookies = &gcookie.Cookies{}
var _ gcore.HTTPRouter = &grouter.HTTPRouter{}
var _ gcore.APIServer = &ghttpserver.APIServer{}
var _ gcore.HTTPServer = &ghttpserver.HTTPServer{}
var _ gcore.Service = &ghttpserver.HTTPServer{}
var _ gcore.Service = &ghttpserver.APIServer{}
var _ gcore.App = &gapp.GMCApp{}
var _ gcore.Template = &gtemplate.Template{}
var _ gcore.View = &gview.View{}
var _ gcore.Cache = &gcache.FileCache{}
var _ gcore.Cache = &gcache.MemCache{}
var _ gcore.Cache = &gcache.RedisCache{}
var _ gcore.Ctx = &gctx.Ctx{}
var _ gcore.Database = &gdb.MySQLDB{}
var _ gcore.Database = &gdb.SQLite3DB{}
var _ gcore.DatabaseGroup = &gdb.MySQLDBGroup{}
var _ gcore.DatabaseGroup = &gdb.SQLite3DBGroup{}
var _ gcore.ActiveRecord = &gdb.MySQLActiveRecord{}
var _ gcore.ActiveRecord = &gdb.SQLite3ActiveRecord{}
var _ gcore.ResultSet = &gdb.ResultSet{}
var _ gcore.Error = &gerror.Error{}
var _ gcore.Database = &gdb.SQLite3DB{}
var _ gcore.Config = &gconfig.Config{}
