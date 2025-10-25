# GMC Error 模块

## 简介

GMC Error 模块提供了增强的错误处理功能，为 Go 语言的标准错误接口添加了堆栈跟踪（Stack Trace）支持。这对于调试和理解错误发生时的执行状态非常有用。

## 功能特性

- **堆栈跟踪**：自动捕获错误发生时的调用堆栈
- **错误包装**：支持包装现有错误并添加上下文
- **错误前缀**：为错误添加描述性前缀
- **Panic 恢复**：提供多种 panic 恢复机制
- **安全调用**：Try/Catch 风格的错误处理
- **兼容标准库**：实现标准 error 接口
- **详细信息**：支持打印完整的错误堆栈
- **Go 1.13+ 支持**：兼容 errors.Is 和 errors.As

## 安装

```bash
go get github.com/snail007/gmc/module/error
```

## 快速开始

### 创建带堆栈的错误

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/module/error"
)

func main() {
    // 创建新错误
    err := gerror.New("something went wrong")
    
    // 打印错误信息
    fmt.Println(err.Error())
    
    // 打印完整堆栈
    fmt.Println(err.ErrorStack())
}
```

### 包装现有错误

```go
package main

import (
    "errors"
    "fmt"
    "github.com/snail007/gmc/module/error"
)

func doSomething() error {
    return errors.New("database connection failed")
}

func main() {
    err := doSomething()
    if err != nil {
        // 包装错误并添加堆栈
        wrappedErr := gerror.Wrap(err)
        fmt.Println(wrappedErr.ErrorStack())
    }
}
```

### 使用错误前缀

```go
package main

import (
    "errors"
    "github.com/snail007/gmc/module/error"
)

func processOrder(orderID string) error {
    err := validateOrder(orderID)
    if err != nil {
        // 添加前缀提供上下文
        return gerror.New().WrapPrefix(err, "process order failed", 0)
    }
    return nil
}

func validateOrder(orderID string) error {
    return errors.New("invalid order ID")
}

func main() {
    err := processOrder("12345")
    if err != nil {
        // 输出: process order failed: invalid order ID
        println(err.Error())
    }
}
```

### Panic 恢复

```go
package main

import (
    "fmt"
    gcore "github.com/snail007/gmc/core"
    "github.com/snail007/gmc/module/error"
)

func riskyFunction() {
    defer gerror.Recover(func(err interface{}) {
        // 处理 panic
        if e, ok := err.(gcore.Error); ok {
            fmt.Println("Recovered from panic:")
            fmt.Println(e.ErrorStack())
        }
    })
    
    // 这里会 panic
    panic("something went wrong!")
}

func main() {
    riskyFunction()
    fmt.Println("Program continues...")
}
```

### Try/Catch 风格错误处理

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/module/error"
)

func main() {
    // Try 捕获 panic 并转换为 error
    err := gerror.Try(func() {
        // 可能会 panic 的代码
        panic("oops!")
    })
    
    if err != nil {
        fmt.Println("Caught:", err)
    }
    
    // TryWithStack 返回带堆栈的错误
    errWithStack := gerror.TryWithStack(func() {
        panic("error with stack!")
    })
    
    if errWithStack != nil {
        fmt.Println(errWithStack.ErrorStack())
    }
}
```

### 静默恢复 Panic

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/module/error"
)

func mayPanic() {
    defer gerror.RecoverNop() // 静默恢复，不做任何处理
    
    panic("this will be ignored")
}

func mayPanicWithCleanup() {
    defer gerror.RecoverNopFunc(func() {
        // 恢复后执行清理
        fmt.Println("Cleanup after panic")
    })
    
    panic("this will be recovered and cleanup executed")
}

func main() {
    mayPanic()
    fmt.Println("After mayPanic")
    
    mayPanicWithCleanup()
    fmt.Println("After mayPanicWithCleanup")
}
```

## API 参考

### 创建错误

```go
// 创建新的错误实例
func New(e ...interface{}) gcore.Error

