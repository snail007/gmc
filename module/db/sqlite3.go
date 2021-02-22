// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/snail007/gmc/core"
	makeutil "github.com/snail007/gmc/internal/util/make"
	gmap "github.com/snail007/gmc/util/map"
	"net/url"
	"reflect"
	"strings"
	"time"

	gosqlcipher "github.com/snail007/go-sqlcipher"
)

type SQLite3DBGroup struct {
	defaultConfigKey string
	config           map[string]SQLite3DBConfig
	dbGroup          map[string]*SQLite3DB
	cache            gcore.DBCache
}

func NewSQLite3DBGroupCache(defaultConfigName string, cache gcore.DBCache) (group *SQLite3DBGroup) {
	group = &SQLite3DBGroup{}
	group.defaultConfigKey = defaultConfigName
	group.config = map[string]SQLite3DBConfig{}
	group.dbGroup = map[string]*SQLite3DB{}
	group.cache = cache
	return
}
func NewSQLite3DBGroup(defaultConfigName string) (group *SQLite3DBGroup) {
	group = &SQLite3DBGroup{}
	group.defaultConfigKey = defaultConfigName
	group.config = map[string]SQLite3DBConfig{}
	group.dbGroup = map[string]*SQLite3DB{}
	return
}
func (g *SQLite3DBGroup) RegistGroup(cfg interface{}) (err error) {
	g.config = cfg.(map[string]SQLite3DBConfig)
	for name, config := range g.config {
		if config.Cache == nil {
			config.Cache = g.cache
		}
		g.Regist(name, config)
		if err != nil {
			return
		}
	}
	return
}
func (g *SQLite3DBGroup) Regist(name string, cfgI interface{}) (err error) {
	var db SQLite3DB
	cfg := cfgI.(SQLite3DBConfig)
	if cfg.Cache == nil {
		cfg.Cache = g.cache
	}
	db, err = NewSQLite3DB(cfg)
	if err != nil {
		return
	}
	g.config[name] = cfg
	g.dbGroup[name] = &db
	return
}
func (g *SQLite3DBGroup) DB(name ...string) (db gcore.Database) {
	key := ""
	if len(name) == 0 {
		key = g.defaultConfigKey
	} else {
		key = name[0]
	}
	db0, ok := g.dbGroup[key]
	if ok {
		return db0
	}
	return
}

type SQLite3DB struct {
	Config   SQLite3DBConfig
	ConnPool *sql.DB
	DSN      string
}

func NewSQLite3DB(config SQLite3DBConfig) (db SQLite3DB, err error) {
	db = SQLite3DB{}
	err = db.init(config)
	return
}
func (db *SQLite3DB) init(config SQLite3DBConfig) (err error) {
	db.Config = config
	db.DSN = db.getDSN()
	db.ConnPool, err = db.getDB()
	return
}

