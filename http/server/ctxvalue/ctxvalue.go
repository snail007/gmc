package ctxvalue

import (
	gmcconfig "github.com/snail007/gmc/config/gmc"
	gmcrouter "github.com/snail007/gmc/http/router"
	gmcsession "github.com/snail007/gmc/http/session"
	gmctemplate "github.com/snail007/gmc/http/template"
)

type CtxValue struct {
	Tpl          *gmctemplate.Template
	SessionStore gmcsession.Store
	Router       *gmcrouter.HTTPRouter
	Config       *gmcconfig.GMCConfig
	AppConfig    *gmcconfig.GMCConfig
}
type (
	ctxValueKey struct{}
)

// CtxValueKey is the request context key under which global HTTP Server object are stored.
var CtxValueKey = ctxValueKey{}
