# grand 包

## 简介

grand 包提供了随机数生成工具。

## 功能特性

- **随机字符串**：生成随机字符串
- **随机数字**：生成随机数字
- **UUID**：生成 UUID
- **线程安全**：所有操作都是线程安全的

## 安装

```bash
go get github.com/snail007/gmc/util/rand
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/rand"
)

func main() {
    r := grand.New()
    
    // 随机字符串
    str := grand.String(10)
    fmt.Println("Random string:", str)
    
    // 随机数字
    num := r.Int()
    fmt.Println("Random number:", num)
    
    // 指定范围的随机数
    rangeNum := r.IntRange(1, 100)
    fmt.Println("Random in range:", rangeNum)
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
