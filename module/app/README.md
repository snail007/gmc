# GMC App 模块

## 简介

GMC App 模块是 GMC 框架的核心应用管理模块，提供应用的生命周期管理、热重载、服务管理、配置加载等功能。

## 功能特性

- **应用生命周期管理**：统一管理应用的初始化、启动、停止流程
- **热重载（Hot Reload）**：支持零停机时间的应用热更新（仅 Linux）
- **服务管理**：管理多个服务的启动和停止
- **配置管理**：集成配置加载和管理
- **日志管理**：集成日志系统
- **数据库管理**：自动初始化数据库连接
- **缓存管理**：自动初始化缓存系统
- **国际化支持**：集成 i18n 功能
- **优雅关闭**：支持优雅停止服务

## 安装

```bash
go get github.com/snail007/gmc/module/app
```

## 快速开始

### 创建默认应用

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    // 创建默认应用，自动加载 app.toml 配置文件
    app := gmc.New.AppDefault()
    
    // 运行应用
    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### 创建自定义应用

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

func main() {
    // 创建新应用
    app := gmc.New.App()
    
    // 设置配置文件
    cfg := gmc.New.Config()
    cfg.SetConfigFile("./config/app.toml")
    err := cfg.ReadInConfig()
    if err != nil {
        panic(err)
    }
    app.SetConfig(cfg)
    
    // 运行应用
    err = app.Run()
    if err != nil {
        panic(err)
    }
}
```

### 添加自定义服务

```go
package main

