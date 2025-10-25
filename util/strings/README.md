# gstrings 包

## 简介

gstrings 包提供了字符串处理相关的工具函数。

## 功能特性

- **字符串转换**：各种格式转换
- **字符串处理**：分割、连接、替换
- **随机字符串**：生成随机字符串
- **编码解码**：Base64、URL 编码等

## 安装

```bash
go get github.com/snail007/gmc/util/strings
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/strings"
)

func main() {
    // 驼峰转下划线
    snake := gstrings.CamelToSnake("HelloWorld")
    fmt.Println(snake) // hello_world
    
    // 下划线转驼峰
    camel := gstrings.SnakeToCamel("hello_world")
    fmt.Println(camel) // HelloWorld
    
    // 首字母大写
    title := gstrings.Title("hello")
    fmt.Println(title) // Hello
    
    // 首字母小写
    lower := gstrings.Untitle("Hello")
    fmt.Println(lower) // hello
}
```

## API 参考

- `CamelToSnake(s) string`：驼峰转下划线
- `SnakeToCamel(s) string`：下划线转驼峰
- `Title(s) string`：首字母大写
- `Untitle(s) string`：首字母小写
- `Random(n) string`：随机字符串
- `Contains(s, substr) bool`：是否包含子串
- `HasPrefix(s, prefix) bool`：是否有前缀
- `HasSuffix(s, suffix) bool`：是否有后缀

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
