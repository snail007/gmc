package gmcdb

import (
	"database/sql"
	"reflect"
)

type Cache interface {
	Set(key string, val []byte, expire uint) (err error)
	Get(key string) (data []byte, err error)
}

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
	Join(table, as, on, type_ string) ActiveRecord
	Limit(limit ...int) ActiveRecord
	OrderBy(column, type_ string) ActiveRecord
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

type Database interface {
	AR() (ar ActiveRecord)
	Stats() sql.DBStats
	Begin() (tx *sql.Tx, err error)
	ExecTx(ar ActiveRecord, tx *sql.Tx) (rs *ResultSet, err error)
	ExecSQLTx(tx *sql.Tx, sqlStr string, values ...interface{}) (rs *ResultSet, err error)
	Exec(ar ActiveRecord) (rs *ResultSet, err error)
	ExecSQL(sqlStr string, values ...interface{}) (rs *ResultSet, err error)
	QuerySQL(sqlStr string, values ...interface{}) (rs *ResultSet, err error)
	Query(ar ActiveRecord) (rs *ResultSet, err error)
}

type DatabaseGroup interface {
	RegistGroup(cfg interface{}) (err error)
	Regist(name string, cfg interface{}) (err error)
	DB(name ...string) (db Database)
}

func isArray(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array
}
func isBool(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Bool
}