func (db *SQLite3DB) getDSN() string {
	//_pragma_key=x'%s'&_pragma_cipher_page_size=4096
	if db.Config.Password == "" {
		return fmt.Sprintf("file:%s?cache=%s&mode=%s",
			url.QueryEscape(db.Config.Database),
			url.QueryEscape(db.Config.CacheMode),
			url.QueryEscape(db.Config.OpenMode),
		)
	}
	return fmt.Sprintf("file:%s?cache=%s&mode=%s&_pragma_key=x'%s'&_pragma_cipher_page_size=4096",
		url.QueryEscape(db.Config.Database),
		url.QueryEscape(db.Config.CacheMode),
		url.QueryEscape(db.Config.OpenMode),
		url.QueryEscape(db.md5(db.Config.Password)),
	)
}
func (db *SQLite3DB) getDB() (connPool *sql.DB, err error) {
	connPool, err = sql.Open("sqlite3", db.getDSN())
	if err != nil {
		return
	}
	connPool.SetMaxIdleConns(0)
	connPool.SetMaxOpenConns(0)
	return
}
func (db *SQLite3DB) md5(password string) string {
	s := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return strings.ToUpper(s + s)
}
func (db *SQLite3DB) AR() gcore.ActiveRecord {
	ar := new(SQLite3ActiveRecord)
	ar.Reset()
	ar.tablePrefix = db.Config.TablePrefix
	ar.tablePrefixSQLIdentifier = db.Config.TablePrefixSQLIdentifier
	return ar
}
func (db *SQLite3DB) IsEncrypted() bool {
	ok, e := gosqlcipher.IsEncrypted(db.Config.Database)
	if e != nil {
		return false
	}
	return ok
}
func (db *SQLite3DB) Stats() sql.DBStats {
	return db.ConnPool.Stats()
}
func (db *SQLite3DB) Begin() (tx *sql.Tx, err error) {
	return db.ConnPool.Begin()
}
func (db *SQLite3DB) ExecTx(ar0 gcore.ActiveRecord, tx *sql.Tx) (rs gcore.ResultSet, err error) {
	ar := ar0.(*SQLite3ActiveRecord)
	return db.execSQLTx(ar.SQL(), len(ar.arInsertBatch), tx, ar.values...)
}
func (db *SQLite3DB) ExecSQLTx(tx *sql.Tx, sqlStr string, values ...interface{}) (rs gcore.ResultSet, err error) {
	return db.execSQLTx(sqlStr, 0, tx, values...)
}
func (db *SQLite3DB) execSQLTx(sqlStr string, arInsertBatchCnt int, tx *sql.Tx, values ...interface{}) (rs gcore.ResultSet, err error) {
	start := time.Now().UnixNano()
	var stmt *sql.Stmt
	var result sql.Result

	stmt, err = tx.Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()
	result, err = stmt.Exec(values...)
	if err != nil {
		return
	}
	rsRaw := new(ResultSet)
	rsRaw.rowsAffected, err = result.RowsAffected()
	rsRaw.lastInsertID, err = result.LastInsertId()
	rsRaw.timeUsed = int((start - time.Now().UnixNano()) / 1e6)
	rsRaw.sql = sqlStr
	if err != nil {
		return
	}
	l := int64(arInsertBatchCnt)
	if l > 1 {
		rsRaw.lastInsertID = rsRaw.lastInsertID - +1
		rsRaw.rowsAffected = l
	}
	rs = rsRaw
	return
}
func (db *SQLite3DB) Exec(ar0 gcore.ActiveRecord) (rs gcore.ResultSet, err error) {
	ar := ar0.(*SQLite3ActiveRecord)
	return db.execSQL(ar.SQL(), len(ar.arInsertBatch), ar.values...)
}
func (db *SQLite3DB) ExecSQL(sqlStr string, values ...interface{}) (rs gcore.ResultSet, err error) {
	return db.execSQL(sqlStr, 0, values...)
}
func (db *SQLite3DB) execSQL(sqlStr string, arInsertBatchCnt int, values ...interface{}) (rs gcore.ResultSet, err error) {
	start := time.Now().UnixNano()
	var stmt *sql.Stmt
	var result sql.Result

	stmt, err = db.ConnPool.Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()
	result, err = stmt.Exec(values...)
	if err != nil {
		return
	}
	rsRaw := new(ResultSet)
	rsRaw.rowsAffected, err = result.RowsAffected()
	rsRaw.lastInsertID, err = result.LastInsertId()
	rsRaw.timeUsed = int((start - time.Now().UnixNano()) / 1e6)
	rsRaw.sql = sqlStr
	l := int64(arInsertBatchCnt)
	if l > 1 {
		rsRaw.lastInsertID = rsRaw.lastInsertID - +1
		rsRaw.rowsAffected = l
	}
	rs = rsRaw
	return
}
func (db *SQLite3DB) QuerySQL(sqlStr string, values ...interface{}) (rs gcore.ResultSet, err error) {
	start := time.Now().UnixNano()
	var results []map[string][]byte
	var stmt *sql.Stmt
	stmt, err = db.ConnPool.Prepare(sqlStr)
	if err != nil {
		return
	}
	defer stmt.Close()
	var rows *sql.Rows
	rows, err = stmt.Query(values...)
	if err != nil {
		return
	}
	defer rows.Close()
	cols, e := rows.Columns()
	if e != nil {
		return nil, e
	}
	closCnt := len(cols)

	// scans := make([]interface{},closCnt)
	var scans []interface{}
	scans = makeutil.GetX(scans, uint64(len(cols)), func() interface{} {
		a := make([]interface{}, closCnt)
		for i := 0; i < closCnt; i++ {
			a[i] = new([]byte)
		}
		return a
	}).([]interface{})
	defer func() {
		for i := 0; i < closCnt; i++ {
			scans[i] = new([]byte)
		}
		makeutil.PutX(scans, uint64(len(cols)))
	}()

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return
		}
		row := map[string][]byte{}
		for i := range cols {
			row[cols[i]] = *(scans[i].(*[]byte))
		}
		results = append(results, row)
	}
	rsRaw := NewResultSet(&results)
	rsRaw.timeUsed = int((start - time.Now().UnixNano()) / 1e6)
	rsRaw.sql = sqlStr
	rs = rsRaw
	return
}
func (db *SQLite3DB) Query(ar0 gcore.ActiveRecord) (rs gcore.ResultSet, err error) {
	ar := ar0.(*SQLite3ActiveRecord)
	start := time.Now().UnixNano()
	var results []map[string][]byte
	if ar.cacheKey != "" {
		var data []byte
		data, err = db.Config.Cache.Get(ar.cacheKey)
		if err == nil {
			d := gob.NewDecoder(bytes.NewReader(data))
			err = d.Decode(&results)
			if err != nil {
				return
			}
		}
	}
	if results == nil || len(results) == 0 {
		sqlStr := ar.SQL()
		var stmt *sql.Stmt
		stmt, err = db.ConnPool.Prepare(sqlStr)
		if err != nil {
			return
		}
		defer stmt.Close()
		var rows *sql.Rows
		rows, err = stmt.Query(ar.values...)
		if err != nil {
			return
		}
		defer rows.Close()
		cols := []string{}
		cols, err = rows.Columns()
		if err != nil {
			return
		}
		closCnt := len(cols)

		// scans := make([]interface{},closCnt)
		var scans []interface{}
		scans = makeutil.GetX(scans, uint64(len(cols)), func() interface{} {
			a := make([]interface{}, closCnt)
			for i := 0; i < closCnt; i++ {
				a[i] = new([]byte)
			}
			return a
		}).([]interface{})
		defer func() {
			for i := 0; i < closCnt; i++ {
				scans[i] = new([]byte)
			}
			makeutil.PutX(scans, uint64(len(cols)))
		}()

		for rows.Next() {
			err = rows.Scan(scans...)
			if err != nil {
				return
			}
			row := map[string][]byte{}
			for i := range cols {
				row[cols[i]] = *(scans[i].(*[]byte))
			}
			results = append(results, row)
		}
		if ar.cacheKey != "" {
			b := new(bytes.Buffer)
			e := gob.NewEncoder(b)
			err = e.Encode(results)
			if err != nil {
				return
			}
			err = db.Config.Cache.Set(ar.cacheKey, b.Bytes(), ar.cacheSeconds)
			if err != nil {
				return
			}
		}
	}
	rsRaw := NewResultSet(&results)
	rsRaw.timeUsed = int((start - time.Now().UnixNano()) / 1e6)
	rsRaw.sql = ar.SQL()
	rs = rsRaw
	return
}

