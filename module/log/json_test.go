// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

// 辅助函数：解析 JSON 日志为 Map
func parseLogEntry(t *testing.T, logStr string) map[string]interface{} {
	var entry map[string]interface{}
	err := json.Unmarshal([]byte(logStr), &entry)
	assert.NoError(t, err)
	return entry
}

// 测试日志级别输出
func TestLogLevels(t *testing.T) {
	logger := NewJSONLogger()

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{}) string
		expected string
	}{
		{"Trace", logger.Tracef, "trace"},
		{"Debug", logger.Debugf, "debug"},
		{"Info", logger.Infof, "info"},
		{"Warn", logger.Warnf, "warn"},
		{"Error", logger.Errorf, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logStr := tt.logFunc("test message")
			entry := parseLogEntry(t, logStr)
			assert.Equal(t, strings.ToUpper(tt.expected), entry["log_level"])
			assert.Contains(t, entry, "log_time")
		})
	}
}

// 测试字段类型序列化
func TestFieldTypes(t *testing.T) {
	logger := NewJSONLogger()
	testTime := time.Now()

	logStr := logger.
		Uint8("uint8_val", 10).
		Int("int_val", -5).
		Float32("float32_val", 3.14).
		String("string_val", "gopher").
		Bool("bool_val", true).
		Time("time_val", testTime).
		Duration("duration_val", 5*time.Second).
		Any("any_val", map[string]int{"key": 1}).
		Infof("fields test")

	entry := parseLogEntry(t, logStr)

	// 验证基础类型
	assert.Equal(t, float64(10), entry["uint8_val"]) // JSON 数字统一为 float64
	assert.Equal(t, float64(-5), entry["int_val"])
	assert.InEpsilon(t, 3.14, entry["float32_val"], 0.001)
	assert.Equal(t, "gopher", entry["string_val"])
	assert.Equal(t, true, entry["bool_val"])

	// 验证时间格式 (RFC3339 含毫秒)
	_, err := time.Parse("2006-01-02 15:04:05.000 -07", entry["time_val"].(string))
	assert.NoError(t, err)

	// 验证复杂类型
	assert.Equal(t, "5s", entry["duration_val"])
	assert.Equal(t, map[string]interface{}{"key": float64(1)}, entry["any_val"])
}

func NewJSONLogger() *JSONLogger {
	return NewLogger().JSON()
}

// 测试 Panic 日志
func TestPanicLog(t *testing.T) {
	logger := NewJSONLogger()

	defer func() {
		if r := recover(); r != nil {
			logStr := r.(string)
			entry := parseLogEntry(t, logStr)
			assert.Equal(t, "PANIC", entry["log_level"])
		} else {
			t.Fatal("Expected panic did not occur")
		}
	}()

	logger.Panic("critical error")
}

// 测试日志格式字符串
func TestFormattedLog(t *testing.T) {
	logger := NewJSONLogger()
	logStr := logger.Infof("User %s logged in from %s", "alice", "192.168.1.1")

	entry := parseLogEntry(t, logStr)
	assert.Equal(t, "User alice logged in from 192.168.1.1", entry["log_msg"])
}

// 测试必填字段
func TestRequiredFields(t *testing.T) {
	logger := NewJSONLogger()
	logStr := logger.Infof("system check")

	entry := parseLogEntry(t, logStr)
	assert.Contains(t, entry, "log_time")
	assert.Equal(t, "INFO", entry["log_level"])
	assert.Equal(t, "system check", entry["log_msg"])
}

// 测试并发安全 (通过 -race 参数检测)
func TestConcurrentLogging(t *testing.T) {
	logger := NewLogger()
	const goroutines = 10

	done := make(chan bool)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				logger.JSON().Int("goroutine", id).
					Int("iteration", j).
					Debugf("concurrent log %d", j)
			}
			done <- true
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}
}
