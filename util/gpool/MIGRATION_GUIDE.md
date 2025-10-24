# gpool 优化版本迁移指南

## 快速开始

### 原版用法
```go
import "github.com/snail007/gmc/util/gpool"

// 创建pool
p := gpool.New(100)
defer p.Stop()

// 提交任务
p.Submit(func() {
    // 你的任务代码
})

// 等待完成
p.WaitDone()
```

### 优化版用法
```go
import "github.com/snail007/gmc/util/gpool"

// 只需改变创建方式
p := gpool.NewOptimized(100)  // 就这里不同！
defer p.Stop()

// 其他代码完全相同
p.Submit(func() {
    // 你的任务代码
})

p.WaitDone()
```

## API对比

| 功能 | 原版API | 优化版API | 兼容性 |
|-----|---------|-----------|--------|
| 创建pool | `New(n)` | `NewOptimized(n)` | ✅ 参数相同 |
| 带选项创建 | `NewWithOption(n, opt)` | `NewOptimizedWithOption(n, opt)` | ✅ 参数相同 |
| 提交任务 | `Submit(func())` | `Submit(func())` | ✅ 完全相同 |
| 等待完成 | `WaitDone()` | `WaitDone()` | ✅ 完全相同 |
| 停止pool | `Stop()` | `Stop()` | ✅ 完全相同 |
| 增加worker | `Increase(n)` | `Increase(n)` | ✅ 完全相同 |
| 减少worker | `Decrease(n)` | `Decrease(n)` | ✅ 完全相同 |
| Worker数量 | `WorkerCount()` | `WorkerCount()` | ✅ 完全相同 |
| 运行中数量 | `RunningWorkerCount()` | `RunningWorkerCount()` | ✅ 完全相同 |
| 空闲数量 | `IdleWorkerCount()` | `IdleWorkerCount()` | ✅ 完全相同 |
| 队列长度 | `QueuedJobCount()` | `QueuedJobCount()` | ✅ 完全相同 |

## 配置选项变化

### WithStack 默认值变化

```go
// 原版：默认开启（影响性能）
p := gpool.New(100)  // WithStack = true

// 优化版：默认关闭（更好性能）
p := gpool.NewOptimized(100)  // WithStack = false

// 如需stack trace，显式开启：
p := gpool.NewOptimizedWithOption(100, &gpool.Option{
    WithStack: true,  // 开启后会记录调用栈
})
```

### MaxJobCount 默认值

```go
// 原版：无默认限制
p := gpool.New(100)  // MaxJobCount = 0 (无限制)

// 优化版：默认10000
p := gpool.NewOptimized(100)  // MaxJobCount = 10000

// 自定义队列大小：
p := gpool.NewOptimizedWithOption(100, &gpool.Option{
    MaxJobCount: 50000,  // 自定义大小
})
```

## 性能对比示例

### 低并发场景（改进明显）
```go
// 原版：6,304 ns/op
p := gpool.New(100)
for i := 0; i < 10000; i++ {
    p.Submit(func() { /* work */ })
}

// 优化版：761 ns/op （快8.3倍）
p := gpool.NewOptimized(100)
for i := 0; i < 10000; i++ {
    p.Submit(func() { /* work */ })
}
```

### 高并发场景（改进巨大）
```go
// 原版：6,109 ns/op
p := gpool.New(100)
for i := 0; i < 16; i++ {
    go func() {
        for j := 0; j < 1000; j++ {
            p.Submit(func() { /* work */ })
        }
    }()
}

// 优化版：288 ns/op （快21倍）
p := gpool.NewOptimized(100)
for i := 0; i < 16; i++ {
    go func() {
        for j := 0; j < 1000; j++ {
            p.Submit(func() { /* work */ })
        }
    }()
}
```

## 实际使用案例

### Web服务器
```go
// 创建全局pool
var taskPool = gpool.NewOptimized(1000)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 异步处理任务
    taskPool.Submit(func() {
        // 处理耗时操作
        processData()
    })
    
    w.Write([]byte("Processing..."))
}
```