const (
	OpenModeReadOnly        = "ro"
	OpenModeReadWrite       = "rw"
	OpenModeReadWriteCreate = "rwc"
	OpenModeMemory          = "memory"
	CacheModeShared         = "shared"
	CacheModePrivate        = "private"
	SyncModeOff             = 0
	SyncModeNormal          = 1
	SyncModeFull            = 2
	SyncModeExtra           = 3
)

type SQLite3DBConfig struct {
	Database                 string
	TablePrefix              string
	TablePrefixSQLIdentifier string
	Cache                    gcore.DBCache
	SyncMode                 int
	OpenMode                 string
	CacheMode                string
	Password                 string
}

func NewSQLite3DBConfigWith(dbfilename, password, openMode, cacheMode string, syncMode int) (cfg SQLite3DBConfig) {
	cfg = NewSQLite3DBConfig()
	cfg.Database = dbfilename
	cfg.OpenMode = openMode
	cfg.CacheMode = cacheMode
	cfg.SyncMode = syncMode
	cfg.Password = password
	return
}
func NewSQLite3DBConfig() SQLite3DBConfig {
	return SQLite3DBConfig{
		OpenMode:                 OpenModeReadWrite,
		CacheMode:                CacheModeShared,
		SyncMode:                 SyncModeOff,
		Database:                 "test",
		TablePrefix:              "",
		TablePrefixSQLIdentifier: "",
		Password:                 "",
	}
}

type SQLite3ActiveRecord struct {
	arSelect                 [][]interface{}
	arFrom                   []string
	arJoin                   [][]string
	arWhere                  [][]interface{}
	arGroupBy                []string
	arHaving                 [][]interface{}
	arOrderBy                map[string]string
	arLimit                  string
	arSet                    map[string][]interface{}
	arUpdateBatch            []interface{}
	arInsert                 gmap.M
	arInsertBatch            []gmap.M
	asTable                  map[string]bool
	values                   []interface{}
	sqlType                  string
	currentSQL               string
	tablePrefix              string
	tablePrefixSQLIdentifier string
	cacheKey                 string
	cacheSeconds             uint
}

