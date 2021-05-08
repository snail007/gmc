// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gmc_test

import (
	"bytes"
	"fmt"
	"github.com/snail007/gmc"
	gcore "github.com/snail007/gmc/core"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestGMCImplements(t *testing.T) {
	assert := assert2.New(t)
	ctx := gmc.New.Ctx()
	for _, v := range []struct {
		factory func() (obj interface{}, err error)
		impl    interface{}
		msg     string
	}{
		{func() (obj interface{}, err error) {
			obj, err = gmc.New.APIServer(ctx, "")
			return
		}, (*gcore.APIServer)(nil), "default api server"},
		{func() (obj interface{}, err error) {
			obj = gmc.New.HTTPServer(ctx)
			return
		}, (*gcore.HTTPServer)(nil), "default web server"},
		{func() (obj interface{}, err error) {
			obj = gmc.New.App()
			return
		}, (*gcore.App)(nil), "default app"},
		{func() (obj interface{}, err error) {
			obj = gmc.New.AppDefault()
			return
		}, (*gcore.App)(nil), "default appDefault"},
		{func() (obj interface{}, err error) {
			obj = gcore.ProviderController()(ctx)
			return
		}, (*gcore.Controller)(nil), "default controller"},
		{func() (obj interface{}, err error) {
			cfg := gmc.New.Config()
			ctx.SetConfig(cfg)
			cfg.SetConfigType("toml")
			cfg.ReadConfig(bytes.NewReader([]byte(`
				[database]
				default="mysql"
				[[database.mysql]]
				enable=true
				id="default"
				host="127.0.0.1"
				port="3306"
				username="user"
				password="user"
				database="test"
				prefix=""
				prefix_sql_holder="__PREFIX__"
				charset="utf8"
				collate="utf8_general_ci"
				maxidle=30
				maxconns=200
				timeout=15000
				readtimeout=15000
				writetimeout=15000
				maxlifetimeseconds=1800`)))
			err = gmc.DB.Init(cfg)
			if err != nil {
				return
			}
			obj = gmc.DB.MySQL()
			if obj != gmc.DB.DB() {
				err = fmt.Errorf("gmc.DB() not match gmc.DB.MySQL()")
				return
			}
			db, err := gcore.ProviderDatabase()(ctx)
			if err != nil {
				return
			}
			if db != gmc.DB.DB() {
				err = fmt.Errorf("gcore.Provider.DB() not match gmc.DB.DB()")
				return
			}
			return
		}, (*gcore.Database)(nil), "default database mysql"},
		{func() (obj interface{}, err error) {
			cfg := gmc.New.Config()
			ctx.SetConfig(cfg)
			cfg.SetConfigType("toml")
			cfg.ReadConfig(bytes.NewReader([]byte(`
				[database]
				default="sqlite3"
				[[database.sqlite3]]
				enable=true
				id="default"
				database="test.db"
				password=""
				prefix=""
				prefix_sql_holder="__PREFIX__"
				# sync mode: 0:OFF, 1:NORMAL, 2:FULL, 3:EXTRA
				syncmode=0
				# open mode: ro,rw,rwc,memory
				openmode="rw"
				# cache mode: shared,private
				cachemode="shared"`)))
			err = gmc.DB.Init(cfg)
			if err != nil {
				return
			}
			obj = gmc.DB.SQLite3()
			if obj != gmc.DB.DB() {
				err = fmt.Errorf("gmc.DB() not match gmc.DB.SQLite3()")
				return
			}
			db, err := gcore.ProviderDatabase()(ctx)
			if err != nil {
				return
			}
			if db != gmc.DB.DB() {
				err = fmt.Errorf("gcore.Provider.DB() not match gmc.DB.DB()")
				return
			}
			return
		}, (*gcore.Database)(nil), "default sqlite3 mysql"},
	} {
		obj, err := v.factory()
		assert.Nil(err)
		assert.NotNil(obj)
		assert.Implements(v.impl, obj)
	}
}
