# GMC 框架架构说明

本文档说明 GMC 框架的包组织结构、依赖关系以及如何正确引用框架的各个模块。

## 包组织结构

GMC 框架采用模块化设计，主要包含以下几个部分：

```
gmc/
├── core/               # 核心接口定义
├── using/              # 提供者注册（依赖注入）
│   ├── basic/         # 基础组件（app, config, error, log, ctx）
│   ├── cache/         # 缓存组件
│   ├── db/            # 数据库组件
│   └── web/           # Web 组件（session, template, router, server）
├── module/            # 功能模块实现
│   ├── app/           # 应用程序框架
│   ├── cache/         # 缓存实现（Memory, Redis, File）
│   ├── config/        # 配置管理
│   ├── ctx/           # 上下文
│   ├── db/            # 数据库（MySQL, SQLite3）
│   ├── error/         # 错误处理
│   ├── i18n/          # 国际化
│   ├── log/           # 日志
│   └── middleware/    # 中间件
│       └── accesslog/ # 访问日志中间件
├── http/              # HTTP 相关
│   ├── controller/    # 控制器
│   ├── cookie/        # Cookie 处理
│   ├── router/        # 路由
│   ├── server/        # HTTP 服务器
│   ├── session/       # 会话管理
│   ├── template/      # 模板引擎
│   └── view/          # 视图
└── util/              # 工具包（独立使用）
    ├── args/          # 命令行参数解析
    ├── batch/         # 批处理
    ├── bytes/         # 字节处理
    ├── captcha/       # 验证码
    ├── cast/          # 类型转换
    ├── collection/    # 集合操作
    ├── compress/      # 压缩/解压
    ├── cond/          # 条件变量
    ├── env/           # 环境变量
    ├── exec/          # 命令执行
    ├── file/          # 文件操作
    ├── func/          # 函数工具
    ├── gpool/         # 协程池
    ├── hash/          # 哈希
    ├── http/          # HTTP 工具
    ├── json/          # JSON 工具
    ├── linklist/      # 链表
    ├── list/          # 列表
    ├── loop/          # 循环控制
    ├── map/           # Map 工具
    ├── net/           # 网络工具
    ├── os/            # 操作系统工具
    ├── paginator/     # 分页
    ├── pprof/         # 性能分析
    ├── process/       # 进程管理
    ├── proxy/         # 代理
    ├── rand/          # 随机数
    ├── rate/          # 限流
    ├── reflect/       # 反射工具
    ├── set/           # 集合
    ├── strings/       # 字符串工具
    ├── sync/          # 同步工具
    ├── testing/       # 测试工具
    ├── url/           # URL 工具
    └── value/         # 值处理
```

## Using 包说明

`using` 包实现了依赖注入机制，负责将各个模块的具体实现注册到核心接口。

### using/basic

提供基础组件的提供者注册：
- App（应用程序框架）
- Config（配置管理）
- Error（错误处理）
- Logger（日志）
- Ctx（上下文）

### using/cache

提供缓存组件的提供者注册，依赖 `using/basic`。

### using/db

提供数据库组件的提供者注册，依赖 `using/basic`。

### using/web

提供 Web 组件的提供者注册，依赖 `using/basic`，包括：
- Session（会话管理）
- SessionStorage（会话存储）
- View（视图）
- Template（模板引擎）
- HTTPRouter（路由）
- Cookies（Cookie 处理）
- I18n（国际化）
- Controller（控制器）
- HTTPServer（HTTP 服务器）
- APIServer（API 服务器）

## 如何引用 GMC 包

根据使用场景，有不同的引用方式：

### 场景 1：完整的 Web 应用

**推荐方式：** 直接引用 `github.com/snail007/gmc`

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

func main() {
    // 创建应用
    app := gmc.New.App()
    
    // 使用配置
    cfg := app.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    // 初始化数据库
    gmc.DB.Init(cfg)
    
    // 初始化缓存
    gmc.Cache.Init(cfg)
    
    // 创建 Web 服务器
    s := gmc.New.HTTPServer(app.Ctx())
    s.Init(cfg)
    
    // 配置路由
    r := s.Router()
    r.Controller("/user", new(UserController))
    
    // 运行应用
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

**说明：** `github.com/snail007/gmc` 已经包含了 `using/cache`、`using/db` 和 `using/web`，所以可以直接使用 `gmc.DB`、`gmc.Cache` 等。

### 场景 2：仅使用基础模块（不需要 Web 功能）

如果只需要使用配置、日志、错误处理等基础功能，可以直接引用相应的模块包：

```go
package main

import (
    "github.com/snail007/gmc/module/config"
    "github.com/snail007/gmc/module/log"
    "github.com/snail007/gmc/module/error"
)

func main() {
    // 配置
    cfg := gconfig.New()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    // 日志
    logger := glog.New()
    logger.Info("Application started")
    
    // 错误处理
    err := gerror.New()
    err.Wrap(someError, "Failed to process")
}
```

### 场景 3：仅使用工具包

工具包（`util/` 目录下的包）都是独立的，可以单独引用：

```go
package main

import (
    "github.com/snail007/gmc/util/gpool"
    "github.com/snail007/gmc/util/cast"
    "github.com/snail007/gmc/util/captcha"
)

func main() {
    // 协程池
    pool := gpool.New(10)
    pool.Submit(func() {
        // 任务
    })
    
    // 类型转换
    str := gcast.ToString(123)
    
    // 验证码
    cap := gcaptcha.NewDefault()
    img, code := cap.Create(4, gcaptcha.NUM)
}
```

### 场景 4：只需要数据库功能

```go
package main

import (
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    // 方式 1：使用 gmc 包（推荐）
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    gmc.DB.Init(cfg)
    db := gmc.DB.DB()
    
    // 方式 2：仅使用 using/db
    // import _ "github.com/snail007/gmc/using/db"
    // 然后使用 gdb 包的功能
}
```

### 场景 5：只需要缓存功能

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    gmc.Cache.Init(cfg)
    cache := gmc.Cache.Cache()
    
    cache.Set("key", "value", 3600)
}
```

## 依赖关系图

```
github.com/snail007/gmc
  └── using/web (自动导入)
       └── using/basic
  └── using/db (自动导入)
       └── using/basic
  └── using/cache (自动导入)
       └── using/basic

