package strings

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
