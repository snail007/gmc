<div align="center">

# GMC

<img src="/doc/images/logo2.png" width="200" alt="GMC Logo"/>

### 🚀 现代化的 Go Web & API 开发框架

一个智能、灵活、高性能的 Golang Web 和 API 开发框架

[![Actions Status](https://github.com/snail007/gmc/workflows/build/badge.svg)](https://github.com/snail007/gmc/actions)
[![codecov](https://codecov.io/gh/snail007/gmc/branch/master/graph/badge.svg)](https://codecov.io/gh/snail007/gmc)
[![Go Report](https://goreportcard.com/badge/github.com/snail007/gmc)](https://goreportcard.com/report/github.com/snail007/gmc)
[![API Reference](https://img.shields.io/badge/go.dev-reference-blue)](https://pkg.go.dev/github.com/snail007/gmc)
[![LICENSE](https://img.shields.io/github/license/snail007/gmc)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/snail007/gmc)](go.mod)

[English](README_EN.md) | 简体中文

[📖 完整文档](https://snail007.github.io/gmc/zh/) | [🎯 快速开始](#-%E5%BF%AB%E9%80%9F%E5%BC%80%E5%A7%8B) | [💡 特性](#-%E6%A0%B8%E5%BF%83%E7%89%B9%E6%80%A7) | [🔧 示例](#-%E7%A4%BA%E4%BE%8B%E4%BB%A3%E7%A0%81)

</div>

---

## 📋 目录

- [简介](#-%E7%AE%80%E4%BB%8B)
- [核心特性](#-%E6%A0%B8%E5%BF%83%E7%89%B9%E6%80%A7)
- [安装](#-%E5%AE%89%E8%A3%85)
- [快速开始](#-%E5%BF%AB%E9%80%9F%E5%BC%80%E5%A7%8B)
- [架构设计](#%EF%B8%8F-%E6%9E%B6%E6%9E%84%E8%AE%BE%E8%AE%A1)
- [示例代码](#-%E7%A4%BA%E4%BE%8B%E4%BB%A3%E7%A0%81)
- [核心组件](#-%E6%A0%B8%E5%BF%83%E7%BB%84%E4%BB%B6)
- [工具包](#%EF%B8%8F-%E5%B7%A5%E5%85%B7%E5%8C%85)
- [配置说明](#%EF%B8%8F-%E9%85%8D%E7%BD%AE%E8%AF%B4%E6%98%8E)
- [性能测试](#-%E6%80%A7%E8%83%BD%E6%B5%8B%E8%AF%95)
- [项目结构](#-%E9%A1%B9%E7%9B%AE%E7%BB%93%E6%9E%84)
- [贡献指南](#-%E8%B4%A1%E7%8C%AE%E6%8C%87%E5%8D%97)
- [许可证](#-%E8%AE%B8%E5%8F%AF%E8%AF%81)
- [联系我们](#-%E8%81%94%E7%B3%BB%E6%88%91%E4%BB%AC)

---

## 🎯 简介

**GMC**（Go Micro Container）是一个面向现代 Web 开发的全栈 Golang 框架。它致力于提供：

- 🎨 **高生产力** - 用更少的代码完成更多的功能
- ⚡ **高性能** - 基于高性能路由和优化的中间件
- 🧩 **模块化** - 清晰的架构和完善的依赖注入
- 🛠️ **工具丰富** - 60+ 开箱即用的实用工具包
- 📦 **易于使用** - 简洁的 API 设计和详细的文档

GMC 不仅是一个 Web 框架，更是一个完整的开发工具集，适用于从小型 API 到大型企业级应用的各种场景。

---

## ✨ 核心特性

### 🌐 Web & API 开发
- **RESTful API** - 快速构建 RESTful 风格的 API 服务
- **MVC 架构** - 完整的 MVC 模式支持，清晰的代码组织
- **路由系统** - 高性能路由引擎，支持路由分组、参数、中间件
- **控制器** - 优雅的控制器设计，支持依赖注入
- **模板引擎** - 内置模板引擎，支持布局、继承、自定义函数

### 🗄️ 数据处理
- **多数据库支持** - MySQL、SQLite3 开箱即用
- **ORM 集成** - 优雅的数据库操作接口
- **缓存系统** - Memory、Redis、File 多种缓存后端
- **会话管理** - 灵活的 Session 管理机制

### 🔧 开发工具
- **配置管理** - 支持 TOML、JSON、YAML 等多种配置格式
- **日志系统** - 分级日志、异步写入、自动轮转
- **错误处理** - 完善的错误堆栈和错误链
- **国际化** - i18n 支持，轻松实现多语言
- **验证码** - 内置验证码生成器
- **分页器** - 开箱即用的分页组件

### ⚙️ 高级功能
- **中间件** - 灵活的中间件系统
- **协程池** - 高性能 Goroutine 池管理
- **限流器** - 双算法限流（滑动窗口/令牌桶），支持 API 限流、带宽控制
- **性能分析** - pprof 集成，便捷的性能分析
- **进程管理** - 守护进程、优雅重启支持
- **依赖注入** - 清晰的依赖注入机制
- **热编译** - 开发时自动编译重启（gmct run）
- **资源打包** - 静态文件、模板、i18n 打包进二进制（gmct）

### 🛠️ 实用工具库（60+）
涵盖文件操作、网络工具、加密哈希、类型转换、集合操作、压缩解压、JSON 处理等各个方面。

### 🔨 GMCT 工具链
- **项目生成** - 一键生成 Web/API 项目脚手架
- **热编译** - 开发时自动编译和重启
- **资源打包** - 将静态文件、模板、i18n 打包进二进制
- **项目管理** - 简化开发流程的各种工具

---

## 📦 安装

### 环境要求

- Go 1.16 或更高版本

### 安装框架

```bash
go get -u github.com/snail007/gmc
```

### 安装 GMCT 工具链

**GMCT** 是 GMC 的官方命令行工具，提供项目脚手架、热编译、资源打包等强大功能：

```bash
# 安装 gmct
go install github.com/snail007/gmct@latest

# 验证安装
gmct version
```

#### GMCT 快速安装（Linux/macOS）

```bash
# Linux AMD64
bash -c "$(curl -L https://github.com/snail007/gmct/raw/master/install.sh)" @ linux-amd64

# Linux ARM64
bash -c "$(curl -L https://github.com/snail007/gmct/raw/master/install.sh)" @ linux-arm64

# macOS - 请从 Release 页面下载
# https://github.com/snail007/gmct/releases
```

📖 **GMCT 完整文档**: [https://github.com/snail007/gmct](https://github.com/snail007/gmct)

### 验证安装

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
)

func main() {
    fmt.Println("GMC framework installed successfully!")
}
```

---

## 🚀 快速开始

### 使用 GMCT 创建项目（推荐）

GMCT 是 GMC 的官方工具链，可以快速生成项目脚手架：

```bash
# 创建 Web 项目
mkdir myapp && cd myapp
gmct new web

# 或创建 API 项目
gmct new api

# 热编译模式运行（开发时推荐）
gmct run

# 访问 http://localhost:7080
```

生成的项目结构：
```
myapp/
├── conf/
│   └── app.toml          # 配置文件
├── controller/
│   └── demo.go           # 控制器
├── initialize/
│   └── initialize.go     # 初始化
├── router/
│   └── router.go         # 路由
├── static/               # 静态文件
├── views/                # 模板文件
├── grun.toml            # GMCT 配置
└── main.go              # 入口文件
```

### 手动创建项目

### 1. 创建一个简单的 API 服务

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
    gmap "github.com/snail007/gmc/util/map"
)

func main() {
    // 创建 API 服务器
    api, _ := gmc.New.APIServer(gmc.New.Ctx(), ":8080")

    // 注册路由
    api.API("/", func(c gmc.C) {
        c.Write(gmap.M{
            "code":    0,
            "message": "Hello GMC!",
            "data":    nil,
        })
    })

    // 创建应用并运行
    app := gmc.New.App()
    app.AddService(gcore.ServiceItem{
        Service: api.(gcore.Service),
    })
    
    app.Run()
}
```

运行后访问 `http://localhost:8080/` 即可看到返回的 JSON 数据。

### 2. 创建一个 Web 应用

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type HomeController struct {
    gmc.Controller
}

func (c *HomeController) Index() {
    c.Write("Welcome to GMC!")
}

func main() {
    // 创建应用
    app := gmc.New.App()
    
    // 创建 HTTP 服务器
    s := gmc.New.HTTPServer(app.Ctx())
    s.Router().Controller("/", new(HomeController))
    
    // 添加服务并运行
    app.AddService(gcore.ServiceItem{
        Service: s,
    })
    
    app.Run()
}
```

---

## 🏗️ 架构设计

GMC 采用清晰的模块化架构，主要由以下几部分组成：

```
gmc/
├── core/               # 核心接口定义
├── module/             # 功能模块实现
│   ├── app/           # 应用程序框架
│   ├── cache/         # 缓存（Memory, Redis, File）
│   ├── config/        # 配置管理
│   ├── db/            # 数据库（MySQL, SQLite3）
│   ├── log/           # 日志系统
│   ├── i18n/          # 国际化
│   └── middleware/    # 中间件
├── http/              # HTTP 相关
│   ├── server/        # HTTP/API 服务器
│   ├── router/        # 路由
│   ├── controller/    # 控制器
│   ├── session/       # 会话管理
│   ├── template/      # 模板引擎
│   └── cookie/        # Cookie 处理
├── util/              # 工具包（60+ 独立工具）
│   ├── gpool/         # 协程池
│   ├── captcha/       # 验证码
│   ├── cast/          # 类型转换
│   ├── compress/      # 压缩/解压
│   ├── file/          # 文件操作
│   ├── http/          # HTTP 工具
│   ├── json/          # JSON 工具
│   ├── rate/          # 限流器
│   └── ...            # 更多工具
└── using/             # 依赖注入注册
```

详细架构说明请参考 [ARCHITECTURE.md](ARCHITECTURE.md)

---

## 🔨 GMCT 工具链

GMCT 是 GMC 框架的官方命令行工具，提供项目脚手架、热编译、资源打包等强大功能，极大提升开发效率。

### 🎯 主要功能

#### 1. 项目生成

快速生成标准化的项目结构：

```bash
# 生成 Web 项目（MVC 架构）
gmct new web

# 生成 API 项目（轻量级）
gmct new api

# 指定包名
gmct new web --pkg github.com/yourname/myapp
```

#### 2. 热编译开发

开发时自动监听文件变化，自动编译和重启：

```bash
# 热编译模式运行
gmct run

# 配置文件 grun.toml
[run]
# 监听的文件扩展名
watch_ext = [".go", ".toml"]
# 排除的目录
exclude_dir = ["vendor", ".git"]
# 编译命令
build_cmd = "go build -o tmp/app"
# 运行命令
run_cmd = "./tmp/app"
```

#### 3. 资源打包

将静态文件、模板、i18n 文件打包进二进制，实现单文件部署：

```bash
# 打包模板文件
gmct tpl --dir ./views

# 打包静态文件
gmct static --dir ./static

# 打包 i18n 文件
gmct i18n --dir ./i18n

# 清理打包文件
gmct static --clean
gmct tpl --clean
gmct i18n --clean
```

打包后，你的应用可以编译成单个二进制文件，无需携带任何资源文件。

#### 4. 项目信息

```bash
# 查看版本
gmct version

# 查看帮助
gmct help

# 查看具体命令帮助
gmct new --help
gmct run --help
```

### 📋 GMCT 命令列表

| 命令 | 说明 | 示例 |
|------|------|------|
| `gmct new` | 创建新项目 | `gmct new web` |
| `gmct run` | 热编译运行 | `gmct run` |
| `gmct tpl` | 打包模板 | `gmct tpl --dir ./views` |
| `gmct static` | 打包静态文件 | `gmct static --dir ./static` |
| `gmct i18n` | 打包国际化文件 | `gmct i18n --dir ./i18n` |
| `gmct version` | 查看版本 | `gmct version` |
| `gmct help` | 查看帮助 | `gmct help` |

### 🎬 完整开发流程示例

```bash
# 1. 安装 GMCT
go install github.com/snail007/gmct@latest

# 2. 创建新项目
mkdir mywebapp && cd mywebapp
gmct new web --pkg github.com/me/mywebapp

# 3. 热编译开发
gmct run
# 修改代码后自动重新编译和重启

# 4. 打包资源（生产环境）
gmct static --dir ./static
gmct tpl --dir ./views
gmct i18n --dir ./i18n

# 5. 编译发布
go build -ldflags "-s -w" -o myapp

# 6. 部署
./myapp
# 单个二进制文件，包含所有资源
```

### ⚙️ 配置文件 grun.toml

GMCT 运行配置文件示例：

```toml
[run]
# 监听的文件扩展名
watch_ext = [".go", ".toml", ".html", ".js", ".css"]

# 排除的目录
exclude_dir = [
    "vendor",
    ".git",
    ".idea",
    "tmp",
    "bin",
]

# 编译前执行的命令
before_build = []

# 编译命令
build_cmd = "go build -o tmp/app"

# 运行命令
run_cmd = "./tmp/app"

# 运行后执行的命令
after_run = []

# 延迟重启时间（毫秒）
restart_delay = 1000
```

### 🌟 GMCT 优势

1. **提升开发效率** - 热编译省去手动重启的麻烦
2. **标准化项目** - 统一的项目结构，便于团队协作
3. **简化部署** - 资源打包后单文件部署
4. **降低学习成本** - 开箱即用的最佳实践
5. **灵活配置** - 可自定义编译和运行流程

📖 **完整文档**: [GMCT 工具链仓库](https://github.com/snail007/gmct)

---

## 💡 示例代码

### API 路由

```go
api, _ := gmc.New.APIServer(gmc.New.Ctx(), ":8080")

// GET 请求
api.API("/user/:id", func(c gmc.C) {
    id := c.Param().ByName("id")
    c.Write(gmap.M{
        "user_id": id,
        "name":    "John Doe",
    })
})

// POST 请求
api.API("/user", func(c gmc.C) {
    name := c.Request().FormValue("name")
    // 处理业务逻辑
    c.Write(gmap.M{"status": "created", "name": name})
}, "POST")
```

### 控制器

```go
type UserController struct {
    gmc.Controller
}

func (c *UserController) List() {
    users := []string{"Alice", "Bob", "Charlie"}
    c.Write(users)
}

func (c *UserController) Detail() {
    id := c.Param().ByName("id")
    c.Write("User ID: " + id)
}

// 在路由中注册
router.Controller("/user", new(UserController))
```

### 数据库操作

```go
// 初始化数据库
gmc.DB.Init(cfg)
db := gmc.DB.DB()

// 使用 ActiveRecord 查询（推荐）
ar := db.AR()
ar.From("users").Where(gdb.M{"age >": 18}).OrderBy("created_at", "DESC")
rs, err := db.Query(ar)

// 插入数据
ar = db.AR()
ar.Insert("users", gdb.M{
    "name":  "John",
    "email": "john@example.com",
    "age":   25,
})
result, err := db.Exec(ar)
lastID := result.LastInsertId()

// 更新数据
ar = db.AR()
ar.Update("users", gdb.M{"age": 26}, gdb.M{"id": lastID})
db.Exec(ar)
```

📖 **详细文档**: [数据库模块完整指南](module/db/README.md)

### 缓存使用

```go
// 初始化缓存
gmc.Cache.Init(cfg)
cache := gmc.Cache.Cache()

// 设置缓存（60秒过期）
cache.Set("key", "value", 60)

// 获取缓存
value, exists := cache.Get("key")

// 删除缓存
cache.Del("key")
```

📖 **详细文档**: [缓存模块完整指南](module/cache/README.md)

### 协程池

```go
import "github.com/snail007/gmc/util/gpool"

// 创建协程池（最多10个并发）
pool := gpool.New(10)

// 提交任务
for i := 0; i < 100; i++ {
    pool.Submit(func() {
        // 执行任务
    })
}

// 等待所有任务完成
pool.Wait()
```

📖 **详细文档**: [协程池使用指南](util/gpool/README.md)

### 验证码生成

```go
import "github.com/snail007/gmc/util/captcha"

// 创建验证码
cap := gcaptcha.NewDefault()
img, code := cap.Create(4, gcaptcha.NUM)

// img 是验证码图片数据
// code 是验证码文本
```

📖 **详细文档**: [验证码工具使用指南](util/captcha/README.md)

### 限流器

GMC 提供两种高性能限流器：

**滑动窗口限流器** - 适用于严格控制 QPS，如 API 接口限流：
```go
import "github.com/snail007/gmc/util/rate"

// 创建滑动窗口限流器（每秒最多100个请求）
limiter := grate.NewSlidingWindowLimiter(100, time.Second)

if limiter.Allow() {
    // 处理请求
} else {
    // 请求被限流
}
```

**令牌桶限流器** - 适用于带宽控制、流量整形：
```go
// 创建令牌桶限流器（每秒100个令牌，支持突发200个）
limiter := grate.NewTokenBucketBurstLimiter(100, time.Second, 200)

// 阻塞等待令牌
if err := limiter.Wait(ctx); err != nil {
    return err
}
// 处理请求
```

**使用场景：**
- ✅ **滑动窗口**：API 限流、防刷、严格 QPS 控制
- ✅ **令牌桶**：带宽限制、文件传输、消息队列消费

📖 **详细文档**: [限流器使用指南](util/rate/README.md)

---

### 🔗 更多示例和文档

#### 核心模块
- [应用框架 (App)](module/app/README.md) - 应用程序生命周期管理
- [配置管理 (Config)](module/config/README.md) - 多格式配置文件支持
- [日志系统 (Log)](module/log/README.md) - 强大的日志功能
- [错误处理 (Error)](module/error/README.md) - 错误堆栈和错误链
- [国际化 (i18n)](module/i18n/README.md) - 多语言支持
- [中间件 (Middleware)](module/middleware/README.md) - 中间件系统

#### 工具包（部分）
- [文件操作 (File)](util/file/README.md) - 文件读写、复制、移动等
- [类型转换 (Cast)](util/cast/README.md) - 各种类型之间的转换
- [JSON工具 (JSON)](util/json/README.md) - 高性能 JSON 处理
- [压缩工具 (Compress)](util/compress/README.md) - gzip、tar、zip 等
- [HTTP工具 (HTTP)](util/http/README.md) - HTTP 客户端工具
- [网络工具 (Net)](util/net/README.md) - 网络相关工具函数
- [哈希工具 (Hash)](util/hash/README.md) - MD5、SHA、bcrypt 等
- [字符串工具 (Strings)](util/strings/README.md) - 字符串处理工具
- [集合工具 (Collection)](util/collection/README.md) - 集合操作
- [性能分析 (Pprof)](util/pprof/README.md) - 性能分析工具

**📚 查看所有工具包**: [util/](util/)

**🎓 完整示例**: [demos/](demos/) 目录包含了各种使用场景的完整示例代码

---

## 🧩 核心组件

### HTTP 服务器

GMC 提供两种 HTTP 服务器：

- **HTTPServer** - 完整的 Web 服务器，支持 MVC、模板、Session 等
- **APIServer** - 轻量级 API 服务器，专注于 RESTful API 开发

### 路由系统

- 高性能路由匹配
- 支持路径参数 `/user/:id`
- 支持通配符 `/files/*filepath`
- 路由分组和中间件
- RESTful 路由设计

### 中间件

```go
// 全局中间件
api.Use(func(c gmc.C, next func()) {
    // 前置处理
    start := time.Now()
    
    next() // 调用下一个处理器
    
    // 后置处理
    duration := time.Since(start)
    log.Printf("Request took %v", duration)
})
```

### 模板引擎

```go
// 渲染模板
c.View().Render("user/profile", gmap.M{
    "name": "John",
    "age":  25,
})
```

---

## 🛠️ 工具包

GMC 提供 60+ 独立的工具包，可以在任何 Go 项目中单独使用：

| 分类 | 工具包 | 说明 |
|------|--------|------|
| 🔢 **数据处理** | cast | 类型转换 |
| | json | JSON 操作 |
| | collection | 集合操作 |
| | set | 集合数据结构 |
| | list | 列表操作 |
| | map | Map 工具 |
| 📁 **文件 & I/O** | file | 文件操作 |
| | compress | 压缩/解压（gzip, tar, zip, xz） |
| | bytes | 字节处理 |
| 🌐 **网络** | http | HTTP 客户端工具 |
| | net | 网络工具 |
| | proxy | 代理工具 |
| | url | URL 处理 |
| 🔐 **安全** | hash | 哈希（MD5, SHA, bcrypt） |
| | captcha | 验证码生成 |
| ⚡ **并发** | gpool | 协程池 |
| | sync | 同步工具 |
| | rate | 限流器（滑动窗口/令牌桶） |
| | loop | 循环控制 |
| 🔧 **系统** | process | 进程管理 |
| | exec | 命令执行 |
| | os | 操作系统工具 |
| | env | 环境变量 |
| 📊 **其他** | paginator | 分页器 |
| | pprof | 性能分析 |
| | args | 参数解析 |
| | rand | 随机数 |

单独使用工具包示例：

```go
import "github.com/snail007/gmc/util/cast"

// 类型转换
str := gcast.ToString(123)
num := gcast.ToInt("456")
```

---

## ⚙️ 配置说明

GMC 支持多种配置格式（TOML、JSON、YAML）。推荐使用 TOML 格式。

### 基本配置示例（app.toml）

```toml
# HTTP 服务器配置
[httpserver]
listen = ":8080"
tlsenable = false
tlscert = ""
tlskey = ""

# 模板配置
[template]
dir = "views"
ext = ".html"

# 数据库配置
[database]
default = "mysql"

[database.mysql]
enable = true
driver = "mysql"
dsn = "root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True"
maxidle = 10
maxconns = 100
maxlifetimeseconds = 3600

# 缓存配置
[cache]
default = "redis"

[cache.redis]
enable = true
address = "127.0.0.1:6379"
password = ""
db = 0
timeout = 10

# 日志配置
[log]
level = "info"
output = "stdout"

# Session 配置
[session]
store = "memory"
ttl = 3600
```

### 加载配置

```go
cfg := gmc.New.Config()
cfg.SetConfigFile("app.toml")
err := cfg.ReadInConfig()
```

---

## 📊 性能测试

GMC 在性能测试中表现优异：

```bash
# 运行基准测试
go test -bench=. -benchmem ./...
```

主要性能指标：

- **路由性能** - 高速路由匹配，支持数万路由规模
- **并发处理** - 协程池优化，高效的并发任务调度
- **内存占用** - 优化的内存分配，降低 GC 压力
- **吞吐量** - 高并发下保持稳定的吞吐量

---

## 📂 项目结构

推荐的项目结构：

```
myapp/
├── main.go              # 应用入口
├── app.toml            # 配置文件
├── controller/         # 控制器
│   ├── home.go
│   └── user.go
├── model/              # 数据模型
│   └── user.go
├── service/            # 业务逻辑层
│   └── user_service.go
├── middleware/         # 自定义中间件
│   └── auth.go
├── router/             # 路由配置
│   └── router.go
├── initialize/         # 初始化逻辑
│   └── init.go
├── views/              # 模板文件
│   ├── layout.html
│   └── home/
│       └── index.html
└── static/             # 静态文件
    ├── css/
    ├── js/
    └── images/
```

---

## 🤝 贡献指南

我们欢迎所有形式的贡献！在提交 PR 之前，请确保：

### 代码规范

1. **注释** - 为公共函数和类型添加清晰的注释
2. **测试** - 测试覆盖率应达到 90% 以上
3. **示例** - 为公共函数提供使用示例
4. **基准测试** - 为性能关键代码添加基准测试

### 包必需文件

每个包应包含以下文件（`xxx` 为包名）：

| 文件 | 说明 |
|------|------|
| xxx.go | 主文件 |
| xxx_test.go | 单元测试（覆盖率 > 90%） |
| example_test.go | 示例代码 |
| benchmark_test.go | 基准测试 |
| doc.go | 包说明文档 |
| README.md | 测试和基准测试结果 |

可以参考 `util/gpool` 包来了解详细的代码规范。

### 提交流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📝 文档

- **完整文档**: [https://snail007.github.io/gmc/zh/](https://snail007.github.io/gmc/zh/)
- **API 文档**: [https://pkg.go.dev/github.com/snail007/gmc](https://pkg.go.dev/github.com/snail007/gmc)
- **架构说明**: [ARCHITECTURE.md](ARCHITECTURE.md)
- **示例代码**: [demos/](demos/)

---

## 📜 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

---

## 💬 联系我们

- **GitHub Issues**: [https://github.com/snail007/gmc/issues](https://github.com/snail007/gmc/issues)
- **GitHub Discussions**: [https://github.com/snail007/gmc/discussions](https://github.com/snail007/gmc/discussions)

---

## ⭐ Star 历史

如果这个项目对你有帮助，请给我们一个 Star ⭐

[![Star History Chart](https://api.star-history.com/svg?repos=snail007/gmc&type=Date)](https://star-history.com/#snail007/gmc&Date)

---

## 🙏 致谢

感谢所有为 GMC 做出贡献的开发者！

---

<div align="center">

**[⬆ 回到顶部](#gmc)**

Made with ❤️ by the GMC Team

</div>