func (ar *SQLite3ActiveRecord) Cache(key string, seconds uint) gcore.ActiveRecord {
	ar.cacheKey = key
	ar.cacheSeconds = seconds
	return ar
}
func (ar *SQLite3ActiveRecord) getValues() []interface{} {
	return ar.values
}
func (ar *SQLite3ActiveRecord) Reset() {
	ar.arSelect = [][]interface{}{}
	ar.arFrom = []string{}
	ar.arJoin = [][]string{}
	ar.arWhere = [][]interface{}{}
	ar.arGroupBy = []string{}
	ar.arHaving = [][]interface{}{}
	ar.arOrderBy = map[string]string{}
	ar.arLimit = ""
	ar.arSet = map[string][]interface{}{}
	ar.arUpdateBatch = []interface{}{}
	ar.arInsert = gmap.M{}
	ar.arInsertBatch = []gmap.M{}
	ar.asTable = map[string]bool{}
	ar.values = []interface{}{}
	ar.sqlType = "select"
	ar.currentSQL = ""
	ar.cacheKey = ""
	ar.cacheSeconds = 0
}

func (ar *SQLite3ActiveRecord) Select(columns string) gcore.ActiveRecord {
	return ar._select(columns, true)
}
func (ar *SQLite3ActiveRecord) SelectNoWrap(columns string) gcore.ActiveRecord {
	return ar._select(columns, false)
}

func (ar *SQLite3ActiveRecord) _select(columns string, wrap bool) gcore.ActiveRecord {
	for _, column := range strings.Split(columns, ",") {
		ar.arSelect = append(ar.arSelect, []interface{}{column, wrap})
	}
	return ar
}
func (ar *SQLite3ActiveRecord) From(from string) gcore.ActiveRecord {
	ar.FromAs(from, "")
	return ar
}
func (ar *SQLite3ActiveRecord) FromAs(from, as string) gcore.ActiveRecord {
	ar.arFrom = []string{from, as}
	if as != "" {
		ar.asTable[as] = true
	}
	return ar
}

func (ar *SQLite3ActiveRecord) Join(table, as, on, typ string) gcore.ActiveRecord {
	ar.arJoin = append(ar.arJoin, []string{table, as, on, typ})
	return ar
}
func (ar *SQLite3ActiveRecord) Where(where gmap.M) gcore.ActiveRecord {
	if len(where) > 0 {
		ar.WhereWrap(where, "AND", "")
	}
	return ar
}
func (ar *SQLite3ActiveRecord) WhereWrap(where gmap.M, leftWrap, rightWrap string) gcore.ActiveRecord {
	if len(where) > 0 {
		ar.arWhere = append(ar.arWhere, []interface{}{where, leftWrap, rightWrap, len(ar.arWhere)})
	}
	return ar
}
func (ar *SQLite3ActiveRecord) GroupBy(column string) gcore.ActiveRecord {
	for _, columnCurrent := range strings.Split(column, ",") {
		ar.arGroupBy = append(ar.arGroupBy, strings.TrimSpace(columnCurrent))
	}
	return ar
}
func (ar *SQLite3ActiveRecord) Having(having string) gcore.ActiveRecord {
	ar.HavingWrap(having, "AND", "")
	return ar
}
func (ar *SQLite3ActiveRecord) HavingWrap(having, leftWrap, rightWrap string) gcore.ActiveRecord {
	ar.arHaving = append(ar.arHaving, []interface{}{having, leftWrap, rightWrap, len(ar.arHaving)})
	return ar
}

func (ar *SQLite3ActiveRecord) OrderBy(column, typ string) gcore.ActiveRecord {
	ar.arOrderBy[column] = typ
	return ar
}

//Limit Limit(offset,count) or Limit(count)
func (ar *SQLite3ActiveRecord) Limit(limit ...int) gcore.ActiveRecord {
	if len(limit) == 1 {
		ar.arLimit = fmt.Sprintf("%d", limit[0])

	} else if len(limit) == 2 {
		ar.arLimit = fmt.Sprintf("%d,%d", limit[0], limit[1])
	} else {
		ar.arLimit = ""
	}
	return ar
}

func (ar *SQLite3ActiveRecord) Insert(table string, data gmap.M) gcore.ActiveRecord {
	ar.sqlType = "insert"
	ar.arInsert = data
	ar.From(table)
	return ar
}
func (ar *SQLite3ActiveRecord) Replace(table string, data gmap.M) gcore.ActiveRecord {
	ar.sqlType = "replace"
	ar.arInsert = data
	ar.From(table)
	return ar
}

func (ar *SQLite3ActiveRecord) InsertBatch(table string, data []gmap.M) gcore.ActiveRecord {
	ar.sqlType = "insertBatch"
	ar.arInsertBatch = data
	ar.From(table)
	return ar
}
func (ar *SQLite3ActiveRecord) ReplaceBatch(table string, data []gmap.M) gcore.ActiveRecord {
	ar.InsertBatch(table, data)
	ar.sqlType = "replaceBatch"
	return ar
}

