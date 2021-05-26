// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package db

import (
	"bytes"
	"fmt"
	gconfig "github.com/snail007/gmc/module/config"
	gdb "github.com/snail007/gmc/module/db"

	gcore "github.com/snail007/gmc/core"
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
			err = gdb.Init(cfg)
			if err != nil {
				return
			}
			obj = gdb.DBMySQL()
			if obj != gdb.DB() {
				err = fmt.Errorf("gmc.DB() not match gmc.DB.MySQL()")
				return
			}
			db, err := gcore.ProviderDatabase()(ctx)
			if err != nil {
				return
			}
			if db != gdb.DB() {
				err = fmt.Errorf("gcore.Provider.DB() not match gmc.DB.DB()")
				return
			}
			return
		}, (*gcore.Database)(nil), "default database mysql"},
		{func() (obj interface{}, err error) {
			cfg := gconfig.New()
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
			err = gdb.Init(cfg)
			if err != nil {
				return
			}
			obj = gdb.DBSQLite3()
			if obj != gdb.DB() {
				err = fmt.Errorf("gmc.DB() not match gmc.DB.SQLite3()")
				return
			}
			db, err := gcore.ProviderDatabase()(ctx)
			if err != nil {
				return
			}
			if db != gdb.DB() {
				err = fmt.Errorf("gcore.Provider.DB() not match gmc.DB.DB()")
				return
			}
			return
		}, (*gcore.Database)(nil), "default database sqlite3"},
	} {
		obj, err := v.factory()
		assert.Nilf(err, v.msg)
		assert.NotNilf(obj, v.msg)
		assert.Implementsf(v.impl, obj, v.msg)
	}
}
