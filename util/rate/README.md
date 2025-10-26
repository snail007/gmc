# grate 包 - 限流器

## 简介

grate 包提供了两种高性能的限流器实现，用于控制请求速率和带宽限制。支持高并发场景，所有实现都是线程安全的。

## 功能特性

- **滑动窗口限流器**：精确控制时间窗口内的请求数量上限
- **令牌桶限流器**：平滑限流，支持突发流量，适用于带宽控制
- **并发安全**：所有限流器都是线程安全的
- **灵活配置**：可配置速率、窗口大小、突发容量
- **高性能**：基于原子操作，支持高并发场景

## 安装

```bash
go get github.com/snail007/gmc/util/rate
```

## 两种限流器对比

### 滑动窗口限流器 (SlidingWindowLimiter)

**算法特点：**
- 严格限制时间窗口内的请求数量
- 时间窗口精确划分，无法超过设定上限
- 适合需要严格控制 QPS 的场景

**适用场景：**
- ✅ **API 接口限流**：严格限制每秒/每分钟的 API 调用次数
- ✅ **防止恶意请求**：限制同一 IP 在时间窗口内的请求次数
- ✅ **第三方服务调用**：遵守第三方 API 的速率限制（如每秒 100 次）
- ✅ **秒杀活动**：严格控制秒杀接口的并发请求数
- ✅ **数据库查询限流**：防止数据库过载

**优点：**
- 精确控制，不会超过限制
- 实时响应，立即生效

**缺点：**
- 在窗口边界可能出现流量突刺
- 不支持短时突发流量

### 令牌桶限流器 (TokenBucketLimiter)

**算法特点：**
- 按固定速率生成令牌，请求消耗令牌
- 允许一定程度的突发流量（burst）
- 平滑限流，流量更均匀

**适用场景：**
- ✅ **带宽限制**：控制下载/上传速率（如限制每秒 1MB）
- ✅ **流量整形**：平滑网络流量，避免瞬时峰值
- ✅ **消息队列消费**：控制消息处理速率
- ✅ **文件上传下载**：限制文件传输速率
- ✅ **视频流传输**：控制视频数据传输速率
- ✅ **需要支持突发流量的场景**：如短时间内允许更多请求

**优点：**
- 支持突发流量（burst）
- 流量更平滑，避免锯齿状
- 适合带宽控制

**缺点：**
- 控制相对宽松，可能短时间内超过平均速率

## 快速开始

### 滑动窗口限流器示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/rate"
)

func main() {
    // 创建限流器：每秒最多 10 个请求
    limiter := grate.NewSlidingWindowLimiter(10, time.Second)
    
    for i := 0; i < 20; i++ {
        if limiter.Allow() {
            fmt.Println("Request", i, "allowed")
        } else {
            fmt.Println("Request", i, "rejected")
        }
        time.Sleep(50 * time.Millisecond)
    }
}
```

### 令牌桶限流器示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/rate"
)

func main() {
    // 创建限流器：每秒 10 个令牌，突发容量 10
    limiter := grate.NewTokenBucketLimiter(10, time.Second)
    
    // 支持突发流量的限流器：每秒 10 个令牌，突发容量 20
    burstLimiter := grate.NewTokenBucketBurstLimiter(10, time.Second, 20)
    
    for i := 0; i < 15; i++ {
        if limiter.Allow() {
            fmt.Println("Request", i, "allowed")
        } else {
            fmt.Println("Request", i, "rejected")
        }
    }
}
```

## API 参考

### 滑动窗口限流器 (SlidingWindowLimiter)

#### 创建限流器

```go
func NewSlidingWindowLimiter(capacity int, duration time.Duration) *SlidingWindowLimiter
```

创建一个新的滑动窗口限流器。

**参数：**
- `capacity`：时间窗口内允许的最大请求数
- `duration`：时间窗口大小

**示例：**
```go
// 每秒最多 100 个请求
limiter := grate.NewSlidingWindowLimiter(100, time.Second)

// 每分钟最多 1000 个请求
limiter := grate.NewSlidingWindowLimiter(1000, time.Minute)

// 每小时最多 10000 个请求
limiter := grate.NewSlidingWindowLimiter(10000, time.Hour)
```

#### 方法

##### Allow

```go
func (l *SlidingWindowLimiter) Allow() bool
```

检查是否允许一个请求。

**返回值：**
- `true`：请求被允许
- `false`：请求被拒绝（超过限流）

**示例：**
```go
if limiter.Allow() {
    // 处理请求
    handleRequest()
} else {
    // 返回限流错误
    return errors.New("rate limit exceeded")
}
```

