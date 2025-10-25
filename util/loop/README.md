# gloop 包

## 简介

gloop 包提供了事件循环（Event Loop）实现，用于异步事件处理。

## 功能特性

- **事件循环**：非阻塞事件处理
- **定时任务**：支持定时和延时任务
- **异步执行**：在循环中执行异步任务
- **优雅关闭**：支持优雅停止循环

## 安装

```bash
go get github.com/snail007/gmc/util/loop
```

## 快速开始

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/loop"
)

func main() {
    loop := gloop.New()
    
    // 添加任务
    loop.Run(func() {
        fmt.Println("Task executed")
    })
    
    // 延时任务
    loop.RunAfter(2*time.Second, func() {
        fmt.Println("Delayed task")
    })
    
    // 定时任务
    loop.RunEvery(1*time.Second, func() bool {
        fmt.Println("Periodic task")
        return true // 返回 true 继续，false 停止
    })
    
    // 启动循环
    loop.Start()
    
    time.Sleep(10 * time.Second)
    loop.Stop()
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
