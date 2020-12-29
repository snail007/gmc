package gcore

import (
	"net"
	"net/http"
	"time"
)

type Ctx interface {
	Template() Template
	SetTemplate(template Template)
	I18n() I18n
	SetI18n(I18n)
	Config() Config
	SetConfig(config Config)
	Logger() Logger
	SetLogger(logger Logger)
	App() App
	SetApp(app App)
	Clone() Ctx
	CloneWithHTTP(w http.ResponseWriter, r *http.Request, ps ...Params) Ctx
	APIServer() APIServer
	SetAPIServer(apiServer APIServer)
	WebServer() HTTPServer
	SetWebServer(webServer HTTPServer)
	LocalAddr() string
	SetLocalAddr(localAddr string)
	Param() Params
	Request() *http.Request
	SetRequest(request *http.Request)
	Response() http.ResponseWriter
	SetResponse(response http.ResponseWriter)
	SetParam(param Params)
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
	Conn() net.Conn
	SetConn(conn net.Conn)
	RemoteAddr() string
	SetRemoteAddr(remoteAddr string)
}
