// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package web

import (
	"bytes"
	gcore "github.com/snail007/gmc/core"
	gconfig "github.com/snail007/gmc/module/config"
	gctx "github.com/snail007/gmc/module/ctx"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestGMCImplements(t *testing.T) {
	assert := assert2.New(t)
	ctx := gctx.NewCtx()
	for _, v := range []struct {
		factory func() (obj interface{}, err error)
		impl    interface{}
		msg     string
	}{
		{func() (obj interface{}, err error) {
			cfg := gconfig.New()
			ctx.SetConfig(cfg)
			cfg.SetConfigType("toml")
			cfg.ReadConfig(bytes.NewReader([]byte(`
				[session]
				# turn on/off session
				enable=true
				store="memory"
				cookiename="gmcsid"
				ttl=3600
				[session.memory]
				gctime=300`)))
			obj,err = gcore.ProviderSessionStorage()(ctx)
			return
		}, (*gcore.SessionStorage)(nil), "default ProviderSessionStorage server"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderSession()()
			return
		}, (*gcore.Session)(nil), "default ProviderSession"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderView()(nil,nil)
			return
		}, (*gcore.View)(nil), "default ProviderView"},
		{func() (obj interface{}, err error) {
			obj,err = gcore.ProviderTemplate()(ctx,".")
			return
		}, (*gcore.Template)(nil), "default ProviderTemplate"},
		{func() (obj interface{}, err error) {
			cfg := gconfig.New()
			ctx.SetConfig(cfg)
			cfg.SetConfigType("toml")
			cfg.ReadConfig(bytes.NewReader([]byte(`
				[template]
				dir="views"
				ext=".html"
				delimiterleft="{{"
				delimiterright="}}"
				layout="layout"`)))
			obj,err = gcore.ProviderTemplate()(ctx,".")
			return
		}, (*gcore.Template)(nil), "default ProviderTemplate, form [template]"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderHTTPRouter()(ctx)
			return
		}, (*gcore.HTTPRouter)(nil), "default ProviderHTTPRouter"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderCookies()(ctx)
			return
		}, (*gcore.Cookies)(nil), "default ProviderCookies"},
		{func() (obj interface{}, err error) {
			obj,err = gcore.ProviderI18n()(ctx)
			return
		}, (*gcore.I18n)(nil), "default I18nProvider"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderController()(ctx)
			return
		}, (*gcore.Controller)(nil), "default ProviderController"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderHTTPServer()(ctx)
			return
		}, (*gcore.HTTPServer)(nil), "default ProviderHTTPServer"},
		{func() (obj interface{}, err error) {
			obj,err = gcore.ProviderAPIServer()(ctx,":0")
			return
		}, (*gcore.APIServer)(nil), "default ProviderAPIServer"},
	} {
		obj, err := v.factory()
		assert.Nilf(err,v.msg)
		assert.NotNilf(obj,v.msg)
		assert.Implementsf(v.impl, obj,v.msg)
	}
}
