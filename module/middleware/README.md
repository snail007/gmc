# GMC Middleware 模块

## 简介

GMC Middleware 模块提供 HTTP 中间件支持。目前内置了 **AccessLog（访问日志）** 中间件，用于记录 HTTP 请求的详细信息。

## 中间件架构与生命周期

### 架构概览

GMC 的中间件架构允许在请求处理的不同阶段插入自定义逻辑。请求从客户端进入后，会依次经过不同优先级的中间件层，最终到达控制器处理，响应则按相反顺序返回。

<p align="center">
  <img src="../../doc/images/http-and-api-server-architecture.png" alt="GMC Middleware Architecture" width="800"/>
</p>

### 生命周期说明

中间件在 HTTP 请求处理的整个生命周期中扮演关键角色：

1. **请求到达** → HTTP Server 接收请求
2. **Middleware0** → 路由匹配前执行（全局保护层）
3. **路由匹配** → 匹配到对应的路由和控制器
4. **Middleware1** → 控制器方法执行前（认证/预处理层）
5. **Controller** → 执行控制器业务逻辑
6. **Middleware2** → 控制器方法执行后（响应处理层）
7. **Middleware3** → 响应返回前（日志/统计层）
8. **响应返回** → 将响应发送给客户端

**关键特点：**
- 中间件按级别顺序执行
- 每层中间件可以选择继续或停止处理
- 返回 `true` 停止后续处理，返回 `false` 继续
- Middleware3 在响应后执行，可获取完整的请求信息

## 内置中间件

### AccessLog - 访问日志中间件

记录 HTTP 请求的详细信息，支持自定义格式、异步写入、日志轮转和压缩。

**功能特性：**
- ✅ 自定义日志格式和占位符
- ✅ 异步写入，不影响性能
- ✅ 自动日志轮转（按日期/小时）
- ✅ Gzip 压缩历史日志
- ✅ 智能获取客户端真实 IP
- ✅ 高性能协程池处理

📖 **详细文档和配置**: [**AccessLog 中间件指南**](accesslog/README.md)

#### 快速开始

```go
import (
    "github.com/snail007/gmc"
    "github.com/snail007/gmc/module/middleware/accesslog"
)

func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    // 添加访问日志中间件
    s.AddMiddleware3(accesslog.NewFromConfig(s.Config()))
    
    // 配置路由...
    app.Run()
}
```

#### 配置示例（app.toml）

```toml
[accesslog]
dir = "./logs"
filename = "access_%Y%m%d.log"
gzip = true
format = "$req_time $client_ip $uri $status_code ${time_used}ms"
```

## 中间件级别详解

GMC 支持 4 个中间件级别，每个级别在请求处理流程中的不同位置执行：

### Middleware0 - 路由前执行

**执行时机：** 在路由匹配之前，最先执行  
**函数签名：** `func(c gmc.C, s *gmc.HTTPServer) bool`  
**返回值：** `true` 停止处理，`false` 继续执行

```go
s.AddMiddleware0(func(c gmc.C, s *gmc.HTTPServer) bool {
    // 最先执行，可用于：全局限流、IP 黑名单、请求预处理
    // 在这里可以访问原始请求，但还未匹配路由
    
    // 示例：IP 黑名单
    if isBlocked(c.ClientIP()) {
        c.WriteHeader(403)
        c.Write("Forbidden")
        return true // 停止后续处理
    }
    
    return false // 继续处理
})
```

**适用场景：**
- 全局限流和防护
- IP 黑白名单
- 请求签名验证
- 请求日志记录（开始时间）

### Middleware1 - 控制器前执行

**执行时机：** 路由匹配后，控制器方法执行前  
**函数签名：** `func(c gmc.C) bool`  
**返回值：** `true` 停止处理，`false` 继续执行

```go
s.AddMiddleware1(func(c gmc.C) bool {
    // 已经匹配到路由，可以进行认证、权限检查等
    
    // 示例：认证检查
    token := c.Request().Header.Get("Authorization")
    if token == "" {
        c.WriteHeader(401)
        c.WriteJSON(gmc.M{"error": "Unauthorized"})
        return true // 停止后续处理
    }
    
    // 验证 token 并设置用户信息
    // userID := validateToken(token)
    // c.Set("user_id", userID)
    
    return false // 继续处理
})
```

**适用场景：**
- 用户认证
- 权限检查
- 参数验证
- 请求数据预处理

### Middleware2 - 控制器后执行

