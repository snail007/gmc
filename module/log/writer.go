// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	"compress/gzip"
	"github.com/fatih/color"
	gcore "github.com/snail007/gmc/core"
	gbytes "github.com/snail007/gmc/util/bytes"
	gcast "github.com/snail007/gmc/util/cast"
	gfile "github.com/snail007/gmc/util/file"
	ghash "github.com/snail007/gmc/util/hash"
	grand "github.com/snail007/gmc/util/rand"
	gatomic "github.com/snail007/gmc/util/sync/atomic"
	gonce "github.com/snail007/gmc/util/sync/once"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	timeNowFunc                               = time.Now
	_                      gcore.LoggerWriter = &FileWriter{}
	listFiles              sync.Map
	listFilesOnce          sync.Once
	listFilesCLeanInterval = time.Hour
)

type listFile struct {
	logsDir    string
	maxBackups int
	lock       *sync.Mutex
}

func initListFilesDaemon() {
	isOnce := false
	listFilesOnce.Do(func() {
		isOnce = true
	})
	if !isOnce {
		return
	}
	go func() {
		for {
			time.Sleep(listFilesCLeanInterval)
			listFiles.Range(func(listFilePath_, v interface{}) bool {
				item := v.(*listFile)
				item.lock.Lock()
				defer item.lock.Unlock()
				listFilePath := listFilePath_.(string)
				var files []string
				for _, f := range strings.Split(string(gfile.Bytes(listFilePath)), "\n") {
					if strings.TrimSpace(f) == "" {
						continue
					}
					f = filepath.Join(item.logsDir, filepath.Clean(f))
					if gfile.Exists(f) {
						files = append(files, f)
					}
				}
				if len(files) <= item.maxBackups {
					return true
				}
				filesLen := len(files)
				idx := filesLen - item.maxBackups
				deleteList := files[:idx]
				keepList := files[idx:]
				for _, file := range deleteList {
					os.Remove(file)
				}
				for i, file := range keepList {
					keepList[i] = strings.TrimPrefix(file, item.logsDir+string(filepath.Separator))
				}
				os.Remove(listFilePath)
				gfile.WriteString(listFilePath, strings.Join(keepList, "\n"), false)
				return true
			})
		}
	}()
}

// FileWriterOption
// 1. the current log file final path is: LogsDir/Filename
// 2. the current backup log file final path is: LogsDir/ArchiveDir/Filename
// ArchiveDir and  Filename can contain:
// %Y: represents the year, will be replaced such as: 2012
// %m: represents the month, will be replaced such as: 12
// %d: represents the day of month, will be replaced such as: 31
// %h: represents the hour, will be replaced such as: 23
// %i: represents the minute, will be replaced such as: 36
// %s: represents the second, will be replaced such as: 00
// Example of ArchiveDir: "%Y%m%d", ArchiveDir will be split by day. ArchiveDir can't contain: /
// Example of Filename: "access_log_%Y%m%d%h.log", logfile will be split by hour.
type FileWriterOption struct {
	// Filename log file name, if AliasFilename is empty,
	// this will also be used as the current log file name.
	Filename string
	// LogsDir sets the log files store in
	LogsDir string
	// ArchiveDir sets the backups directory in the LogsDir
	ArchiveDir string
	// IsGzip sets if compress the backups
	IsGzip bool
	// AliasFilename current log file name
	AliasFilename string
	// MaxSize  max size of log file, example: 1M, 2G, empty no limit.
	MaxSize string
	// MaxBackups max backups of log files, 0 no limit
	MaxBackups int
	maxSize    int64
}

type FileWriter struct {
	filepath        string
	aliasFilepath   string
	file            *os.File
	archiveDir      string
	initLock        *sync.Mutex
	opt             *FileWriterOption
	writeBytesCount *gatomic.Int64
	fileListPath    string
	fileListLock    *sync.Mutex
}

func NewFileWriter(opt *FileWriterOption) (w *FileWriter) {
	opt.LogsDir = gfile.Abs(opt.LogsDir)
	w, err := NewFileWriterE(opt)
	if err != nil {
		Warnf("[FileWriter]: new writer fail, error: %s", err)
		return nil
	}
	return
}

func NewFileWriterE(opt *FileWriterOption) (w *FileWriter, err error) {
	if opt.MaxSize != "" {
		size, e := gbytes.ParseSize(opt.MaxSize)
		if e != nil {
			return nil, e
		}
		opt.maxSize = int64(size)
	}
	w = &FileWriter{
		opt:             opt,
		writeBytesCount: gatomic.NewInt64(0),
		fileListLock:    new(sync.Mutex),
		initLock:        new(sync.Mutex),
	}
	err = w.init()
	return
}

func (s *FileWriter) init() (err error) {
	if !existsDir(s.opt.LogsDir) {
		err = os.MkdirAll(s.opt.LogsDir, 0755)
		if err != nil {
			return
		}
	}
	if s.opt.MaxBackups > 0 {
		s.fileListPath = filepath.Join(s.opt.LogsDir, "."+ghash.MD5(s.opt.LogsDir+s.opt.Filename)+".file_list")
		if !gfile.Exists(s.fileListPath) {
			e := gfile.WriteString(s.fileListPath, "", false)
			if e != nil {
				return e
			}
			os.Chmod(s.fileListPath, 0600)
		}
		listFiles.Store(s.fileListPath, &listFile{
			maxBackups: s.opt.MaxBackups,
			logsDir:    s.opt.LogsDir,
			lock:       s.fileListLock,
		})
		initListFilesDaemon()
	}
	_, _, _, err = s.initLogfile()
	return
}

