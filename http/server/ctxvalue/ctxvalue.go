package ctxvalue

import (
	gmcconfig "github.com/snail007/gmc/config/gmc"
	"github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/http/session"
	"github.com/snail007/gmc/http/template"
)

type CtxValue struct {
	Tpl          *template.Template
	SessionStore session.Store
	Router       *router.HTTPRouter
	Config       *gmcconfig.GMCConfig
	AppConfig    *gmcconfig.GMCConfig
}
type (
	ctxValueKey struct{}
)

// CtxValueKey is the request context key under which global HTTP Server object are stored.
var CtxValueKey = ctxValueKey{}
