# GMC HTTP Server

## 简介

GMC HTTP Server 是一个功能强大的 HTTP 服务器实现，提供了完整的 Web 服务和 API 服务功能。支持路由、模板、会话、中间件、静态文件服务、HTTPS、HTTP/2 等特性。

## 功能特性

- **双服务器模式**：支持 HTTPServer（完整 Web 功能）和 APIServer（轻量级 API）
- **强大的路由**：基于高性能的 httprouter，支持 RESTful 路由
- **模板引擎**：支持多种模板引擎和自定义函数
- **会话管理**：内置会话支持，多种存储后端
- **中间件系统**：4 级中间件，灵活的请求处理流程
- **静态文件服务**：支持静态文件、嵌入式文件、Gzip 压缩
- **HTTPS/HTTP2**：完整的 TLS 和 HTTP/2 支持
- **优雅关闭**：支持优雅关闭和重启
- **i18n 支持**：内置国际化支持
- **错误处理**：可自定义 404、500 错误处理
- **连接统计**：实时连接数统计

## 安装

```bash
go get github.com/snail007/gmc
```

## 快速开始

### 创建简单的 Web 服务器

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    // 创建应用
    app := gmc.New.App()
    
    // 加载配置
    cfg := app.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    // 创建 HTTP 服务器
    s := gmc.New.HTTPServer(app.Ctx())
    s.Init(cfg)
    
    // 配置路由
    s.Router().HandlerFunc("GET", "/", func(c gmc.C) {
        c.Write("Hello GMC!")
    })
    
    // 添加服务到应用
    app.AddService(gmc.ServiceItem{
        Service: s,
    })
    
    // 运行应用
    app.Run()
}
```

### 创建 API 服务器

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

func main() {
    app := gmc.New.App()
    cfg := app.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    // 创建 API 服务器（轻量级，没有模板等功能）
    api, err := gmc.New.APIServer(app.Ctx(), ":8080")
    if err != nil {
        panic(err)
    }
    
    // 配置路由
    r := api.Router()
    r.HandlerFunc("GET", "/api/users", func(c gmc.C) {
        c.JSON(gcore.M{
            "users": []string{"Alice", "Bob"},
        })
    })
    
    // 添加服务
    app.AddService(gmc.ServiceItem{
        Service: api,
    })
    
    app.Run()
}
```

### 使用控制器

```go
package main

import (
    "github.com/snail007/gmc"
)

type UserController struct {
    gmc.Controller
}

func (c *UserController) List() {
    c.Write("User list")
}

func (c *UserController) Detail() {
    id := c.Query("id")
    c.Write("User detail: " + id)
}

func main() {
    app := gmc.New.App()
    cfg := app.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    s := gmc.New.HTTPServer(app.Ctx())
    s.Init(cfg)
    
    // 绑定控制器
    // 访问: /user/list, /user/detail
    s.Router().Controller("/user", new(UserController))
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

## 配置文件

### app.toml 基本配置

```toml
[httpserver]
# 监听地址
listen = ":8080"

# 优雅关闭超时时间（秒）
graceshutdowntimeout = 15

# 静态文件目录
static = "./static"

# 模板目录
template = "./views"

# 会话配置
[session]
ttl = 3600
store = "memory"

# 日志配置
[log]
level = "info"
output = "stdout"
```

### HTTPS 配置

```toml
[httpserver]
listen = ":8443"

# TLS 配置
[httpserver.tls]
# 启用 TLS
enable = true
# 证书文件
cert = "./cert.pem"
# 密钥文件
key = "./key.pem"
# 客户端证书认证
[httpserver.tls.clientauth]
enable = false
ca = ""
```

### HTTP/2 配置

HTTP/2 在启用 TLS 后自动启用，无需额外配置。

## 中间件系统

GMC 提供 4 级中间件，按执行顺序：

### Middleware0 - 路由前中间件

在路由匹配之前执行，可以用于全局过滤：

```go
s.AddMiddleware0(func(c gmc.C) (isStop bool) {
    // 记录所有请求
    c.Logger().Info(c.Request().URL.Path)
    
    // 返回 true 停止后续处理
    if c.Request().URL.Path == "/blocked" {
        c.WriteHeader(403)
        c.Write("Forbidden")
        return true
    }
    
    return false
})
```

### Middleware1 - 路由后、控制器前

路由匹配后、控制器方法执行前：

```go
s.AddMiddleware1(func(c gmc.C) (isStop bool) {
    // 身份验证
    token := c.Header("Authorization")
    if token == "" {
        c.WriteHeader(401)
        c.JSON(gcore.M{"error": "Unauthorized"})
        return true
    }
    
    return false
})
```

### Middleware2 - 控制器后

控制器方法执行后：

```go
s.AddMiddleware2(func(c gmc.C) (isStop bool) {
    // 添加响应头
    c.SetHeader("X-Powered-By", "GMC")
    return false
})
```

### Middleware3 - 最后执行

所有处理完成后，适合日志记录：

```go
s.AddMiddleware3(func(c gmc.C) (isStop bool) {
    // 记录响应状态和耗时
    c.Logger().Infof("%s %d %dms", 
        c.Request().URL.Path, 
        c.StatusCode(), 
        c.TimeUsed()/1000000)
    return false
})
```

### 使用内置中间件

```go
import "github.com/snail007/gmc/module/middleware/accesslog"

