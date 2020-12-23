package ghttputil

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	writeByteCnt int64
}

func NewResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	if _, ok := w.(*ResponseWriter); ok {
		return w
	}
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     200,
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
