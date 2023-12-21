// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package gfile

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Exists checks whether a file or directory exists.
// It returns false when the file or directory does not exists.
func Exists(path string) bool {
	_, e := os.Stat(path)
	return e == nil
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exists.
func IsFile(file string) bool {
	f, e := os.Stat(file)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsDir checks whether the path is a directory,
// it returns false when it's not a directory or does not exists.
func IsDir(file string) bool {
	f, e := os.Stat(file)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// IsLink returns true if path is a Symlink
func IsLink(path string) bool {
	s, err := os.Lstat(path)
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

// FileMTime get file modified time
func FileMTime(file string) (time.Time, error) {
	f, e := os.Stat(file)
	if e != nil {
		return time.Time{}, e
	}
	return f.ModTime(), nil
}

// HomeDir returns the home directory.
//
// Return "" if the home directory is empty.
func HomeDir() string {
	if v := os.Getenv("HOME"); v != "" { // For Unix/Linux
		return v
	} else if v := os.Getenv("HOMEPATH"); v != "" { // For Windows
		return v
	}
	return ""
}

var homeDir = HomeDir()

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

// FileName returns file name end of the path without extension.
func FileName(file string) string {
	basename := filepath.Base(file)
	if !strings.Contains(basename, ".") {
		return basename
	}
	return strings.TrimSuffix(basename, filepath.Ext(file))
}

// BaseName returns file name end of the path with extension.
func BaseName(file string) string {
	return filepath.Base(file)
}

// ReadAll returns string contents of file,
// if read fail, returns empty.
func ReadAll(file string) string {
	return string(Bytes(file))
}

// Bytes returns the []byte contents of file,
// if read fail, returns nil.
func Bytes(file string) (d []byte) {
	d, _ = ioutil.ReadFile(file)
	return
}

// Write writes []byte to file.
func Write(file string, data []byte, append bool) (err error) {
	mode := os.O_CREATE | os.O_WRONLY
	if append {
		mode |= os.O_APPEND
	}
	f, err := os.OpenFile(file, mode, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	return
}

// WriteString writes string to file.
func WriteString(file string, data string, append bool) (err error) {
	return Write(file, []byte(data), append)
}

// Copy source path file to destination file path, when mkdir is true and destination path not exists,
// the destination dir will be created
func Copy(srcFile, dstFile string, mkdir bool) error {
	sourceFile, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	dstPath := filepath.Dir(dstFile)
	if !IsDir(dstPath) {
		if mkdir {
			err = os.MkdirAll(dstPath, 0755)
		}
	}
	if err != nil {
		return err
	}
	destinationFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