func InitRouter(s *gmc.HTTPServer) {
    // 访问日志中间件
    s.AddMiddleware3(accesslog.NewFromConfig(s.Config()))
    
    // 其他路由配置...
}
```

## 静态文件服务

### 提供静态文件

```go
// 方法 1：配置文件
// app.toml
// [httpserver]
// static = "./static"
// staticurlpath = "/static"

// 方法 2：代码设置
s.SetStaticDir("./static")
s.SetStaticUrlPath("/static")

// 访问: http://localhost:8080/static/css/style.css
// 映射到: ./static/css/style.css
```

### 嵌入式静态文件

```go
package main

import (
    "embed"
    "github.com/snail007/gmc"
)

//go:embed static/*
var staticFS embed.FS

func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    // 从嵌入式文件系统提供静态文件
    s.Ext(".html")
    s.Router().HandlerAny("/static/*filepath", func(c gmc.C) {
        filepath := c.Param("filepath")
        data, _ := staticFS.ReadFile("static" + filepath)
        c.Write(string(data))
    })
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

## 模板渲染

### 使用模板

```go
type HomeController struct {
    gmc.Controller
}

func (c *HomeController) Index() {
    data := gcore.M{
        "title": "Welcome",
        "name":  "GMC",
    }
    c.View("index", data)
}
```

### 添加模板函数

```go
s.AddFuncMap(map[string]interface{}{
    "add": func(a, b int) int {
        return a + b
    },
    "upper": func(s string) string {
        return strings.ToUpper(s)
    },
})
```

模板中使用：
```html
<h1>{{.title}}</h1>
<p>{{upper .name}}</p>
<p>Sum: {{add 1 2}}</p>
```

## 会话管理

```go
type LoginController struct {
    gmc.Controller
}

func (c *LoginController) Login() {
    username := c.PostForm("username")
    password := c.PostForm("password")
    
    // 验证...
    
    // 启动会话
    c.SessionStart()
    defer c.SessionDestroy()
    
    // 设置会话数据
    c.Session().Set("user_id", "123")
    c.Session().Set("username", username)
    
    c.JSON(gcore.M{"status": "ok"})
}

func (c *LoginController) Profile() {
    c.SessionStart()
    defer c.SessionDestroy()
    
    // 读取会话数据
    userID, ok := c.Session().Get("user_id")
    if !ok {
        c.WriteHeader(401)
        c.JSON(gcore.M{"error": "Not logged in"})
        return
    }
    
    c.JSON(gcore.M{
        "user_id": userID,
    })
}
```

## 错误处理

### 自定义 404 处理

```go
s.SetNotFoundHandler(func(c gmc.C, tpl gcore.Template) {
    c.WriteHeader(404)
    c.View("404", gcore.M{
        "path": c.Request().URL.Path,
    })
})
```

### 自定义 500 处理

```go
s.SetErrorHandler(func(c gmc.C, tpl gcore.Template, err interface{}) {
    c.WriteHeader(500)
    c.Logger().Errorf("Internal error: %v", err)
    c.View("500", gcore.M{
        "error": err,
    })
})
```

## HTTPS 和 HTTP/2

### 启用 HTTPS

```toml
[httpserver]
listen = ":8443"

[httpserver.tls]
enable = true
cert = "./cert.pem"
key = "./key.pem"
```

### 生成自签名证书（测试用）

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

### 客户端证书认证

```toml
[httpserver.tls]
enable = true
cert = "./server.pem"
key = "./server-key.pem"

[httpserver.tls.clientauth]
enable = true
ca = "./ca.pem"
```

## 高级功能

### 自定义监听器

