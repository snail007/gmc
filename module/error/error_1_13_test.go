package gerror

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIs(t *testing.T) {
	// 创建自定义错误
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	// 创建自定义错误的包装错误
	wrappedErr1 := &Error{Err: err1}
	wrappedErr2 := &Error{Err: err2}
	wrappedErr3 := &Error{Err: err3}

	// 测试基本错误与原始错误相等的情况
	result := Is(err1, err1)
	assert.True(t, result)

	// 测试基本错误与不同的原始错误不相等的情况
	result = Is(err1, err2)
	assert.False(t, result)

	// 测试包装错误与原始错误相等的情况
	result = Is(wrappedErr1, err1)
	assert.True(t, result)

	// 测试包装错误与不同的原始错误不相等的情况
	result = Is(wrappedErr1, err2)
	assert.False(t, result)

	// 测试包装错误与包装错误的原始错误相等的情况
	result = Is(wrappedErr1, wrappedErr1)
	assert.True(t, result)

	// 测试包装错误与不同的包装错误的原始错误不相等的情况
	result = Is(wrappedErr1, wrappedErr2)
	assert.False(t, result)

	// 原始错误为多层包装错误的情况
	result = Is(wrappedErr2, err2)
	assert.True(t, result)

	// 包装错误为多层包装错误的情况
	result = Is(wrappedErr2, wrappedErr2)
	assert.True(t, result)

	// 包装错误为多层包装错误但不相等的情况
	result = Is(wrappedErr2, wrappedErr3)
	assert.False(t, result)
}
