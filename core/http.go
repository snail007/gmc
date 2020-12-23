package gcore

import (
	"io"
	"net/http"
	"time"
)

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (path variables).
type Handle func(http.ResponseWriter, *http.Request, Params)

type Param interface {
	Key() string
	Value() string
}

type Params interface {
	ByName(name string) string
	SetParam(p Param)
	AppendParam(p Param)
	First() Params
	Len() int
	Truncate(end int)
	Set(int, Param)
}

type HTTPRouter interface {
	Group(ns string) HTTPRouter
	Namespace() string
	Ext(ext string) HTTPRouter
	Controller(urlPath string, obj interface{}, ext ...string)
	PrintRouteTable(w io.Writer)
	RouteTable() (table map[string][]string)
	ControllerMethod(urlPath string, obj interface{}, method string)
	Handle(method, path string, handle Handle)
	HandleAny(path string, handle Handle)
	Handler(method, path string, handler http.Handler)
	HandlerAny(path string, handler http.Handler)
	HandlerFunc(method, path string, handler http.HandlerFunc)
	HandlerFuncAny(path string, handler http.HandlerFunc)
}
type Ctx interface {
	LocalAddr() string
	SetLocalAddr(localAddr string)
	Param() Params
	Request() *http.Request
	SetRequest(request *http.Request)
	Response() http.ResponseWriter
	SetResponse(response http.ResponseWriter)
	SetParam(param Params) Ctx
	TimeUsed() time.Duration
	SetTimeUsed(t time.Duration) Ctx
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
	Set(name, val string, options ...*CookieOptions) Cookies
	Remove(name string, options ...*CookieOptions)
}
