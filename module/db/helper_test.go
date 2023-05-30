// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"os"
	"testing"

	gcore "github.com/snail007/gmc/core"

	"github.com/stretchr/testify/assert"
)

func TestRegistMysql(t *testing.T) {
	assert := assert.New(t)
	cfg := gcore.ProviderConfig()()
	cfg.SetConfigFile("../app/app.toml")
	err := cfg.ReadInConfig()
	assert.Nil(err)
	err = Init(cfg)
	assert.Nil(err)
	assert.NotNil(DBMySQL())
	assert.EqualValues(1, DBMySQL().Stats().OpenConnections)
}

func TestRegistMysql_1(t *testing.T) {
	assert := assert.New(t)
	cfg := gcore.ProviderConfig()()
	cfg.SetConfigFile("../app/app.toml")
	err := cfg.ReadInConfig()
	assert.Nil(err)
	err = Init(cfg)
	assert.Nil(err)
	db := DBMySQL()
	assert.NotNil(db)
	ar := db.AR().From("test")
	_, err = db.Query(ar)
	assert.Nil(err)
	// fmt.Println(rs.Rows(), err)
	// t.Fail()
	//groupMySQL.SQLite3DB().ConnPool.Close()
}

func TestRegistMysql_2(t *testing.T) {
	assert := assert.New(t)
	err := InitFromFile("../app/app.toml")
	assert.Nil(err)
	db := DB().(*MySQLDB)
	db.Config.TablePrefix = "gmc_"
	db.Config.TablePrefixSQLIdentifier = "__PREFIX__"
	assert.NotNil(db)
	db.ExecSQL("create table __PREFIX__test(id int)")
	ar := db.AR().From("test")
	_, err = db.Query(ar)
	assert.Nil(err)
	_, err = db.Query(db.AR().Raw("select * from __PREFIX__test"))
	assert.Nil(err)
	_, err = db.ExecSQL("drop table __PREFIX__test")
	ar = db.AR().From("test")
	_, err = db.QuerySQL(ar.SQL())
	assert.NotNil(err)
	// fmt.Println(rs.Rows(), err)
	// t.Fail()
	//groupMySQL.SQLite3DB().ConnPool.Close()
}

func TestRegistSQLite3(t *testing.T) {
	os.Remove("test.db")
	assert := assert.New(t)
	cfg := NewSQLite3DBConfig()
	cfg.OpenMode = OpenModeReadWriteCreate
	cfg.Database = "test.db"
	cfg.TablePrefix = "gmc_"
	cfg.TablePrefixSQLIdentifier = "__PREFIX__"
	db, err := NewSQLite3DB(cfg)
	assert.Nil(err)
	_, err = db.ExecSQL("create table __PREFIX__test(id int)")
	assert.Nil(err)
	_, err = db.Query(db.AR().Raw("select * from __PREFIX__test"))
	assert.Nil(err)
	_, err = db.QuerySQL(db.AR().Raw("select * from __PREFIX__test").SQL())
	assert.Nil(err)
	db.ConnPool.Close()
	os.Remove("test.db")
}
func TestRegistSQLite3_1(t *testing.T) {
	os.Remove("test.db")
	assert := assert.New(t)
	err := InitFromFile("testdata/app_db_sqlite3.toml")
	assert.Nil(err)
	db := DB().(*SQLite3DB)
	_, err = db.ExecSQL("create table test(id int)")
	assert.Nil(err)
	db.ConnPool.Close()
	os.Remove("test.db")
}
