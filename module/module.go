package module

import (
 	gcore "github.com/snail007/gmc/core"
	gcontroller "github.com/snail007/gmc/http/controller"
	gcookie "github.com/snail007/gmc/http/cookie"
	grouter "github.com/snail007/gmc/http/router"
	ghttpserver "github.com/snail007/gmc/http/server"
	gsession "github.com/snail007/gmc/http/session"
	gfilestore "github.com/snail007/gmc/http/session/filestore"
	gmemorystore "github.com/snail007/gmc/http/session/memorystore"
	gredisstore "github.com/snail007/gmc/http/session/redisstore"
	gapp "github.com/snail007/gmc/module/app"
)

// Make sure all modules conform with the correct interface

var _ gcore.SessionStorage = &gmemorystore.MemoryStore{}
var _ gcore.SessionStorage = &gfilestore.FileStore{}
var _ gcore.SessionStorage = &gredisstore.RedisStore{}
var _ gcore.Session = &gsession.Session{}
var _ gcore.Controller =&gcontroller.Controller{}
var _ gcore.Cookies =&gcookie.Cookies{}
var _ gcore.HTTPRouter =&grouter.HTTPRouter{}
var _ gcore.APIServer = &ghttpserver.APIServer{}
var _ gcore.HTTPServer = &ghttpserver.HTTPServer{}
var _ gcore.GMCApp = &gapp.GMCApp{}






