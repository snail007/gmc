// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileWriter struct {
	filename string
	dir      string
	filepath string
	file     *os.File
	isGzip   bool
}

func NewFileWriter(filename string, dir string, isGzip bool) (w *FileWriter) {
	if !existsDir(dir) {
		os.MkdirAll(dir, 0755)
	}
	logger := New()
	filename0 := filepath.Join(dir, timeFormatText(time.Now(), filename))
	w0, err := os.OpenFile(filename0, os.O_CREATE|os.O_APPEND, 0700)
	if err != nil {
		logger.Warnf("new writer fail, error: %s", err)
		return
	}
	w = &FileWriter{
		filename: filename,
		filepath: filename0,
		dir:      dir,
		file:     w0,
		isGzip:   isGzip,
	}
	return
}

func (f *FileWriter) Write(p []byte) (n int, err error) {
	filename0 := filepath.Join(f.dir, timeFormatText(time.Now(), f.filename))
	if filename0 != f.filepath {
		oldFilepath := f.filepath
		f.filepath = filename0
		f.file.Close()
		if f.isGzip {
			go func() {
				fgz, e := os.OpenFile(oldFilepath+".gz", os.O_CREATE|os.O_WRONLY, 0700)
				if e != nil {
					fmt.Printf("[FileWriter] WARN: compress log file fail, error :%v\n", e)
					return
				}
				defer fgz.Close()
				flog, e := os.OpenFile(oldFilepath, os.O_RDONLY, 0700)
				if e != nil {
					fmt.Printf("[FileWriter] WARN: compress log file fail, error :%v\n", e)
					return
				}
				defer flog.Close()
				zipWriter := gzip.NewWriter(fgz)
				defer zipWriter.Close()
				_, e = io.Copy(zipWriter, flog)
				if e != nil {
					fmt.Printf("[FileWriter] WARN: compress log file fail, error :%v\n", e)
					return
				}
				os.Remove(oldFilepath)
			}()
		}
	}
	f.file, err = os.OpenFile(filename0, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0700)
	if err != nil {
		return
	}
	n, err = f.file.Write(p)
	return
}