util/* (独立，无依赖)

module/config (基础模块，可独立使用)
module/log (基础模块，可独立使用)
module/error (基础模块，可独立使用)

module/cache (需要 using/cache 或 gmc)
module/db (需要 using/db 或 gmc)
module/i18n (需要 using/web 或 gmc)

http/server (需要 using/web 或 gmc)
http/router (需要 using/web 或 gmc)
http/session (需要 using/web 或 gmc)
http/template (需要 using/web 或 gmc)
http/controller (需要 using/web 或 gmc)
http/cookie (可独立使用或配合 gmc)
```

## 最佳实践

### 1. 推荐的项目结构

```
myapp/
├── main.go              # 入口文件，引用 github.com/snail007/gmc
├── app.toml            # 配置文件
├── controller/         # 控制器
│   ├── user.go
│   └── product.go
├── model/              # 数据模型
│   ├── user.go
│   └── product.go
├── router/             # 路由配置
│   └── router.go
├── initialize/         # 初始化
│   └── init.go
├── service/            # 业务逻辑
│   ├── user_service.go
│   └── product_service.go
├── middleware/         # 自定义中间件
│   └── auth.go
├── views/              # 模板文件
│   ├── layout.html
│   └── user/
│       ├── list.html
│       └── detail.html
└── static/             # 静态文件
    ├── css/
    ├── js/
    └── images/
```

### 2. main.go 典型结构

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
    "myapp/initialize"
    "myapp/router"
)

func main() {
    // 创建应用
    app := gmc.New.App()
    
    // 加载配置
    cfg := app.Config()
    cfg.SetConfigFile("app.toml")
    if err := cfg.ReadInConfig(); err != nil {
        panic(err)
    }
    
    // 创建 HTTP 服务器
    s := gmc.New.HTTPServer(app.Ctx())
    if err := s.Init(cfg); err != nil {
        panic(err)
    }
    
    // 添加服务
    app.AddService(gmc.ServiceItem{
        Service: s,
        AfterInit: func(s gmc.Service, cfg gcore.Config) (err error) {
            // 初始化数据库、缓存等
            if err = initialize.Initialize(s.(*gmc.HTTPServer)); err != nil {
                return
            }
            // 初始化路由
            router.InitRouter(s.(*gmc.HTTPServer))
            return
        },
    })
    
    // 运行应用
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

### 3. 不要混用引用方式

❌ **不推荐：**
```go
import (
    "github.com/snail007/gmc"
    "github.com/snail007/gmc/module/db"
    _ "github.com/snail007/gmc/using/db"  // 不需要，gmc 已包含
)
```

✅ **推荐：**
```go
import (
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"  // 仅作为类型引用
)
```

## 常见问题

### Q1: 为什么需要 using 包？

A: `using` 包实现了依赖注入，将具体的实现注册到核心接口。这样做的好处是：
- 模块解耦，核心接口不依赖具体实现
- 可以方便地替换实现（如替换缓存后端）
- 支持按需引入模块

### Q2: 什么时候需要显式引用 using 包？

A: 几乎不需要。直接使用 `github.com/snail007/gmc` 即可，它已经包含了所有必要的 using 包。

### Q3: util 包可以单独使用吗？

A: 可以。`util/` 目录下的所有包都是独立的工具包，可以单独引用和使用，不需要任何其他依赖。

### Q4: 如何选择引用方式？

A: 
- **完整 Web 应用** → `github.com/snail007/gmc`
- **仅基础功能** → 直接引用 `module/config`、`module/log` 等
- **仅工具函数** → 直接引用 `util/*` 包

### Q5: 为什么有些示例使用 module 包，有些使用 gmc 包？

A: 
- 使用 `gmc.New.XXX()` 是推荐方式，适用于完整应用
- 直接使用 `module/xxx` 包适用于只需要该模块功能的场景

## 总结

- **完整 Web 应用**：使用 `github.com/snail007/gmc`，一站式解决方案
- **工具包**：使用 `github.com/snail007/gmc/util/*`，独立无依赖
- **基础模块**：可以单独使用 `module/config`、`module/log` 等
- **不要**：手动引用 `using` 包，除非你完全理解依赖注入机制

选择合适的引用方式，可以让代码更清晰、依赖更简单。
