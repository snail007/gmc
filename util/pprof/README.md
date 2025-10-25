# gpprof 包

## 简介

gpprof 包提供了性能分析（pprof）工具的便捷封装。

## 功能特性

- **CPU 分析**：CPU 性能分析
- **内存分析**：内存使用分析
- **便捷启动**：简化 pprof 服务启动

## 安装

```bash
go get github.com/snail007/gmc/util/pprof
```

## 快速开始

```go
package main

import (
    "github.com/snail007/gmc/util/pprof"
)

func main() {
    // 启动 pprof HTTP 服务
    gpprof.Start(":6060")
    
    // 访问 http://localhost:6060/debug/pprof/
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
