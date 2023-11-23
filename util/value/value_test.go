package gvalue

import (
	"github.com/snail007/gmc/util/cast"
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
