// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gdb

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type a struct {
	B map[string]interface{} `json:"B"`
}
type c struct {
	B B `json:"B"`
}

type B struct {
	C string `json:"C"`
}

type cc struct {
	B *B `json:"B"`
}

func TestResultSet_mapToStruct(t *testing.T) {
	m := map[string]string{"B": `{"C":{"1":"1"}}`}
	rs := ResultSet{}
	r, e := rs.mapToStruct(m, a{})
	assert.Nil(t, e)
	assert.Equal(t, "1", r.(a).B["C"].(map[string]interface{})["1"])
}

func TestResultSet_mapToStruct1(t *testing.T) {
	m := map[string]string{"B": `{"C":"1"}`}
	rs := ResultSet{}
	r, e := rs.mapToStruct(m, c{})
	assert.Nil(t, e)
	assert.Equal(t, "1", r.(c).B.C)
}

func TestResultSet_mapToStruct2(t *testing.T) {
	m := map[string]string{"B": `{"C":"1"}`}
	rs := ResultSet{}
	r, e := rs.mapToStruct(m, cc{})
	assert.Nil(t, e)
	assert.Equal(t, "1", r.(cc).B.C)
}
func TestNewResultSet(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
		{"id": []byte("2"), "name": []byte("Jane")},
	}
	rs := NewResultSet(&rawRows)

	if rs.rawRows != &rawRows {
		t.Error("Expected rawRows to be set")
	}
}

func TestResultSet_SQL(t *testing.T) {
	rawRows := []map[string][]byte{}
	rs := NewResultSet(&rawRows)
	rs.sql = "SELECT * FROM users"

	sql := rs.SQL()

	if sql != "SELECT * FROM users" {
		t.Errorf("Expected SQL to be 'SELECT * FROM users', got '%s'", sql)
	}
}

func TestResultSet_Len(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
		{"id": []byte("2"), "name": []byte("Jane")},
	}
	rs := NewResultSet(&rawRows)

	length := rs.Len()

	if length != 2 {
		t.Errorf("Expected length to be 2, got %d", length)
	}
}

func TestResultSet_LastInsertID(t *testing.T) {
	rs := &ResultSet{lastInsertID: 10}

	lastInsertID := rs.LastInsertID()

	if lastInsertID != 10 {
		t.Errorf("Expected LastInsertID to be 10, got %d", lastInsertID)
	}
}

func TestResultSet_RowsAffected(t *testing.T) {
	rs := &ResultSet{rowsAffected: 5}

	rowsAffected := rs.RowsAffected()

	if rowsAffected != 5 {
		t.Errorf("Expected RowsAffected to be 5, got %d", rowsAffected)
	}
}

func TestResultSet_TimeUsed(t *testing.T) {
	rs := &ResultSet{timeUsed: 100}

	timeUsed := rs.TimeUsed()

	if timeUsed != 100 {
		t.Errorf("Expected TimeUsed to be 100, got %d", timeUsed)
	}
}

func TestResultSet_MapRows(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
		{"id": []byte("2"), "name": []byte("Jane")},
	}
	rs := NewResultSet(&rawRows)

	rowsMap := rs.MapRows("id")

	if len(rowsMap) != 2 {
		t.Errorf("Expected rowsMap length to be 2, got %d", len(rowsMap))
	}

	if rowsMap["1"]["name"] != "John" {
		t.Error("Expected rowsMap['1']['name'] to be 'John'")
	}

	if rowsMap["2"]["name"] != "Jane" {
		t.Error("Expected rowsMap['2']['name'] to be 'Jane'")
	}
}

func TestResultSet_MapStructs(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
		{"id": []byte("2"), "name": []byte("Jane")},
	}
	rs := NewResultSet(&rawRows)

	type Person struct {
		ID   int    `column:"id"`
		Name string `column:"name"`
	}

	structsMap, err := rs.MapStructs("id", Person{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(structsMap) != 2 {
		t.Errorf("Expected structsMap length to be 2, got %d", len(structsMap))
	}

	person1, ok := structsMap["1"].(Person)
	if !ok {
		t.Error("Expected person1 to be of type Person")
	}

	if person1.ID != 1 {
		t.Error("Expected person1.ID to be 1")
	}

	if person1.Name != "John" {
		t.Error("Expected person1.Name to be 'John'")
	}

	person2, ok := structsMap["2"].(Person)
	if !ok {
		t.Error("Expected person2 to be of type Person")
	}

	if person2.ID != 2 {
		t.Error("Expected person2.ID to be 2")
	}

	if person2.Name != "Jane" {
		t.Error("Expected person2.Name to be 'Jane'")
	}
}

func TestResultSet_Rows(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
		{"id": []byte("2"), "name": []byte("Jane")},
	}
	rs := NewResultSet(&rawRows)

	rows := rs.Rows()

	if len(rows) != 2 {
		t.Errorf("Expected rows length to be 2, got %d", len(rows))
	}

	if rows[0]["name"] != "John" {
		t.Error("Expected rows[0]['name'] to be 'John'")
	}

	if rows[1]["name"] != "Jane" {
		t.Error("Expected rows[1]['name'] to be 'Jane'")
	}
}

