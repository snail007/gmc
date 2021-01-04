// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcache

// import (
// 	"os"
// 	"testing"

// 	gsqlite3 "github.com/snail007/gmc/db/sqlite3"

// 	gconfig "github.com/snail007/gmc/config/gmc"

// 	"github.com/stretchr/testify/assert"
// )

// func TestRegistMysql(t *testing.T) {
// 	assert := assert.New(t)
// 	cfg := gconfig.New()
// 	cfg.SetConfigFile("../../app/app.toml")
// 	err := cfg.ReadInConfig()
// 	assert.Nil(err)
// 	err = RegistGroup(cfg)
// 	assert.Nil(err)
// 	assert.EqualValues(1, DBMySQL().Stats().OpenConnections)
// 	groupMySQL.DB().ConnPool.Close()

// }

// func TestRegistMysql_1(t *testing.T) {
// 	assert := assert.New(t)
// 	cfg := gconfig.New()
// 	cfg.SetConfigFile("../../app/app.toml")
// 	err := cfg.ReadInConfig()
// 	assert.Nil(err)
// 	err = RegistGroup(cfg)
// 	assert.Nil(err)
// 	db := DBMySQL()
// 	assert.NotNil(db)
// 	ar := db.AR().From("test")
// 	rs, err := db.Query(ar)
// 	assert.Nil(err)
// 	fmt.Println(rs.Rows(), err)
// 	t.Fail()
// }
// func TestRegistSQLite3(t *testing.T) {
// 	os.Remove("test.db")
// 	assert := assert.New(t)
// 	cfg := gsqlite3.NewDBConfig()
// 	cfg.OpenMode = gsqlite3.OPEN_MODE_READ_WRITE_CREATE
// 	cfg.Password = "123"
// 	cfg.Database = "test.db"
// 	db, err := gsqlite3.NewDB(cfg)
// 	assert.Nil(err)
// 	_, err = db.ExecSQL("create table test(id int)")
// 	assert.Nil(err)
// 	assert.True(gsqlite3.IsEncrypted("test.db"))
// 	db.ConnPool.Close()
// 	os.Remove("test.db")
// }
// func TestRegistSQLite3_1(t *testing.T) {
// 	os.Remove("test.db")
// 	assert := assert.New(t)
// 	cfg := gsqlite3.NewDBConfig()
// 	cfg.OpenMode = gsqlite3.OPEN_MODE_READ_WRITE_CREATE
// 	cfg.Password = ""
// 	cfg.Database = "test.db"
// 	db, err := gsqlite3.NewDB(cfg)
// 	assert.Nil(err)
// 	_, err = db.ExecSQL("create table test(id int)")
// 	assert.Nil(err)
// 	assert.False(gsqlite3.IsEncrypted("test.db"))
// 	db.ConnPool.Close()
// 	os.Remove("test.db")
// }
