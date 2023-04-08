// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"compress/gzip"
	"fmt"
	gfile "github.com/snail007/gmc/util/file"
	gonce "github.com/snail007/gmc/util/sync/once"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileWriter struct {
	filename   string
	logsDir    string
	filepath   string
	file       *os.File
	isGzip     bool
	openLock   sync.Mutex
	opened     *bool
	archiveDir string
}

func NewFileWriter(filename string, logsDir string, archiveDir string, isGzip bool) (w *FileWriter) {
	logsDir, _ = filepath.Abs(logsDir)
	if !existsDir(logsDir) {
		os.MkdirAll(logsDir, 0755)
	}
	logger := New()
	filename0 := filepath.Join(logsDir, timeFormatText(time.Now(), filename))
	w0, err := os.OpenFile(filename0, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	if err != nil {
		logger.Warnf("[FileWriter] WARN: new writer fail, error: %s", err)
		return
	}
	w = &FileWriter{
		filename:   filename,
		filepath:   filename0,
		logsDir:    logsDir,
		file:       w0,
		isGzip:     isGzip,
		archiveDir: archiveDir,
	}
	return
}

func (f *FileWriter) Write(p []byte) (n int, err error) {
	filename0 := filepath.Join(f.logsDir, timeFormatText(time.Now(), f.filename))
	if filename0 != f.filepath {
		oldFilepath := f.filepath
		gonce.OnceDo(oldFilepath, func() {
			f.filepath = filename0
			f.file, err = os.OpenFile(filename0, os.O_CREATE|os.O_WRONLY, 0700)
		})
		if err != nil {
			return
		}
		go func() {
			toMoveFile := oldFilepath
			if f.isGzip {
				flog, e := os.OpenFile(oldFilepath, os.O_RDONLY, 0700)
				if e != nil {
					fmt.Printf("[FileWriter] WARN: compress log file fail, error :%v\n", e)
					return
				}
				defer flog.Close()
				gzFile := oldFilepath + ".gz"
				fgz, e := os.OpenFile(gzFile, os.O_CREATE|os.O_WRONLY, 0700)
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
				toMoveFile = gzFile
			}
			f.Move(toMoveFile)
		}()
	}
	n, err = f.file.Write(p)
	return
}
func (f *FileWriter) Move(oldPath string) {
	archiveDir := timeFormatText(time.Now(), f.archiveDir)
	if archiveDir != "" {
		archiveDir = filepath.Join(f.logsDir, archiveDir)
		if !gfile.Exists(archiveDir) {
			e := os.MkdirAll(archiveDir, 0755)
			if e != nil {
				fmt.Printf("[FileWriter] WARN: create archive dir fail, dir: %s, error :%v\n", archiveDir, e)
				return
			}
		}
		newFile := filepath.Join(archiveDir, filepath.Base(oldPath))
		if e := os.Rename(oldPath, newFile); e != nil {
			fmt.Printf("[FileWriter] WARN: move log file to archive dir fail, file: %s, dst: %s, error :%v\n", oldPath, newFile, e)
		}
	}
}
