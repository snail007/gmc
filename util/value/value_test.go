package gvalue

import (
	"encoding/base64"
	"github.com/snail007/gmc/util/cast"
	gmap "github.com/snail007/gmc/util/map"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestIsNil(t *testing.T) {
	var a *http.Client
	var b *http.Client
	var c interface{}
	var d chan bool
	var e map[string]string
	var f []byte
	var g *int
	b = nil
	assert.True(t, IsNil(a))
	assert.True(t, IsNil(b))
	assert.True(t, IsNil(c))
	assert.True(t, IsNil(d))
	assert.True(t, IsNil(e))
	assert.True(t, IsNil(f))
	assert.True(t, IsNil(g))
	assert.True(t, IsNil(nil))
	assert.False(t, IsNil(123))
	assert.False(t, IsNil(1.1))
	assert.False(t, IsNil(http.Client{}))
	assert.False(t, IsNil(map[string]string{}))
	assert.False(t, IsNil(new(int)))
	assert.False(t, IsNil(new(http.Client)))
}

func TestValue_xxx(t *testing.T) {
	var t0 interface{}
	t0 = gcast.ToTime("2000-01-01 13:00:00")
	t1 := New(t0)
	assert.Equal(t, t0, t1.Time())
	assert.NotNil(t, t1.cacheTime)
	assert.Equal(t, t0, t1.Time())

	t0 = time.Duration(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Duration())
	assert.NotNil(t, t1.cacheDuration)
	assert.Equal(t, t0, t1.Duration())

	t0 = int(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Int())
	assert.NotNil(t, t1.cacheInt)
	assert.Equal(t, t0, t1.Int())

	t0 = int8(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Int8())
	assert.NotNil(t, t1.cacheInt8)
	assert.Equal(t, t0, t1.Int8())

	t0 = int32(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Int32())
	assert.NotNil(t, t1.cacheInt32)
	assert.Equal(t, t0, t1.Int32())

	t0 = int64(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Int64())
	assert.NotNil(t, t1.cacheInt64)
	assert.Equal(t, t0, t1.Int64())

	t0 = float32(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Float32())
	assert.NotNil(t, t1.cacheFloat32)
	assert.Equal(t, t0, t1.Float32())

	t0 = float64(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Float64())
	assert.NotNil(t, t1.cacheFloat64)
	assert.Equal(t, t0, t1.Float64())

	t0 = uint(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Uint())
	assert.NotNil(t, t1.cacheUint)
	assert.Equal(t, t0, t1.Uint())

	t0 = uint8(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Uint8())
	assert.NotNil(t, t1.cacheUint8)
	assert.Equal(t, t0, t1.Uint8())

	t0 = uint32(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Uint32())
	assert.NotNil(t, t1.cacheUint32)
	assert.Equal(t, t0, t1.Uint32())

	t0 = uint64(123)
	t1 = New(t0)
	assert.Equal(t, t0, t1.Uint64())
	assert.NotNil(t, t1.cacheUint64)
	assert.Equal(t, t0, t1.Uint64())

	t0 = true
	t1 = New(t0)
	assert.Equal(t, t0, t1.Bool())
	assert.NotNil(t, t1.cacheBool)
	assert.Equal(t, t0, t1.Bool())

	t0 = "123"
	t1 = New(t0)
	assert.Equal(t, t0, t1.String())
	assert.NotNil(t, t1.cacheString)
	assert.Equal(t, t0, t1.String())

	t0 = []string{"123"}
	t1 = New(t0)
	assert.Equal(t, t0, t1.StringSlice())
	assert.NotNil(t, t1.cacheStringSlice)
	assert.Equal(t, t0, t1.StringSlice())

	t0 = []byte("123")
	t1 = New(t0)
	assert.Equal(t, t0, t1.Bytes())
	assert.NotNil(t, t1.cacheBytes)
	assert.Equal(t, t0, t1.Bytes())

	t0 = 123
	t1 = New(t0)
	assert.Equal(t, t0, t1.Val())

	t0 = nil
	t1 = New(t0)
	assert.Nil(t, t1.Val())

	t0 = map[string]interface{}{"123": "123"}
	t1 = New(t0)
	assert.Equal(t, t0, t1.Map())
	assert.NotNil(t, t1.cacheMap)
	assert.Equal(t, t0, t1.Map())

	t0 = map[string]string{"123": "123"}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapString())
	assert.NotNil(t, t1.cacheMapString)
	assert.Equal(t, t0, t1.MapString())

	t0 = map[string]bool{"123": true}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapBool())
	assert.NotNil(t, t1.MapBool())
	assert.Equal(t, t0, t1.MapBool())

	t0 = map[string]int{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapInt())
	assert.NotNil(t, t1.cacheMapInt)
	assert.Equal(t, t0, t1.MapInt())

	t0 = map[string]int8{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapInt8())
	assert.NotNil(t, t1.cacheMapInt8)
	assert.Equal(t, t0, t1.MapInt8())

	t0 = map[string]int32{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapInt32())
	assert.NotNil(t, t1.cacheMapInt32)
	assert.Equal(t, t0, t1.MapInt32())

	t0 = map[string]int64{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapInt64())
	assert.NotNil(t, t1.cacheMapInt64)
	assert.Equal(t, t0, t1.MapInt64())

	t0 = map[string]uint{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapUint())
	assert.NotNil(t, t1.cacheMapUint)
	assert.Equal(t, t0, t1.MapUint())

	t0 = map[string]uint8{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapUint8())
	assert.NotNil(t, t1.MapUint8())
	assert.Equal(t, t0, t1.MapUint8())

	t0 = map[string]uint32{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapUint32())
	assert.NotNil(t, t1.MapUint32())
	assert.Equal(t, t0, t1.MapUint32())

	t0 = map[string]uint64{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapUint64())
	assert.NotNil(t, t1.MapUint64())
	assert.Equal(t, t0, t1.MapUint64())

	t0 = map[string]float32{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapFloat32())
	assert.NotNil(t, t1.cacheMapFloat32)
	assert.Equal(t, t0, t1.MapFloat32())

	t0 = map[string]float64{"123": 123}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapFloat64())
	assert.NotNil(t, t1.cacheMapFloat64)
	assert.Equal(t, t0, t1.MapFloat64())

	t0 = map[string][]string{"123": {"123"}}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapStringSlice())
	assert.NotNil(t, t1.cacheMapStringSlice)
	assert.Equal(t, t0, t1.MapStringSlice())

	t0 = map[string][]interface{}{"123": {"123"}}
	t1 = New(t0)
	assert.Equal(t, t0, t1.MapSlice())
	assert.NotNil(t, t1.cacheMapSlice)
	assert.Equal(t, t0, t1.MapSlice())
}

