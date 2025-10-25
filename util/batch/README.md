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

### 使用 Context 取消任务

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
            for i := 0; i < 10; i++ {
                // 检查是否被取消
                select {
                case <-ctx.Done():
                    return nil, ctx.Err()
                default:
                    time.Sleep(100 * time.Millisecond)
                    fmt.Printf("Working... %d\n", i)
                }
            }
            return "completed", nil
        },
    )
    
    // WaitFirstSuccess 会在首个任务成功后取消其他任务
    value, err := executor.WaitFirstSuccess()
    fmt.Printf("Result: %v, Error: %v\n", value, err)
}
```

## API 参考

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
2. **任务取消**：使用 `WaitFirstSuccess` 或 `WaitFirstDone` 时，其他任务会收到取消信号
3. **Panic 处理**：任务中的 panic 会被捕获并转换为 error
4. **结果顺序**：`WaitAll` 返回的结果顺序与任务添加顺序可能不同
5. **资源清理**：确保在任务函数中正确处理资源清理

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

### 2. 检查 Context 取消

```go
executor.AppendTask(func(ctx context.Context) (interface{}, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // 执行任务
    }
    return result, nil
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
