# GMC Cache 模块

## 简介

GMC Cache 模块提供统一的缓存接口，支持多种缓存后端，包括 Redis、内存缓存和文件缓存。支持多数据源配置和管理。

## 功能特性

- **多种缓存后端**：支持 Redis、Memory、File 三种缓存类型
- **多数据源支持**：可同时配置和使用多个缓存实例
- **统一接口**：提供统一的缓存操作 API
- **键前缀**：支持为每个缓存设置键前缀
- **连接池**：Redis 缓存支持连接池配置
- **自动过期**：支持键的自动过期
- **调试模式**：支持调试日志输出

## 安装

```bash
go get github.com/snail007/gmc/module/cache
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc"
)

func main() {
    // 加载配置
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    err := cfg.ReadInConfig()
    if err != nil {
        panic(err)
    }
    
    // 初始化缓存
    err = gmc.Cache.Init(cfg)
    if err != nil {
        panic(err)
    }
    
    // 获取默认缓存实例
    cache := gmc.Cache.Cache()
    
    // 设置键值（5 秒过期）
    err = cache.Set("key", "value", 5*time.Second)
    if err != nil {
        panic(err)
    }
    
    // 获取值
    value, err := cache.Get("key")
    if err != nil {
        panic(err)
    }
    fmt.Println("Value:", value)
    
    // 删除键
    err = cache.Del("key")
    if err != nil {
        panic(err)
    }
}
```

### 使用多个缓存实例

```go
package main

import (
    "github.com/snail007/gmc"
    "time"
)

func main() {
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    gmc.Cache.Init(cfg)
    
    // 使用默认缓存
    defaultCache := gmc.Cache.Cache()
    defaultCache.Set("key1", "value1", time.Minute)
    
    // 使用指定 ID 的缓存
    sessionCache := gmc.Cache.Cache("session")
    sessionCache.Set("session_id", "user123", time.Hour)
    
    // 使用不同类型的缓存
    redisCache := gmc.Cache.Redis("redis1")
    memoryCache := gmc.Cache.Memory("mem1")
    fileCache := gmc.Cache.File("file1")
}
```

### Redis 缓存

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/module/cache"
)

func main() {
    // 创建 Redis 缓存配置
    cfg := &gcache.RedisCacheConfig{
        Addr:            "127.0.0.1:6379",
        Password:        "",
        DBNum:           0,
        Prefix:          "myapp:",
        MaxIdle:         10,
        MaxActive:       30,
        IdleTimeout:     300 * time.Second,
        MaxConnLifetime: 3600 * time.Second,
        Wait:            true,
        Timeout:         10 * time.Second,
        Debug:           false,
    }
    
    // 创建 Redis 缓存实例
    cache := gcache.NewRedisCache(cfg)
    
    // 使用缓存
    cache.Set("user:1", "Alice", time.Hour)
    value, _ := cache.Get("user:1")
    fmt.Println(value)
}
```

### 内存缓存

```go
package main

import (
    "time"
    "github.com/snail007/gmc/module/cache"
)

func main() {
    // 创建内存缓存配置
    cfg := &gcache.MemCacheConfig{
        CleanupInterval: 10 * time.Minute,
    }
    
    // 创建内存缓存实例
    cache := gcache.NewMemCache(cfg)
    
    // 使用缓存
    cache.Set("temp", "data", 5*time.Minute)
    value, _ := cache.Get("temp")
}
```

### 文件缓存

```go
package main

import (
    "time"
    "github.com/snail007/gmc/module/cache"
)

func main() {
    // 创建文件缓存配置
    cfg := &gcache.FileCacheConfig{
        Dir:             "./cache",
        CleanupInterval: 10 * time.Minute,
    }
    
    // 创建文件缓存实例
    cache, err := gcache.NewFileCache(cfg)
    if err != nil {
        panic(err)
    }
    
    // 使用缓存
    cache.Set("file_key", "content", time.Hour)
    value, _ := cache.Get("file_key")
}
```

## 配置文件

### 完整配置示例

```toml
[cache]
# 默认缓存 ID
default = "default"

# Redis 缓存配置
[[cache.redis]]
enable = true
id = "default"
address = "127.0.0.1:6379"
password = ""
dbnum = 0
prefix = "myapp:"
# 超时时间（秒）
timeout = 10
# 最大空闲连接数
maxidle = 10
# 最大活动连接数
maxactive = 30
# 空闲连接超时（秒）
idletimeout = 300
# 连接最大生命周期（秒）
maxconnlifetime = 3600
# 是否等待可用连接
wait = true
# 调试模式
debug = false

