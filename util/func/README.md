# gfunc 包

## 简介

gfunc 包提供了安全的函数调用和错误处理工具，主要用于 panic 捕获和错误恢复。

## 功能特性

- **安全调用**：捕获函数执行中的 panic
- **错误包装**：将 panic 转换为带堆栈的错误
- **恢复机制**：多种 panic 恢复策略
- **检查错误**：提供类似异常抛出的错误检查机制

## 安装

```bash
go get github.com/snail007/gmc/util/func
```

## 快速开始

### 安全调用函数

```go
package main

import (
    "fmt"
    gfunc "github.com/snail007/gmc/util/func"
)

func main() {
    // SafetyCall 捕获 panic 并返回 error
    err := gfunc.SafetyCall(func() {
        panic("something went wrong")
    })
    
    if err != nil {
        fmt.Println("捕获到错误:", err)
    }
}
```

### 带堆栈的错误

```go
package main

import (
    "fmt"
    gcore "github.com/snail007/gmc/core"
    gfunc "github.com/snail007/gmc/util/func"
)

func main() {
    // SafetyCallError 返回带堆栈的 gcore.Error
    err := gfunc.SafetyCallError(func() {
        panic("error with stack")
    })
    
    if err != nil {
        // 打印堆栈信息
        fmt.Println(err.(gcore.Error).ErrorStack())
    }
}
```

### 自定义恢复处理

```go
package main

import (
    "fmt"
    gcore "github.com/snail007/gmc/core"
    gfunc "github.com/snail007/gmc/util/func"
)

func main() {
    defer gfunc.Recover(func(err gcore.Error) {
        fmt.Println("捕获到 panic:", err.Error())
        fmt.Println("堆栈:", err.ErrorStack())
    })
    
    panic("critical error")
}
```

### 检查错误并抛出

```go
package main

import (
    gfunc "github.com/snail007/gmc/util/func"
)

func processData() (err error) {
    // 捕获 CheckError 抛出的 panic
    defer gfunc.CatchCheckError()
    
    // 如果 err 不为 nil，会 panic
    gfunc.CheckError(someOperation())
    gfunc.CheckError(anotherOperation())
    
    return nil
}

func someOperation() error {
    return nil // 或返回错误
}

func anotherOperation() error {
    return nil
}
```

## API 参考

### SafetyCall

```go
func SafetyCall(f func()) error
```

安全调用函数，捕获 panic 并返回 error。

### SafetyCallError

```go
func SafetyCallError(f func()) gcore.Error
```

安全调用函数，捕获 panic 并返回带堆栈的 gcore.Error。

### Recover

```go
func Recover(f func(gcore.Error))
```

在 defer 中使用，捕获 panic 并调用自定义处理函数。

### RecoverNop

```go
func RecoverNop()
```

在 defer 中使用，捕获 panic 但不做任何处理。

### RecoverNopAndFunc  

```go
func RecoverNopAndFunc(f func())
```

捕获 panic 后调用指定函数。

### CheckError

```go
func CheckError(err error)
```

检查错误，如果不为 nil 则 panic。配合 `CatchCheckError` 使用。

### CatchCheckError

```go
func CatchCheckError()
```

在 defer 中使用，捕获 `CheckError` 抛出的 panic。

## 使用场景

1. **HTTP 处理器**：防止 panic 导致服务崩溃
2. **Goroutine**：安全执行并发任务
3. **插件系统**：隔离第三方代码的 panic
4. **错误链**：简化错误检查代码

## 注意事项

1. `CheckError` 必须与 `CatchCheckError` 配合使用
2. `Recover` 系列函数必须在 defer 中调用
3. 堆栈信息会增加内存开销

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
