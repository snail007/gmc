package gcond

import (
	gcast "github.com/snail007/gmc/util/cast"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCond(t *testing.T) {
	// 测试 check 为 true 的情况
	result := Cond(true, "ok", "fail")
	if result.String() != "ok" {
		t.Errorf("Cond(true) = %v; want 'ok'", result)
	}

	// 测试 check 为 false 的情况
	result = Cond(false, "ok", "fail")
	if result.String() != "fail" {
		t.Errorf("Cond(false) = %v; want 'fail'", result)
	}
}

func TestCondFn(t *testing.T) {
	// 测试 check 为 true 的情况
	result := CondFn(true, func() interface{} { return "ok" }, func() interface{} { return "fail" })
	if result.String() != "ok" {
		t.Errorf("CondFn(true) = %v; want 'ok'", result)
	}

	// 测试 check 为 false 的情况
	result = CondFn(false, func() interface{} { return "ok" }, func() interface{} { return "fail" })
	if result.String() != "fail" {
		t.Errorf("CondFn(false) = %v; want 'fail'", result)
	}
}

func TestValue_xxx(t *testing.T) {
	var t0 interface{}
	t0 = gcast.ToTime("2000-01-01 13:00:00")
	t1 := NewValue(t0)
	assert.Equal(t, t0, t1.Time())

	t0 = time.Duration(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Duration())

	t0 = int(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Int())

	t0 = int32(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Int32())

	t0 = int64(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Int64())

	t0 = float32(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Float32())

	t0 = float64(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Float64())

	t0 = uint(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Uint())

	t0 = uint8(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Uint8())

	t0 = uint32(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Uint32())

	t0 = uint64(123)
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Uint64())

	t0 = true
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Bool())

	t0 = "123"
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.String())

	t0 = 123
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Val())

	t0 = map[string]interface{}{"123": "123"}
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.Map())

	t0 = map[string]string{"123": "123"}
	t1 = NewValue(t0)
	assert.Equal(t, t0, t1.MapSS())
}
