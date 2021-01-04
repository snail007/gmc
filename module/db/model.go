// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"fmt"
	"github.com/snail007/gmc/core"
	"github.com/snail007/gmc/util/cast"
	"sync"
)

type Model struct {
	db         gcore.Database
	table      string
	primaryKey string
	once       *sync.Once
}

func Table(table string, db ...interface{}) *Model {
	m := &Model{
		table:      table,
		primaryKey: table + "_id",
		once:       &sync.Once{},
	}
	var dbDefault interface{}
	if len(db) == 1 {
		dbDefault = db[0]
	} else {
		dbDefault = DB()
	}
	switch v := dbDefault.(type) {
	case string:
		m.db = DB(v)
	case *MySQLDB:
		m.db = v
	case *SQLite3DB:
		m.db = v
	}
	if m.db == nil {
		panic(gcore.Providers.Error("")().New((fmt.Errorf("table db arguments must be 'db string ID' or *gmysql.SQLite3DB or *gsqlite3.SQLite3DB"))))
	}
	return m
}

func (s *Model) QuerySQL(sql string, values ...interface{}) (ret []map[string]string, error error) {
	db := s.db
	ar := db.AR().Raw(sql, values...)
	rs, err := db.Query(ar)
	if err != nil {
		return nil, err
	}
	ret = rs.Rows()
	return
}

func (s *Model) ExecSQL(sql string, values ...interface{}) (lastInsertID, rowsAffected int64, error error) {
	db := s.db
	ar := db.AR().Raw(sql, values...)
	rs, err := db.Exec(ar)
	if err != nil {
		return 0, 0, err
	}
	rowsAffected = rs.RowsAffected()
	lastInsertID = rs.LastInsertID()
	return
}

func (s *Model) GetByID(id string) (ret map[string]string, error error) {
	return s.GetByIDWithFields("*", id)
}

func (s *Model) GetByIDWithFields(fields string, id string) (ret map[string]string, error error) {
	db := s.db
	rs, err := db.Query(db.AR().Select(fields).From(s.table).Where(map[string]interface{}{
		s.primaryKey: id,
	}).Limit(0, 1))
	if err != nil {
		return nil, err
	}
	ret = rs.Row()
	return
}

func (s *Model) GetBy(where map[string]interface{}) (ret map[string]string, error error) {
	return s.GetByWithFields("*", where)
}

func (s *Model) GetByWithFields(fields string, where map[string]interface{}) (ret map[string]string, error error) {
	db := s.db
	rs, err := db.Query(db.AR().Select(fields).From(s.table).Where(where).Limit(0, 1))
	if err != nil {
		return nil, err
	}
	ret = rs.Row()
	return
}

func (s *Model) MGetByIDs(ids []string, orderBy ...interface{}) (ret map[string]string, error error) {
	return s.MGetByIDsWithFields("*", ids, orderBy)
}

func (s *Model) MGetByIDsWithFields(fields string, ids []string, orderBy ...interface{}) (ret map[string]string, error error) {
	db := s.db
	ar := db.AR().Select(fields).From(s.table).Where(map[string]interface{}{
		s.primaryKey: ids,
	})
	if col, by := s.OrderBy(orderBy...); col != "" {
		ar.OrderBy(col, by)
	}
	rs, err := db.Query(ar)
	if err != nil {
		return nil, err
	}
	ret = rs.Row()
	return
}

func (s *Model) GetAll(orderBy ...interface{}) (ret []map[string]string, error error) {
	return s.GetAllWithFields("*", orderBy)
}

func (s *Model) GetAllWithFields(fields string, orderBy ...interface{}) (ret []map[string]string, error error) {
	return s.MGetByWithFields(fields, nil, orderBy...)
}

func (s *Model) MGetBy(where map[string]interface{}, orderBy ...interface{}) (ret []map[string]string, error error) {
	return s.MGetByWithFields("*", where, orderBy...)
}

func (s *Model) MGetByWithFields(fields string, where map[string]interface{}, orderBy ...interface{}) (ret []map[string]string, error error) {
	db := s.db
	ar := db.AR().Select(fields).From(s.table).Where(where).Limit(0, 1)
	if col, by := s.OrderBy(orderBy...); col != "" {
		ar.OrderBy(col, by)
	}
	rs, err := db.Query(ar)
	if err != nil {
		return nil, err
	}
	ret = rs.Rows()
	return
}

