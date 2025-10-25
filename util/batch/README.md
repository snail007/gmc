# gbatch 包

## 简介

gbatch 包提供了一个高效的批量任务执行器，支持并发执行多个任务，并提供多种等待策略。该包使用 goroutine 池来控制并发数，避免创建过多的 goroutine。

## 功能特性

- **并发控制**：可配置工作线程数量，控制并发执行的任务数
- **多种等待策略**：支持等待所有任务、等待首个成功、等待首个完成
- **错误处理**：支持 panic 捕获和自定义错误处理
- **上下文支持**：任务可访问 context.Context 进行取消和超时控制
- **线程安全**：所有操作都是线程安全的

## 安装

```bash
go get github.com/snail007/gmc/util/batch
```

## 快速开始

### 基本使用 - 等待所有任务完成

```go
package main

import (
    "context"
    "fmt"
    "github.com/snail007/gmc/util/batch"
)

func main() {
    executor := gbatch.NewBatchExecutor()
    
    // 设置工作线程数
    executor.SetWorkers(5)
    
    // 添加任务
    executor.AppendTask(
        func(ctx context.Context) (interface{}, error) {
            return "task1 result", nil
        },
        func(ctx context.Context) (interface{}, error) {
            return "task2 result", nil
        },
        func(ctx context.Context) (interface{}, error) {
            return nil, fmt.Errorf("task3 error")
        },
    )
    
    // 等待所有任务完成
    results := executor.WaitAll()
    
    // 处理结果
    for i, result := range results {
        fmt.Printf("Task %d: ", i)
        if result.Err() != nil {
            fmt.Printf("Error: %v\n", result.Err())
        } else {
            fmt.Printf("Value: %v\n", result.Value())
        }
    }
}
```

### 等待首个成功的任务

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/snail007/gmc/util/batch"
)

func main() {
    executor := gbatch.NewBatchExecutor()
    
    // 添加多个任务，只要有一个成功就返回
    executor.AppendTask(
        func(ctx context.Context) (interface{}, error) {
            time.Sleep(2 * time.Second)
            return "slow task", nil
        },
        func(ctx context.Context) (interface{}, error) {
            time.Sleep(100 * time.Millisecond)
            return "fast task", nil
        },
        func(ctx context.Context) (interface{}, error) {
            return nil, fmt.Errorf("failed task")
        },
    )
    
    // 等待首个成功的任务
    value, err := executor.WaitFirstSuccess()
    if err != nil {
        fmt.Printf("All tasks failed: %v\n", err)
    } else {
        fmt.Printf("First success: %v\n", value)
    }
}
```

### 等待首个完成的任务

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/snail007/gmc/util/batch"
)

func main() {
    executor := gbatch.NewBatchExecutor()
    
    executor.AppendTask(
        func(ctx context.Context) (interface{}, error) {
            time.Sleep(2 * time.Second)
            return "task1", nil
        },
        func(ctx context.Context) (interface{}, error) {
            time.Sleep(100 * time.Millisecond)
            return "task2", nil
        },
    )
    
    // 等待首个完成的任务（无论成功还是失败）
    value, err := executor.WaitFirstDone()
    fmt.Printf("First done - Value: %v, Error: %v\n", value, err)
}
```

### 自定义 Panic 处理

```go
package main

import (
    "context"
    "fmt"
    "github.com/snail007/gmc/util/batch"
)

func main() {
    executor := gbatch.NewBatchExecutor()
    
    // 设置自定义 panic 处理器
    executor.SetPanicHandler(func(e interface{}) {
        fmt.Printf("Task panicked: %v\n", e)
        // 可以在这里记录日志、发送告警等
    })
    
    executor.AppendTask(
        func(ctx context.Context) (interface{}, error) {
            panic("something went wrong")
        },
    )
    
    results := executor.WaitAll()
    // panic 会被捕获并转换为 error
    fmt.Printf("Error: %v\n", results[0].Err())
}
```

## Context 使用指南

### Context 在不同场景下的行为

gbatch 包中，context 的行为在不同的等待策略下有所不同：

#### 1. WaitAll - 所有任务共享同一个 Context