func TestValue_xxx2(t *testing.T) {
	val := New(nil)
	assert.Equal(t, val.Val(), nil)
	assert.Equal(t, val.Int(), 0)
	assert.Equal(t, val.Int8(), int8(0))
	assert.Equal(t, val.Int32(), int32(0))
	assert.Equal(t, val.Int64(), int64(0))
	assert.Equal(t, val.Uint(), uint(0))
	assert.Equal(t, val.Uint8(), uint8(0))
	assert.Equal(t, val.Uint32(), uint32(0))
	assert.Equal(t, val.Uint64(), uint64(0))
	assert.Equal(t, val.Float32(), float32(0.0))
	assert.Equal(t, val.Float64(), 0.0)
	assert.Equal(t, val.Bool(), false)
	assert.Equal(t, val.String(), "")
	assert.Equal(t, val.Duration(), time.Duration(0))
	assert.Equal(t, val.Time(), time.Time{})
	assert.Nil(t, val.Bytes())
	assert.Nil(t, val.MapSlice())
	assert.Nil(t, val.MapStringSlice())
	assert.Nil(t, val.StringSlice())
	assert.Nil(t, val.MapUint64())
	assert.Nil(t, val.MapUint32())
	assert.Nil(t, val.MapUint8())
	assert.Nil(t, val.MapUint())
	assert.Nil(t, val.MapInt64())
	assert.Nil(t, val.MapInt32())
	assert.Nil(t, val.MapInt8())
	assert.Nil(t, val.MapInt())
	assert.Nil(t, val.MapString())
	assert.Nil(t, val.MapBool())
	assert.Nil(t, val.Map())
	assert.Nil(t, val.MapFloat64())
	assert.Nil(t, val.MapFloat32())
}

