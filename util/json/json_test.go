// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gjson

import (
	"bytes"
	"fmt"
	gctx "github.com/snail007/gmc/module/ctx"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResult_Parse(t *testing.T) {
	j := NewResult(`{"code":1,"data":"1","message":"2"}`)
	assert.NotNil(t, j)
	assert.Equal(t, 1, j.Code())
	assert.Equal(t, "1", j.Data())
	assert.Equal(t, "2", j.Message())
	j = NewResult([]byte(`{"code":1,"data":"1","message":"2"}`))
	assert.NotNil(t, j)
	assert.Equal(t, 1, j.Code())
	assert.Equal(t, "1", j.Data())
	assert.Equal(t, "2", j.Message())
	j = NewResult(1, 2, 3)
	assert.Equal(t, 1, j.Code())
	assert.Equal(t, "2", j.Message())
	assert.Equal(t, 3, j.Data())
	j = NewResult([]byte(`{`))
	assert.Nil(t, j)
}

func TestJSONResult_NewResult(t *testing.T) {
	result := NewResult()
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Code())
	assert.Equal(t, "", result.Message())
	assert.Nil(t, result.Data())
}

func TestJSONResult_NewResult_WithArgs(t *testing.T) {
	result := NewResult(200, "OK", "data")
	assert.NotNil(t, result)
	assert.Equal(t, 200, result.Code())
	assert.Equal(t, "OK", result.Message())
	assert.Equal(t, "data", result.Data())
}

func TestJSONResult_NewResult_WithInvalidData(t *testing.T) {
	result := NewResult([]byte("invalid"))
	assert.Nil(t, result)
}

func TestJSONResult_Set(t *testing.T) {
	result := NewResult().Set("key", "value")
	assert.NotNil(t, result)
	assert.Equal(t, "value", result.DataMap().(map[string]interface{})["key"])
}

func TestJSONResult_ToJSON(t *testing.T) {
	result := NewResult(200, "OK", "data")
	jsonData := result.ToJSON()

	expected := []byte(`{"code":200,"data":"data","message":"OK"}`)
	assert.Equal(t, expected, jsonData)
}

func TestJSONResult_WriteTo(t *testing.T) {
	result := NewResult(200, "OK", "data")
	buffer := bytes.NewBuffer(nil)
	err := result.WriteTo(buffer)

	assert.Nil(t, err)

	jsonData := buffer.Bytes()

	expected := []byte(`{"code":200,"data":"data","message":"OK"}`)
	assert.Equal(t, expected, jsonData)
}

func TestJSONResult_SetCode(t *testing.T) {
	result := NewResult()
	result.SetCode(200)
	assert.Equal(t, 200, result.Code())
}

func TestJSONResult_SetMessage(t *testing.T) {
	result := NewResult()
	result.SetMessage("Success")
	assert.Equal(t, "Success", result.Message())
}

func TestJSONResult_SetData(t *testing.T) {
	result := NewResult()
	result.SetData("data")
	assert.Equal(t, "data", result.Data())
}

func TestJSONResult_WriteToCtx(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.NewCtx().CloneWithHTTP(w, r)
		result := NewResult(200, "OK", "data")
		result.WriteToCtx(ctx)
	})

	handler.ServeHTTP(recorder, req)

	expected := `{"code":200,"data":"data","message":"OK"}`
	assert.Equal(t, expected, recorder.Body.String())
}

func TestJSONResult_Success(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.NewCtxWithHTTP(w, r)
		result := NewResultCtx(ctx)
		result.Success("success")
	})

	handler.ServeHTTP(recorder, req)

	expected := `{"code":0,"data":"success","message":""}`
	assert.Equal(t, expected, recorder.Body.String())
}

func TestJSONResult_Fail(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	recorder := httptest.NewRecorder()

	// 处理请求
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.NewCtxWithHTTP(w, r)
		result := NewResultCtx(ctx)
		result.Fail("fail")
	})

	handler.ServeHTTP(recorder, req)

	expected := `{"code":1,"data":null,"message":"fail"}`
	assert.Equal(t, expected, recorder.Body.String())
}

func TestBuilderOperations(t *testing.T) {
	builder := NewBuilder(`{"name": "John", "age": 30, "city": "New York"}`)

	err := builder.Set("name", "Doe")
	if err != nil {
		t.Errorf("Set method failed: %v", err)
	}

	err = builder.Delete("age")
	if err != nil {
		t.Errorf("Delete method failed: %v", err)
	}

	err = builder.SetRaw("country", `"USA"`)
	if err != nil {
		t.Errorf("SetRaw method failed: %v", err)
	}

	err = builder.SetRaw("country", `USA`)
	if err == nil {
		t.Errorf("SetRaw method failed: %v", err)
	}

	result := builder.Get("name")
	if result.String() != "Doe" {
		t.Errorf("Get method failed, expected 'Doe', got '%s'", result.String())
	}
	assert.Equal(t, result.Path(), "name")

	results := builder.GetMany("name", "country")
	if len(results) != 2 {
		t.Errorf("GetMany method failed, expected 2 results, got %d", len(results))
	}
	if results[1].String() != "USA" {
		t.Errorf("GetMany method failed, expected 'USA', got '%s'", results[1].String())
	}
	assert.Nil(t, result.Paths())
}