##### AllowN

```go
func (l *SlidingWindowLimiter) AllowN(n int32) bool
```

检查是否允许 n 个请求。

**参数：**
- `n`：请求数量

**返回值：**
- `true`：请求被允许
- `false`：请求被拒绝

**示例：**
```go
// 批量操作，一次消耗 5 个配额
if limiter.AllowN(5) {
    processBatch(5)
}
```

##### Capacity

```go
func (l *SlidingWindowLimiter) Capacity() int
```

获取限流器的容量（时间窗口内允许的最大请求数）。

##### Duration

```go
func (l *SlidingWindowLimiter) Duration() time.Duration
```

获取时间窗口大小。

### 令牌桶限流器 (TokenBucketLimiter)

#### 创建限流器

```go
func NewTokenBucketLimiter(count int, duration time.Duration) *TokenBucketLimiter
```

创建一个新的令牌桶限流器，突发容量等于速率。

**参数：**
- `count`：在 duration 时间内生成的令牌数量
- `duration`：时间周期

**示例：**
```go
// 每秒 100 个令牌，突发容量 100
limiter := grate.NewTokenBucketLimiter(100, time.Second)

// 每秒 1MB 带宽限制（假设每个令牌代表 1KB）
limiter := grate.NewTokenBucketLimiter(1024, time.Second)
```

```go
func NewTokenBucketBurstLimiter(count int, duration time.Duration, burst int) *TokenBucketLimiter
```

创建一个新的令牌桶限流器，支持自定义突发容量。

**参数：**
- `count`：在 duration 时间内生成的令牌数量
- `duration`：时间周期
- `burst`：令牌桶容量（突发流量上限）

**示例：**
```go
// 每秒 100 个令牌，但允许突发 200 个
limiter := grate.NewTokenBucketBurstLimiter(100, time.Second, 200)
```

#### 方法

TokenBucketLimiter 继承了 `golang.org/x/time/rate.Limiter` 的所有方法：

##### Allow

```go
func (l *TokenBucketLimiter) Allow() bool
```

检查是否允许一个请求（立即返回）。

##### Wait

```go
func (l *TokenBucketLimiter) Wait(ctx context.Context) error
```

等待直到有可用的令牌（阻塞等待）。

**示例：**
```go
// 阻塞等待直到可以处理请求
if err := limiter.Wait(ctx); err != nil {
    return err
}
handleRequest()
```

##### Reserve

```go
func (l *TokenBucketLimiter) Reserve() *rate.Reservation
```

预留令牌，返回预留信息。

**示例：**
```go
r := limiter.Reserve()
if !r.OK() {
    return errors.New("rate limit exceeded")
}
time.Sleep(r.Delay())
handleRequest()
```

## 使用场景详解

### 场景 1：API 接口限流（使用滑动窗口）

```go
package main

import (
    "net/http"
    "github.com/snail007/gmc/util/rate"
)

var limiter = grate.NewSlidingWindowLimiter(100, time.Second)

func apiHandler(w http.ResponseWriter, r *http.Request) {
    if !limiter.Allow() {
        http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
        return
    }
    
    // 处理正常请求
    w.Write([]byte("OK"))
}
```

### 场景 2：IP 限流（使用滑动窗口）

```go
package main

import (
    "net/http"
    "sync"
    "github.com/snail007/gmc/util/rate"
)

var (
    limiters = make(map[string]*grate.SlidingWindowLimiter)
    mu       sync.RWMutex
)

func getLimiter(ip string) *grate.SlidingWindowLimiter {
    mu.RLock()
    limiter, exists := limiters[ip]
    mu.RUnlock()
    
    if !exists {
        mu.Lock()
        limiter = grate.NewSlidingWindowLimiter(10, time.Minute)
        limiters[ip] = limiter
        mu.Unlock()
    }
    
    return limiter
}

func handler(w http.ResponseWriter, r *http.Request) {
    ip := r.RemoteAddr
    limiter := getLimiter(ip)
    
    if !limiter.Allow() {
        http.Error(w, "Too many requests", http.StatusTooManyRequests)
        return
    }
    
    // 处理请求
}
```

### 场景 3：带宽限制（使用令牌桶）

