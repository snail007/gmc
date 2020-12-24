package gcore

import (
	gtemplate "github.com/snail007/gmc/http/template"
	gconfig "github.com/snail007/gmc/util/config"
	"io"
	"net"
	"net/http"
	"time"
)

type (
	// Handle is a function that can be registered to a route to handle HTTP
	// requests. Like http.HandlerFunc, but has a third parameter for the values of
	// wildcards (path variables).
	Handle func(http.ResponseWriter, *http.Request, Params)
	// Alias of api middleware
	MiddlewareAPI func(ctx Ctx, api APIServer) (isStop bool)
	// Alias of web middleware
	MiddlewareWeb func(ctx Ctx, server HTTPServer) (isStop bool)
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}

// MatchedRoutePathParam is the Param name under which the path of the matched
// route is stored, if Router.SaveMatchedRoutePath is set.
var MatchedRoutePathParam = "$matchedRoutePath"

// MatchedRoutePath retrieves the path of the matched route.
// Router.SaveMatchedRoutePath must have been enabled when the respective
// handler was added, otherwise this function always returns an empty string.
func (ps Params) MatchedRoutePath() string {
	return ps.ByName(MatchedRoutePathParam)
}

type Ctx interface {
	LocalAddr() string
	SetLocalAddr(localAddr string)
	Param() Params
	Request() *http.Request
	Response() http.ResponseWriter
	TimeUsed() time.Duration
	SetTimeUsed(t time.Duration)
	Write(data ...interface{}) (n int, err error)
	WriteE(data ...interface{}) (n int, err error)
	WriteHeader(statusCode int)
	StatusCode() int
	WriteCount() int64
	IsPOST() bool
	IsGET() bool
	IsPUT() bool
	IsDELETE() bool
	IsPATCH() bool
	IsHEAD() bool
	IsOPTIONS() bool
	IsAJAX() bool
	Stop(msg ...interface{})
	ClientIP() (ip string)
	NewPager(perPage int, total int64) Paginator
	GET(key string, Default ...string) (val string)
	POST(key string, Default ...string) (val string)
	Redirect(url string) (val string)
}

type CtxValue struct {
	Tpl          Template
	SessionStore SessionStorage
	Router       HTTPRouter
	Config       *gconfig.Config
	AppConfig    *gconfig.Config
	Logger       Logger
}

type (
	ctxValueKey struct{}
)

// CtxValueKey is the request context key under which global HTTP Server object are stored.
var CtxValueKey = ctxValueKey{}

// Options is used to setting cookie.
type CookieOptions struct {
	MaxAge   int    // optional
	Path     string // optional, default to "/"
	Domain   string // optional
	Secure   bool   // optional
	HTTPOnly bool   // optional, default to `true``
}

var DefaultCookieOptions = &CookieOptions{
	MaxAge:   0,
	Path:     "/",
	Domain:   "",
	Secure:   false,
	HTTPOnly: true,
}

type Cookies interface {
	Get(name string, signed ...bool) (value string, err error)
	Set(name, val string, options ...*CookieOptions)
	Remove(name string, options ...*CookieOptions)
}

type SessionStorage interface {
	Load(sessionID string) (session Session, isExits bool)
	Save(session Session) (err error)
	Delete(sessionID string) (err error)
}

type Session interface {
	Set(k interface{}, v interface{})
	Get(k interface{}) (value interface{})
	Delete(k interface{}) (err error)
	Destroy() (err error)
	Values() (data map[interface{}]interface{})
	IsDestroy() bool
	SessionID() (sessionID string)
	TouchTime() (time int64)
	Touch()
	Serialize() (str string, err error)
	Unserialize(data string) (err error)
}

type Template interface {
	Delims(left, right string) *gtemplate.Template
	Funcs(funcMap map[string]interface{}) *gtemplate.Template
	String() string
	Extension(ext string) *gtemplate.Template
	Execute(name string, data interface{}) (output []byte, err error)
	Parse() (err error)
}

