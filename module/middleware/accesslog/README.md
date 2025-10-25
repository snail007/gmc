# GMC AccessLog 中间件

## 简介

GMC AccessLog 中间件提供了 HTTP 访问日志记录功能，可以记录每个 HTTP 请求的详细信息，包括请求时间、客户端 IP、URI、响应状态码、处理时间等。支持自定义日志格式、异步写入、日志轮转和压缩。

## 功能特性

- **自定义格式**：支持自定义日志格式和占位符
- **异步写入**：异步日志写入，不影响请求性能
- **日志轮转**：支持按日期自动轮转日志文件
- **日志压缩**：支持 gzip 压缩历史日志
- **丰富信息**：记录客户端真实 IP、响应时间、状态码等
- **高性能**：使用协程池处理日志，最多缓冲 100000 条日志

## 安装

```bash
go get github.com/snail007/gmc/module/middleware/accesslog
```

## 快速开始

### 在 GMC Web 应用中使用

```go
package main

import (
    "github.com/snail007/gmc"
    "github.com/snail007/gmc/module/middleware/accesslog"
)

func main() {
    // 创建 HTTP 服务器
    app := gmc.New.App()
    
    // 加载配置
    cfg := app.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    // 创建 Web 服务器
    s := gmc.New.HTTPServer(app.Ctx())
    
    // 添加访问日志中间件（使用 Middleware3 级别）
    s.AddMiddleware3(accesslog.NewFromConfig(s.Config()))
    
    // 配置路由...
    
    // 启动服务器
    app.AddService(gmc.ServiceItem{
        Service: s,
        AfterInit: func(s gmc.Service, cfg gcore.Config) (err error) {
            // 初始化路由等
            return
        },
    })
    
    app.Run()
}
```

### 在路由初始化中使用

```go
package router

import (
    "github.com/snail007/gmc"
    "github.com/snail007/gmc/module/middleware/accesslog"
)

func InitRouter(s *gmc.HTTPServer) {
    // 添加访问日志中间件
    s.AddMiddleware3(accesslog.NewFromConfig(s.Config()))
    
    // 配置其他路由...
    r := s.Router()
    r.Controller("/api", new(APIController))
}
```

## 配置文件

### app.toml 配置

在配置文件中添加 `[accesslog]` 段：

```toml
##############################################################
# Web & API 访问日志中间件配置
##############################################################
[accesslog]
# 日志文件目录
dir = "./logs"

# 归档目录（可选，为空则不归档）
archive_dir = ""

# 日志文件名（支持时间占位符）
# 可用占位符：
#   %Y: 年份（2024）
#   %m: 月份（01-12）
#   %d: 日期（01-31）
#   %H: 小时（00-23）
filename = "access_%Y%m%d.log"

# 是否 gzip 压缩旧日志
gzip = true

# 日志格式（支持自定义占位符）
# 可用占位符：
#   $host         : 请求的 Host（包含端口），如 domain:port
#   $uri          : 请求路径
#   $query        : 完整查询字符串
#   $status_code  : HTTP 响应状态码
#   $time_used    : 请求处理耗时（毫秒）
#   $req_time     : 请求时间，格式：2024-10-25 15:33:55
#   $client_ip    : 客户端真实 IP（从 X-Forwarded-For, X-Real-IP 或 RemoteAddr 获取）
#   $remote_addr  : 远程地址（包含端口）
#   $local_addr   : 本地服务地址
format = "$req_time $host $uri?$query $status_code ${time_used}ms"
```

### 配置示例

#### 简单格式

```toml
[accesslog]
dir = "./logs"
filename = "access.log"
gzip = false
format = "$req_time $uri $status_code"
```

输出示例：
```
2024-10-25 15:33:55 /api/users 200
2024-10-25 15:34:01 /api/products 404
```

#### 详细格式

```toml
[accesslog]
dir = "./logs"
filename = "access_%Y%m%d.log"
gzip = true
format = "$req_time $client_ip $host $uri?$query $status_code ${time_used}ms"
```

输出示例：
```
2024-10-25 15:33:55 192.168.1.100 api.example.com /api/users?page=1 200 45ms
2024-10-25 15:34:01 192.168.1.101 api.example.com /api/products?id=123 404 12ms
```

#### Apache Combined 风格

```toml
[accesslog]
dir = "./logs"
filename = "access_%Y%m%d.log"
gzip = true
format = "$client_ip - - [$req_time] \"$uri\" $status_code ${time_used}ms"
```

输出示例：
```
192.168.1.100 - - [2024-10-25 15:33:55] "/api/users" 200 45ms
```

## API 参考

### NewFromConfig

```go
func NewFromConfig(c gcore.Config) gcore.Middleware
```

从配置创建访问日志中间件。

**参数：**
- `c`: GMC Config 对象

**返回：**
- `gcore.Middleware`: 中间件函数

## 占位符说明

### 请求相关

- **$host**：请求的 Host 头，包含端口号
  - 示例：`api.example.com:8080` 或 `example.com`

- **$uri**：请求的 URI 路径（不含查询字符串）
  - 示例：`/api/users` 或 `/index.html`

