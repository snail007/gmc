// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gctx

import (
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/util/paginator"
	"net"
	"net/http"
	"strings"
	"time"

	ghttputil "github.com/snail007/gmc/internal/util/http"
)

type Ctx struct {
	response   http.ResponseWriter
	request    *http.Request
	param      gcore.Params
	timeUsed   time.Duration
	localAddr  string
	apiServer  gcore.APIServer
	webServer  gcore.HTTPServer
	app        gcore.App
	logger     gcore.Logger
	config     gcore.Config
	i18n       gcore.I18n
	template   gcore.Template
	remoteAddr string
	conn       net.Conn
}

func (this *Ctx) Conn() net.Conn {
	return this.conn
}

func (this *Ctx) SetConn(conn net.Conn) {
	this.conn = conn
}

func (this *Ctx) RemoteAddr() string {
	return this.remoteAddr
}

func (this *Ctx) SetRemoteAddr(remoteAddr string) {
	this.remoteAddr = remoteAddr
}

func (this *Ctx) Template() gcore.Template {
	return this.template
}

func (this *Ctx) SetTemplate(template gcore.Template) {
	this.template = template
}

func (this *Ctx) I18n() gcore.I18n {
	return this.i18n
}

func (this *Ctx) SetI18n(i18n gcore.I18n) {
	this.i18n = i18n
}

func (this *Ctx) Config() gcore.Config {
	if this.config != nil {
		return this.config
	} else if this.app != nil {
		this.config = this.app.Config()
	} else if this.webServer != nil {
		this.config = this.webServer.Config()
	}
	if this.config == nil {
		this.config = gcore.Providers.Config("")()
	}
	return this.config
}

func (this *Ctx) SetConfig(config gcore.Config) {
	this.config = config
}

func (this *Ctx) Logger() gcore.Logger {
	return this.logger
}

func (this *Ctx) SetLogger(logger gcore.Logger) {
	this.logger = logger
}

func (this *Ctx) App() gcore.App {
	return this.app
}

func (this *Ctx) SetApp(app gcore.App) {
	this.app = app
}

func (this *Ctx) Clone() gcore.Ctx {
	ps := this.param
	if ps == nil {
		ps = gcore.Params{}
	}
	return &Ctx{
		app:        this.app,
		apiServer:  this.apiServer,
		webServer:  this.webServer,
		response:   this.response,
		request:    this.request,
		param:      ps,
		timeUsed:   this.timeUsed,
		localAddr:  this.localAddr,
		remoteAddr: this.remoteAddr,
		logger:     this.logger,
		config:     this.config,
		i18n:       this.i18n,
		template:   this.template,
		conn:       this.conn,
	}
}

func (this *Ctx) CloneWithHTTP(w http.ResponseWriter, r *http.Request, ps ...gcore.Params) gcore.Ctx {
	var ps0 gcore.Params
	if len(ps) > 0 && ps[0] != nil {
		ps0 = ps[0]
	}
	if ps0 == nil {
		ps0 = gcore.Params{}
	}
	c := this.Clone()
	c.SetResponse(w)
	c.SetRequest(r)
	c.SetParam(ps0)
	return c
}

func (this *Ctx) APIServer() gcore.APIServer {
	return this.apiServer
}

func (this *Ctx) SetAPIServer(apiServer gcore.APIServer) {
	this.apiServer = apiServer
}

func (this *Ctx) WebServer() gcore.HTTPServer {
	return this.webServer
}

func (this *Ctx) SetWebServer(webServer gcore.HTTPServer) {
	this.webServer = webServer
}

func (this *Ctx) LocalAddr() string {
	return this.localAddr
}

func (this *Ctx) SetLocalAddr(localAddr string) {
	this.localAddr = localAddr
}

func (this *Ctx) Param() gcore.Params {
	return this.param
}

func (this *Ctx) Request() *http.Request {
	return this.request
}

func (this *Ctx) SetRequest(request *http.Request) {
	this.request = request
}

func (this *Ctx) Response() http.ResponseWriter {
	return this.response
}

func (this *Ctx) SetResponse(response http.ResponseWriter) {
	this.response = response
}

func (this *Ctx) SetParam(param gcore.Params) {
	if param == nil {
		param = gcore.Params{}
	}
	this.param = param
	return
}

// acquires the method cost time, only for middleware2 and middleware3.
func (this *Ctx) TimeUsed() time.Duration {
	return this.timeUsed
}

// sets the method cost time, only for middleware2 and middleware3, do not call this.
func (this *Ctx) SetTimeUsed(t time.Duration) {
	this.timeUsed = t
}

