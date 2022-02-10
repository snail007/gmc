// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"fmt"
	"strings"
	"testing"
	"time"

	gcore "github.com/snail007/gmc/core"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
)

func ar() *MySQLActiveRecord {
	ar := new(MySQLActiveRecord)
	ar.Reset()
	return ar
}
func TestFrom(t *testing.T) {
	want := "SELECT * \nFROM `test`"
	got := strings.TrimSpace(ar().From("test").SQL())
	if want != got {
		t.Errorf("TestFrom , except:%s , got:%s", want, got)
	}
}
func TestFromAs(t *testing.T) {
	want := "SELECT * \nFROM `test` AS `asname`"
	got := strings.TrimSpace(ar().FromAs("test", "asname").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestSelect(t *testing.T) {
	want := "SELECT `a`,`b` \nFROM `test`"
	got := strings.TrimSpace(ar().From("test").Select("a,b").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestJoin(t *testing.T) {
	want := "SELECT `u`.`a`,`test`.`b` \nFROM `test` LEFT JOIN `user` AS `u` ON `u`.`a`=`test`.`a`"
	got := strings.TrimSpace(ar().From("test").Select("u.a,test.b").Join("user", "u", "u.a=test.a", "LEFT").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestWhere(t *testing.T) {
	_ar := ar()
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
func TestGroupBy(t *testing.T) {
	want := "SELECT * \nFROM `test`  \nGROUP BY `name`,`uid`"
	got := strings.TrimSpace(ar().From("test").GroupBy("name,uid").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestHaving(t *testing.T) {
	want := "SELECT * \nFROM `test`  \nGROUP BY `name`,`uid` \nHAVING count(uid)>3"
	got := strings.TrimSpace(ar().From("test").GroupBy("name,uid").Having("count(uid)>3").SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestOrderBy(t *testing.T) {
	want1 := "SELECT * \nFROM `test`    \nORDER BY `id` DESC,`name` ASC"
	want2 := "SELECT * \nFROM `test`    \nORDER BY `name` ASC,`id` DESC"
	got := strings.TrimSpace(ar().From("test").OrderBy("id", "desc").OrderBy("name", "asc").SQL())
	if (want1 != got) && (want2 != got) {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want1, got)
	}
}
func TestLimit(t *testing.T) {
	want := "SELECT * \nFROM `test`     \nLIMIT 0,3"
	got := strings.TrimSpace(ar().From("test").Limit(0, 3).SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestInsert(t *testing.T) {
	_ar := ar()
	want := "INSERT INTO  `test` (`name`,`gid`,`addr`,`is_delete`) \nVALUES (?,?,?,?)"
	got := strings.TrimSpace(_ar.Insert("test", map[string]interface{}{
		"1:name":      "admin",
		"2:gid":       33,
		"3:addr":      nil,
		"4:is_delete": false,
	}).Limit(0, 3).SQL())
	//fmt.Println(_ar.Values())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestReplace(t *testing.T) {
	_ar := ar()
	want := "REPLACE INTO  `test` (`name`,`gid`,`addr`,`is_delete`) \nVALUES (?,?,?,?)"
	got := strings.TrimSpace(_ar.Replace("test", map[string]interface{}{
		"1:name":      "admin",
		"2:gid":       33,
		"3:addr":      nil,
		"4:is_delete": false,
	}).Limit(0, 3).SQL())
	//fmt.Println(_ar.Values())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestInsertBatch(t *testing.T) {
	_ar := ar()
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
func TestReplaceBatch(t *testing.T) {
	_ar := ar()
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
func TestDelete(t *testing.T) {
	want := "DELETE FROM  `test`"
	got := strings.TrimSpace(ar().Delete("test", nil).SQL())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}
func TestUpdate(t *testing.T) {
	_ar := ar()
	want := "UPDATE  `test` \nSET `addr` = NULL"
	got := strings.TrimSpace(_ar.Update("test", map[string]interface{}{
		"addr": nil,
	}, nil).SQL())
	//fmt.Println(_ar.Values())
	if want != got {
		t.Errorf("\n==> Except : \n%s\n==> Got : \n%s", want, got)
	}
}

func TestUpdateBatch(t *testing.T) {
	_ar := ar()
	want := "UPDATE  `test` \nSET `name` = CASE \nWHEN `gid` = ? THEN ? \nWHEN `gid` = ? THEN ? \nELSE `name` END \nWHERE `gid` IN (?,?)"
	got := strings.TrimSpace(_ar.UpdateBatch("test", []map[string]interface{}{
		{
			"name": "admin11",
			"gid":  22,
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
func Test(t *testing.T) {
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

func TestCount(t *testing.T) {
	db := db()
	db.ExecSQL("drop table count_test")
	db.ExecSQL("create table count_test(id int)")
	for i := 1; i <= 3; i++ {
		db.Exec(db.AR().Insert("count_test", gmap.M{
			"id": i,
		}))
	}
	count, err := Table("count_test", db).Count(nil)
	assert.Nil(t, err)
	assert.Equal(t, 3, int(count))
	db.ExecSQL("delete table count_test")
}
func db() gcore.Database {
	group := NewMySQLDBGroup("default")
	group.Regist("default", NewMySQLDBConfigWith("127.0.0.1", 3306, "test", "root", "admin"))
	return group.DB()
}

type User struct {
	Name       string    `column:"name"`
	ID         int       `column:"id"`
	Weight     uint      `column:"weight"`
	Height     float32   `column:"height"`
	Sex        bool      `column:"sex"`
	CreateTime time.Time `column:"create_time"`
	Foo        string    `column:"foo"`
}

var rawRows = []map[string][]byte{
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

func TestStruct(t *testing.T) {
	assert := assert.New(t)
	rs := NewResultSet(&rawRows)
	s, err := rs.Struct(User{})
	assert.Nil(err)
	assert.Equal("jack", s.(User).Name)
	assert.Equal(int(229), s.(User).ID)
	assert.Equal(uint(60), s.(User).Weight)
	assert.Equal(float32(160.3), s.(User).Height)
	assert.True(s.(User).Sex)
	assert.Equal("2017-10-10 17:00:09 +0800 CST", s.(User).CreateTime.In(time.FixedZone("CST", 3600*8)).String())
}
func TestStructs(t *testing.T) {
	assert := assert.New(t)
	rs := NewResultSet(&rawRows)
	sts, err := rs.Structs(User{})
	assert.Nil(err)
	for _, s := range sts {
		assert.Equal("jack", s.(User).Name)
		assert.Equal(int(229), s.(User).ID)
		assert.Equal(uint(60), s.(User).Weight)
		assert.Equal(float32(160.3), s.(User).Height)
		assert.True(s.(User).Sex)
		assert.Equal("2017-10-10 17:00:09 +0800 CST", s.(User).CreateTime.In(time.FixedZone("CST", 3600*8)).String())
	}
}
func TestMapStructs(t *testing.T) {
	assert := assert.New(t)
	rs := NewResultSet(&rawRows)
	sts, err := rs.MapStructs("pid", User{})
	assert.Nil(err)
	for _, s := range sts {
		assert.Equal("jack", s.(User).Name)
		assert.Equal(int(229), s.(User).ID)
		assert.Equal(uint(60), s.(User).Weight)
		assert.Equal(float32(160.3), s.(User).Height)
		assert.True(s.(User).Sex)
		assert.Equal("2017-10-10 17:00:09 +0800 CST", s.(User).CreateTime.In(time.FixedZone("CST", 3600*8)).String())
	}
}
func TestUpdateBatch0(t *testing.T) {
	// assert := assert.New(t)
	ar := ar().UpdateBatch("test", []map[string]interface{}{
		{
			"id":      "id1",
			"gid":     22,
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