func (s *FileWriter) addFileToList(filePath string) {
	s.fileListLock.Lock()
	defer s.fileListLock.Unlock()
	filePath = strings.TrimPrefix(filePath, gfile.Abs(s.opt.LogsDir)+string(os.PathSeparator))
	gfile.WriteString(s.fileListPath, filePath+"\n", true)
}

func (s *FileWriter) getAliasFilepath() string {
	if s.opt.AliasFilename != "" {
		return filepath.Join(s.opt.LogsDir, s.opt.AliasFilename)
	}
	return s.getFilepath(true)
}

func (s *FileWriter) getFilepath(joinLogsPath bool) (filePath string) {
	defer func() {
		if s.opt.maxSize > 0 && s.writeBytesCount.Val() > s.opt.maxSize {
			filePath += "." + grand.String(8)
			return
		}
	}()
	if joinLogsPath {
		return filepath.Join(s.opt.LogsDir, timeFormatText(timeNowFunc(), s.opt.Filename))
	}
	return timeFormatText(timeNowFunc(), s.opt.Filename)
}

func (s *FileWriter) initLogfile() (changed bool, oldFilepath, oldArchiveDir string, err error) {
	s.initLock.Lock()
	defer s.initLock.Unlock()
	newFilePath := s.getFilepath(true)
	oldFilepath = s.filepath
	oldArchiveDir = s.archiveDir
	oldAliasFilepath := s.aliasFilepath
	changed = false
	if newFilePath == oldFilepath {
		return
	}
	changed = true
	isInit := false
	gonce.OnceDo(newFilePath, func() {
		isInit = true
	})
	if !isInit {
		return
	}
	s.filepath = newFilePath
	s.archiveDir = s.getArchivePath()
	s.aliasFilepath = s.getAliasFilepath()
	s.writeBytesCount.SetVal(0)
	if s.file != nil {
		s.file.Close()
	}
	if oldFilepath != "" && s.opt.AliasFilename != "" {
		os.Rename(oldAliasFilepath, oldFilepath)
	}
	s.file, err = os.OpenFile(s.aliasFilepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	return
}

func (s *FileWriter) compress(src string) (gzFilePath string, err error) {
	flog, err := os.OpenFile(src, os.O_RDONLY, 0700)
	if err != nil {
		return
	}
	defer flog.Close()
	gzFilePath = src + ".gz"
	fgz, err := os.OpenFile(gzFilePath, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return
	}
	defer fgz.Close()
	zipWriter := gzip.NewWriter(fgz)
	defer zipWriter.Close()
	_, err = io.Copy(zipWriter, flog)
	if err != nil {
		return
	}
	os.Remove(src)
	return
}

func (s *FileWriter) Write(p []byte, _ gcore.LogLevel) (n int, err error) {
	changed, oldFilePath, oldFileArchivePath, err := s.initLogfile()
	if err != nil {
		return
	}
	if changed {
		go func() {
			src := oldFilePath
			if s.opt.IsGzip {
				gzFilePath, e := s.compress(oldFilePath)
				if e != nil {
					Errorf("[FileWriter] compress fail, file: %s error :%v", oldFilePath, e)
					return
				}
				src = gzFilePath
			}
			e := s.move(src, oldFileArchivePath)
			if e != nil {
				Warnf("[FileWriter] move fail, file: %s error :%v", src, e)
			}
		}()
	}
	n, err = s.file.Write(p)
	s.writeBytesCount.Increase(int64(n))
	return
}
func (s *FileWriter) getArchivePath() string {
	archiveDir := timeFormatText(timeNowFunc(), s.opt.ArchiveDir)
	if archiveDir == "" {
		return ""
	}
	return filepath.Join(s.opt.LogsDir, archiveDir)
}
func (s *FileWriter) move(oldPath, oldArchiveDir string) (err error) {
	f := oldPath
	if oldArchiveDir != "" {
		if !gfile.Exists(oldArchiveDir) {
			err = os.MkdirAll(oldArchiveDir, 0755)
			if err != nil {
				return
			}
		}
		newFile := filepath.Join(oldArchiveDir, filepath.Base(oldPath))
		if e := os.Rename(oldPath, newFile); e != nil {
			return e
		}
		f = newFile
	}
	s.addFileToList(f)
	return
}

func (s *FileWriter) ToWriter() io.Writer {
	return NewIOWriter(s)
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
	defer func() {
		if err == nil {
			n = len(p)
		}
	}()
	switch level {
	case gcore.LogLevelTrace, gcore.LogLeveDebug:
		return s.w.Write([]byte(color.WhiteString("%s", string(p))))
	case gcore.LogLeveInfo:
		return s.w.Write([]byte(color.GreenString("%s", string(p))))
	case gcore.LogLeveWarn:
		return s.w.Write([]byte(color.YellowString("%s", string(p))))
	case gcore.LogLeveError, gcore.LogLevePanic, gcore.LogLeveFatal:
		return s.w.Write([]byte(color.RedString("%s", string(p))))
	default:
		return s.w.Write(p)
	}
}

func NewConsoleWriter() *ConsoleWriter {
	w := &ConsoleWriter{
		w:            os.Stdout,
		disableColor: runtime.GOOS == "windows" || gcast.ToBool(os.Getenv("DISABLE_CONSOLE_COLOR")),
	}
	w.redirected = w.isRedirected()
	return w
}
