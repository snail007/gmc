package log

import (
	"compress/gzip"
	gconfig "github.com/snail007/gmc/config"
	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/util"
	"io"
	"os"
	"path/filepath"
	"time"
)

func NewFromConfig(c *gconfig.Config) (l gcore.Logger) {
	l = NewGMCLog()
	cfg := c.Sub("log")
	l.SetLevel(gcore.LOG_LEVEL(cfg.GetInt("level")))
	if cfg.GetBool("async") {
		l.EnableAsync()
	}
	output := cfg.GetIntSlice("output")
	var writers []io.Writer
	for _, v := range output {
		switch v {
		case 0:
			writers = append(writers, os.Stdout)
		case 1:
			w0 := NewFileWriter(cfg.GetString("filename"),
				cfg.GetString("dir"), cfg.GetBool("gzip"))
			w0.SetLogger(l)
			writers = append(writers, w0)
		}
	}
	if len(writers) == 1 {
		l.SetOutput(writers[0])
	} else if len(writers) > 1 {
		l.SetOutput(io.MultiWriter(writers...))
	}
	return
}

type FileWriter struct {
	filename string
	dir      string
	filepath string
	file     *os.File
	isGzip   bool
	logger   gcore.Logger
}

func NewFileWriter(filename string, dir string, isGzip bool) (w *FileWriter) {
	if !gutil.ExistsDir(dir) {
		os.MkdirAll(dir, 0755)
	}
	logger := NewGMCLog()
	filename0 := filepath.Join(dir, gutil.TimeFormatText(time.Now(), filename))
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
		logger:   logger,
	}
	return
}

func (f *FileWriter) Filename() string {
	return f.filename
}

func (f *FileWriter) SetFilename(filename string) {
	f.filename = filename
}

func (f *FileWriter) Dir() string {
	return f.dir
}

func (f *FileWriter) SetDir(dir string) {
	f.dir = dir
}
func (f *FileWriter) SetLogger(logger gcore.Logger) {
	f.logger = logger
}
func (f *FileWriter) Write(p []byte) (n int, err error) {
	filename0 := filepath.Join(f.dir, gutil.TimeFormatText(time.Now(), f.filename))
	if filename0 != f.filepath {
		oldFilepath := f.filepath
		f.filepath = filename0
		f.file.Close()
		if f.isGzip {
			go func() {
				fgz, e := os.OpenFile(oldFilepath+".gz", os.O_CREATE|os.O_WRONLY, 0700)
				if e != nil {
					f.logger.Warnf("compress log file fail, error :%v", e)
					return
				}
				defer fgz.Close()
				flog, e := os.OpenFile(oldFilepath, os.O_RDONLY, 0700)
				if e != nil {
					f.logger.Warnf("compress log file fail, error :%v", e)
					return
				}
				defer flog.Close()
				zipWriter := gzip.NewWriter(fgz)
				defer zipWriter.Close()
				_, e = io.Copy(zipWriter, flog)
				if e != nil {
					f.logger.Warnf("compress log file fail, error :%v", e)
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
