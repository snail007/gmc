# gpaginator 包

## 简介

gpaginator 包提供了分页器功能，用于数据分页展示。

## 功能特性

- **分页计算**：计算总页数、偏移量等
- **页码生成**：生成页码列表
- **边界处理**：自动处理边界情况

## 安装

```bash
go get github.com/snail007/gmc/util/paginator
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/paginator"
)

func main() {
    // 创建分页器：总数 100，每页 10 条，当前第 3 页
    p := gpaginator.New(100, 10, 3)
    
    fmt.Println("总页数:", p.TotalPage())
    fmt.Println("当前页:", p.CurrentPage())
    fmt.Println("偏移量:", p.Offset())
    fmt.Println("每页数量:", p.PerPage())
    
    // 获取页码列表
    pages := p.Pages()
    fmt.Println("页码:", pages)
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
