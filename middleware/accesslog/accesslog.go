package accesslog

import (
	"fmt"
	"github.com/snail007/gmc"
	gconfig "github.com/snail007/gmc/config"
	log2 "github.com/snail007/gmc/gmc/log"
	grouter "github.com/snail007/gmc/http/router"
	"github.com/snail007/gmc/util/cast"
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

func NewWebFromConfig(c *gconfig.Config) gmc.MiddlewareWeb {
	a := newFromConfig(c)
	return func(ctx gcore.Ctx, server *gmc.HTTPServer) (isStop bool) {
		go log(ctx, a)
		return false
	}
}

func NewAPIFromConfig(c *gconfig.Config) gmc.MiddlewareAPI {
	a := newFromConfig(c)
	return func(ctx gcore.Ctx, server *gmc.APIServer) (isStop bool) {
		go log(ctx, a)
		return false
	}
}

func log(ctx gcore.Ctx, logger *accesslog) {
	rule := [][]string{
		[]string{"$host", ctx.Request.Host},
		[]string{"$uri", ctx.Request.URL.RequestURI()},
		[]string{"$time_used", cast.ToString(int(ctx.TimeUsed() / time.Millisecond))},
		[]string{"$status_code", cast.ToString(ctx.StatusCode())},
		[]string{"$query", ctx.Request.URL.Query().Encode()},
		[]string{"$req_time", time.Now().Format("2006-01-02 15:04:05")},
		[]string{"$client_ip", ctx.ClientIP()},
		[]string{"$remote_addr", ctx.Request.RemoteAddr},
		[]string{"$local_addr", ctx.LocalAddr},
	}
	str := logger.format
	for _, v := range rule {
		key := fmt.Sprintf("${%s}", v[0][1:])
		str = strings.Replace(str, key, v[1], 1)
		str = strings.Replace(str, v[0], v[1], 1)
	}
	logger.logger.Write(str + "\n")
}
