# gset 包

## 简介

gset 包提供了线程安全的集合（Set）实现。

## 功能特性

- **无重复元素**：自动去重
- **线程安全**：所有操作都是线程安全的
- **集合操作**：并集、交集、差集
- **成员检查**：快速检查元素是否存在

## 安装

```bash
go get github.com/snail007/gmc/util/set
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/set"
)

func main() {
    set := gset.New()
    
    // 添加元素
    set.Add("apple")
    set.Add("banana")
    set.Add("apple") // 重复元素不会添加
    
    // 检查是否存在
    if set.Has("apple") {
        fmt.Println("Set contains apple")
    }
    
    // 获取大小
    fmt.Println("Size:", set.Len())
    
    // 遍历
    set.Range(func(item interface{}) bool {
        fmt.Println(item)
        return true
    })
    
    // 删除
    set.Remove("banana")
    
    // 转换为切片
    items := set.Items()
    fmt.Println(items)
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