- **$query**：完整的查询字符串
  - 示例：`page=1&size=10` 或空字符串

### 响应相关

- **$status_code**：HTTP 响应状态码
  - 示例：`200`, `404`, `500`

- **$time_used**：请求处理耗时（毫秒）
  - 示例：`45`, `123`, `1002`

### 时间相关

- **$req_time**：请求时间，格式：`2006-01-02 15:04:05`
  - 示例：`2024-10-25 15:33:55`

### IP 地址相关

- **$client_ip**：客户端真实 IP（按优先级从以下来源获取）
  1. X-Forwarded-For 头（第一个 IP）
  2. X-Real-IP 头
  3. RemoteAddr（去除端口号）
  - 示例：`192.168.1.100` 或 `10.0.0.5`

- **$remote_addr**：请求的远程地址（包含端口）
  - 示例：`192.168.1.100:54321`

- **$local_addr**：服务器本地监听地址
  - 示例：`127.0.0.1:8080`

## 使用示例

### 示例 1：标准 Web 应用日志

```toml
[accesslog]
dir = "./logs"
filename = "web_%Y%m%d.log"
gzip = true
format = "$req_time [$client_ip] $uri $status_code ${time_used}ms"
```

### 示例 2：API 服务日志

```toml
[accesslog]
dir = "/var/log/myapi"
filename = "api_%Y%m%d_%H.log"  # 按小时轮转
gzip = true
format = "$req_time $client_ip $host $uri?$query $status_code ${time_used}ms"
```

### 示例 3：性能监控日志

```toml
[accesslog]
dir = "./logs/performance"
filename = "perf_%Y%m%d.log"
gzip = true
# 重点记录响应时间
format = "${time_used}ms $status_code $uri $client_ip"
```

使用此格式可以快速发现慢请求：
```
1523ms 200 /api/heavy-query 192.168.1.100
45ms 200 /api/users 192.168.1.101
2034ms 500 /api/timeout 192.168.1.102
```

### 示例 4：调试详细日志

```toml
[accesslog]
dir = "./logs/debug"
filename = "debug_%Y%m%d.log"
gzip = false
format = "[$req_time] RemoteAddr=$remote_addr ClientIP=$client_ip LocalAddr=$local_addr Host=$host URI=$uri Query=$query Status=$status_code Time=${time_used}ms"
```

## 中间件级别

GMC 支持多个中间件级别，访问日志通常使用 **Middleware3**：

```go
s.AddMiddleware0(...)  // 最先执行，在路由匹配之前
s.AddMiddleware1(...)  // 路由匹配后，控制器方法前
s.AddMiddleware2(...)  // 控制器方法后
s.AddMiddleware3(accesslog.NewFromConfig(s.Config()))  // 最后执行
```

使用 Middleware3 的原因：
1. 可以获取完整的响应状态码
2. 可以计算准确的请求处理时间
3. 不会影响业务逻辑的执行

## 性能考虑

1. **异步写入**：日志异步写入，不阻塞请求处理
2. **协程池**：使用协程池复用 goroutine，避免频繁创建
3. **缓冲队列**：最多缓冲 100000 条日志，防止内存溢出
4. **批量写入**：日志批量写入磁盘，提高 I/O 效率

## 日志轮转

日志文件名支持时间占位符，实现自动轮转：

- **按天轮转**：`access_%Y%m%d.log` → `access_20241025.log`
- **按小时轮转**：`access_%Y%m%d_%H.log` → `access_20241025_15.log`
- **按月轮转**：`access_%Y%m.log` → `access_202410.log`

## 日志压缩

设置 `gzip = true` 后，历史日志文件会自动压缩为 `.gz` 格式，节省磁盘空间。

## 客户端 IP 获取

中间件智能识别客户端真实 IP，按以下优先级：

1. **X-Forwarded-For** 头（反向代理场景）
2. **X-Real-IP** 头（Nginx 等设置）
3. **RemoteAddr**（直连场景）

适用于以下架构：
```
Client → Nginx/LB → GMC Server
```

## 最佳实践

1. **生产环境启用 gzip**：节省磁盘空间
2. **按天轮转日志**：便于日志管理和归档
3. **使用合适的格式**：根据日志分析需求定制格式
4. **监控日志大小**：定期清理旧日志
5. **分离不同类型日志**：Web 和 API 使用不同的日志文件

## 日志分析

可以使用标准工具分析访问日志：

### 统计状态码分布

```bash
cat access_20241025.log | awk '{print $4}' | sort | uniq -c
```

### 查找慢请求（>1000ms）

```bash
cat access_20241025.log | awk '$5 ~ /ms$/ && $5+0 > 1000'
```

### 统计访问最多的 URI

```bash
cat access_20241025.log | awk '{print $3}' | sort | uniq -c | sort -rn | head -10
```

## 注意事项

1. 确保日志目录有写权限
2. 监控磁盘空间，避免日志填满磁盘
3. 日志格式中的占位符大小写敏感
4. 异步日志有缓冲，服务器异常退出可能丢失部分日志

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [GMC HTTP Server](../../http/server/README.md)
- [GMC 中间件](../middleware/README.md)
- [GMC Log 模块](../log/README.md)
