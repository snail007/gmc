// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

package accesslog

import (
	"fmt"
	"github.com/snail007/gmc/util/gpool"
	"strings"
	"time"

	gcore "github.com/snail007/gmc/core"
	"github.com/snail007/gmc/module/log"
	"github.com/snail007/gmc/util/cast"
)

var (
	pool = gpool.New(1)
)

func init() {
	pool.SetMaxJobCount(100000)
}

type accesslog struct {
	logger *glog.Logger
	format string
}

func newFromConfig(c gcore.Config) *accesslog {
	cfg := c.Sub("accesslog")
	logger := glog.New().(*glog.Logger)
	logger.SetFlag(gcore.LogFlagShort)
	logger.SetOutput(glog.NewFileWriter(&glog.FileWriterOption{
		Filename:      cfg.GetString("filename"),
		LogsDir:       cfg.GetString("dir"),
		ArchiveDir:    cfg.GetString("archive_dir"),
		IsGzip:        cfg.GetBool("gzip"),
		AliasFilename: cfg.GetString("filename_alias"),
	}))
	logger.EnableAsync()
	return &accesslog{
		format: cfg.GetString("format"),
		logger: logger,
	}
}

func NewFromConfig(c gcore.Config) gcore.Middleware {
	a := newFromConfig(c)
	return func(ctx gcore.Ctx) (isStop bool) {
		pool.Submit(func() {
			log(ctx, a)
		})
		return false
	}
}

func log(ctx gcore.Ctx, logger *accesslog) {
	rule := [][]string{
		{"$host", ctx.Request().Host},
		{"$uri", ctx.Request().URL.RequestURI()},
		{"$time_used", gcast.ToString(int(ctx.TimeUsed() / time.Millisecond))},
		{"$status_code", gcast.ToString(ctx.StatusCode())},
		{"$query", ctx.Request().URL.Query().Encode()},
		{"$req_time", time.Now().Format("2006-01-02 15:04:05")},
		{"$client_ip", ctx.ClientIP()},
		{"$remote_addr", ctx.Request().RemoteAddr},
		{"$local_addr", ctx.LocalAddr()},
	}
	str := logger.format
	for _, v := range rule {
		key := fmt.Sprintf("${%s}", v[0][1:])
		str = strings.Replace(str, key, v[1], 1)
		str = strings.Replace(str, v[0], v[1], 1)
	}
	logger.logger.WriteRaw(str+"\n", gcore.LogLeveInfo)
}
