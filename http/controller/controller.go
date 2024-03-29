// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcontroller

import (
	"fmt"
	"net"
	"net/http"
	"net/textproto"

	gcore "github.com/snail007/gmc/core"
	ghttputil "github.com/snail007/gmc/internal/util/http"
	gcast "github.com/snail007/gmc/util/cast"
)

// Controller gmc web and api controller implements gcore.Controller.
type Controller struct {
	Response     http.ResponseWriter
	Request      *http.Request
	Param        gcore.Params
	Session      gcore.Session
	Tpl          gcore.Template
	I18n         gcore.I18n
	SessionStore gcore.SessionStorage
	Router       gcore.HTTPRouter
	Config       gcore.Config
	Cookie       gcore.Cookies
	Ctx          gcore.Ctx
	View         gcore.View
	Lang         string
	Logger       gcore.Logger
}

func (this *Controller) GetParam() gcore.Params {
	return this.Param
}

func (this *Controller) GetSession() gcore.Session {
	return this.Session
}

func (this *Controller) GetTemplate() gcore.Template {
	return this.Tpl
}

func (this *Controller) GetI18n() gcore.I18n {
	return this.I18n
}

func (this *Controller) GetSessionStore() gcore.SessionStorage {
	return this.SessionStore
}

func (this *Controller) GetRouter() gcore.HTTPRouter {
	return this.Router
}

func (this *Controller) GetCookie() gcore.Cookies {
	return this.Cookie
}

func (this *Controller) GetView() gcore.View {
	return this.View
}

func (this *Controller) GetLang() string {
	return this.Lang
}

func (this *Controller) GetLogger() gcore.Logger {
	return this.Logger
}

func (this *Controller) GetConfig() gcore.Config {
	return this.Config
}

//MethodCallPre called before controller method and Before() if have.
func (this *Controller) MethodCallPre(ctx gcore.Ctx) {
	// 1. init basic objects
	this.Ctx = ctx
	this.Response = ctx.Response()
	this.Request = ctx.Request()
	this.Param = ctx.Param()
	this.Tpl = ctx.WebServer().Tpl()
	this.SessionStore = ctx.WebServer().SessionStore()
	this.Router = ctx.WebServer().Router()
	this.Config = ctx.WebServer().Config()
	this.Logger = ctx.WebServer().Logger()
	this.View = gcore.ProviderView()(ctx.Response(), ctx.Template())
	this.Cookie = gcore.ProviderCookies()(ctx)
	// 2.init stuff below
	this.View.SetLayoutDir(this.Config.GetString("template.layout"))

	//init lang
	this.initLang()

	//init GPSC, delayed, called in view render
	this.View.OnRenderOnce(func() {
		this.initGPSC()
	})
}

// initLang parse browser's accept language to i18n lang flag.
func (this *Controller) initLang() {
	if this.Config.GetBool("i18n.enable") {
		this.Lang = "none"
		t, e := this.Ctx.I18n().MatchAcceptLanguageT(this.Request)
		if e == nil {
			this.Lang = t.String()
			this.View.Set("Lang", this.Lang)
		}
		this.I18n = this.Ctx.I18n().Clone(this.Lang)
	}
}