import (
    "fmt"
    "net"
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

// 自定义服务需要实现 gcore.Service 接口
type MyService struct {
    cfg      gcore.Config
    logger   gcore.Logger
    listener net.Listener
}

func (s *MyService) Init(cfg gcore.Config) error {
    s.cfg = cfg
    fmt.Println("MyService initialized")
    return nil
}

func (s *MyService) Start() error {
    fmt.Println("MyService started")
    // 启动服务逻辑
    return nil
}

func (s *MyService) Stop() {
    fmt.Println("MyService stopped")
}

func (s *MyService) GracefulStop() {
    fmt.Println("MyService gracefully stopped")
    s.Stop()
}

func (s *MyService) SetLog(logger gcore.Logger) {
    s.logger = logger
}

func (s *MyService) InjectListeners(listeners []net.Listener) {
    if len(listeners) > 0 {
        s.listener = listeners[0]
    }
}

func (s *MyService) Listeners() []net.Listener {
    if s.listener != nil {
        return []net.Listener{s.listener}
    }
    return nil
}

func main() {
    app := gmc.New.AppDefault()
    
    // 添加自定义服务
    app.AddService(gcore.ServiceItem{
        Service: &MyService{},
        BeforeInit: func(s *gcore.ServiceItem) error {
            fmt.Println("Before init service")
            return nil
        },
        BeforeStart: func(s *gcore.ServiceItem) error {
            fmt.Println("Before start service")
            return nil
        },
    })
    
    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### 应用生命周期钩子

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

func main() {
    app := gmc.New.AppDefault()
    
    // 在应用运行前执行
    app.OnRun(func(cfg gcore.Config) error {
        fmt.Println("App is about to run")
        // 可以在这里进行初始化操作
        return nil
    })
    
    // 在应用关闭时执行
    app.OnShutdown(func() {
        fmt.Println("App is shutting down")
        // 可以在这里进行清理操作
    })
    
    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

## 热重载（Hot Reload）

热重载允许在不停止服务的情况下更新应用代码，仅支持 Linux 系统。

### 启用热重载

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    app := gmc.New.AppDefault()
    
    // 添加支持热重载的服务
    // 服务需要实现 InjectListeners 和 Listeners 方法
    
    err := app.Run()
    if err != nil {
        panic(err)
    }
}
```

### 触发热重载

在 Linux 系统上，发送 USR2 信号触发热重载：

```bash
# 发送热重载信号
pkill -USR2 yourappname

# 或使用 kill 命令
kill -USR2 <pid>
```

### 热重载流程

1. **接收信号**：应用接收到 USR2 信号
2. **保存状态**：调用所有服务的 `Listeners()` 获取监听器
3. **启动子进程**：启动新进程并传递监听器文件描述符
4. **初始化新进程**：
   - 调用 `InjectListeners()` 注入监听器
   - 调用 `Init()` 初始化服务
   - 调用 `Start()` 启动服务
5. **停止旧进程**：调用 `GracefulStop()` 优雅停止旧服务

### Service 接口

```go
type Service interface {
    // 初始化服务
    Init(cfg gcore.Config) error
    
    // 启动服务
    Start() error
    
    // 停止服务
    Stop()
    
    // 优雅停止服务
    GracefulStop()
    
    // 设置日志记录器
    SetLog(logger gcore.Logger)
    
    // 注入监听器（用于热重载）
    InjectListeners(listeners []net.Listener)
    
    // 获取监听器（用于热重载）
    Listeners() []net.Listener
}
```

## 配置文件

默认配置文件 `app.toml` 示例：

```toml
[app]
# 应用名称
name = "myapp"
# 应用版本
version = "1.0.0"

[log]
# 日志级别: debug, info, warn, error
level = "info"
# 日志输出: stdout, stderr, file
output = "stdout"

[database]
default = "default"

[[database.mysql]]
enable = true
id = "default"
host = "127.0.0.1"
port = 3306
username = "root"
password = ""
database = "test"
prefix = ""
charset = "utf8mb4"
collate = "utf8mb4_general_ci"
maxidle = 10
maxconns = 100
timeout = 3000

[cache]
default = "default"

[[cache.redis]]
enable = true
id = "default"
address = "127.0.0.1:6379"
password = ""
dbnum = 0
```

## API 参考

### 创建应用

```go
// 创建新应用
func New() gcore.App

// 创建默认应用（自动加载配置）
func Default() gcore.App

// 工厂函数
func NewApp(isDefault bool) gcore.App
```

### 应用方法

```go
// 设置配置
SetConfig(cfg gcore.Config)

// 获取配置
Config() gcore.Config

// 添加服务
AddService(item gcore.ServiceItem)

// 运行应用
Run() error

// 设置阻塞模式
SetBlock(isBlock bool)

// 添加运行前回调
OnRun(fn func(gcore.Config) error)

// 添加关闭回调
OnShutdown(fn func())

// 附加额外配置
AttachConfig(cfg gcore.Config) string
AttachConfigFile(file string) (string, error)

// 获取上下文
Ctx() gcore.Ctx

// 设置上下文
SetCtx(ctx gcore.Ctx)

// 获取日志记录器
Logger(name string) gcore.Logger
```

## 使用场景

1. **Web 应用**：构建 HTTP/HTTPS 服务器
2. **API 服务**：构建 RESTful API 服务
3. **微服务**：构建微服务应用
4. **TCP 服务**：构建 TCP 服务器
5. **定时任务**：管理定时任务服务
6. **混合服务**：同时运行多种服务

## 最佳实践

### 1. 优雅关闭

```go
app.OnShutdown(func() {
    // 关闭数据库连接
    db := gmc.DB.DB()
    if db != nil {
        db.Close()
    }
    
    // 关闭缓存连接
    cache := gmc.Cache.Cache()
    if cache != nil {
        cache.Close()
    }
    
    // 等待异步日志写入完成
    logger := app.Logger("")
    if logger.Async() {
        logger.WaitAsyncDone()
    }
})
```

### 2. 错误处理

```go
err := app.Run()
if err != nil {
    log.Printf("Application error: %v", err)
    os.Exit(1)
}
```

### 3. 配置管理

```go
// 使用环境变量覆盖配置
// 设置环境变量前缀
os.Setenv("ENV_PREFIX", "MYAPP")
// 环境变量 MYAPP_DATABASE_HOST 会覆盖 database.host
```

## 注意事项

1. **热重载限制**：仅支持 Linux 系统
2. **监听器传递**：服务必须正确实现 `InjectListeners` 和 `Listeners` 方法才能支持热重载
3. **优雅关闭**：确保服务实现 `GracefulStop` 方法，避免数据丢失
4. **配置文件**：默认搜索路径为 `.`、`conf`、`config` 目录
5. **阻塞模式**：默认阻塞运行，可通过 `SetBlock(false)` 改为非阻塞

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [配置模块文档](../config/)
- [日志模块文档](../log/)
- [数据库模块文档](../db/)