func TestAnyValue_xxx(t *testing.T) {
	val := NewAny(nil)
	assert.Equal(t, val.Val(), nil)
	assert.Equal(t, val.Int(), 0)
	assert.Equal(t, val.Int8(), int8(0))
	assert.Equal(t, val.Int32(), int32(0))
	assert.Equal(t, val.Int64(), int64(0))
	assert.Equal(t, val.Uint(), uint(0))
	assert.Equal(t, val.Uint8(), uint8(0))
	assert.Equal(t, val.Uint32(), uint32(0))
	assert.Equal(t, val.Uint64(), uint64(0))
	assert.Equal(t, val.Float32(), float32(0.0))
	assert.Equal(t, val.Float64(), 0.0)
	assert.Equal(t, val.Bool(), false)
	assert.Equal(t, val.String(), "")
	assert.Equal(t, val.Duration(), time.Duration(0))
	assert.Equal(t, val.Time(), time.Time{})
	assert.Nil(t, val.Bytes())
	assert.Nil(t, val.StringSlice())
	assert.Nil(t, val.MapUint64())
	assert.Nil(t, val.MapUint32())
	assert.Nil(t, val.MapUint8())
	assert.Nil(t, val.MapUint())
	assert.Nil(t, val.MapInt64())
	assert.Nil(t, val.MapInt32())
	assert.Nil(t, val.MapInt8())
	assert.Nil(t, val.MapInt())
	assert.Nil(t, val.MapString())
	assert.Nil(t, val.MapBool())
	assert.Nil(t, val.Map())
	assert.Nil(t, val.MapFloat64())
	assert.Nil(t, val.MapFloat32())
	assert.Nil(t, val.IntSlice())
	assert.Nil(t, val.Int8Slice())
	assert.Nil(t, val.Int32Slice())
	assert.Nil(t, val.Int64Slice())
	assert.Nil(t, val.UintSlice())
	assert.Nil(t, val.Uint8Slice())
	assert.Nil(t, val.Uint32Slice())
	assert.Nil(t, val.Uint64Slice())
	assert.Nil(t, val.Float32Slice())
	assert.Nil(t, val.Float64Slice())
	assert.Nil(t, val.BoolSlice())
	assert.Nil(t, val.DurationSlice())
}

func TestMust(t *testing.T) {
	assert.Equal(t, 123, Must(gcast.ToIntE("123")).Int())
	a, e := gcast.ToIntE("123abc")
	assert.Equal(t, 0, a)
	assert.NotNil(t, e)
	assert.Nil(t, Must(a, e))
}

