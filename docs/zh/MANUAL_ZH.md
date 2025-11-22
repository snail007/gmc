# GMC 框架完整使用手册

<p align="center">
  <img src="https://raw.githubusercontent.com/snail007/gmc/master/doc/images/logo2.png" width="300" alt="GMC Logo"/>
</p>

<p align="center">
  <a href="https://github.com/snail007/gmc/actions"><img src="https://github.com/snail007/gmc/workflows/build/badge.svg" alt="Build Status"></a>
  <a href="https://codecov.io/gh/snail007/gmc"><img src="https://codecov.io/gh/snail007/gmc/branch/master/graph/badge.svg" alt="codecov"></a>
  <a href="https://goreportcard.com/report/github.com/snail007/gmc"><img src="https://goreportcard.com/badge/github.com/snail007/gmc" alt="Go Report"></a>
  <a href="https://pkg.go.dev/github.com/snail007/gmc"><img src="https://img.shields.io/badge/go.dev-reference-blue" alt="API Reference"></a>
</p>

> **📖 语言切换**: **中文** | [🌐 English Version](../MANUAL.md)
> 
> **说明**: 本文档为完整的中文版本，包含所有功能和详细说明。英文版本提供核心功能和基础用法。

## 什么是 GMC？

GMC（Go Micro Controller）是一个智能、灵活、高性能的 Golang Web 和 API 开发框架。GMC 的目标是实现高性能、高生产力，让开发者用更少的代码完成更多的事情。

### 核心特性

- 🚀 **高性能**: 基于高效的路由引擎，性能卓越
- 🎯 **简单易用**: 直观的 API 设计，学习曲线平缓
- 🔧 **强大工具链**: 提供完整的开发工具，一键生成项目
- 📦 **模块化设计**: 纯接口抽象，组件可随意替换
- 🔄 **热编译**: 开发时自动编译重启，提升开发效率
- 💾 **资源嵌入**: 支持将静态文件、视图、i18n 文件打包进二进制
- 🌍 **国际化支持**: 内置完整的多语言解决方案
- 📝 **丰富的文档**: 详细的文档和示例代码

### 为什么选择 GMC？

1. **完整的解决方案**: 从项目生成、开发、测试到部署，提供全生命周期支持
2. **最佳实践内置**: 框架设计遵循 Go 语言和 Web 开发最佳实践
3. **灵活可扩展**: Provider 模式让你可以轻松替换任何组件
4. **活跃的社区**: 持续维护和更新，快速响应问题

---

## 快速开始

### 环境要求

- Go 1.16 或更高版本
- 支持的操作系统: Linux、macOS、Windows

### 安装 GMC

使用 `go get` 命令安装 GMC 框架：

```bash
go get -u github.com/snail007/gmc
```

### 安装 GMCT 工具链

GMCT 是 GMC 的配套工具链，提供项目生成、热编译等功能：

```bash
# 安装 gmct
go install github.com/snail007/gmc/tool/gmct@latest

# 验证安装
gmct version
```

### 创建第一个 Web 项目

使用 GMCT 创建一个新的 Web 项目：

```bash
# 创建项目目录
mkdir -p $GOPATH/src/myapp
cd $GOPATH/src/myapp

# 初始化项目
gmct new web

# 或者指定项目类型
gmct new --type web
```

### 运行项目

```bash
# 开发模式运行（支持热编译）
gmct run

# 或直接运行
go run main.go
```

访问 `http://localhost:7080` 查看运行结果。

### 项目结构说明

生成的项目默认目录结构如下：

```text
myapp/
├── conf/
│   └── app.toml          # 主配置文件
├── controller/
│   └── demo.go           # 示例控制器
├── initialize/
│   └── initialize.go     # 初始化逻辑
├── router/
│   └── router.go         # 路由配置
├── static/
│   └── jquery.js         # 静态文件
├── views/
│   └── welcome.html      # 视图模板
├── go.mod                # Go 模块文件
├── go.sum                # Go 依赖锁定文件
├── grun.toml             # GMCT 运行配置
└── main.go               # 程序入口
```

#### 文件说明

