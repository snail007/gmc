package gcond

import (
	"testing"
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
