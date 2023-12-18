package gstrings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytesRef(t *testing.T) {
	a := "abc"
	b := BytesRef(a)
	assert.Equal(t, []byte("abc"), b)
}

func TestStringRef(t *testing.T) {
	a := []byte("abc")
	b := StringRef(a)
	assert.Equal(t, "abc", b)
}

func TestHasPrefixAny(t *testing.T) {
	assert.True(t, HasPrefixAny("http://", "https://", "http://"))
	assert.False(t, HasPrefixAny("ftp://", "https://", "http://"))
}

func TestHasHTTPPrefix(t *testing.T) {
	assert.True(t, HasHTTPPrefix("http://"))
	assert.False(t, HasHTTPPrefix("ftp://"))
}
