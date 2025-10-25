# grate 包

## 简介

grate 包提供了限流器实现，用于控制请求速率。

## 功能特性

- **滑动窗口限流**：基于滑动窗口算法
- **并发安全**：线程安全的限流器
- **灵活配置**：可配置速率和窗口大小

## 安装

```bash
go get github.com/snail007/gmc/util/rate
```

## 快速开始

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

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
