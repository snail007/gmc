# gexec 包

## 简介

gexec 包提供了一个功能强大的 Linux 命令执行器，支持异步执行、进程分离、超时控制、环境变量设置等高级特性。

**注意：此包仅支持 Linux 系统，不支持 Windows。**

## 功能特性

- **命令执行**：执行 Shell 命令和脚本
- **异步执行**：支持异步执行和回调
- **进程分离（Detach）**：子进程可以在父进程退出后继续运行
- **超时控制**：设置命令执行超时时间
- **环境变量**：自定义命令执行环境
- **工作目录**：指定命令执行目录
- **输出捕获**：捕获命令输出或重定向到自定义 Writer
- **严格模式**：支持 `set -e` 模式，命令失败立即退出
- **生命周期钩子**：支持执行前、执行后、退出后的钩子函数
- **终端类型**：可配置终端类型（TERM 环境变量）

## 安装

```bash
go get github.com/snail007/gmc/util/exec
```

## 快速开始

### 基本命令执行

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/exec"
)

func main() {
    cmd := gexec.NewCommand("echo 'Hello, World!'")
    
    output, err := cmd.Exec()
    if err != nil {
        panic(err)
    }
    
    fmt.Println(output)
}
```

### 带参数的命令

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/exec"
)

func main() {
    cmd := gexec.NewCommand("echo")
    cmd.Args("Hello", "World")
    
    output, err := cmd.Exec()
    if err != nil {
        panic(err)
    }
    
    fmt.Println(output)
}
```

### 异步执行

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/exec"
)

func main() {
    cmd := gexec.NewCommand("sleep 5 && echo 'Done'")
    
    cmd.Async(true).AsyncCallback(func(c *gexec.Command, output string, err error) {
        if err != nil {
            fmt.Println("Error:", err)
        } else {
            fmt.Println("Output:", output)
        }
    })
    
    output, err := cmd.Exec()
    fmt.Println("Command started asynchronously")
    
    // 主程序继续执行
    time.Sleep(6 * time.Second)
}
```

### 进程分离（Detach）

```go
package main

import (
    "github.com/snail007/gmc/util/exec"
)

func main() {
    // 启动一个守护进程，即使父进程退出也继续运行
    cmd := gexec.NewCommand("nohup myserver > /var/log/myserver.log 2>&1")
    cmd.Detach(true)
    
    _, err := cmd.Exec()
    if err != nil {
        panic(err)
    }
    
    // 父进程可以退出，子进程继续运行
}
```

### 超时控制

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/exec"
)

func main() {
    cmd := gexec.NewCommand("sleep 10")
    cmd.Timeout(2 * time.Second)
    
    _, err := cmd.Exec()
    if err != nil {
        fmt.Println("Command timeout:", err)
    }
}
```

### 自定义环境变量

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/exec"
)

func main() {
    cmd := gexec.NewCommand("echo $MY_VAR")
    cmd.Env(map[string]string{
        "MY_VAR": "Hello from custom env",
    })
    
    output, err := cmd.Exec()
    if err != nil {
        panic(err)
    }
    
    fmt.Println(output)
}
```

### 指定工作目录

```go
package main

import (
    "github.com/snail007/gmc/util/exec"
)

func main() {
    cmd := gexec.NewCommand("pwd")
    cmd.WorkDir("/tmp")
    
    output, _ := cmd.Exec()
    println(output) // /tmp
}
```

### 输出重定向

```go
package main

import (
    "os"
    "github.com/snail007/gmc/util/exec"
)

