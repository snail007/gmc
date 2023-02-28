package gjson

import (
	"encoding/json"
	"fmt"
	gcore "github.com/snail007/gmc/core"
	gcast "github.com/snail007/gmc/util/cast"
	"io"
)

type JSONResult struct {
	data map[string]interface{}
	ctx  gcore.Ctx
}

//NewResultCtx Optional args: code int, message string, data interface{}
func NewResultCtx(ctx gcore.Ctx, d ...interface{}) *JSONResult {
	r := NewResult(d...)
	r.ctx = ctx
	return r
}

//NewResult Optional args: code int, message string, data interface{}
func NewResult(d ...interface{}) *JSONResult {
	code := 0
	message := ""
	var data interface{}
	if len(d) >= 1 {
		code = gcast.ToInt(d[0])
	}
	if len(d) >= 2 {
		message = gcast.ToString(d[1])
	}
	if len(d) >= 1 {
		data = gcast.ToInt(d[2])
	}
	return &JSONResult{
		data: map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
	}
}

func (s *JSONResult) Set(key string, value interface{}) *JSONResult {
	s.data[key] = value
	return s
}

func (s *JSONResult) ToJSON() []byte {
	j, _ := json.Marshal(s.data)
	return j
}

func (s *JSONResult) Code(code int) *JSONResult {
	return s.Set("code", code)
}

func (s *JSONResult) Message(format string, msg ...interface{}) *JSONResult {
	return s.Set("message", fmt.Sprintf(format, msg...))
}

func (s *JSONResult) Data(d interface{}) *JSONResult {
	return s.Set("data", d)
}

func (s *JSONResult) WriteTo(dst io.Writer) (err error) {
	_, err = dst.Write(s.ToJSON())
	return
}

func (s *JSONResult) WriteToCtx(ctx gcore.Ctx) (err error) {
	_, err = ctx.Response().Write(s.ToJSON())
	return
}

//Success worked with NewResultCtx()
func (s *JSONResult) Success(d ...interface{}) (err error) {
	var data interface{}
	if len(d) == 1 {
		data = d[0]
	}
	return s.Data(data).WriteToCtx(s.ctx)
}

//Fail worked with NewResultCtx()
func (s *JSONResult) Fail(format string, v ...interface{}) (err error) {
	return s.Code(1).Message(format, v...).WriteToCtx(s.ctx)
}