func (ar *SQLite3ActiveRecord) Delete(table string, where gmap.M) gcore.ActiveRecord {
	ar.From(table)
	ar.Where(where)
	ar.sqlType = "delete"
	return ar
}
func (ar *SQLite3ActiveRecord) Update(table string, data, where gmap.M) gcore.ActiveRecord {
	ar.From(table)
	ar.Where(where)
	_data := sortMap(data, true)
	for _, val := range _data {
		k, v := val["col"].(string), val["value"]
		if isBool(v) {
			value := 0
			if v.(bool) {
				value = 1
			}
			ar.Set(k, value)
		} else if v == nil {
			ar.SetNoWrap(k, "NULL")
		} else {
			ar.Set(k, v)
		}
	}
	return ar
}
func (ar *SQLite3ActiveRecord) UpdateBatch(table string, values []gmap.M, whereColumn []string) gcore.ActiveRecord {
	ar.From(table)
	ar.sqlType = "updateBatch"
	ar.arUpdateBatch = []interface{}{values, whereColumn}
	if len(values) > 0 {
		for _, whereCol := range whereColumn {
			ids := []interface{}{}
			for _, val := range values {
				ids = append(ids, val[whereCol])
			}
			ar.Where(gmap.M{whereCol: ids})
		}
	}
	return ar
}

func (ar *SQLite3ActiveRecord) Set(column string, value interface{}) gcore.ActiveRecord {
	ar.sqlType = "update"
	ar.arSet[column] = []interface{}{value, true}
	return ar
}
func (ar *SQLite3ActiveRecord) SetNoWrap(column string, value interface{}) gcore.ActiveRecord {
	ar.sqlType = "update"
	ar.arSet[column] = []interface{}{value, false}
	return ar
}
func (ar *SQLite3ActiveRecord) Wrap(v string) string {
	columns := strings.Split(v, ".")
	if len(columns) == 2 {
		return ar.protectIdentifier(ar.checkPrefix(columns[0])) + "." + ar.checkPrefix(columns[1])
	}
	return ar.protectIdentifier(ar.checkPrefix(columns[0]))
}
func (ar *SQLite3ActiveRecord) Raw(sql string, values ...interface{}) gcore.ActiveRecord {
	ar.currentSQL = sql
	if len(values) > 0 {
		ar.values = append(ar.values, values...)
	}
	return ar
}
func (ar *SQLite3ActiveRecord) Values() []interface{} {
	return ar.values
}
func (ar *SQLite3ActiveRecord) SQL() string {
	if ar.currentSQL != "" {
		return ar.currentSQL
	}
	switch ar.sqlType {
	case "select":
		ar.currentSQL = ar.getSelectSQL()
	case "update":
		ar.currentSQL = ar.getUpdateSQL()
	case "updateBatch":
		ar.currentSQL = ar.getUpdateBatchSQL()
	case "insert":
		ar.currentSQL = ar.getInsertSQL()
	case "insertBatch":
		ar.currentSQL = ar.getInsertBatchSQL()
	case "replace":
		ar.currentSQL = ar.getReplaceSQL()
	case "replaceBatch":
		ar.currentSQL = ar.getReplaceBatchSQL()
	case "delete":
		ar.currentSQL = ar.getDeleteSQL()
	}
	return ar.currentSQL
}
func (ar *SQLite3ActiveRecord) getUpdateSQL() string {
	SQL := []string{"UPDATE "}
	SQL = append(SQL, ar.getFrom())
	SQL = append(SQL, "\nSET")
	SQL = append(SQL, ar.compileSet())
	SQL = append(SQL, ar.getWhere())
	orderBy := strings.TrimSpace(ar.compileOrderBy())
	if orderBy != "" {
		SQL = append(SQL, fmt.Sprintf("\nORDER BY %s", orderBy))
	}
	SQL = append(SQL, ar.getLimit())
	return strings.Join(SQL, " ")
}

