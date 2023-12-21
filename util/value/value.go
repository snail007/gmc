package gvalue

import (
	gbytes "github.com/snail007/gmc/util/bytes"
	gcast "github.com/snail007/gmc/util/cast"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"reflect"
	"strings"
	"time"
)

// Must ignore err, return *Must of v
func Must(v interface{}, err error) (value *Value) {
	if err != nil || IsNil(v) {
		return nil
	}
	return New(v)
}

// MustAny ignore err, return *AnyValue of v
func MustAny(v interface{}, err error) (value *AnyValue) {
	if err != nil || IsNil(v) {
		return nil
	}
	return NewAny(v)
}

// IsNil if object is nil or a nil pointer variable or nil builtin pointer variable(chan,func,interface,map,ptr,slice) return true, others false.
func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	isNilAbleKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice},
		kind)

	if isNilAbleKind && value.IsNil() {
		return true
	}
	return false
}

// IsEmpty if object is IsNil with true, or is empty string, return true, others false.
func IsEmpty(object interface{}) bool {
	if IsNil(object) {
		return true
	}
	if v, ok := object.(string); ok && v == "" {
		return true
	}
	return false
}

var numberPrinter = message.NewPrinter(language.English)

// FormatNumber return split by comma format of number v.
// v must be intX, uintX or intX string, uintX string,
// example: 10000 -> 10,000
func FormatNumber(v interface{}) string {
	switch val := v.(type) {
	case int, int8, int32, int64, uint, uint8, uint32, uint64:
		return numberPrinter.Sprintf("%d", val)
	case string:
		return numberPrinter.Sprintf("%d", gcast.ToInt64(val))
	}
	return ""
}

// ParseNumber convert human-readable number format v to int64,
// example: 100,123 -> 100123
func ParseNumber(v string) int64 {
	return gcast.ToInt64(strings.Replace(v, ",", "", -1))
}

// FormatByteSize convert byte size v to human-readable format, for example 1024 -> 1K
func FormatByteSize(v uint64) string {
	return Must(gbytes.SizeStr(v)).String()
}

// ParseByteSize convert human-readable byte size format v to byte size,
// v can be: 1, 1B, 1KB, 1.1KB, 1MB, 1GB, 1TB, 1PB, 1EB
func ParseByteSize(v string) uint64 {
	return Must(gbytes.ParseSize(v)).Uint64()
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
	val              interface{}
	cacheInt         *int
	cacheInt8        *int8
	cacheInt32       *int32
	cacheInt64       *int64
	cacheUint        *uint
	cacheUint8       *uint8
	cacheUint32      *uint32
	cacheUint64      *uint64
	cacheFloat32     *float32
	cacheFloat64     *float64
	cacheString      *string
	cacheStringSlice []string
	cacheBool        *bool
	cacheBytes       []byte
	cacheDuration    *time.Duration
	cacheTime        *time.Time

	cacheMap            map[string]interface{}
	cacheMapSlice       map[string][]interface{}
	cacheMapBool        map[string]bool
	cacheMapString      map[string]string
	cacheMapStringSlice map[string][]string
	cacheMapInt         map[string]int
	cacheMapInt8        map[string]int8
	cacheMapInt32       map[string]int32
	cacheMapInt64       map[string]int64
	cacheMapUint        map[string]uint
	cacheMapUint8       map[string]uint8
	cacheMapUint32      map[string]uint32
	cacheMapUint64      map[string]uint64
	cacheMapFloat32     map[string]float32
	cacheMapFloat64     map[string]float64
}

// New wrap the val with type *Value, type val and *Value.XXX() must be same.
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
	if s.cacheInt != nil {
		return *s.cacheInt
	}
	v := s.val.(int)
	s.cacheInt = &v
	return v
}

func (s *Value) Int8() int8 {
	if s.cacheInt8 != nil {
		return *s.cacheInt8
	}
	v := s.val.(int8)
	s.cacheInt8 = &v
	return v
}