# 会话缓存（Redis）
[[cache.redis]]
enable = true
id = "session"
address = "127.0.0.1:6379"
password = ""
dbnum = 1
prefix = "sess:"
timeout = 10
maxidle = 5
maxactive = 20

# 内存缓存配置
[[cache.memory]]
enable = true
id = "mem1"
# 清理过期键的间隔（秒）
cleanupinterval = 600

# 文件缓存配置
[[cache.file]]
enable = true
id = "file1"
# 缓存文件存储目录
dir = "./cache"
# 清理过期文件的间隔（秒）
cleanupinterval = 600
```

## API 参考

### 全局函数

```go
// 初始化缓存系统
func Init(cfg gcore.Config) error

// 获取默认缓存或指定 ID 的缓存
func Cache(id ...string) gcore.Cache

// 获取 Redis 缓存
func Redis(id ...string) gcore.Cache

// 获取内存缓存
func Memory(id ...string) gcore.Cache

// 获取文件缓存
func File(id ...string) gcore.Cache

// 设置日志记录器
func SetLogger(logger gcore.Logger)
```

### Cache 接口

```go
type Cache interface {
    // 设置键值
    Set(key string, value interface{}, ttl time.Duration) error
    
    // 获取值
    Get(key string) (interface{}, error)
    
    // 删除键
    Del(key string) error
    
    // 检查键是否存在
    Exists(key string) (bool, error)
    
    // 设置过期时间
    Expire(key string, ttl time.Duration) error
    
    // 增加数值
    Incr(key string, delta int64) (int64, error)
    
    // 减少数值
    Decr(key string, delta int64) (int64, error)
    
    // 获取多个键的值
    MGet(keys ...string) ([]interface{}, error)
    
    // 设置多个键值
    MSet(items map[string]interface{}, ttl time.Duration) error
    
    // 关闭缓存连接
    Close() error
}
```

### 配置结构

#### RedisCacheConfig

```go
type RedisCacheConfig struct {
    Debug           bool
    Prefix          string
    Logger          gcore.Logger
    Addr            string
    Password        string
    DBNum           int
    MaxIdle         int
    MaxActive       int
    IdleTimeout     time.Duration
    Wait            bool
    MaxConnLifetime time.Duration
    Timeout         time.Duration
}
```

#### MemCacheConfig

```go
type MemCacheConfig struct {
    CleanupInterval time.Duration
}
```

#### FileCacheConfig

```go
type FileCacheConfig struct {
    Dir             string
    CleanupInterval time.Duration
}
```

## 使用场景

1. **会话存储**：存储用户会话数据
2. **页面缓存**：缓存渲染后的 HTML 页面
3. **API 响应缓存**：缓存 API 查询结果
4. **数据库查询缓存**：缓存频繁查询的数据
5. **限流计数**：使用 Incr/Decr 实现限流
6. **临时数据存储**：存储临时文件上传信息

## 最佳实践

### 1. 使用键前缀

```go
// 为不同模块使用不同前缀，避免键冲突
cfg := &gcache.RedisCacheConfig{
    Prefix: "user:",
    // ...
}
```

### 2. 设置合理的过期时间

```go
// 根据数据特性设置过期时间
cache.Set("hot_data", value, 5*time.Minute)    // 热数据
cache.Set("cold_data", value, 1*time.Hour)     // 冷数据
cache.Set("static_data", value, 24*time.Hour)  // 静态数据
```

### 3. 错误处理

```go
value, err := cache.Get("key")
if err != nil {
    if err.Error() == "key not found" {
        // 键不存在，从数据库加载
        value = loadFromDatabase()
        cache.Set("key", value, time.Hour)
    } else {
        // 其他错误
        log.Printf("Cache error: %v", err)
    }
}
```

### 4. 批量操作

```go
// 使用 MSet 批量设置，提高性能
items := map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
}
cache.MSet(items, time.Hour)

// 使用 MGet 批量获取
values, err := cache.MGet("key1", "key2", "key3")
```

## 性能考虑

1. **连接池配置**：合理设置 MaxIdle 和 MaxActive
2. **键前缀长度**：前缀不要过长，影响性能
3. **值大小**：避免缓存过大的对象
4. **过期时间**：合理设置过期时间，避免内存浪费
5. **清理间隔**：内存和文件缓存的清理间隔影响性能

## 注意事项

1. **Redis 连接**：确保 Redis 服务可访问
2. **文件权限**：文件缓存需要对目录有写权限
3. **内存限制**：内存缓存受进程内存限制
4. **并发安全**：所有缓存实现都是并发安全的
5. **序列化**：复杂对象会被序列化存储

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [Redis 文档](https://redis.io/documentation)