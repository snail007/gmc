# WaitFirst 使用示例

## 问题说明

在使用 `WaitFirstSuccess` 或 `WaitFirstDone` 时，task 内部可能会监听 `ctx.Done()` 来执行清理操作。当第一个 task 成功完成后：

1. **所有 task 的 context 都会被取消**（包括第一个成功的）
2. 所有 task 内部监听 `ctx.Done()` 的 goroutine 都会收到取消信号
3. 但是我们需要区分：
   - **第一个成功的 task**：不应该执行失败清理逻辑（虽然收到取消信号）
   - **其他被取消的 task**：应该执行失败清理逻辑

## 解决方案

使用 `IsFirstSuccess(ctx)` 或 `IsFirstDone(ctx)` 函数来判断当前 task 是否是第一个成功/完成的。

**关键点**：
- 所有 task 的 context 都从 `rootCtx` 派生
- 第一个完成时，调用 `rootCancel()` 会自动取消所有派生的 context
- 在 `WaitFirstSuccess` 中使用 `IsFirstSuccess(ctx)` 判断
- 在 `WaitFirstDone` 中使用 `IsFirstDone(ctx)` 判断
- 第一个完成的 task 虽然收到取消信号，但判断函数返回 `true`，不执行清理逻辑
- 其他 task 判断函数返回 `false`，执行清理逻辑

**实现简洁**：
- 不需要管理多个 cancel 函数
- 直接调用 `rootCancel()` 即可取消所有 task
- 通过 context.Value 传递状态标识
- 提供两个独立函数，语义更清晰：`IsFirstSuccess` 和 `IsFirstDone`

## 使用示例

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
	
	// Task 1: 模拟一个需要 2 秒的任务
	executor.AppendTask(func(ctx context.Context) (interface{}, error) {
		// 启动后台 goroutine 监听取消
		done := make(chan struct{})
		go func() {
			defer close(done)
			select {
			case <-ctx.Done():
				// 检查是否是第一个完成的
				if !gbatch.IsFirstSuccess(ctx) {
					fmt.Println("Task 1: 被取消，执行清理操作...")
					// 执行清理操作，如关闭连接、释放资源等
				} else {
					fmt.Println("Task 1: 成功完成，不执行清理操作")
				}
			case <-done:
				return
			}
		}()
		
		time.Sleep(2 * time.Second)
		fmt.Println("Task 1: 完成")
		return "Task 1 result", nil
	})
	
	// Task 2: 模拟一个需要 1 秒的任务（会先完成）
	executor.AppendTask(func(ctx context.Context) (interface{}, error) {
		// 启动后台 goroutine 监听取消
		done := make(chan struct{})
		go func() {
			defer close(done)
			select {
			case <-ctx.Done():
				// 检查是否是第一个完成的
				if !gbatch.IsFirstSuccess(ctx) {
					fmt.Println("Task 2: 被取消，执行清理操作...")
					// 执行清理操作
				} else {
					fmt.Println("Task 2: 成功完成，不执行清理操作")
				}
			case <-done:
				return
			}
		}()
		
		time.Sleep(1 * time.Second)
		fmt.Println("Task 2: 完成")
		return "Task 2 result", nil
	})
	
	// 等待第一个成功的任务
	result, err := executor.WaitFirstSuccess()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("第一个完成的结果: %v\n", result)
	
	// 等待一段时间，让清理操作完成
	time.Sleep(2 * time.Second)
}
```

## 输出示例

```
Task 2: 完成
第一个完成的结果: Task 2 result
Task 2: 成功完成，不执行清理操作
Task 1: 被取消，执行清理操作...
```

**说明**：
- Task 2 先完成，成为第一个成功的 task
- Task 2 的 context 也会被取消，但 `IsFirstSuccess` 返回 `true`，所以不执行清理
- Task 1 的 context 被取消，`IsFirstSuccess` 返回 `false`，执行清理逻辑

## 注意事项

1. `IsFirstSuccess` 只在 `WaitFirstSuccess` 和 `WaitFirstDone` 中有效
2. 在 `WaitAll` 中，所有 task 使用相同的 context，不会被取消
3. **所有 task 的 context 都会被取消**，包括第一个完成的，这样可以避免 goroutine 泄漏
4. 第一个完成的 task 通过 `IsFirstSuccess` 返回 `true` 来区分
5. 建议在监听 `ctx.Done()` 时总是检查 `IsFirstSuccess`，避免误执行清理逻辑
