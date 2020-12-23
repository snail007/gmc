// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gcontroller

import (
	"fmt"
	gconfig "github.com/snail007/gmc/config"
	gcore "github.com/snail007/gmc/core"
	gi18n "github.com/snail007/gmc/gmc/i18n"
	gview "github.com/snail007/gmc/http/view"
	"github.com/snail007/gmc/util/cast"
	ghttputil "github.com/snail007/gmc/util/http"
	"net/http"

	gcookie "github.com/snail007/gmc/http/cookie"
	grouter "github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/http/server/ctxvalue"
	gsession "github.com/snail007/gmc/http/session"
	gtemplate "github.com/snail007/gmc/http/template"
)

type Controller struct {
	Response     http.ResponseWriter
	Request      *http.Request
	Param        grouter.Params
	Session      *gsession.Session
	Tpl          *gtemplate.Template
	SessionStore gsession.Store
	Router       *grouter.HTTPRouter
	Config       *gconfig.Config
	Cookie       *gcookie.Cookies
	Ctx          gcore.Ctx
	View         *gview.View
	Lang         string
	Logger       gcore.Logger
}

//MethodCallPre called before controller method and Before() if have.
func (this *Controller) MethodCallPre(w http.ResponseWriter, r *http.Request, ps grouter.Params) {
	// 1. init basic objects
	ctxvalue := r.Context().Value(ctxvalue.CtxValueKey).(ctxvalue.CtxValue)
	this.Response = w
	this.Request = r
	this.Param = ps
	this.Tpl = ctxvalue.Tpl
	this.SessionStore = ctxvalue.SessionStore
	this.Router = ctxvalue.Router
	this.Config = ctxvalue.Config
	this.Logger = ctxvalue.Logger

	this.View = gview.New(w, ctxvalue.Tpl)
	this.Cookie = gcookie.New(w, r)
	this.Ctx = grouter.NewCtx(w, r, ps)

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
		t, e := gi18n.MatchAcceptLanguageT(this.Request)
		if e == nil {
			this.Lang = t.String()
			this.View.Set("Lang", this.Lang)
		}
	}
}

// init GPSC variables in views
func (this *Controller) initGPSC() {
	g, p, s, c, u := map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}

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
			s[cast.ToString(k)] = cast.ToString(v)
		}
	}

	// cookie
	for _, v := range this.Request.Cookies() {
		c[v.Name] = v.Value
	}

	// URL
	u0 := this.Request.URL
	u["HOST"] = u0.Host
	u["HOSTNAME"] = u0.Hostname()
	u["PORT"] = u0.Port()
	u["PATH"] = u0.Path
	u["FRAGMENT"] = u0.Fragment
	u["OPAQUE"] = u0.Opaque
	u["RAW_PATH"] = u0.RawPath
	u["RAW_QUERY"] = u0.RawQuery
	u["SCHEME"] = u0.Scheme
	u["USER"] = u0.User.Username()
	u["PASSWORD"], _ = u0.User.Password()
	u["URI"] = u0.RequestURI()
	u["URL"] = u0.String()

	// fill gpsc to data
	data := map[string]interface{}{
		"G": g,
		"P": p,
		"S": s,
		"C": c,
		"U": u,
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
			err:=this.SessionStore.Save(this.Session)
			if err!=nil{
				this.Logger.Warnf("save session fail, %s",err)
			}
		}
	}
}

//Tr translates the key to `this.Lang's` text.
func (this *Controller) Tr(key string, defaultText ...string) string {
	return gi18n.Tr(this.Lang, key, defaultText...)
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
	if this.Session!=nil{
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
		sess := gsession.NewSession()
		sess.Touch()
		this.Cookie.Set(sessionCookieName, sess.SessionID(), &gcookie.Options{
			Path: "/",
			MaxAge: this.Config.GetInt("session.ttl"),
			HTTPOnly:true,
		})
		err = this.SessionStore.Save(sess)
		this.Session = sess
	}

	return
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