func TestAnyValue(t *testing.T) {
	// Test NewAny method
	val := NewAny(10)
	assert.Equal(t, val.Val(), int(10))
	assert.Equal(t, val.Int(), int(10))
	assert.Equal(t, val.Int8(), int8(10))
	assert.Equal(t, val.Int32(), int32(10))
	assert.Equal(t, val.Int64(), int64(10))
	assert.Equal(t, val.Uint(), uint(10))
	assert.Equal(t, val.Uint8(), uint8(10))
	assert.Equal(t, val.Uint32(), uint32(10))
	assert.Equal(t, val.Uint64(), uint64(10))
	assert.Equal(t, val.Float32(), float32(10.0))
	assert.Equal(t, val.Float64(), float64(10.0))
	assert.Equal(t, val.Bool(), true)
	assert.Equal(t, val.String(), "10")
	assert.Equal(t, val.Duration(), time.Duration(10))
	assert.Equal(t, val.Time(), time.Unix(10, 0))
	assert.Equal(t, val.Bytes(), []byte("10"))

	assert.Equal(t, val.Val(), int(10))
	assert.Equal(t, val.Int(), int(10))
	assert.Equal(t, val.Int8(), int8(10))
	assert.Equal(t, val.Int32(), int32(10))
	assert.Equal(t, val.Int64(), int64(10))
	assert.Equal(t, val.Uint(), uint(10))
	assert.Equal(t, val.Uint8(), uint8(10))
	assert.Equal(t, val.Uint32(), uint32(10))
	assert.Equal(t, val.Uint64(), uint64(10))
	assert.Equal(t, val.Float32(), float32(10.0))
	assert.Equal(t, val.Float64(), float64(10.0))
	assert.Equal(t, val.Bool(), true)
	assert.Equal(t, val.String(), "10")
	assert.Equal(t, val.Duration(), time.Duration(10))
	assert.Equal(t, val.Time(), time.Unix(10, 0))
	assert.Equal(t, val.Bytes(), []byte("10"))

	val = NewAny([]int{10})
	assert.Equal(t, val.IntSlice(), []int{10})
	assert.Equal(t, val.Int8Slice(), []int8{10})
	assert.Equal(t, val.Int32Slice(), []int32{10})
	assert.Equal(t, val.Int64Slice(), []int64{10})
	assert.Equal(t, val.UintSlice(), []uint{10})
	assert.Equal(t, val.Uint8Slice(), []uint8{10})
	assert.Equal(t, val.Uint32Slice(), []uint32{10})
	assert.Equal(t, val.Uint64Slice(), []uint64{10})
	assert.Equal(t, val.Float32Slice(), []float32{float32(10.0)})
	assert.Equal(t, val.Float64Slice(), []float64{float64(10.0)})
	assert.Equal(t, val.BoolSlice(), []bool{true})
	assert.Equal(t, val.StringSlice(), []string{"10"})
	assert.Equal(t, val.DurationSlice(), []time.Duration{time.Duration(10)})

	assert.Equal(t, val.IntSlice(), []int{10})
	assert.Equal(t, val.Int8Slice(), []int8{10})
	assert.Equal(t, val.Int32Slice(), []int32{10})
	assert.Equal(t, val.Int64Slice(), []int64{10})
	assert.Equal(t, val.UintSlice(), []uint{10})
	assert.Equal(t, val.Uint8Slice(), []uint8{10})
	assert.Equal(t, val.Uint32Slice(), []uint32{10})
	assert.Equal(t, val.Uint64Slice(), []uint64{10})
	assert.Equal(t, val.Float32Slice(), []float32{float32(10.0)})
	assert.Equal(t, val.Float64Slice(), []float64{float64(10.0)})
	assert.Equal(t, val.BoolSlice(), []bool{true})
	assert.Equal(t, val.StringSlice(), []string{"10"})
	assert.Equal(t, val.DurationSlice(), []time.Duration{time.Duration(10)})

	val = NewAny(map[string]interface{}{"a": "10"})
	assert.Equal(t, val.MapString(), map[string]string{"a": "10"})
	assert.Equal(t, val.Map(), map[string]interface{}{"a": "10"})
	assert.Equal(t, val.MapBool(), map[string]bool{"a": false})
	assert.Equal(t, val.MapInt(), map[string]int{"a": 10})
	assert.Equal(t, val.MapInt8(), map[string]int8{"a": 10})
	assert.Equal(t, val.MapInt32(), map[string]int32{"a": 10})
	assert.Equal(t, val.MapInt64(), map[string]int64{"a": 10})
	assert.Equal(t, val.MapUint(), map[string]uint{"a": 10})
	assert.Equal(t, val.MapUint8(), map[string]uint8{"a": 10})
	assert.Equal(t, val.MapUint32(), map[string]uint32{"a": 10})
	assert.Equal(t, val.MapUint64(), map[string]uint64{"a": 10})
	assert.Equal(t, val.MapFloat32(), map[string]float32{"a": 10})
	assert.Equal(t, val.MapFloat64(), map[string]float64{"a": 10})

	assert.Equal(t, val.MapString(), map[string]string{"a": "10"})
	assert.Equal(t, val.Map(), map[string]interface{}{"a": "10"})
	assert.Equal(t, val.MapBool(), map[string]bool{"a": false})
	assert.Equal(t, val.MapInt(), map[string]int{"a": 10})
	assert.Equal(t, val.MapInt8(), map[string]int8{"a": 10})
	assert.Equal(t, val.MapInt32(), map[string]int32{"a": 10})
	assert.Equal(t, val.MapInt64(), map[string]int64{"a": 10})
	assert.Equal(t, val.MapUint(), map[string]uint{"a": 10})
	assert.Equal(t, val.MapUint8(), map[string]uint8{"a": 10})
	assert.Equal(t, val.MapUint32(), map[string]uint32{"a": 10})
	assert.Equal(t, val.MapUint64(), map[string]uint64{"a": 10})
	assert.Equal(t, val.MapFloat32(), map[string]float32{"a": 10})
	assert.Equal(t, val.MapFloat64(), map[string]float64{"a": 10})

	val = NewAny(nil)
	assert.Nil(t, val.Val())

	assert.Nil(t, val.Val())
}