func TestResultSet_Structs(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
		{"id": []byte("2"), "name": []byte("Jane")},
	}
	rs := NewResultSet(&rawRows)

	type Person struct {
		ID   int    `column:"id"`
		Name string `column:"name"`
	}

	structs, err := rs.Structs(Person{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(structs) != 2 {
		t.Errorf("Expected structs length to be 2, got %d", len(structs))
	}

	person1, ok := structs[0].(Person)
	if !ok {
		t.Error("Expected person1 to be of type Person")
	}

	if person1.ID != 1 {
		t.Error("Expected person1.ID to be 1")
	}

	if person1.Name != "John" {
		t.Error("Expected person1.Name to be 'John'")
	}

	person2, ok := structs[1].(Person)
	if !ok {
		t.Error("Expected person2 to be of type Person")
	}

	if person2.ID != 2 {
		t.Error("Expected person2.ID to be 2")
	}

	if person2.Name != "Jane" {
		t.Error("Expected person2.Name to be 'Jane'")
	}
}

func TestResultSet_Row(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
	}
	rs := NewResultSet(&rawRows)

	row := rs.Row()

	if len(row) != 2 {
		t.Errorf("Expected row length to be 2, got %d", len(row))
	}

	if row["name"] != "John" {
		t.Error("Expected row['name'] to be 'John'")
	}
}

func TestResultSet_Struct(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
	}
	rs := NewResultSet(&rawRows)

	type Person struct {
		ID   int    `column:"id"`
		Name string `column:"name"`
	}

	person, err := rs.Struct(Person{})

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	p, ok := person.(Person)
	if !ok {
		t.Error("Expected person to be of type Person")
	}

	if p.ID != 1 {
		t.Error("Expected person.ID to be 1")
	}

	if p.Name != "John" {
		t.Error("Expected person.Name to be 'John'")
	}
}

func TestResultSet_Values(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
		{"id": []byte("2"), "name": []byte("Jane")},
	}
	rs := NewResultSet(&rawRows)

	values := rs.Values("name")

	if len(values) != 2 {
		t.Errorf("Expected values length to be 2, got %d", len(values))
	}

	if values[0] != "John" {
		t.Error("Expected values[0] to be 'John'")
	}

	if values[1] != "Jane" {
		t.Error("Expected values[1] to be 'Jane'")
	}
}

func TestResultSet_MapValues(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
		{"id": []byte("2"), "name": []byte("Jane")},
	}
	rs := NewResultSet(&rawRows)

	values := rs.MapValues("id", "name")

	if len(values) != 2 {
		t.Errorf("Expected values length to be 2, got %d", len(values))
	}

	if values["1"] != "John" {
		t.Error("Expected values['1'] to be 'John'")
	}

	if values["2"] != "Jane" {
		t.Error("Expected values['2'] to be 'Jane'")
	}
}

func TestResultSet_Value(t *testing.T) {
	rawRows := []map[string][]byte{
		{"id": []byte("1"), "name": []byte("John")},
	}
	rs := NewResultSet(&rawRows)

	value := rs.Value("name")

	if value != "John" {
		t.Error("Expected value to be 'John'")
	}
}

type TestStruct0 struct {
	ID          int       `column:"id"`
	Name        string    `column:"name"`
	Age         int       `column:"age"`
	Score       float64   `column:"score"`
	IsPassed    bool      `column:"is_passed"`
	CreatedAt   time.Time `column:"created_at"`
	Description string    `column:"description"`
	Bytes       []byte    `column:"bytes"`
}

func TestMapToStruct(t *testing.T) {
	mapData := map[string]string{
		"id":          "1",
		"name":        "John",
		"age":         "30",
		"score":       "9.8",
		"is_passed":   "true",
		"created_at":  "2023-05-29T10:00:00Z",
		"description": "Lorem ipsum",
		"bytes":       "SGVsbG8gd29ybGQ=",
	}
	rs := &ResultSet{}
	result, err := rs.mapToStruct(mapData, TestStruct0{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := TestStruct0{
		ID:          1,
		Name:        "John",
		Age:         30,
		Score:       9.8,
		IsPassed:    true,
		CreatedAt:   time.Date(2023, 5, 29, 10, 0, 0, 0, time.UTC),
		Description: "Lorem ipsum",
		Bytes:       []byte("SGVsbG8gd29ybGQ="),
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result. Expected: %v, but got: %v", expected, result)
	}
}