// 包装错误并添加堆栈
func Wrap(e interface{}) gcore.Error

// 获取错误的堆栈信息
func Stack(e interface{}) string
```

### Error 类型方法

```go
// 创建新错误
New(e interface{}) gcore.Error

// 包装错误（从调用处开始堆栈）
Wrap(e interface{}) gcore.Error

// 包装错误（指定跳过的堆栈帧数）
WrapN(e interface{}, skip int) gcore.Error

// 包装错误并添加前缀
WrapPrefix(e interface{}, prefix string, skip int) gcore.Error
WrapPrefixN(e interface{}, prefix string, skip int) gcore.Error

// 格式化创建错误
Errorf(format string, a ...interface{}) gcore.Error

// 获取错误信息
Error() string

// 获取带堆栈的错误信息
ErrorStack() string

// 获取堆栈信息
Stack() []byte

// 获取堆栈帧
StackFrames() []gcore.StackFrame
```

### Panic 恢复函数

```go
// 恢复 panic 并执行回调
func Recover(f func(err interface{}))

// 静默恢复 panic
func RecoverNop()

// 恢复 panic 并执行清理函数
func RecoverNopFunc(f func())

// Try/Catch 风格错误处理
func Try(f func()) error
func TryWithStack(f func()) gcore.Error

// 解析 recover() 返回值为 error
func ParseRecover(e interface{}) error
```

### 堆栈帧信息

```go
type StackFrame struct {
    File           string // 文件路径
    LineNumber     int    // 行号
    Name           string // 函数名
    Package        string // 包名
    ProgramCounter uintptr
}
```

## 配置

### 设置最大堆栈深度

```go
package main

import (
    "github.com/snail007/gmc/module/error"
)

func main() {
    // 默认是 50，可以根据需要调整
    gerror.MaxStackDepth = 100
}
```

## 实用示例

### 1. HTTP 处理器错误恢复

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/snail007/gmc/module/error"
)

func handler(w http.ResponseWriter, r *http.Request) {
    defer gerror.Recover(func(err interface{}) {
        // 记录错误日志
        if e, ok := err.(error); ok {
            fmt.Printf("Panic in handler: %v\n", e)
        }
        
        // 返回 500 错误
        http.Error(w, "Internal Server Error", 500)
    })
    
    // 处理请求的代码
    // 如果 panic，会被上面的 Recover 捕获
}
```

### 2. 数据库操作错误包装

```go
package main

import (
    "database/sql"
    gcore "github.com/snail007/gmc/core"
    "github.com/snail007/gmc/module/error"
)

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) GetUser(id int) (user *User, err error) {
    defer func() {
        if e := recover(); e != nil {
            err = gerror.Wrap(e)
        }
    }()
    
    row := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
    user = &User{}
    err = row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, gerror.New().WrapPrefix(err, "failed to get user", 0)
    }
    
    return user, nil
}

type User struct {
    ID    int
    Name  string
    Email string
}
```

### 3. 链式错误处理

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/module/error"
)

func level3() error {
    return gerror.New("level 3 error")
}

func level2() error {
    err := level3()
    if err != nil {
        return gerror.New().WrapPrefix(err, "level 2", 0)
    }
    return nil
}

func level1() error {
    err := level2()
    if err != nil {
        return gerror.New().WrapPrefix(err, "level 1", 0)
    }
    return nil
}

func main() {
    err := level1()
    if err != nil {
        fmt.Println(err.Error())
        // 输出: level 1: level 2: level 3 error
        
        fmt.Println("\nFull stack trace:")
        if e, ok := err.(gcore.Error); ok {
            fmt.Println(e.ErrorStack())
        }
    }
}
```

### 4. Goroutine 错误处理

```go
package main

import (
    "fmt"
    "sync"
    "github.com/snail007/gmc/module/error"
)

