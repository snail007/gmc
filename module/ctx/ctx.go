// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gctx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	gcore "github.com/snail007/gmc/core"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/snail007/gmc/util/paginator"

	ghttputil "github.com/snail007/gmc/internal/util/http"
)

type Ctx struct {
	response         http.ResponseWriter
	request          *http.Request
	param            gcore.Params
	timeUsed         time.Duration
	localAddr        string
	apiServer        gcore.APIServer
	webServer        gcore.HTTPServer
	app              gcore.App
	logger           gcore.Logger
	config           gcore.Config
	i18n             gcore.I18n
	template         gcore.Template
	remoteAddr       string
	conn             net.Conn
	metadata         *gmap.Map
	controllerMethod string
}

func (this *Ctx) IsTLSRequest() bool {
	return this.request.TLS != nil
}

func (this *Ctx) ControllerMethod() string {
	return this.controllerMethod
}

func (this *Ctx) SetControllerMethod(controllerMethod string) {
	this.controllerMethod = controllerMethod
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
		this.config = gcore.ProviderConfig()()
	}
	return this.config
}

func (this *Ctx) SetConfig(config gcore.Config) {
	this.config = config
}

func (this *Ctx) Logger() gcore.Logger {
	if this.logger != nil {
		return this.logger
	} else if this.app != nil && this.app.Logger() != nil {
		this.logger = this.app.Logger()
	} else if this.webServer != nil && this.webServer.Logger() != nil {
		this.logger = this.webServer.Logger()
	} else if this.apiServer != nil && this.apiServer.Logger() != nil {
		this.logger = this.apiServer.Logger()
	}
	if this.logger == nil {
		this.logger = gcore.ProviderLogger()(this, "")
	}
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
	c := &Ctx{
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
		metadata:   this.metadata.Clone(),
	}
	paramCopy := make([]gcore.Param, len(this.param))
	copy(paramCopy, this.param)
	c.param = paramCopy
	return c
}

func (this *Ctx) CloneWithHTTP(w http.ResponseWriter, r *http.Request, ps ...gcore.Params) gcore.Ctx {
	var ps0 gcore.Params
	if len(ps) > 0 && ps[0] != nil {
		ps0 = ps[0]
	}
	if ps0 == nil {
		ps0 = gcore.Params{}
	}
	c := this.Clone().(*Ctx)
	c.response = w
	c.request = r
	c.param = ps0
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

// GetParam returns the value of the URL param.
// It is a shortcut for c.Param().ByName(key)
//     router.GET("/user/:id", func(c *gcore.Ctx) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func (this *Ctx) GetParam(key string) string {
	return this.param.ByName(key)
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

// TimeUsed acquires the method cost time, only for middleware2 and middleware3.
func (this *Ctx) TimeUsed() time.Duration {
	return this.timeUsed
}

// SetTimeUsed sets the method cost time, only for middleware2 and middleware3, do not call this.
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
	this.response.WriteHeader(statusCode)
}

// StatusCode returns http code in response, if not set, default is 200.
func (this *Ctx) StatusCode() int {
	return ghttputil.StatusCode(this.Response())
}

// Status sets the HTTP response code.
func (this *Ctx) Status(code int) {
	this.response.WriteHeader(code)
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

// IsAJAX returns true if the request is jquery AJAX request.
func (this *Ctx) IsAJAX() bool {
	return strings.EqualFold(this.Header("X-Requested-With"), "XMLHttpRequest")
}

// IsWebsocket returns true if the request headers indicate that a websocket
// handshake is being initiated by the client.
func (this *Ctx) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(this.Header("Connection")), "upgrade") &&
		strings.EqualFold(this.Header("Upgrade"), "websocket") {
		return true
	}
	return false
}

// Stop will exit controller method or api handle function at once.
func (this *Ctx) Stop(msg ...interface{}) {
	ghttputil.Stop(this.Response(), msg...)
}

// StopJSON will exit controller method or api handle function at once.
// And  writes the status code and a JSON body.
// It also sets the Content-Type as "application/json".
func (this *Ctx) StopJSON(code int, msg interface{}) {
	this.JSON(code, msg)
	this.Stop()
}

var (
	ipOnce        = sync.Once{}
	ipNetwork     []*net.IPNet
	ipFetchMaxTry = 60
)

type CloudFlareIPS struct {
	Result struct {
		Ipv4Cidrs []string `json:"ipv4_cidrs"`
		Ipv6Cidrs []string `json:"ipv6_cidrs"`
		Etag      string   `json:"etag"`
	} `json:"result"`
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
}

