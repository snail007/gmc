package grand

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
)

// String returns the length random string.
func String(length int) string {
	b := make([]byte, length/2+1)
	crand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// New returns a new rand.Rand object of safe source.
func New() *rand.Rand {
	b := make([]byte, 8)
	crand.Read(b)
	return rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(b))))
}

// IntString returns the length random string only contains number and starts with no zero.
// If length is 1, it returns 0-9 randomly, if length is 0, it returns empty.
func IntString(length int) string {
	if length == 0 {
		return ""
	}
	r := New()
	var rIntChar= func() string{
		return strconv.FormatInt(r.Int63n(10), 10)
	}
	if length==1{
		return rIntChar()
	}
	str := rIntChar()
	for str == "0" {
		str = rIntChar()
	}
	for i := 0; i < length-1; i++ {
		str+=rIntChar()
	}
	return str
}