func (s *Value) Int32() int32 {
	if s.cacheInt32 != nil {
		return *s.cacheInt32
	}
	v := s.val.(int32)
	s.cacheInt32 = &v
	return v
}

func (s *Value) Int64() int64 {
	if s.cacheInt64 != nil {
		return *s.cacheInt64
	}
	v := s.val.(int64)
	s.cacheInt64 = &v
	return v
}

func (s *Value) Uint() uint {
	if s.cacheUint != nil {
		return *s.cacheUint
	}
	v := s.val.(uint)
	s.cacheUint = &v
	return v
}
func (s *Value) Uint8() uint8 {
	if s.cacheUint8 != nil {
		return *s.cacheUint8
	}
	v := s.val.(uint8)
	s.cacheUint8 = &v
	return v
}

func (s *Value) Uint32() uint32 {
	if s.cacheUint32 != nil {
		return *s.cacheUint32
	}
	v := s.val.(uint32)
	s.cacheUint32 = &v
	return v
}

func (s *Value) Uint64() uint64 {
	if s.cacheUint64 != nil {
		return *s.cacheUint64
	}
	v := s.val.(uint64)
	s.cacheUint64 = &v
	return v
}

func (s *Value) Float32() float32 {
	if s.cacheFloat32 != nil {
		return *s.cacheFloat32
	}
	v := s.val.(float32)
	s.cacheFloat32 = &v
	return v
}

func (s *Value) Float64() float64 {
	if s.cacheFloat64 != nil {
		return *s.cacheFloat64
	}
	v := s.val.(float64)
	s.cacheFloat64 = &v
	return v
}

func (s *Value) Bool() bool {
	if s.cacheBool != nil {
		return *s.cacheBool
	}
	v := s.val.(bool)
	s.cacheBool = &v
	return v
}

func (s *Value) Bytes() []byte {
	if s.cacheBytes != nil {
		return s.cacheBytes
	}
	s.cacheBytes = s.val.([]byte)
	return s.cacheBytes
}

func (s *Value) String() string {
	if s.cacheString != nil {
		return *s.cacheString
	}
	v := s.val.(string)
	s.cacheString = &v
	return v
}

func (s *Value) StringSlice() []string {
	if s.cacheStringSlice != nil {
		return s.cacheStringSlice
	}
	if s.val == nil {
		return nil
	}
	v := s.val.([]string)
	s.cacheStringSlice = v
	return v
}

func (s *Value) Duration() time.Duration {
	if s.cacheDuration != nil {
		return *s.cacheDuration
	}
	v := s.val.(time.Duration)
	s.cacheDuration = &v
	return v
}

func (s *Value) Time() time.Time {
	if s.cacheTime != nil {
		return *s.cacheTime
	}
	v := s.val.(time.Time)
	s.cacheTime = &v
	return v
}

func (s *Value) Map() map[string]interface{} {
	if s.cacheMap != nil {
		return s.cacheMap
	}
	s.cacheMap = s.val.(map[string]interface{})
	return s.cacheMap
}

func (s *Value) MapString() map[string]string {
	if s.cacheMapString != nil {
		return s.cacheMapString
	}
	s.cacheMapString = s.val.(map[string]string)
	return s.cacheMapString
}

func (s *Value) MapBool() map[string]bool {
	if s.cacheMapBool != nil {
		return s.cacheMapBool
	}
	s.cacheMapBool = s.val.(map[string]bool)
	return s.cacheMapBool
}

func (s *Value) MapSlice() map[string][]interface{} {
	if s.cacheMapSlice != nil {
		return s.cacheMapSlice
	}
	s.cacheMapSlice = s.val.(map[string][]interface{})
	return s.cacheMapSlice
}

func (s *Value) MapStringSlice() map[string][]string {
	if s.cacheMapStringSlice != nil {
		return s.cacheMapStringSlice
	}
	s.cacheMapStringSlice = s.val.(map[string][]string)
	return s.cacheMapStringSlice
}

func (s *Value) MapInt() map[string]int {
	if s.cacheMapInt != nil {
		return s.cacheMapInt
	}
	s.cacheMapInt = s.val.(map[string]int)
	return s.cacheMapInt
}