**执行时机：** 控制器方法执行后，响应返回前  
**函数签名：** `func(c gmc.C) bool`  
**返回值：** `true` 停止处理，`false` 继续执行

```go
s.AddMiddleware2(func(c gmc.C) bool {
    // 控制器已执行，可以处理响应数据
    
    // 示例：添加响应头
    c.Response().Header().Set("X-Response-Time", 
        time.Since(c.Get("start_time").(time.Time)).String())
    
    return false
})
```

**适用场景：**
- 响应数据转换
- 添加响应头
- 数据加密
- 响应缓存

### Middleware3 - 响应后执行

**执行时机：** 响应返回前，最后执行  
**函数签名：** `func(c gmc.C, status int, message string)`  
**无返回值**（响应已生成）

```go
s.AddMiddleware3(func(c gmc.C, status int, message string) {
    // 可获取完整的响应状态码和信息
    // 通常用于日志记录和统计
    
    // 示例：记录完整请求信息
    duration := time.Since(c.Get("start_time").(time.Time))
    fmt.Printf("[%d] %s %s - %v\n", 
        status,
        c.Request().Method,
        c.Request().URL.Path,
        duration)
})
```

**适用场景：**
- 访问日志记录（推荐）
- 性能统计
- 审计日志
- 监控指标上报

### 执行顺序示意

```
请求 → Middleware0 → 路由匹配 → Middleware1 → 控制器 → Middleware2 → Middleware3 → 响应
  ↓         ↓              ↓              ↓           ↓           ↓             ↓          ↓
阻止?      阻止?          找到路由?       阻止?      执行业务    处理响应      记录日志    返回
```

### 中间件选择建议

| 类型 | 推荐级别 | 说明 |
|------|---------|------|
| 限流/防护 | Middleware0 | 最早拦截，保护服务器 |
| 认证/鉴权 | Middleware1 | 路由后，业务前 |
| 数据处理 | Middleware2 | 业务后，响应前 |
| 日志记录 | Middleware3 | 最后执行，记录完整信息 |

## 常见中间件实现示例

虽然 GMC 核心只提供了 AccessLog，但你可以轻松实现其他常见中间件：

### 1. CORS 跨域中间件

```go
func CORSMiddleware() gcore.Middleware0 {
    return func(c gmc.C, s *gmc.HTTPServer) bool {
        c.Response().Header().Set("Access-Control-Allow-Origin", "*")
        c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // 处理 OPTIONS 预检请求
        if c.Request().Method == "OPTIONS" {
            c.WriteHeader(204)
            return true // 停止后续处理
        }
        return false
    }
}

// 使用
s.AddMiddleware0(CORSMiddleware())
```

### 2. 认证中间件

```go
func AuthMiddleware() gcore.Middleware1 {
    return func(c gmc.C) bool {
        token := c.Request().Header.Get("Authorization")
        
        if token == "" {
            c.WriteHeader(401)
            c.WriteJSON(gmc.M{"error": "Unauthorized"})
            return true // 停止后续处理
        }
        
        // 验证 token...
        // userID := validateToken(token)
        // c.Set("user_id", userID)
        
        return false
    }
}

// 使用
s.AddMiddleware1(AuthMiddleware())
```

### 3. Recovery 错误恢复中间件

```go
func RecoveryMiddleware() gcore.Middleware1 {
    return func(c gmc.C) bool {
        defer func() {
            if err := recover(); err != nil {
                // 记录错误
                fmt.Printf("Panic recovered: %v\n", err)
                
                // 返回 500 错误
                c.WriteHeader(500)
                c.WriteJSON(gmc.M{
                    "error": "Internal Server Error",
                })
            }
        }()
        return false
    }
}

// 使用
s.AddMiddleware1(RecoveryMiddleware())
```

### 4. 请求计时中间件

```go
func TimingMiddleware() gcore.Middleware1 {
    return func(c gmc.C) bool {
        start := time.Now()
        c.Set("start_time", start)
        return false
    }
}

func TimingLogMiddleware() gcore.Middleware3 {
    return func(c gmc.C, status int, message string) {
        if start, ok := c.Get("start_time").(time.Time); ok {
            duration := time.Since(start)
            fmt.Printf("[%s] %s - %d - %v\n", 
                c.Request().Method, 
                c.Request().URL.Path, 
                status, 
                duration)
        }
    }
}

// 使用
s.AddMiddleware1(TimingMiddleware())
s.AddMiddleware3(TimingLogMiddleware())
```

### 5. 简单限流中间件

