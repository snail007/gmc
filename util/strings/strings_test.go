package gstrings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytesRef(t *testing.T) {
	a := "abc"
	b := StringToBytes(a)
	assert.Equal(t, []byte("abc"), b)
}

func TestStringRef(t *testing.T) {
	a := []byte("abc")
	b := BytesToString(a)
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

func TestHasSuffix(t *testing.T) {
	assert.True(t, HasSuffixAny("a.txt", ".log", ".txt"))
	assert.False(t, HasSuffixAny("a.log", ".log1", ".txt1"))
}
