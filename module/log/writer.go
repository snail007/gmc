// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"compress/gzip"
	gcore "github.com/snail007/gmc/core"
	gcast "github.com/snail007/gmc/util/cast"
	gfile "github.com/snail007/gmc/util/file"
	gonce "github.com/snail007/gmc/util/sync/once"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var timeNowFunc = time.Now

type FileWriter struct {
	filepath   string
	file       *os.File
	archiveDir string
	openLock   sync.Mutex
	opt        *FileWriterOption
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
	s.archiveDir = s.getArchiveDir()
	s.file, err = os.OpenFile(s.getAltFilepath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	return
}

func (s *FileWriter) getAltFilepath() string {
	filename := timeFormatText(timeNowFunc(), s.opt.Filename)
	if s.opt.AliasFilename != "" {
		filename = s.opt.AliasFilename
	}
	return filepath.Join(s.opt.LogsDir, filename)
}

func (s *FileWriter) getRawFilepath() string {
	return filepath.Join(s.opt.LogsDir, timeFormatText(timeNowFunc(), s.opt.Filename))
}

func (s *FileWriter) Write(p []byte) (n int, err error) {
	filename0 := s.getRawFilepath()
	if filename0 != s.filepath {
		oldFilepath := s.filepath
		oldArchiveDir := s.archiveDir
		gonce.OnceDo(oldFilepath, func() {
			s.filepath = s.getRawFilepath()
			s.archiveDir = s.getArchiveDir()
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
					Errorf("[FileWriter] open log file fail, file: %s error :%v", oldFilepath, e)
					return
				}
				defer flog.Close()
				gzFile := oldFilepath + ".gz"
				fgz, e := os.OpenFile(gzFile, os.O_CREATE|os.O_WRONLY, 0700)
				if e != nil {
					Errorf("[FileWriter] create gzip log file fail, file: %s, error :%v", gzFile, e)
					return
				}
				defer fgz.Close()
				zipWriter := gzip.NewWriter(fgz)
				defer zipWriter.Close()
				_, e = io.Copy(zipWriter, flog)
				if e != nil {
					Errorf("[FileWriter] WARN: write gzip log file fail, error :%v\n", e)
					return
				}
				os.Remove(oldFilepath)
				toMoveFile = gzFile
			}
			s.Move(toMoveFile, oldArchiveDir)
		}()
	}
	n, err = s.file.Write(p)
	return
}
func (s *FileWriter) getArchiveDir() string {
	archiveDir := timeFormatText(timeNowFunc(), s.opt.ArchiveDir)
	if archiveDir == "" {
		return ""
	}
	return filepath.Join(s.opt.LogsDir, archiveDir)
}
func (s *FileWriter) Move(oldPath, oldArchiveDir string) {
	if oldArchiveDir != "" {
		if !gfile.Exists(oldArchiveDir) {
			e := os.MkdirAll(oldArchiveDir, 0755)
			if e != nil {
				Warnf("[FileWriter] create archive dir fail, dir: %s, error :%v\n", oldArchiveDir, e)
				return
			}
		}
		newFile := filepath.Join(oldArchiveDir, filepath.Base(oldPath))
		if e := os.Rename(oldPath, newFile); e != nil {
			Warnf("[FileWriter] move log file to archive dir fail, file: %s, dst: %s, error :%v\n", oldPath, newFile, e)
		}
	}
}

type ConsoleWriter struct {
	w            io.Writer
	redirected   bool
	disableColor bool
}

func (s *ConsoleWriter) isRedirected() bool {
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		mode := fileInfo.Mode()
		if mode&os.ModeCharDevice == 0 {
			return true
		}
	}
	if fileInfo, err := os.Stderr.Stat(); err == nil {
		mode := fileInfo.Mode()
		if mode&os.ModeCharDevice == 0 {
			return true
		}
	}
	return false
}

func (s *ConsoleWriter) Write(p []byte, level gcore.LogLevel) (n int, err error) {
	if s.isRedirected() || s.disableColor {
		return s.w.Write(p)
	}
	switch level {
	case gcore.LogLeveInfo:
		return s.w.Write([]byte(Green.Color(string(p))))
	case gcore.LogLeveWarn:
		return s.w.Write([]byte(Yellow.Color(string(p))))
	case gcore.LogLeveError, gcore.LogLevePanic, gcore.LogLeveFatal:
		return s.w.Write([]byte(Red.Color(string(p))))
	default:
		return s.w.Write(p)
	}
}

func (s *ConsoleWriter) Writer() io.Writer {
	return s.w
}

func NewConsoleWriter() *ConsoleWriter {
	w := &ConsoleWriter{
		w:            os.Stdout,
		disableColor: gcast.ToBool(os.Getenv("DISABLE_CONSOLE_COLOR")),
	}
	w.redirected = w.isRedirected()
	return w
}