func (s *Value) MapInt8() map[string]int8 {
	if s.cacheMapInt8 != nil {
		return s.cacheMapInt8
	}
	s.cacheMapInt8 = s.val.(map[string]int8)
	return s.cacheMapInt8
}

func (s *Value) MapInt32() map[string]int32 {
	if s.cacheMapInt32 != nil {
		return s.cacheMapInt32
	}
	s.cacheMapInt32 = s.val.(map[string]int32)
	return s.cacheMapInt32
}

func (s *Value) MapInt64() map[string]int64 {
	if s.cacheMapInt64 != nil {
		return s.cacheMapInt64
	}
	s.cacheMapInt64 = s.val.(map[string]int64)
	return s.cacheMapInt64
}

func (s *Value) MapUint() map[string]uint {
	if s.cacheMapUint != nil {
		return s.cacheMapUint
	}
	s.cacheMapUint = s.val.(map[string]uint)
	return s.cacheMapUint
}

func (s *Value) MapUint8() map[string]uint8 {
	if s.cacheMapUint8 != nil {
		return s.cacheMapUint8
	}
	s.cacheMapUint8 = s.val.(map[string]uint8)
	return s.cacheMapUint8
}

func (s *Value) MapUint32() map[string]uint32 {
	if s.cacheMapUint32 != nil {
		return s.cacheMapUint32
	}
	s.cacheMapUint32 = s.val.(map[string]uint32)
	return s.cacheMapUint32
}

func (s *Value) MapUint64() map[string]uint64 {
	if s.cacheMapUint64 != nil {
		return s.cacheMapUint64
	}
	s.cacheMapUint64 = s.val.(map[string]uint64)
	return s.cacheMapUint64
}

func (s *Value) MapFloat32() map[string]float32 {
	if s.cacheMapFloat32 != nil {
		return s.cacheMapFloat32
	}
	s.cacheMapFloat32 = s.val.(map[string]float32)
	return s.cacheMapFloat32
}

func (s *Value) MapFloat64() map[string]float64 {
	if s.cacheMapFloat64 != nil {
		return s.cacheMapFloat64
	}
	s.cacheMapFloat64 = s.val.(map[string]float64)
	return s.cacheMapFloat64
}

type AnyValue struct {
	val          interface{}
	cacheInt     *int
	cacheInt8    *int8
	cacheInt32   *int32
	cacheInt64   *int64
	cacheUint    *uint
	cacheUint8   *uint8
	cacheUint32  *uint32
	cacheUint64  *uint64
	cacheFloat32 *float32
	cacheFloat64 *float64

	cacheIntSlice     []int
	cacheInt8Slice    []int8
	cacheInt32Slice   []int32
	cacheInt64Slice   []int64
	cacheUintSlice    []uint
	cacheUint8Slice   []uint8
	cacheUint32Slice  []uint32
	cacheUint64Slice  []uint64
	cacheFloat32Slice []float32
	cacheFloat64Slice []float64

	cacheBool     *bool
	cacheDuration *time.Duration
	cacheTime     *time.Time
	cacheString   *string

	cacheBoolSlice     []bool
	cacheDurationSlice []time.Duration
	cacheTimeSlice     []time.Time
	cacheStringSlice   []string

	cacheMap        map[string]interface{}
	cacheMapBool    map[string]bool
	cacheMapString  map[string]string
	cacheMapInt     map[string]int
	cacheMapInt8    map[string]int8
	cacheMapInt32   map[string]int32
	cacheMapInt64   map[string]int64
	cacheMapUint    map[string]uint
	cacheMapUint8   map[string]uint8
	cacheMapUint32  map[string]uint32
	cacheMapUint64  map[string]uint64
	cacheMapFloat32 map[string]float32
	cacheMapFloat64 map[string]float64
}

// NewAny wrap the val with type *AnyValue, *AnyValue.XXX() convert val to type XXX.
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
	if s.cacheInt != nil {
		return *s.cacheInt
	}
	v := gcast.ToInt(s.val)
	s.cacheInt = &v
	return v
}

