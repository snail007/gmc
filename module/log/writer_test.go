// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"fmt"
	"github.com/fatih/color"
	gcore "github.com/snail007/gmc/core"
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
	w := NewFileWriter(&FileWriterOption{
		Filename: "logs-%h%i%s.log",
		LogsDir:  dir,
	}).ToWriter()
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
	w := NewFileWriter(&FileWriterOption{
		Filename: "foo.log",
		LogsDir:  dir,
	})
	assert.Nil(w)
}

func TestNewFileWriter2(t *testing.T) {
	dir := "fw1logs"
	defer os.RemoveAll(dir)
	assert := assert2.New(t)
	w := NewFileWriter(&FileWriterOption{
		Filename:      "logs-%h%i%s.log",
		LogsDir:       dir,
		AliasFilename: "app.log",
	}).ToWriter()
	assert.FileExists(filepath.Join(dir, "app.log"))
	time.Sleep(time.Second)
	w.Write([]byte("hello"))
	time.Sleep(time.Second)
	w.Write([]byte("hello"))
	fs, err := filepath.Glob(dir + "/*.log")
	assert.Nil(err)
	t.Log(fs)
	assert.Len(fs, 3)
}
func TestNewFileWriter_Gzip(t *testing.T) {
	dir := "fwgzlogs"
	os.RemoveAll(dir)
	defer func() {
		os.RemoveAll(dir)
	}()
	assert := assert2.New(t)
	w := NewFileWriter(&FileWriterOption{
		Filename: "logs-%h%i%s.log",
		LogsDir:  dir,
		IsGzip:   true,
	}).ToWriter()
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
	archiveDir := "%Y%m%d"
	fileName := "a_%h%i%s.log"
	os.RemoveAll(dir)
	defer func() {
		os.RemoveAll(dir)
	}()
	assert := assert2.New(t)
	w := NewFileWriter(&FileWriterOption{
		Filename:   fileName,
		LogsDir:    dir,
		ArchiveDir: archiveDir,
		IsGzip:     true,
	}).ToWriter()
	assert.Implements((*io.Writer)(nil), w)
	_, err := w.Write([]byte("abc"))
	assert.Nil(err)
	time.Sleep(time.Second * 2)
	f := dir + "/" + timeFormatText(time.Now(), fileName)
	_, err = w.Write([]byte("abc"))
	assert.Nil(err)
	assert.DirExists(gfile.Abs(dir))
	filepath.Glob(dir + "/*")
	d := filepath.Join(gfile.Abs(dir), timeFormatText(time.Now(), archiveDir))
	time.Sleep(time.Second * 2)
	assert.DirExists(d)
	t.Log(f)
	assert.Contains(gfile.ReadAll(f), "abc")
}

func TestWrite1(t *testing.T) {
	dir := "fwwlogs"
	archiveDir := "%Y%m%d"
	fileName := "a_%Y%m%d-%h%i%s.log"
	os.RemoveAll(dir)
	defer func() {
		os.RemoveAll(dir)
	}()
	assert := assert2.New(t)
	timeNowFunc = func() time.Time {
		return time.Now().Add(-time.Hour * 24)
	}
	w := NewFileWriter(&FileWriterOption{
		Filename:   fileName,
		LogsDir:    dir,
		ArchiveDir: archiveDir,
		IsGzip:     true,
	}).ToWriter()
	w.Write([]byte("abc"))
	timeNowFunc = func() time.Time {
		return time.Now()
	}
	time.Sleep(time.Second)
	w.Write([]byte("abc"))
	time.Sleep(time.Second * 3)
	p := filepath.Join(dir, timeFormatText(time.Now().Add(-time.Hour*24), archiveDir))
	assert.DirExists(p)
}

func TestColor(t *testing.T) {
	l := New()
	l.SetLevel(gcore.LogLevelTrace)
	l.Trace("trace")
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")
	fmt.Println(color.BlueString("%s", "abc"))
	fmt.Println(color.CyanString("%s", "abc"))
	fmt.Println(color.MagentaString("%s", "abc"))
	fmt.Println(color.WhiteString("%s", "abc"))
	t.Log("done")
}