func (ar *SQLite3ActiveRecord) getUpdateBatchSQL() string {
	SQL := []string{"UPDATE "}
	SQL = append(SQL, ar.getFrom())
	SQL = append(SQL, "\nSET")
	SQL = append(SQL, ar.compileUpdateBatch())
	SQL = append(SQL, ar.getWhere())
	return strings.Join(SQL, " ")
}
func (ar *SQLite3ActiveRecord) getInsertSQL() string {
	SQL := []string{"INSERT INTO "}
	SQL = append(SQL, ar.getFrom())
	SQL = append(SQL, ar.compileInsert())
	return strings.Join(SQL, " ")
}
func (ar *SQLite3ActiveRecord) getReplaceSQL() string {
	SQL := []string{"REPLACE INTO "}
	SQL = append(SQL, ar.getFrom())
	SQL = append(SQL, ar.compileInsert())
	return strings.Join(SQL, " ")
}
func (ar *SQLite3ActiveRecord) getInsertBatchSQL() string {
	SQL := []string{"INSERT INTO "}
	SQL = append(SQL, ar.getFrom())
	SQL = append(SQL, ar.compileInsertBatch())
	return strings.Join(SQL, " ")
}
func (ar *SQLite3ActiveRecord) getReplaceBatchSQL() string {
	SQL := []string{"REPLACE INTO "}
	SQL = append(SQL, ar.getFrom())
	SQL = append(SQL, ar.compileInsertBatch())
	return strings.Join(SQL, " ")
}
func (ar *SQLite3ActiveRecord) getDeleteSQL() string {
	SQL := []string{"DELETE FROM "}
	SQL = append(SQL, ar.getFrom())
	SQL = append(SQL, ar.getWhere())
	orderBy := strings.TrimSpace(ar.compileOrderBy())
	if orderBy != "" {
		SQL = append(SQL, fmt.Sprintf("\nORDER BY %s", orderBy))
	}
	SQL = append(SQL, ar.getLimit())
	return strings.Join(SQL, " ")
}
func (ar *SQLite3ActiveRecord) getSelectSQL() string {
	from := ar.getFrom()
	where := ar.getWhere()
	having := ""
	for _, w := range ar.arHaving {
		having += ar.compileWhere(w[0], w[1].(string), w[2].(string), w[3].(int))
	}
	having = strings.TrimSpace(having)
	if having != "" {
		having = fmt.Sprintf("\nHAVING %s", having)
	}
	groupBy := strings.TrimSpace(ar.compileGroupBy())
	if groupBy != "" {
		groupBy = fmt.Sprintf("\nGROUP BY %s", groupBy)
	}
	orderBy := strings.TrimSpace(ar.compileOrderBy())
	if orderBy != "" {
		orderBy = fmt.Sprintf("\nORDER BY %s", orderBy)
	}
	limit := ar.getLimit()
	Select := ar.compileSelect()
	return fmt.Sprintf("SELECT %s \nFROM %s %s %s %s %s %s", Select, from, where, groupBy, having, orderBy, limit)
}
func (ar *SQLite3ActiveRecord) compileUpdateBatch() string {
	_values, _index := ar.arUpdateBatch[0], ar.arUpdateBatch[1]
	index := _index.([]string)
	values := _values.([]gmap.M)
	columns := []string{}
	for _, val := range sortMap(values[0], true) {
		k := val["col"].(string)
		_continue := false
		for _, v1 := range index {
			if k == v1 {
				_continue = true
				break
			}
		}
		if _continue {
			continue
		}
		columns = append(columns, k)
	}
	str := ""
	for _, column := range columns {
		_column := column
		realColumnArr := strings.Split(column, " ")
		if len(realColumnArr) == 2 {
			_column = realColumnArr[0]
		}
		str += fmt.Sprintf("%s = CASE \n", ar.protectIdentifier(_column))
		for _, row := range values {
			_when := []string{}
			for _, col := range index {
				_when = append(_when, fmt.Sprintf("%s = ?", ar.protectIdentifier(col)))
				ar.values = append(ar.values, row[col])
			}
			_whenStr := strings.Join(_when, " AND ")
			if len(realColumnArr) == 2 {
				str += fmt.Sprintf("WHEN %s THEN %s %s ? \n", _whenStr, ar.protectIdentifier(_column), realColumnArr[1])
			} else {
				str += fmt.Sprintf("WHEN %s THEN ? \n", _whenStr)
			}
			ar.values = append(ar.values, row[column])
		}
		str += fmt.Sprintf("ELSE %s END,", ar.protectIdentifier(_column))
	}
	return strings.TrimRight(str, " ,")
}

