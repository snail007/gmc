// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	gcore "github.com/snail007/gmc/core"
	gmap "github.com/snail007/gmc/util/map"
	gvalue "github.com/snail007/gmc/util/value"
	"time"
)

type ResultSet struct {
	rawRows      *[]map[string][]byte
	lastInsertID int64
	rowsAffected int64
	// TimeUsed milliseconds used by execute the SQL statement associated to the result set
	timeUsed time.Duration
	// SQL statement associated to the result set
	sql string
}

func NewResultSet(rawRows *[]map[string][]byte) (rs *ResultSet) {
	rs = &ResultSet{}
	if rawRows != nil {
		rs.rawRows = rawRows
	} else {
		rs.rawRows = &([]map[string][]byte{})
	}
	return
}

func (rs *ResultSet) SQL() string {
	return rs.sql
}

func (rs *ResultSet) Len() int {
	return len(*rs.rawRows)
}

func (rs *ResultSet) LastInsertID() int64 {
	return rs.lastInsertID
}

func (rs *ResultSet) RowsAffected() int64 {
	return rs.rowsAffected
}

func (rs *ResultSet) TimeUsed() time.Duration {
	return rs.timeUsed
}

func (rs *ResultSet) MapRows(keyColumn string) (rowsMap map[string]map[string]string) {
	rowsMap = map[string]map[string]string{}
	for _, row := range *rs.rawRows {
		newRow := map[string]string{}
		for k, v := range row {
			newRow[k] = string(v)
		}
		rowsMap[newRow[keyColumn]] = newRow
	}
	return
}
func (rs *ResultSet) MapStructs(keyColumn string, strucT interface{}, tagName ...string) (structsMap map[string]interface{}, err error) {
	structsMap = map[string]interface{}{}
	for _, row := range *rs.rawRows {
		newRow := map[string]string{}
		for k, v := range row {
			newRow[k] = string(v)
		}
		var _struct interface{}
		_struct, err = rs.mapToStruct(newRow, strucT, tagName...)
		if err != nil {
			return nil, err
		}
		structsMap[newRow[keyColumn]] = _struct
	}
	return
}
func (rs *ResultSet) Rows() (rows []map[string]string) {
	rows = []map[string]string{}
	for _, row := range *rs.rawRows {
		newRow := map[string]string{}
		for k, v := range row {
			newRow[k] = string(v)
		}
		rows = append(rows, newRow)
	}
	return
}
func (rs *ResultSet) Structs(strucT interface{}, tagName ...string) (structs []interface{}, err error) {
	structs = []interface{}{}
	for _, row := range *rs.rawRows {
		newRow := map[string]string{}
		for k, v := range row {
			newRow[k] = string(v)
		}
		var _struct interface{}
		_struct, err = rs.mapToStruct(newRow, strucT, tagName...)
		if err != nil {
			return nil, err
		}
		structs = append(structs, _struct)
	}
	return structs, nil
}
func (rs *ResultSet) Row() (row map[string]string) {
	row = map[string]string{}
	if rs.Len() > 0 {
		row = map[string]string{}
		for k, v := range (*rs.rawRows)[0] {
			row[k] = string(v)
		}
	}
	return
}
func (rs *ResultSet) Struct(strucT interface{}, tagName ...string) (Struct interface{}, err error) {
	if rs.Len() > 0 {
		return rs.mapToStruct(rs.Row(), strucT, tagName...)
	}
	return nil, gcore.ProviderError()().New("rs is empty")
}
func (rs *ResultSet) Values(column string) (values []string) {
	values = []string{}
	for _, row := range *rs.rawRows {
		values = append(values, string(row[column]))
	}
	return
}
func (rs *ResultSet) MapValues(keyColumn, valueColumn string) (values map[string]string) {
	values = map[string]string{}
	for _, row := range *rs.rawRows {
		values[string(row[keyColumn])] = string(row[valueColumn])
	}
	return
}
func (rs *ResultSet) Value(column string) (value string) {
	row := rs.Row()
	if row != nil {
		value, _ = row[column]
	}
	return
}

func (rs *ResultSet) mapToStruct(mapData map[string]string, structValue interface{}, tagName ...string) (_struct interface{}, err error) {
	tag := "column"
	if len(tagName) == 1 {
		tag = tagName[0]
	}
	return gvalue.MapToStructWithTag(gmap.ToAny(mapData), structValue, tag)
}
