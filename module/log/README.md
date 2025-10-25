# GMC Log 模块

## 简介

GMC Log 模块提供强大的日志记录功能，支持多种输出格式、异步日志、日志分级、日志分组等特性。

## 功能特性

- **多级别日志**：支持 Trace、Debug、Info、Warn、Error、Panic、Fatal
- **多种格式**：支持 Text、JSON 格式
- **异步日志**：支持异步写入，提高性能
- **日志分组**：支持多个独立的日志实例
- **灵活输出**：支持输出到文件、标准输出、自定义 Writer
- **调用栈**：错误日志自动记录调用栈
- **日志轮转**：支持按大小和时间轮转
- **颜色输出**：终端输出支持颜色

## 快速开始

### 基本使用

```go
package main

import (
    "github.com/snail007/gmc/module/log"
)

func main() {
    logger := glog.New()
    
    // 不同级别的日志
    logger.Trace("trace message")
    logger.Debug("debug message")
    logger.Info("info message")
    logger.Warn("warn message")
    logger.Error("error message")
    
    // 格式化输出
    logger.Infof("User %s logged in", "John")
    logger.Errorf("Failed to connect: %v", err)
}
```

### 从配置初始化

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    // 从配置初始化日志
    logger := gmc.New.Logger(cfg.Sub("log"))
    
    logger.Info("Application started")
}
```

### 异步日志

```go
package main

import (
    "github.com/snail007/gmc/module/log"
)

func main() {
    logger := glog.New()
    logger.EnableAsync() // 启用异步
    
    defer logger.WaitAsyncDone() // 等待异步日志写入完成
    
    for i := 0; i < 1000; i++ {
        logger.Infof("Message %d", i)
    }
}
```

### 日志分组

```go
package main

import (
    "github.com/snail007/gmc/module/log"
)

func main() {
    // 访问日志
    accessLog := glog.NewLogger("access")
    accessLog.Info("User visited homepage")
    
    // 错误日志
    errorLog := glog.NewLogger("error")
    errorLog.Error("Database connection failed")
    
    // 业务日志
    bizLog := glog.NewLogger("business")
    bizLog.Info("Order created")
}
```

## 配置文件

### app.toml 日志配置

```toml
[log]
# 日志输出: stdout, stderr, file
output = "stdout"

# 日志级别: trace, debug, info, warn, error, panic, fatal
level = "info"

# 日志格式: text, json
format = "text"

# 启用异步日志
async = false

# 文件输出配置
[[log.file]]
name = "app"
level = "info"
# 日志文件路径
path = "./logs/app.log"
# 最大大小（MB）
maxsize = 100
# 保留文件数
maxbackups = 10
# 保留天数
maxage = 30
# 是否压缩
compress = true

# 错误日志文件
[[log.file]]
name = "error"
level = "error"
path = "./logs/error.log"
maxsize = 100
maxbackups = 10
maxage = 90
compress = true
```

## API 参考

### 创建 Logger

```go
// 创建新 Logger
func New() gcore.Logger

// 创建命名 Logger
func NewLogger(name string) gcore.Logger

// 从配置创建
func NewFromConfig(cfg gcore.Config, prefix string) gcore.Logger
```

### 日志方法

```go
// 基本日志方法
Trace(args ...interface{})
Debug(args ...interface{})
Info(args ...interface{})
Warn(args ...interface{})
Error(args ...interface{})
Panic(args ...interface{})
Fatal(args ...interface{})

// 格式化日志方法
Tracef(format string, args ...interface{})
Debugf(format string, args ...interface{})
Infof(format string, args ...interface{})
Warnf(format string, args ...interface{})
Errorf(format string, args ...interface{})
Panicf(format string, args ...interface{})
Fatalf(format string, args ...interface{})

// 带字段的日志
WithFields(fields map[string]interface{}) gcore.Logger
```

### 配置方法

```go
// 设置日志级别
SetLevel(level string)

// 启用/禁用异步
EnableAsync()
DisableAsync()

// 是否异步
Async() bool

// 等待异步日志完成
WaitAsyncDone()

// 设置输出
SetOutput(w io.Writer)

// 设置格式
SetFormat(format string) // "text" or "json"
```

## 使用场景

1. **应用日志**：记录应用运行状态
2. **访问日志**：记录 HTTP 请求
3. **错误日志**：记录错误和异常
4. **审计日志**：记录关键操作
5. **调试日志**：开发调试

## 最佳实践

### 1. 结构化日志

```go
logger.WithFields(map[string]interface{}{
    "user_id": 123,
    "action":  "login",
    "ip":      "192.168.1.1",
}).Info("User logged in")
```

### 2. 错误日志带上下文

```go
if err != nil {
    logger.WithFields(map[string]interface{}{
        "error":  err.Error(),
        "method": "ProcessOrder",
        "order_id": orderId,
    }).Error("Failed to process order")
}
```

### 3. 使用不同日志级别

```go
// 开发环境
logger.SetLevel("debug")

// 生产环境
logger.SetLevel("info")

// 仅错误
logger.SetLevel("error")
```

### 4. 异步日志提升性能

```go
logger := glog.New()
logger.EnableAsync()

// 确保程序退出前日志写入完成
defer logger.WaitAsyncDone()
```

## 性能考虑

- **异步日志**：高并发场景使用异步日志
- **日志级别**：生产环境设置为 info 或 warn
- **日志轮转**：及时清理旧日志文件
- **JSON 格式**：结构化日志便于分析但性能稍差

## 注意事项

1. **Fatal 和 Panic**：会导致程序终止或 panic
2. **异步日志**：程序退出前需调用 `WaitAsyncDone()`
3. **文件权限**：确保对日志目录有写权限
4. **日志大小**：注意日志文件大小，及时轮转

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
