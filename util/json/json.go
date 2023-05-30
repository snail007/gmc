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
	if len(d) == 1 {
		var b []byte
		switch d[0].(type) {
		case string:
			b = []byte(d[0].(string))
		case []byte:
			b = d[0].([]byte)
		}
		if len(b) > 0 {
			s := &JSONResult{
				data: map[string]interface{}{},
			}
			if e := json.Unmarshal(b, &s.data); e != nil {
				return nil
			}
			return s
		}
	}
	code := 0
	message := ""
	var data interface{}
	if len(d) >= 1 {
		code = gcast.ToInt(d[0])
	}
	if len(d) >= 2 {
		message = gcast.ToString(d[1])
	}
	if len(d) >= 3 {
		data = d[2]
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

func (s *JSONResult) SetCode(code int) *JSONResult {
	return s.Set("code", code)
}

func (s *JSONResult) SetMessage(format string, msg ...interface{}) *JSONResult {
	return s.Set("message", fmt.Sprintf(format, msg...))
}

func (s *JSONResult) SetData(d interface{}) *JSONResult {
	return s.Set("data", d)
}

func (s *JSONResult) Code() int {
	return gcast.ToInt(s.data["code"])
}

func (s *JSONResult) Message() string {
	return gcast.ToString(s.data["message"])
}

func (s *JSONResult) Data() interface{} {
	return s.data["data"]
}

func (s *JSONResult) DataMap() interface{} {
	return s.data
}

func (s *JSONResult) WriteTo(dst io.Writer) (err error) {
	_, err = dst.Write(s.ToJSON())
	return
}

func (s *JSONResult) WriteToCtx(ctx gcore.Ctx) (err error) {
	_, err = ctx.Response().Write(s.ToJSON())
	return
}

//Success only worked with NewResultCtx()
func (s *JSONResult) Success(d ...interface{}) (err error) {
	var data interface{}
	if len(d) == 1 {
		data = d[0]
	}
	return s.SetData(data).WriteToCtx(s.ctx)
}

//Fail only worked with NewResultCtx()
func (s *JSONResult) Fail(format string, v ...interface{}) (err error) {
	return s.SetCode(1).SetMessage(format, v...).WriteToCtx(s.ctx)
}
