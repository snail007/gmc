# glist 包

## 简介

glist 包提供了线程安全的列表（List）容器，支持常见的列表操作。

## 功能特性

- **线程安全**：所有操作都是线程安全的
- **动态大小**：自动扩展
- **丰富操作**：添加、删除、插入、查找等
- **批量操作**：支持批量添加和合并

## 安装

```bash
go get github.com/snail007/gmc/util/list
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/list"
)

func main() {
    list := glist.New()
    
    // 添加元素
    list.Add("apple", "banana", "cherry")
    
    // 获取元素
    fmt.Println(list.Get(0)) // apple
    
    // 设置元素
    list.Set(1, "blueberry")
    
    // 删除元素
    list.Remove(2)
    
    // 遍历
    list.Range(func(index int, value interface{}) bool {
        fmt.Printf("%d: %v\n", index, value)
        return true
    })
    
    // 长度
    fmt.Println("Length:", list.Len())
}
```

## API 参考

- `New() *List`：创建列表
- `Add(v ...interface{})`：添加元素
- `AddFront(v interface{})`：添加到头部
- `Get(idx int) interface{}`：获取元素
- `Set(idx int, v interface{})`：设置元素
- `Remove(idx int)`：删除元素
- `Len() int`：获取长度
- `Range(func(int, interface{}) bool)`：遍历
- `Clear()`：清空
- `Clone() *List`：克隆

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
