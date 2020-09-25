package gmcdb

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

type Cache interface {
	Set(key string, val []byte, expire uint) (err error)
	Get(key string) (data []byte, err error)
}

type ResultSet struct {
	rawRows      *[]map[string][]byte
	LastInsertId int64
	RowsAffected int64
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
func (rs *ResultSet) Len() int {
	return len(*rs.rawRows)
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
func (rs *ResultSet) MapStructs(keyColumn string, strucT interface{}) (structsMap map[string]interface{}, err error) {
	structsMap = map[string]interface{}{}
	for _, row := range *rs.rawRows {
		newRow := map[string]string{}
		for k, v := range row {
			newRow[k] = string(v)
		}
		var _struct interface{}
		_struct, err = rs.mapToStruct(newRow, strucT)
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
func (rs *ResultSet) Structs(strucT interface{}) (structs []interface{}, err error) {
	structs = []interface{}{}
	for _, row := range *rs.rawRows {
		newRow := map[string]string{}
		for k, v := range row {
			newRow[k] = string(v)
		}
		var _struct interface{}
		_struct, err = rs.mapToStruct(newRow, strucT)
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
func (rs *ResultSet) Struct(strucT interface{}) (Struct interface{}, err error) {
	if rs.Len() > 0 {
		return rs.mapToStruct(rs.Row(), strucT)
	}
	return nil, errors.New("rs is empty")
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
func (rs *ResultSet) mapToStruct(mapData map[string]string, Struct interface{}) (struCt interface{}, err error) {
	rv := reflect.New(reflect.TypeOf(Struct)).Elem()
	if reflect.TypeOf(Struct).Kind() != reflect.Struct {
		return nil, errors.New("v must be struct")
	}
	fieldType := rv.Type()
	for i, fieldCount := 0, rv.NumField(); i < fieldCount; i++ {
		fieldVal := rv.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		structField := fieldType.Field(i)
		structTag := structField.Tag
		name := structTag.Get("column")

		if _, ok := mapData[name]; !ok {
			continue
		}
		switch structField.Type.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
			if val, err := strconv.ParseUint(mapData[name], 10, 64); err == nil {
				fieldVal.SetUint(val)
			}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if val, err := strconv.ParseInt(mapData[name], 10, 64); err == nil {
				fieldVal.SetInt(val)
			}
		case reflect.String:
			fieldVal.SetString(mapData[name])
		case reflect.Bool:
			val := false
			if mapData[name] == "1" {
				val = true
			}
			fieldVal.SetBool(val)
		case reflect.Float32, reflect.Float64:
			if val, err := strconv.ParseFloat(mapData[name], 64); err == nil {
				fieldVal.SetFloat(val)
			}
		case reflect.Struct:
			if structField.Type.Name() == "Time" {
				local, _ := time.LoadLocation("Local")
				val, err := time.ParseInLocation("2006-01-02 15:04:05", mapData[name], local)
				if err == nil {
					fieldVal.Set(reflect.ValueOf(val))
				}
			}
		}
	}
	return rv.Interface(), err
}
