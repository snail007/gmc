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

### 应用生命周期与钩子

GMC App 的生命周期设计精巧，通过一系列钩子方法，你可以在应用的不同阶段注入自定义逻辑。下面是完整的生命周期流程和钩子说明。

#### 启动流程 (Startup Sequence)

当调用 `app.Run()` 时，应用按以下顺序启动：

```
+------------------------------------+
| 1. 解析配置文件 (app.toml)         |
+------------------------------------+
                 |
                 ▼
+------------------------------------+
| 2. 初始化核心组件                  |
|    (Log, DB, Cache, i18n)        |
+------------------------------------+
                 |
                 ▼
+------------------------------------+
| 3. 执行 app.OnRun() 钩子           |
+------------------------------------+
                 |
                 ▼
+------------------------------------+
| 4. 遍历并初始化所有服务 (Services) |
|    +----------------------------+  |
|    | a. 调用 service.BeforeInit |  |
|    +----------------------------+  |
|    | b. 调用 service.Init()     |  |
|    +----------------------------+  |
|    | c. 调用 service.AfterInit  |  |
|    +----------------------------+  |
|    | d. 调用 service.Start()    |  |
|    +----------------------------+  |
+------------------------------------+
                 |
                 ▼
+------------------------------------+
| 5. 监听并等待关闭信号              |
|    (阻塞运行)                    |
+------------------------------------+
```

#### 关闭流程 (Shutdown Sequence)

当应用接收到关闭信号 (如 `Ctrl+C`) 时，`app.Stop()` 被调用：

1.  **执行 `app.OnShutdown()` 钩子**：按注册顺序执行所有应用级别的关闭钩子，用于资源清理等。
2.  **停止所有服务**：遍历所有服务，调用每个服务的 `Stop()` 方法。

#### 热重载流程 (Hot Reload Sequence)

在非 Windows 系统上，当应用接收到 `SIGUSR2` 信号时：

1.  **旧进程**：获取所有服务的网络监听器 (`Listeners`)，并启动一个带监听器文件描述符的新子进程。
2.  **新进程**：通过 `InjectListeners` 方法接收网络监听器，并完整地执行一遍**启动流程**。
3.  **旧进程**：在新进程成功启动后，调用所有服务的 `GracefulStop()` 方法优雅退出，实现零停机更新。

#### 钩子方法详解

-   **`app.OnRun(func(gcore.Config) error)`**
    -   **触发时机**：在核心组件 (log, db, cache) 初始化之后，但在任何服务 (`Service`) 初始化之前。
    -   **主要用途**：进行应用级别的初始化检查、数据准备等。

-   **`app.OnShutdown(func())`**
    -   **触发时机**：在应用关闭流程开始时，在任何服务被停止之前。
    -   **主要用途**：执行最终的清理工作，如关闭数据库连接、等待异步日志刷盘等。

-   **`service.BeforeInit(func(*gcore.ServiceItem, gcore.Config) error)`**
    -   **触发时机**：在单个服务 `Init()` 方法被调用之前。
    -   **主要用途**：对特定服务进行预配置或检查。

-   **`service.AfterInit(func(*gcore.ServiceItem) error)`**
    -   **触发时机**：在单个服务 `Init()` 方法成功调用之后，但在 `Start()` 之前。
    -   **主要用途**：在服务初始化后、启动前进行一些操作，例如注册路由、添加中间件等，这是最常用的服务钩子。

#### 完整示例

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

func main() {
    // 1. 创建 App
    app := gmc.New.AppDefault()

    // 2. 注册 App 级别的钩子
    app.OnRun(func(cfg gcore.Config) error {
        fmt.Println("-> Hook: app.OnRun")
        // 此时可以访问已初始化的配置和核心组件
        return nil
    })

    app.OnShutdown(func() {
        fmt.Println("-> Hook: app.OnShutdown")
        // 执行最后的清理工作
    })

    // 3. 添加服务并注册服务级别的钩子
    app.AddService(gcore.ServiceItem{
        Service: gmc.New.HTTPServer(app.Ctx()).(gcore.Service),
        BeforeInit: func(s *gcore.ServiceItem, cfg gcore.Config) error {
            fmt.Println("--> Hook: service.BeforeInit")
            return nil
        },
        AfterInit: func(s *gcore.ServiceItem) error {
            fmt.Println("--> Hook: service.AfterInit")
            // 这是注册路由和中间件的最佳位置
            httpServer := s.Service.(gcore.HTTPServer)
            httpServer.Router().GET("/", func(ctx gmc.Ctx) {
                ctx.Write("Hello GMC!")
            })
            return nil
        },
    })

    // 4. 运行 App
    fmt.Println("App is starting...")
    if err := app.Run(); err != nil {
        panic(err)
    }
}