```go
s.SetListenerFactory(func(addr string) (net.Listener, error) {
    // 自定义 Listener 配置
    lc := net.ListenConfig{
        KeepAlive: 30 * time.Second,
    }
    return lc.Listen(context.Background(), "tcp", addr)
})
```

### 连接计数

```go
count := s.ActiveConnCount()
fmt.Printf("Active connections: %d\n", count)
```

### 优雅关闭

```go
// 配置文件设置超时
[httpserver]
graceshutdowntimeout = 15

// 或代码设置
s.SetConfig(cfg)

// 服务器会在接收到 SIGTERM/SIGINT 信号时优雅关闭
```

## API 参考

### HTTPServer 主要方法

```go
// 初始化
func (s *HTTPServer) Init(cfg gcore.Config) error

// 路由
func (s *HTTPServer) Router() gcore.HTTPRouter

// 中间件
func (s *HTTPServer) AddMiddleware0(m gcore.Middleware)
func (s *HTTPServer) AddMiddleware1(m gcore.Middleware)
func (s *HTTPServer) AddMiddleware2(m gcore.Middleware)
func (s *HTTPServer) AddMiddleware3(m gcore.Middleware)

// 静态文件
func (s *HTTPServer) SetStaticDir(dir string)
func (s *HTTPServer) SetStaticUrlPath(path string)

// 错误处理
func (s *HTTPServer) SetNotFoundHandler(fn func(ctx gcore.Ctx, tpl gcore.Template))
func (s *HTTPServer) SetErrorHandler(fn func(ctx gcore.Ctx, tpl gcore.Template, err interface{}))

// 模板
func (s *HTTPServer) AddFuncMap(f map[string]interface{})

// 配置
func (s *HTTPServer) Config() gcore.Config
func (s *HTTPServer) SetConfig(c gcore.Config)

// 其他
func (s *HTTPServer) Logger() gcore.Logger
func (s *HTTPServer) ActiveConnCount() int64
func (s *HTTPServer) Close()
```

### APIServer

APIServer 是 HTTPServer 的轻量级版本，移除了模板、会话等功能，适合纯 API 服务：

```go
// 创建 API 服务器
func NewAPIServer(ctx gcore.Ctx, address string) (gcore.APIServer, error)

// API 服务器有相同的路由和中间件功能，但没有：
// - 模板引擎
// - 会话管理
// - 静态文件服务
```

## 最佳实践

### 1. 项目结构

```
myapp/
├── main.go              # 入口
├── app.toml            # 配置
├── controller/         # 控制器
│   ├── user.go
│   └── product.go
├── router/             # 路由配置
│   └── router.go
├── initialize/         # 初始化
│   └── init.go
├── model/              # 数据模型
├── views/              # 模板
└── static/             # 静态文件
```

### 2. 路由组织

```go
// router/router.go
package router

import "github.com/snail007/gmc"

func InitRouter(s *gmc.HTTPServer) {
    r := s.Router()
    
    // 静态页面
    r.Controller("/", new(controller.Home))
    
    // API 路由
    api := r.Group("/api")
    {
        api.Controller("/users", new(controller.User))
        api.Controller("/products", new(controller.Product))
    }
    
    // 管理后台
    admin := r.Group("/admin")
    {
        admin.Controller("/dashboard", new(controller.Dashboard))
    }
}
```

### 3. 中间件使用

```go
// 按需使用不同级别的中间件
s.AddMiddleware0(corsMiddleware)      // 全局 CORS
s.AddMiddleware1(authMiddleware)      // 需要路由信息的认证
s.AddMiddleware2(responseMiddleware)  // 修改响应
s.AddMiddleware3(accesslog.NewFromConfig(s.Config())) // 访问日志
```

### 4. 错误恢复

控制器中的 panic 会被自动捕获，触发 500 错误处理器。

## 性能优化

1. **启用 Gzip**：自动压缩响应（配置 `gzip = true`）
2. **静态文件缓存**：设置适当的 Cache-Control 头
3. **连接池**：数据库连接池配置
4. **异步日志**：使用异步日志减少 I/O 阻塞
5. **HTTP/2**：启用 HTTP/2 提升性能

## 注意事项

1. **端口权限**：监听 1024 以下端口需要 root 权限
2. **优雅关闭**：设置合理的 graceshutdowntimeout
3. **日志记录**：生产环境使用异步日志
4. **HTTPS**：生产环境务必使用 HTTPS
5. **会话安全**：使用 HttpOnly 和 Secure Cookie

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [GMC Router](../router/README.md)
- [GMC Controller](../controller/README.md)
- [GMC Session](../session/README.md)
- [GMC Template](../template/README.md)
