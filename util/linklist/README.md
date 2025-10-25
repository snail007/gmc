# glinklist 包

## 简介

glinklist 包提供了线程安全的双向链表实现。

## 功能特性

- **双向链表**：支持前后遍历
- **线程安全**：所有操作都是线程安全的
- **灵活插入**：支持头部、尾部、指定位置插入

## 安装

```bash
go get github.com/snail007/gmc/util/linklist
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/linklist"
)

func main() {
    list := glinklist.New()
    
    // 添加元素
    list.Append(1)
    list.Append(2)
    list.Append(3)
    
    // 遍历
    list.Range(func(v interface{}) bool {
        fmt.Println(v)
        return true
    })
    
    // 获取长度
    fmt.Println("Length:", list.Len())
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