// 预期的输出顺序:
// App is starting...
// -> Hook: app.OnRun
// --> Hook: service.BeforeInit
// (服务自身的 Init 方法被调用)
// --> Hook: service.AfterInit
// (服务自身的 Start 方法被调用)
// gmc app started done.
// (应用阻塞运行，等待关闭信号...)
// (接收到关闭信号后)
// -> Hook: app.OnShutdown
// (服务自身的 Stop 方法被调用)
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

提示：如果使用 `gmc.New.AppDefault()` 创建应用，并且希望运行 `APIServer`，则需要在 `app.toml` 中添加 `[apiserver]` 配置块。

默认配置文件 `app.toml` 示例：

```toml
# GMC 默认配置文件 app.toml

############################################################
# HTTP 服务配置
############################################################
[httpserver]
listen=":7080"
tlsenable=false
tlscert="conf/server.crt"
tlskey="conf/server.key"
tlsclientauth=false
tlsclientsca="./conf/clintsca.crt"
printroute=true
showerrorstack=true

############################################################
# 静态文件服务配置
############################################################
[static]
dir="static"
urlpath="/static/"

#############################################################
# 日志配置
#############################################################
[log]
level=0
output=[0,1]
dir="./logs"
archive_dir=""
filename="web_%Y%m%d.log"
gzip=true
async=true

#############################################################
# i18n (国际化) 配置
#############################################################
[i18n]
enable=false
dir="i18n"
default="zh-CN"

#############################################################
# 视图/模板配置
#############################################################
[template]
dir="views"
ext=".html"
delimiterleft="{{"
delimiterright="}}"
layout="layout"

########################################################
# Session 配置
########################################################
[session]
enable=true
store="memory"
cookiename="gmcsid"
ttl=3600

[session.file]
dir="{tmp}"
gctime=300
prefix=".gmcsession_"

[session.memory]
gctime=300

[session.redis]
debug=false
address="127.0.0.1:6379"
prefix=""
password=""
timeout=10
dbnum=0
maxidle=10
maxactive=30
idletimeout=300
maxconnlifetime=3600
wait=false

############################################################
# 缓存配置
############################################################
[cache]
default="redis"

[[cache.redis]]
debug=true
enable=true
id="default"
address="127.0.0.1:6379"
prefix=""
password=""
timeout=10
dbnum=0
maxidle=10
maxactive=30
idletimeout=300
maxconnlifetime=3600
wait=false

[[cache.memory]]
enable=true
id="default"
cleanupinterval=30

[[cache.file]]
enable=true
id="default"
dir="{tmp}"
cleanupinterval=30

########################################################
# 数据库配置
########################################################
[database]
default="mysql"

[[database.mysql]]
enable=true
id="default"
host="127.0.0.1"
port="3306"
username="user"
password="user"
database="test"
prefix=""
prefix_sql_holder="__PREFIX__"
charset="utf8"
collate="utf8_general_ci"
maxidle=30
maxconns=200
timeout=15000
readtimeout=15000
writetimeout=15000
maxlifetimeseconds=1800

[[database.sqlite3]]
enable=false
id="default"
database="test.db"
password=""
prefix=""
prefix_sql_holder="__PREFIX__"
syncmode=0
openmode="rw"
cachemode="shared"

##############################################################
# 访问日志中间件配置
##############################################################
[accesslog]
dir = "./logs"
archive_dir = ""
filename="access_%Y%m%d.log"
gzip=true
format="$req_time $client_ip $host $uri?$query $status_code ${time_used}ms"

##############################################################
# 前端代理配置
##############################################################
[frontend]
#type="proxy"
#ips=["192.168.1.1","192.168.0.0/16"]
#header=""


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