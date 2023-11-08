package gvalue

import (
	gcast "github.com/snail007/gmc/util/cast"
	"reflect"
	"time"
)

func Must(v interface{}, err error) (value *Value) {
	if err != nil || IsNil(v) {
		return nil
	}
	return New(v)
}

func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	isNilableKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice},
		kind)

	if isNilableKind && value.IsNil() {
		return true
	}
	return false
}

func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for i := 0; i < len(kinds); i++ {
		if kind == kinds[i] {
			return true
		}
	}
	return false
}

type Value struct {
	val interface{}
}

func New(val interface{}) *Value {
	v := val
	if IsNil(val) {
		v = nil
	}
	return &Value{val: v}
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

func (s *Value) MapString() map[string]string {
	return s.val.(map[string]string)
}

func (s *Value) MapBool() map[string]bool {
	return s.val.(map[string]bool)
}

func (s *Value) MapSlice() map[string][]interface{} {
	return s.val.(map[string][]interface{})
}

func (s *Value) MapStringSlice() map[string][]string {
	return s.val.(map[string][]string)
}

func (s *Value) MapInt() map[string]int {
	return s.val.(map[string]int)
}

func (s *Value) MapInt64() map[string]int64 {
	return s.val.(map[string]int64)
}

func (s *Value) MapUint() map[string]uint {
	return s.val.(map[string]uint)
}

func (s *Value) MapUint8() map[string]uint8 {
	return s.val.(map[string]uint8)
}

func (s *Value) MapUint32() map[string]uint32 {
	return s.val.(map[string]uint32)
}

func (s *Value) MapUint64() map[string]uint64 {
	return s.val.(map[string]uint64)
}

func (s *Value) MapFloat32() map[string]float32 {
	return s.val.(map[string]float32)
}

func (s *Value) MapFloat64() map[string]float64 {
	return s.val.(map[string]float64)
}

type AnyValue struct {
	val interface{}
}

func NewAny(val interface{}) *AnyValue {
	v := val
	if IsNil(val) {
		v = nil
	}
	return &AnyValue{val: v}
}

func (s *AnyValue) Val() interface{} {
	return s.val
}

func (s *AnyValue) Int() int {
	return gcast.ToInt(s.val)
}

func (s *AnyValue) IntSlice() []int {
	return gcast.ToIntSlice(s.val)
}

func (s *AnyValue) Int32() int32 {
	return gcast.ToInt32(s.val)
}

func (s *AnyValue) Int32Slice() []int32 {
	var a []int32
	walkSlice(s.val, func(v string) {
		a = append(a, gcast.ToInt32(v))
	})
	return a
}

func (s *AnyValue) Int64() int64 {
	return gcast.ToInt64(s.val)
}

func (s *AnyValue) Int64Slice() []int64 {
	var a []int64
	walkSlice(s.val, func(v string) {
		a = append(a, gcast.ToInt64(v))
	})
	return a
}

func (s *AnyValue) Uint() uint {
	return gcast.ToUint(s.val)
}

func (s *AnyValue) UintSlice() []uint {
	var a []uint
	walkSlice(s.val, func(v string) {
		a = append(a, gcast.ToUint(v))
	})
	return a
}

func (s *AnyValue) Uint8() uint8 {
	return gcast.ToUint8(s.val)
}

func (s *AnyValue) Uint8Slice() []uint8 {
	var a []uint8
	walkSlice(s.val, func(v string) {
		a = append(a, gcast.ToUint8(v))
	})
	return a
}

func (s *AnyValue) Uint32() uint32 {
	return gcast.ToUint32(s.val)
}

func (s *AnyValue) Uint32Slice() []uint32 {
	var a []uint32
	walkSlice(s.val, func(v string) {
		a = append(a, gcast.ToUint32(v))
	})
	return a
}

func (s *AnyValue) Uint64() uint64 {
	return gcast.ToUint64(s.val)
}

func (s *AnyValue) Uint64Slice() []uint64 {
	var a []uint64
	walkSlice(s.val, func(v string) {
		a = append(a, gcast.ToUint64(v))
	})
	return a
}

func (s *AnyValue) Float32() float32 {
	return gcast.ToFloat32(s.val)
}

func (s *AnyValue) Float32Slice() []float32 {
	var a []float32
	walkSlice(s.val, func(v string) {
		a = append(a, gcast.ToFloat32(v))
	})
	return a
}

func (s *AnyValue) Float64() float64 {
	return gcast.ToFloat64(s.val)
}

func (s *AnyValue) Float64Slice() []float64 {
	var a []float64
	walkSlice(s.val, func(v string) {
		a = append(a, gcast.ToFloat64(v))
	})
	return a
}

func (s *AnyValue) Bool() bool {
	return gcast.ToBool(s.val)
}

func (s *AnyValue) BoolSlice() []bool {
	return gcast.ToBoolSlice(s.val)
}

func (s *AnyValue) String() string {
	return gcast.ToString(s.val)
}

func (s *AnyValue) StringSlice() []string {
	return gcast.ToStringSlice(s.val)
}

func (s *AnyValue) Duration() time.Duration {
	return gcast.ToDuration(s.val)
}

func (s *AnyValue) DurationSlice() []time.Duration {
	return gcast.ToDurationSlice(s.val)
}

func (s *AnyValue) Time() time.Time {
	return gcast.ToTime(s.val)
}

func (s *AnyValue) Map() map[string]interface{} {
	m := map[string]interface{}{}
	walkMap(s.val, func(k, v string) {
		m[k] = v
	})
	return m
}

func (s *AnyValue) MapString() map[string]string {
	return gcast.ToStringMapString(s.val)
}

func (s *AnyValue) MapBool() map[string]bool {
	return gcast.ToStringMapBool(s.val)
}

func (s *AnyValue) MapInt() map[string]int {
	return gcast.ToStringMapInt(s.val)
}

func (s *AnyValue) MapInt32() map[string]int32 {
	m := map[string]int32{}
	walkMap(s.val, func(k, v string) {
		m[k] = gcast.ToInt32(v)
	})
	return m
}

func (s *AnyValue) MapInt64() map[string]int64 {
	return gcast.ToStringMapInt64(s.val)
}

func (s *AnyValue) MapUint() map[string]uint {
	m := map[string]uint{}
	walkMap(s.val, func(k, v string) {
		m[k] = gcast.ToUint(v)
	})
	return m
}

func (s *AnyValue) MapUint8() map[string]uint8 {
	m := map[string]uint8{}
	walkMap(s.val, func(k, v string) {
		m[k] = gcast.ToUint8(v)
	})
	return m
}

func (s *AnyValue) MapUint32() map[string]uint32 {
	m := map[string]uint32{}
	walkMap(s.val, func(k, v string) {
		m[k] = gcast.ToUint32(v)
	})
	return m
}

func (s *AnyValue) MapUint64() map[string]uint64 {
	m := map[string]uint64{}
	walkMap(s.val, func(k, v string) {
		m[k] = gcast.ToUint64(v)
	})
	return m
}

func (s *AnyValue) MapFloat32() map[string]float32 {
	m := map[string]float32{}
	walkMap(s.val, func(k, v string) {
		m[k] = gcast.ToFloat32(v)
	})
	return m
}

func (s *AnyValue) MapFloat64() map[string]float64 {
	m := map[string]float64{}
	walkMap(s.val, func(k, v string) {
		m[k] = gcast.ToFloat64(v)
	})
	return m
}

func walkMap(val interface{}, f func(k, v string)) {
	for k, v := range gcast.ToStringMapString(val) {
		f(k, v)
	}
}

func walkSlice(val interface{}, f func(v string)) {
	for _, v := range gcast.ToStringSlice(val) {
		f(v)
	}
}
