package gos

import (
	grand "github.com/snail007/gmc/util/rand"
	"os"
	"path/filepath"
)

func TempFile(prefix, suffix string) string {
	return filepath.Join(os.TempDir(), prefix+grand.String(32)+suffix)
}