func main() {
    file, _ := os.Create("/tmp/output.log")
    defer file.Close()
    
    cmd := gexec.NewCommand("echo 'Logging to file'")
    cmd.Output(file)
    
    cmd.Exec()
}
```

## API 参考

### NewCommand

```go
func NewCommand(cmd string) *Command
```

创建一个新的命令执行器。

**参数：**
- `cmd`：要执行的 Shell 命令

**返回值：**
- `*Command`：命令执行器实例

**注意：**
- 仅支持 Linux 系统
- Windows 系统会 panic

### Command 方法

#### 基本配置

##### Args

```go
func (s *Command) Args(args ...string) *Command
```

设置命令参数。

##### WorkDir

```go
func (s *Command) WorkDir(workDir string) *Command
```

设置命令执行的工作目录。

##### Env

```go
func (s *Command) Env(env map[string]string) *Command
```

设置环境变量。会与系统环境变量合并。

##### TermType

```go
func (s *Command) TermType(termType TermType) *Command
```

设置终端类型（TERM 环境变量）。

**终端类型常量：**
- `TermXterm`："xterm"（默认）
- `TermXterm256Color`："xterm-256color"
- `TermVt100`："vt100"
- `TermRxvt`："rxvt"
- `TermGnome`："gnome-terminal"
- `TermKonsole`："konsole"
- `TermTmux`："tmux"
- `TermScreen`："screen"
- `TermAnsi`："ansi"
- `TermNull`：空（不设置）

#### 执行控制

##### StrictMode

```go
func (s *Command) StrictMode(strictMode bool) *Command
```

设置严格模式。启用后相当于在脚本开头添加 `set -e`。

**默认：** true

##### Timeout

```go
func (s *Command) Timeout(timeout time.Duration) *Command
```

设置命令执行超时时间。

##### Async

```go
func (s *Command) Async(async bool) *Command
```

设置是否异步执行。异步模式会创建 goroutine 等待命令完成。

**注意：** 异步模式下，父进程退出时子进程也会退出。

##### Detach

```go
func (s *Command) Detach(detach bool) *Command
```

设置进程分离模式。启用后子进程会脱离父进程，父进程退出后子进程继续运行。

**要求：** Go 1.20+

##### AsyncCallback

```go
func (s *Command) AsyncCallback(callback func(cmd *Command, output string, err error)) *Command
```

设置异步执行完成时的回调函数。

#### 输出控制

##### Output

```go
func (s *Command) Output(w io.Writer) *Command
```

设置输出重定向目标。不设置时输出会被捕获并作为返回值。

##### Log

```go
func (s *Command) Log(log gcore.Logger) *Command
```

设置日志记录器，用于记录错误信息。

#### 生命周期钩子

##### BeforeExec

```go
func (s *Command) BeforeExec(f func(command *Command, cmd *exec.Cmd)) *Command
```

设置命令执行前的钩子函数。

##### AfterExec

```go
func (s *Command) AfterExec(f func(command *Command, cmd *exec.Cmd, err error)) *Command
```

设置命令启动后的钩子函数（Start 之后立即调用）。

##### AfterExited

```go
func (s *Command) AfterExited(f func(command *Command, cmd *exec.Cmd, err error)) *Command
```

设置命令退出后的钩子函数（Wait 之后调用）。

#### 执行和控制

##### Exec

```go
func (s *Command) Exec() (output string, err error)
```

执行命令。同步执行时等待命令完成，异步执行时立即返回。

**返回值：**
- `output`：命令输出（仅同步模式且未设置 Output 时有效）
- `err`：执行错误

##### ExecAsync

```go
func (s *Command) ExecAsync() error
```

异步执行命令，不等待命令完成。

**注意：** 不能与 `Async(true)` 同时使用。

##### Kill

```go
func (s *Command) Kill()
```

杀死正在执行的命令进程。

##### Cmd

```go
func (s *Command) Cmd() *exec.Cmd
```

获取底层的 `*exec.Cmd` 对象。

## 使用场景

### 1. 系统运维脚本

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/exec"
)

func main() {
    script := `
    df -h
    free -m
    uptime
    `
    
    cmd := gexec.NewCommand(script)
    output, err := cmd.Exec()
    if err != nil {
        panic(err)
    }
    
    fmt.Println(output)
}
```

### 2. 启动守护进程

```go
package main

import (
    "github.com/snail007/gmc/util/exec"
)

func startDaemon() error {
    cmd := gexec.NewCommand(`
    cd /app
    nohup ./myapp > /var/log/myapp.log 2>&1 &
    `)
    cmd.Detach(true)
    
    _, err := cmd.Exec()
    return err
}
```

### 3. 批量操作

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/exec"
)

func backupDatabases() {
    databases := []string{"db1", "db2", "db3"}
    
    for _, db := range databases {
        cmd := gexec.NewCommand(fmt.Sprintf(
            "mysqldump -u root %s > /backup/%s.sql", db, db))
        
        if _, err := cmd.Exec(); err != nil {
            fmt.Printf("Backup %s failed: %v\n", db, err)
        }
    }
}
```

### 4. 定时任务

```go
package main

import (
    "time"
    "github.com/snail007/gmc/util/exec"
)

func runCron() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        cmd := gexec.NewCommand("./cleanup.sh")
        cmd.Async(true).AsyncCallback(func(c *gexec.Command, output string, err error) {
            if err != nil {
                // 记录错误
            }
        })
        cmd.Exec()
    }
}
```

## 注意事项

1. **仅支持 Linux**：Windows 系统会 panic
2. **Shell 脚本**：命令会在 bash 中执行
3. **临时文件**：命令会被写入临时文件 `/tmp/tmp_*.sh`
4. **严格模式**：默认启用 `set -e`，任何错误都会导致脚本退出
5. **进程管理**：Async 和 ExecAsync 不能同时使用
6. **Detach 要求**：进程分离需要 Go 1.20+
7. **输出捕获**：设置 Output 后不会返回输出字符串

## 最佳实践

### 1. 错误处理

```go
cmd := gexec.NewCommand("risky-command")
output, err := cmd.Exec()
if err != nil {
    log.Printf("Command failed: %v, output: %s", err, output)
    return
}
```

### 2. 超时设置

```go
cmd := gexec.NewCommand("long-running-task")
cmd.Timeout(5 * time.Minute)
```

### 3. 日志记录

```go
cmd := gexec.NewCommand("important-command")
cmd.Log(logger)
```

### 4. 资源清理

```go
cmd := gexec.NewCommand("streaming-command")
defer cmd.Kill()
```

## 依赖

- `github.com/snail007/gmc/core`：核心接口
- `github.com/snail007/gmc/module/error`：错误处理
- `github.com/snail007/gmc/util/file`：文件操作
- `github.com/snail007/gmc/util/rand`：随机数生成

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