在 `WaitAll()` 中，所有任务使用相同的 `rootCtx`，**不会被自动取消**。

```go
executor := gbatch.NewBatchExecutor()

executor.AppendTask(
    func(ctx context.Context) (interface{}, error) {
        // ctx 在所有任务完成前不会被取消
        // 可以安全地启动后台 goroutine
        go func() {
            <-ctx.Done() // 不会收到取消信号（除非外部取消）
            cleanup()
        }()
        
        // 执行任务
        return doWork(), nil
    },
)

results := executor.WaitAll()
// 所有任务正常完成，context 不会被取消
```

**使用建议**：
- 在 `WaitAll` 中，task 可以放心使用 context 进行超时控制
- 不需要担心其他 task 的完成状态
- Context 主要用于外部取消和超时控制

#### 2. WaitFirstSuccess / WaitFirstDone - 需要判断 IsFirstSuccess

在 `WaitFirstSuccess()` 和 `WaitFirstDone()` 中：
- **所有 task 的 context 都会被取消**（当第一个成功/完成时）
- 第一个完成的 task 可以通过 `IsFirstSuccess(ctx)` 判断
- 其他未完成的 task 收到取消信号后应该停止执行

```go
import "github.com/snail007/gmc/util/batch"

executor := gbatch.NewBatchExecutor()

executor.AppendTask(
    func(ctx context.Context) (interface{}, error) {
        // 启动后台 goroutine 监听取消
        done := make(chan struct{})
        go func() {
            defer close(done)
            select {
            case <-ctx.Done():
                // ⚠️ 关键：检查是否是第一个完成的
                if !gbatch.IsFirstSuccess(ctx) {
                    // 不是第一个完成的，执行清理逻辑
                    log.Println("Task cancelled, cleaning up...")
                    closeConnections()
                    releaseResources()
                } else {
                    // 是第一个完成的，不执行失败清理
                    log.Println("Task completed successfully")
                }
            case <-done:
                return
            }
        }()
        
        // 执行任务，支持取消
        for i := 0; i < 10; i++ {
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            default:
                // 执行工作
                time.Sleep(100 * time.Millisecond)
            }
        }
        
        return "result", nil
    },
)

// 当第一个任务成功/完成后：
// 1. 第一个完成的 task 的 IsFirstSuccess(ctx) 返回 true
// 2. 其他 task 的 IsFirstSuccess(ctx) 返回 false
// 3. 所有 task 的 ctx.Done() 都会收到信号
value, err := executor.WaitFirstSuccess()
```

**为什么所有 context 都要被取消？**
- 避免 goroutine 泄漏：如果第一个完成的 task 的 context 不被取消，其内部监听 `ctx.Done()` 的 goroutine 会永远阻塞
- 统一清理：所有 task 都能收到清理信号，只是清理逻辑不同

**IsFirstSuccess 的作用：**
- 返回 `true`：当前 task 是第一个成功/完成的，不应执行失败清理
- 返回 `false`：当前 task 不是第一个完成的，应该执行清理逻辑

### 完整示例：正确使用 Context

#### 示例 1：WaitAll - 简单使用

```go
executor := gbatch.NewBatchExecutor()

executor.AppendTask(
    func(ctx context.Context) (interface{}, error) {
        // WaitAll 中不需要担心 context 被取消
        return fetchData(ctx), nil
    },
    func(ctx context.Context) (interface{}, error) {
        return processData(ctx), nil
    },
)

results := executor.WaitAll()
```

#### 示例 2：WaitFirstSuccess - 处理资源清理

```go
executor := gbatch.NewBatchExecutor()

// 任务 1：尝试从主数据源获取
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    conn := openConnection("primary")
    defer func() {
        // 总是关闭连接
        conn.Close()
    }()
    
    // 后台监听取消，执行额外清理
    go func() {
        <-ctx.Done()
        if !gbatch.IsFirstSuccess(ctx) {
            // 不是第一个成功的，执行清理
            log.Println("Primary source cancelled, cleaning up...")
            cancelPendingRequests(conn)
        }
    }()
    
    return fetchFromPrimary(ctx, conn)
})

// 任务 2：尝试从备份数据源获取
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    conn := openConnection("backup")
    defer conn.Close()
    
    go func() {
        <-ctx.Done()
        if !gbatch.IsFirstSuccess(ctx) {
            log.Println("Backup source cancelled, cleaning up...")
            cancelPendingRequests(conn)
        }
    }()
    
    return fetchFromBackup(ctx, conn)
})

// 使用第一个成功的结果
data, err := executor.WaitFirstSuccess()
if err != nil {
    log.Fatal("All sources failed:", err)
}
fmt.Println("Got data:", data)
```