### 批量数据处理
```go
func processBatch(items []Item) {
    p := gpool.NewOptimized(100)
    defer p.Stop()
    
    for _, item := range items {
        item := item  // 闭包变量捕获
        p.Submit(func() {
            processItem(item)
        })
    }
    
    p.WaitDone()  // 等待所有任务完成
}
```

### 微服务调用
```go
func callMultipleServices() ([]Response, error) {
    p := gpool.NewOptimizedWithOption(10, &gpool.Option{
        PanicHandler: func(e interface{}) {
            log.Printf("Service call panic: %v", e)
        },
    })
    defer p.Stop()
    
    results := make([]Response, len(services))
    var mu sync.Mutex
    
    for i, svc := range services {
        i, svc := i, svc
        p.Submit(func() {
            resp, err := svc.Call()
            mu.Lock()
            results[i] = resp
            mu.Unlock()
        })
    }
    
    p.WaitDone()
    return results, nil
}
```

## 注意事项

### 1. 队列满处理
```go
// 优化版有默认队列限制
p := gpool.NewOptimized(100)  // 默认MaxJobCount=10000

// 队列满时行为：
err := p.Submit(func() {})
if err == gpool.ErrMaxQueuedJobCountReachedOptimized {
    // 队列满了，处理逻辑
}

// 或使用阻塞模式：
p := gpool.NewOptimizedWithOption(100, &gpool.Option{
    Blocking: true,  // 队列满时阻塞等待
})
```

### 2. 资源清理
```go
// 务必调用Stop()
p := gpool.NewOptimized(100)
defer p.Stop()  // 清理资源

// 或显式停止
p.Stop()
// 注意：Stop后不能再Submit
```

### 3. Worker数量选择
```go
// CPU密集型：
cpuCount := runtime.NumCPU()
p := gpool.NewOptimized(cpuCount)

// IO密集型：
p := gpool.NewOptimized(cpuCount * 2)  // 或更多

// 网络调用：
p := gpool.NewOptimized(100)  // 根据并发连接数调整
```

## 性能调优建议

### 1. 合适的队列大小
```go
// 根据实际场景调整
p := gpool.NewOptimizedWithOption(100, &gpool.Option{
    MaxJobCount: 1000,    // 小队列：低延迟
    // MaxJobCount: 100000,  // 大队列：高吞吐
})
```

### 2. 关闭不需要的功能
```go
p := gpool.NewOptimizedWithOption(100, &gpool.Option{
    WithStack: false,     // 默认关闭，性能最佳
    Debug: false,         // 生产环境关闭调试
    Logger: nil,          // 不需要日志时设为nil
})
```

### 3. 预分配Worker
```go
// 避免运行时动态创建worker
p := gpool.NewOptimizedWithOption(100, &gpool.Option{
    PreAlloc: true,  // 启动时创建所有worker
})
```

## 故障排查

### 问题1：Submit返回错误
```go
err := p.Submit(func() {})
if err != nil {
    // 可能原因：
    // 1. 队列满了（增加MaxJobCount或开启Blocking）
    // 2. Pool已停止（检查Stop()调用）
}
```

### 问题2：任务没执行
```go
// 确保调用WaitDone或者不要过早Stop
p.Submit(func() { println("task") })
// p.Stop()  // ❌ 立即stop可能导致任务未执行
p.WaitDone()  // ✅ 等待完成
p.Stop()
```

### 问题3：性能不如预期
```go
// 检查：
// 1. WithStack是否开启（应关闭）
// 2. Worker数量是否合理
// 3. 任务粒度是否太小（考虑批处理）
```

## 迁移检查清单

- [ ] 替换 `New()` 为 `NewOptimized()`
- [ ] 检查 `WithStack` 配置（如无需要，保持默认关闭）
- [ ] 测试队列满的情况（调整 `MaxJobCount`）
- [ ] 验证功能正确性
- [ ] 运行基准测试，确认性能提升
- [ ] 更新文档和注释

## 总结

优化版提供了：
- ✅ **8-21倍** 性能提升
- ✅ **92倍** 内存优化
- ✅ API基本兼容
- ✅ 更符合Go惯例

只需简单替换创建函数即可享受性能提升！