func TestJSONArray_Append(t *testing.T) {
	arr := NewJSONArray("[123]")
	assert.Equal(t, "123", arr.Get("0").String())

	obj := NewJSONObject(map[string]string{"name": "456"})
	arr.Append(obj)
	assert.Equal(t, "456", arr.Get("1.name").String())
	assert.Equal(t, "456", arr.Get("1").ToJSONObject().Get("name").String())

	obj = NewJSONObject(nil)
	obj.Set("name", "789")
	arr.Append(*obj)
	assert.Equal(t, "789", arr.Get("2.name").String())

	obja := NewJSONArray(nil)
	obja.Append("000", "111")
	arr.Append(obja)
	assert.Equal(t, "000", arr.Get("3.0").String())
	assert.Equal(t, "111", arr.Get("3.1").String())
	assert.Equal(t, "000", arr.Get("3").ToJSONArray().Get("0").String())

	assert.Equal(t, int64(4), arr.Len())

	obja = NewJSONArray([]string{"0000", "1111"})
	arr.Append(*obja)
	assert.Equal(t, "0000", arr.Get("4.0").String())
	assert.Equal(t, "1111", arr.Get("4.1").String())
	assert.Equal(t, "0000", arr.Get("4").ToJSONArray().Get("0").String())

	obj = NewJSONObject(`{"name":"111"}`)
	arr.Append(obj)
	assert.Equal(t, "111", arr.Get("5.name").String())

	obj = NewJSONObject([]byte(`{"name":"222"}`))
	arr.Append(obj)
	assert.Equal(t, "222", arr.Get("6.name").String())

	assert.Nil(t, NewJSONObject("{,abc"))
	assert.Nil(t, NewJSONArray("{,abc"))
	assert.Nil(t, NewBuilder("{,abc"))
}

func TestJSONArray_Merge(t *testing.T) {
	a := NewJSONArray([]int{123})
	arr := NewJSONArray(nil)
	assert.Nil(t, arr.Merge(a))
	assert.Nil(t, arr.Merge(*a))
	assert.Nil(t, arr.Merge([]string{"abc", "111"}))
	assert.Equal(t, int64(123), arr.Get("0").Int())
	assert.Equal(t, int64(123), arr.Get("1").Int())
	assert.Equal(t, "abc", arr.Get("2").String())
	assert.Equal(t, "111", arr.Get("3").String())
}

func TestBuilder_AsJSONObject(t *testing.T) {
	a := NewBuilder(`[]`)
	assert.Nil(t, a.JSONObject())
	assert.NotNil(t, a.JSONArray())
	assert.Equal(t, "[]", a.String())
	assert.Error(t, a.JSONArray().Append(http.Client{}))
	a = NewBuilder(`{}`)
	assert.Nil(t, a.JSONArray())
	assert.NotNil(t, a.JSONObject())
	assert.Equal(t, "{}", a.String())
}

func TestJSONArray_Last(t *testing.T) {
	a := NewJSONArray(nil)
	assert.False(t, a.First().Exists())
	assert.False(t, a.Last().Exists())
	assert.Empty(t, a.First().String())
	assert.Empty(t, a.Last().String())

	a.Append("123")
	assert.Equal(t, "123", a.First().String())
	assert.Equal(t, "123", a.Last().String())

	a.Append("456")
	assert.Equal(t, "123", a.First().String())
	assert.Equal(t, "456", a.Last().String())

}

func TestNewJSONObjectE(t *testing.T) {
	_, err := NewJSONObjectE([]string{})
	assert.Error(t, err)
	_, err = NewJSONArrayE(map[string]string{})
	assert.Error(t, err)
}

func TestBuilder_Interface(t *testing.T) {
	a := NewBuilder(nil)
	assert.Nil(t, a.Interface())

	a = NewBuilder("")
	assert.Nil(t, a)

	a = NewBuilder(`{"a":"b"}`)
	assert.Equal(t, map[string]interface{}{"a": "b"}, a.Interface())
	a = NewBuilder(map[string]interface{}{"a": "b"})
	assert.Equal(t, map[string]interface{}{"a": "b"}, a.Interface())

	a = NewBuilder(`[123]`)
	assert.Equal(t, []interface{}{float64(123)}, a.Interface())
	a = NewBuilder([]interface{}{123})
	assert.Equal(t, []interface{}{float64(123)}, a.Interface())

	r := Parse(`{"a":"b"}`)
	assert.Equal(t, "b", r.Get("a").String())

}
