package gerror

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsePanic(t *testing.T) {
	text := `panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x109ea93]

goroutine 1 [running]:
main.(*foo).destruct(0xc208067e98)
	/path/to/main.go:22 +0x151
`

	err, parseErr := ParsePanic(text)
	if parseErr != nil {
		t.Errorf("Error parsing panic: %v", parseErr)
	}
	assert.NotNil(t, err)
}
