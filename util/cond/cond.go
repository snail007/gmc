package gcond

import (
	"time"
)

func Cond(check bool, ok interface{}, fail interface{}) *Value {
	if check {
		return NewValue(ok)
	}
	return NewValue(fail)
}

func CondFn(check bool, ok func() interface{}, fail func() interface{}) *Value {
	if check {
		return NewValue(ok())
	}
	return NewValue(fail())
}

type Value struct {
	val interface{}
}

func NewValue(val interface{}) *Value {
	return &Value{val: val}
}

func (s *Value) Val() interface{} {
	return s.val
}

func (s *Value) Int() int {
	return s.val.(int)
}

func (s *Value) Int32() int32 {
	return s.val.(int32)
}

func (s *Value) Int64() int64 {
	return s.val.(int64)
}

func (s *Value) Uint() uint {
	return s.val.(uint)
}
func (s *Value) Uint8() uint8 {
	return s.val.(uint8)
}

func (s *Value) Uint32() uint32 {
	return s.val.(uint32)
}

func (s *Value) Uint64() uint64 {
	return s.val.(uint64)
}

func (s *Value) Float32() float32 {
	return s.val.(float32)
}

func (s *Value) Float64() float64 {
	return s.val.(float64)
}

func (s *Value) Bool() bool {
	return s.val.(bool)
}

func (s *Value) String() string {
	return s.val.(string)
}

func (s *Value) Duration() time.Duration {
	return s.val.(time.Duration)
}

func (s *Value) Time() time.Time {
	return s.val.(time.Time)
}

func (s *Value) Map() map[string]interface{} {
	return s.val.(map[string]interface{})
}

func (s *Value) MapSS() map[string]string {
	return s.val.(map[string]string)
}
