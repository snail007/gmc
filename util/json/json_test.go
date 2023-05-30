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
	// 创建一个模拟的 HTTP 请求
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	// 创建一个 ResponseRecorder 来记录响应
	recorder := httptest.NewRecorder()

	// 处理请求
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.NewCtx().CloneWithHTTP(w, r)
		result := NewResult(200, "OK", "data")
		result.WriteToCtx(ctx)
	})

	// 将请求发送到处理器
	handler.ServeHTTP(recorder, req)

	expected := `{"code":200,"data":"data","message":"OK"}`
	assert.Equal(t, expected, recorder.Body.String())
}

func TestJSONResult_Success(t *testing.T) {
	// 创建一个模拟的 HTTP 请求
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	// 创建一个 ResponseRecorder 来记录响应
	recorder := httptest.NewRecorder()

	// 处理请求
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.NewCtxWithHTTP(w, r)
		result := NewResultCtx(ctx)
		result.Success("success")
	})

	// 将请求发送到处理器
	handler.ServeHTTP(recorder, req)

	expected := `{"code":0,"data":"success","message":""}`
	assert.Equal(t, expected, recorder.Body.String())
}

func TestJSONResult_Fail(t *testing.T) {
	// 创建一个模拟的 HTTP 请求
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	// 创建一个 ResponseRecorder 来记录响应
	recorder := httptest.NewRecorder()

	// 处理请求
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.NewCtxWithHTTP(w, r)
		result := NewResultCtx(ctx)
		result.Fail("fail")
	})

	// 将请求发送到处理器
	handler.ServeHTTP(recorder, req)

	expected := `{"code":1,"data":null,"message":"fail"}`
	assert.Equal(t, expected, recorder.Body.String())
}
