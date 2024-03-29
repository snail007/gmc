// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package glog

import (
	gcore "github.com/snail007/gmc/core"
	"io"
	"os"
	"strings"
	"time"
)

func NewFromConfig(c gcore.Config, prefix ...string) (l gcore.Logger) {
	prefix0 := ""
	if len(prefix) == 1 {
		prefix0 = prefix[0]
	}
	l = New(prefix0)
	cfg := c.Sub("log")
	if cfg == nil {
		return
	}
	l.SetLevel(gcore.LogLevel(cfg.GetInt("level")))
	if cfg.GetBool("async") {
		l.EnableAsync()
	}
	output := cfg.GetIntSlice("output")
	var writers []gcore.LoggerWriter
	for _, v := range output {
		switch v {
		case 0:
			writers = append(writers, NewConsoleWriter())
		case 1:
			w0 := NewFileWriter(&FileWriterOption{
				Filename:      cfg.GetString("filename"),
				LogsDir:       cfg.GetString("dir"),
				ArchiveDir:    cfg.GetString("archive_dir"),
				IsGzip:        cfg.GetBool("gzip"),
				AliasFilename: cfg.GetString("filename_alias"),
			})
			writers = append(writers, w0)
		}
	}
	if len(writers) == 1 {
		l.SetOutput(writers[0])
	} else if len(writers) > 1 {
		l.SetOutput(newMultiWriter(writers...))
	}
	return
}

func existsDir(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	stat, _ := f.Stat()
	return stat != nil && stat.IsDir()
}

func timeFormatText(t time.Time, text string) string {
	d := map[string]string{
		"Y": t.Format("2006"),
		"y": t.Format("06"),
		"m": t.Format("01"),
		"n": t.Format("1"),
		"d": t.Format("02"),
		"j": t.Format("2"),
		"H": t.Format("15"),
		"h": t.Format("03"),
		"g": t.Format("3"),
		"i": t.Format("04"),
		"s": t.Format("05"),
	}
	for k, v := range d {
		k = "%" + k
		text = strings.Replace(text, k, v, 1)
	}
	return text
}

type LoggerWriter struct {
	w io.Writer
}

func (s *LoggerWriter) Write(p []byte, _ gcore.LogLevel) (n int, err error) {
	return s.w.Write(p)
}

func NewLoggerWriter(w io.Writer) gcore.LoggerWriter {
	return &LoggerWriter{w: w}
}

type IOWriter struct {
	w gcore.LoggerWriter
}

func NewIOWriter(w gcore.LoggerWriter) io.Writer {
	return &IOWriter{w: w}
}
func (s *IOWriter) Write(p []byte) (n int, err error) {
	return s.w.Write(p, gcore.LogLeveInfo)
}