func (s *Model) DeleteBy(where map[string]interface{}) (cnt int64, err error) {
	db := s.db
	rs, err := db.Exec(db.AR().Delete(s.table, where))
	if err != nil {
		return 0, err
	}
	cnt = rs.RowsAffected()
	return
}

func (s *Model) DeleteByIDs(ids []string) (cnt int64, err error) {
	db := s.db
	rs, err := db.Exec(db.AR().Delete(s.table, map[string]interface{}{
		s.primaryKey: ids,
	}))
	if err != nil {
		return 0, err
	}
	cnt = rs.RowsAffected()
	return
}

func (s *Model) Insert(data map[string]interface{}) (lastInsertID int64, err error) {
	db := s.db
	rs, err := db.Exec(db.AR().Insert(s.table, data))
	if err != nil {
		return 0, err
	}
	lastInsertID = rs.LastInsertID()
	return
}

func (s *Model) InsertBatch(data []map[string]interface{}) (cnt, lastInsertID int64, err error) {
	db := s.db
	rs, err := db.Exec(db.AR().InsertBatch(s.table, data))
	if err != nil {
		return 0, 0, err
	}
	lastInsertID = rs.LastInsertID()
	cnt = rs.RowsAffected()
	return
}

func (s *Model) UpdateByIDs(ids []string, data map[string]interface{}) (cnt int64, err error) {
	db := s.db
	rs, err := db.Exec(db.AR().Update(s.table, data, map[string]interface{}{
		s.primaryKey: ids,
	}))
	if err != nil {
		return 0, err
	}
	cnt = rs.RowsAffected()
	return
}

func (s *Model) UpdateBy(where, data map[string]interface{}) (cnt int64, err error) {
	db := s.db
	rs, err := db.Exec(db.AR().Update(s.table, data, where))
	if err != nil {
		return 0, err
	}
	cnt = rs.RowsAffected()
	return
}

func (s *Model) Page(where map[string]interface{}, offset, length int, orderBy ...interface{}) (ret []map[string]string, total int, err error) {
	return s.PageWithFields("*", where, offset, length, orderBy)
}

func (s *Model) PageWithFields(fields string, where map[string]interface{}, offset, length int, orderBy ...interface{}) (ret []map[string]string, total int, err error) {
	db := s.db
	ar := db.AR().Select("count(*) as total").From(s.table)
	if len(where) > 0 {
		ar.Where(where)
	}
	rs, err := db.Query(ar)
	if err != nil {
		return nil, 0, err
	}
	total = gcast.ToInt(rs.Value("total"))

	ar = db.AR().Select(fields).From(s.table).Where(where).Limit(offset, length)
	if len(where) > 0 {
		ar.Where(where)
	}
	if col, by := s.OrderBy(orderBy...); col != "" {
		ar.OrderBy(col, by)
	}
	rs, err = db.Query(ar)
	if err != nil {
		return nil, 0, err
	}
	ret = rs.Rows()
	return
}

func (s *Model) List(where map[string]interface{}, offset, length int, orderBy ...interface{}) (ret []map[string]string, err error) {
	return s.ListWithFields("*", where, offset, length, orderBy)
}

func (s *Model) ListWithFields(fields string, where map[string]interface{}, offset, length int, orderBy ...interface{}) (ret []map[string]string, err error) {
	db := s.db
	ar := db.AR().Select(fields).From(s.table).Where(where).Limit(offset, length)
	if len(where) > 0 {
		ar.Where(where)
	}
	if col, by := s.OrderBy(orderBy...); col != "" {
		ar.OrderBy(col, by)
	}
	rs, err := db.Query(ar)
	if err != nil {
		return nil, err
	}
	ret = rs.Rows()
	return
}

func (s *Model) OrderBy(orderBy ...interface{}) (col, by string) {
	if len(orderBy) > 0 {
		switch val := orderBy[0].(type) {
		case map[string]interface{}:
			for k, v := range val {
				col, by = k, v.(string)
				break
			}
		case map[string]string:
			for k, v := range val {
				col, by = k, v
				break
			}
		}
	}
	return
}
