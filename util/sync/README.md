# gsync 包

## 简介

gsync 包提供了并发同步工具。

## 功能特性

- **WaitGroup 增强**：带超时的 WaitGroup
- **Once 增强**：可重置的 Once
- **锁工具**：便捷的锁操作
- **并发控制**：并发数量控制

## 安装

```bash
go get github.com/snail007/gmc/util/sync
```

## 快速开始

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/sync"
)

func main() {
    // 带超时的 WaitGroup
    wg := gsync.NewWaitGroupTimeout(5 * time.Second)
    
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            time.Sleep(time.Second)
            fmt.Println("Task", n, "done")
        }(i)
    }
    
    if wg.Wait() {
        fmt.Println("All tasks completed")
    } else {
        fmt.Println("Timeout")
    }
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
