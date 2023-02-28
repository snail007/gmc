package gjson

import (
	"encoding/json"
	"io"
)

type JSONResult struct {
	data map[string]interface{}
}

func NewResult() *JSONResult {
	return &JSONResult{
		data: map[string]interface{}{
			"code":    0,
			"message": "",
			"data":    nil,
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

func (s *JSONResult) Data(msg string) *JSONResult {
	return s.Set("message", msg)
}

func (s *JSONResult) WriteTo(dst io.Writer) (err error) {
	_, err = dst.Write(s.ToJSON())
	return
}
