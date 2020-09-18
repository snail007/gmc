// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/snail007/gmc/config/app"

	"github.com/snail007/gmc/http/cookie"
	"github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/http/session"
)

type Controller struct {
	Response http.ResponseWriter
	Request  *http.Request
	Args     router.Params
	Session  *session.Session
}

//MethodCallPre__ called before controller method and Before__() if have.
func (this *Controller) MethodCallPre__(w http.ResponseWriter, r *http.Request, ps router.Params) {
	this.Response = w
	this.Request = r
	this.Args = ps
}

//MethodCallPost__ called after controller method and After__() if have.
func (this *Controller) MethodCallPost__() {
	if appconfig.SessionStore != nil && this.Session != nil && this.Session.IsDestory() {
		appconfig.SessionStore.Delete(this.Session.SessionID())
	}
}

//Die will prevent to call After__() if have, and MethodCallPost__()
func (this *Controller) Die() {
	panic("__DIE__")
}

//Stop will exit controller method at once
func (this *Controller) Stop() {
	panic("__STOP__")
}
func (this *Controller) SessionStart() (err error) {
	if appconfig.SessionStore == nil {
		err = fmt.Errorf("session is disabled")
		return
	}
	cookies := cookie.New(this.Response, this.Request)
	sid, _ := cookies.Get(appconfig.SessionCookieName)
	var isExists bool
	if sid != "" {
		this.Session, isExists = appconfig.SessionStore.Load(sid)
	}
	if !isExists {
		sess := session.NewSession()
		sess.Touch()
		cookies.Set(appconfig.SessionCookieName, sess.SessionID(), &cookie.Options{
			MaxAge: appconfig.FileConfig.GetInt("session.ttl"),
		})
		err = appconfig.SessionStore.Save(sess)

	}
	return
}

func (this *Controller) SessionDestory() (err error) {
	if appconfig.SessionStore == nil {
		err = fmt.Errorf("session is disabled")
		return
	}
	if this.Session != nil {
		this.Session.Destory()
	}
	return
}

func (this *Controller) Write(data interface{}) (err error) {
	switch v := data.(type) {
	case []byte:
		this.Response.Write(v)
	case string:
		this.Response.Write([]byte(v))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		this.Response.Write([]byte(fmt.Sprintf("%d", v)))
	case bool:
		str := "true"
		if !v {
			str = "false"
		}
		this.Response.Write([]byte(str))
	case float32, float64:
		this.Response.Write([]byte(fmt.Sprintf("%f", v)))
	case error:
		this.Response.Write([]byte(v.Error()))
	default:
		t := reflect.TypeOf(v)
		//map, slice
		jsonType := []string{"[", "map["}
		found := false
		vTypeStr := t.String()
		for _, typ := range jsonType {
			if strings.HasPrefix(vTypeStr, typ) {
				found = true
				var b []byte
				b, err = json.Marshal(v)
				if err == nil {
					this.Response.Write(b)
				}
				break
			}
		}
		if !found {
			this.Response.Write([]byte(fmt.Sprintf("unsupported type to write: %s", t.String())))
		}
	}
	return
}
