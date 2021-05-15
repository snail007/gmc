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
	var rIntChar= func(max int64,add int64) string{
		return strconv.FormatInt(r.Int63n(max)+add, 10)
	}
	if length==1{
		return rIntChar(10,0)
	}
	str := rIntChar(9,1)
	for i := 0; i < length-1; i++ {
		str+=rIntChar(10,0)
	}
	return str
}