func (s *AnyValue) IntSlice() []int {
	if s.cacheIntSlice != nil {
		return s.cacheIntSlice
	}
	v := gcast.ToIntSlice(s.val)
	s.cacheIntSlice = v
	return v
}

func (s *AnyValue) Int8() int8 {
	if s.cacheInt8 != nil {
		return *s.cacheInt8
	}
	v := gcast.ToInt8(s.val)
	s.cacheInt8 = &v
	return v
}

func (s *AnyValue) Int8Slice() []int8 {
	if s.cacheInt8Slice != nil {
		return s.cacheInt8Slice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheInt8Slice = append(s.cacheInt8Slice, gcast.ToInt8(v))
	})
	return s.cacheInt8Slice
}

func (s *AnyValue) Int32() int32 {
	if s.cacheInt32 != nil {
		return *s.cacheInt32
	}
	v := gcast.ToInt32(s.val)
	s.cacheInt32 = &v
	return v
}

func (s *AnyValue) Int32Slice() []int32 {
	if s.cacheInt32Slice != nil {
		return s.cacheInt32Slice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheInt32Slice = append(s.cacheInt32Slice, gcast.ToInt32(v))
	})
	return s.cacheInt32Slice
}

func (s *AnyValue) Int64() int64 {
	if s.cacheInt64 != nil {
		return *s.cacheInt64
	}
	v := gcast.ToInt64(s.val)
	s.cacheInt64 = &v
	return v
}

func (s *AnyValue) Int64Slice() []int64 {
	if s.cacheInt64Slice != nil {
		return s.cacheInt64Slice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheInt64Slice = append(s.cacheInt64Slice, gcast.ToInt64(v))
	})
	return s.cacheInt64Slice
}

func (s *AnyValue) Uint() uint {
	if s.cacheUint != nil {
		return *s.cacheUint
	}
	v := gcast.ToUint(s.val)
	s.cacheUint = &v
	return v
}

func (s *AnyValue) UintSlice() []uint {
	if s.cacheUintSlice != nil {
		return s.cacheUintSlice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheUintSlice = append(s.cacheUintSlice, gcast.ToUint(v))
	})
	return s.cacheUintSlice
}

func (s *AnyValue) Uint8() uint8 {
	if s.cacheUint8 != nil {
		return *s.cacheUint8
	}
	v := gcast.ToUint8(s.val)
	s.cacheUint8 = &v
	return v
}

func (s *AnyValue) Uint8Slice() []uint8 {
	if s.cacheUint8Slice != nil {
		return s.cacheUint8Slice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheUint8Slice = append(s.cacheUint8Slice, gcast.ToUint8(v))
	})
	return s.cacheUint8Slice
}

func (s *AnyValue) Uint32() uint32 {
	if s.cacheUint32 != nil {
		return *s.cacheUint32
	}
	v := gcast.ToUint32(s.val)
	s.cacheUint32 = &v
	return v
}

func (s *AnyValue) Uint32Slice() []uint32 {
	if s.cacheUint32Slice != nil {
		return s.cacheUint32Slice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheUint32Slice = append(s.cacheUint32Slice, gcast.ToUint32(v))
	})
	return s.cacheUint32Slice
}

func (s *AnyValue) Uint64() uint64 {
	if s.cacheUint64 != nil {
		return *s.cacheUint64
	}
	v := gcast.ToUint64(s.val)
	s.cacheUint64 = &v
	return v
}

func (s *AnyValue) Uint64Slice() []uint64 {
	if s.cacheUint64Slice != nil {
		return s.cacheUint64Slice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheUint64Slice = append(s.cacheUint64Slice, gcast.ToUint64(v))
	})
	return s.cacheUint64Slice
}

func (s *AnyValue) Float32() float32 {
	if s.cacheFloat32 != nil {
		return *s.cacheFloat32
	}
	v := gcast.ToFloat32(s.val)
	s.cacheFloat32 = &v
	return v
}

