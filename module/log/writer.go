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
	"sync"
	"time"
)

type FileWriter struct {
	filename string
	dir      string
	filepath string
	file     *os.File
	isGzip   bool
	openLock sync.Mutex
	opened   *bool
}

func NewFileWriter(filename string, dir string, isGzip bool) (w *FileWriter) {
	if !existsDir(dir) {
		os.MkdirAll(dir, 0755)
	}
	logger := New()
	filename0 := filepath.Join(dir, timeFormatText(time.Now(), filename))
	w0, err := os.OpenFile(filename0, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	if err != nil {
		logger.Warnf("[FileWriter] WARN: new writer fail, error: %s", err)
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
		// close old log file
		f.file.Close()
		// open new logging file only once by locking.
		f.openLock.Lock()
		_, err = f.file.Write([]byte(""))
		// Double checking, avoid multiple open the new file.
		// If file is closed, err will be: file already closed, indicates file is old.
		if err != nil {
			f.file, err = os.OpenFile(filename0, os.O_CREATE|os.O_WRONLY, 0700)
		}
		f.openLock.Unlock()
		if err != nil {
			return
		}
		if f.isGzip {
			go func() {
				flog, e := os.OpenFile(oldFilepath, os.O_RDONLY, 0700)
				if e != nil {
					fmt.Printf("[FileWriter] WARN: compress log file fail, error :%v\n", e)
					return
				}
				defer flog.Close()

				fgz, e := os.OpenFile(oldFilepath+".gz", os.O_CREATE|os.O_WRONLY, 0700)
				if e != nil {
					fmt.Printf("[FileWriter] WARN: compress log file fail, error :%v\n", e)
					return
				}
				defer fgz.Close()

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
	n, err = f.file.Write(p)
	return
}
