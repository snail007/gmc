# gtesting 包

## 简介

gtesting 包提供了测试相关的工具函数。

## 功能特性

- **断言工具**：便捷的测试断言
- **Mock 工具**：模拟对象
- **测试辅助**：测试辅助函数

## 安装

```bash
go get github.com/snail007/gmc/util/testing
```

## 快速开始

```go
package main

import (
    "testing"
    "github.com/snail007/gmc/util/testing"
)

func TestExample(t *testing.T) {
    // 断言相等
    gtesting.Equal(t, 1+1, 2)
    
    // 断言不为 nil
    gtesting.NotNil(t, &struct{}{})
    
    // 断言为真
    gtesting.True(t, true)
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
