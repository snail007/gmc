package gmclog

import (
	"compress/gzip"
	gmcconfig "github.com/snail007/gmc/config"
	gmccore "github.com/snail007/gmc/core"
	fileutil "github.com/snail007/gmc/util/file"
	timeutil "github.com/snail007/gmc/util/time"
	"io"
	"os"
	"path/filepath"
	"time"
)

func NewFromConfig(c *gmcconfig.Config) (l gmccore.Logger, err error) {
	l = NewGMCLog()
	cfg := c.Sub("log")
	l.SetLevel(gmccore.LOG_LEVEL(cfg.GetInt("level")))
	output := cfg.GetIntSlice("output")
	var writers []io.Writer
	for _, v := range output {
		switch v {
		case 0:
			writers=append(writers,os.Stdout)
		case 1:
			w0, err := NewFileWriter(cfg.GetString("filename"),
				cfg.GetString("dir"), cfg.GetBool("gzip"))
			if err != nil {
				return nil, err
			}
			w0.SetLogger(l)
			writers=append(writers,w0)
		}
	}
	if len(writers)==1{
		l.SetOutput(writers[0])
	}else if len(writers)>1{
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
	logger   gmccore.Logger
}

func NewFileWriter(filename string, dir string, isGzip bool) (w *FileWriter, err error) {
	if !fileutil.ExistsDir(dir) {
		os.MkdirAll(dir, 0755)
	}
	filename0 := filepath.Join(dir, timeutil.TimeFormatText(time.Now(), filename))
	w0, err := os.OpenFile(filename0, os.O_CREATE|os.O_APPEND, 0700)
	if err != nil {
		return
	}
	w = &FileWriter{
		filename: filename,
		filepath: filename0,
		dir:      dir,
		file:     w0,
		isGzip:   isGzip,
		logger:   NewGMCLog(),
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
func (f *FileWriter) SetLogger(logger gmccore.Logger) {
	f.logger = logger
}
func (f *FileWriter) Write(p []byte) (n int, err error) {
	filename0 := filepath.Join(f.dir, timeutil.TimeFormatText(time.Now(), f.filename))
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