type View interface {
	Err() error
	Set(key string, val interface{}) View
	SetMap(d map[string]interface{}) View
	Render(tpl string, data ...map[string]interface{}) View
	RenderR(tpl string, data ...map[string]interface{}) (d []byte)
	Layout(l string) View
	Stop()
	OnRenderOnce(f func()) View
	SetLayoutDir(layoutDir string)
}
type HTTPRouter interface {
	Group(ns string) HTTPRouter
	Namespace() string
	Ext(ext string)
	Controller(urlPath string, obj Controller, ext ...string)
	PrintRouteTable(w io.Writer)
	RouteTable() (table map[string][]string)
	ControllerMethod(urlPath string, obj Controller, method string)
	Handle(method, path string, handle Handle)
	HandleAny(path string, handle Handle)
	Handler(method, path string, handler http.Handler)
	HandlerAny(path string, handler http.Handler)
	HandlerFunc(method, path string, handler http.HandlerFunc)
	HandlerFuncAny(path string, handler http.HandlerFunc)
	GET(path string, handle Handle)
	HEAD(path string, handle Handle)
	OPTIONS(path string, handle Handle)
	POST(path string, handle Handle)
	PUT(path string, handle Handle)
	PATCH(path string, handle Handle)
	DELETE(path string, handle Handle)
	ServeFiles(path string, root http.FileSystem)
	Lookup(method, path string) (Handle, Params, bool)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
type APIServer interface {
	Server() *http.Server
	Router() HTTPRouter
	SetTLSFile(certFile, keyFile string)
	SetLogger(l Logger)
	Logger() Logger
	AddMiddleware0(m func(ctx Ctx, server APIServer) (isStop bool))
	AddMiddleware1(m func(ctx Ctx, server APIServer) (isStop bool))
	AddMiddleware2(m func(ctx Ctx, server APIServer) (isStop bool))
	AddMiddleware3(m func(ctx Ctx, server APIServer) (isStop bool))
	Handle404(handle func(ctx Ctx))
	Handle500(handle func(ctx Ctx, err interface{}))
	ShowErrorStack(isShow bool)
	Ext(ext string)
	API(path string, handle func(ctx Ctx), ext ...string)
	Group(path string) APIServer
	PrintRouteTable(w io.Writer)
	ActiveConnCount() int64
	SetLog(l Logger)
	Listeners() []net.Listener
	Listener() net.Listener
}

type HTTPServer interface {
	SetHandler40x(fn func(ctx Ctx, tpl Template))
	SetHandler50x(fn func(ctx Ctx, tpl Template, err interface{}))
	AddFuncMap(f map[string]interface{})
	SetConfig(c *gconfig.Config)
	Config() *gconfig.Config
	ActiveConnCount() int64
	Listener() net.Listener
	Listeners() []net.Listener
	Server() *http.Server
	SetLogger(l Logger)
	Logger() Logger
	SetRouter(r HTTPRouter)
	Router() HTTPRouter
	SetTpl(t Template)
	Tpl() Template
	SetSessionStore(st SessionStorage)
	SessionStore() SessionStorage
	AddMiddleware0(m func(ctx Ctx, server HTTPServer) (isStop bool))
	AddMiddleware1(m func(ctx Ctx, server HTTPServer) (isStop bool))
	AddMiddleware2(m func(ctx Ctx, server HTTPServer) (isStop bool))
	AddMiddleware3(m func(ctx Ctx, server HTTPServer) (isStop bool))
	PrintRouteTable(w io.Writer)
	SetLog(l Logger)
}

type paramsKey struct{}

// ParamsKey is the request context key under which URL params are stored.
var ParamsKey = paramsKey{}

type Controller interface {
	MethodCallPre(w http.ResponseWriter, r *http.Request, ps Params)
	MethodCallPost()
	Tr(key string, defaultText ...string) string
	Die(msg ...interface{})
	Stop(msg ...interface{})
	StopE(err interface{}, fn ...func())
	SessionStart() (err error)
	SessionDestroy() (err error)
	Write(data ...interface{}) (n int, err error)
	WriteE(data ...interface{}) (n int, err error)
}
