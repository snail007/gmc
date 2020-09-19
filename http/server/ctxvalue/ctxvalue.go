package ctxvalue

import (
	"github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/http/session"
	"github.com/snail007/gmc/http/template"
	"github.com/spf13/viper"
)

type CtxValue struct {
	Tpl          *template.Template
	SessionStore session.Store
	Router       *router.HTTPRouter
	Config       *viper.Viper
}
type (
	ctxValueKey struct{}
)

// CtxValueKey is the request context key under which global HTTP Server object are stored.
var CtxValueKey = ctxValueKey{}