func (ar *SQLite3ActiveRecord) compileInsert() string {
	var columns = []string{}
	var values = []string{}
	data := sortMap(ar.arInsert, true)
	for _, val := range data {
		k, v := val["col"].(string), val["value"]
		columns = append(columns, ar.protectIdentifier(k))
		values = append(values, "?")
		ar.values = append(ar.values, v)
	}
	if len(columns) > 0 {
		return fmt.Sprintf("(%s) \nVALUES (%s)", strings.Join(columns, ","), strings.Join(values, ","))
	}
	return ""
}
func (ar *SQLite3ActiveRecord) compileInsertBatch() string {
	var columns []string
	var values []string
	data := sortMap(ar.arInsertBatch[0], true)
	for _, val := range data {
		col := val["col"].(string)
		columns = append(columns, ar.protectIdentifier(col))
	}
	for _, row := range ar.arInsertBatch {
		_values := []string{}
		for _, col := range columns {
			_values = append(_values, "?")
			ar.values = append(ar.values, row[strings.Trim(col, "`")])
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(_values, ",")))
	}
	return fmt.Sprintf("(%s) \nVALUES %s", strings.Join(columns, ","), strings.Join(values, ","))
}
func (ar *SQLite3ActiveRecord) compileSet() string {
	set := []string{}
	for key, _value := range ar.arSet {
		value, wrap := _value[0], _value[1]
		_column := key
		op := ""
		realColumnArr := strings.Split(key, " ")
		if len(realColumnArr) == 2 {
			_column = realColumnArr[0]
			op = realColumnArr[1]
		}
		if wrap.(bool) {
			if op != "" {
				set = append(set, fmt.Sprintf("%s = %s %s ?", ar.protectIdentifier(_column), ar.protectIdentifier(_column), op))
			} else {
				set = append(set, fmt.Sprintf("%s = ?", ar.protectIdentifier(_column)))
			}
			ar.values = append(ar.values, value)
		} else {
			set = append(set, fmt.Sprintf("%s = %s", ar.protectIdentifier(_column), value))
		}
	}
	return strings.Join(set, ",")
}
func (ar *SQLite3ActiveRecord) compileGroupBy() string {
	groupBy := []string{}
	for _, key := range ar.arGroupBy {
		_key := strings.Split(key, ".")
		if len(_key) == 2 {
			groupBy = append(groupBy, fmt.Sprintf("%s.%s", ar.protectIdentifier(ar.checkPrefix(_key[0])), ar.protectIdentifier(_key[1])))
		} else {
			groupBy = append(groupBy, fmt.Sprintf("%s", ar.protectIdentifier(_key[0])))
		}
	}
	return strings.Join(groupBy, ",")
}

func (ar *SQLite3ActiveRecord) compileOrderBy() string {
	orderBy := []string{}
	for _, val := range sortMapSS(ar.arOrderBy, true) {
		key := val["col"].(string)
		Type := strings.ToUpper(val["value"].(string))
		_key := strings.Split(key, ".")
		if len(_key) == 2 {
			orderBy = append(orderBy, fmt.Sprintf("%s.%s %s", ar.protectIdentifier(ar.checkPrefix(_key[0])), ar.protectIdentifier(_key[1]), Type))

		} else {
			orderBy = append(orderBy, fmt.Sprintf("%s %s", ar.protectIdentifier(_key[0]), Type))
		}
	}
	return strings.Join(orderBy, ",")
}
func (ar *SQLite3ActiveRecord) compileWhere(where0 interface{}, leftWrap, rightWrap string, index int) string {

	_where := []string{}
	if index == 0 {
		str := strings.ToUpper(strings.TrimSpace(leftWrap))
		if strings.Contains(str, "AND") || strings.Contains(str, "OR") {
			leftWrap = ""
		}
	}
	if reflect.TypeOf(where0).Kind() == reflect.String {
		return fmt.Sprintf(" %s %s %s ", leftWrap, where0, rightWrap)
	}
	where := sortMap(where0.(gmap.M), true)
	for _, val := range where {
		key, value := val["col"].(string), val["value"]
		k := ""
		k = strings.TrimSpace(key)
		_key := strings.SplitN(k, " ", 2)
		op := ""
		if len(_key) == 2 {
			op = _key[1]
		}
		keys := strings.Split(_key[0], ".")
		if len(keys) == 2 {
			k = ar.protectIdentifier(ar.checkPrefix(keys[0])) + "." + ar.protectIdentifier(keys[1])
		} else {
			k = ar.protectIdentifier(keys[0])
		}

		if isArray(value) {
			if op != "" {
				op += " IN"
			} else {
				op = "IN"
			}
			op = strings.ToUpper(op)
			l := reflect.ValueOf(value).Len()

			_v := []string{}
			for i := 0; i < l; i++ {
				_v = append(_v, "?")
			}
			_where = append(_where, fmt.Sprintf("%s %s (%s)", k, op, strings.Join(_v, ",")))
			for _, v := range *ar.interface2Slice(value) {
				ar.values = append(ar.values, v)
			}
		} else if isBool(value) {
			if op == "" {
				op = "="
			}
			op = strings.ToUpper(op)
			_v := 0
			if value.(bool) {
				_v = 1
			}
			_where = append(_where, fmt.Sprintf("%s %s ?", k, op))
			ar.values = append(ar.values, _v)
		} else if value == nil {
			if op == "" {
				op = "IS"
			}
			op = strings.ToUpper(op)
			_where = append(_where, fmt.Sprintf("%s %s NULL", k, op))
		} else {
			if op == "" {
				op = "="
			}
			op = strings.ToUpper(op)
			_where = append(_where, fmt.Sprintf("%s %s ?", k, op))
			ar.values = append(ar.values, value)
		}
	}
	return fmt.Sprintf(" %s %s %s ", leftWrap, strings.Join(_where, " AND "), rightWrap)
}
func (ar *SQLite3ActiveRecord) interface2Slice(data interface{}) (arr *[]interface{}) {
	arr = &[]interface{}{}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			e := val.Index(i)
			*arr = append(*arr, e.Interface())
		}
	}
	return
}
func (ar *SQLite3ActiveRecord) compileSelect() string {
	selects := ar.arSelect
	columns := []string{}
	if len(selects) == 0 {
		selects = append(selects, []interface{}{"*", true})
	}
	for _, v := range selects {
		protect := v[1].(bool)
		value := strings.TrimSpace(v[0].(string))
		if value != "*" {
			info := strings.Split(value, ".")
			if len(info) == 2 {
				_v := ar.checkPrefix(info[0])
				if protect {
					info[0] = ar.protectIdentifier(_v)
					info[1] = ar.protectIdentifier(info[1])
				} else {
					info[0] = _v
				}
				value = strings.Join(info, ".")
			} else if protect {
				value = ar.protectIdentifier(value)
			}
		}
		columns = append(columns, value)
	}
	return strings.Join(columns, ",")
}

