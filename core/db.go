// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gcore

import (
	"database/sql"
)

// DBCache gmc abstract db cache layer, to cache db query.
type DBCache interface {
	Set(key string, val []byte, expire uint) (err error)
	Get(key string) (data []byte, err error)
}

// ActiveRecord gmc abstract sql statement builder layer.
type ActiveRecord interface {
	Cache(key string, seconds uint) ActiveRecord
	Delete(table string, where map[string]interface{}) ActiveRecord
	From(from string) ActiveRecord
	FromAs(from, as string) ActiveRecord
	GroupBy(column string) ActiveRecord
	Having(having string) ActiveRecord
	HavingWrap(having, leftWrap, rightWrap string) ActiveRecord
	Insert(table string, data map[string]interface{}) ActiveRecord
	InsertBatch(table string, data []map[string]interface{}) ActiveRecord
	Join(table, as, on, typ string) ActiveRecord
	Limit(limit ...int) ActiveRecord
	OrderBy(column, typ string) ActiveRecord
	Raw(sql string, values ...interface{}) ActiveRecord
	Replace(table string, data map[string]interface{}) ActiveRecord
	ReplaceBatch(table string, data []map[string]interface{}) ActiveRecord
	Reset()
	Select(columns string) ActiveRecord
	SelectNoWrap(columns string) ActiveRecord
	Set(column string, value interface{}) ActiveRecord
	SetNoWrap(column string, value interface{}) ActiveRecord
	SQL() string
	Update(table string, data, where map[string]interface{}) ActiveRecord
	UpdateBatch(table string, values []map[string]interface{}, whereColumn []string) ActiveRecord
	Values() []interface{}
	Where(where map[string]interface{}) ActiveRecord
	WhereWrap(where map[string]interface{}, leftWrap, rightWrap string) ActiveRecord
	Wrap(v string) string
}

// Database gmc abstract db access layer, to connect to db and execute sql.
type Database interface {
	AR() (ar ActiveRecord)
	Stats() sql.DBStats
	Begin() (tx *sql.Tx, err error)
	ExecTx(ar ActiveRecord, tx *sql.Tx) (rs ResultSet, err error)
	ExecSQLTx(tx *sql.Tx, sqlStr string, values ...interface{}) (rs ResultSet, err error)
	Exec(ar ActiveRecord) (rs ResultSet, err error)
	ExecSQL(sqlStr string, values ...interface{}) (rs ResultSet, err error)
	QuerySQL(sqlStr string, values ...interface{}) (rs ResultSet, err error)
	Query(ar ActiveRecord) (rs ResultSet, err error)
}

type DatabaseGroup interface {
	RegistGroup(cfg interface{}) (err error)
	Regist(name string, cfg interface{}) (err error)
	DB(name ...string) (db Database)
}

type ResultSet interface {
	SQL() string
	Len() int
	LastInsertID() int64
	RowsAffected() int64
	TimeUsed() int
	MapRows(keyColumn string) (rowsMap map[string]map[string]string)
	MapStructs(keyColumn string, strucT interface{}, tagName ...string) (structsMap map[string]interface{}, err error)
	Rows() (rows []map[string]string)
	Structs(strucT interface{}, tagName ...string) (structs []interface{}, err error)
	Row() (row map[string]string)
	Struct(strucT interface{}, tagName ...string) (Struct interface{}, err error)
	Values(column string) (values []string)
	MapValues(keyColumn, valueColumn string) (values map[string]string)
	Value(column string) (value string)
}
