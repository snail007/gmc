package gos

import (
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTempFile(t *testing.T) {
	assert2.Contains(t, TempFile("prefix", "suffix"), "prefix")
	assert2.Contains(t, TempFile("prefix", "suffix"), "suffix")
	assert2.Equal(t, strings.TrimSuffix(os.TempDir(), string(filepath.Separator)),
		strings.TrimSuffix(filepath.Dir(TempFile("prefix", "suffix")), string(filepath.Separator)))
	assert2.Equal(t, 44, len(filepath.Base(TempFile("prefix", "suffix"))))
}

func TestIsOSxxx1(t *testing.T) {
	assert2.True(t, IsWindows() || IsUnix() || IsMaxOS())
}
func TestIsOSxxx2(t *testing.T) {
	assert2.True(t, IsUnix() || IsMaxOS() || IsWindows())
}
func TestIsOSxxx3(t *testing.T) {
	assert2.True(t, IsMaxOS() || IsWindows() || IsUnix())
}

func TestCreateTempFile(t *testing.T) {
	f, err := CreateTempFile("", ".suffix")
	assert2.NoError(t, err)
	assert2.IsType(t, &os.File{}, f)
	assert2.Contains(t, f.Name(), ".suffix")
	if f != nil {
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()
	}
}
