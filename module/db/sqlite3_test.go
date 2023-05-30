// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"fmt"
	gcast "github.com/snail007/gmc/util/cast"
	"strings"
	"testing"
	"time"

	gcore "github.com/snail007/gmc/core"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
)

func arSqlite3() *SQLite3ActiveRecord {
	ar := new(SQLite3ActiveRecord)
	ar.Reset()
	return ar
}
func TestFrom1(t *testing.T) {
	want := "SELECT * \nFROM `test`"
	got := strings.TrimSpace(arSqlite3().From("test").SQL())
	if want != got {
		t.Errorf("TestFrom , except:%s , got:%s", want, got)
	}
}
func TestFromAs1(t *testing.T) {
	want := "SELECT * \nFROM `test` AS `asname`"
	got := strings.TrimSpace(arSqlite3().FromAs("test", "asname").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestSelect1(t *testing.T) {
	want := "SELECT `a`,`b` \nFROM `test`"
	got := strings.TrimSpace(arSqlite3().From("test").Select("a,b").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestJoin1(t *testing.T) {
	want := "SELECT `u`.`a`,`test`.`b` \nFROM `test` LEFT JOIN `user` AS `u` ON `u`.`a`=`test`.`a`"
	got := strings.TrimSpace(arSqlite3().From("test").Select("u.a,test.b").Join("user", "u", "u.a=test.a", "LEFT").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestWhere2(t *testing.T) {
	_ar := arSqlite3()
	want := "SELECT * \nFROM `test` \nWHERE `addr` = ? AND `name` = ?"
	got := strings.TrimSpace(_ar.From("test").Where(map[string]interface{}{
		"2:name": "kitty",
		"1:addr": "hk",
	}).SQL())
	t.Log(want)
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestWhere11(t *testing.T) {
	_ar := arSqlite3()
	want := "SELECT * \nFROM `test` \nWHERE `name` >= `age`"
	got := strings.TrimSpace(_ar.From("test").Where(map[string]interface{}{
		":`name` >= `age`": "",
	}).SQL())
	assert.Len(t, _ar.values, 0)
	t.Log(want)
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestGroupBy1(t *testing.T) {
	want := "SELECT * \nFROM `test`  \nGROUP BY `name`,`uid`"
	got := strings.TrimSpace(arSqlite3().From("test").GroupBy("name,uid").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestHaving1(t *testing.T) {
	want := "SELECT * \nFROM `test`  \nGROUP BY `name`,`uid` \nHAVING count(uid)>3"
	got := strings.TrimSpace(arSqlite3().From("test").GroupBy("name,uid").Having("count(uid)>3").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestOrderBy1(t *testing.T) {
	want := "SELECT * \nFROM `test`    \nORDER BY `id` DESC,`name` ASC"
	got := strings.TrimSpace(arSqlite3().From("test").OrderBy("id", "desc").OrderBy("name", "asc").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestLimit1(t *testing.T) {
	want := "SELECT * \nFROM `test`     \nLIMIT 0,3"
	got := strings.TrimSpace(arSqlite3().From("test").Limit(0, 3).SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestInsert1(t *testing.T) {
	_ar := arSqlite3()
	want := "INSERT INTO  `test` (`name`,`gid`,`addr`,`is_delete`) \nVALUES (?,?,?,?)"
	got := strings.TrimSpace(_ar.Insert("test", map[string]interface{}{
		"1:name":      "admin",
		"2:gid":       true,
		"3:addr":      nil,
		"4:is_delete": false,
	}).Limit(3).SQL())
	//fmt.Println(_ar.Values())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
	assert.Equal(t, []interface{}{"admin", true, nil, false}, _ar.GetValues())
}
func TestReplace1(t *testing.T) {
	_ar := arSqlite3()
	want := "REPLACE INTO  `test` (`name`,`gid`,`addr`,`is_delete`) \nVALUES (?,?,?,?)"
	got := strings.TrimSpace(_ar.Replace("test", map[string]interface{}{
		"1:name":      "admin",
		"2:gid":       true,
		"3:addr":      nil,
		"4:is_delete": false,
	}).Limit(0, 3).SQL())
	//fmt.Println(_ar.Values())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestInsertBatch1(t *testing.T) {
	_ar := arSqlite3()
	want := "INSERT INTO  `test` (`name`) \nVALUES (?),(?)"
	got := strings.TrimSpace(_ar.InsertBatch("test", []map[string]interface{}{
		{
			"name": "admin11",
		},
		{
			"name": "admin",
		},
	}).SQL())

	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestReplaceBatch1(t *testing.T) {
	_ar := arSqlite3()
	want := "REPLACE INTO  `test` (`name`) \nVALUES (?),(?)"
	got := strings.TrimSpace(_ar.ReplaceBatch("test", []map[string]interface{}{
		{
			"name": "admin11",
		},
		{
			"name": "admin",
		},
	}).SQL())

	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestDelete1(t *testing.T) {
	want := "DELETE FROM  `test`"
	got := strings.TrimSpace(arSqlite3().Delete("test", nil).SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestUpdate1(t *testing.T) {
	_ar := arSqlite3()
	want := "UPDATE  `test` \nSET `addr` = NULL"
	got := strings.TrimSpace(_ar.Update("test", map[string]interface{}{
		"addr": nil,
	}, nil).SQL())
	//fmt.Println(_ar.Values())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestUpdateBatch1(t *testing.T) {
	_ar := arSqlite3()
	want := "UPDATE  `test` \nSET `name` = CASE \nWHEN `gid` = ? THEN ? \nWHEN `gid` = ? THEN ? \nELSE `name` END \nWHERE `gid` IN (?,?)"
	got := strings.TrimSpace(_ar.UpdateBatch("test", []map[string]interface{}{
		{
			"name": "admin11",
			"gid":  true,
		},
		{
			"name": "admin",
			"gid":  33,
		},
	}, []string{"gid"}).SQL())
	//fmt.Println(_ar.Values())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func Test1(t *testing.T) {
	group := NewMySQLDBGroup("default")
	group.Regist("default", NewMySQLDBConfigWith("127.0.0.1", 3306, "test", "root", "admin"))
	group.Regist("blog", NewMySQLDBConfigWith("127.0.0.1", 3306, "test", "root", "admin"))
	group.Regist("www", NewMySQLDBConfigWith("127.0.0.1", 3306, "test", "root", "admin"))
	db := group.DB("www")
	if db != nil {
		rs, err := db.Query(db.AR().From("test"))
		if err != nil {
			t.Errorf("ERR:%s", err)
		} else {
			fmt.Println(rs.Rows())
		}
	} else {
		fmt.Printf("db group config of name %s not found", "www")
	}
}

func TestCount1(t *testing.T) {
	db := db1()
	db.ExecSQL("drop table count_test")
	db.ExecSQL("create table count_test(id int)")
	for i := 1; i <= 3; i++ {
		db.Exec(db.AR().Cache("", 0).Insert("count_test", gmap.M{
			"id": i,
		}))
	}
	count, err := Table("count_test", db).Count(nil)
	assert.Nil(t, err)
	assert.Equal(t, 3, int(count))
	db.ExecSQL("delete table count_test")
}
func db1() gcore.Database {
	group := NewSQLite3DBGroup("default")
	cfg := NewSQLite3DBConfigWith("test_ar.db", OpenModeReadWriteCreate, CacheModeShared, SyncModeOff)
	group.Regist("default", cfg)
	db := group.DB()
	return db
}

type User1 struct {
	Name       string    `column:"name"`
	ID         int       `column:"id"`
	Weight     uint      `column:"weight"`
	Height     float32   `column:"height"`
	Sex        bool      `column:"sex"`
	CreateTime time.Time `column:"create_time"`
	Foo        string    `column:"foo"`
}

var rawRows1 = []map[string][]byte{
	{
		"name":        []byte("jack"),
		"id":          []byte("229"),
		"weight":      []byte("60"),
		"height":      []byte("160.3"),
		"sex":         []byte("1"),
		"create_time": []byte("2017-10-10 09:00:09"),
		"pid":         []byte("1"),
	},
	{
		"name":        []byte("jack"),
		"id":          []byte("229"),
		"weight":      []byte("60"),
		"height":      []byte("160.3"),
		"sex":         []byte("1"),
		"create_time": []byte("2017-10-10 09:00:09"),
		"pid":         []byte("2"),
	},
}
var (
	timeTest1, _ = gcast.StringToDateInDefaultLocation("2017-10-10 09:00:09", time.Local)
)

func TestStruct1(t *testing.T) {
	assert := assert.New(t)
	rs := NewResultSet(&rawRows1)
	s, err := rs.Struct(User1{})
	assert.Nil(err)
	assert.Equal("jack", s.(User1).Name)
	assert.Equal(int(229), s.(User1).ID)
	assert.Equal(uint(60), s.(User1).Weight)
	assert.Equal(float32(160.3), s.(User1).Height)
	assert.True(s.(User1).Sex)
	assert.Equal(timeTest1.Unix(), s.(User1).CreateTime.Unix())
}
func TestStructs1(t *testing.T) {
	assert := assert.New(t)
	rs := NewResultSet(&rawRows1)
	sts, err := rs.Structs(User1{})
	assert.Nil(err)
	for _, s := range sts {
		assert.Equal("jack", s.(User1).Name)
		assert.Equal(int(229), s.(User1).ID)
		assert.Equal(uint(60), s.(User1).Weight)
		assert.Equal(float32(160.3), s.(User1).Height)
		assert.True(s.(User1).Sex)
		assert.Equal(timeTest1.Unix(), s.(User1).CreateTime.Unix())
	}
}
func TestMapStructs1(t *testing.T) {
	assert := assert.New(t)
	rs := NewResultSet(&rawRows1)
	sts, err := rs.MapStructs("pid", User1{})
	assert.Nil(err)
	for _, s := range sts {
		assert.Equal("jack", s.(User1).Name)
		assert.Equal(int(229), s.(User1).ID)
		assert.Equal(uint(60), s.(User1).Weight)
		assert.Equal(float32(160.3), s.(User1).Height)
		assert.True(s.(User1).Sex)
		assert.Equal(timeTest1.Unix(), s.(User1).CreateTime.Unix())
	}
}
func TestUpdateBatch01(t *testing.T) {
	// assert := assert.New(t)
	ar := arSqlite3().UpdateBatch("test", []map[string]interface{}{
		{
			"id":      "id1",
			"gid":     true,
			"name":    "test1",
			"score +": 1,
		}, {
			"id":      "id2",
			"gid":     33,
			"name":    "test2",
			"score +": 2,
		},
	}, []string{"id", "gid"})
	fmt.Println(ar.SQL(), ar.Values())
	// assert.Fail("")
}

func TestMySQLActiveRecord_WhereRaw1(t *testing.T) {
	_ar := arSqlite3()
	want := "SELECT * \nFROM `test` \nWHERE `name` >= `age`"
	got := strings.TrimSpace(_ar.From("test").WhereRaw("`name` >= `age`").SQL())
	assert.Len(t, _ar.values, 0)
	t.Log(want)
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
