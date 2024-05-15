package gcore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogLevel_String(t *testing.T) {
	l := LogLevelTrace
	assert.Equal(t, "TRACE", l.String())
	l = LogLeveDebug
	assert.Equal(t, "DEBUG", l.String())
	l = LogLeveInfo
	assert.Equal(t, "INFO", l.String())
	l = LogLeveWarn
	assert.Equal(t, "WARN", l.String())
	l = LogLeveError
	assert.Equal(t, "ERROR", l.String())
	l = LogLevePanic
	assert.Equal(t, "PANIC", l.String())
	l = LogLeveFatal
	assert.Equal(t, "FATAL", l.String())
	l = LogLeveNone
	assert.Equal(t, "NONE", l.String())
}
