package ctxvalue

import (
	gconfig "github.com/snail007/gmc/config"
	gcore "github.com/snail007/gmc/core"
	grouter "github.com/snail007/gmc/http/router"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
)

type CtxValue struct {
	Tpl          *gtemplate.Template
	SessionStore gsession.Store
	Router       *grouter.HTTPRouter
	Config       *gconfig.Config
	AppConfig    *gconfig.Config
	Logger       gcore.Logger
}

type (
	ctxValueKey struct{}
)

// CtxValueKey is the request context key under which global HTTP Server object are stored.
var CtxValueKey = ctxValueKey{}
