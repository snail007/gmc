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
	filepath string
	file     *os.File
	openLock sync.Mutex
	opt      *FileWriterOption
}
type FileWriterOption struct {
	Filename      string
	LogsDir       string
	ArchiveDir    string
	IsGzip        bool
	AliasFilename string
}

func NewFileWriter(opt *FileWriterOption) (w *FileWriter) {
	w, err := NewFileWriterE(opt)
	if err != nil {
		Warnf("[FileWriter]: new writer fail, error: %s", err)
		return nil
	}
	return
}

func NewFileWriterE(opt *FileWriterOption) (w *FileWriter, err error) {
	w = &FileWriter{opt: opt}
	err = w.init()
	return
}

func (s *FileWriter) init() (err error) {
	logsDir, _ := filepath.Abs(s.opt.LogsDir)
	if !existsDir(logsDir) {
		err = os.MkdirAll(logsDir, 0755)
		if err != nil {
			return
		}
	}
	s.filepath = s.getRawFilepath()
	s.file, err = os.OpenFile(s.getAltFilepath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	return
}

func (s *FileWriter) getAltFilepath() string {
	filename := timeFormatText(time.Now(), s.opt.Filename)
	if s.opt.AliasFilename != "" {
		filename = s.opt.AliasFilename
	}
	return filepath.Join(s.opt.LogsDir, filename)
}

func (s *FileWriter) getRawFilepath() string {
	return filepath.Join(s.opt.LogsDir, timeFormatText(time.Now(), s.opt.Filename))
}

func (s *FileWriter) Write(p []byte) (n int, err error) {
	filename0 := s.getRawFilepath()
	if filename0 != s.filepath {
		oldFilepath := s.filepath
		gonce.OnceDo(oldFilepath, func() {
			s.filepath = filename0
			if s.file != nil {
				s.file.Close()
			}
			if s.opt.AliasFilename != "" {
				os.Rename(s.getAltFilepath(), oldFilepath)
			}
			s.file, err = os.OpenFile(s.getAltFilepath(), os.O_CREATE|os.O_WRONLY, 0700)
		})
		if err != nil {
			return
		}
		go func() {
			toMoveFile := oldFilepath
			if s.opt.IsGzip {
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
			s.Move(toMoveFile)
		}()
	}
	n, err = s.file.Write(p)
	return
}
func (s *FileWriter) Move(oldPath string) {
	archiveDir := timeFormatText(time.Now(), s.opt.ArchiveDir)
	if archiveDir != "" {
		archiveDir = filepath.Join(s.opt.LogsDir, archiveDir)
		if !gfile.Exists(archiveDir) {
			e := os.MkdirAll(archiveDir, 0755)
			if e != nil {
				Warnf("[FileWriter] create archive dir fail, dir: %s, error :%v\n", archiveDir, e)
				return
			}
		}
		newFile := filepath.Join(archiveDir, filepath.Base(oldPath))
		if e := os.Rename(oldPath, newFile); e != nil {
			Warnf("[FileWriter] move log file to archive dir fail, file: %s, dst: %s, error :%v\n", oldPath, newFile, e)
		}
	}
}
