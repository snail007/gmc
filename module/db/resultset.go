// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	gcore "github.com/snail007/gmc/core"
	gcast "github.com/snail007/gmc/util/cast"
	"reflect"
	"strings"
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

var typeOfBytes = reflect.TypeOf([]byte(nil))

func (rs *ResultSet) mapToStruct(mapData map[string]string, Struct interface{}, tagName ...string) (struCt interface{}, err error) {
	tag := "column"
	if len(tagName) == 1 {
		tag = tagName[0]
	}
	rv := reflect.New(reflect.TypeOf(Struct)).Elem()
	if reflect.TypeOf(Struct).Kind() != reflect.Struct {
		return nil, errors.New("v must be struct")
	}
	structType := rv.Type()
	var value interface{}
	for i, fieldCount := 0, rv.NumField(); i < fieldCount; i++ {
		fieldVal := rv.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		field := structType.Field(i)
		fieldType := field.Type
		fieldKind := fieldType.Kind()
		if fieldKind == reflect.Ptr {
			fieldType = fieldType.Elem()
			fieldKind = fieldType.Kind()
		}
		col := strings.Split(field.Tag.Get(tag), ",")[0]
		val, ok := mapData[col]
		if !ok {
			val, ok = mapData[field.Name]
		}
		if !ok {
			continue
		}
	BREAK:
		switch fieldKind {
		case reflect.Uint8:
			value = gcast.ToUint8(val)
		case reflect.Uint16:
			value = gcast.ToUint16(val)
		case reflect.Uint32:
			value = gcast.ToUint32(val)
		case reflect.Uint64:
			value = gcast.ToUint64(val)
		case reflect.Uint:
			value = gcast.ToUint(val)
		case reflect.Int8:
			value = gcast.ToInt8(val)
		case reflect.Int16:
			value = gcast.ToInt16(val)
		case reflect.Int32:
			value = gcast.ToInt32(val)
		case reflect.Int64:
			value = gcast.ToInt64(val)
		case reflect.Int:
			value = gcast.ToInt(val)
		case reflect.String, reflect.Interface:
			value = gcast.ToString(val)
		case reflect.Slice:
			if fieldVal.Type() == typeOfBytes {
				value = []byte(gcast.ToString(val))
			}
		case reflect.Bool:
			value = gcast.ToBool(val)
		case reflect.Float32:
			value = gcast.ToFloat32(val)
		case reflect.Float64:
			value = gcast.ToFloat64(val)
		case reflect.Map, reflect.Struct:
			switch field.Type.Name() {
			case "Time":
				unix, e := gcast.ToInt64E(val)
				if e == nil {
					value = time.Unix(unix, 0).In(time.Local)
				} else if v, e := gcast.StringToDateInDefaultLocation(gcast.ToString(val), time.Local); e == nil {
					value = v
				} else {
					err = e
				}
			default:
				d := []byte(gcast.ToString(val))
				if !json.Valid(d) {
					err = fmt.Errorf("convert json string to map field fail, json format error, field: %s, type: %s", field.Name, fieldKind.String())
					break BREAK
				}
				var iv interface{}
				ivIsPtr := false
				if fieldKind == reflect.Struct {
					ivIsPtr = true
					iv = reflect.New(field.Type).Interface()
				}
				e := json.Unmarshal(d, &iv)
				if e != nil {
					err = fmt.Errorf("unspported json to map or struct field fail, field: %s, type: %s", field.Name, fieldKind.String())
					break BREAK
				}
				if ivIsPtr {
					value = reflect.ValueOf(iv).Elem().Interface()
				} else {
					value = iv
				}
			}
		default:
			err = fmt.Errorf("unspported struct field type, field: %s, type: %s", field.Name, fieldKind.String())
			return nil, err
		}
		rValue := reflect.ValueOf(value)
		if !rValue.IsValid() {
			e := fmt.Errorf("unspported field: %s, type: %s", field.Name, fieldKind.String())
			if err != nil {
				e = errors.Wrapf(err, "convert to field error, field: %s, type: %s", field.Name, field.Type.String())
			}
			return nil, e
		}
		fieldVal.Set(rValue)
	}
	return rv.Interface(), err
}
