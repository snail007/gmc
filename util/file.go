// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More infomation at https://github.com/snail007/gmc

package gutil

import (
	"io"
	"io/ioutil"
	"os"
)

func ExistsFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	if stat.IsDir() {
		return false
	}
	return true
}

func ExistsDir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	return true
}

func Exists(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	_, err = f.Stat()
	if err != nil {
		return false
	}
	return true
}

func IsEmptyDir(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}

func GetString(file string) (str string, err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	str = string(b)
	return
}

func MustGetString(file string) (str string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	str = string(b)
	return
}

func GetBytes(file string) (b []byte, err error) {
	b, err = ioutil.ReadFile(file)
	return
}

func MustGetBytes(file string) (b []byte) {
	b, _ = ioutil.ReadFile(file)
	return
}