```go
package main

import (
    "context"
    "io"
    "net/http"
    "github.com/snail007/gmc/util/rate"
)

// 限制下载速率为 1MB/s
var downloadLimiter = grate.NewTokenBucketLimiter(1024*1024, time.Second)

type limitedReader struct {
    r       io.Reader
    limiter *grate.TokenBucketLimiter
}

func (lr *limitedReader) Read(p []byte) (n int, err error) {
    n, err = lr.r.Read(p)
    if n > 0 {
        // 等待令牌，控制读取速率
        lr.limiter.WaitN(context.Background(), n)
    }
    return
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
    file, _ := os.Open("large-file.zip")
    defer file.Close()
    
    // 使用限速 reader
    lr := &limitedReader{
        r:       file,
        limiter: downloadLimiter,
    }
    
    io.Copy(w, lr)
}
```

### 场景 4：消息队列消费限流（使用令牌桶）

```go
package main

import (
    "context"
    "github.com/snail007/gmc/util/rate"
)

// 每秒处理 50 条消息
var consumerLimiter = grate.NewTokenBucketLimiter(50, time.Second)

func consumeMessages(ctx context.Context, messages <-chan Message) {
    for msg := range messages {
        // 等待令牌，控制消费速率
        if err := consumerLimiter.Wait(ctx); err != nil {
            return
        }
        
        // 处理消息
        processMessage(msg)
    }
}
```

### 场景 5：第三方 API 调用限流（使用滑动窗口）

```go
package main

import (
    "github.com/snail007/gmc/util/rate"
)

// 某第三方 API 限制：每分钟 60 次调用
var apiLimiter = grate.NewSlidingWindowLimiter(60, time.Minute)

func callThirdPartyAPI(endpoint string, data interface{}) error {
    if !apiLimiter.Allow() {
        return errors.New("rate limit: please wait before next call")
    }
    
    // 调用第三方 API
    return doAPICall(endpoint, data)
}
```

### 场景 6：文件上传限流（使用令牌桶）

```go
package main

import (
    "context"
    "io"
    "net/http"
    "github.com/snail007/gmc/util/rate"
)

// 限制上传速率为 2MB/s，允许突发 5MB
var uploadLimiter = grate.NewTokenBucketBurstLimiter(2*1024*1024, time.Second, 5*1024*1024)

type limitedWriter struct {
    w       io.Writer
    limiter *grate.TokenBucketLimiter
}

func (lw *limitedWriter) Write(p []byte) (n int, err error) {
    // 等待令牌
    if err := lw.limiter.WaitN(context.Background(), len(p)); err != nil {
        return 0, err
    }
    return lw.w.Write(p)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    file, _ := os.Create("upload.dat")
    defer file.Close()
    
    lw := &limitedWriter{
        w:       file,
        limiter: uploadLimiter,
    }
    
    io.Copy(lw, r.Body)
}
```

## 并发安全

所有限流器都是线程安全的，可以在多个 goroutine 中并发使用：

```go
limiter := grate.NewSlidingWindowLimiter(100, time.Second)

// 在多个 goroutine 中并发调用
for i := 0; i < 1000; i++ {
    go func() {
        if limiter.Allow() {
            // 安全处理
        }
    }()
}
```

## 性能基准测试

运行基准测试：

```bash
# 滑动窗口限流器基准测试
go test -bench=BenchmarkSlidingWindow -benchmem

# 令牌桶限流器基准测试  
go test -bench=BenchmarkTokenBucket -benchmem
```

## 选择指南

### 何时使用滑动窗口限流器？

- ✅ 需要**严格限制**时间窗口内的请求数量
- ✅ **API 接口限流**，遵守明确的速率限制
- ✅ **防刷**、**防爬虫**场景
- ✅ 需要**精确控制 QPS**

### 何时使用令牌桶限流器？

- ✅ 需要**平滑限流**，允许短时突发
- ✅ **带宽限制**、流量整形
- ✅ **文件传输**速率控制
- ✅ **消息队列**消费限流
- ✅ 需要支持**突发流量**的场景

## 注意事项

1. **选择合适的限流器**：
   - 严格限流 → 滑动窗口
   - 带宽控制/平滑限流 → 令牌桶

2. **时间窗口大小**：
   - 滑动窗口：内部分成 10 个子窗口，时间窗口不宜过小
   - 建议最小 100ms

3. **突发容量**：
   - 令牌桶允许突发，burst 设置需要合理
   - 过大的 burst 可能导致瞬时压力

4. **内存占用**：
   - 滑动窗口：每个限流器占用固定内存（10 个 bucket）
   - 令牌桶：内存占用极小

5. **并发性能**：
   - 两种限流器都基于原子操作，支持高并发
   - 滑动窗口在极高并发下性能略优

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [golang.org/x/time/rate 文档](https://pkg.go.dev/golang.org/x/time/rate)
