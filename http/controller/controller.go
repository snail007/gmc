// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmccontroller

import (
	"fmt"
	gmcview "github.com/snail007/gmc/http/view"
	gmci18n "github.com/snail007/gmc/i18n"
	"github.com/snail007/gmc/util/castutil"
	"net/http"

	gmcconfig "github.com/snail007/gmc/config"

	gmchttputil "github.com/snail007/gmc/util/httputil"

	gmccookie "github.com/snail007/gmc/http/cookie"
	gmcrouter "github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/http/server/ctxvalue"
	gmcsession "github.com/snail007/gmc/http/session"
	gmctemplate "github.com/snail007/gmc/http/template"
)

type IController interface {
	Response__() http.ResponseWriter
	Request__() *http.Request
	Args__() gmcrouter.Params
	Write(data ...interface{}) (n int, err error)
	Session__() *gmcsession.Session
	Tpl__() *gmctemplate.Template
	SessionStore__() gmcsession.Store
	Router__() *gmcrouter.HTTPRouter
	Config__() *gmcconfig.Config
	Cookie__() *gmccookie.Cookies
}
type Controller struct {
	Response     http.ResponseWriter
	Request      *http.Request
	Param        gmcrouter.Params
	Session      *gmcsession.Session
	Tpl          *gmctemplate.Template
	SessionStore gmcsession.Store
	Router       *gmcrouter.HTTPRouter
	Config       *gmcconfig.Config
	Cookie       *gmccookie.Cookies
	Ctx          *gmcrouter.Ctx
	View         *gmcview.View
	Lang         string
}

func (this *Controller) Response__() http.ResponseWriter {
	return this.Response
}
func (this *Controller) Request__() *http.Request {
	return this.Request
}

func (this *Controller) Args__() gmcrouter.Params {
	return this.Param
}
func (this *Controller) Session__() *gmcsession.Session {
	return this.Session
}
func (this *Controller) Tpl__() *gmctemplate.Template {
	return this.Tpl
}
func (this *Controller) SessionStore__() gmcsession.Store {
	return this.SessionStore
}
func (this *Controller) Router__() *gmcrouter.HTTPRouter {
	return this.Router
}
func (this *Controller) Config__() *gmcconfig.Config {
	return this.Config
}
func (this *Controller) Cookie__() *gmccookie.Cookies {
	return this.Cookie
}

//MethodCallPre__ called before controller method and Before__() if have.
func (this *Controller) MethodCallPre__(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params) {
	// 1. init basic objects
	ctxvalue := r.Context().Value(ctxvalue.CtxValueKey).(ctxvalue.CtxValue)
	this.Response = w
	this.Request = r
	this.Param = ps
	this.Tpl = ctxvalue.Tpl
	this.SessionStore = ctxvalue.SessionStore
	this.Router = ctxvalue.Router
	this.Config = ctxvalue.Config

	this.View = gmcview.New(w, ctxvalue.Tpl)
	this.Cookie = gmccookie.New(w, r)
	this.Ctx = gmcrouter.NewCtx(w, r, ps)

	// 2.init stuff below

	//init lang
	this.initLang()

	//init GPSC
	this.initGPSC()
}

// initLang parse browser's accept language to i18n lang flag.
func (this *Controller) initLang() {
	if this.Config.GetBool("i18n.enable") {
		this.Lang = "none"
		t, e := gmci18n.MatchAcceptLanguageT(this.Request)
		if e == nil {
			this.Lang = t.String()
			this.View.Set("Lang", this.Lang)
		}
	}
}

// init GPSC variables in views
func (this *Controller) initGPSC() {
	g, p, s, c := map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}

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

	// cookie
	for _, v := range this.Request.Cookies() {
		c[v.Name] = v.Value
	}

	// fill gpsc to data
	data := map[string]interface{}{
		"G": g,
		"P": p,
		"S": s,
		"C": c,
	}

	// set data to view
	this.View.SetMap(data)
}

//MethodCallPost__ called after controller method and After__() if have.
func (this *Controller) MethodCallPost__() {
	if this.SessionStore != nil && this.Session != nil {
		if this.Session.IsDestory() {
			this.SessionStore.Delete(this.Session.SessionID())
		} else {
			this.SessionStore.Save(this.Session)
		}
	}
}

//Tr translates the key to `this.Lang's` text.
func (this *Controller) Tr(key string, defaultText ...string) string {
	return gmci18n.Tr(this.Lang, key, defaultText...)
}

//Die will prevent to call After__() if have, and MethodCallPost__()
func (this *Controller) Die(msg ...interface{}) {
	gmchttputil.Die(this.Response, msg...)
}

//Stop will exit controller method at once
func (this *Controller) Stop(msg ...interface{}) {
	gmchttputil.Stop(this.Response, msg...)
}
func (this *Controller) SessionStart() (err error) {
	if this.SessionStore == nil {
		err = fmt.Errorf("session is disabled")
		return
	}
	sessionCookieName := this.Config.GetString("session.cookiename")
	sid, _ := this.Cookie.Get(sessionCookieName)
	var isExists bool
	if sid != "" {
		this.Session, isExists = this.SessionStore.Load(sid)
	}
	if !isExists {
		sess := gmcsession.NewSession()
		sess.Touch()
		this.Cookie.Set(sessionCookieName, sess.SessionID(), &gmccookie.Options{
			MaxAge: this.Config.GetInt("session.ttl"),
		})
		err = this.SessionStore.Save(sess)
		this.Session = sess
	}

	// fill session data to view
	s:=map[string]string{}
	if this.SessionStore != nil {
		for k, v := range this.Session.Values() {
			s[castutil.ToString(k)] = castutil.ToString(v)
		}
	}
	this.View.SetMap(map[string]interface{}{
		"S":s,
	})

	return
}

func (this *Controller) SessionDestory() (err error) {
	if this.SessionStore == nil {
		err = fmt.Errorf("session is disabled")
		return
	}
	if this.Session != nil {
		this.Session.Destory()
	}
	return
}

func (this *Controller) Write(data ...interface{}) (n int, err error) {
	return gmchttputil.Write(this.Response, data...)
}