#### 示例 3：WaitFirstDone - 超时控制

```go
executor := gbatch.NewBatchExecutor()

// 任务 1：正常执行
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    go func() {
        <-ctx.Done()
        // 使用 IsFirstDone 判断
        if !gbatch.IsFirstDone(ctx) {
            log.Println("Task 1 cancelled, cleaning up...")
            cleanup()
        }
    }()
    return longRunningTask(ctx)
})

// 任务 2：超时控制
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    time.Sleep(5 * time.Second)
    return nil, errors.New("timeout")
})

// 谁先完成就用谁的结果
value, err := executor.WaitFirstDone()
```

### 常见模式

#### 模式 1：支持取消的长时间任务

```go
func(ctx context.Context) (interface{}, error) {
    result := make(chan interface{}, 1)
    errCh := make(chan error, 1)
    
    go func() {
        // 执行实际工作
        data, err := doWork()
        if err != nil {
            errCh <- err
            return
        }
        result <- data
    }()
    
    // 等待完成或取消
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case err := <-errCh:
        return nil, err
    case data := <-result:
        return data, nil
    }
}
```

#### 模式 2：批量操作支持取消

```go
func(ctx context.Context) (interface{}, error) {
    items := getItems()
    results := make([]Result, 0, len(items))
    
    for _, item := range items {
        // 每次循环检查取消
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        
        result := processItem(item)
        results = append(results, result)
    }
    
    return results, nil
}
```

#### 模式 3：资源清理

```go
func(ctx context.Context) (interface{}, error) {
    resource := acquireResource()
    
    // 监听取消信号进行清理
    cleanupDone := make(chan struct{})
    go func() {
        defer close(cleanupDone)
        <-ctx.Done()
        
        // 检查是否需要清理
        if !gbatch.IsFirstSuccess(ctx) {
            resource.Cleanup()
        }
    }()
    
    defer func() {
        resource.Release()
        <-cleanupDone // 等待清理完成
    }()
    
    return resource.Process()
}
```

### 关键要点

| 场景 | Context 行为 | 使用函数 | 说明 |
|-----|------------|---------|------|
| **WaitAll** | 不会被取消 | ❌ 不需要 | 所有 task 共享 rootCtx，正常完成不会取消 |
| **WaitFirstSuccess** | ✅ 会被取消 | `IsFirstSuccess(ctx)` | 判断是否是第一个成功的 |
| **WaitFirstDone** | ✅ 会被取消 | `IsFirstDone(ctx)` | 判断是否是第一个完成的 |

**记住**：
1. `WaitAll` 中放心使用 context，不会被自动取消
2. `WaitFirstSuccess` 中使用 `IsFirstSuccess(ctx)` 判断是否执行清理
3. `WaitFirstDone` 中使用 `IsFirstDone(ctx)` 判断是否执行清理
4. 所有 context 都会被取消是为了避免 goroutine 泄漏
5. 使用对应场景的判断函数，语义更清晰

## API 参考

### IsFirstSuccess 函数

```go
func IsFirstSuccess(ctx context.Context) bool
```

检查当前 task 是否是第一个成功的。**仅在 `WaitFirstSuccess` 中有效**。

**返回值：**
- `true`：当前 task 是第一个成功的，不应该执行失败清理逻辑
- `false`：当前 task 不是第一个成功的，应该执行清理逻辑

**使用场景：**

在 `WaitFirstSuccess` 中，当第一个任务成功后，所有任务的 context 都会被取消。使用此函数可以区分：
- 第一个成功的 task：虽然收到 `ctx.Done()` 信号，但不应执行失败清理
- 其他被取消的 task：应该执行清理逻辑