func (ar *SQLite3ActiveRecord) checkPrefix(v string) string {
	if strings.Contains(v, "(") || strings.Contains(v, ")") || strings.TrimSpace(v) == "*" {
		return v
	}
	if ar.tablePrefix != "" && !strings.Contains(v, ar.tablePrefix) {
		if _, exists := ar.asTable[v]; !exists {
			return ar.tablePrefix + v
		}
	}
	return v
}
func (ar *SQLite3ActiveRecord) protectIdentifier(v string) string {
	if strings.Contains(v, "(") || strings.Contains(v, ")") || strings.TrimSpace(v) == "*" {
		return v
	}
	values := strings.Split(v, " ")
	if len(values) == 3 && strings.ToLower(values[1]) == "as" {
		return fmt.Sprintf("`%s` AS `%s`", values[0], values[2])
	}
	return fmt.Sprintf("`%s`", v)
}
func (ar *SQLite3ActiveRecord) compileFrom(from, as string) string {
	if as != "" {
		ar.asTable[as] = true
		as = " AS " + ar.protectIdentifier(as) + " "
	}
	return ar.protectIdentifier(ar.checkPrefix(from)) + as
}
func (ar *SQLite3ActiveRecord) compileJoin(table, as, on, typ string) string {
	tableUsed := ""
	if as != "" {
		ar.asTable[table] = true
		tableUsed = ar.protectIdentifier(ar.checkPrefix(table)) + " AS " + ar.protectIdentifier(as)
	} else {
		tableUsed = ar.protectIdentifier(ar.checkPrefix(table))
	}
	a := strings.Split(on, "=")
	if len(a) == 2 {
		left := strings.Split(a[0], ".")
		right := strings.Split(a[1], ".")
		left[0] = ar.protectIdentifier(ar.checkPrefix(left[0]))
		left[1] = ar.protectIdentifier(left[1])
		right[0] = ar.protectIdentifier(ar.checkPrefix(right[0]))
		right[1] = ar.protectIdentifier(right[1])
		on = strings.Join(left, ".") + "=" + strings.Join(right, ".")
	}
	return fmt.Sprintf(" %s JOIN %s ON %s ", typ, tableUsed, on)
}

func (ar *SQLite3ActiveRecord) getFrom() string {
	table := ar.compileFrom(ar.arFrom[0], ar.arFrom[1])
	for _, v := range ar.arJoin {
		table += ar.compileJoin(v[0], v[1], v[2], v[3])
	}
	return table
}
func (ar *SQLite3ActiveRecord) getLimit() string {
	limit := ar.arLimit
	if limit != "" {
		limit = fmt.Sprintf("\nLIMIT %s", limit)
	}
	return limit
}
func (ar *SQLite3ActiveRecord) getWhere() string {
	where := []string{}
	hasEmptyIn := false

	for _, v := range ar.arWhere {
		for _, value := range v[0].(gmap.M) {
			if isArray(value) && reflect.ValueOf(value).Len() == 0 {
				hasEmptyIn = true
				break
			}
		}
		if hasEmptyIn {
			break
		}
		where = append(where, ar.compileWhere(v[0].(gmap.M), v[1].(string), v[2].(string), v[3].(int)))
	}
	if hasEmptyIn {
		return "WHERE 0"
	}
	allWhere := strings.TrimSpace(strings.Join(where, ""))
	if allWhere != "" {
		allWhere = fmt.Sprintf("\nWHERE %s", allWhere)
	}
	return allWhere
}

func IsEncrypted(file string) bool {
	ok, e := gosqlcipher.IsEncrypted(file)
	if e != nil {
		return false
	}
	return ok
}