func worker(id int, wg *sync.WaitGroup, errChan chan<- error) {
    defer wg.Done()
    defer func() {
        if e := recover(); e != nil {
            errChan <- gerror.Wrap(e)
        }
    }()
    
    // 模拟可能出错的工作
    if id == 3 {
        panic(fmt.Sprintf("worker %d failed", id))
    }
    
    fmt.Printf("Worker %d completed\n", id)
}

func main() {
    var wg sync.WaitGroup
    errChan := make(chan error, 5)
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, &wg, errChan)
    }
    
    wg.Wait()
    close(errChan)
    
    // 处理错误
    for err := range errChan {
        if e, ok := err.(gcore.Error); ok {
            fmt.Println("Error:", e.Error())
            fmt.Println("Stack:", e.ErrorStack())
        }
    }
}
```

### 5. 中间件错误处理

```go
package main

import (
    gcore "github.com/snail007/gmc/core"
    "github.com/snail007/gmc/module/error"
)

func ErrorRecoveryMiddleware(ctx gcore.Ctx, next func()) {
    defer gerror.Recover(func(err interface{}) {
        // 记录错误
        if e, ok := err.(gcore.Error); ok {
            ctx.Logger().Error(e.ErrorStack())
        }
        
        // 返回错误响应
        ctx.WriteJSON(gcore.M{
            "error": "Internal Server Error",
        })
        ctx.StatusCode(500)
    })
    
    next()
}
```

## 使用场景

1. **Web 应用**：捕获处理器中的 panic
2. **并发程序**：安全处理 goroutine 中的错误
3. **数据库操作**：为数据库错误添加上下文
4. **API 服务**：提供详细的错误堆栈用于调试
5. **中间件**：统一的错误恢复和处理
6. **微服务**：错误跟踪和调试

## 最佳实践

### 1. 始终在 defer 中使用

```go
func riskyFunction() {
    defer gerror.Recover(func(err interface{}) {
        // 处理错误
    })
    
    // 可能 panic 的代码
}
```

### 2. 为错误添加上下文

```go
if err != nil {
    return gerror.New().WrapPrefix(err, "failed to process user", 0)
}
```

### 3. 在生产环境隐藏堆栈

```go
if err != nil {
    if isDevelopment {
        // 开发环境显示完整堆栈
        log.Println(err.(gcore.Error).ErrorStack())
    } else {
        // 生产环境只显示错误信息
        log.Println(err.Error())
    }
}
```

### 4. 使用 Try 简化错误处理

```go
err := gerror.Try(func() {
    // 可能 panic 的代码
    mustDoSomething()
    mustDoAnotherThing()
})
```

## 性能考虑

1. **堆栈捕获开销**：捕获堆栈有一定性能开销，建议仅在错误发生时使用
2. **堆栈深度**：调整 `MaxStackDepth` 以平衡详细程度和性能
3. **生产环境**：考虑在生产环境禁用详细堆栈输出
4. **内存使用**：堆栈信息会占用额外内存

## 与标准库对比

### 标准 error

```go
err := errors.New("something went wrong")
// 只有错误消息，没有堆栈信息
```

### GMC Error

```go
err := gerror.New("something went wrong")
// 有错误消息 + 完整调用堆栈
fmt.Println(err.ErrorStack())
```

## 兼容性

- **Go 版本**：支持 Go 1.13+
- **errors.Is**：支持标准库的 errors.Is 检查
- **errors.As**：支持标准库的 errors.As 类型断言
- **标准 error 接口**：完全兼容，可以无缝替换

## 依赖

- 基于 [go-errors/errors](https://github.com/go-errors/errors) 增强
- 无其他外部依赖

## 注意事项

1. **Panic 恢复**：恢复函数必须在 defer 中调用
2. **堆栈深度**：默认最大 50 层，超过会被截断
3. **性能影响**：捕获堆栈有性能开销，不建议在热路径使用
4. **Go 版本**：某些功能需要 Go 1.13+ 支持

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [go-errors/errors](https://github.com/go-errors/errors)
- [Go 官方 errors 包](https://pkg.go.dev/errors)