```go
import "sync"

func RateLimitMiddleware(maxRequests int, window time.Duration) gcore.Middleware0 {
    var (
        mu       sync.Mutex
        requests = make(map[string][]time.Time)
    )
    
    return func(c gmc.C, s *gmc.HTTPServer) bool {
        ip := c.ClientIP()
        now := time.Now()
        
        mu.Lock()
        defer mu.Unlock()
        
        // 清理过期记录
        if times, ok := requests[ip]; ok {
            var valid []time.Time
            for _, t := range times {
                if now.Sub(t) < window {
                    valid = append(valid, t)
                }
            }
            requests[ip] = valid
        }
        
        // 检查限流
        if len(requests[ip]) >= maxRequests {
            c.WriteHeader(429)
            c.WriteJSON(gmc.M{"error": "Too Many Requests"})
            return true
        }
        
        // 记录请求
        requests[ip] = append(requests[ip], now)
        return false
    }
}

// 使用：每分钟最多 100 个请求
s.AddMiddleware0(RateLimitMiddleware(100, time.Minute))
```

## 中间件组合示例

### Web 应用推荐组合

```go
func SetupWebMiddleware(s *gmc.HTTPServer) {
    // 1. 错误恢复（最先）
    s.AddMiddleware1(RecoveryMiddleware())
    
    // 2. 请求计时
    s.AddMiddleware1(TimingMiddleware())
    
    // 3. 认证（可选，某些路由跳过）
    s.AddMiddleware1(func(c gmc.C) bool {
        // 公开路径不需要认证
        if strings.HasPrefix(c.Request().URL.Path, "/public/") {
            return false
        }
        return AuthMiddleware()(c)
    })
    
    // 4. 访问日志（最后）
    s.AddMiddleware3(accesslog.NewFromConfig(s.Config()))
    s.AddMiddleware3(TimingLogMiddleware())
}
```

### API 服务推荐组合

```go
func SetupAPIMiddleware(s *gmc.HTTPServer) {
    // 1. CORS（最先）
    s.AddMiddleware0(CORSMiddleware())
    
    // 2. 限流
    s.AddMiddleware0(RateLimitMiddleware(1000, time.Hour))
    
    // 3. 错误恢复
    s.AddMiddleware1(RecoveryMiddleware())
    
    // 4. 认证
    s.AddMiddleware1(AuthMiddleware())
    
    // 5. 访问日志
    s.AddMiddleware3(accesslog.NewFromConfig(s.Config()))
}
```

## 使用第三方中间件

GMC 的中间件接口简单灵活，可以轻松适配第三方中间件或将标准 HTTP Handler 包装为 GMC 中间件。

### 包装标准 HTTP Handler

```go
func WrapHTTPHandler(handler http.Handler) gcore.Middleware0 {
    return func(c gmc.C, s *gmc.HTTPServer) bool {
        handler.ServeHTTP(c.Response(), c.Request())
        return false
    }
}

// 使用第三方包
import "github.com/some/middleware"

s.AddMiddleware0(WrapHTTPHandler(middleware.NewSomeMiddleware()))
```

## 最佳实践

1. **合理选择中间件级别**
   - Middleware0: 路由前，用于全局保护（限流、黑名单）
   - Middleware1: 业务前，用于认证和预处理
   - Middleware2: 业务后，用于响应处理
   - Middleware3: 最后，用于日志和统计

2. **性能考虑**
   - 避免在中间件中执行耗时操作
   - 使用异步处理日志等非关键操作
   - 合理使用缓存减少重复计算

3. **错误处理**
   - 使用 Recovery 中间件捕获 panic
   - 返回友好的错误信息给客户端
   - 记录详细错误日志便于排查

4. **安全性**
   - 在 Middleware0 层面实现限流和防护
   - 及时更新认证 token
   - 记录异常访问行为

5. **可维护性**
   - 将中间件逻辑封装成独立函数
   - 使用配置文件控制中间件行为
   - 编写单元测试验证中间件功能

## 性能提示

- ✅ 按需加载中间件，避免不必要的处理
- ✅ 使用条件判断跳过特定路径
- ✅ 缓存认证结果，减少数据库查询
- ✅ 日志使用异步写入
- ✅ 限流使用高效的数据结构

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [HTTP Server 文档](../../http/server/README.md)
- [路由文档](../../http/router/README.md)
- [控制器文档](../../http/controller/README.md)
- [AccessLog 中间件详细文档](accesslog/README.md)
