// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcore

import (
	"mime/multipart"
	"net"
	"net/http"
	"time"
)

// Ctx gmc context layer in gmc web.
type Ctx interface {
	Conn() net.Conn
	SetConn(conn net.Conn)
	RemoteAddr() string
	SetRemoteAddr(remoteAddr string)
	Template() Template
	SetTemplate(template Template)
	I18n() I18n
	SetI18n(i18n I18n)
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
	GetParam(key string) string
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
	Status(code int)
	WriteCount() int64
	IsPOST() bool
	IsGET() bool
	IsPUT() bool
	IsDELETE() bool
	IsPATCH() bool
	IsHEAD() bool
	IsOPTIONS() bool
	IsAJAX() bool
	IsWebsocket() bool
	Stop(msg ...interface{})
	StopJSON(code int, msg interface{})
	ClientIP() (ip string)
	NewPager(perPage int, total int64) Paginator
	GET(key string, Default ...string) (val string)
	GETArray(key string, Default ...string) (val []string)
	GETData() (data map[string]string)
	POST(key string, Default ...string) (val string)
	POSTArray(key string, Default ...string) (val []string)
	POSTData() (data map[string]string)
	GetPost(key string, Default ...string) (val string)
	Host() (host string)
	JSON(code int, data interface{}) (err error)
	PrettyJSON(code int, data interface{}) (err error)
	JSONP(code int, data interface{}) (err error)
	Redirect(url string) (val string)
	SetHeader(key, value string)
	Header(key string) string
	RequestBody() ([]byte, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)
	Cookie(name string) string
	WriteFile(filepath string)
	WriteFileFromFS(filepath string, fs http.FileSystem)
	WriteFileAttachment(filepath, filename string)
	FullPath() string
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, exists bool)
	MustGet(key interface{}) (v interface{})
	FormFile(name string, maxMultipartMemory int64) (*multipart.FileHeader, error)
	MultipartForm(maxMultipartMemory int64) (*multipart.Form, error)
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
	ControllerMethod() string
	SetControllerMethod(controllerMethod string)
	IsTLSRequest() bool
}

const ctxKeyInResponseWriter = "CtxKeyInResponseWriter"

func SetCtx(w http.ResponseWriter, c Ctx) http.ResponseWriter {
	if v, ok := w.(ResponseWriter); ok {
		v.SetData(ctxKeyInResponseWriter, c)
	}
	return w
}

func GetCtx(w http.ResponseWriter) Ctx {
	if v, ok := w.(ResponseWriter); ok {
		ctx := v.Data(ctxKeyInResponseWriter)
		if ctx == nil {
			return nil
		}
		return ctx.(Ctx)
	}
	return nil
}