**示例：**

```go
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    go func() {
        <-ctx.Done()
        // 检查是否是第一个成功的
        if !gbatch.IsFirstSuccess(ctx) {
            // 不是第一个成功的，执行清理
            cleanupResources()
        }
    }()
    
    return doWork(), nil
})

// 使用 WaitFirstSuccess
result, err := executor.WaitFirstSuccess()
```

**注意：**
- 在 `WaitAll` 中，此函数始终返回 `false`
- 在 `WaitFirstDone` 中，应该使用 `IsFirstDone` 而不是此函数
- 只有在 context 收到取消信号后才需要判断

### IsFirstDone 函数

```go
func IsFirstDone(ctx context.Context) bool
```

检查当前 task 是否是第一个完成的（无论成功或失败）。**仅在 `WaitFirstDone` 中有效**。

**返回值：**
- `true`：当前 task 是第一个完成的，不应该执行失败清理逻辑
- `false`：当前 task 不是第一个完成的，应该执行清理逻辑

**使用场景：**

在 `WaitFirstDone` 中，当第一个任务完成后，所有任务的 context 都会被取消。使用此函数可以区分：
- 第一个完成的 task：虽然收到 `ctx.Done()` 信号，但不应执行失败清理
- 其他被取消的 task：应该执行清理逻辑

**示例：**

```go
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    go func() {
        <-ctx.Done()
        // 检查是否是第一个完成的
        if !gbatch.IsFirstDone(ctx) {
            // 不是第一个完成的，执行清理
            cleanupResources()
        }
    }()
    
    return doWork(), nil
})

// 使用 WaitFirstDone
result, err := executor.WaitFirstDone()
```

**注意：**
- 在 `WaitAll` 中，此函数始终返回 `false`
- 在 `WaitFirstSuccess` 中，应该使用 `IsFirstSuccess` 而不是此函数
- 只有在 context 收到取消信号后才需要判断

### Executor 类型

批量任务执行器。

#### 创建执行器

```go
func NewBatchExecutor() *Executor
```

创建一个新的批量任务执行器，默认使用 10 个工作线程。

#### 方法

##### SetWorkers

```go
func (s *Executor) SetWorkers(workersCnt int)
```

设置工作线程数量（goroutine 池大小）。

**参数：**
- `workersCnt`：工作线程数量

**示例：**
```go
executor.SetWorkers(20) // 使用 20 个工作线程
```

##### AppendTask

```go
func (s *Executor) AppendTask(tasks ...task)
```

添加一个或多个任务到执行器。

**参数：**
- `tasks`：任务函数，签名为 `func(ctx context.Context) (value interface{}, err error)`

**示例：**
```go
executor.AppendTask(
    func(ctx context.Context) (interface{}, error) {
        return "result", nil
    },
)
```

##### SetPanicHandler

```go
func (s *Executor) SetPanicHandler(panicHandler func(e interface{}))
```

设置自定义的 panic 处理函数。

**参数：**
- `panicHandler`：panic 处理函数

##### WaitAll

```go
func (s *Executor) WaitAll() []taskResult
```

等待所有任务完成，返回所有任务的结果。

**返回值：**
- `[]taskResult`：所有任务的结果数组

**示例：**
```go
results := executor.WaitAll()
for _, result := range results {
    if result.Err() != nil {
        // 处理错误
    } else {
        // 使用结果值 result.Value()
    }
}
```

##### WaitFirstSuccess

```go
func (s *Executor) WaitFirstSuccess() (value interface{}, err error)
```

等待首个成功的任务。如果所有任务都失败，返回最后一个任务的错误。

**返回值：**
- `value`：首个成功任务的返回值
- `err`：如果所有任务都失败，返回最后一个错误

##### WaitFirstDone

```go
func (s *Executor) WaitFirstDone() (value interface{}, err error)
```

等待首个完成的任务（无论成功还是失败）。

**返回值：**
- `value`：首个完成任务的返回值
- `err`：任务的错误（如果有）

### taskResult 类型

任务执行结果。

#### 方法

##### Value

```go
func (t taskResult) Value() interface{}
```

获取任务的返回值。

##### Err