func (s *AnyValue) Float32Slice() []float32 {
	if s.cacheFloat32Slice != nil {
		return s.cacheFloat32Slice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheFloat32Slice = append(s.cacheFloat32Slice, gcast.ToFloat32(v))
	})
	return s.cacheFloat32Slice
}

func (s *AnyValue) Float64() float64 {
	if s.cacheFloat64 != nil {
		return *s.cacheFloat64
	}
	v := gcast.ToFloat64(s.val)
	s.cacheFloat64 = &v
	return v
}

func (s *AnyValue) Float64Slice() []float64 {
	if s.cacheFloat64Slice != nil {
		return s.cacheFloat64Slice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheFloat64Slice = append(s.cacheFloat64Slice, gcast.ToFloat64(v))
	})
	return s.cacheFloat64Slice
}

func (s *AnyValue) Bool() bool {
	if s.cacheBool != nil {
		return *s.cacheBool
	}
	v := gcast.ToBool(s.val)
	s.cacheBool = &v
	return v
}

func (s *AnyValue) BoolSlice() []bool {
	if s.cacheBoolSlice != nil {
		return s.cacheBoolSlice
	}
	walkSlice(s.val, func(v interface{}) {
		s.cacheBoolSlice = append(s.cacheBoolSlice, gcast.ToBool(v))
	})
	return s.cacheBoolSlice
}

func (s *AnyValue) String() string {
	if s.cacheString != nil {
		return *s.cacheString
	}
	v := gcast.ToString(s.val)
	s.cacheString = &v
	return v
}

func (s *AnyValue) StringSlice() []string {
	if s.cacheStringSlice != nil {
		return s.cacheStringSlice
	}
	s.cacheStringSlice = gcast.ToStringSlice(s.val)
	return s.cacheStringSlice
}

func (s *AnyValue) Duration() time.Duration {
	if s.cacheDuration != nil {
		return *s.cacheDuration
	}
	v := gcast.ToDuration(s.val)
	s.cacheDuration = &v
	return v
}

func (s *AnyValue) DurationSlice() []time.Duration {
	if s.cacheDurationSlice != nil {
		return s.cacheDurationSlice
	}
	s.cacheDurationSlice = gcast.ToDurationSlice(s.val)
	return s.cacheDurationSlice
}

func (s *AnyValue) Time() time.Time {
	if s.cacheTime != nil {
		return *s.cacheTime
	}
	v := gcast.ToTime(s.val)
	s.cacheTime = &v
	return v
}

func (s *AnyValue) Map() map[string]interface{} {
	if s.cacheMap != nil {
		return s.cacheMap
	}
	m := map[string]interface{}{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = v
	})
	s.cacheMap = m
	return m
}

func (s *AnyValue) MapString() map[string]string {
	if s.cacheMapString != nil {
		return s.cacheMapString
	}
	s.cacheMapString = gcast.ToStringMapString(s.val)
	return s.cacheMapString
}

func (s *AnyValue) MapBool() map[string]bool {
	if s.cacheMapBool != nil {
		return s.cacheMapBool
	}
	s.cacheMapBool = gcast.ToStringMapBool(s.val)
	return s.cacheMapBool
}

func (s *AnyValue) MapInt() map[string]int {
	if s.cacheMapInt != nil {
		return s.cacheMapInt
	}
	s.cacheMapInt = gcast.ToStringMapInt(s.val)
	return s.cacheMapInt
}

func (s *AnyValue) MapInt8() map[string]int8 {
	if s.cacheMapInt8 != nil {
		return s.cacheMapInt8
	}
	m := map[string]int8{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = gcast.ToInt8(v)
	})
	s.cacheMapInt8 = m
	return m
}

func (s *AnyValue) MapInt32() map[string]int32 {
	if s.cacheMapInt32 != nil {
		return s.cacheMapInt32
	}
	m := map[string]int32{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = gcast.ToInt32(v)
	})
	s.cacheMapInt32 = m
	return m
}