func TestMustAny(t *testing.T) {
	assert.Equal(t, 123, MustAny(gcast.ToUint64E("123")).Int())
	val := NewAny(nil)
	assert.Nil(t, val.Val())
	a, e := gcast.ToIntE("123abc")
	assert.Equal(t, 0, a)
	assert.NotNil(t, e)
	assert.Nil(t, MustAny(a, e))
}

func TestIsEmpty(t *testing.T) {
	var a *http.Client
	assert.True(t, IsEmpty(nil))
	assert.True(t, IsEmpty(""))
	assert.True(t, IsEmpty(a))
	assert.False(t, IsEmpty(false))
	assert.False(t, IsEmpty(0))
	b := 0
	walkSlice(nil, func(v interface{}) {
		b = 1
	})
	assert.Equal(t, 0, b)
	walkMap(nil, func(k string, v interface{}) {
		b = 2
	})
	assert.Equal(t, 0, b)
}

func TestPrettyNumber(t *testing.T) {
	assert.Equal(t, "123,456", FormatNumber(123456))
	assert.Equal(t, "-123,456", FormatNumber(-123456))
	assert.Equal(t, "123,456", FormatNumber("123456"))
	assert.Empty(t, FormatNumber(nil))
}

func TestPrettyByteSize(t *testing.T) {
	assert.Equal(t, "0 B", FormatByteSize(0))
}

func TestParseByteSize(t *testing.T) {
	assert.Equal(t, uint64(512), ParseByteSize("0.5K"))
	assert.Equal(t, uint64(1024), ParseByteSize("1K"))
}

func TestParseNumber(t *testing.T) {
	assert.Equal(t, int64(123000), ParseNumber("123,000"))
}

func TestGetValueAt(t *testing.T) {
	var a interface{}
	a = 123
	assert.Equal(t, 456, GetValueAt(a, 1, 456).Int())
	a = []int{1, 2, 3}
	assert.Equal(t, 2, GetValueAt(a, 1, 3).Int())
	assert.Equal(t, 3, GetValueAt(a, 100, 3).Int())
	assert.Equal(t, 3, GetValueAt(a, -1, 3).Int())

	a = []int{1, 2, 3}
	assert.Equal(t, 2, GetValueAt(a, 1, 3).Val())
}

func TestKeys(t *testing.T) {
	assert.Equal(t, nil, Keys(nil).Val())
	assert.Equal(t, []string(nil), Keys(nil).StringSlice())
	var a interface{}
	a = 123
	assert.Equal(t, nil, Keys(a).Val())
	assert.Equal(t, []string(nil), Keys(a).StringSlice())
	a = map[string]string{"a": "1"}
	assert.Equal(t, []string{"a"}, Keys(a).StringSlice())
	a = map[string]string{}
	assert.Equal(t, []string(nil), Keys(a).StringSlice())
}

func TestContains(t *testing.T) {
	assert.Equal(t, false, Contains(nil, 1))
	var a interface{}
	a = 123
	assert.Equal(t, false, Contains(a, 1))
	a = []int{1, 2, 3}
	assert.Equal(t, true, Contains(a, 1))
	assert.Equal(t, 2, IndexOf(a, 3))
	a = [3]int{1, 2, 3}
	assert.Equal(t, true, Contains(a, 1))
	assert.Equal(t, 2, IndexOf(a, 3))
	a = []interface{}{1, 2, 3, nil}
	assert.Equal(t, 3, IndexOf(a, nil))
	assert.Equal(t, -1, IndexOf(a, 0))
}

type YourStruct struct {
	Uint8Field     uint8
	Uint16Field    uint16
	Uint32Field    uint32
	Uint64Field    uint64
	UintField      uint
	Int8Field      int8
	Int16Field     int16
	Int32Field     int32
	Int64Field     int64
	IntField       int
	StringField    string
	BytesField     []byte
	RawBytesField  []byte
	BoolField      bool
	Float32Field   float32
	Float64Field   float64
	TimeField      time.Time
	Int64TimeField time.Time
	ErrTimeField   time.Time
	StructField    YourInnerStruct
	PtrStructField *YourInnerPtrStruct
	ErrStructField YourInnerStruct
	MapKey         int8 `mkey:"key1"`
	privateField   int
	//map field only  map[string]interface{} supported
	MapField    map[string]interface{}
	IgnoreField int `mkey:"-"`
}

