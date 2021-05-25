package gos

import (
	grand "github.com/snail007/gmc/util/rand"
	"os"
	"path/filepath"
	"runtime"
)

// TempFile returns a path string, the path is join temp path and a random string.
// such as is: /tmp/prefix.xxx.suffix. xxx is a length 32 random string, prefix and suffix are arguments.
func TempFile(prefix, suffix string) string {
	return filepath.Join(os.TempDir(), prefix+grand.String(32)+suffix)
}

// CreateTempFile creates a file in temp path, file name is randomly.
func CreateTempFile(prefix, suffix string) (file *os.File, err error) {
	f := TempFile(prefix, suffix)
	file, err = os.OpenFile(f, os.O_CREATE|os.O_RDWR, 0755)
	return
}

// IsWindows returns true indicates the os platform is Microsoft Windows.
// Others returns false.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsUnix returns true indicates the os platform is Unix or linux based System.
// Others returns false.
func IsUnix() bool {
	return filepath.Separator == '/'
}

// IsMaxOS returns true indicates the os platform is Apple Mac OS.
// Others returns false.
func IsMaxOS() bool {
	return runtime.GOOS == "darwin"
}
