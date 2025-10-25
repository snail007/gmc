# GMC Core 核心包

[![Go Reference](https://pkg.go.dev/badge/github.com/snail007/gmc/core.svg)](https://pkg.go.dev/github.com/snail007/gmc/core)

## 概述

`core` 包（gcore）是 GMC 框架的基础抽象层。它定义了所有核心接口和契约，使 GMC 能够实现模块化架构，允许在不改变应用代码的情况下轻松替换和扩展不同的实现。

这个包**只包含接口定义**，不提供具体实现。它作为 GMC 各个组件之间的契约，实现了松耦合和高灵活性。

## 核心特性

- 🎯 **纯接口设计**: 仅包含接口定义，对实现零依赖
- 🔌 **Provider 模式**: 灵活的 provider 系统，用于注册和获取组件实现
- 🏗️ **模块化架构**: 所有框架组件清晰的关注点分离
- 🔄 **可插拔组件**: 无需修改代码即可轻松替换实现
- 📦 **零外部依赖**: 核心接口没有外部包依赖

## 架构

core 包为以下主要组件定义了接口：

### 应用层
- **App**: 应用生命周期管理和服务编排
- **Service**: 用于构建可插拔应用服务的服务接口
- **ServiceItem**: 带有生命周期钩子的服务容器
- **Ctx**: 携带所有必要组件的请求上下文

### HTTP/Web 层
- **HTTPServer**: 具有模板、会话和路由支持的 Web 服务器
- **APIServer**: 用于构建 Web 服务的 RESTful API 服务器
- **HTTPRouter**: 高性能 HTTP 路由
- **Controller**: MVC 风格的控制器基础接口
- **Middleware**: 请求/响应处理中间件
- **Handler**: HTTP 请求处理函数类型
- **View**: 模板视图渲染接口
- **Template**: 模板解析和执行接口

### 数据库层
- **Database**: 数据库访问和查询执行
- **DatabaseGroup**: 多数据库连接管理
- **ActiveRecord**: 具有链式方法的 SQL 语句构建器
- **ResultSet**: 查询结果处理和映射
- **DBCache**: 数据库查询结果缓存

### 配置层
- **Config**: 应用配置管理
- **SubConfig**: 嵌套配置支持
- **CommonConfig**: 基本配置操作

### 日志层
- **Logger**: 多级别结构化日志
- **LoggerWriter**: 自定义日志输出写入器
- **LogLevel**: 日志严重级别（TRACE、DEBUG、INFO、WARN、ERROR、PANIC、FATAL）
- **LogFlag**: 日志格式标志

### 会话和状态管理
- **Session**: 用户会话数据存储和管理
- **SessionStorage**: 会话持久化层
- **Cookies**: HTTP cookie 操作

### 缓存层
- **Cache**: 支持 TTL 的通用缓存接口
- 操作: Get、Set、Delete、Incr/Decr、批量操作

### 国际化
- **I18n**: 多语言支持和翻译管理
- Accept-Language 头解析和匹配
- 模板安全的翻译

### 错误处理
- **Error**: 带堆栈跟踪的增强错误
- **StackFrame**: 用于调试的堆栈帧信息

### 工具类
- **Paginator**: 列表和搜索结果的分页助手
- **Params**: URL 参数处理
- **ResponseWriter**: 增强的 HTTP 响应写入器

## Provider 系统

core 包实现了一个复杂的 provider 模式，用于组件注册和获取：

### Provider 类型

所有主要组件都有相应的 provider 函数：

```go
type SessionProvider func() Session
type CacheProvider func(ctx Ctx) (Cache, error)
type DatabaseProvider func(ctx Ctx) (Database, error)
type LoggerProvider func(ctx Ctx, prefix string) Logger
type ConfigProvider func() Config
type TemplateProvider func(ctx Ctx, rootDir string) (Template, error)
// ... 还有更多
```

### Provider 注册

可以使用自定义键注册组件：

```go
// 注册自定义实现
gcore.RegisterLogger("mylogger", myLoggerProvider)
gcore.RegisterCache("redis", redisProvider)
gcore.RegisterDatabase("postgres", postgresProvider)

// 使用默认键注册
gcore.RegisterLogger(gcore.DefaultProviderKey, defaultLoggerProvider)
```

### Provider 获取

通过键获取 provider 或使用自动解析：

```go
// 通过键获取特定的 provider
logger := gcore.ProviderLogger("mylogger")

// 获取默认 provider（自动解析）
cache := gcore.ProviderCache()

// Auto-provider 自动使用最后注册的 provider
db := gcore.ProviderDatabase()
```

### AutoProvider

`AutoProvider` 实现了自动 provider 解析，采用最后注册优先策略：

```go
// 创建自定义 provider 注册表
autoProvider := gcore.NewAutoProvider()

// 注册多个实现
autoProvider.RegisterCache("memory", memCacheProvider)
autoProvider.RegisterCache("redis", redisCacheProvider)

// 默认使用最后注册的
cache := autoProvider.Cache() // 返回 redis provider
```

## 核心接口

### App 接口

应用生命周期和配置管理：

```go
type App interface {
    SetConfigFile(file string)
    SetConfig(cfg Config)
    AttachConfigFile(id, file string)
    Config(id ...string) Config
    Run() (err error)
    OnRun(fn func(Config) (err error))
    OnShutdown(fn func())
    AddService(item ServiceItem)
    SetLogger(logger Logger)
    Logger() Logger
    Stop()
    Ctx() Ctx
    SetCtx(Ctx)
}
```

### HTTPServer 接口

具有全栈功能的 Web 服务器：

```go
type HTTPServer interface {
    SetNotFoundHandler(fn func(ctx Ctx, tpl Template))
    SetErrorHandler(fn func(ctx Ctx, tpl Template, err interface{}))
    SetRouter(r HTTPRouter)
    Router() HTTPRouter
    SetTpl(t Template)
    Tpl() Template
    SetSessionStore(st SessionStorage)
    SessionStore() SessionStorage
    AddMiddleware0(m Middleware)
    AddMiddleware1(m Middleware)
    AddMiddleware2(m Middleware)
    AddMiddleware3(m Middleware)
    Listen() (err error)
    ListenTLS() (err error)
    // ... more methods
}
```

### Database 接口

使用 Active Record 模式的数据库操作：

```go
type Database interface {
    AR() (ar ActiveRecord)
    Stats() sql.DBStats
    Begin() (tx *sql.Tx, err error)
    Exec(ar ActiveRecord) (rs ResultSet, err error)
    ExecSQL(sqlStr string, values ...interface{}) (rs ResultSet, err error)
    Query(ar ActiveRecord) (rs ResultSet, err error)
    QuerySQL(sqlStr string, values ...interface{}) (rs ResultSet, err error)
}
```

### Cache 接口

具有常见操作的通用缓存：

```go
type Cache interface {
    Has(key string) (bool, error)
    Get(key string) (string, error)
    Set(key string, value string, ttl time.Duration) error
    Del(key string) error
    GetMulti(keys []string) (map[string]string, error)
    SetMulti(values map[string]string, ttl time.Duration) error
    Incr(key string) (int64, error)
    Decr(key string) (int64, error)
    Clear() error
}
```

### Logger 接口

多级别结构化日志：

```go
type Logger interface {
    Trace(v ...interface{})
    Debug(v ...interface{})
    Info(v ...interface{})
    Warn(v ...interface{})
    Error(v ...interface{})
    Panic(v ...interface{})
    Fatal(v ...interface{})
    
    // Formatted variants
    Tracef(format string, v ...interface{})
    Debugf(format string, v ...interface{})
    // ... more methods
    
    Level() LogLevel
    SetLevel(LogLevel)
    With(name string) Logger
    Writer() LoggerWriter
    AddWriter(LoggerWriter) Logger
    EnableAsync()
    WaitAsyncDone()
}
```

### Ctx 接口

包含所有依赖项的请求上下文：

```go
type Ctx interface {
    // 核心组件
    App() App
    Config() Config
    Logger() Logger
    Template() Template
    I18n() I18n
    
    // HTTP 组件
    Request() *http.Request
    Response() http.ResponseWriter
    Param() Params
    
    // 请求方法
    IsPOST() bool
    IsGET() bool
    IsAJAX() bool
    IsWebsocket() bool
    
    // 响应方法
    Write(data ...interface{}) (n int, err error)
    JSON(code int, data interface{}) (err error)
    JSONP(code int, data interface{}) (err error)
    Redirect(url string) string
    WriteFile(filepath string)
    
    // 请求数据
    GET(key string, Default ...string) string
    POST(key string, Default ...string) string
    Cookie(name string) string
    ClientIP() string
    
    // 会话和状态
    Set(key interface{}, value interface{})
    Get(key interface{}) (interface{}, bool)
    
    Clone() Ctx
    // ... 更多方法
}
```

## 使用示例

### 实现自定义缓存 Provider

```go
package mycache

import (
    "github.com/snail007/gmc/core"
    "time"
)

// 实现 Cache 接口
type MyCache struct {
    data map[string]string
}

func (c *MyCache) Get(key string) (string, error) {
    return c.data[key], nil
}

func (c *MyCache) Set(key string, value string, ttl time.Duration) error {
    c.data[key] = value
    return nil
}

// ... 实现其他 Cache 方法

// 创建 provider 函数
func NewMyCacheProvider(ctx gcore.Ctx) (gcore.Cache, error) {
    return &MyCache{
        data: make(map[string]string),
    }, nil
}

// 注册 provider
func init() {
    gcore.RegisterCache("mycache", NewMyCacheProvider)
}
```

### 实现自定义日志 Provider

```go
package mylogger

import "github.com/snail007/gmc/core"

type MyLogger struct {
    level gcore.LogLevel
    prefix string
}

func (l *MyLogger) Info(v ...interface{}) {
    // 你的日志实现
}

func (l *MyLogger) SetLevel(level gcore.LogLevel) {
    l.level = level
}

// ... 实现其他 Logger 方法

func NewMyLoggerProvider(ctx gcore.Ctx, prefix string) gcore.Logger {
    return &MyLogger{
        level: gcore.LogLeveInfo,
        prefix: prefix,
    }
}

func init() {
    gcore.RegisterLogger("mylogger", NewMyLoggerProvider)
}
```

### 使用 Provider 系统

```go
package main

import (
    "github.com/snail007/gmc/core"
)

func main() {
    // 获取 providers
    loggerProvider := gcore.ProviderLogger()
    cacheProvider := gcore.ProviderCache()
    
    // 创建实例
    logger := loggerProvider(nil, "myapp")
    cache, _ := cacheProvider(nil)
    
    // 使用组件
    logger.Info("应用已启动")
    cache.Set("key", "value", 0)
}
```

## 设计原则

### 1. 接口隔离
每个接口专注于特定的职责，使实现更简单、更易于维护。

### 2. 依赖倒置
高层模块依赖于抽象（接口）而不是具体实现。

### 3. Provider 模式
动态组件解析允许运行时配置和测试灵活性。

### 4. 上下文传播
`Ctx` 接口携带所有必要的依赖项，消除全局状态。

### 5. 中间件链
中间件系统允许灵活的请求/响应处理管道。

## 组件关系

```
┌─────────────────────────────────────────────┐
│                    App                      │
│          (生命周期和服务管理)                │
└──────────────────┬──────────────────────────┘
                   │
         ┌─────────┴─────────┐
         │                   │
    ┌────▼─────┐      ┌─────▼──────┐
    │HTTPServer│      │ APIServer  │
    │          │      │            │
    └────┬─────┘      └─────┬──────┘
         │                  │
    ┌────▼──────────────────▼────┐
    │      HTTPRouter             │
    │        (路由管理)           │
    └────┬────────────────────────┘
         │
    ┌────▼────┐
    │   Ctx   │◄──────────────────┐
    │         │                   │
    └────┬────┘                   │
         │                        │
    ┌────▼─────────────────────┐  │
    │    组件（通过 Ctx）      │  │
    │  - Config                │  │
    │  - Logger                │  │
    │  - Database              │  │
    │  - Cache                 │  │
    │  - Session               │  │
    │  - Template              │──┘
    │  - I18n                  │
    └──────────────────────────┘
```

## HTTP 请求流程

```
1. 客户端请求
   ↓
2. HTTPServer/APIServer 接收请求
   ↓
3. Middleware0（最高优先级）
   ↓
4. Middleware1
   ↓
5. Middleware2
   ↓
6. Middleware3（最低优先级）
   ↓
7. HTTPRouter 匹配路由
   ↓
8. Controller/Handler 执行
   ↓
9. View/Template 渲染（用于 Web）
   ↓
10. 响应发送到客户端
```

## 日志级别

日志支持以下级别（从低到高优先级）：

- **TRACE**: 细粒度调试信息
- **DEBUG**: 开发调试信息
- **INFO**: 信息性消息
- **WARN**: 警告消息
- **ERROR**: 错误消息
- **PANIC**: Panic 级别错误（可恢复）
- **FATAL**: 致命错误（应用退出）
- **NONE**: 禁用所有日志

## 测试

由于这个包只包含接口，测试由具体实现完成。在为你的实现编写测试时：

1. 确保所有接口方法都已实现
2. 测试边界情况和错误条件
3. 验证正确的资源清理
4. 在适用时测试并发访问

测试结构示例：

```go
func TestCacheImplementation(t *testing.T) {
    cache := NewMyCache()
    
    // 测试 Set/Get
    err := cache.Set("key", "value", 0)
    assert.NoError(t, err)
    
    val, err := cache.Get("key")
    assert.NoError(t, err)
    assert.Equal(t, "value", val)
    
    // 测试其他方法...
}
```

## 线程安全

接口实现应该记录它们的线程安全保证。一般来说：

- **Logger**: 应该是线程安全的
- **Cache**: 应该是线程安全的
- **Database**: 连接池是线程安全的；单个连接可能不是
- **Config**: 读操作是线程安全的；写操作应该同步
- **Ctx**: 不是线程安全的；每个 goroutine 创建新的上下文

## 性能考虑

### Provider 查找
- Provider 查找使用基于 map 的存储进行优化（O(1)）
- AutoProvider 缓存最后注册的 providers

### 上下文克隆
- 使用 `Ctx.Clone()` 进行轻量级上下文复制
- 避免在热路径中进行不必要的克隆

### 中间件
- 保持中间件逻辑轻量级
- 在需要时使用提前返回停止处理
- 按优先级排序中间件（最重要的在前）

## 最佳实践

### 1. Provider 注册
```go
// 在 init() 函数中注册 providers
func init() {
    gcore.RegisterLogger("default", NewDefaultLogger)
}
```

### 2. 上下文使用
```go
// 始终传递上下文，不要存储它
func HandleRequest(ctx gcore.Ctx) {
    logger := ctx.Logger()
    logger.Info("处理请求")
}
```

### 3. 错误处理
```go
// 使用 Error 接口获取丰富的错误信息
err := gcore.ProviderError().New("操作失败")
logger.Error(err.StackError(err))
```

### 4. 资源清理
```go
// 在 Stop/GracefulStop 中实现正确的清理
func (s *MyService) GracefulStop() {
    s.shutdown()
    s.releaseResources()
}
```

### 5. 配置
```go
// 使用 Sub() 的嵌套配置
dbConfig := cfg.Sub("database")
host := dbConfig.GetString("host")
```

## 相关包

core 包由各种 GMC 包实现：

- `github.com/snail007/gmc/module/app` - App 实现
- `github.com/snail007/gmc/module/log` - Logger 实现
- `github.com/snail007/gmc/module/config` - Config 实现
- `github.com/snail007/gmc/http` - HTTP 服务器实现
- `github.com/snail007/gmc/module/cache` - Cache 实现
- `github.com/snail007/gmc/module/db` - Database 实现

## 贡献

向此包添加新接口时：

1. 保持接口专注和最小化
2. 使用清晰的 godoc 注释记录所有方法
3. 考虑向后兼容性
4. 添加相应的 provider 类型
5. 使用示例更新此 README

## 许可证

Copyright 2020 The GMC Author. All rights reserved.
使用此源代码受 MIT 风格许可证约束，可以在 LICENSE 文件中找到。

## 链接

- **文档**: https://snail007.github.io/gmc/
- **GitHub**: https://github.com/snail007/gmc
- **Go Reference**: https://pkg.go.dev/github.com/snail007/gmc/core
