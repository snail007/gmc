// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog_test

import (
	glog "github.com/snail007/gmc/module/log"
	gfile "github.com/snail007/gmc/util/file"
	assert2 "github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewFileWriter(t *testing.T) {
	dir := "fwlogs"
	os.RemoveAll(dir)
	defer func() {
		os.RemoveAll(dir)
	}()
	assert := assert2.New(t)
	w := glog.NewFileWriter("logs-%h%i%s.log", dir, false)
	assert.Implements((*io.Writer)(nil), w)
	assert.DirExists(dir)
	fs, err := filepath.Glob(dir + "/*.log")
	assert.Nil(err)
	assert.Len(fs, 1)
	time.Sleep(time.Second)
	_, err = w.Write([]byte("\n"))
	assert.Nil(err)
	fs, err = filepath.Glob(dir + "/*.log")
	assert.Nil(err)
	assert.Len(fs, 2)
}

func TestNewFileWriter1(t *testing.T) {
	dir := "fw1logs"
	assert := assert2.New(t)
	f, err := os.Create(dir)
	assert.Nil(err)
	f.Close()
	defer os.Remove(f.Name())
	w := glog.NewFileWriter("foo.log", dir, false)
	assert.Nil(w)
}

func TestNewFileWriter_Gzip(t *testing.T) {
	dir := "fwgzlogs"
	os.RemoveAll(dir)
	defer func() {
		os.RemoveAll(dir)
	}()
	assert := assert2.New(t)
	w := glog.NewFileWriter("logs-%h%i%s.log", dir, true)
	assert.Implements((*io.Writer)(nil), w)
	assert.DirExists(dir)
	fs, err := filepath.Glob(dir + "/*.log")
	assert.Nil(err)
	assert.Len(fs, 1)
	time.Sleep(time.Second)
	_, err = w.Write([]byte("\n"))
	assert.Nil(err)
	time.Sleep(time.Second)
	fs, err = filepath.Glob(dir + "/*.gz")
	assert.Nil(err)
	assert.Len(fs, 1)
}

func TestWrite(t *testing.T) {
	dir := "fwwlogs"
	os.RemoveAll(dir)
	defer func() {
		os.RemoveAll(dir)
	}()
	assert := assert2.New(t)
	w := glog.NewFileWriter("a.log", dir, false)
	assert.Implements((*io.Writer)(nil), w)
	assert.DirExists(dir)
	_, err := w.Write([]byte("abc"))
	assert.Nil(err)
	assert.Contains(gfile.ReadAll(dir+"/a.log"), "abc")
}