// ClientIP acquires the real client ip, search in X-Forwarded-For, X-Real-IP, request.RemoteAddr.
func (this *Ctx) ClientIP() (ip string) {
	clientIP, _, _ := net.SplitHostPort(this.Request().RemoteAddr)
	frontType := this.Config().GetString("frontend.type")
	headerKey := this.Config().GetString("frontend.header")
	ipOnce.Do(func() {
		var cloudflare = func() (err error) {
			client := http.Client{
				Timeout: time.Second * 10,
			}
			resp, err := client.Get("https://api.cloudflare.com/client/v4/ips")
			if err != nil {
				err = fmt.Errorf("fetch cloudflare ips fail, error:%s ", err)
				return
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				err = fmt.Errorf("read cloudflare ips fail, error:%s ", err)
				return
			}
			var cloudFlareIPS CloudFlareIPS
			err = json.Unmarshal(b, &cloudFlareIPS)
			if err != nil {
				err = fmt.Errorf("parse cloudflare ips fail, content:%s, error:%s", string(b), err)
				return
			}
			if !cloudFlareIPS.Success {
				err = fmt.Errorf("parse cloudflare ips fail, content:%s, error:%s", string(b), err)
				return
			}
			for _, v := range append(cloudFlareIPS.Result.Ipv4Cidrs, cloudFlareIPS.Result.Ipv6Cidrs...) {
				_, cloudFlareIPsNetwork, _ := net.ParseCIDR(v)
				if cloudFlareIPsNetwork != nil {
					ipNetwork = append(ipNetwork, cloudFlareIPsNetwork)
				}
			}
			return
		}
		switch frontType {
		case "cloudflare":
			err := cloudflare()
			if err != nil {
				this.Logger().Warn(err)
				go func() {
					var tryCnt = 0
					for {
						err := cloudflare()
						if err != nil {
							tryCnt++
							if tryCnt >= ipFetchMaxTry {
								this.Logger().Warn("fetch cloudflare ip range fail, max try %d reached", ipFetchMaxTry)
								return
							}
							time.Sleep(time.Second * 30)
						} else {
							return
						}
					}
				}()
			}
		case "proxy":
			ipsProxy := this.Config().GetStringSlice("frontend.ips")
			if len(ipsProxy) > 0 {
				for _, v := range ipsProxy {
					if !strings.Contains(v, "/") {
						v += "/32"
					}
					_, cloudFlareIPsNetwork, _ := net.ParseCIDR(v)
					if cloudFlareIPsNetwork != nil {
						ipNetwork = append(ipNetwork, cloudFlareIPsNetwork)
					}
				}
			}
		}
	})
	check := false
	if len(ipNetwork) > 0 {
		clientIPObj := net.ParseIP(clientIP)
		for _, ipn := range ipNetwork {
			if ipn.Contains(clientIPObj) {
				check = true
				break
			}
		}
	}
	if check {
		switch frontType {
		case "cloudflare":
			if headerKey == "" {
				headerKey = "CF-Connecting-IP"
			}
			ip = this.Header(headerKey)
		case "proxy":
			if headerKey == "" {
				headerKey = "X-Forwarded-For"
			}
			if strings.ToLower(headerKey) == strings.ToLower("X-Forwarded-For") {
				xForwardedFor := this.Header(headerKey)
				if xForwardedFor != "" {
					xForwardedFor = strings.Replace(xForwardedFor, ":", " ", -1)
					xForwardedFor = strings.Replace(xForwardedFor, ",", " ", -1)
					proxyIps := strings.Fields(xForwardedFor)
					if len(proxyIps) > 0 {
						ip = proxyIps[0]
					}
				}
			} else {
				ip = this.Header(headerKey)
			}
		}
	}
	if ip == "" {
		// RemoteAddr
		ip = clientIP
	}
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

// GETArray returns a slice of strings for a given query key.
func (this *Ctx) GETArray(key string, Default ...string) (val []string) {
	this.request.ParseForm()
	val = this.Request().URL.Query()[key]
	if len(val) == 0 && len(Default) > 0 {
		return Default
	}
	return
}

// GETData gets full k,v query data from request url.
func (this *Ctx) GETData() (data map[string]string) {
	data = make(map[string]string)
	for k, v := range this.Request().URL.Query() {
		if len(v) > 0 {
			data[k] = v[0]
		}
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

// POSTArray returns a slice of strings for a given form key.
// The length of the slice depends on the number of params with the given key.
func (this *Ctx) POSTArray(key string, Default ...string) (val []string) {
	this.Request().ParseForm()
	val = this.Request().PostForm[key]
	if len(val) == 0 && len(Default) > 0 {
		return Default
	}
	return
}

// POSTData gets full k,v query data from request FORM.
func (this *Ctx) POSTData() (data map[string]string) {
	data = make(map[string]string)
	this.Request().ParseForm()
	for k, v := range this.Request().PostForm {
		if len(v) > 0 {
			data[k] = v[0]
		}
	}
	return
}

// GetPost gets the first value associated with the given key in GET or POST.
func (this *Ctx) GetPost(key string, Default ...string) (val string) {
	val = this.GET(key, "")
	if val == "" {
		return this.POST(key, Default...)
	}
	return
}

// Host returns the HOST in requested HTTP HEADER.
func (this *Ctx) Host() (host string) {
	return this.request.Host
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (this *Ctx) JSON(code int, data interface{}) (err error) {
	this.Status(code)
	this.SetHeader("Content-Type", "application/json")
	_, err = this.Write(data)
	return
}

// PrettyJSON serializes the given struct as pretty JSON (indented) into the response body.
// It also sets the Content-Type as "application/json". WARNING: we recommend to use this only for
// development purposes since printing pretty JSON is more CPU and bandwidth consuming.
// Use Ctx.JSON() instead.
func (this *Ctx) PrettyJSON(code int, data interface{}) (err error) {
	this.Status(code)
	this.SetHeader("Content-Type", "application/json")
	_, err = ghttputil.WritePretty(this.Response(), data)
	return
}

// JSONP serializes the given struct as JSON into the response body.
// It sets the Content-Type as "application/javascript".
func (this *Ctx) JSONP(code int, data interface{}) (err error) {
	callback := this.GET("callback", "_jsonp")
	this.response.WriteHeader(code)
	this.response.Header().Set("Content-Type", "application/javascript")
	var buf = &bytes.Buffer{}
	buf.Write([]byte(callback + "("))
	_, err = ghttputil.Write(buf, data)
	if err != nil {
		return
	}
	buf.Write([]byte(")"))
	_, err = this.Write(buf.Bytes())
	return
}

// Redirect redirects to the url, using http location header, and sets http code 302
func (this *Ctx) Redirect(url string) (val string) {
	http.Redirect(this.Response(), this.Request(), url, http.StatusFound)
	ghttputil.JustDie()
	return
}

// SetHeader is a intelligent shortcut for ctx.Response().Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `ctx.Response().Header().Del(key)`
func (this *Ctx) SetHeader(key, value string) {
	if value == "" {
		this.response.Header().Del(key)
		return
	}
	this.response.Header().Set(key, value)
}

// Header returns value from request headers.
func (this *Ctx) Header(key string) string {
	return this.request.Header.Get(key)
}

// RequestBody returns raw request body data.
func (this *Ctx) RequestBody() ([]byte, error) {
	return ioutil.ReadAll(this.request.Body)
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (this *Ctx) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(this.response, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (this *Ctx) Cookie(name string) string {
	cookie, err := this.request.Cookie(name)
	if err != nil {
		return ""
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val
}

// WriteFile writes the specified file into the body stream in an efficient way.
func (this *Ctx) WriteFile(filepath string) {
	http.ServeFile(this.response, this.request, filepath)
}

// WriteFileFromFS writes the specified file from http.FileSystem into the body stream in an efficient way.
func (this *Ctx) WriteFileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		this.request.URL.Path = old
	}(this.request.URL.Path)

	this.request.URL.Path = filepath

	http.FileServer(fs).ServeHTTP(this.response, this.request)
}

// WriteFileAttachment writes the specified file into the body stream in an efficient way
// On the client side, the file will typically be downloaded with the given filename
func (this *Ctx) WriteFileAttachment(filepath, filename string) {
	this.response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(this.response, this.request, filepath)
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//     router.GET("/user/:id", func(c gcore.Ctx) {
//         c.FullPath() == "/user/:id" // true
//     })
func (this *Ctx) FullPath() string {
	return this.param.MatchedRoutePath()
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (this *Ctx) Set(key interface{}, value interface{}) {
	this.metadata.Store(key, value)
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (this *Ctx) Get(key interface{}) (value interface{}, exists bool) {
	return this.metadata.Load(key)
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (this *Ctx) MustGet(key interface{}) (v interface{}) {
	v, _ = this.metadata.Load(key)
	return
}

// FormFile returns the first file for the provided form key.
// maxMultipartMemory limits the request form parser using memory byte size.
func (this *Ctx) FormFile(name string, maxMultipartMemory int64) (*multipart.FileHeader, error) {
	if maxMultipartMemory <= 0 {
		maxMultipartMemory = 32 << 20
	}
	if this.request.MultipartForm == nil {
		if err := this.request.ParseMultipartForm(maxMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := this.request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

// MultipartForm is the parsed multipart form, including file uploads.
// maxMultipartMemory limits the request form parser using memory byte size.
func (this *Ctx) MultipartForm(maxMultipartMemory int64) (*multipart.Form, error) {
	if maxMultipartMemory <= 0 {
		maxMultipartMemory = 32 << 20
	}
	err := this.request.ParseMultipartForm(maxMultipartMemory)
	return this.request.MultipartForm, err
}

// SaveUploadedFile uploads the form file to specific dst.
func (this *Ctx) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func NewCtx() *Ctx {
	return &Ctx{
		param:    gcore.Params{},
		metadata: gmap.New(),
	}
}

func NewCtxWithHTTP(w http.ResponseWriter, r *http.Request) *Ctx {
	c := NewCtx()
	c.SetResponse(w)
	c.SetRequest(r)
	return c
}

func NewCtxFromConfig(c gcore.Config) *Ctx {
	return &Ctx{
		config:   c,
		param:    gcore.Params{},
		metadata: gmap.New(),
	}
}

func NewCtxFromConfigFile(file string) (ctx *Ctx, err error) {
	c := gcore.ProviderConfig()()
	c.SetConfigFile(file)
	err = c.ReadInConfig()
	if err != nil {
		return
	}
	ctx = &Ctx{
		config:   c,
		param:    gcore.Params{},
		metadata: gmap.New(),
	}
	return
}
