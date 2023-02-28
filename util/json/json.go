package gjson

import (
	"encoding/json"
	gcast "github.com/snail007/gmc/util/cast"
	"io"
)

type JSONResult struct {
	data map[string]interface{}
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

func (s *JSONResult) Message(msg string) *JSONResult {
	return s.Set("message", msg)
}

func (s *JSONResult) Data(msg interface{}) *JSONResult {
	return s.Set("message", msg)
}

func (s *JSONResult) WriteTo(dst io.Writer) (err error) {
	_, err = dst.Write(s.ToJSON())
	return
}