func (s *AnyValue) MapInt64() map[string]int64 {
	if s.cacheMapInt64 != nil {
		return s.cacheMapInt64
	}
	s.cacheMapInt64 = gcast.ToStringMapInt64(s.val)
	return s.cacheMapInt64
}

func (s *AnyValue) MapUint() map[string]uint {
	if s.cacheMapUint != nil {
		return s.cacheMapUint
	}
	m := map[string]uint{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = gcast.ToUint(v)
	})
	s.cacheMapUint = m
	return m
}

func (s *AnyValue) MapUint8() map[string]uint8 {
	if s.cacheMapUint8 != nil {
		return s.cacheMapUint8
	}
	m := map[string]uint8{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = gcast.ToUint8(v)
	})
	s.cacheMapUint8 = m
	return m
}

func (s *AnyValue) MapUint32() map[string]uint32 {
	if s.cacheMapUint32 != nil {
		return s.cacheMapUint32
	}
	m := map[string]uint32{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = gcast.ToUint32(v)
	})
	s.cacheMapUint32 = m
	return m
}

func (s *AnyValue) MapUint64() map[string]uint64 {
	if s.cacheMapUint64 != nil {
		return s.cacheMapUint64
	}
	m := map[string]uint64{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = gcast.ToUint64(v)
	})
	s.cacheMapUint64 = m
	return m
}

func (s *AnyValue) MapFloat32() map[string]float32 {
	if s.cacheMapFloat32 != nil {
		return s.cacheMapFloat32
	}
	m := map[string]float32{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = gcast.ToFloat32(v)
	})
	s.cacheMapFloat32 = m
	return m
}

func (s *AnyValue) MapFloat64() map[string]float64 {
	if s.cacheMapFloat64 != nil {
		return s.cacheMapFloat64
	}
	m := map[string]float64{}
	walkMap(s.val, func(k string, v interface{}) {
		m[k] = gcast.ToFloat64(v)
	})
	s.cacheMapFloat64 = m
	return m
}

func walkMap(val interface{}, f func(k string, v interface{})) {
	value := reflect.ValueOf(val)
	if value.Kind() != reflect.Map {
		return
	}
	keys := value.MapKeys()
	for _, k := range keys {
		v := value.MapIndex(k)
		f(gcast.ToString(k.Interface()), v.Interface())
	}
}

func walkSlice(val interface{}, f func(v interface{})) {
	if IsNil(val) {
		return
	}
	kind := reflect.TypeOf(val).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(val)
		for j := 0; j < s.Len(); j++ {
			f(s.Index(j).Interface())
		}
	}
}

// GetValueAt sliceOrArray is slice or array, idx is the element to get, if idx invalid or sliceOrArray is not slice or array, the
// default value will be used
func GetValueAt(sliceOrArray interface{}, idx int, defaultValue interface{}) *Value {
	value := reflect.ValueOf(sliceOrArray)

	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return New(defaultValue)
	}

	if idx < 0 || idx >= value.Len() {
		return New(defaultValue)
	}

	return New(value.Index(idx).Interface())
}

// Keys input _map is a map, returns its keys as slice, AnyValue.xxxSlice to get the typed slice.
func Keys(_map interface{}) *AnyValue {
	value := reflect.ValueOf(_map)

	if value.Kind() != reflect.Map {
		return NewAny(nil)
	}

	keys := value.MapKeys()
	result := make([]interface{}, len(keys))

	for i, key := range keys {
		result[i] = key.Interface()
	}

	return NewAny(result)
}

// Contains check the input slice if contains the value
func Contains(sliceInterface interface{}, value interface{}) bool {
	return IndexOf(sliceInterface, value) >= 0
}

// IndexOf returns the index of value in slice, if not found -1 returned
func IndexOf(sliceInterface interface{}, value interface{}) int {
	sliceValue := reflect.ValueOf(sliceInterface)

	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Array {
		return -1
	}

	for i := 0; i < sliceValue.Len(); i++ {
		element := sliceValue.Index(i).Interface()
		if reflect.DeepEqual(element, value) {
			return i
		}
	}

	return -1
}