```go
func (t taskResult) Err() error
```

获取任务的错误。

## 任务函数

任务函数的签名为：

```go
func(ctx context.Context) (value interface{}, err error)
```

**参数：**
- `ctx`：上下文对象，可用于取消和超时控制

**返回值：**
- `value`：任务的返回值，可以是任意类型
- `err`：任务执行过程中的错误

## 使用场景

1. **并行 HTTP 请求**：同时请求多个 API 端点
2. **数据库批量操作**：并行执行多个数据库查询
3. **文件批量处理**：并发处理多个文件
4. **负载均衡**：向多个服务器发送请求，使用首个响应
5. **容错重试**：尝试多个数据源，使用首个成功的结果

## 注意事项

1. **工作线程数**：合理设置工作线程数，避免创建过多 goroutine
2. **Context 取消行为**：
   - `WaitAll`：所有任务共享 rootCtx，正常完成不会被取消
   - `WaitFirstSuccess`：所有任务的 context 都会被取消（包括第一个成功的），使用 `IsFirstSuccess(ctx)` 判断
   - `WaitFirstDone`：所有任务的 context 都会被取消（包括第一个完成的），使用 `IsFirstDone(ctx)` 判断
3. **使用对应的判断函数**：
   - `WaitFirstSuccess` → 使用 `IsFirstSuccess(ctx)` 
   - `WaitFirstDone` → 使用 `IsFirstDone(ctx)`
   - 不要混用，否则会得到错误的结果
4. **Panic 处理**：任务中的 panic 会被捕获并转换为 error
5. **结果顺序**：`WaitAll` 返回的结果顺序与任务添加顺序可能不同
6. **资源清理**：确保在任务函数中正确处理资源清理，避免资源泄漏
7. **Goroutine 泄漏**：在 `WaitFirstSuccess/Done` 中，所有 context 都会被取消以避免 goroutine 泄漏

## 最佳实践

### 1. 合理设置并发数

```go
// 根据任务类型设置合适的并发数
executor := gbatch.NewBatchExecutor()
if isCPUBound {
    executor.SetWorkers(runtime.NumCPU())
} else {
    executor.SetWorkers(100) // IO 密集型可以设置更多
}
```

### 2. 正确使用 Context

#### WaitAll - 无需特殊处理

```go
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    // WaitAll 中 context 不会被自动取消
    // 可以正常使用 context 进行超时控制
    return doWork(ctx), nil
})
```

#### WaitFirstSuccess/Done - 必须使用 IsFirstSuccess

```go
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    // 启动后台清理 goroutine
    go func() {
        <-ctx.Done()
        // ⚠️ 关键：判断是否是第一个完成的
        if !gbatch.IsFirstSuccess(ctx) {
            // 不是第一个完成的，执行清理
            cleanupResources()
        }
    }()
    
    // 支持取消的主逻辑
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        return doWork(), nil
    }
})
```

### 3. 错误处理

```go
results := executor.WaitAll()
var errors []error
var values []interface{}

for _, result := range results {
    if result.Err() != nil {
        errors = append(errors, result.Err())
    } else {
        values = append(values, result.Value())
    }
}
```

### 4. 资源管理

```go
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    // 获取资源
    conn := getConnection()
    
    // 确保资源被释放
    defer conn.Close()
    
    // 监听取消执行额外清理
    cleanupDone := make(chan struct{})
    go func() {
        defer close(cleanupDone)
        <-ctx.Done()
        if !gbatch.IsFirstSuccess(ctx) {
            conn.CancelPendingRequests()
        }
    }()
    defer func() { <-cleanupDone }()
    
    return conn.DoWork(ctx)
})
```

## 性能考虑

- 使用 goroutine 池避免频繁创建和销毁 goroutine
- 默认工作线程数为 10，可根据实际情况调整
- 对于 IO 密集型任务，可以设置较大的工作线程数
- 对于 CPU 密集型任务，建议设置为 CPU 核心数

## 依赖

- `github.com/snail007/gmc/module/error`：错误处理
- `github.com/snail007/gmc/util/gpool`：goroutine 池
- `github.com/snail007/gmc/util/list`：线程安全的列表

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
