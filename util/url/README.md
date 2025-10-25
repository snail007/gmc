# gurl 包

## 简介

gurl 包提供了 URL 处理相关的工具函数。

## 功能特性

- **URL 解析**：解析 URL 各部分
- **参数处理**：查询参数的编码解码
- **URL 构建**：构建 URL
- **URL 验证**：验证 URL 合法性

## 安装

```bash
go get github.com/snail007/gmc/util/url
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/url"
)

func main() {
    // 解析 URL
    u, _ := gurl.Parse("https://example.com/path?key=value")
    fmt.Println("Host:", u.Host)
    fmt.Println("Path:", u.Path)
    
    // 构建查询参数
    params := map[string]string{
        "name": "John Doe",
        "age":  "30",
    }
    query := gurl.BuildQuery(params)
    fmt.Println("Query:", query) // name=John+Doe&age=30
    
    // 编码 URL 部分
    encoded := gurl.PathEscape("hello world")
    fmt.Println("Encoded:", encoded) // hello%20world
}
```

## API 参考

- `Parse(rawurl) (*URL, error)`：解析 URL
- `BuildQuery(params) string`：构建查询字符串
- `PathEscape(s) string`：转义路径部分
- `QueryEscape(s) string`：转义查询参数
- `IsValid(rawurl) bool`：验证 URL 是否合法

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
