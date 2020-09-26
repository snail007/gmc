package gmchttputil

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
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

func (this *ResponseWriter) StatusCode() int {
	return this.statusCode
}

func StatusCode(w http.ResponseWriter) (code int) {
	if v, ok := w.(*ResponseWriter); ok {
		code = v.statusCode
	}
	return
}