// Write output data to response
func (this *Ctx) Write(data ...interface{}) (n int, err error) {
	return ghttputil.Write(this.Response(), data...)
}

// WriteE outputs data to response and sets http status code 500
func (this *Ctx) WriteE(data ...interface{}) (n int, err error) {
	this.Response().WriteHeader(http.StatusInternalServerError)
	return ghttputil.Write(this.Response(), data...)
}

// WriteHeader sets http code in response
func (this *Ctx) WriteHeader(statusCode int) {
	this.Response().WriteHeader(statusCode)
}

// StatusCode returns http code in response, if not set, default is 200.
func (this *Ctx) StatusCode() int {
	return ghttputil.StatusCode(this.Response())
}

// WriteCount acquires outgoing bytes count by writer
func (this *Ctx) WriteCount() int64 {
	return ghttputil.WriteCount(this.Response())
}

// IsPOST returns true if the request is POST request.
func (this *Ctx) IsPOST() bool {
	return http.MethodPost == this.Request().Method
}

// IsGET returns true if the request is GET request.
func (this *Ctx) IsGET() bool {
	return http.MethodGet == this.Request().Method
}

// IsPUT returns true if the request is PUT request.
func (this *Ctx) IsPUT() bool {
	return http.MethodPut == this.Request().Method
}

// IsDELETE returns true if the request is DELETE request.
func (this *Ctx) IsDELETE() bool {
	return http.MethodDelete == this.Request().Method
}

// IsPATCH returns true if the request is PATCH request.
func (this *Ctx) IsPATCH() bool {
	return http.MethodPatch == this.Request().Method
}

// IsHEAD returns true if the request is HEAD request.
func (this *Ctx) IsHEAD() bool {
	return http.MethodHead == this.Request().Method
}

// IsOPTIONS returns true if the request is OPTIONS request.
func (this *Ctx) IsOPTIONS() bool {
	return http.MethodOptions == this.Request().Method
}

// IsOPTIONS returns true if the request is jquery AJAX request.
func (this *Ctx) IsAJAX() bool {
	return strings.EqualFold(this.Request().Header.Get("X-Requested-With"), "XMLHttpRequest")
}

// Stop will exit controller method or api handle function at once
func (this *Ctx) Stop(msg ...interface{}) {
	ghttputil.Stop(this.Response(), msg...)
}

// ClientIP acquires the real client ip, search in X-Forwarded-For, X-Real-IP, request.RemoteAddr.
func (this *Ctx) ClientIP() (ip string) {
	// X-Forwarded-For
	xForwardedFor := this.Request().Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		proxyIps := strings.Split(strings.Replace(xForwardedFor, " ", "", -1), ",")
		if len(proxyIps) > 0 {
			ip = proxyIps[0]
			return
		}
	}
	// X-Real-IP
	ip = this.Request().Header.Get("X-Real-IP")
	if ip != "" {
		return
	}
	// RemoteAddr
	ip, _, _ = net.SplitHostPort(this.Request().RemoteAddr)
	return
}

// NewPager create a new paginator used for template
func (this *Ctx) NewPager(perPage int, total int64) gcore.Paginator {
	return paginator.NewPaginator(this.Request(), perPage, total, "page")
}

// GET gets the first value associated with the given key in url query.
func (this *Ctx) GET(key string, Default ...string) (val string) {
	val = this.Request().URL.Query().Get(key)
	if val == "" && len(Default) > 0 {
		return Default[0]
	}
	return
}

// POST gets the first value associated with the given key in post body.
func (this *Ctx) POST(key string, Default ...string) (val string) {
	val = this.Request().PostFormValue(key)
	if val == "" && len(Default) > 0 {
		return Default[0]
	}
	return
}

// Redirect redirects to the url, using http location header, and sets http code 302
func (this *Ctx) Redirect(url string) (val string) {
	http.Redirect(this.Response(), this.Request(), url, http.StatusFound)
	ghttputil.JustDie()
	return
}

func NewCtx() *Ctx {
	return &Ctx{
		param: gcore.Params{},
	}
}

func NewCtxFromConfig(c gcore.Config) *Ctx {
	return &Ctx{
		config: c,
		param:  gcore.Params{},
	}
}

func NewCtxFromConfigFile(file string) (ctx *Ctx, err error) {
	c := gcore.Providers.Config("")()
	c.SetConfigFile(file)
	err = c.ReadInConfig()
	if err != nil {
		return
	}
	ctx = &Ctx{
		config: c,
		param:  gcore.Params{},
	}
	return
}
