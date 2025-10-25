# gfile 包

## 简介

gfile 包提供了丰富的文件和目录操作工具函数，简化文件系统操作。

## 功能特性

- **文件检查**：检查文件/目录是否存在、是否为文件/目录/链接
- **文件信息**：获取文件大小、修改时间
- **路径处理**：支持 `~` 和 `$HOME` 展开
- **文件读写**：便捷的文件读写函数
- **文件复制**：支持文件拷贝

## 安装

```bash
go get github.com/snail007/gmc/util/file
```

## 快速开始

### 文件检查

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/file"
)

func main() {
    // 检查文件是否存在
    if gfile.Exists("/tmp/test.txt") {
        fmt.Println("文件存在")
    }
    
    // 检查是否为文件
    if gfile.IsFile("/tmp/test.txt") {
        fmt.Println("是文件")
    }
    
    // 检查是否为目录
    if gfile.IsDir("/tmp") {
        fmt.Println("是目录")
    }
    
    // 检查是否为符号链接
    if gfile.IsLink("/tmp/link") {
        fmt.Println("是符号链接")
    }
    
    // 检查目录是否为空
    if gfile.IsEmptyDir("/tmp/empty") {
        fmt.Println("目录为空")
    }
}
```

### 文件读写

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/file"
)

func main() {
    // 读取文件为字符串
    content := gfile.ReadAll("/tmp/test.txt")
    fmt.Println(content)
    
    // 读取文件为字节数组
    data := gfile.Bytes("/tmp/test.txt")
    fmt.Printf("读取了 %d 字节\n", len(data))
    
    // 写入字符串到文件（覆盖）
    err := gfile.WriteString("/tmp/output.txt", "Hello, World!", false)
    if err != nil {
        panic(err)
    }
    
    // 追加字符串到文件
    err = gfile.WriteString("/tmp/output.txt", "\nNew Line", true)
    if err != nil {
        panic(err)
    }
    
    // 写入字节数组
    err = gfile.Write("/tmp/data.bin", []byte{0x01, 0x02, 0x03}, false)
    if err != nil {
        panic(err)
    }
}
```

### 路径处理

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/file"
)

func main() {
    // 获取家目录
    home := gfile.HomeDir()
    fmt.Println("Home:", home)
    
    // 展开路径中的 ~ 和 $HOME
    path := gfile.Abs("~/Documents/file.txt")
    fmt.Println("绝对路径:", path)
    
    // $HOME 也会被展开
    path = gfile.Abs("$HOME/Documents/file.txt")
    fmt.Println("绝对路径:", path)
    
    // 获取文件名（不含扩展名）
    name := gfile.FileName("/path/to/file.txt")
    fmt.Println("文件名:", name) // "file"
    
    // 获取文件名（含扩展名）
    basename := gfile.BaseName("/path/to/file.txt")
    fmt.Println("基础名:", basename) // "file.txt"
}
```

### 文件信息

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/file"
)

func main() {
    // 获取文件大小
    size, err := gfile.FileSize("/tmp/test.txt")
    if err != nil {
        panic(err)
    }
    fmt.Printf("文件大小: %d 字节\n", size)
    
    // 获取文件修改时间
    mtime, err := gfile.FileMTime("/tmp/test.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println("修改时间:", mtime)
}
```

### 文件复制

```go
package main

import (
    "github.com/snail007/gmc/util/file"
)

func main() {
    // 复制文件，如果目标目录不存在则创建
    err := gfile.Copy("/tmp/source.txt", "/tmp/backup/dest.txt", true)
    if err != nil {
        panic(err)
    }
}
```

### 符号链接处理

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/file"
)

func main() {
    linkPath := "/tmp/mylink"
    
    if gfile.IsLink(linkPath) {
        // 解析符号链接
        target, err := gfile.ResolveLink(linkPath)
        if err != nil {
            panic(err)
        }
        fmt.Println("链接指向:", target)
    }
}
```

## API 参考

### 文件检查

```go
func Exists(path string) bool
func IsFile(file string) bool
func IsDir(file string) bool
func IsLink(path string) bool
func IsEmptyDir(path string) bool
```

### 文件信息

```go
func FileSize(file string) (int64, error)
func FileMTime(file string) (time.Time, error)
func HomeDir() string
```

### 路径处理

```go
func Abs(p string) string
func FileName(file string) string
func BaseName(file string) string
```

### 文件读取

```go
func ReadAll(file string) string
func Bytes(file string) []byte
```

### 文件写入

```go
func Write(file string, data []byte, append bool) error
func WriteString(file string, data string, append bool) error
```

### 文件操作

```go
func Copy(srcFile, dstFile string, mkdir bool) error
func ResolveLink(path string) (string, error)
```

## 使用场景

1. **配置文件处理**：读写应用配置
2. **日志文件**：追加日志内容
3. **文件备份**：复制文件到备份目录
4. **路径处理**：处理用户输入的路径（支持 `~`）
5. **文件监控**：检查文件变化

## 注意事项

1. `ReadAll` 和 `Bytes` 失败时返回空值而不是错误
2. 文件权限默认为 0755
3. `Abs` 函数出错时返回原始路径
4. `Copy` 可以自动创建目标目录

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
