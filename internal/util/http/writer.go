// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package ghttputil

import (
	"net/http"
	"sync"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	writeByteCnt int64
	data         *sync.Map
}

func (this *ResponseWriter) ClearData() {
	this.data = &sync.Map{}
}

func (this *ResponseWriter) Data(k interface{}) interface{} {
	v, ok := this.data.Load(k)
	if !ok {
		return nil
	}
	return v
}

func (this *ResponseWriter) SetData(k interface{}, v interface{}) {
	this.data.Store(k, v)
}

func NewResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	if _, ok := w.(*ResponseWriter); ok {
		return w
	}
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     200,
		data:           &sync.Map{},
	}
}

func (this *ResponseWriter) WriteHeader(status int) {
	this.statusCode = status
	this.ResponseWriter.WriteHeader(status)
}
func (this *ResponseWriter) Write(b []byte) (n int, err error) {
	n, err = this.ResponseWriter.Write(b)
	if n > 0 {
		this.writeByteCnt += int64(n)
	}
	return
}

//WriteCount acquires outgoing bytes count by writer
func (this *ResponseWriter) WriteCount() int64 {
	return this.writeByteCnt
}
func (this *ResponseWriter) StatusCode() int {
	return this.statusCode
}

func StatusCode(w http.ResponseWriter) (code int) {
	if v, ok := w.(*ResponseWriter); ok {
		code = v.statusCode
	}
	return
}
func WriteCount(w http.ResponseWriter) int64 {
	if v, ok := w.(*ResponseWriter); ok {
		return v.writeByteCnt
	}
	return 0
}
