package gdb

import (
	gconfig "github.com/snail007/gmc/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistMysql(t *testing.T) {
	assert := assert.New(t)
	cfg := gconfig.NewConfig()
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
	cfg := gconfig.NewConfig()
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
func TestRegistSQLite3(t *testing.T) {
	os.Remove("test.db")
	assert := assert.New(t)
	cfg := NewSQLite3DBConfig()
	cfg.OpenMode = OPEN_MODE_READ_WRITE_CREATE
	cfg.Password = "123"
	cfg.Database = "test.db"
	db, err := NewSQLite3DB(cfg)
	assert.Nil(err)
	_, err = db.ExecSQL("create table test(id int)")
	assert.Nil(err)
	assert.True(IsEncrypted("test.db"))
	db.ConnPool.Close()
	os.Remove("test.db")
}
func TestRegistSQLite3_1(t *testing.T) {
	os.Remove("test.db")
	assert := assert.New(t)
	cfg := NewSQLite3DBConfig()
	cfg.OpenMode = OPEN_MODE_READ_WRITE_CREATE
	cfg.Password = ""
	cfg.Database = "test.db"
	db, err := NewSQLite3DB(cfg)
	assert.Nil(err)
	_, err = db.ExecSQL("create table test(id int)")
	assert.Nil(err)
	assert.False(IsEncrypted("test.db"))
	db.ConnPool.Close()
	os.Remove("test.db")
}