// init GPSC variables in views
func (this *Controller) initGPSC() {
	g, p, s, c, u, h := map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}

	// get
	for k, v := range this.Request.URL.Query() {
		g[k] = ""
		if len(v) > 0 {
			g[k] = v[0]
		}
	}

	// post
	err := this.Request.ParseForm()
	if err == nil {
		for k, v := range this.Request.PostForm {
			p[k] = ""
			if len(v) > 0 {
				p[k] = v[0]
			}
		}
	}

	// session
	if this.SessionStore != nil && this.Session != nil {
		for k, v := range this.Session.Values() {
			s[gcast.ToString(k)] = gcast.ToString(v)
		}
	}

	// cookie
	for _, v := range this.Request.Cookies() {
		c[v.Name] = v.Value
	}

	// URL
	scheme := "http"
	if this.Ctx.IsWebsocket() {
		scheme = "ws"
	}
	if this.Request.TLS != nil {
		scheme += "s"
	}
	u0 := this.Request.URL
	u["HOST"] = this.Request.Host
	u["HOSTNAME"], u["PORT"], _ = net.SplitHostPort(this.Request.Host)
	u["PATH"] = u0.Path
	u["FRAGMENT"] = u0.Fragment
	u["OPAQUE"] = u0.Opaque
	u["RAW_PATH"] = u0.RawPath
	u["RAW_QUERY"] = u0.RawQuery
	u["SCHEME"] = scheme
	u["USER"], u["PASSWORD"], _ = this.Request.BasicAuth()
	u["URI"] = u0.RequestURI()
	u["URL"] = u0.String()

	// HTTP HEADER
	for k, v := range this.Request.Header {
		h[k] = ""
		if len(v) > 0 {
			h[textproto.CanonicalMIMEHeaderKey(k)] = v[0]
		}
	}

	// fill gpsc to data
	data := map[string]interface{}{
		"G": g,
		"P": p,
		"S": s,
		"C": c,
		"U": u,
		"H": h,
	}

	// set data to view
	this.View.SetMap(data)
}

//MethodCallPost called after controller method and After() if have.
func (this *Controller) MethodCallPost() {
	if this.SessionStore != nil && this.Session != nil {
		if this.Session.IsDestroy() {
			this.SessionStore.Delete(this.Session.SessionID())
		} else {
			err := this.SessionStore.Save(this.Session)
			if err != nil {
				this.Logger.Warnf("save session fail, %s", err)
			}
		}
	}
}

//Tr translates the key to `this.Lang's` text.
func (this *Controller) Tr(key string, defaultText ...string) string {
	return this.I18n.Tr(this.Lang, key, defaultText...)
}

//Die will prevent to call After() if have, and MethodCallPost()
func (this *Controller) Die(msg ...interface{}) {
	ghttputil.Die(this.Response, msg...)
}

//Stop will exit controller method at once
func (this *Controller) Stop(msg ...interface{}) {
	ghttputil.Stop(this.Response, msg...)
}

// StopE will exit controller method if error is not nil.
// First argument is an error.
// Secondary argument is fail function, it be called if error is not nil.
// Third argument is success function, it be called if error is nil.
func (this *Controller) StopE(err interface{}, fn ...func()) {
	ghttputil.StopE(err, fn...)
}

func (this *Controller) SessionStart() (err error) {
	if this.SessionStore == nil {
		err = fmt.Errorf("session is disabled")
		return
	}
	if this.Session != nil {
		//already started
		return
	}
	sessionCookieName := this.Config.GetString("session.cookiename")
	sid, _ := this.Cookie.Get(sessionCookieName)
	var isExists bool
	if sid != "" {
		this.Session, isExists = this.SessionStore.Load(sid)
	}
	if !isExists {
		sess := gcore.ProviderSession()()
		sess.Touch()
		this.Cookie.Set(sessionCookieName, sess.SessionID(), &gcore.CookieOptions{
			Path:     "/",
			MaxAge:   this.Config.GetInt("session.ttl"),
			HTTPOnly: true,
		})
		err = this.SessionStore.Save(sess)
		this.Session = sess
	}

	return
}
func (this *Controller) GetCtx() gcore.Ctx {
	return this.Ctx
}
func (this *Controller) SessionDestroy() (err error) {
	if this.SessionStore == nil {
		err = fmt.Errorf("session is disabled")
		return
	}
	if this.Session != nil {
		this.Session.Destroy()
		sessionCookieName := this.Config.GetString("session.cookiename")
		this.Cookie.Remove(sessionCookieName)
	}
	return
}

func (this *Controller) Write(data ...interface{}) (n int, err error) {
	return ghttputil.Write(this.Response, data...)
}

func (this *Controller) WriteE(data ...interface{}) (n int, err error) {
	this.Response.WriteHeader(http.StatusInternalServerError)
	return ghttputil.Write(this.Response, data...)
}
