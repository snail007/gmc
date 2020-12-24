package accesslog

import (
	"fmt"
	gcore "github.com/snail007/gmc/core"
	log2 "github.com/snail007/gmc/module/log"
	"github.com/snail007/gmc/util/cast"
	gconfig "github.com/snail007/gmc/util/config"
	"strings"
	"time"
)

type accesslog struct {
	logger *log2.GMCLog
	format string
}

func newFromConfig(c *gconfig.Config) *accesslog {
	cfg := c.Sub("accesslog")
	logger := log2.NewGMCLog().(*log2.GMCLog)
	logger.SetFlags(0)
	logger.SetOutput(log2.NewFileWriter(cfg.GetString("filename"),
		cfg.GetString("dir"), cfg.GetBool("gzip")))
	logger.EnableAsync()
	return &accesslog{
		format: cfg.GetString("format"),
		logger: logger,
	}
}

func NewWebFromConfig(c *gconfig.Config) gcore.MiddlewareWeb {
	a := newFromConfig(c)
	return func(ctx gcore.Ctx, server gcore.HTTPServer) (isStop bool) {
		go log(ctx, a)
		return false
	}
}

func NewAPIFromConfig(c *gconfig.Config) gcore.MiddlewareAPI {
	a := newFromConfig(c)
	return func(ctx gcore.Ctx, server gcore.APIServer) (isStop bool) {
		go log(ctx, a)
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
	logger.logger.Write(str + "\n")
}