- **conf/app.toml**: 项目的主配置文件，包含服务器、数据库、缓存等所有配置
- **controller/**: 存放控制器文件，处理业务逻辑
- **initialize/**: 项目初始化代码，如路由注册、服务配置等
- **router/**: 路由配置文件，定义 URL 和控制器的映射关系
- **static/**: 静态资源目录（JS、CSS、图片等）
- **views/**: 视图模板目录
- **main.go**: 应用程序入口，启动 Web 服务

### Hello World 示例

让我们创建一个简单的 Hello World 示例来了解 GMC 的基本使用。

#### 1. 创建控制器

编辑 `controller/demo.go`:

```go
package controller

import (
"github.com/snail007/gmc"
)

type Demo struct {
gmc.Controller
}

// Hello 方法会响应 /demo/hello 请求
func (this *Demo) Hello() {
this.Write("Hello GMC! 欢迎使用 GMC 框架！")
}

// JSON 响应示例
func (this *Demo) JsonDemo() {
data := map[string]interface{}{
"message": "Hello GMC",
"status":  "success",
"code":    200,
}
this.Ctx.JSON(200, data)
}

// 获取参数示例
func (this *Demo) GetParams() {
name := this.Ctx.GET("name", "Guest")
this.Write("Hello, " + name + "!")
}
```

#### 2. 配置路由

编辑 `router/router.go`:

```go
package router

import (
"myapp/controller"
"github.com/snail007/gmc"
)

func Init(s gmc.HTTPServer) {
// 获取路由对象
r := s.Router()

// 绑定控制器
r.Controller("/demo", new(controller.Demo))

// 或者绑定单个方法
r.ControllerMethod("/", new(controller.Demo), "Index")
}
```

#### 3. 运行并测试

```bash
# 运行项目
gmct run

# 在另一个终端测试
curl http://localhost:7080/demo/hello
# 输出: Hello GMC! 欢迎使用 GMC 框架！

curl http://localhost:7080/demo/jsondemo
# 输出: {"code":200,"message":"Hello GMC","status":"success"}

curl http://localhost:7080/demo/getparams?name=张三
# 输出: Hello, 张三!
```

> **HTTP Server 详细文档：** [http/server/README.md](https://github.com/snail007/gmc/blob/master/http/server/README.md) - 查看完整的服务器配置、TLS 设置、性能优化等

---

## 核心概念

GMC 采用分层架构、Provider 模式和完整的生命周期管理，提供灵活可扩展的应用框架。

> **详细文档：** [module/app/README.md](https://github.com/snail007/gmc/blob/master/module/app/README.md) - 查看完整的应用架构、生命周期管理、热重载等

### 架构设计

GMC 采用分层架构设计，主要包含以下几层：

```
┌─────────────────────────────────────┐
│         Application Layer           │
│        (App, Service, Ctx)          │
└──────────────┬──────────────────────┘
               │
┌──────────────┴──────────────────────┐
│       HTTP/API Server Layer         │
│   (HTTPServer, APIServer, Router)   │
└──────────────┬──────────────────────┘
               │
┌──────────────┴──────────────────────┐
│      Business Logic Layer           │
│  (Controller, Handler, Middleware)  │
└──────────────┬──────────────────────┘
               │
┌──────────────┴──────────────────────┐
│       Data Access Layer             │
│  (Database, Cache, Session, etc.)   │
└─────────────────────────────────────┘
```

### Provider 模式

GMC 使用 Provider 模式来管理组件的创建和注册。这种模式的优势：

1. **松耦合**: 业务代码不直接依赖具体实现
2. **可测试**: 方便进行单元测试时替换为 Mock 对象
3. **可扩展**: 轻松切换或添加新的实现

#### Provider 注册

```go
// 注册自定义缓存 Provider
gcore.RegisterCache("redis", func(ctx gcore.Ctx) (gcore.Cache, error) {
    // 创建并返回 Redis 缓存实例
    return NewRedisCache(ctx.Config()), nil
})

// 注册自定义日志 Provider
gcore.RegisterLogger("mylogger", func(ctx gcore.Ctx, prefix string) gcore.Logger {
    return NewMyLogger(prefix)
})
```

#### Provider 使用

```go
// 获取已注册的 Provider
cacheProvider := gcore.ProviderCache()
cache, err := cacheProvider(ctx)

// 或使用默认 Provider
logger := gcore.ProviderLogger()(ctx, "myapp")
logger.Info("应用启动")
```

### 生命周期

GMC 应用的生命周期包含以下阶段：

```
┌──────────────┐
│  创建 App    │
└──────┬───────┘
       │
┌──────▼───────┐
│  配置加载    │
└──────┬───────┘
       │
┌──────▼───────┐
│  服务初始化  │  (BeforeInit Hook)
└──────┬───────┘
       │
┌──────▼───────┐
│  启动服务    │  (Start)
└──────┬───────┘
       │
┌──────▼───────┐
│  运行中      │
└──────┬───────┘
       │
┌──────▼───────┐
│  停止服务    │  (GracefulStop)
└──────────────┘
```

#### 生命周期钩子

```go
// 创建应用
app := gmc.New.App()

// 设置配置文件
app.SetConfigFile("conf/app.toml")

// OnRun 钩子：在服务启动前执行
app.OnRun(func(cfg gcore.Config) error {
    // 初始化数据库连接
    // 注册路由
    // 其他初始化操作
    return nil
})

// OnShutdown 钩子：在服务停止时执行
app.OnShutdown(func() {
    // 关闭数据库连接
    // 清理资源
    fmt.Println("应用正在关闭...")
})

// 启动应用
app.Run()
```

### 资源嵌入

**推荐方式：使用 Go embed 功能**

GMC 推荐使用 Go 1.16+ 原生的 `embed` 功能将静态资源和视图模板直接打包到二进制文件中，实现单文件部署。这是标准、类型安全的方式。

> **⚠️ 注意：** 不再推荐使用 `gmct tpl`、`gmct static`、`gmct i18n` 等打包命令。请使用下面介绍的 `embed` 方式。

**embed 的优势：**
- ✅ Go 原生功能，无需额外工具
- ✅ 编译时类型检查
- ✅ IDE 支持完善
- ✅ 标准化、易维护

关键在于，你需要**显式导入**包含 `embed.FS` 变量的包，并在代码中**直接使用**这些变量。

> **资源打包指南-详细文档：** [资源打包指南](https://github.com/snail007/gmc/blob/master/doc/资源打包指南.md) - 查看三种资源打包的完整说明


#### 嵌入静态文件

> **详细文档：** [http/server/README.md](https://github.com/snail007/gmc/blob/master/http/server/README.md) - 查看静态文件服务和嵌入的完整说明

1.  在 `static` 文件夹中创建 `static.go` 文件，并导出一个 `embed.FS` 变量：

```go
package static

import (
	"embed"
)

//go:embed *
var StaticFS embed.FS
```

**重要提示**: 当使用 `go:embed` 嵌入静态资源时，为了避免框架优先从本地目录加载文件，应将 `app.toml` 中 `[static]` 配置块下的 `dir` 设置为空，即 `dir = ""`。


#### 嵌入视图文件

> **详细文档：** 
> - [http/template/README.md](https://github.com/snail007/gmc/blob/master/http/template/README.md) - 模板引擎完整说明
> - [http/view/README.md](https://github.com/snail007/gmc/blob/master/http/view/README.md) - 视图渲染完整文档

1.  在 `views` 文件夹中创建 `views.go` 文件，并导出一个 `embed.FS` 变量：

```go
package views

import (
	"embed"
)

//go:embed *
var ViewFS embed.FS
```

#### 嵌入 i18n 文件

> **详细文档：** [module/i18n/README.md](https://github.com/snail007/gmc/blob/master/module/i18n/README.md) - 查看国际化完整使用指南

GMC 提供了简单的 API 来嵌入 i18n 国际化文件：

1.  在 `i18n` 文件夹中创建 `i18n.go` 文件：

```go
package i18n

import "embed"

//go:embed *.toml
var I18nFS embed.FS
```

2.  在 `main.go` 中初始化：

```go
import (
    gi18n "github.com/snail007/gmc/module/i18n"
    "myapp/i18n"
)

func main() {
    // 初始化嵌入的 i18n 文件，设置默认语言
    err := gi18n.InitEmbedFS(i18n.I18nFS, "zh-CN")
    if err != nil {
        panic(err)
    }
    
    // 继续初始化应用...
}
```

**重要提示**: 使用 `InitEmbedFS` 时，应将 `app.toml` 中 `[i18n]` 的 `enable` 设置为 `false`。

查看 [i18n 模块文档](https://github.com/snail007/gmc/blob/master/module/i18n/README.md) 了解更多详情。

#### 完整示例

下面我们提供两种方式来初始化并使用嵌入的资源（包含静态文件、视图和 i18n）。

**通用文件:**

*   **项目结构:**

    ```text
    /myapp
    ├── go.mod
    ├── static/
    │   ├── css/
    │   │   └── style.css
    │   └── static.go
    ├── views/
    │   ├── index.html
    │   └── views.go
    ├── i18n/
    │   ├── zh-CN.toml
    │   ├── en-US.toml
    │   └── i18n.go
    └── main.go
    ```

*   **`static/static.go`:**

    ```go
    package static
    import "embed"

    //go:embed *
    var StaticFS embed.FS
    ```

*   **`views/views.go`:**

    ```go
    package views
    import "embed"

    //go:embed *
    var ViewFS embed.FS
    ```

*   **`i18n/i18n.go`:**

    ```go
    package i18n
    import "embed"

    //go:embed *.toml
    var I18nFS embed.FS
    ```

---

**方式一：直接使用 HTTPServer (简单直接)**

*   **`main.go`**

    ```go
    package main

    import (
    	"github.com/snail007/gmc"
    	gtemplate "github.com/snail007/gmc/http/template"
    	gi18n "github.com/snail007/gmc/module/i18n"

    	// 显式导入 static、views 和 i18n 包
    	"myapp/static"
    	"myapp/views"
    	"myapp/i18n"
    )

    func main() {
    	// 1. 初始化嵌入的 i18n 文件
    	err := gi18n.InitEmbedFS(i18n.I18nFS, "zh-CN")
    	if err != nil {
    		panic(err)
    	}

    	// 2. 创建一个 HTTP 服务器
    	s := gmc.New.HTTPServer(gmc.New.CtxDefault())

    	// 3. 注册嵌入的静态文件
    	s.ServeEmbedFS(static.StaticFS, "/static")

    	// 4. 注册嵌入的视图文件
    	tpl := gtemplate.NewEmbedTemplateFS(s.Tpl(), views.ViewFS, ".")
    	if err := tpl.Parse(); err != nil {
    		s.Logger().Panicf("解析模板失败: %s", err)
    	}

    	// 5. 设置路由并启动
    	s.Router().GET("/", func(ctx gmc.Ctx) {
    		ctx.View.Render("index.html")
    	})
    	s.Run()
    }
    ```

---

**方式二：使用 App 管理服务 (推荐用于复杂应用)**

*   **`main.go`**

    ```go
    package main

    import (
    	"github.com/snail007/gmc"
    	gcore "github.com/snail007/gmc/core"
    	gtemplate "github.com/snail007/gmc/http/template"
    	gi18n "github.com/snail007/gmc/module/i18n"

    	// 显式导入 static、views 和 i18n 包
    	"myapp/static"
    	"myapp/views"
    	"myapp/i18n"
    )

    func main() {
    	// 1. 初始化嵌入的 i18n 文件
    	err := gi18n.InitEmbedFS(i18n.I18nFS, "zh-CN")
    	if err != nil {
    		panic(err)
    	}

    	// 2. 创建应用
    	app := gmc.New.App()
    	app.AddService(gcore.ServiceItem{
    		Service: gmc.New.HTTPServer(app.Ctx()).(gcore.Service),
    		AfterInit: func(s *gcore.ServiceItem) (err error) {
    			httpServer := s.Service.(*gmc.HTTPServer)

    			// 注册静态文件，直接使用导入的 static.StaticFS
    			httpServer.ServeEmbedFS(static.StaticFS, "/static")

    			// 注册视图文件，直接使用导入的 views.ViewFS
    			tpl := gtemplate.NewEmbedTemplateFS(httpServer.Tpl(), views.ViewFS, ".")
    			if err = tpl.Parse(); err != nil {
    				return
    			}

    			// 注册路由
    			httpServer.Router().GET("/", func(ctx gmc.Ctx) {
    				ctx.View.Render("index.html")
    			})
    			return
    		},
    	})
    	app.Run()
    }
    ```
---

## 配置

GMC 使用强大的配置管理模块，基于 Viper 封装，支持多种配置格式、环境变量、配置热加载等特性。

> **详细文档：** [module/config/README.md](https://github.com/snail007/gmc/blob/master/module/config/README.md) - 查看完整的 API 文档、高级用法和最佳实践

### 配置文件

GMC 使用 TOML 格式的配置文件。默认配置文件是 `conf/app.toml`。

#### 基本配置结构

```toml
# GMC 默认配置文件 app.toml

############################################################
# HTTP 服务配置
############################################################
[httpserver]
# 监听地址和端口
listen=":7080"
# 是否启用 TLS (HTTPS)
tlsenable=false
# TLS 证书文件路径
tlscert="conf/server.crt"
# TLS 密钥文件路径
tlskey="conf/server.key"
# 是否开启客户端证书认证 (双向TLS)
tlsclientauth=false
# 客户端 CA 证书路径
tlsclientsca="./conf/clintsca.crt"
# 是否在启动时打印路由表
printroute=true
# 是否在发生 panic 时在浏览器中显示错误和调用栈
showerrorstack=true

############################################################
# 静态文件服务配置 (当不使用 embed 嵌入时)
############################################################
[static]
# 静态文件目录的本地路径
dir="static"
# 访问静态文件的 URL 前缀
urlpath="/static/"

#############################################################
# 日志配置
#############################################################
[log]
# 日志级别: 1-7 分别对应 TRACE, DEBUG, INFO, WARN, ERROR, PANIC, NONE
# 7 表示不输出任何日志
level=3 # 默认为 INFO
# 日志输出目标: 0 表示控制台, 1 表示文件
output=[0,1]
# 日志文件存放目录 (仅当 output 包含 1 时有效)
dir="./logs"
# 归档目录，如果设置，过期的日志文件会被移动到这里
archive_dir=""
# 日志文件名，支持占位符: %Y(年), %m(月), %d(日), %H(时)
filename="web_%Y%m%d.log"
# 是否启用 gzip 压缩日志文件
gzip=true
# 是否开启异步日志，开启后需要确保在程序退出前调用 logger.WaitAsyncDone()
async=true

#############################################################
# i18n (国际化) 配置
#############################################################
[i18n]
# 是否启用 i18n
enable=false
# 语言文件目录
dir="i18n"
# 默认语言 (文件名，不含扩展名，如 zh-CN.toml)
default="zh-CN"

#############################################################
# 视图/模板配置
#############################################################
[template]
# 模板文件目录 (当不使用 embed 嵌入时)
dir="views"
# 模板文件扩展名
ext=".html"
# 模板语法分隔符
delimiterleft="{{"
delimiterright="}}"
# 布局(layout)文件所在的子目录名
layout="layout"

########################################################
# Session 配置
########################################################
[session]
# 是否启用 Session
enable=true
# 存储引擎: "file", "memory", "redis"
store="memory"
# Session ID 存储在 Cookie 中的名称
cookiename="gmcsid"
# Session 过期时间 (秒)
ttl=3600

# 文件存储引擎配置
[session.file]
# Session 文件存放目录, {tmp} 是系统临时目录的占位符
dir="{tmp}"
# GC (垃圾回收) 周期 (秒)
gctime=300
# Session 文件前缀
prefix=".gmcsession_"

# 内存存储引擎配置
[session.memory]
# GC 周期 (秒)
gctime=300

# Redis 存储引擎配置
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
# 默认使用的缓存实例 ID
default="default"

# Redis 缓存实例配置 (可以有多个)
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

# 内存缓存实例配置
[[cache.memory]]
enable=true
id="default"
# 清理周期 (秒)
cleanupinterval=30

# 文件缓存实例配置
[[cache.file]]
enable=true
id="default"
# 缓存目录, {tmp} 是系统临时目录的占位符
dir="{tmp}"
# 清理周期 (秒)
cleanupinterval=30

########################################################
# 数据库配置
########################################################
[database]
# 默认使用的数据库实例 ID
default="default"

# MySQL 实例配置 (可以有多个)
[[database.mysql]]
enable=true
id="default"
host="127.0.0.1"
port="3306"
username="user"
password="user"
database="test"
# 表前缀
prefix=""
# SQL语句中表前缀的占位符
prefix_sql_holder="__PREFIX__"
charset="utf8"
collate="utf8_general_ci"
maxidle=30
maxconns=200
timeout=15000
readtimeout=15000
writetimeout=15000
maxlifetimeseconds=1800

# SQLite 实例配置
[[database.sqlite3]]
enable=false
id="default"
database="test.db"
# 如果密码不为空，数据库将被加密
password=""
prefix=""
prefix_sql_holder="__PREFIX__"
# 同步模式: 0:OFF, 1:NORMAL, 2:FULL, 3:EXTRA
syncmode=0
# 打开模式: ro,rw,rwc,memory
openmode="rw"
# 缓存模式: shared,private
cachemode="shared"

# SQLite 实例配置
[[database.sqlite3]]
enable=false
id="default"
database="test.db"
# 如果密码不为空，数据库将被加密
password=""
prefix=""
prefix_sql_holder="__PREFIX__"
# 同步模式: 0:OFF, 1:NORMAL, 2:FULL, 3:EXTRA
syncmode=0
# 打开模式: ro,rw,rwc,memory
openmode="rw"
# 缓存模式: shared,private
cachemode="shared"

##############################################################
# Web & API 访问日志中间件配置
##############################################################
[accesslog]
dir = "./logs"
archive_dir = ""
# 日志文件名，支持占位符
filename="access_%Y%m%d.log"
gzip=true
# 日志格式, 可用占位符:
# $host: URL中的主机名(含端口)
# $uri: 请求路径
# $query: URL中的查询字符串
# $status_code: 响应的 HTTP 状态码
# $time_used: 请求处理耗时(毫秒)
# $req_time: 请求时间, 格式: 2020-10-55 15:33:55
# $client_ip: 客户端真实IP
# $remote_addr: 客户端地址(含端口)
# $local_addr: 服务端被访问的地址
format="$req_time $client_ip $host $uri?$query $status_code ${time_used}ms"

##############################################################
# 前端代理配置 (用于安全地获取客户端IP)
##############################################################
[frontend]
# 代理类型: "cloudflare", "proxy"
# 当类型为 cloudflare, gmc 会自动获取 Cloudflare 的 IP 段来验证请求头
# 当类型为 proxy, 你需要手动在下面的 ips 字段中提供你的代理服务器IP地址
#type="proxy"
# 代理服务器的 IP 或 CIDR 地址段
#ips=["192.168.1.1","192.168.0.0/16"]
# 用于获取真实IP的请求头字段
# cloudflare 可用: True-Client-IP, CF-Connecting-IP (默认)
# proxy 可用: X-Real-IP, X-Forwarded-For (默认)
#header=""
```

#### API 服务配置示例 (api.toml)

对于纯 API 服务，配置可以更精简。如果使用默认应用(`gmc.New.AppDefault()`)并希望运行 `APIServer`，需要在 `app.toml` 中添加 `[apiserver]` 配置块。

```toml
# GMC API服务配置文件 api.toml

############################################################
# API 服务配置
############################################################
[apiserver]
# 监听地址和端口
listen=":7081"
# 是否在启动时打印路由表
printroute=true
# 是否在发生 panic 时显示错误和调用栈
showerrorstack=true

#############################################################
# 日志配置
#############################################################
[log]
# 日志级别: 1-7 (INFO, WARN, ERROR 等)
level=3
# 输出目标: 0-控制台
output=[0]

############################################################
# 缓存配置 (按需启用)
############################################################
[cache]
default="default"

[[cache.redis]]
enable=false
id="default"
address="127.0.0.1:6379"

########################################################
# 数据库配置 (按需启用)
########################################################
[database]
default="default"

[[database.mysql]]
enable=false
id="default"
host="127.0.0.1"
port="3306"
username="user"
password="user"
database="test"

[[database.sqlite3]]
enable=false
id="default"
database="test.db"
password=""
```


### 配置加载

#### 加载主配置

```go
// 设置配置文件路径
app.SetConfigFile("conf/app.toml")

// 或者直接设置配置对象
cfg := gconfig.New()
cfg.SetConfigFile("conf/app.toml")
cfg.ReadInConfig()
app.SetConfig(cfg)
```

#### 附加额外配置

```go
// 附加其他配置文件
app.AttachConfigFile("database", "conf/database.toml")
app.AttachConfigFile("redis", "conf/redis.toml")

// 使用附加配置
dbCfg := app.Config("database")
host := dbCfg.GetString("host")
```

### 读取配置

#### 在代码中读取配置

```go
// 在控制器中
func (this *Demo) Index() {
    cfg := this.Config
    
    // 读取字符串
    appName := cfg.GetString("app.name")
    
    // 读取整数
    port := cfg.GetInt("httpserver.port")
    
    // 读取布尔值
    debug := cfg.GetBool("app.debug")
    
    // 读取子配置
    dbCfg := cfg.Sub("database")
    driver := dbCfg.GetString("driver")
    
    // 设置默认值
    timeout := cfg.GetInt("app.timeout", 30)
}
```

#### 配置类型转换

```go
// 基本类型
stringVal := cfg.GetString("key")
intVal := cfg.GetInt("key")
int64Val := cfg.GetInt64("key")
floatVal := cfg.GetFloat64("key")
boolVal := cfg.GetBool("key")

// 时间类型
duration := cfg.GetDuration("timeout") // 如: "30s", "5m"
timeVal := cfg.GetTime("start_time")

// 切片类型
intSlice := cfg.GetIntSlice("ports")
stringSlice := cfg.GetStringSlice("hosts")

// Map 类型
stringMap := cfg.GetStringMap("database")
stringMapString := cfg.GetStringMapString("headers")
```

### 环境变量

GMC 支持通过环境变量覆盖配置：

```go
// 启用自动环境变量
cfg.AutomaticEnv()

// 设置环境变量前缀
cfg.SetEnvPrefix("MYAPP")

// 绑定特定环境变量
cfg.BindEnv("database.host", "DB_HOST")
cfg.BindEnv("database.port", "DB_PORT")
```

使用示例：

```bash
# 设置环境变量
export MYAPP_DATABASE_HOST=192.168.1.100
export MYAPP_DATABASE_PORT=3306

# 运行应用
./myapp
```

### 配置最佳实践

1. **分离环境配置**: 开发、测试、生产使用不同的配置文件
2. **敏感信息**: 数据库密码等敏感信息使用环境变量
3. **配置验证**: 启动时验证必要的配置项
4. **合理默认值**: 为可选配置提供合理的默认值

```go
// 配置验证示例
func validateConfig(cfg gcore.Config) error {
    if cfg.GetString("database.dsn") == "" {
        return errors.New("database.dsn is required")
    }
    
    if cfg.GetInt("httpserver.port") == 0 {
        return errors.New("httpserver.port must be set")
    }
    
    return nil
}
```

---

## 路由

GMC 提供灵活强大的路由系统，支持多种路由绑定方式、路由参数、路由组、中间件等特性。

> **详细文档：** [http/router/README.md](https://github.com/snail007/gmc/blob/master/http/router/README.md) - 查看完整的路由 API、高级路由模式和最佳实践

### 基础路由

GMC 提供灵活强大的路由系统，支持多种路由绑定方式。

#### 绑定控制器

```go
func InitRouter(s gmc.HTTPServer) {
    r := s.Router()
    
    // 绑定控制器，自动识别所有公开方法
    // 访问路径：/user/list, /user/create, /user/update 等
    r.Controller("/user", new(controller.User))
    
    // 带 URL 后缀的控制器
    // 访问路径：/api/list.json, /api/create.json
    r.Controller("/api", new(controller.API), ".json")
}
```

#### 绑定单个方法

```go
// 绑定控制器的特定方法
r.ControllerMethod("/", new(controller.Index), "Home")
r.ControllerMethod("/about", new(controller.Index), "About")
```

#### 绑定处理函数

```go
// Handle 函数签名: func(w http.ResponseWriter, r *http.Request, ps gcore.Params)
func Hello(w http.ResponseWriter, r *http.Request, ps gcore.Params) {
    w.Write([]byte("Hello GMC!"))
}

// 绑定到路由
r.GET("/hello", Hello)
r.POST("/hello", Hello)
r.Handle("GET", "/hello", Hello)
r.HandleAny("/hello", Hello) // 支持所有 HTTP 方法
```

#### HTTP 方法路由

```go
// RESTful 风格路由
r.GET("/users", ListUsers)           // 获取用户列表
r.POST("/users", CreateUser)         // 创建用户
r.GET("/users/:id", GetUser)         // 获取单个用户
r.PUT("/users/:id", UpdateUser)      // 更新用户
r.PATCH("/users/:id", PatchUser)     // 部分更新用户
r.DELETE("/users/:id", DeleteUser)   // 删除用户

// 其他 HTTP 方法
r.HEAD("/users", HeadUsers)
r.OPTIONS("/users", OptionsUsers)
```

### 路由参数

#### 命名参数

```go
// 定义带参数的路由
r.GET("/user/:name", func(w http.ResponseWriter, r *http.Request, ps gcore.Params) {
    // 获取参数
    name := ps.ByName("name")
    w.Write([]byte("Hello " + name))
})

// 多个参数
r.GET("/post/:category/:id", func(w http.ResponseWriter, r *http.Request, ps gcore.Params) {
    category := ps.ByName("category")
    id := ps.ByName("id")
    // 处理逻辑...
})
```

#### 在控制器中获取参数

```go
type User struct {
    gmc.Controller
}

func (this *User) Profile() {
    // 方法 1: 通过 Param 获取
    userID := this.Param.ByName("id")
    
    // 方法 2: 通过 Ctx 获取
    userID := this.Ctx.GetParam("id")
    
    this.Write("User ID: " + userID)
}

// 路由配置
// r.GET("/user/:id/profile", ...) 或
// r.Controller("/user/:id", new(controller.User))
```

#### 通配符参数

```go
// 捕获所有路径
r.GET("/files/*filepath", func(w http.ResponseWriter, r *http.Request, ps gcore.Params) {
    filepath := ps.ByName("filepath")
    // 例如: /files/docs/manual.pdf -> filepath = "docs/manual.pdf"
})

// 静态文件服务
r.ServeFiles("/static/*filepath", http.Dir("public"))
```

### 路由组

路由组允许你为一组路由共享相同的前缀和中间件。

#### 基本路由组

```go
// 创建 API 路由组
apiGroup := r.Group("/api")
{
    apiGroup.GET("/users", ListUsers)
    apiGroup.POST("/users", CreateUser)
    apiGroup.GET("/users/:id", GetUser)
}

// 等价于
// r.GET("/api/users", ListUsers)
// r.POST("/api/users", CreateUser)
// r.GET("/api/users/:id", GetUser)
```

#### 嵌套路由组

```go
api := r.Group("/api")
{
    v1 := api.Group("/v1")
    {
        v1.GET("/users", V1ListUsers)
        v1.POST("/users", V1CreateUser)
    }
    
    v2 := api.Group("/v2")
    {
        v2.GET("/users", V2ListUsers)
        v2.POST("/users", V2CreateUser)
    }
}

// 生成的路由:
// /api/v1/users
// /api/v2/users
```

#### 路由组与控制器

```go
// 为控制器设置命名空间
admin := r.Group("/admin")
admin.Controller("/user", new(controller.AdminUser))
admin.Controller("/post", new(controller.AdminPost))

// 访问路径：
// /admin/user/list
// /admin/post/create
```

### 中间件

中间件是在请求到达控制器之前或之后执行的代码，GMC 提供 4 个优先级的中间件层。

> **详细文档：** [module/middleware/README.md](https://github.com/snail007/gmc/blob/master/module/middleware/README.md) - 查看内置中间件、自定义中间件开发指南

#### 中间件架构

GMC 的中间件架构允许在请求处理的不同阶段插入自定义逻辑：

<p align="center">
  <img src="https://raw.githubusercontent.com/snail007/gmc/master/doc/images/http-and-api-server-architecture.png" alt="GMC Middleware Architecture" width="800"/>
</p>

如图所示，请求从客户端进入后，会依次经过不同优先级的中间件层，最终到达控制器处理，响应则按相反顺序返回。

#### 全局中间件

GMC 提供 4 个级别的中间件，按优先级从高到低：

```go
func InitMiddleware(s gmc.HTTPServer) {
    // 优先级 0 - 最高优先级
    s.AddMiddleware0(func(ctx gcore.Ctx) bool {
        // 记录请求开始时间
        ctx.Set("start_time", time.Now())
        return false // 返回 false 继续处理，true 则停止
    })
    
    // 优先级 1
    s.AddMiddleware1(AuthMiddleware)
    
    // 优先级 2
    s.AddMiddleware2(LogMiddleware)
    
    // 优先级 3 - 最低优先级
    s.AddMiddleware3(func(ctx gcore.Ctx) bool {
        // 在响应前添加自定义头
        ctx.Response().Header().Set("X-Custom-Header", "value")
        return false
    })
}
```

#### 编写中间件

```go
// 认证中间件
func AuthMiddleware(ctx gcore.Ctx) bool {
    token := ctx.Header("Authorization")
    
    if token == "" {
        ctx.WriteHeader(401)
        ctx.Write("Unauthorized")
        return true // 停止后续处理
    }
    
    // 验证 token
    user, err := ValidateToken(token)
    if err != nil {
        ctx.WriteHeader(401)
        ctx.JSON(401, map[string]string{
            "error": "Invalid token",
        })
        return true
    }
    
    // 将用户信息存储到上下文
    ctx.Set("user", user)
    return false // 继续处理
}

// 日志中间件
func LogMiddleware(ctx gcore.Ctx) bool {
    // 记录请求信息
    logger := ctx.Logger()
    logger.Infof("Request: %s %s from %s",
        ctx.Request().Method,
        ctx.Request().URL.Path,
        ctx.ClientIP(),
    )
    return false
}

// CORS 中间件
func CORSMiddleware(ctx gcore.Ctx) bool {
    ctx.SetHeader("Access-Control-Allow-Origin", "*")
    ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // 处理预检请求
    if ctx.Request().Method == "OPTIONS" {
        ctx.WriteHeader(204)
        return true
    }
    
    return false
}
```

#### 中间件使用场景

1. **认证授权**: 验证用户身份和权限
2. **日志记录**: 记录请求和响应信息
3. **性能监控**: 统计请求处理时间
4. **CORS 处理**: 处理跨域请求
5. **限流控制**: 防止 API 滥用
6. **数据压缩**: 压缩响应数据
7. **错误恢复**: 捕获 panic 并优雅处理

### 静态文件服务

```go
// 服务静态文件目录
r.ServeFiles("/static/*filepath", http.Dir("public"))

// 访问: http://localhost:7080/static/css/style.css
// 映射到: public/css/style.css
```

### 路由打印

查看所有已注册的路由：

```go
// 在初始化后打印路由表
r.PrintRouteTable(os.Stdout)

// 输出示例：
// GET    /                    controller.Index.Home
// GET    /user/:id            controller.User.Profile
// POST   /user/create         controller.User.Create
```

---

## 控制器

GMC 控制器提供完整的 HTTP 请求处理能力，包括生命周期钩子、请求解析、响应输出、视图渲染等功能。

> **详细文档：** [http/controller/README.md](https://github.com/snail007/gmc/blob/master/http/controller/README.md) - 查看完整的控制器 API、生命周期详解和高级用法

### 创建控制器

#### 基本控制器

```go
package controller

import (
    "github.com/snail007/gmc"
)

// User 用户控制器
type User struct {
    gmc.Controller
}

// List 用户列表页
func (this *User) List() {
    this.Write("用户列表")
}

// Detail 用户详情页
func (this *User) Detail() {
    userID := this.Param.ByName("id")
    this.Write("用户ID: " + userID)
}

// Create 创建用户
func (this *User) Create() {
    // POST 请求处理
    if !this.IsPOST() {
        this.Write("请使用 POST 方法")
        return
    }
    
    name := this.Ctx.POST("name")
    email := this.Ctx.POST("email")
    
    // 保存用户逻辑...
    
    this.Ctx.JSON(200, map[string]interface{}{
        "message": "用户创建成功",
        "data": map[string]string{
            "name": name,
            "email": email,
        },
    })
}
```

### 控制器规则

#### 方法命名规则

1. **公开方法**: 只有首字母大写的公开方法才能被路由访问
2. **忽略的方法**: 以 `_` 或 `__` 结尾的方法会被忽略
3. **保留方法名**: 以下方法名被 GMC 保留，不要使用：
   - `Before()`, `After()` - 生命周期钩子
   - `MethodCallPre()`, `MethodCallPost()` - 内部钩子
   - `Stop()`, `Die()` - 流程控制
   - `Write()`, `WriteE()` - 输出方法
   - 以 `Get` 开头的 getter 方法

```go
// ✅ 可以访问
func (this *User) Index() {}
func (this *User) UserList() {}

// ❌ 不能访问
func (this *User) index() {}      // 私有方法
func (this *User) Helper_() {}    // 以 _ 结尾
func (this *User) Private__() {}  // 以 __ 结尾

// ⚠️ 不要使用这些名称
func (this *User) Before() {}     // 生命周期钩子，有特殊用途
func (this *User) Write() {}      // 冲突
```

### 生命周期钩子

#### Before 方法

在控制器方法执行前调用，可用于：
- 权限验证
- 参数预处理
- 日志记录

```go
// Before 在所有方法执行前调用
func (this *User) Before() {
    // 检查用户登录状态
    if !this.IsLoggedIn() {
        this.Ctx.Redirect("/login")
        this.Stop() // 停止后续处理
        return
    }
    
    // 记录访问日志
    this.Logger.Infof("User访问: %s", this.Ctx.Request().URL.Path)
    
    // 设置公共数据
    this.Ctx.Set("current_user", this.GetCurrentUser())
}
```

#### After 方法

在控制器方法执行后调用，可用于：
- 资源清理
- 统一响应处理
- 性能统计

```go
// After 在所有方法执行后调用
func (this *User) After() {
    // 计算执行时间
    if startTime, exists := this.Ctx.Get("start_time"); exists {
        elapsed := time.Since(startTime.(time.Time))
        this.Logger.Infof("请求耗时: %v", elapsed)
    }
    
    // 清理资源
    this.CleanupResources()
}
```

#### 流程控制

```go
func (this *User) AdminOnly() {
    if !this.IsAdmin() {
        this.Ctx.JSON(403, map[string]string{
            "error": "需要管理员权限",
        })
        this.Stop() // 停止执行当前方法，但会执行 After()
        return
    }
    
    // 管理员操作...
}

func (this *User) Critical() {
    if !this.ValidateRequest() {
        this.Ctx.JSON(400, map[string]string{
            "error": "请求验证失败",
        })
        this.Die() // 停止所有后续处理，包括 After()
        return
    }
    
    // 关键操作...
}
```

### 请求处理

#### 获取 GET 参数

```go
func (this *User) Search() {
    // 获取单个参数
    keyword := this.Ctx.GET("q")
    
    // 获取参数，带默认值
    page := this.Ctx.GET("page", "1")
    
    // 获取数组参数 (如: ?tags=go&tags=web)
    tags := this.Ctx.GETArray("tags")
    
    // 获取所有 GET 参数
    params := this.Ctx.GETData()
    
    // 使用参数...
}
```

#### 获取 POST 数据

```go
func (this *User) Update() {
    // 获取单个 POST 参数
    name := this.Ctx.POST("name")
    email := this.Ctx.POST("email")
    
    // 获取参数，带默认值
    gender := this.Ctx.POST("gender", "未知")
    
    // 获取数组参数
    interests := this.Ctx.POSTArray("interests")
    
    // 获取所有 POST 数据
    allData := this.Ctx.POSTData()
    
    // 优先从 POST 获取，没有则从 GET 获取
    value := this.Ctx.GetPost("key", "default")
}
```

#### 获取 JSON 请求体

```go
func (this *User) CreateUser() {
    // 读取原始请求体
    body, err := this.Ctx.RequestBody()
    if err != nil {
        this.Ctx.JSON(400, map[string]string{"error": "无法读取请求体"})
        return
    }
    
    // 解析 JSON
    var user struct {
        Name  string `json:"name"`
        Email string `json:"email"`
        Age   int    `json:"age"`
    }
    
    if err := json.Unmarshal(body, &user); err != nil {
        this.Ctx.JSON(400, map[string]string{"error": "JSON 格式错误"})
        return
    }
    
    // 使用数据...
    this.Ctx.JSON(200, map[string]interface{}{
        "message": "用户创建成功",
        "user": user,
    })
}
```

#### 判断请求方法

```go
func (this *User) HandleRequest() {
    if this.Ctx.IsGET() {
        // 处理 GET 请求
        this.ShowForm()
        return
    }
    
    if this.Ctx.IsPOST() {
        // 处理 POST 请求
        this.ProcessForm()
        return
    }
    
    if this.Ctx.IsPUT() {
        // 处理 PUT 请求
        this.UpdateResource()
        return
    }
    
    if this.Ctx.IsDELETE() {
        // 处理 DELETE 请求
        this.DeleteResource()
        return
    }
    
    // 其他判断
    if this.Ctx.IsAJAX() {
        // AJAX 请求
    }
    
    if this.Ctx.IsWebsocket() {
        // WebSocket 请求
    }
}
```

### 响应输出

#### 文本输出

```go
func (this *User) Hello() {
    // 基本输出
    this.Write("Hello World")
    
    // 格式化输出
    this.Write(fmt.Sprintf("Hello %s", name))
    
    // 多参数输出
    this.Write("Hello", " ", "World")
    
    // 输出带错误处理
    n, err := this.WriteE("Hello")
    if err != nil {
        this.Logger.Error(err)
    }
}
```

#### JSON 响应

```go
func (this *User) GetUser() {
    user := map[string]interface{}{
        "id":    1,
        "name":  "张三",
        "email": "zhangsan@example.com",
    }
    
    // 标准 JSON 响应
    this.Ctx.JSON(200, user)
    
    // 格式化的 JSON（开发环境友好）
    this.Ctx.PrettyJSON(200, user)
    
    // JSONP 响应（跨域）
    this.Ctx.JSONP(200, user)
}

// RESTful API 响应封装
func (this *User) ApiResponse(code int, message string, data interface{}) {
    response := map[string]interface{}{
        "code":    code,
        "message": message,
        "data":    data,
    }
    this.Ctx.JSON(code, response)
}

// 使用
func (this *User) List() {
    users := this.GetAllUsers()
    this.ApiResponse(200, "success", users)
}
```

#### 重定向

```go
func (this *User) Login() {
    // 临时重定向（302）
    this.Ctx.Redirect("/dashboard")
    
    // 永久重定向，需要手动设置
    this.Ctx.WriteHeader(301)
    this.Ctx.SetHeader("Location", "/new-url")
}
```

#### 文件下载

```go
func (this *User) Download() {
    // 直接输出文件
    this.Ctx.WriteFile("/path/to/file.pdf")
    
    // 指定下载文件名
    this.Ctx.WriteFileAttachment("/path/to/file.pdf", "report.pdf")
    
    // 从自定义文件系统读取
    fs := http.Dir("uploads")
    this.Ctx.WriteFileFromFS("document.pdf", fs)
}
```

#### 设置响应头

```go
func (this *User) CustomResponse() {
    // 设置单个响应头
    this.Ctx.SetHeader("Content-Type", "application/json")
    this.Ctx.SetHeader("X-Custom-Header", "value")
    
    // 设置状态码
    this.Ctx.WriteHeader(404)
    
    // 设置 Cookie
    this.Ctx.SetCookie("session_id", "abc123", 3600, "/", "", false, true)
    
    // 输出内容
    this.Write("Custom Response")
}
```

### 访问控制器成员

控制器提供了丰富的成员变量，方便访问各种功能：

```go
func (this *User) Example() {
    // HTTP 请求和响应
    req := this.Request    // *http.Request
    res := this.Response   // http.ResponseWriter
    
    // 路由参数
    params := this.Param   // gcore.Params
    
    // 上下文对象
    ctx := this.Ctx        // gcore.Ctx
    
    // 配置对象
    cfg := this.Config     // gcore.Config
    
    // 日志对象
    log := this.Logger     // gcore.Logger
    
    // 模板对象
    tpl := this.Tpl        // gcore.Template
    
    // Session 对象（需要先 SessionStart）
    sess := this.Session   // gcore.Session
    
    // Cookie 对象
    cookie := this.Cookie  // gcore.Cookies
    
    // 视图对象
    view := this.View      // gcore.View
    
    // 路由对象
    router := this.Router  // gcore.HTTPRouter
    
    // 国际化对象
    i18n := this.Lang      // gcore.I18n
}
```

---

## 请求与响应

### 获取输入

#### 表单数据

```go
func (this *User) HandleForm() {
    // GET 参数
    search := this.Ctx.GET("q")
    page := this.Ctx.GET("page", "1")
    
    // POST 数据
    username := this.Ctx.POST("username")
    password := this.Ctx.POST("password")
    
    // 优先 POST，其次 GET
    token := this.Ctx.GetPost("token")
    
    // 获取所有表单数据
    postData := this.Ctx.POSTData()    // map[string]string
    getData := this.Ctx.GETData()      // map[string]string
}
```

#### 获取请求头

```go
func (this *User) Headers() {
    // 获取单个请求头
    contentType := this.Ctx.Header("Content-Type")
    userAgent := this.Ctx.Header("User-Agent")
    auth := this.Ctx.Header("Authorization")
    
    // 获取所有请求头
    headers := this.Request.Header
    for key, values := range headers {
        fmt.Printf("%s: %v\n", key, values)
    }
}
```

#### Cookie 操作

> **详细文档：** [http/cookie/README.md](https://github.com/snail007/gmc/blob/master/http/cookie/README.md) - 查看完整的 Cookie API、安全选项和最佳实践

```go
func (this *User) CookieDemo() {
    // 读取 Cookie
    sessionID := this.Ctx.Cookie("session_id")
    
    // 设置 Cookie
    // SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
    this.Ctx.SetCookie("user_id", "123", 3600, "/", "", false, true)
    
    // 使用 Cookie 对象
    this.Cookie.Set("token", "abc123", &gcore.CookieOptions{
        MaxAge:   7200,
        Path:     "/",
        Domain:   "",
        Secure:   false,
        HTTPOnly: true,
    })
    
    // 读取 Cookie
    token, err := this.Cookie.Get("token")
    if err != nil {
        this.Logger.Error(err)
    }
    
    // 删除 Cookie
    this.Cookie.Remove("token")
}
```

### 文件上传

#### 单文件上传

```go
func (this *User) Upload() {
    // 获取上传的文件
    // FormFile(fieldName, maxMemory)
    file, err := this.Ctx.FormFile("avatar", 10<<20) // 10MB
    if err != nil {
        this.Ctx.JSON(400, map[string]string{
            "error": "文件上传失败: " + err.Error(),
        })
        return
    }
    
    // 生成保存路径
    filename := fmt.Sprintf("uploads/%s", file.Filename)
    
    // 保存文件
    if err := this.Ctx.SaveUploadedFile(file, filename); err != nil {
        this.Ctx.JSON(500, map[string]string{
            "error": "文件保存失败: " + err.Error(),
        })
        return
    }
    
    this.Ctx.JSON(200, map[string]interface{}{
        "message":  "上传成功",
        "filename": filename,
        "size":     file.Size,
    })
}
```

#### 多文件上传

```go
func (this *User) MultiUpload() {
    // 获取 multipart 表单
    form, err := this.Ctx.MultipartForm(32 << 20) // 32MB
    if err != nil {
        this.Ctx.JSON(400, map[string]string{
            "error": "解析表单失败: " + err.Error(),
        })
        return
    }
    
    // 获取多个文件
    files := form.File["files"]
    
    var savedFiles []string
    for _, file := range files {
        filename := fmt.Sprintf("uploads/%s", file.Filename)
        
        if err := this.Ctx.SaveUploadedFile(file, filename); err != nil {
            this.Logger.Errorf("保存文件失败: %v", err)
            continue
        }
        
        savedFiles = append(savedFiles, filename)
    }
    
    this.Ctx.JSON(200, map[string]interface{}{
        "message": "上传成功",
        "files":   savedFiles,
        "count":   len(savedFiles),
    })
}
```

#### 文件验证

```go
func (this *User) ValidateUpload() {
    file, err := this.Ctx.FormFile("file", 5<<20) // 5MB
    if err != nil {
        this.Ctx.JSON(400, map[string]string{"error": "文件上传失败"})
        return
    }
    
    // 验证文件大小
    if file.Size > 5*1024*1024 {
        this.Ctx.JSON(400, map[string]string{"error": "文件不能超过5MB"})
        return
    }
    
    // 验证文件类型
    allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif"}
    ext := strings.ToLower(filepath.Ext(file.Filename))
    
    allowed := false
    for _, t := range allowedTypes {
        if ext == t {
            allowed = true
            break
        }
    }
    
    if !allowed {
        this.Ctx.JSON(400, map[string]string{
            "error": "只允许上传图片文件",
        })
        return
    }
    
    // 保存文件...
}
```

### JSON 响应

#### 标准 JSON 响应

```go
func (this *User) JsonResponse() {
    data := map[string]interface{}{
        "code":    200,
        "message": "success",
        "data": map[string]interface{}{
            "user_id": 123,
            "name":    "张三",
        },
    }
    
    // 紧凑 JSON
    this.Ctx.JSON(200, data)
    
    // 格式化 JSON（便于调试）
    this.Ctx.PrettyJSON(200, data)
}
```

#### RESTful API 响应封装

```go
// 定义响应结构
type ApiResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func (this *User) Success(data interface{}, message ...string) {
    msg := "success"
    if len(message) > 0 {
        msg = message[0]
    }
    
    this.Ctx.JSON(200, ApiResponse{
        Code:    200,
        Message: msg,
        Data:    data,
    })
}

// 错误响应
func (this *User) Error(code int, message string) {
    this.Ctx.JSON(code, ApiResponse{
        Code:    code,
        Message: message,
    })
}

// 分页响应
func (this *User) PageResponse(items interface{}, total int64, page, pageSize int) {
    this.Ctx.JSON(200, map[string]interface{}{
        "code":    200,
        "message": "success",
        "data": map[string]interface{}{
            "items":     items,
            "total":     total,
            "page":      page,
            "page_size": pageSize,
            "total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
        },
    })
}

// 使用示例
func (this *User) List() {
    users, err := this.GetUsers()
    if err != nil {
        this.Error(500, "获取用户列表失败")
        return
    }
    
    this.Success(users)
}
```

### 重定向

```go
func (this *User) RedirectDemo() {
    // 临时重定向（302）
    this.Ctx.Redirect("/dashboard")
    
    // 根据条件重定向
    if !this.IsLoggedIn() {
        this.Ctx.Redirect("/login")
        this.Stop()
        return
    }
    
    // 重定向到外部URL
    this.Ctx.Redirect("https://www.example.com")
    
    // 永久重定向（301）
    this.Ctx.WriteHeader(301)
    this.Ctx.SetHeader("Location", "/new-url")
}
```

### 响应类型

#### HTML 响应

```go
func (this *User) HtmlResponse() {
    html := `
    <!DOCTYPE html>
    <html>
    <head><title>Test</title></head>
    <body><h1>Hello GMC</h1></body>
    </html>
    `
    this.Ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
    this.Write(html)
}
```

#### XML 响应

```go
func (this *User) XmlResponse() {
    xml := `<?xml version="1.0" encoding="UTF-8"?>
    <response>
        <code>200</code>
        <message>success</message>
    </response>`
    
    this.Ctx.SetHeader("Content-Type", "application/xml; charset=utf-8")
    this.Write(xml)
}
```

#### 二进制响应

```go
func (this *User) BinaryResponse() {
    data := []byte{0x89, 0x50, 0x4E, 0x47...} // PNG header
    
    this.Ctx.SetHeader("Content-Type", "image/png")
    this.Ctx.SetHeader("Content-Length", fmt.Sprintf("%d", len(data)))
    this.Ctx.Response().Write(data)
}
```

### 客户端信息

```go
func (this *User) ClientInfo() {
    // 获取客户端 IP
    ip := this.Ctx.ClientIP()
    
    // 获取完整的远程地址
    remoteAddr := this.Ctx.RemoteAddr()
    
    // 获取 Host
    host := this.Ctx.Host()
    
    // 判断请求类型
    isAjax := this.Ctx.IsAJAX()
    isWS := this.Ctx.IsWebsocket()
    isTLS := this.Ctx.IsTLSRequest()
    
    this.Logger.Infof("Client IP: %s, IsAJAX: %v", ip, isAjax)
}
```

---

*由于文档内容非常长，这只是第一部分。文档将继续包含视图模板、数据库、缓存、Session、日志、国际化、API 开发、测试、部署等完整章节...*


## 视图模板

GMC 的模板引擎基于 Go 标准库的 `text/template`，增强了模板继承、包含、自定义函数等功能。

> **详细文档：** 
> - [http/template/README.md](https://github.com/snail007/gmc/blob/master/http/template/README.md) - 模板引擎详细说明
> - [http/view/README.md](https://github.com/snail007/gmc/blob/master/http/view/README.md) - 视图渲染完整文档

### 模板配置

GMC 的模板引擎基于 Go 标准库的 `text/template`，并进行了功能增强。

#### 配置文件

在 `conf/app.toml` 中配置模板：

```toml
[template]
dir = "views"              # 模板文件目录
ext = ".html"              # 模板文件扩展名
delimiterleft = "{{"       # 左分隔符
delimiterright = "}}"      # 右分隔符
```

#### 目录结构

推荐的模板目录结构：

```text
views/
├── layout/
│   ├── base.html          # 基础布局
│   └── admin.html         # 管理后台布局
├── common/
│   ├── header.html        # 公共头部
│   ├── footer.html        # 公共底部
│   └── sidebar.html       # 侧边栏
├── user/
│   ├── list.html          # 用户列表
│   ├── detail.html        # 用户详情
│   └── edit.html          # 编辑用户
└── home.html              # 首页
```

### 模板语法

#### 变量输出

```html
<!-- 输出变量 -->
<h1>{{.Title}}</h1>
<p>{{.Content}}</p>

<!-- 访问结构体字段 -->
<div>
    <p>用户名: {{.User.Name}}</p>
    <p>邮箱: {{.User.Email}}</p>
</div>

<!-- 访问 Map -->
<p>{{.Data.Key1}}</p>

<!-- HTML 转义输出（默认） -->
<p>{{.Content}}</p>

<!-- 不转义输出 -->
<p>{{.Content | html}}</p>
```

#### 控制结构

```html
<!-- if 条件判断 -->
{{if .IsLoggedIn}}
    <p>欢迎回来，{{.Username}}！</p>
{{else}}
    <p>请先登录</p>
{{end}}

<!-- if-else if-else -->
{{if eq .Role "admin"}}
    <p>管理员</p>
{{else if eq .Role "user"}}
    <p>普通用户</p>
{{else}}
    <p>访客</p>
{{end}}

<!-- range 循环 -->
<ul>
{{range .Users}}
    <li>{{.Name}} - {{.Email}}</li>
{{end}}
</ul>

<!-- range 带索引 -->
{{range $index, $user := .Users}}
    <p>{{$index}}: {{$user.Name}}</p>
{{end}}

<!-- range Map -->
{{range $key, $value := .Data}}
    <p>{{$key}}: {{$value}}</p>
{{end}}

<!-- with 设置作用域 -->
{{with .User}}
    <p>姓名: {{.Name}}</p>
    <p>年龄: {{.Age}}</p>
{{end}}
```

#### 比较运算

```html
<!-- 相等 -->
{{if eq .Status "active"}}激活{{end}}

<!-- 不等 -->
{{if ne .Count 0}}有数据{{end}}

<!-- 小于 -->
{{if lt .Age 18}}未成年{{end}}

<!-- 小于等于 -->
{{if le .Score 60}}不及格{{end}}

<!-- 大于 -->
{{if gt .Price 100}}贵{{end}}

<!-- 大于等于 -->
{{if ge .Level 5}}高级{{end}}
```

#### 逻辑运算

```html
<!-- and -->
{{if and .IsLoggedIn .IsAdmin}}
    管理员已登录
{{end}}

<!-- or -->
{{if or .IsAdmin .IsModerator}}
    有权限
{{end}}

<!-- not -->
{{if not .IsDeleted}}
    <p>显示内容</p>
{{end}}
```

### 模板继承

GMC 使用 Layout 布局系统，通过 `{{.GMC_LAYOUT_CONTENT}}` 占位符实现模板继承。

#### 定义布局文件

`views/layout/page.html`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.title}}</title>
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body>
    <header>
        <h1>我的网站</h1>
        <nav>
            <a href="/">首页</a>
            <a href="/user">用户</a>
            <a href="/about">关于</a>
        </nav>
    </header>
    
    <main>
        {{.GMC_LAYOUT_CONTENT}}
    </main>
    
    <footer>
        <p>&copy; 2024 我的网站</p>
    </footer>
    
    <script src="/static/js/main.js"></script>
</body>
</html>
```

**说明：**
- `{{.GMC_LAYOUT_CONTENT}}` 是 GMC 的特殊占位符，会被实际内容模板替换
- 布局文件通常放在 `views/layout/` 目录下

#### 创建内容模板

`views/user/list.html`:

```html
<div class="user-list">
    <h1>用户列表</h1>
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>姓名</th>
                <th>邮箱</th>
                <th>操作</th>
            </tr>
        </thead>
        <tbody>
        {{range .users}}
            <tr>
                <td>{{.ID}}</td>
                <td>{{.Name}}</td>
                <td>{{.Email}}</td>
                <td>
                    <a href="/user/{{.ID}}">查看</a>
                    <a href="/user/{{.ID}}/edit">编辑</a>
                </td>
            </tr>
        {{end}}
        </tbody>
    </table>
</div>
```

`views/welcome.html`:

```html
<div class="welcome">
    <h2>{{.title}}</h2>
    <p>{{.message}}</p>
</div>
```

#### 在控制器中使用布局

**方式一：指定布局名称（不带扩展名）**

```go
func (this *User) List() {
    users := this.GetAllUsers()
    
    // 使用 Layout 方法指定布局，然后渲染内容模板
    this.View.Layout("page").Render("user/list", map[string]interface{}{
        "title": "用户列表",
        "users": users,
    })
}
```

**方式二：指定布局完整路径**

```go
func (this *Demo) Welcome() {
    // 可以带扩展名
    this.View.Layout("page.html").Render("welcome.html", map[string]interface{}{
        "title":   "欢迎",
        "message": "欢迎使用 GMC 框架",
    })
}
```

**方式三：使用相对路径**

```go
func (this *Demo) Index() {
    // 使用 layout 子目录
    this.View.Layout("layout/page").Render("welcome", map[string]interface{}{
        "title": "首页",
    })
}
```

### 模板包含

GMC 支持模板包含（Include），可以将公共的模板片段复用。框架会自动处理模板文件的加载，不需要使用 `{{define}}`。

#### 创建公共组件

`views/common/header.html`:

```html
<header>
    <nav>
        <ul>
            <li><a href="/">首页</a></li>
            <li><a href="/about">关于</a></li>
            <li><a href="/contact">联系</a></li>
        </ul>
    </nav>
</header>
```

`views/common/footer.html`:

```html
<footer>
    <p>&copy; 2024 我的网站. 保留所有权利.</p>
</footer>
```

#### 包含组件

在其他模板中使用 `{{template}}` 包含组件：

`views/user/list.html`:

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>用户列表</title>
</head>
<body>
    <!-- 包含 header，传递当前数据 -->
    {{template "common/header.html" .}}
    
    <main>
        <h1>用户列表</h1>
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>姓名</th>
                    <th>邮箱</th>
                </tr>
            </thead>
            <tbody>
            {{range .users}}
                <tr>
                    <td>{{.ID}}</td>
                    <td>{{.Name}}</td>
                    <td>{{.Email}}</td>
                </tr>
            {{end}}
            </tbody>
        </table>
    </main>
    
    <!-- 包含 footer -->
    {{template "common/footer.html" .}}
</body>
</html>
```

#### 包含分页组件

`views/paginator/default.html`:

```html
{{ $pager := .paginator }}
<ul class="pagination">
    {{if $pager.HasPrev}}
        <li><a href="{{$pager.PageLinkFirst}}">首页</a></li>
        <li><a href="{{$pager.PageLinkPrev}}">上一页</a></li>
    {{end}}
    
    {{range $index, $page := $pager.Pages}}
        <li{{if $pager.IsActive .}} class="active"{{end}}>
            <a href="{{$pager.PageLink $page}}">{{$page}}</a>
        </li>
    {{end}}
    
    {{if $pager.HasNext}}
        <li><a href="{{$pager.PageLinkNext}}">下一页</a></li>
        <li><a href="{{$pager.PageLinkLast}}">尾页</a></li>
    {{end}}
</ul>
```

在列表页面中使用：

```html
<div class="user-list">
    <!-- 用户列表内容 -->
</div>

<!-- 包含分页组件 -->
{{template "paginator/default.html" .}}
```

#### 说明

- **文件路径**: 使用相对于 views 目录的路径，包含扩展名
- **数据传递**: 使用 `.` 传递当前所有数据，或传递特定变量
- **自动加载**: GMC 会自动加载和解析所有模板文件，无需 `{{define}}`
- **命名约定**: 建议使用清晰的目录结构组织公共组件

### 自定义函数

#### 注册模板函数

在 `initialize/initialize.go` 中注册：

```go
func Initialize(s *gmc.HTTPServer) error {
    // 定义自定义函数
    funcMap := map[string]interface{}{
        // 字符串长度
        "strlen": func(str string) int {
            return len(str)
        },
        
        // 格式化日期
        "formatDate": func(t time.Time) string {
            return t.Format("2006-01-02 15:04:05")
        },
        
        // 字符串截取
        "substr": func(s string, start, length int) string {
            if start < 0 || start >= len(s) {
                return ""
            }
            end := start + length
            if end > len(s) {
                end = len(s)
            }
            return s[start:end]
        },
        
        // HTML 转义
        "escape": func(s string) string {
            return html.EscapeString(s)
        },
        
        // 转大写
        "upper": func(s string) string {
            return strings.ToUpper(s)
        },
        
        // 数字格式化
        "formatNumber": func(n float64) string {
            return fmt.Sprintf("%.2f", n)
        },
    }
    
    // 注册函数到模板引擎
    s.AddFuncMap(funcMap)
    
    return nil
}
```

#### 在模板中使用自定义函数

`views/func.html`:

```html
<div>
    <p>名称: {{.name}}</p>
    <p>长度: {{strlen .name}}</p>
    <p>大写: {{upper .name}}</p>
</div>

<div>
    <p>创建时间: {{formatDate .created_at}}</p>
    <p>价格: ￥{{formatNumber .price}}</p>
</div>
```

**控制器代码：**

```go
func (this *Demo) Func() {
    this.View.Set("name", "hello")
    this.View.Render("func")
}
        
        // 截断字符串
        "truncate": func(s string, length int) string {
            if len(s) <= length {
                return s
            }
            return s[:length] + "..."
        },
        
        // 数字格式化
        "formatNumber": func(n int) string {
            return fmt.Sprintf("%,d", n)
        },
        
        // URL 生成
        "url": func(path string) string {
            return "/app" + path
        },
        
        // 判断是否在切片中
        "in": func(item interface{}, slice []interface{}) bool {
            for _, v := range slice {
                if v == item {
                    return true
                }
            }
            return false
        },
    })
}
```

#### 在模板中使用

```html
<!-- 格式化日期 -->
<p>创建时间: {{formatDate .CreatedAt}}</p>

<!-- 截断字符串 -->
<p>{{truncate .Description 100}}</p>

<!-- 格式化数字 -->
<p>价格: ¥{{formatNumber .Price}}</p>

<!-- URL 生成 -->
<a href="{{url "/user/profile"}}">个人资料</a>

<!-- 条件判断 -->
{{if in .CurrentPage .ActivePages}}
    <li class="active">{{.CurrentPage}}</li>
{{end}}
```

### 视图对象使用

#### 基本使用

```go
func (this *User) Profile() {
    user := this.GetUser()
    
    // 方式 1: 使用 View 对象
    this.View.
        Set("user", user).
        Set("title", "用户资料").
        Render("user/profile")
    
    // 方式 2: 使用 Map
    this.View.Render("user/profile", map[string]interface{}{
        "user":  user,
        "title": "用户资料",
    })
}
```

#### 设置布局

```go
func (this *User) Index() {
    // 使用默认布局
    this.View.Layout("layout/base").Render("user/index")
    
    // 不使用布局
    this.View.Layout("").Render("user/ajax")
}
```

#### 链式调用

```go
func (this *User) Dashboard() {
    this.View.
        Layout("layout/admin").
        Set("title", "控制台").
        Set("stats", this.GetStats()).
        Set("charts", this.GetCharts()).
        Render("admin/dashboard")
}
```

### 模板错误处理

```go
func (this *User) SafeRender() {
    this.View.Render("user/profile")
    
    // 检查渲染错误
    if err := this.View.Err(); err != nil {
        this.Logger.Errorf("模板渲染错误: %v", err)
        this.Ctx.WriteHeader(500)
        this.Write("页面渲染失败")
        return
    }
}
```

### 模板最佳实践

1. **目录组织**: 按功能模块组织模板文件
2. **公共组件**: 提取可复用的组件
3. **布局继承**: 使用布局减少重复代码
4. **数据准备**: 在控制器中准备好所有数据
5. **错误处理**: 优雅处理模板错误
6. **性能优化**: 避免在模板中进行复杂计算

---

## 数据库

GMC 提供强大的数据库操作能力，基于 GORM 封装，支持 MySQL、PostgreSQL、SQLite 等多种数据库。

> **详细文档：** [module/db/README.md](https://github.com/snail007/gmc/blob/master/module/db/README.md) - 查看完整的数据库 API、查询构建器、事务处理等

### 配置连接

#### 数据库配置

在 `conf/app.toml` 中配置数据库：

```toml
[database]
enable = true
driver = "mysql"
dsn = "root:password@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
maxidle = 10
maxconns = 100
maxlifetimeseconds = 3600
timeout = 5000
# debug = true  # 开启 SQL 调试

# 多数据库配置
[[database.groups]]
name = "default"
driver = "mysql"
dsn = "root:password@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True"

[[database.groups]]
name = "analytics"
driver = "mysql"
dsn = "root:password@tcp(127.0.0.1:3306)/analytics?charset=utf8mb4"
```

#### 支持的数据库

- MySQL / MariaDB
- PostgreSQL
- SQLite
- SQL Server

#### 初始化数据库

```go
func InitDatabase(cfg gcore.Config) error {
    // 获取数据库配置
    dbCfg := cfg.Sub("database")
    
    if !dbCfg.GetBool("enable") {
        return nil
    }
    
    // 数据库会在第一次使用时自动连接
    return nil
}
```

### 查询构建器

GMC 提供了强大的 ActiveRecord 风格查询构建器。

#### 获取数据库对象

```go
func (this *User) GetDB() gcore.Database {
    // 从上下文获取数据库
    db, err := gcore.ProviderDatabase()(this.Ctx)
    if err != nil {
        this.Logger.Errorf("获取数据库失败: %v", err)
        return nil
    }
    return db
}
```

#### SELECT 查询

```go
func (this *User) QueryExamples() {
    db := this.GetDB()
    
    // 基本查询
    ar := db.AR()
    ar.Select("*").From("users")
    result, err := db.Query(ar)
    
    // 指定字段
    ar = db.AR()
    ar.Select("id, name, email").From("users")
    result, _ = db.Query(ar)
    
    // WHERE 条件
    ar = db.AR()
    ar.Select("*").
        From("users").
        Where(map[string]interface{}{
            "status": "active",
            "age >":  18,
        })
    result, _ = db.Query(ar)
    
    // 复杂条件
    ar = db.AR()
    ar.Select("*").
        From("users").
        Where(map[string]interface{}{
            "status":     "active",
            "age >=":     18,
            "age <=":     60,
            "city IN":    []string{"北京", "上海", "深圳"},
            "name LIKE":  "%张%",
        })
    result, _ = db.Query(ar)
    
    // ORDER BY
    ar = db.AR()
    ar.Select("*").
        From("users").
        OrderBy("created_at", "DESC")
    result, _ = db.Query(ar)
    
    // LIMIT
    ar = db.AR()
    ar.Select("*").
        From("users").
        Limit(10, 0) // LIMIT 10 OFFSET 0
    result, _ = db.Query(ar)
    
    // JOIN
    ar = db.AR()
    ar.Select("u.*, p.title").
        From("users u").
        Join("posts p", "p", "u.id=p.user_id", "LEFT")
    result, _ = db.Query(ar)
    
    // GROUP BY
    ar = db.AR()
    ar.Select("city, COUNT(*) as count").
        From("users").
        GroupBy("city").
        Having("COUNT(*) > 10")
    result, _ = db.Query(ar)
}
```

#### INSERT 操作

```go
func (this *User) InsertExamples() {
    db := this.GetDB()
    
    // 单条插入
    ar := db.AR()
    ar.Insert("users", map[string]interface{}{
        "name":       "张三",
        "email":      "zhangsan@example.com",
        "created_at": time.Now(),
    })
    result, err := db.Exec(ar)
    if err != nil {
        this.Logger.Error(err)
        return
    }
    
    // 获取插入 ID
    lastID := result.LastInsertID()
    this.Logger.Infof("插入成功，ID: %d", lastID)
    
    // 批量插入
    ar = db.AR()
    users := []map[string]interface{}{
        {"name": "李四", "email": "lisi@example.com"},
        {"name": "王五", "email": "wangwu@example.com"},
    }
    ar.InsertBatch("users", users)
    result, _ = db.Exec(ar)
}
```

#### UPDATE 操作

```go
func (this *User) UpdateExamples() {
    db := this.GetDB()
    
    // 更新数据
    ar := db.AR()
    ar.Update("users",
        map[string]interface{}{
            "name":       "张三三",
            "updated_at": time.Now(),
        },
        map[string]interface{}{
            "id": 1,
        },
    )
    result, err := db.Exec(ar)
    
    // 获取受影响的行数
    affected := result.RowsAffected()
    this.Logger.Infof("更新了 %d 行", affected)
    
    // 批量更新
    ar = db.AR()
    updates := []map[string]interface{}{
        {"id": 1, "status": "active"},
        {"id": 2, "status": "inactive"},
    }
    ar.UpdateBatch("users", updates, []string{"id"})
    db.Exec(ar)
}
```

#### DELETE 操作

```go
func (this *User) DeleteExamples() {
    db := this.GetDB()
    
    // 删除数据
    ar := db.AR()
    ar.Delete("users", map[string]interface{}{
        "id": 1,
    })
    result, err := db.Exec(ar)
    
    // 条件删除
    ar = db.AR()
    ar.Delete("users", map[string]interface{}{
        "status":         "inactive",
        "created_at <":   time.Now().AddDate(0, -6, 0),
    })
    db.Exec(ar)
}
```

#### 原始 SQL

```go
func (this *User) RawSQL() {
    db := this.GetDB()
    
    // 查询
    sql := "SELECT * FROM users WHERE age > ? AND city = ?"
    result, err := db.QuerySQL(sql, 18, "北京")
    
    // 执行
    sql = "UPDATE users SET status = ? WHERE id = ?"
    result, err = db.ExecSQL(sql, "active", 123)
}
```

### 事务处理

#### 基本事务

```go
func (this *User) TransactionExample() {
    db := this.GetDB()
    
    // 开始事务
    tx, err := db.Begin()
    if err != nil {
        this.Logger.Error(err)
        return
    }
    
    // 延迟回滚或提交
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            this.Logger.Errorf("事务回滚: %v", r)
        }
    }()
    
    // 执行操作
    ar := db.AR()
    ar.Insert("users", map[string]interface{}{
        "name": "测试用户",
    })
    result, err := db.ExecTx(ar, tx)
    if err != nil {
        tx.Rollback()
        return
    }
    
    userID := result.LastInsertID()
    
    // 第二个操作
    ar = db.AR()
    ar.Insert("profiles", map[string]interface{}{
        "user_id": userID,
        "bio":     "个人简介",
    })
    _, err = db.ExecTx(ar, tx)
    if err != nil {
        tx.Rollback()
        return
    }
    
    // 提交事务
    if err := tx.Commit(); err != nil {
        this.Logger.Error(err)
        return
    }
    
    this.Logger.Info("事务提交成功")
}
```

#### 事务封装

```go
func (this *User) WithTransaction(fn func(tx *sql.Tx) error) error {
    db := this.GetDB()
    
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}

// 使用
func (this *User) CreateUserWithProfile() {
    err := this.WithTransaction(func(tx *sql.Tx) error {
        // 创建用户
        ar := db.AR()
        ar.Insert("users", map[string]interface{}{
            "name": "张三",
        })
        result, err := db.ExecTx(ar, tx)
        if err != nil {
            return err
        }
        
        // 创建资料
        userID := result.LastInsertID()
        ar = db.AR()
        ar.Insert("profiles", map[string]interface{}{
            "user_id": userID,
            "bio":     "简介",
        })
        _, err = db.ExecTx(ar, tx)
        return err
    })
    
    if err != nil {
        this.Logger.Error(err)
        this.Ctx.JSON(500, map[string]string{"error": "操作失败"})
        return
    }
    
    this.Ctx.JSON(200, map[string]string{"message": "成功"})
}
```

### 结果集处理

#### 获取数据行

```go
func (this *User) ResultExamples() {
    db := this.GetDB()
    
    ar := db.AR()
    ar.Select("*").From("users")
    result, err := db.Query(ar)
    if err != nil {
        return
    }
    
    // 获取所有行（map 切片）
    rows := result.Rows()
    for _, row := range rows {
        fmt.Printf("ID: %s, Name: %s\n", row["id"], row["name"])
    }
    
    // 获取单行
    row := result.Row()
    if row != nil {
        fmt.Printf("用户: %s\n", row["name"])
    }
    
    // 获取指定列的所有值
    names := result.Values("name")
    for _, name := range names {
        fmt.Println(name)
    }
    
    // 获取键值对 map
    userMap := result.MapValues("id", "name")
    // userMap: {"1": "张三", "2": "李四"}
}
```

#### 映射到结构体

```go
type User struct {
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
}

func (this *UserController) GetUsers() {
    db := this.GetDB()
    
    ar := db.AR()
    ar.Select("*").From("users")
    result, _ := db.Query(ar)
    
    // 映射到结构体切片
    users, err := result.Structs(&User{})
    if err != nil {
        this.Logger.Error(err)
        return
    }
    
    for _, u := range users {
        user := u.(*User)
        fmt.Printf("用户: %s (%s)\n", user.Name, user.Email)
    }
    
    // 映射单个结构体
    user, err := result.Struct(&User{})
    if err != nil {
        this.Logger.Error(err)
        return
    }
    
    u := user.(*User)
    fmt.Printf("用户: %s\n", u.Name)
}
```

#### 使用键构建 Map

```go
func (this *User) MapResults() {
    db := this.GetDB()
    
    ar := db.AR()
    ar.Select("*").From("users")
    result, _ := db.Query(ar)
    
    // 以 id 为键的 map[string]map[string]string
    usersMap := result.MapRows("id")
    user1 := usersMap["1"]
    fmt.Printf("用户1: %s\n", user1["name"])
    
    // 以 id 为键映射到结构体
    structsMap, _ := result.MapStructs("id", &User{})
    user := structsMap["1"].(*User)
    fmt.Printf("用户: %s\n", user.Name)
}
```

### 查询缓存

```go
func (this *User) CachedQuery() {
    db := this.GetDB()
    
    // 使用缓存（60秒）
    ar := db.AR()
    ar.Cache("users:list", 60).
        Select("*").
        From("users").
        Where(map[string]interface{}{
            "status": "active",
        })
    
    result, err := db.Query(ar)
    // 第一次查询会访问数据库并缓存结果
    // 60秒内的相同查询会直接从缓存返回
}
```

### 数据库连接池

```go
func (this *User) PoolStats() {
    db := this.GetDB()
    
    // 获取连接池统计
    stats := db.Stats()
    
    this.Logger.Infof("连接池状态:")
    this.Logger.Infof("  打开连接数: %d", stats.OpenConnections)
    this.Logger.Infof("  使用中: %d", stats.InUse)
    this.Logger.Infof("  空闲: %d", stats.Idle)
    this.Logger.Infof("  等待: %d", stats.WaitCount)
}
```

### 多数据库

```go
func InitMultiDatabase(cfg gcore.Config) error {
    // 获取数据库组管理器
    dbGroup, err := gcore.ProviderDatabaseGroup()(nil)
    if err != nil {
        return err
    }
    
    // 注册数据库组
    err = dbGroup.RegistGroup(cfg.Sub("database"))
    if err != nil {
        return err
    }
    
    return nil
}

// 使用不同的数据库
func (this *User) MultiDBExample() {
    dbGroup, _ := gcore.ProviderDatabaseGroup()(this.Ctx)
    
    // 使用默认数据库
    defaultDB := dbGroup.DB()
    
    // 使用分析数据库
    analyticsDB := dbGroup.DB("analytics")
    
    // 执行查询
    ar := defaultDB.AR()
    ar.Select("*").From("users")
    defaultDB.Query(ar)
    
    ar = analyticsDB.AR()
    ar.Select("*").From("events")
    analyticsDB.Query(ar)
}
```

### 数据库最佳实践

1. **连接池配置**: 根据负载合理设置连接池大小
2. **索引优化**: 为常用查询字段添加索引
3. **事务使用**: 需要原子性操作时使用事务
4. **预防 SQL 注入**: 使用参数化查询，避免拼接 SQL
5. **查询优化**: 只查询需要的字段，避免 SELECT *
6. **批量操作**: 大量数据使用批量插入/更新
7. **缓存使用**: 对热点数据使用缓存
8. **错误处理**: 妥善处理数据库错误

---

## 缓存

GMC Cache 模块提供统一的缓存接口，支持 Redis、内存缓存、文件缓存等多种后端，支持多数据源配置和管理。

> **详细文档：** [module/cache/README.md](https://github.com/snail007/gmc/blob/master/module/cache/README.md) - 查看完整的配置选项、高级功能和使用示例

### 缓存配置

#### 配置文件

在 `conf/app.toml` 中配置缓存：

```toml
[cache]
enable = true

# 内存缓存
[[cache.stores]]
store = "memory"
cleanupintervalseconds = 60

# Redis 缓存
[[cache.stores]]
store = "redis"
address = "127.0.0.1:6379"
password = ""
prefix = "myapp:"
db = 0
timeout = 5000
maxidle = 10
maxactive = 100
```

#### 支持的缓存驱动

- **Memory**: 内存缓存，适合开发和小规模应用
- **Redis**: 生产环境推荐，支持分布式
- **File**: 文件缓存，简单场景使用

### 缓存操作

#### 获取缓存对象

```go
func (this *User) GetCache() gcore.Cache {
    cache, err := gcore.ProviderCache()(this.Ctx)
    if err != nil {
        this.Logger.Errorf("获取缓存失败: %v", err)
        return nil
    }
    return cache
}
```

#### 基本操作

```go
func (this *User) CacheBasics() {
    cache := this.GetCache()
    
    // 设置缓存（60秒过期）
    err := cache.Set("user:1", "张三", 60*time.Second)
    if err != nil {
        this.Logger.Error(err)
    }
    
    // 获取缓存
    value, err := cache.Get("user:1")
    if err != nil {
        this.Logger.Error(err)
    }
    fmt.Printf("用户: %s\n", value)
    
    // 检查缓存是否存在
    exists, err := cache.Has("user:1")
    if exists {
        fmt.Println("缓存存在")
    }
    
    // 删除缓存
    err = cache.Del("user:1")
    
    // 清空所有缓存
    err = cache.Clear()
}
```

#### 批量操作

```go
func (this *User) BatchCache() {
    cache := this.GetCache()
    
    // 批量设置
    values := map[string]string{
        "user:1": "张三",
        "user:2": "李四",
        "user:3": "王五",
    }
    err := cache.SetMulti(values, 300*time.Second)
    
    // 批量获取
    keys := []string{"user:1", "user:2", "user:3"}
    results, err := cache.GetMulti(keys)
    for key, value := range results {
        fmt.Printf("%s: %s\n", key, value)
    }
    
    // 批量删除
    err = cache.DelMulti(keys)
}
```

#### 计数器操作

```go
func (this *User) CounterCache() {
    cache := this.GetCache()
    
    // 初始化计数器
    cache.Set("page:views", "0", 0) // 0 表示永不过期
    
    // 自增
    newValue, err := cache.Incr("page:views")
    fmt.Printf("访问量: %d\n", newValue)
    
    // 增加指定值
    newValue, err = cache.IncrN("page:views", 10)
    
    // 自减
    newValue, err = cache.Decr("page:views")
    
    // 减少指定值
    newValue, err = cache.DecrN("page:views", 5)
}
```

### 缓存模式

#### 缓存穿透防护

```go
func (this *User) GetUserWithCache(userID int64) (*User, error) {
    cache := this.GetCache()
    cacheKey := fmt.Sprintf("user:%d", userID)
    
    // 1. 尝试从缓存获取
    cached, err := cache.Get(cacheKey)
    if err == nil && cached != "" {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // 2. 从数据库获取
    user, err := this.GetUserFromDB(userID)
    if err != nil {
        return nil, err
    }
    
    // 3. 如果用户不存在，缓存空值（防止缓存穿透）
    if user == nil {
        cache.Set(cacheKey, "null", 60*time.Second)
        return nil, nil
    }
    
    // 4. 缓存结果
    data, _ := json.Marshal(user)
    cache.Set(cacheKey, string(data), 300*time.Second)
    
    return user, nil
}
```

#### 缓存更新策略

```go
// Cache Aside 模式
func (this *User) UpdateUser(userID int64, data map[string]interface{}) error {
    // 1. 更新数据库
    err := this.UpdateUserInDB(userID, data)
    if err != nil {
        return err
    }
    
    // 2. 删除缓存
    cache := this.GetCache()
    cacheKey := fmt.Sprintf("user:%d", userID)
    cache.Del(cacheKey)
    
    return nil
}

// Write Through 模式
func (this *User) UpdateUserWriteThrough(userID int64, data map[string]interface{}) error {
    // 1. 更新数据库
    err := this.UpdateUserInDB(userID, data)
    if err != nil {
        return err
    }
    
    // 2. 更新缓存
    user, _ := this.GetUserFromDB(userID)
    cache := this.GetCache()
    cacheKey := fmt.Sprintf("user:%d", userID)
    jsonData, _ := json.Marshal(user)
    cache.Set(cacheKey, string(jsonData), 300*time.Second)
    
    return nil
}
```

#### 缓存预热

```go
func (this *User) WarmupCache() {
    cache := this.GetCache()
    
    // 获取热点数据
    hotUsers := this.GetHotUsers(100)
    
    // 预热缓存
    for _, user := range hotUsers {
        cacheKey := fmt.Sprintf("user:%d", user.ID)
        jsonData, _ := json.Marshal(user)
        cache.Set(cacheKey, string(jsonData), 3600*time.Second)
    }
    
    this.Logger.Info("缓存预热完成")
}
```

### 分布式锁

使用 Redis 实现分布式锁：

```go
func (this *User) WithLock(key string, timeout time.Duration, fn func() error) error {
    cache := this.GetCache()
    lockKey := "lock:" + key
    lockValue := fmt.Sprintf("%d", time.Now().UnixNano())
    
    // 尝试获取锁
    err := cache.Set(lockKey, lockValue, timeout)
    if err != nil {
        return errors.New("获取锁失败")
    }
    
    // 确保释放锁
    defer func() {
        // 验证锁还是自己的再删除
        if value, _ := cache.Get(lockKey); value == lockValue {
            cache.Del(lockKey)
        }
    }()
    
    // 执行业务逻辑
    return fn()
}

// 使用示例
func (this *User) ProcessOrder(orderID int64) {
    lockKey := fmt.Sprintf("order:%d", orderID)
    
    err := this.WithLock(lockKey, 10*time.Second, func() error {
        // 处理订单的业务逻辑
        return this.DoProcessOrder(orderID)
    })
    
    if err != nil {
        this.Logger.Error(err)
    }
}
```

### 自定义驱动

#### 实现缓存接口

```go
package mycache

import (
    "time"
    "github.com/snail007/gmc/core"
)

type MyCache struct {
    // 你的实现
}

func (c *MyCache) Has(key string) (bool, error) {
    // 实现
    return false, nil
}

func (c *MyCache) Get(key string) (string, error) {
    // 实现
    return "", nil
}

func (c *MyCache) Set(key string, value string, ttl time.Duration) error {
    // 实现
    return nil
}

func (c *MyCache) Del(key string) error {
    // 实现
    return nil
}

func (c *MyCache) GetMulti(keys []string) (map[string]string, error) {
    // 实现
    return nil, nil
}

func (c *MyCache) SetMulti(values map[string]string, ttl time.Duration) error {
    // 实现
    return nil
}

func (c *MyCache) DelMulti(keys []string) error {
    // 实现
    return nil
}

func (c *MyCache) Incr(key string) (int64, error) {
    // 实现
    return 0, nil
}

func (c *MyCache) Decr(key string) (int64, error) {
    // 实现
    return 0, nil
}

func (c *MyCache) IncrN(key string, n int64) (int64, error) {
    // 实现
    return 0, nil
}

func (c *MyCache) DecrN(key string, n int64) (int64, error) {
    // 实现
    return 0, nil
}

func (c *MyCache) Clear() error {
    // 实现
    return nil
}

func (c *MyCache) String() string {
    return "MyCache Driver"
}
```

#### 注册驱动

```go
func init() {
    gcore.RegisterCache("mycache", func(ctx gcore.Ctx) (gcore.Cache, error) {
        return &MyCache{}, nil
    })
}
```

### 缓存最佳实践

1. **合理的过期时间**: 根据数据特性设置合适的 TTL
2. **缓存键命名**: 使用统一的命名规范，如 `prefix:type:id`
3. **防止穿透**: 缓存空值或使用布隆过滤器
4. **防止雪崩**: 设置随机过期时间，避免同时失效
5. **防止击穿**: 使用分布式锁保护热点数据
6. **监控统计**: 监控缓存命中率和性能指标
7. **容量规划**: 合理规划缓存容量，避免内存溢出
8. **序列化**: 复杂对象使用 JSON 或 MessagePack 序列化

---

## Session

GMC Session 模块提供灵活的会话管理，支持多种存储后端（内存、Redis、文件），内置安全特性。

> **详细文档：** [http/session/README.md](https://github.com/snail007/gmc/blob/master/http/session/README.md) - 查看完整的 Session API、安全配置和高级用法

### Session 配置

#### 配置文件

在 `conf/app.toml` 中配置 Session：

```toml
[session]
enable = true
store = "memory"           # memory, redis, file
ttl = 3600                 # Session 有效期（秒）
cookiename = "gmc_session" # Cookie 名称
cookiedomain = ""          # Cookie 域名
cookiepath = "/"           # Cookie 路径
cookiesecure = false       # 是否仅 HTTPS
cookiehttponly = true      # HttpOnly 标志

# Redis 存储配置（当 store = "redis"）
[session.redis]
address = "127.0.0.1:6379"
password = ""
db = 0
prefix = "session:"

# 文件存储配置（当 store = "file"）
[session.file]
dir = "sessions"           # 存储目录
```

### 使用 Session

#### 启动 Session

```go
func (this *User) Login() {
    // 启动 Session
    err := this.SessionStart()
    if err != nil {
        this.Logger.Error(err)
        this.Ctx.JSON(500, map[string]string{"error": "Session 启动失败"})
        return
    }
    
    // 验证登录...
    if this.ValidateLogin() {
        // 存储用户信息到 Session
        this.Session.Set("user_id", 123)
        this.Session.Set("username", "zhangsan")
        this.Session.Set("role", "admin")
        this.Session.Set("login_time", time.Now())
        
        this.Ctx.JSON(200, map[string]string{"message": "登录成功"})
    }
}
```

#### 读取 Session

```go
func (this *User) Profile() {
    // 启动 Session
    err := this.SessionStart()
    if err != nil {
        this.Ctx.Redirect("/login")
        return
    }
    
    // 获取 Session 数据
    userID := this.Session.Get("user_id")
    if userID == nil {
        this.Ctx.Redirect("/login")
        return
    }
    
    // 类型断言
    uid := userID.(int)
    username := this.Session.Get("username").(string)
    
    // 显示个人资料
    this.View.Set("user_id", uid).
        Set("username", username).
        Render("user/profile")
}
```

#### 删除 Session 数据

```go
func (this *User) Logout() {
    err := this.SessionStart()
    if err != nil {
        return
    }
    
    // 方式 1: 销毁整个 Session
    this.SessionDestroy()
    
    // 方式 2: 删除特定的键
    // this.Session.Delete("user_id")
    // this.Session.Delete("username")
    
    this.Ctx.Redirect("/")
}
```

### Session 中间件

创建登录检查中间件：

```go
func AuthMiddleware(ctx gcore.Ctx) bool {
    // 白名单路径
    whitelist := []string{"/login", "/register", "/"}
    path := ctx.Request().URL.Path
    
    for _, p := range whitelist {
        if p == path {
            return false // 继续处理
        }
    }
    
    // 启动 Session
    sess, err := ctx.SessionStart()
    if err != nil {
        ctx.Redirect("/login")
        return true // 停止处理
    }
    
    // 检查登录状态
    userID := sess.Get("user_id")
    if userID == nil {
        ctx.Redirect("/login")
        return true
    }
    
    // 将用户信息存到上下文
    ctx.Set("user_id", userID)
    ctx.Set("username", sess.Get("username"))
    
    return false // 继续处理
}

// 注册中间件
func InitMiddleware(s gmc.HTTPServer) {
    s.AddMiddleware1(AuthMiddleware)
}
```

### Session 封装

#### Session 工具类

```go
package session

import (
    "github.com/snail007/gmc/core"
)

type SessionHelper struct {
    ctx gcore.Ctx
}

func New(ctx gcore.Ctx) *SessionHelper {
    return &SessionHelper{ctx: ctx}
}

func (s *SessionHelper) Start() (gcore.Session, error) {
    return s.ctx.SessionStart()
}

func (s *SessionHelper) IsLoggedIn() bool {
    sess, err := s.Start()
    if err != nil {
        return false
    }
    return sess.Get("user_id") != nil
}

func (s *SessionHelper) GetUserID() int64 {
    sess, err := s.Start()
    if err != nil {
        return 0
    }
    
    if uid := sess.Get("user_id"); uid != nil {
        return uid.(int64)
    }
    return 0
}

func (s *SessionHelper) GetUsername() string {
    sess, err := s.Start()
    if err != nil {
        return ""
    }
    
    if name := sess.Get("username"); name != nil {
        return name.(string)
    }
    return ""
}

func (s *SessionHelper) SetUser(userID int64, username string, role string) {
    sess, _ := s.Start()
    sess.Set("user_id", userID)
    sess.Set("username", username)
    sess.Set("role", role)
    sess.Set("login_at", time.Now().Unix())
}

func (s *SessionHelper) Logout() {
    sess, _ := s.Start()
    sess.Destroy()
}

// 使用示例
func (this *User) Login() {
    sh := session.New(this.Ctx)
    
    // 验证登录
    user, err := this.ValidateLogin()
    if err != nil {
        this.Ctx.JSON(401, map[string]string{"error": "登录失败"})
        return
    }
    
    // 设置 Session
    sh.SetUser(user.ID, user.Name, user.Role)
    
    this.Ctx.JSON(200, map[string]string{"message": "登录成功"})
}

func (this *User) Dashboard() {
    sh := session.New(this.Ctx)
    
    if !sh.IsLoggedIn() {
        this.Ctx.Redirect("/login")
        return
    }
    
    userID := sh.GetUserID()
    username := sh.GetUsername()
    
    // 显示控制台...
}
```

### Session 存储

#### 自定义 Session 存储

```go
package storage

import (
    "github.com/snail007/gmc/core"
    "time"
)

type CustomStorage struct {
    // 你的存储实现
}

func (s *CustomStorage) Load(sessionID string) (gcore.Session, bool) {
    // 从存储加载 Session
    return nil, false
}

func (s *CustomStorage) Save(session gcore.Session) error {
    // 保存 Session 到存储
    return nil
}

func (s *CustomStorage) Delete(sessionID string) error {
    // 从存储删除 Session
    return nil
}

// 注册存储驱动
func init() {
    gcore.RegisterSessionStorage("custom", func(ctx gcore.Ctx) (gcore.SessionStorage, error) {
        return &CustomStorage{}, nil
    })
}
```

### Session 安全

#### Session 固定攻击防护

```go
func (this *User) Login() {
    // 登录成功后重新生成 Session ID
    err := this.SessionStart()
    if err != nil {
        return
    }
    
    // 验证登录
    if this.ValidateLogin() {
        // 获取旧的 Session 数据
        oldData := this.Session.Values()
        
        // 销毁旧 Session
        this.SessionDestroy()
        
        // 创建新 Session
        this.SessionStart()
        
        // 恢复数据
        for k, v := range oldData {
            this.Session.Set(k, v)
        }
        
        // 设置新的用户信息
        this.Session.Set("user_id", user.ID)
        this.Session.Set("username", user.Name)
    }
}
```

#### Session 超时检查

```go
func SessionTimeoutMiddleware(ctx gcore.Ctx) bool {
    sess, err := ctx.SessionStart()
    if err != nil {
        return false
    }
    
    // 检查最后活动时间
    lastActive := sess.Get("last_active")
    if lastActive != nil {
        lastTime := lastActive.(time.Time)
        
        // 超过 30 分钟未活动，清除 Session
        if time.Since(lastTime) > 30*time.Minute {
            sess.Destroy()
            ctx.Redirect("/login?timeout=1")
            return true
        }
    }
    
    // 更新最后活动时间
    sess.Set("last_active", time.Now())
    
    return false
}
```

### Session 最佳实践

1. **安全的 Session ID**: 使用足够长且随机的 Session ID
2. **HttpOnly Cookie**: 防止 XSS 攻击窃取 Session
3. **HTTPS**: 生产环境使用 HTTPS 传输 Session Cookie
4. **Session 超时**: 设置合理的过期时间
5. **重新生成 ID**: 登录后重新生成 Session ID
6. **最小化存储**: 只在 Session 中存储必要的数据
7. **分布式存储**: 多服务器部署使用 Redis 等共享存储
8. **定期清理**: 定期清理过期的 Session 数据

---

## 日志

GMC Log 模块提供强大的日志记录功能，支持多级别、多种输出格式、异步日志、日志轮转、结构化日志等特性。

> **详细文档：** [module/log/README.md](https://github.com/snail007/gmc/blob/master/module/log/README.md) - 查看完整的 API 文档、输出格式、日志轮转配置等

### 日志配置

#### 配置文件

在 `conf/app.toml` 中配置日志：

```toml
[log]
level = "info"              # 日志级别: trace, debug, info, warn, error, fatal
output = "console"          # 输出方式: console, file, both
async = false               # 是否异步写入
filename = "logs/app.log"   # 日志文件路径
maxsize = 100               # 单个文件最大大小（MB）
maxbackups = 10             # 保留的旧文件数量
maxage = 30                 # 文件保留天数
compress = true             # 是否压缩归档
```

#### 日志级别

GMC 支持以下日志级别（从低到高）：

- **Trace**: 跟踪级别，最详细的信息
- **Debug**: 调试信息
- **Info**: 一般信息
- **Warn**: 警告信息
- **Error**: 错误信息
- **Fatal**: 致命错误，记录后程序退出

### 基本使用

#### 在控制器中使用

```go
func (this *User) Example() {
    // 基本日志
    this.Logger.Info("用户访问了首页")
    this.Logger.Debug("调试信息")
    this.Logger.Warn("警告信息")
    this.Logger.Error("错误信息")
    
    // 格式化日志
    this.Logger.Infof("用户 %s 登录成功", username)
    this.Logger.Debugf("请求参数: %+v", params)
    this.Logger.Warnf("访问频率过高: %d 次/分钟", rate)
    this.Logger.Errorf("数据库错误: %v", err)
    
    // 带字段的结构化日志
    this.Logger.With("user_id", 123).
        With("action", "login").
        Info("用户登录")
    
    // 记录错误堆栈
    if err != nil {
        this.Logger.WithError(err).Error("操作失败")
    }
}
```

#### 在其他包中使用

```go
package service

import (
    "github.com/snail007/gmc"
)

type UserService struct {
    logger gcore.Logger
}

func NewUserService(logger gcore.Logger) *UserService {
    return &UserService{
        logger: logger,
    }
}

func (s *UserService) CreateUser(user *User) error {
    s.logger.Infof("创建用户: %s", user.Name)
    
    err := s.saveUser(user)
    if err != nil {
        s.logger.Errorf("保存用户失败: %v", err)
        return err
    }
    
    s.logger.With("user_id", user.ID).
        Info("用户创建成功")
    
    return nil
}
```

### 结构化日志

#### 使用字段

```go
func (this *User) Login() {
    username := this.Ctx.POST("username")
    
    // 添加多个字段
    logger := this.Logger.
        With("username", username).
        With("ip", this.Ctx.ClientIP()).
        With("user_agent", this.Ctx.Header("User-Agent"))
    
    // 验证登录
    user, err := this.ValidateLogin(username)
    if err != nil {
        logger.Error("登录失败")
        return
    }
    
    logger.With("user_id", user.ID).Info("登录成功")
}
```

#### 使用 Map 添加字段

```go
func (this *User) LogRequest() {
    fields := map[string]interface{}{
        "method":     this.Request.Method,
        "path":       this.Request.URL.Path,
        "query":      this.Request.URL.RawQuery,
        "client_ip":  this.Ctx.ClientIP(),
        "duration":   time.Since(startTime),
    }
    
    this.Logger.WithFields(fields).Info("请求完成")
}
```

### 日志中间件

创建访问日志中间件：

```go
func AccessLogMiddleware(ctx gcore.Ctx) bool {
    startTime := time.Now()
    logger := ctx.Logger()
    
    // 记录请求信息
    logger.With("method", ctx.Request().Method).
        With("path", ctx.Request().URL.Path).
        With("client_ip", ctx.ClientIP()).
        Info("收到请求")
    
    // 继续处理请求
    defer func() {
        duration := time.Since(startTime)
        
        // 记录响应信息
        logger.With("duration", duration).
            With("status", ctx.StatusCode()).
            Info("请求完成")
    }()
    
    return false
}

// 注册中间件
func InitMiddleware(s gmc.HTTPServer) {
    s.AddMiddleware1(AccessLogMiddleware)
}
```

### 自定义日志格式

#### 创建自定义 Logger

```go
package logger

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/core"
)

type CustomLogger struct {
    prefix string
    level  string
}

func NewCustomLogger(prefix string) *CustomLogger {
    return &CustomLogger{
        prefix: prefix,
        level:  "info",
    }
}

func (l *CustomLogger) Info(msg string) {
    l.log("INFO", msg)
}

func (l *CustomLogger) Infof(format string, args ...interface{}) {
    l.log("INFO", fmt.Sprintf(format, args...))
}

func (l *CustomLogger) Debug(msg string) {
    l.log("DEBUG", msg)
}

func (l *CustomLogger) Debugf(format string, args ...interface{}) {
    l.log("DEBUG", fmt.Sprintf(format, args...))
}

func (l *CustomLogger) Warn(msg string) {
    l.log("WARN", msg)
}

func (l *CustomLogger) Warnf(format string, args ...interface{}) {
    l.log("WARN", fmt.Sprintf(format, args...))
}

func (l *CustomLogger) Error(msg string) {
    l.log("ERROR", msg)
}

func (l *CustomLogger) Errorf(format string, args ...interface{}) {
    l.log("ERROR", fmt.Sprintf(format, args...))
}

func (l *CustomLogger) log(level, msg string) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("[%s] [%s] %s: %s\n", timestamp, level, l.prefix, msg)
}

// 实现其他 gcore.Logger 接口方法...
```

#### 注册自定义 Logger

```go
func init() {
    gcore.RegisterLogger("custom", func(ctx gcore.Ctx, prefix string) gcore.Logger {
        return NewCustomLogger(prefix)
    })
}
```

### 日志轮转

GMC 默认支持日志文件轮转，通过配置实现：

```toml
[log]
output = "file"
filename = "logs/app.log"
maxsize = 100               # 单个文件最大 100MB
maxbackups = 10             # 保留 10 个备份文件
maxage = 30                 # 保留 30 天
compress = true             # 压缩旧文件
```

生成的日志文件：
```
logs/
├── app.log                 # 当前日志
├── app-2024-01-01.log.gz  # 归档日志
├── app-2024-01-02.log.gz
└── app-2024-01-03.log.gz
```

### 日志最佳实践

1. **合适的级别**: 开发用 Debug，生产用 Info 或 Warn
2. **结构化日志**: 使用字段而不是字符串拼接
3. **避免敏感信息**: 不要记录密码、密钥等敏感数据
4. **上下文信息**: 添加足够的上下文便于问题定位
5. **性能考虑**: 生产环境使用异步日志
6. **日志归档**: 定期清理或归档旧日志
7. **统一格式**: 团队统一日志格式规范
8. **错误追踪**: 记录错误堆栈和相关上下文

---

## 国际化

GMC I18n 模块提供完整的国际化支持，支持多语言文件、占位符、复数形式等特性。

> **详细文档：** [module/i18n/README.md](https://github.com/snail007/gmc/blob/master/module/i18n/README.md) - 查看完整的国际化 API、语言文件格式和高级用法

### 配置国际化

#### 配置文件

在 `conf/app.toml` 中启用国际化：

```toml
[i18n]
enable = true
dir = "i18n"                # 语言文件目录
default = "zh-CN"           # 默认语言
```

#### 语言文件

在 `i18n` 目录创建语言文件：

`i18n/zh-CN.toml`:
```toml
hello = "你好"
welcome = "欢迎使用 GMC 框架"
user_not_found = "用户不存在"
login_success = "登录成功"
login_failed = "登录失败：用户名或密码错误"

# 支持占位符
greeting = "你好，%s！"
items_count = "共有 %d 个项目"
user_info = "用户：%s，邮箱：%s"
```

`i18n/en-US.toml`:
```toml
hello = "Hello"
welcome = "Welcome to GMC Framework"
user_not_found = "User not found"
login_success = "Login successful"
login_failed = "Login failed: incorrect username or password"

greeting = "Hello, %s!"
items_count = "Total %d items"
user_info = "User: %s, Email: %s"
```

### 使用翻译

#### 在控制器中使用

```go
func (this *User) Index() {
    // 获取翻译
    welcome := this.Tr("welcome")
    this.Write(welcome)
    
    // 带占位符的翻译
    username := "张三"
    greeting := this.Tr("greeting", username)
    // 输出: "你好，张三！" 或 "Hello, 张三!"
    
    // 多个占位符
    info := this.Tr("user_info", "张三", "zhangsan@example.com")
    
    // 设置视图变量
    this.View.Set("welcome", this.Tr("welcome"))
    this.View.Render("index")
}
```

#### 在模板中使用

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{tr .Lang "welcome"}}</title>
</head>
<body>
    <h1>{{tr .Lang "hello"}}</h1>
    
    <!-- 带占位符 -->
    <p>{{printf (trs .Lang "greeting") .Username}}</p>
    
    <!-- 或使用 tr（返回 HTML） -->
    <div>{{tr .Lang "welcome"}}</div>
    
    <!-- 使用 trs（返回字符串） -->
    <input placeholder="{{trs .Lang "enter_name"}}">
</body>
</html>
```

### 语言切换

#### 根据 HTTP 头自动切换

GMC 会自动根据 `Accept-Language` 请求头选择语言：

```go
func (this *User) Index() {
    // this.Lang 已经根据客户端语言设置
    // zh-CN, en-US, ja-JP 等
    
    this.Logger.Infof("当前语言: %s", this.Lang)
}
```

#### 手动设置语言

```go
func (this *User) ChangeLanguage() {
    lang := this.Ctx.GET("lang") // 如: zh-CN, en-US
    
    // 验证语言
    supportedLangs := []string{"zh-CN", "en-US", "ja-JP"}
    if !contains(supportedLangs, lang) {
        lang = "zh-CN"
    }
    
    // 保存到 Session
    this.SessionStart()
    this.Session.Set("lang", lang)
    
    // 保存到 Cookie（30天）
    this.Ctx.SetCookie("lang", lang, 30*24*3600, "/", "", false, false)
    
    this.Ctx.JSON(200, map[string]string{"message": "语言已切换"})
}
```

#### 语言检测中间件

```go
func LanguageMiddleware(ctx gcore.Ctx) bool {
    var lang string
    
    // 1. 优先从 URL 参数获取
    lang = ctx.GET("lang")
    
    // 2. 从 Cookie 获取
    if lang == "" {
        lang = ctx.Cookie("lang")
    }
    
    // 3. 从 Session 获取
    if lang == "" {
        sess, _ := ctx.SessionStart()
        if l := sess.Get("lang"); l != nil {
            lang = l.(string)
        }
    }
    
    // 4. 从 Accept-Language 获取
    if lang == "" {
        lang = ctx.Header("Accept-Language")
        // 解析并选择最佳匹配
        lang = parseBestLanguage(lang)
    }
    
    // 5. 使用默认语言
    if lang == "" {
        lang = "zh-CN"
    }
    
    // 设置到上下文
    ctx.Set("lang", lang)
    
    return false
}
```

### 复数形式

不同语言的复数规则不同，可以这样处理：

`i18n/zh-CN.toml`:
```toml
apple_count_zero = "没有苹果"
apple_count_one = "有 1 个苹果"
apple_count_other = "有 %d 个苹果"
```

`i18n/en-US.toml`:
```toml
apple_count_zero = "No apples"
apple_count_one = "1 apple"
apple_count_other = "%d apples"
```

使用：

```go
func (this *User) GetAppleMessage(count int) string {
    var key string
    
    if count == 0 {
        key = "apple_count_zero"
    } else if count == 1 {
        key = "apple_count_one"
    } else {
        key = "apple_count_other"
    }
    
    if count <= 1 {
        return this.Tr(key)
    }
    return this.Tr(key, count)
}
```

### 国际化最佳实践

1. **语言代码**: 使用标准的 BCP 47 语言标签（zh-CN, en-US）
2. **关键字命名**: 使用清晰的关键字，如 `user.create` 而不是 `uc`
3. **默认语言**: 始终提供默认语言的完整翻译
4. **占位符**: 为动态内容使用占位符
5. **上下文**: 相同文字不同含义时使用不同的关键字
6. **测试**: 测试所有语言版本的显示效果
7. **文档**: 维护翻译关键字文档
8. **工具**: 使用工具检查缺失的翻译

**注意**: GMC i18n 目前只支持单层目录结构，所有语言文件必须直接放在 `i18n` 目录下，如 `i18n/zh-CN.toml`、`i18n/en-US.toml`。不支持子目录结构。

---

## API 开发

### 创建 API 服务

#### 简单 API 项目

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    // 创建 API 服务器
    api := gmc.New.APIServer(":8080")
    
    // 注册路由
    api.API("/hello", func(c gmc.C) {
        c.JSON(200, map[string]interface{}{
            "message": "Hello GMC",
            "code":    200,
        })
    })
    
    // 启动服务
    if err := api.Run(); err != nil {
        panic(err)
    }
}
```

#### 完整 API 项目

使用 GMCT 生成：

```bash
gmct new api --pkg myapp
cd $GOPATH/src/myapp
gmct run
```

### RESTful API

#### 标准 REST 接口

```go
package handler

import (
    "github.com/snail007/gmc"
)

type UserHandler struct{}

// GET /api/users - 获取用户列表
func (h *UserHandler) List(c gmc.C) {
    page := c.GET("page", "1")
    pageSize := c.GET("page_size", "20")
    
    users, total, err := GetUsers(page, pageSize)
    if err != nil {
        c.JSON(500, map[string]interface{}{
            "code":    500,
            "message": "获取用户列表失败",
            "error":   err.Error(),
        })
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "code":    200,
        "message": "success",
        "data": map[string]interface{}{
            "items": users,
            "total": total,
            "page":  page,
        },
    })
}

// GET /api/users/:id - 获取单个用户
func (h *UserHandler) Get(c gmc.C) {
    id := c.Param("id")
    
    user, err := GetUserByID(id)
    if err != nil {
        c.JSON(404, map[string]interface{}{
            "code":    404,
            "message": "用户不存在",
        })
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "code":    200,
        "message": "success",
        "data":    user,
    })
}

// POST /api/users - 创建用户
func (h *UserHandler) Create(c gmc.C) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(400, map[string]interface{}{
            "code":    400,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }
    
    // 验证
    if err := user.Validate(); err != nil {
        c.JSON(400, map[string]interface{}{
            "code":    400,
            "message": "数据验证失败",
            "errors":  err,
        })
        return
    }
    
    // 创建
    if err := CreateUser(&user); err != nil {
        c.JSON(500, map[string]interface{}{
            "code":    500,
            "message": "创建用户失败",
        })
        return
    }
    
    c.JSON(201, map[string]interface{}{
        "code":    201,
        "message": "创建成功",
        "data":    user,
    })
}

// PUT /api/users/:id - 更新用户
func (h *UserHandler) Update(c gmc.C) {
    id := c.Param("id")
    
    var user User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(400, map[string]interface{}{
            "code":    400,
            "message": "请求参数错误",
        })
        return
    }
    
    user.ID = id
    if err := UpdateUser(&user); err != nil {
        c.JSON(500, map[string]interface{}{
            "code":    500,
            "message": "更新失败",
        })
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "code":    200,
        "message": "更新成功",
        "data":    user,
    })
}

// DELETE /api/users/:id - 删除用户
func (h *UserHandler) Delete(c gmc.C) {
    id := c.Param("id")
    
    if err := DeleteUser(id); err != nil {
        c.JSON(500, map[string]interface{}{
            "code":    500,
            "message": "删除失败",
        })
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "code":    200,
        "message": "删除成功",
    })
}

// 注册路由
func RegisterUserRoutes(api gmc.APIServer) {
    handler := &UserHandler{}
    
    api.API("GET", "/api/users", handler.List)
    api.API("GET", "/api/users/:id", handler.Get)
    api.API("POST", "/api/users", handler.Create)
    api.API("PUT", "/api/users/:id", handler.Update)
    api.API("DELETE", "/api/users/:id", handler.Delete)
}
```

### API 认证

#### JWT 认证

```go
package middleware

import (
    "strings"
    "github.com/snail007/gmc"
    "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your-secret-key")

// JWT 认证中间件
func JWTAuth(c gmc.C) {
    // 白名单
    whitelist := []string{"/api/login", "/api/register"}
    path := c.Request().URL.Path
    
    for _, p := range whitelist {
        if p == path {
            return
        }
    }
    
    // 获取 Token
    authHeader := c.Header("Authorization")
    if authHeader == "" {
        c.JSON(401, map[string]interface{}{
            "code":    401,
            "message": "未授权：缺少 Authorization 头",
        })
        c.Stop()
        return
    }
    
    // 解析 Token
    parts := strings.SplitN(authHeader, " ", 2)
    if len(parts) != 2 || parts[0] != "Bearer" {
        c.JSON(401, map[string]interface{}{
            "code":    401,
            "message": "未授权：Authorization 格式错误",
        })
        c.Stop()
        return
    }
    
    tokenString := parts[1]
    
    // 验证 Token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil || !token.Valid {
        c.JSON(401, map[string]interface{}{
            "code":    401,
            "message": "未授权：Token 无效",
        })
        c.Stop()
        return
    }
    
    // 提取用户信息
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        c.Set("user_id", claims["user_id"])
        c.Set("username", claims["username"])
    }
}

// 生成 Token
func GenerateToken(userID int64, username string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  userID,
        "username": username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    })
    
    return token.SignedString(jwtSecret)
}

// 登录处理
func Login(c gmc.C) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, map[string]string{"error": "参数错误"})
        return
    }
    
    // 验证用户
    user, err := ValidateUser(req.Username, req.Password)
    if err != nil {
        c.JSON(401, map[string]string{"error": "用户名或密码错误"})
        return
    }
    
    // 生成 Token
    token, err := GenerateToken(user.ID, user.Username)
    if err != nil {
        c.JSON(500, map[string]string{"error": "生成 Token 失败"})
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "code":    200,
        "message": "登录成功",
        "data": map[string]interface{}{
            "token": token,
            "user":  user,
        },
    })
}
```

#### API Key 认证

```go
func APIKeyAuth(c gmc.C) {
    apiKey := c.Header("X-API-Key")
    
    if apiKey == "" {
        c.JSON(401, map[string]string{
            "error": "缺少 API Key",
        })
        c.Stop()
        return
    }
    
    // 验证 API Key
    if !ValidateAPIKey(apiKey) {
        c.JSON(401, map[string]string{
            "error": "无效的 API Key",
        })
        c.Stop()
        return
    }
    
    // 获取 API Key 关联的信息
    app := GetAppByAPIKey(apiKey)
    c.Set("app_id", app.ID)
    c.Set("app_name", app.Name)
}
```

### 请求验证

#### 数据验证

```go
package validator

import (
    "regexp"
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Age      int    `json:"age" validate:"required,gte=0,lte=150"`
    Phone    string `json:"phone" validate:"omitempty,phone"`
}

// 自定义验证规则
func init() {
    validate.RegisterValidation("phone", validatePhone)
}

func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
    return matched
}

// 验证请求
func ValidateCreateUser(req *CreateUserRequest) error {
    return validate.Struct(req)
}

// 格式化验证错误
func FormatValidationErrors(err error) map[string]string {
    errors := make(map[string]string)
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            field := e.Field()
            switch e.Tag() {
            case "required":
                errors[field] = field + " 是必填项"
            case "email":
                errors[field] = "邮箱格式不正确"
            case "min":
                errors[field] = field + " 长度不能少于 " + e.Param()
            case "max":
                errors[field] = field + " 长度不能超过 " + e.Param()
            default:
                errors[field] = field + " 验证失败"
            }
        }
    }
    
    return errors
}

// 使用示例
func CreateUserHandler(c gmc.C) {
    var req CreateUserRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, map[string]interface{}{
            "code":    400,
            "message": "请求参数错误",
            "error":   err.Error(),
        })
        return
    }
    
    // 验证
    if err := ValidateCreateUser(&req); err != nil {
        c.JSON(400, map[string]interface{}{
            "code":    400,
            "message": "数据验证失败",
            "errors":  FormatValidationErrors(err),
        })
        return
    }
    
    // 创建用户...
}
```

### 错误处理

#### 统一错误响应

```go
package response

type ErrorCode int

const (
    CodeSuccess         ErrorCode = 200
    CodeBadRequest      ErrorCode = 400
    CodeUnauthorized    ErrorCode = 401
    CodeForbidden       ErrorCode = 403
    CodeNotFound        ErrorCode = 404
    CodeInternalError   ErrorCode = 500
    CodeServiceUnavailable ErrorCode = 503
)

type Response struct {
    Code    ErrorCode   `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func Success(c gmc.C, data interface{}) {
    c.JSON(200, Response{
        Code:    CodeSuccess,
        Message: "success",
        Data:    data,
    })
}

func Error(c gmc.C, code ErrorCode, message string, err error) {
    resp := Response{
        Code:    code,
        Message: message,
    }
    
    if err != nil {
        resp.Error = err.Error()
    }
    
    c.JSON(int(code), resp)
}

func BadRequest(c gmc.C, message string) {
    Error(c, CodeBadRequest, message, nil)
}

func Unauthorized(c gmc.C, message string) {
    Error(c, CodeUnauthorized, message, nil)
}

func NotFound(c gmc.C, message string) {
    Error(c, CodeNotFound, message, nil)
}

func InternalError(c gmc.C, err error) {
    Error(c, CodeInternalError, "服务器内部错误", err)
}
```

#### 错误恢复中间件

```go
func RecoverMiddleware(c gmc.C) {
    defer func() {
        if err := recover(); err != nil {
            // 记录错误
            c.Logger().Errorf("Panic recovered: %v", err)
            
            // 返回错误响应
            c.JSON(500, map[string]interface{}{
                "code":    500,
                "message": "服务器内部错误",
            })
        }
    }()
}
```

### API 文档

#### Swagger 集成

```go
// 使用 swaggo 生成文档
// @title GMC API
// @version 1.0
// @description GMC 框架 API 接口文档
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// GetUser godoc
// @Summary 获取用户信息
// @Description 根据用户 ID 获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} Response{data=User}
// @Failure 404 {object} Response
// @Router /users/{id} [get]
// @Security ApiKeyAuth
func GetUser(c gmc.C) {
    // 处理逻辑
}
```

### API 版本控制

```go
func RegisterRoutes(api gmc.APIServer) {
    // V1 API
    v1 := api.Group("/api/v1")
    {
        v1.API("GET", "/users", V1GetUsers)
        v1.API("POST", "/users", V1CreateUser)
    }
    
    // V2 API
    v2 := api.Group("/api/v2")
    {
        v2.API("GET", "/users", V2GetUsers)
        v2.API("POST", "/users", V2CreateUser)
    }
}
```

### API 最佳实践

1. **版本控制**: 使用 URL 或 Header 进行版本管理
2. **统一响应**: 使用统一的响应格式
3. **错误处理**: 提供清晰的错误信息
4. **认证授权**: 保护敏感接口
5. **限流**: 防止 API 滥用
6. **文档**: 维护完整的 API 文档
7. **测试**: 编写完整的 API 测试
8. **监控**: 监控 API 性能和错误率

---

## 测试

### 单元测试

#### 测试控制器

```go
package controller_test

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/snail007/gmc"
    "myapp/controller"
)

func TestUserController_Index(t *testing.T) {
    // 创建测试服务器
    s := gmc.New.HTTPServer(gmc.New.Ctx(), ":0")
    
    // 注册路由
    s.Router().Controller("/user", new(controller.User))
    
    // 创建测试请求
    req := httptest.NewRequest("GET", "/user/index", nil)
    w := httptest.NewRecorder()
    
    // 执行请求
    s.ServeHTTP(w, req)
    
    // 验证响应
    if w.Code != http.StatusOK {
        t.Errorf("期望状态码 200，得到 %d", w.Code)
    }
    
    body := w.Body.String()
    if !strings.Contains(body, "用户列表") {
        t.Error("响应内容不符合预期")
    }
}

func TestUserController_Create(t *testing.T) {
    // 准备测试数据
    data := `{"name":"test","email":"test@example.com"}`
    
    req := httptest.NewRequest("POST", "/user/create", strings.NewReader(data))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    
    // 执行请求
    s.ServeHTTP(w, req)
    
    // 验证响应
    if w.Code != http.StatusCreated {
        t.Errorf("期望状态码 201，得到 %d", w.Code)
    }
}
```

#### 测试模型

```go
package model_test

import (
    "testing"
    "myapp/model"
)

func TestUser_Validate(t *testing.T) {
    tests := []struct {
        name    string
        user    *model.User
        wantErr bool
    }{
        {
            name: "有效用户",
            user: &model.User{
                Name:  "test",
                Email: "test@example.com",
                Age:   25,
            },
            wantErr: false,
        },
        {
            name: "邮箱格式错误",
            user: &model.User{
                Name:  "test",
                Email: "invalid-email",
                Age:   25,
            },
            wantErr: true,
        },
        {
            name: "年龄无效",
            user: &model.User{
                Name:  "test",
                Email: "test@example.com",
                Age:   -1,
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### 测试服务层

```go
package service_test

import (
    "testing"
    "myapp/service"
)

func TestUserService_CreateUser(t *testing.T) {
    // 创建测试服务
    svc := service.NewUserService()
    
    // 测试用例
    user := &model.User{
        Name:  "test",
        Email: "test@example.com",
    }
    
    // 执行测试
    err := svc.CreateUser(user)
    
    // 验证结果
    if err != nil {
        t.Errorf("CreateUser() error = %v", err)
    }
    
    if user.ID == 0 {
        t.Error("用户 ID 应该被设置")
    }
}
```

### HTTP 测试

#### 集成测试

```go
package integration_test

import (
    "testing"
    "net/http"
    "encoding/json"
    "bytes"
    "github.com/snail007/gmc"
)

var testServer gmc.HTTPServer

func TestMain(m *testing.M) {
    // 启动测试服务器
    testServer = startTestServer()
    
    // 运行测试
    code := m.Run()
    
    // 清理
    testServer.Shutdown()
    
    os.Exit(code)
}

func TestAPI_Users(t *testing.T) {
    // 测试创建用户
    t.Run("CreateUser", func(t *testing.T) {
        data := map[string]interface{}{
            "name":  "test",
            "email": "test@example.com",
        }
        
        body, _ := json.Marshal(data)
        resp, err := http.Post(
            "http://localhost:8080/api/users",
            "application/json",
            bytes.NewReader(body),
        )
        
        if err != nil {
            t.Fatal(err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusCreated {
            t.Errorf("期望状态码 201，得到 %d", resp.StatusCode)
        }
        
        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        
        if result["code"].(float64) != 201 {
            t.Error("响应码不正确")
        }
    })
    
    // 测试获取用户列表
    t.Run("ListUsers", func(t *testing.T) {
        resp, err := http.Get("http://localhost:8080/api/users")
        if err != nil {
            t.Fatal(err)
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            t.Errorf("期望状态码 200，得到 %d", resp.StatusCode)
        }
    })
}
```

#### API 客户端测试

```go
package client_test

import (
    "testing"
    "myapp/client"
)

func TestAPIClient_GetUser(t *testing.T) {
    client := client.NewAPIClient("http://localhost:8080")
    
    user, err := client.GetUser(1)
    if err != nil {
        t.Fatalf("GetUser() error = %v", err)
    }
    
    if user.ID != 1 {
        t.Errorf("期望用户 ID 为 1，得到 %d", user.ID)
    }
}
```

### 数据库测试

#### 使用测试数据库

```go
package database_test

import (
    "testing"
    "github.com/snail007/gmc"
)

var testDB gcore.Database

func setupTestDB(t *testing.T) {
    // 创建测试数据库配置
    cfg := gmc.New.Config()
    cfg.Set("database.driver", "sqlite3")
    cfg.Set("database.dsn", ":memory:")
    
    // 初始化数据库
    db, err := gcore.ProviderDatabase()(gmc.New.Ctx())
    if err != nil {
        t.Fatal(err)
    }
    
    testDB = db
    
    // 创建测试表
    createTables(testDB)
}

func teardownTestDB() {
    if testDB != nil {
        testDB.Close()
    }
}

func TestDatabase_CreateUser(t *testing.T) {
    setupTestDB(t)
    defer teardownTestDB()
    
    // 插入测试数据
    ar := testDB.AR()
    ar.Insert("users", map[string]interface{}{
        "name":  "test",
        "email": "test@example.com",
    })
    
    result, err := testDB.Exec(ar)
    if err != nil {
        t.Fatalf("插入失败: %v", err)
    }
    
    id := result.LastInsertID()
    if id == 0 {
        t.Error("应该返回插入的 ID")
    }
    
    // 验证数据
    ar = testDB.AR()
    ar.Select("*").From("users").Where(map[string]interface{}{
        "id": id,
    })
    
    queryResult, err := testDB.Query(ar)
    if err != nil {
        t.Fatalf("查询失败: %v", err)
    }
    
    row := queryResult.Row()
    if row["name"] != "test" {
        t.Error("名称不匹配")
    }
}
```

#### 事务测试

```go
func TestDatabase_Transaction(t *testing.T) {
    setupTestDB(t)
    defer teardownTestDB()
    
    tx, err := testDB.Begin()
    if err != nil {
        t.Fatal(err)
    }
    
    // 插入数据
    ar := testDB.AR()
    ar.Insert("users", map[string]interface{}{
        "name": "test",
    })
    
    _, err = testDB.ExecTx(ar, tx)
    if err != nil {
        tx.Rollback()
        t.Fatal(err)
    }
    
    // 回滚
    tx.Rollback()
    
    // 验证数据未插入
    ar = testDB.AR()
    ar.Select("COUNT(*) as count").From("users")
    result, _ := testDB.Query(ar)
    
    row := result.Row()
    if row["count"] != "0" {
        t.Error("事务回滚失败")
    }
}
```

### Mock 测试

#### Mock 数据库

```go
package mock

import (
    "github.com/snail007/gmc/core"
)

type MockDatabase struct {
    QueryFunc func(ar gcore.DBQueryBuilder) (gcore.DBResultSet, error)
    ExecFunc  func(ar gcore.DBQueryBuilder) (gcore.DBResult, error)
}

func (m *MockDatabase) Query(ar gcore.DBQueryBuilder) (gcore.DBResultSet, error) {
    if m.QueryFunc != nil {
        return m.QueryFunc(ar)
    }
    return nil, nil
}

func (m *MockDatabase) Exec(ar gcore.DBQueryBuilder) (gcore.DBResult, error) {
    if m.ExecFunc != nil {
        return m.ExecFunc(ar)
    }
    return nil, nil
}

// 使用 Mock
func TestUserService_GetUser(t *testing.T) {
    mockDB := &MockDatabase{
        QueryFunc: func(ar gcore.DBQueryBuilder) (gcore.DBResultSet, error) {
            // 返回模拟数据
            return &MockResultSet{
                rows: []map[string]string{
                    {"id": "1", "name": "test"},
                },
            }, nil
        },
    }
    
    svc := service.NewUserService(mockDB)
    user, err := svc.GetUser(1)
    
    if err != nil {
        t.Fatal(err)
    }
    
    if user.Name != "test" {
        t.Error("用户名不匹配")
    }
}
```

### 基准测试

```go
func BenchmarkUserController_Index(b *testing.B) {
    s := setupTestServer()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        req := httptest.NewRequest("GET", "/user/index", nil)
        w := httptest.NewRecorder()
        s.ServeHTTP(w, req)
    }
}

func BenchmarkDatabase_Query(b *testing.B) {
    db := setupTestDB(b)
    defer db.Close()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ar := db.AR()
        ar.Select("*").From("users").Limit(10, 0)
        db.Query(ar)
    }
}
```

### 测试覆盖率

```bash
# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html
```

### 测试最佳实践

1. **测试组织**: 按功能模块组织测试文件
2. **命名规范**: 测试函数以 Test 开头
3. **表驱动测试**: 使用表驱动测试提高覆盖率
4. **隔离性**: 每个测试应该独立运行
5. **清理**: 使用 defer 确保资源清理
6. **Mock**: 使用 Mock 隔离外部依赖
7. **覆盖率**: 保持合理的测试覆盖率
8. **持续集成**: 将测试集成到 CI/CD

---

## 部署

### 编译打包

#### 基本编译

```bash
# 编译当前平台
go build -o myapp

# 指定输出路径
go build -o bin/myapp

# 优化编译（减小体积）
go build -ldflags="-s -w" -o myapp

# 查看二进制文件大小
ls -lh myapp
```

#### 交叉编译

```bash
# Linux 64位
GOOS=linux GOARCH=amd64 go build -o myapp-linux

# Windows 64位
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o myapp-mac-arm64

# 多平台编译脚本
#!/bin/bash
platforms=("linux/amd64" "linux/arm64" "windows/amd64" "darwin/amd64" "darwin/arm64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='myapp-'$GOOS'-'$GOARCH
    
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $output_name
    
    if [ $? -ne 0 ]; then
        echo "编译失败: $platform"
        exit 1
    fi
done

echo "所有平台编译完成"
```

### 静态文件嵌入

推荐使用 Go 1.16+ 的 `embed` 功能来嵌入静态资源，而不是使用 GMCT 打包命令。

> **详细文档：** 请查看 [资源嵌入](#资源嵌入) 章节了解如何使用 `embed` 嵌入静态文件、视图和 i18n 文件。

**使用 embed 的优势：**
- ✅ 原生 Go 功能，无需额外工具
- ✅ 类型安全，编译时检查
- ✅ 更好的 IDE 支持
- ✅ 标准化的实现方式

**快速示例：**

```go
package static

import "embed"

//go:embed *
var StaticFS embed.FS
```

详细用法请参考 [资源嵌入](#资源嵌入) 章节。

### 生产环境配置

#### 配置文件优化

生产环境 `conf/app.toml`：

```toml
[app]
debug = false
env = "production"

[httpserver]
listen = ":8080"
tlsenable = true
tlscert = "/path/to/cert.pem"
tlskey = "/path/to/key.pem"
# 生产环境启用 HTTPS

[log]
level = "info"              # 减少日志输出
output = "file"
filename = "/var/log/myapp/app.log"
maxsize = 100
maxbackups = 30
maxage = 90
compress = true
async = true                # 异步日志提升性能

[database]
maxidle = 50
maxconns = 200
maxlifetimeseconds = 3600
timeout = 5000
debug = false               # 关闭 SQL 日志

[cache]
enable = true
# 生产环境使用 Redis
[[cache.stores]]
store = "redis"
address = "redis-server:6379"
password = "your-password"
maxidle = 50
maxactive = 200

[session]
enable = true
store = "redis"             # 分布式 Session
ttl = 7200
cookiesecure = true         # HTTPS only
cookiehttponly = true
```

### Docker 部署

#### Dockerfile

```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myapp .

# 运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 app && adduser -D -u 1000 -G app app

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/myapp .
COPY --from=builder /app/conf ./conf

# 设置权限
RUN chown -R app:app /app

# 切换到非 root 用户
USER app

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./myapp"]
```

#### docker-compose.yml

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped
    networks:
      - app-network

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: myapp
      MYSQL_USER: appuser
      MYSQL_PASSWORD: apppassword
    volumes:
      - mysql-data:/var/lib/mysql
    restart: unless-stopped
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass redispassword
    volumes:
      - redis-data:/data
    restart: unless-stopped
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - app-network

volumes:
  mysql-data:
  redis-data:

networks:
  app-network:
    driver: bridge
```

#### 构建和运行

```bash
# 构建镜像
docker build -t myapp:latest .

# 运行容器
docker run -d -p 8080:8080 --name myapp myapp:latest

# 使用 docker-compose
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

### Systemd 服务

#### 创建服务文件

`/etc/systemd/system/myapp.service`:

```ini
[Unit]
Description=My GMC Application
After=network.target mysql.service redis.service
Wants=mysql.service redis.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/myapp
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=myapp

# 环境变量
Environment="APP_ENV=production"
Environment="GIN_MODE=release"

# 资源限制
LimitNOFILE=65536
LimitNPROC=32768

[Install]
WantedBy=multi-user.target
```

#### 管理服务

```bash
# 重新加载服务配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start myapp

# 停止服务
sudo systemctl stop myapp

# 重启服务
sudo systemctl restart myapp

# 开机自启
sudo systemctl enable myapp

# 查看状态
sudo systemctl status myapp

# 查看日志
sudo journalctl -u myapp -f
```

### Nginx 反向代理

#### nginx.conf

```nginx
upstream myapp {
    server 127.0.0.1:8080;
    # 多实例负载均衡
    # server 127.0.0.1:8081;
    # server 127.0.0.1:8082;
    
    keepalive 32;
}

server {
    listen 80;
    server_name example.com www.example.com;
    
    # HTTPS 重定向
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name example.com www.example.com;
    
    # SSL 证书
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # SSL 配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # 日志
    access_log /var/log/nginx/myapp_access.log;
    error_log /var/log/nginx/myapp_error.log;
    
    # 静态文件
    location /static/ {
        alias /opt/myapp/static/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
    
    # 代理到应用
    location / {
        proxy_pass http://myapp;
        proxy_http_version 1.1;
        
        # 请求头
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # 缓冲
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;
    }
    
    # 健康检查
    location /health {
        proxy_pass http://myapp;
        access_log off;
    }
}
```

### 性能优化

#### 应用层优化

```go
// 1. 使用连接池
cfg.Set("database.maxidle", 50)
cfg.Set("database.maxconns", 200)

// 2. 启用缓存
cfg.Set("cache.enable", true)

// 3. 异步日志
cfg.Set("log.async", true)

// 4. gzip 压缩
func GzipMiddleware(ctx gcore.Ctx) bool {
    // 实现 gzip 压缩
    return false
}
```

#### 系统层优化

```bash
# 增加文件描述符限制
ulimit -n 65536

# 优化 TCP 参数
sysctl -w net.core.somaxconn=32768
sysctl -w net.ipv4.tcp_max_syn_backlog=8192
sysctl -w net.ipv4.tcp_tw_reuse=1

# 设置进程优先级
nice -n -10 ./myapp
```

### 监控和告警

#### Prometheus 集成

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "myapp_requests_total",
            Help: "Total number of requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "myapp_request_duration_seconds",
            Help: "Request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(requestCounter)
    prometheus.MustRegister(requestDuration)
}

// 监控中间件
func MetricsMiddleware(ctx gcore.Ctx) bool {
    start := time.Now()
    
    defer func() {
        duration := time.Since(start).Seconds()
        method := ctx.Request().Method
        path := ctx.Request().URL.Path
        status := fmt.Sprintf("%d", ctx.StatusCode())
        
        requestCounter.WithLabelValues(method, path, status).Inc()
        requestDuration.WithLabelValues(method, path).Observe(duration)
    }()
    
    return false
}

// 暴露 metrics 端点
func RegisterMetrics(s gmc.HTTPServer) {
    s.Router().Handler("GET", "/metrics", promhttp.Handler())
}
```

### 部署最佳实践

1. **版本控制**: 使用语义化版本号
2. **蓝绿部署**: 保证零停机更新
3. **回滚策略**: 准备快速回滚方案
4. **健康检查**: 实现完善的健康检查接口
5. **日志收集**: 集中收集和分析日志
6. **监控告警**: 建立完善的监控体系
7. **备份恢复**: 定期备份数据
8. **文档维护**: 维护部署文档和运维手册

---

## GMCT 工具链

### 安装 GMCT

#### 从源码安装

```bash
# 克隆仓库
git clone https://github.com/snail007/gmct.git
cd gmct

# 编译安装
go install

# 验证安装
gmct version
```

#### go install 安装

```bash
go install github.com/snail007/gmct/cmd/gmct@latest
gmct version
```

### 项目生成

#### 创建 Web 项目

```bash
# 创建完整 Web 项目
gmct new web --pkg github.com/yourname/mywebapp

# 项目结构
mywebapp/
├── conf/
│   └── app.toml
├── controller/
│   └── demo.go
├── initialize/
│   └── initialize.go
├── router/
│   └── router.go
├── static/
│   └── jquery.js
├── views/
│   └── welcome.html
├── go.mod
├── go.sum
├── grun.toml
└── main.go
```

#### 创建 API 项目

```bash
# 创建 API 项目
gmct new api --pkg github.com/yourname/myapi

# 项目结构
myapi/
├── conf/
│   └── app.toml
├── handler/
│   └── demo.go
├── initialize/
│   └── initialize.go
├── go.mod
├── grun.toml
└── main.go
```

#### 创建简单 API

```bash
# 创建轻量级 API
gmct new api-simple --pkg github.com/yourname/simpleapi

# 只包含基本的 API 服务代码
```

### 代码生成

#### 生成控制器

```bash
# 在 controller 目录执行
cd controller

# 生成控制器
gmct controller -n User

# 生成的文件: user.go
# 包含基本的控制器结构和方法
```

生成的代码：

```go
package controller

import (
    "github.com/snail007/gmc"
)

type User struct {
    gmc.Controller
}

func (this *User) Index() {
    this.Write("User.Index")
}

func (this *User) List() {
    this.Write("User.List")
}

func (this *User) Detail() {
    this.Write("User.Detail")
}
```

#### 生成模型

```bash
# 在 model 目录执行
cd model

# 生成 MySQL 模型
gmct model -n user

# 生成 SQLite3 模型
gmct model -n user -t sqlite3

# 强制覆盖
gmct model -n user -f
```

生成的模型包含：
- 表结构定义
- CRUD 方法
- 查询方法
- 关联方法

### 热编译

#### 基本使用

```bash
# 在项目根目录执行
gmct run

# 自动监控文件变化
# 自动重新编译
# 自动重启应用
```

#### 配置 grun.toml

```toml
[build]
# 监控目录
monitor_dirs = ["."]

# 构建命令（为空则使用 go build）
cmd = ""

# 构建参数
args = ["-ldflags", "-s -w"]

# 环境变量
env = ["CGO_ENABLED=0", "GO111MODULE=on"]

# 监控的文件扩展名
include_exts = [".go", ".html", ".toml", ".yaml"]

# 额外监控的文件
include_files = []

# 忽略的文件
exclude_files = ["grun.toml"]

# 忽略的目录
exclude_dirs = ["vendor", ".git", ".idea"]
```

#### 高级配置

```toml
[build]
# 使用自定义构建脚本
cmd = "bash"
args = ["build.sh"]

# 或使用 make
cmd = "make"
args = ["build"]

# 监控多个目录
monitor_dirs = [".", "../shared"]

# ${DIR} 变量代表当前目录
exclude_dirs = ["${DIR}/vendor", "${DIR}/.git"]
```

### 资源打包

**推荐使用 Go embed 功能代替 GMCT 打包命令**

GMC 推荐使用 Go 1.16+ 原生的 `embed` 功能来嵌入静态资源、视图模板和国际化文件，而不是使用 GMCT 的打包命令。

详细的 embed 使用方法请参考 [资源嵌入](#资源嵌入) 章节。

**使用 embed 的优势：**
- ✅ Go 原生功能，无需额外工具
- ✅ 类型安全，编译时检查
- ✅ IDE 支持良好
- ✅ 更标准化的实现

**快速参考：**

```go
// static/static.go
package static
import "embed"
//go:embed *
var StaticFS embed.FS

// views/views.go
package views
import "embed"
//go:embed *
var ViewFS embed.FS

// i18n/i18n.go
package i18n
import "embed"
//go:embed *.toml
var I18nFS embed.FS
```

### 项目模板

#### 使用自定义模板

```bash
# 使用 Git 仓库模板
gmct new web --pkg myapp --template https://github.com/yourname/gmc-template

# 使用本地模板
gmct new web --pkg myapp --template /path/to/template
```

#### 创建项目模板

项目模板结构：

```
gmc-template/
├── template.json         # 模板配置
├── {{.ProjectName}}/     # 项目目录
│   ├── main.go.tpl
│   ├── conf/
│   │   └── app.toml.tpl
│   └── ...
└── README.md
```

`template.json`:

```json
{
  "name": "My GMC Template",
  "description": "Custom GMC project template",
  "version": "1.0.0",
  "variables": {
    "ProjectName": "string",
    "Author": "string",
    "Description": "string"
  }
}
```

### GMCT 最佳实践

1. **版本管理**: 团队使用相同版本的 GMCT
2. **配置共享**: 共享 grun.toml 配置
3. **模板定制**: 根据团队规范定制项目模板
4. **自动化**: 将 GMCT 集成到 CI/CD
5. **文档**: 维护工具使用文档
6. **更新**: 及时更新到最新版本

---

## 进阶主题

### 错误处理

GMC Error 模块提供增强的错误处理功能，支持堆栈跟踪、错误包装、Panic 恢复等特性。

> **详细文档：** [module/error/README.md](https://github.com/snail007/gmc/blob/master/module/error/README.md) - 查看完整的错误处理 API、堆栈跟踪、Try/Catch 模式等

**基本使用：**

```go
import gerror "github.com/snail007/gmc/module/error"

// 创建带堆栈的错误
err := gerror.New("something went wrong")

// 包装现有错误
wrappedErr := gerror.Wrap(existingErr)

// 打印完整堆栈
fmt.Println(err.ErrorStack())
```

---

### 自定义 Provider

Provider 是 GMC 的核心扩展机制，允许你替换或扩展框架的任何组件。

#### 创建自定义缓存 Provider

```go
package cache

import (
    "time"
    "github.com/snail007/gmc/core"
    "github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
    prefix string
}

func NewRedisCache(addr, password, prefix string, db int) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    
    return &RedisCache{
        client: client,
        prefix: prefix,
    }
}

func (c *RedisCache) Get(key string) (string, error) {
    return c.client.Get(ctx, c.prefix+key).Result()
}

func (c *RedisCache) Set(key string, value string, ttl time.Duration) error {
    return c.client.Set(ctx, c.prefix+key, value, ttl).Err()
}

// 实现其他 gcore.Cache 接口方法...

// 注册 Provider
func init() {
    gcore.RegisterCache("redis", func(ctx gcore.Ctx) (gcore.Cache, error) {
        cfg := ctx.Config()
        addr := cfg.GetString("cache.redis.addr")
        password := cfg.GetString("cache.redis.password")
        db := cfg.GetInt("cache.redis.db")
        prefix := cfg.GetString("cache.redis.prefix")
        
        return NewRedisCache(addr, password, prefix, db), nil
    })
}
```

#### 创建自定义日志 Provider

```go
package logger

import (
    "github.com/snail007/gmc/core"
    "go.uber.org/zap"
)

type ZapLogger struct {
    logger *zap.SugaredLogger
}

func NewZapLogger(config *zap.Config) (*ZapLogger, error) {
    logger, err := config.Build()
    if err != nil {
        return nil, err
    }
    
    return &ZapLogger{
        logger: logger.Sugar(),
    }, nil
}

func (l *ZapLogger) Info(msg string) {
    l.logger.Info(msg)
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
    l.logger.Infof(format, args...)
}

// 实现其他 gcore.Logger 接口方法...

// 注册 Provider
func init() {
    gcore.RegisterLogger("zap", func(ctx gcore.Ctx, prefix string) gcore.Logger {
        config := zap.NewProductionConfig()
        logger, _ := NewZapLogger(&config)
        return logger
    })
}
```

### 服务扩展

#### 创建自定义服务

```go
package service

import (
    "context"
    "github.com/snail007/gmc/core"
)

type EmailService struct {
    ctx    gcore.Ctx
    config gcore.Config
    logger gcore.Logger
}

func NewEmailService(ctx gcore.Ctx) *EmailService {
    return &EmailService{
        ctx:    ctx,
        config: ctx.Config(),
        logger: ctx.Logger(),
    }
}

// 实现 gcore.Service 接口
func (s *EmailService) Init(ctx gcore.Ctx) error {
    s.logger.Info("Email service initializing...")
    // 初始化 SMTP 连接等
    return nil
}

func (s *EmailService) Start(ctx gcore.Ctx) error {
    s.logger.Info("Email service started")
    return nil
}

func (s *EmailService) Stop(ctx context.Context) {
    s.logger.Info("Email service stopping...")
}

// 业务方法
func (s *EmailService) SendEmail(to, subject, body string) error {
    // 发送邮件逻辑
    return nil
}

// 注册到应用
func RegisterEmailService(app gcore.App) {
    app.AddService(gcore.ServiceItem{
        Service: NewEmailService(app.Ctx()),
        AfterInit: func(s *gcore.ServiceItem) error {
            // 初始化后的钩子
            return nil
        },
    })
}
```

### 中间件开发

#### 通用中间件模式

```go
package middleware

import (
    "time"
    "github.com/snail007/gmc/core"
)

// 中间件工厂
func NewRateLimiter(rate int) gcore.Middleware {
    limiter := newLimiter(rate)
    
    return func(ctx gcore.Ctx) bool {
        // 检查限流
        if !limiter.Allow() {
            ctx.WriteHeader(429)
            ctx.JSON(429, map[string]string{
                "error": "Too many requests",
            })
            return true // 停止处理
        }
        return false // 继续处理
    }
}

// 带配置的中间件
func NewCORS(cfg CORSConfig) gcore.Middleware {
    return func(ctx gcore.Ctx) bool {
        ctx.SetHeader("Access-Control-Allow-Origin", cfg.AllowOrigin)
        ctx.SetHeader("Access-Control-Allow-Methods", cfg.AllowMethods)
        ctx.SetHeader("Access-Control-Allow-Headers", cfg.AllowHeaders)
        
        if ctx.Request().Method == "OPTIONS" {
            ctx.WriteHeader(204)
            return true
        }
        
        return false
    }
}

// 链式中间件
func Chain(middlewares ...gcore.Middleware) gcore.Middleware {
    return func(ctx gcore.Ctx) bool {
        for _, mw := range middlewares {
            if mw(ctx) {
                return true // 任何一个返回 true 就停止
            }
        }
        return false
    }
}
```

### 性能调优

#### 内存优化

```go
// 1. 对象池复用
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func ProcessData(data []byte) string {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    
    buf.Reset()
    buf.Write(data)
    return buf.String()
}

// 2. 预分配切片
users := make([]User, 0, expectedCount)

// 3. 使用 strings.Builder
var builder strings.Builder
builder.Grow(expectedSize) // 预分配
builder.WriteString("hello")
result := builder.String()
```

#### 并发优化

```go
// 1. 工作池模式
type WorkerPool struct {
    workers   int
    taskQueue chan Task
    wg        sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
    p := &WorkerPool{
        workers:   workers,
        taskQueue: make(chan Task, workers*2),
    }
    
    p.Start()
    return p
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker()
    }
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()
    for task := range p.taskQueue {
        task.Execute()
    }
}

// 2. Context 超时控制
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := DoWork(ctx)

// 3. 批量处理
func BatchProcess(items []Item, batchSize int) error {
    for i := 0; i < len(items); i += batchSize {
        end := i + batchSize
        if end > len(items) {
            end = len(items)
        }
        
        batch := items[i:end]
        if err := processBatch(batch); err != nil {
            return err
        }
    }
    return nil
}
```

#### 数据库优化

```go
// 1. 批量插入
func BatchInsertUsers(users []User) error {
    db := GetDB()
    ar := db.AR()
    
    data := make([]map[string]interface{}, len(users))
    for i, u := range users {
        data[i] = map[string]interface{}{
            "name":  u.Name,
            "email": u.Email,
        }
    }
    
    ar.InsertBatch("users", data)
    _, err := db.Exec(ar)
    return err
}

// 2. 查询优化
// 只查询需要的字段
ar.Select("id, name, email").From("users")

// 使用索引
ar.Where(map[string]interface{}{
    "email": email, // 假设 email 有索引
})

// 3. 连接池调优
cfg.Set("database.maxidle", 50)
cfg.Set("database.maxconns", 200)
cfg.Set("database.maxlifetimeseconds", 3600)
```

### 高级路由

#### 动态路由

```go
// 根据配置动态注册路由
func DynamicRoutes(s gmc.HTTPServer, modules []string) {
    r := s.Router()
    
    for _, module := range modules {
        switch module {
        case "user":
            r.Controller("/user", new(controller.User))
        case "post":
            r.Controller("/post", new(controller.Post))
        case "admin":
            admin := r.Group("/admin")
            admin.Controller("/user", new(controller.AdminUser))
        }
    }
}
```

#### 路由版本控制

```go
// API 版本路由
func RegisterAPIRoutes(api gmc.APIServer) {
    // v1
    v1 := api.Group("/api/v1")
    v1.Middleware(AuthV1Middleware)
    v1.API("GET", "/users", V1GetUsers)
    
    // v2 - 兼容 v1
    v2 := api.Group("/api/v2")
    v2.Middleware(AuthV2Middleware)
    v2.API("GET", "/users", V2GetUsers)
}
```

### 进阶最佳实践

1. **架构设计**: 分层架构，职责清晰
2. **依赖注入**: 使用 Provider 模式管理依赖
3. **接口设计**: 面向接口编程
4. **性能监控**: 使用 pprof 分析性能
5. **代码复用**: 提取公共逻辑到中间件
6. **配置管理**: 环境配置分离
7. **错误处理**: 统一错误处理机制
8. **文档**: 代码注释和 API 文档

---

## 最佳实践

### 项目结构

#### 推荐的目录结构

```
myapp/
├── cmd/                    # 命令行工具
│   └── myapp/
│       └── main.go
├── conf/                   # 配置文件
│   ├── app.toml
│   ├── dev.toml
│   └── prod.toml
├── controller/             # 控制器
│   ├── api/               # API 控制器
│   │   └── user.go
│   ├── web/               # Web 控制器
│   │   └── user.go
│   └── base.go            # 基础控制器
├── model/                  # 数据模型
│   ├── user.go
│   └── post.go
├── service/                # 业务逻辑层
│   ├── user_service.go
│   └── post_service.go
├── repository/             # 数据访问层
│   ├── user_repository.go
│   └── post_repository.go
├── middleware/             # 中间件
│   ├── auth.go
│   ├── logger.go
│   └── cors.go
├── handler/                # API 处理器
│   └── user.go
├── router/                 # 路由配置
│   ├── api.go
│   └── web.go
├── initialize/             # 初始化
│   ├── database.go
│   ├── cache.go
│   └── router.go
├── pkg/                    # 内部包
│   ├── util/
│   │   ├── hash.go
│   │   └── validator.go
│   └── errors/
│       └── errors.go
├── static/                 # 静态文件
│   ├── css/
│   ├── js/
│   └── images/
├── views/                  # 视图模板
│   ├── layout/
│   ├── user/
│   └── common/
├── i18n/                   # 国际化
│   ├── zh-CN.toml
│   └── en-US.toml
├── tests/                  # 测试文件
│   ├── controller/
│   ├── service/
│   └── integration/
├── docs/                   # 文档
│   ├── api.md
│   └── deployment.md
├── scripts/                # 脚本
│   ├── build.sh
│   └── deploy.sh
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── README.md
└── CHANGELOG.md
```

### 代码规范

#### 命名规范

```go
// 1. 包名：小写，简短
package controller
package service

// 2. 接口：以 er 结尾或使用 I 前缀
type Reader interface{}
type IUserService interface{}

// 3. 结构体：驼峰命名
type UserController struct{}
type PostService struct{}

// 4. 方法：驼峰命名，公开方法首字母大写
func (u *User) GetProfile() {}
func (u *User) validateEmail() {}

// 5. 常量：驼峰或全大写
const MaxRetry = 3
const STATUS_ACTIVE = 1

// 6. 变量：驼峰命名
var userCount int
var isActive bool
```

#### 注释规范

```go
// Package controller 提供 HTTP 请求处理控制器
package controller

// User 用户控制器
// 处理用户相关的 HTTP 请求
type User struct {
    gmc.Controller
}

// List 获取用户列表
//
// GET /user/list?page=1&size=20
//
// 参数:
//   - page: 页码，默认 1
//   - size: 每页数量，默认 20
//
// 返回:
//   - users: 用户列表
//   - total: 总数量
func (this *User) List() {
    // 实现代码
}

// CreateUser 创建新用户
// 如果邮箱已存在，返回错误
func CreateUser(name, email string) (*User, error) {
    // 实现代码
    return nil, nil
}
```

#### 错误处理

```go
// 1. 错误定义
var (
    ErrUserNotFound    = errors.New("用户不存在")
    ErrInvalidEmail    = errors.New("邮箱格式错误")
    ErrDuplicateEmail  = errors.New("邮箱已被使用")
)

// 2. 错误包装
func GetUser(id int64) (*User, error) {
    user, err := db.Query(id)
    if err != nil {
        return nil, fmt.Errorf("查询用户失败: %w", err)
    }
    return user, nil
}

// 3. 错误处理
user, err := GetUser(id)
if err != nil {
    if errors.Is(err, ErrUserNotFound) {
        // 处理用户不存在
        return
    }
    // 其他错误
    logger.Error(err)
    return
}

// 4. panic 恢复
func SafeExecute(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic recovered: %v", r)
        }
    }()
    
    fn()
    return nil
}
```

### 安全建议

#### SQL 注入防护

```go
// ❌ 错误示例 - SQL 注入风险
sql := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)
db.QuerySQL(sql)

// ✅ 正确示例 - 使用参数化查询
sql := "SELECT * FROM users WHERE name = ?"
db.QuerySQL(sql, name)

// ✅ 使用查询构建器
ar := db.AR()
ar.Select("*").From("users").Where(map[string]interface{}{
    "name": name,
})
db.Query(ar)
```

#### XSS 防护

```go
// 模板自动转义
// {{.Content}}  - 自动转义
// {{.Content | html}}  - 不转义（谨慎使用）

// 在控制器中处理
import "html"

func (this *User) Display() {
    content := this.Ctx.POST("content")
    
    // 转义 HTML
    safeContent := html.EscapeString(content)
    
    this.View.Set("content", safeContent)
    this.View.Render("display")
}
```

#### CSRF 防护

```go
// CSRF 中间件
func CSRFMiddleware(ctx gcore.Ctx) bool {
    // GET 请求生成 token
    if ctx.IsGET() {
        token := generateCSRFToken()
        ctx.SessionStart()
        ctx.Session.Set("csrf_token", token)
        ctx.Set("csrf_token", token)
        return false
    }
    
    // POST 请求验证 token
    if ctx.IsPOST() {
        formToken := ctx.POST("csrf_token")
        ctx.SessionStart()
        sessionToken := ctx.Session.Get("csrf_token")
        
        if formToken != sessionToken {
            ctx.WriteHeader(403)
            ctx.Write("CSRF token mismatch")
            return true
        }
    }
    
    return false
}

// 模板中使用
// <input type="hidden" name="csrf_token" value="{{.csrf_token}}">
```

#### 密码安全

```go
import "golang.org/x/crypto/bcrypt"

// 密码哈希
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// 密码验证
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// 使用示例
func (this *User) Register() {
    password := this.Ctx.POST("password")
    
    // 哈希密码
    hashedPassword, err := HashPassword(password)
    if err != nil {
        this.Ctx.JSON(500, map[string]string{"error": "密码处理失败"})
        return
    }
    
    // 保存到数据库
    SaveUser(username, hashedPassword)
}
```

#### 敏感信息保护

```go
// 1. 配置文件不提交敏感信息
// .gitignore
conf/prod.toml
conf/.env

// 2. 使用环境变量
dbPassword := os.Getenv("DB_PASSWORD")
apiKey := os.Getenv("API_KEY")

// 3. 日志脱敏
func MaskSensitive(data string) string {
    if len(data) <= 4 {
        return "****"
    }
    return data[:2] + "****" + data[len(data)-2:]
}

logger.Infof("User email: %s", MaskSensitive(email))
// 输出: User email: zh****om
```

### 性能优化

#### 数据库优化

```go
// 1. 使用索引
CREATE INDEX idx_user_email ON users(email);
CREATE INDEX idx_post_user_id ON posts(user_id);

// 2. 批量操作
func BatchCreateUsers(users []User) error {
    ar := db.AR()
    data := make([]map[string]interface{}, len(users))
    for i, u := range users {
        data[i] = map[string]interface{}{
            "name":  u.Name,
            "email": u.Email,
        }
    }
    ar.InsertBatch("users", data)
    _, err := db.Exec(ar)
    return err
}

// 3. 分页查询
func GetUsers(page, pageSize int) ([]User, error) {
    offset := (page - 1) * pageSize
    ar := db.AR()
    ar.Select("*").
        From("users").
        Limit(pageSize, offset).
        OrderBy("created_at", "DESC")
    
    result, err := db.Query(ar)
    if err != nil {
        return nil, err
    }
    
    // 转换为结构体
    users, _ := result.Structs(&User{})
    return users, nil
}
```

#### 缓存策略

```go
// 缓存模式
func GetUserWithCache(id int64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    // 1. 查缓存
    cached, err := cache.Get(cacheKey)
    if err == nil && cached != "" {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }
    
    // 2. 查数据库
    user, err := GetUserFromDB(id)
    if err != nil {
        return nil, err
    }
    
    // 3. 写缓存
    data, _ := json.Marshal(user)
    cache.Set(cacheKey, string(data), 300*time.Second)
    
    return user, nil
}

// 缓存更新
func UpdateUser(id int64, data map[string]interface{}) error {
    // 1. 更新数据库
    err := UpdateUserInDB(id, data)
    if err != nil {
        return err
    }
    
    // 2. 删除缓存
    cacheKey := fmt.Sprintf("user:%d", id)
    cache.Del(cacheKey)
    
    return nil
}
```

### 代码质量

#### 单元测试覆盖

```bash
# 运行测试
go test ./...

# 查看覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 设置覆盖率目标
# 核心业务逻辑: 80%+
# 工具函数: 90%+
# API 接口: 70%+
```

#### 代码审查

```go
// 审查检查项
// 1. 命名是否清晰
// 2. 注释是否完整
// 3. 错误处理是否完善
// 4. 是否有安全风险
// 5. 性能是否有问题
// 6. 测试是否充分
// 7. 代码是否可读
// 8. 是否遵循规范
```

#### 静态分析

```bash
# 安装工具
go install golang.org/x/lint/golint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest

# 运行检查
golint ./...
go vet ./...
staticcheck ./...
gosec ./...

# 集成到 CI
make lint
make security-check
```

### 最佳实践总结

1. **架构设计**
   - 分层清晰，职责单一
   - 依赖注入，降低耦合
   - 面向接口编程

2. **代码质量**
   - 遵循命名规范
   - 完善注释文档
   - 编写单元测试
   - 定期代码审查

3. **安全防护**
   - 参数化查询防注入
   - 输入验证和过滤
   - 输出转义防 XSS
   - CSRF 令牌保护
   - 密码加密存储

4. **性能优化**
   - 合理使用缓存
   - 数据库索引优化
   - 批量操作减少 IO
   - 连接池管理
   - 异步处理

5. **可维护性**
   - 统一错误处理
   - 完整的日志记录
   - 配置环境分离
   - 版本控制规范
   - 文档及时更新

---

## 常见问题

### 安装和配置

**Q: 如何安装 GMC 框架？**

A: 使用 go get 命令安装：
```bash
go get -u github.com/snail007/gmc
```

**Q: GMCT 工具安装失败怎么办？**

A: 确保 Go 版本 1.16+，设置 GOPROXY：
```bash
export GOPROXY=https://goproxy.cn,direct
go install github.com/snail007/gmct/cmd/gmct@latest
```

**Q: 如何更改默认端口？**

A: 修改 `conf/app.toml`：
```toml
[httpserver]
listen = ":8080"  # 改为你需要的端口
```

### 路由问题

**Q: 控制器方法无法访问？**

A: 检查以下几点：
1. 方法名首字母大写
2. 方法不以 `_` 或 `__` 结尾
3. 路由已正确注册
4. URL 路径使用小写

**Q: 如何实现 RESTful 路由？**

A: 使用 HTTP 方法绑定：
```go
r.GET("/users", ListUsers)
r.POST("/users", CreateUser)
r.GET("/users/:id", GetUser)
r.PUT("/users/:id", UpdateUser)
r.DELETE("/users/:id", DeleteUser)
```

### 数据库问题

**Q: 数据库连接失败？**

A: 检查配置和驱动：
```bash
# 确保导入了数据库驱动
import _ "github.com/go-sql-driver/mysql"

# 检查 DSN 配置
dsn = "user:password@tcp(host:port)/dbname?charset=utf8mb4"
```

**Q: 如何使用事务？**

A: 参考事务示例：
```go
tx, err := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 执行操作
_, err = db.ExecTx(ar, tx)
if err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

### Session 问题

**Q: Session 数据丢失？**

A: 检查：
1. 是否调用了 `SessionStart()`
2. Cookie 是否被正确设置
3. Session 存储配置是否正确
4. Session 是否过期

**Q: 多服务器 Session 共享？**

A: 使用 Redis 存储：
```toml
[session]
store = "redis"
[session.redis]
address = "redis-server:6379"
```

### 模板问题

**Q: 模板文件找不到？**

A: 检查：
1. 模板文件路径是否正确
2. 文件扩展名是否匹配配置
3. 是否使用了正确的分隔符

**Q: 模板变量显示 `<no value>`？**

A: 使用 `val` 函数安全输出：
```html
{{val . "name"}}
```

### 部署问题

**Q: 编译后体积太大？**

A: 使用编译优化：
```bash
go build -ldflags="-s -w" -o myapp

# 进一步压缩
upx --best --lzma myapp
```

**Q: 如何实现平滑重启？**

A: 使用信号：
```bash
# Linux 平台
kill -USR2 <pid>

# 或
pkill -USR2 myapp
```

**Q: Docker 容器中文乱码？**

A: 设置环境变量：
```dockerfile
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8
```

### 性能问题

**Q: 如何提升性能？**

A: 性能优化建议：
1. 启用缓存减少数据库查询
2. 使用连接池管理数据库连接
3. 异步处理耗时操作
4. 静态文件使用 CDN
5. 启用 gzip 压缩
6. 数据库查询优化和索引

**Q: 内存占用过高？**

A: 检查：
1. 是否有内存泄漏
2. 连接池配置是否合理
3. 缓存是否设置了过期时间
4. 使用 pprof 分析内存使用

### 开发问题

**Q: 热编译不生效？**

A: 检查：
1. `grun.toml` 配置是否正确
2. 文件扩展名是否在监控列表
3. 是否在项目根目录执行
4. 是否有编译错误

**Q: 如何调试代码？**

A: 调试方法：
1. 使用日志输出
2. 使用 Delve 调试器
3. 使用 IDE 断点调试
4. 使用 pprof 性能分析

### API 开发

**Q: CORS 跨域问题？**

A: 添加 CORS 中间件：
```go
func CORSMiddleware(ctx gcore.Ctx) bool {
    ctx.SetHeader("Access-Control-Allow-Origin", "*")
    ctx.SetHeader("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
    ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type,Authorization")
    
    if ctx.Request().Method == "OPTIONS" {
        ctx.WriteHeader(204)
        return true
    }
    return false
}
```

**Q: 如何实现 JWT 认证？**

A: 参考 JWT 中间件示例（见 API 开发章节）

### 其他问题

**Q: 如何贡献代码？**

A: 
1. Fork 项目
2. 创建特性分支
3. 提交代码
4. 发起 Pull Request
5. 等待审核

**Q: 在哪里获取帮助？**

A:
- GitHub Issues: https://github.com/snail007/gmc/issues
- 文档: https://github.com/snail007/gmc/tree/master/docs
- 示例: https://github.com/snail007/gmc/tree/master/demos

---

## 总结

### GMC 框架特点

GMC 是一个功能全面、性能卓越的 Go Web 框架：

**核心优势：**

1. **简单易用** - 直观的 API 设计，快速上手
2. **高性能** - 优化的路由引擎，高效的请求处理
3. **功能完整** - Web、API、模板、数据库、缓存等开箱即用
4. **强大工具链** - GMCT 提供完整的开发工具
5. **灵活扩展** - Provider 模式支持任意组件替换
6. **生产就绪** - 经过生产环境验证

**适用场景：**

- Web 网站开发
- RESTful API 服务
- 微服务架构
- 管理后台系统
- 快速原型开发

### 学习路径

**初级（1-2周）：**
1. 安装 GMC 和 GMCT
2. 创建第一个 Web 项目
3. 学习路由和控制器
4. 掌握模板系统
5. 了解配置管理

**中级（2-4周）：**
1. 数据库操作和 ORM
2. Session 和 Cookie
3. 缓存使用
4. 中间件开发
5. API 开发
6. 文件上传处理

**高级（1-2月）：**
1. 自定义 Provider
2. 性能优化
3. 安全防护
4. 测试和部署
5. 微服务架构
6. 高并发处理

### 开发建议

**项目开始前：**
- 规划好项目结构
- 确定技术栈和依赖
- 设计数据库表结构
- 定义 API 接口规范

**开发过程中：**
- 遵循代码规范
- 编写单元测试
- 添加必要注释
- 定期代码审查
- 使用版本控制

**部署上线前：**
- 完善错误处理
- 添加日志记录
- 性能压力测试
- 安全漏洞扫描
- 准备回滚方案

### 社区资源

- **官方仓库**: https://github.com/snail007/gmc
- **在线文档**: https://pkg.go.dev/github.com/snail007/gmc
- **示例代码**: https://github.com/snail007/gmc/tree/master/demos
- **问题反馈**: https://github.com/snail007/gmc/issues

### 持续学习

- 关注 GMC 更新日志
- 阅读源码了解实现
- 参与社区讨论
- 贡献代码和文档
- 分享使用经验

### 致谢

感谢所有为 GMC 框架做出贡献的开发者！

---

**GMC - 让 Go Web 开发更简单！**

*最后更新时间: 2024-01*

*文档版本: v1.0*

---

## 常用工具包

GMC 提供了丰富的工具包，涵盖常见的开发需求。

### GPool - 协程池

高性能、并发安全的 Go 协程池，支持动态伸缩、空闲超时、panic 恢复等特性。

**核心特性：**
- ✅ 动态工作协程管理，运行时增减
- ✅ 支持空闲超时自动回收
- ✅ 自动 panic 恢复和自定义处理
- ✅ 提供 OptimizedPool（推荐）和 BasicPool 两种实现
- ✅ 性能提升 30-50%（OptimizedPool）

**快速使用：**

```go
import "github.com/snail007/gmc/util/gpool"

// 推荐：使用 OptimizedPool
pool := gpool.NewOptimized(10)
defer pool.Stop()

pool.Submit(func() {
    // 执行任务
})
pool.WaitDone()
```

**详细文档：** [util/gpool/README.md](https://github.com/snail007/gmc/blob/master/util/gpool/README.md)

---

### Rate - 限流器

提供滑动窗口和令牌桶两种限流算法，支持 API 限流、带宽控制等场景。

**核心特性：**
- ✅ 滑动窗口限流器 - 严格控制时间窗口内的请求数
- ✅ 令牌桶限流器 - 平滑限流，支持突发流量
- ✅ 并发安全，高性能
- ✅ 适用于 API 限流、带宽控制、防刷等场景

**快速使用：**

```go
import "github.com/snail007/gmc/util/rate"

// 滑动窗口：每秒最多 100 个请求
limiter := grate.NewSlidingWindowLimiter(100, time.Second)
if limiter.Allow() {
    // 处理请求
}

// 令牌桶：每秒 10 个令牌，支持突发 20 个
burstLimiter := grate.NewTokenBucketBurstLimiter(10, time.Second, 20)
```

**详细文档：** [util/rate/README.md](https://github.com/snail007/gmc/blob/master/util/rate/README.md)

---

### Captcha - 验证码

纯 Go 实现的验证码生成器，不依赖第三方图形库。

**核心特性：**
- ✅ 使用简单，无需额外依赖
- ✅ 支持多字体、多颜色
- ✅ 可自定义大小、干扰强度
- ✅ 支持数字、字母、混合模式

**快速使用：**

```go
import "github.com/snail007/gmc/util/captcha"

cap := gcaptcha.New()
cap.SetFont("comic.ttf")
cap.SetSize(128, 64)

// 生成 4 位数字验证码
img, str := cap.Create(4, gcaptcha.NUM)

// 自定义验证码内容
img := cap.CreateCustom("hello")
```

**详细文档：** [util/captcha/README.md](https://github.com/snail007/gmc/blob/master/util/captcha/README.md)

---

### 完整工具包列表

GMC 还提供了更多实用工具包，每个包都有详细的 README 文档：

#### 字符串和编码
- **[strings](https://github.com/snail007/gmc/blob/master/util/strings/README.md)** - 字符串处理工具
- **[bytes](https://github.com/snail007/gmc/blob/master/util/bytes/README.md)** - 字节处理工具
- **[cast](https://github.com/snail007/gmc/blob/master/util/cast/README.md)** - 类型转换工具
- **[hash](https://github.com/snail007/gmc/blob/master/util/hash/README.md)** - 哈希和加密工具
- **[compress](https://github.com/snail007/gmc/blob/master/util/compress/README.md)** - 压缩和解压工具

#### 数据结构
- **[collection](https://github.com/snail007/gmc/blob/master/util/collection/README.md)** - 集合操作工具
- **[value](https://github.com/snail007/gmc/blob/master/util/value/README.md)** - 值操作工具

#### 网络和 HTTP
- **[net](https://github.com/snail007/gmc/blob/master/util/net/README.md)** - 网络工具
- **[proxy](https://github.com/snail007/gmc/blob/master/util/proxy/README.md)** - 代理工具
- **[url](https://github.com/snail007/gmc/blob/master/util/url/README.md)** - URL 处理工具

#### 文件和系统
- **[file](https://github.com/snail007/gmc/blob/master/util/file/README.md)** - 文件操作工具
- **[env](https://github.com/snail007/gmc/blob/master/util/env/README.md)** - 环境变量工具

#### 开发和调试
- **[pprof](https://github.com/snail007/gmc/blob/master/util/pprof/README.md)** - 性能分析工具
- **[testing](https://github.com/snail007/gmc/blob/master/util/testing/README.md)** - 测试工具
- **[reflect](https://github.com/snail007/gmc/blob/master/util/reflect/README.md)** - 反射工具

#### 其他工具
- **[cond](https://github.com/snail007/gmc/blob/master/util/cond/README.md)** - 条件判断工具

> **提示：** 点击工具包名称查看详细的 API 文档、使用示例和最佳实践。

### 贡献

如果你开发了有用的工具包，欢迎提交 Pull Request 贡献到 GMC 项目！

---

## 附录

### A. 常用命令

#### GMCT 命令

```bash
# 版本信息
gmct version

# 创建项目
gmct new web --pkg myapp
gmct new api --pkg myapi
gmct new api-simple --pkg simpleapi

# 代码生成
gmct controller -n User
gmct model -n user
gmct model -n user -t sqlite3

# 热编译
gmct run
```

> **注意：** 不再推荐使用 `gmct tpl`、`gmct static`、`gmct i18n` 打包命令。  
> 请使用 Go 原生的 `embed` 功能，详见 [资源嵌入](#资源嵌入) 章节。

#### Go 命令

```bash
# 编译
go build -o myapp
go build -ldflags="-s -w" -o myapp

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o myapp-linux
GOOS=windows GOARCH=amd64 go build -o myapp.exe

# 测试
go test ./...
go test -v ./...
go test -cover ./...
go test -bench=. ./...

# 依赖管理
go mod init
go mod tidy
go mod download
go mod vendor
```

### B. 配置参数

#### HTTP 服务器

```toml
[httpserver]
listen = ":8080"              # 监听地址
tlsenable = false             # 启用 HTTPS
tlscert = ""                  # 证书文件
tlskey = ""                   # 密钥文件
readtimeout = 60              # 读超时（秒）
writetimeout = 60             # 写超时（秒）
idletimeout = 60              # 空闲超时（秒）
maxheaderbytes = 1048576      # 最大请求头（字节）
```

#### 数据库

```toml
[database]
enable = true
driver = "mysql"              # mysql, postgres, sqlite3
dsn = ""                      # 连接字符串
maxidle = 10                  # 最大空闲连接
maxconns = 100                # 最大连接数
maxlifetimeseconds = 3600     # 连接最大生命周期
timeout = 5000                # 连接超时（毫秒）
```

#### 缓存

```toml
[cache]
enable = true

[[cache.stores]]
store = "memory"              # memory, redis, file
cleanupintervalseconds = 60   # 清理间隔
```

#### Session

```toml
[session]
enable = true
store = "memory"              # memory, redis, file
ttl = 3600                    # 有效期（秒）
cookiename = "session_id"     # Cookie 名称
cookiedomain = ""             # Cookie 域名
cookiepath = "/"              # Cookie 路径
cookiesecure = false          # 仅 HTTPS
cookiehttponly = true         # HttpOnly
```

#### 日志

```toml
[log]
level = "info"                # trace, debug, info, warn, error, fatal
output = "console"            # console, file, both
async = false                 # 异步写入
filename = "logs/app.log"     # 日志文件
maxsize = 100                 # 文件大小（MB）
maxbackups = 10               # 备份数量
maxage = 30                   # 保留天数
compress = true               # 压缩归档
```

### C. 错误码

```go
// HTTP 状态码
200 OK                        // 成功
201 Created                   // 创建成功
204 No Content                // 无内容
400 Bad Request               // 请求错误
401 Unauthorized              // 未授权
403 Forbidden                 // 禁止访问
404 Not Found                 // 未找到
405 Method Not Allowed        // 方法不允许
429 Too Many Requests         // 请求过多
500 Internal Server Error     // 服务器错误
502 Bad Gateway               // 网关错误
503 Service Unavailable       // 服务不可用
```

### D. 性能指标

#### 基准测试结果

```
# 测试环境
CPU: Intel Core i7-9700K @ 3.60GHz
RAM: 32GB DDR4
OS: Ubuntu 20.04 LTS
Go: 1.21.0

# 测试结果（Hello World）
Requests/sec: 150,000+
Latency (avg): < 1ms
Latency (p99): < 5ms

# 数据库查询
Simple Query: 10,000 req/s
Complex Query: 5,000 req/s
Transaction: 3,000 req/s
```

### E. 版本兼容性

```
GMC Framework: v1.0.0+
Go Version: 1.16+
MySQL: 5.5+, 8.0+
Redis: 5.0+, 6.0+, 7.0+
SQLite: 3.30+
```

### F. 第三方库推荐

```go
// 数据验证
"github.com/go-playground/validator/v10"

// JWT
"github.com/dgrijalva/jwt-go"

// 加密
"golang.org/x/crypto/bcrypt"

// UUID
"github.com/google/uuid"

// 时间处理
"github.com/jinzhu/now"

// HTTP 客户端
"github.com/go-resty/resty/v2"

// 配置
"github.com/spf13/viper"

// 日志
"go.uber.org/zap"
"github.com/sirupsen/logrus"
```

---

**感谢您使用 GMC 框架！**

如有任何问题或建议，欢迎提交 Issue 或 Pull Request。

Happy Coding! 🚀
