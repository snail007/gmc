# GMC Middleware 模块

## 简介

GMC Middleware 模块提供常用的 HTTP 中间件，用于处理跨域、压缩、访问日志、错误恢复等功能。

## 功能特性

- **CORS**：跨域资源共享支持
- **Gzip**：HTTP 响应压缩
- **Access Log**：访问日志记录
- **Recovery**：Panic 恢复
- **Static**：静态文件服务
- **Auth**：认证中间件
- **RateLimit**：限流中间件

## 快速开始

### CORS 中间件

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    app := gmc.New.AppDefault()
    server := gmc.New.HTTPServer(app.Ctx())
    
    // 添加 CORS 中间件
    server.AddMiddleware(middleware.CORS(&middleware.CORSConfig{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           86400,
    }))
    
    server.Run()
}
```

### Gzip 压缩

```go
server.AddMiddleware(middleware.Gzip(&middleware.GzipConfig{
    Level: 5, // 压缩级别 1-9
}))
```

### 访问日志

```go
server.AddMiddleware(middleware.AccessLog(&middleware.AccessLogConfig{
    Logger: app.Logger("access"),
    Format: "[%s] %s %s %d %s",
}))
```

### Panic 恢复

```go
server.AddMiddleware(middleware.Recovery(&middleware.RecoveryConfig{
    Logger:     app.Logger("error"),
    StackTrace: true,
}))
```

## 使用场景

1. **API 服务**：CORS、压缩、日志、限流
2. **Web 应用**：静态文件、会话、认证
3. **微服务**：日志、监控、限流
4. **网关**：路由、认证、限流

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
