// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gmccontroller

import (
	"fmt"
	"net/http"

	gmcconfig "github.com/snail007/gmc/config/gmc"
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
	Config__() *gmcconfig.GMCConfig
	Cookie__() *gmccookie.Cookies
}
type Controller struct {
	Response     http.ResponseWriter
	Request      *http.Request
	Args         gmcrouter.Params
	Session      *gmcsession.Session
	Tpl          *gmctemplate.Template
	SessionStore gmcsession.Store
	Router       *gmcrouter.HTTPRouter
	Config       *gmcconfig.GMCConfig
	Cookie       *gmccookie.Cookies
}

func (this *Controller) Response__() http.ResponseWriter {
	return this.Response
}
func (this *Controller) Request__() *http.Request {
	return this.Request
}

func (this *Controller) Args__() gmcrouter.Params {
	return this.Args
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
func (this *Controller) Config__() *gmcconfig.GMCConfig {
	return this.Config
}
func (this *Controller) Cookie__() *gmccookie.Cookies {
	return this.Cookie
}

//MethodCallPre__ called before controller method and Before__() if have.
func (this *Controller) MethodCallPre__(w http.ResponseWriter, r *http.Request, ps gmcrouter.Params) {
	this.Response = w
	this.Request = r
	this.Args = ps
	ctxvalue := r.Context().Value(ctxvalue.CtxValueKey).(ctxvalue.CtxValue)
	this.Tpl = ctxvalue.Tpl
	this.SessionStore = ctxvalue.SessionStore
	this.Router = ctxvalue.Router
	this.Config = ctxvalue.Config
	this.Cookie = gmccookie.New(w, r)
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
