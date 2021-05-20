package gos

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTempFile(t *testing.T) {
	assert.Contains(t, TempFile("", ".abc"), ".abc")
	assert.Contains(t, TempFile(".def", ""), ".def")
	f := TempFile(".def", ".abc")
	assert.Contains(t, f, ".def")
	assert.Contains(t, f, ".abc")
}