type YourInnerStruct struct {
	InnerIntField int
	InnerStrField string
}

type YourInnerPtrStruct struct {
	InnerIntField int
	InnerStrField string
}

type YourSecondStruct struct {
	AnotherIntField int
	AnotherStrField string
}

func TestMapToStruct(t *testing.T) {
	mapData := map[string]interface{}{
		"Uint8Field":     uint8(8),
		"Uint16Field":    uint16(16),
		"Uint32Field":    uint32(32),
		"Uint64Field":    uint64(64),
		"UintField":      uint(128),
		"Int8Field":      int8(-8),
		"Int16Field":     int16(-16),
		"Int32Field":     int32(-32),
		"Int64Field":     int64(-64),
		"IntField":       int(-128),
		"StringField":    "test",
		"BytesField":     base64.StdEncoding.EncodeToString([]byte("base64")),
		"RawBytesField":  "base64",
		"BoolField":      true,
		"Float32Field":   float32(3.14),
		"Float64Field":   float64(6.28),
		"TimeField":      time.Now().Format(time.RFC3339),
		"Int64TimeField": 123,
		"StructField": map[string]interface{}{
			"InnerIntField": 42,
			"InnerStrField": "nested",
		},
		"PtrStructField": map[string]interface{}{
			"InnerIntField": 42,
			"InnerStrField": "nested",
		},
		"key1":         int8(8),
		"privateField": 1,
		"MapField":     map[string]string{"abc": "123"},
		"IgnoreField":  1,
	}

	for _, v := range []interface{}{YourStruct{}, new(YourStruct)} {
		result, err := MapToStruct(mapData, v, "mkey")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		yourStruct := result.(YourStruct)
		// 现在，您可以断言各个字段
		assert.Equal(t, uint8(8), yourStruct.Uint8Field)
		assert.Equal(t, uint16(16), yourStruct.Uint16Field)
		assert.Equal(t, uint32(32), yourStruct.Uint32Field)
		assert.Equal(t, uint64(64), yourStruct.Uint64Field)
		assert.Equal(t, uint(128), yourStruct.UintField)
		assert.Equal(t, int8(-8), yourStruct.Int8Field)
		assert.Equal(t, int16(-16), yourStruct.Int16Field)
		assert.Equal(t, int32(-32), yourStruct.Int32Field)
		assert.Equal(t, int64(-64), yourStruct.Int64Field)
		assert.Equal(t, int(-128), yourStruct.IntField)
		assert.Equal(t, "test", yourStruct.StringField)
		assert.Equal(t, []byte("base64"), yourStruct.BytesField)
		assert.Equal(t, []byte("base64"), yourStruct.RawBytesField)
		assert.Equal(t, true, yourStruct.BoolField)
		assert.Equal(t, float32(3.14), yourStruct.Float32Field)
		assert.Equal(t, float64(6.28), yourStruct.Float64Field)
		assert.Equal(t, int64(123), yourStruct.Int64TimeField.Unix())
		assert.Equal(t, 42, yourStruct.StructField.InnerIntField)
		assert.Equal(t, "nested", yourStruct.StructField.InnerStrField)
		assert.Equal(t, 42, yourStruct.PtrStructField.InnerIntField)
		assert.Equal(t, "nested", yourStruct.PtrStructField.InnerStrField)
		assert.Equal(t, 0, yourStruct.privateField)
		assert.Equal(t, "123", yourStruct.MapField["abc"])
		assert.Equal(t, 0, yourStruct.IgnoreField)
	}
	// 测试第二个参数是不同的结构体的情况
	result, err := MapToStruct(mapData, YourSecondStruct{})
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = MapToStruct(mapData, nil)
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = MapToStruct(gmap.M{"ErrTimeField": "abc"}, new(YourStruct))
	assert.Error(t, err)
	assert.Nil(t, result)

	result, err = MapToStruct(gmap.M{"ErrStructField": new(chan bool)}, new(YourStruct))
	assert.Error(t, err)
	assert.Nil(t, result)
}
