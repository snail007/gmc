// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gfile

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// IsExists checks whether a file or directory exists.
// It returns false when the file or directory does not exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exists.
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsDir checks whether the path is a directory,
// it returns false when it's not a directory or does not exists.
func IsDir(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// IsLink returns true if path is a Symlink
func IsLink(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.Mode()&os.ModeSymlink == os.ModeSymlink
}

// IsEmptyDir returns true if the directory contains nothing, otherwise returns false
func IsEmptyDir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}

// ResolveLink resolves the link if it is a Link or return the original path
func ResolveLink(path string) (string, error) {
	if IsLink(path) {
		return os.Readlink(path)
	}
	return path, nil
}

// FileSize get file size as how many bytes
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// GetHomeDir returns the home directory.
//
// Return "" if the home direcotry is empty.
func GetHomeDir() string {
	if v := os.Getenv("HOME"); v != "" { // For Unix/Linux
		return v
	} else if v := os.Getenv("HOMEPATH"); v != "" { // For Windows
		return v
	}
	return ""
}

var homeDir = GetHomeDir()

// Abs is similar to Abs in the std library "path/filepath",
// but firstly convert "~"" and "$HOME" to the home directory.
//
// Return the origin path if there is an error.
func Abs(p string) string {
	p = strings.TrimSpace(p)
	if p != "" && homeDir != "" {
		_len := len(p)
		if p[0] == '~' {
			if _len == 1 || p[1] == '/' || p[1] == '\\' {
				p = filepath.Join(homeDir, p[1:])
			}
		} else if _len >= 5 && p[:5] == "$HOME" {
			if _len == 5 || (p[5] == '/' || p[5] == '\\') {
				p = filepath.Join(homeDir, p[5:])
			}
		}
	}
	if _p, err := filepath.Abs(p); err == nil {
		return _p
	}
	return p
}
