# gbytes 包

## 简介

gbytes 包提供了一系列字节相关的实用工具，包括循环缓冲区、字节池、字节大小解析和格式化、以及便捷的字节写入器。

## 功能特性

- **循环缓冲区（CircularBuffer）**：固定大小的环形缓冲区，支持多读者模式
- **字节池（Pool）**：基于 sync.Pool 的字节切片池，减少内存分配
- **字节大小处理（ByteSize）**：解析和格式化人类可读的字节大小（如 "1.5GB"）
- **IO 写入器（IOWriter）**：便捷的格式化写入工具
- **字节构建器（BytesBuilder）**：类似 strings.Builder 的字节构建工具

## 安装

```bash
go get github.com/snail007/gmc/util/bytes
```

## 快速开始

### 循环缓冲区

循环缓冲区是一个固定大小的环形缓冲区，当缓冲区满时会自动覆盖旧数据：

```go
package main

import (
    "fmt"
    "io"
    "github.com/snail007/gmc/util/bytes"
)

func main() {
    // 创建一个 1KB 的循环缓冲区
    buffer := gbytes.NewCircularBuffer(1024)
    defer buffer.Close()
    
    // 写入数据
    buffer.Write([]byte("Hello, World!"))
    
    // 创建读取器（从当前位置开始读）
    reader := buffer.NewReader()
    defer reader.Close()
    
    // 读取数据
    data := make([]byte, 13)
    n, err := reader.Read(data)
    if err != nil && err != io.EOF {
        panic(err)
    }
    fmt.Printf("Read %d bytes: %s\n", n, string(data[:n]))
}
```

### 字节池

字节池用于重用字节切片，减少 GC 压力：

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/bytes"
)

func main() {
    // 获取或创建一个 1024 字节的池
    pool := gbytes.GetPool(1024)
    
    // 从池中获取字节切片
    buf := pool.Get().([]byte)
    
    // 使用字节切片
    copy(buf, []byte("Hello"))
    fmt.Println(string(buf[:5]))
    
    // 使用完后放回池中
    pool.Put(buf)
}
```

### 字节大小解析和格式化

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/bytes"
)

func main() {
    // 解析人类可读的大小字符串
    size, err := gbytes.ParseSize("1.5GB")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Bytes: %d\n", size) // 1610612736
    
    // 格式化字节大小为人类可读格式
    sizeStr, err := gbytes.SizeStr(1610612736)
    if err != nil {
        panic(err)
    }
    fmt.Println(sizeStr) // "1.5 GB"
    
    // 使用 ByteSize 类型
    var bs gbytes.ByteSize
    bs.Parse("2TB")
    fmt.Printf("Bytes: %d\n", bs.Bytes())
    fmt.Printf("GB: %.2f\n", bs.GBytes())
    fmt.Printf("Human: %s\n", bs.HumanReadable())
}
```

### IO 写入器

```go
package main

import (
    "bytes"
    "github.com/snail007/gmc/util/bytes"
)

func main() {
    buf := &bytes.Buffer{}
    writer := gbytes.NewIOWriter(buf)
    
    // 写入字符串
    writer.WriteStr("Hello, ")
    writer.WriteStr("World!")
    
    // 写入带换行的字符串
    writer.WriteStrLn("New line")
    
    // 格式化写入
    writer.WriteStr("Number: %d", 42)
    
    println(buf.String())
}
```

### 字节构建器

```go
package main

import (
    "github.com/snail007/gmc/util/bytes"
)

func main() {
    builder := gbytes.NewBytesBuilder()
    
    builder.WriteStr("Hello")
    builder.WriteStr(", ")
    builder.WriteStrLn("World!")
    builder.WriteStr("Count: %d", 100)
    
    // 获取结果
    result := builder.String()
    println(result)
    
    // 或获取字节切片
    data := builder.Bytes()
    println(string(data))
}
```

## API 参考

### CircularBuffer - 循环缓冲区

#### 创建和基本操作

```go
func NewCircularBuffer(size int) *CircularBuffer
```

创建指定大小的循环缓冲区。

**方法：**

- `Write(p []byte) (n int, err error)`：写入数据
- `NewReader() io.ReadCloser`：创建从当前位置开始的读取器
- `NewHistoryReader() io.ReadCloser`：创建从缓冲区开始位置的读取器
- `Bytes() []byte`：获取缓冲区内容的拷贝
- `Close() error`：关闭缓冲区
- `Reset()`：重置缓冲区
- `Writer() *IOWriter`：获取 IOWriter 包装器

**示例：**

```go
buffer := gbytes.NewCircularBuffer(4096)

// 写入日志
buffer.Write([]byte("log entry 1\n"))
buffer.Write([]byte("log entry 2\n"))

// 多个读取器可以同时读取
reader1 := buffer.NewReader()
reader2 := buffer.NewHistoryReader() // 从头开始读

// 继续写入会覆盖旧数据
for i := 0; i < 1000; i++ {
    buffer.Write([]byte(fmt.Sprintf("entry %d\n", i)))
}
```

#### 读取器设置

```go
func (b *CircularBuffer) SetReaderDeadline(r io.ReadCloser, deadline time.Time)
func (b *CircularBuffer) ResetReader(r io.ReadCloser)
```

设置读取器的截止时间和重置读取器位置。

### Pool - 字节池

#### 全局池函数

```go
func GetPool(bufSize int) *Pool
func GetPoolCap(bufSize, capSize int) *Pool
```

获取或创建指定大小的字节池。

**示例：**

```go
// 创建 1KB 的池
pool := gbytes.GetPool(1024)

// 创建带容量的池
pool := gbytes.GetPoolCap(512, 1024) // len=512, cap=1024
```

#### 创建自定义池

```go
func NewPool(bufSize int) *Pool
func NewPoolCap(bufSize, capSize int) *Pool
```

**方法：**

- `Get() interface{}`：从池中获取字节切片
- `Put(x interface{})`：将字节切片放回池中

### ByteSize - 字节大小

#### 常量

```go
const (
    B  ByteSize = 1
    KB          = 1024
    MB          = 1024 * KB
    GB          = 1024 * MB
    TB          = 1024 * GB
    PB          = 1024 * TB
    EB          = 1024 * PB
)
```

#### 函数

```go
func ParseSize(s string) (bytes uint64, err error)
func SizeStr(bytes uint64) (s string, err error)
```

解析和格式化字节大小字符串。

**支持的单位：**
- Bytes: `B`, `byte`, `bytes`
- Kilobytes: `K`, `KB`, `kilo`, `kilobyte`, `kilobytes`
- Megabytes: `M`, `MB`, `mega`, `megabyte`, `megabytes`
- Gigabytes: `G`, `GB`, `giga`, `gigabyte`, `gigabytes`
- Terabytes: `T`, `TB`, `tera`, `terabyte`, `terabytes`
- Petabytes: `P`, `PB`, `peta`, `petabyte`, `petabytes`
- Exabytes: `E`, `EB`

#### 方法

```go
func (b *ByteSize) Parse(str string) error
func (b *ByteSize) MustParse(str string) *ByteSize
func (b ByteSize) Bytes() uint64
func (b ByteSize) KBytes() float64
func (b ByteSize) MBytes() float64
func (b ByteSize) GBytes() float64
func (b ByteSize) TBytes() float64
func (b ByteSize) PBytes() float64
func (b ByteSize) EBytes() float64
func (b ByteSize) HumanReadable() string
func (b ByteSize) String() string
```

**示例：**

```go
var size gbytes.ByteSize
size.Parse("1.5GB")

fmt.Println(size.Bytes())         // 1610612736
fmt.Println(size.GBytes())        // 1.5
fmt.Println(size.HumanReadable()) // "1.5 GB"
fmt.Println(size.String())        // "1.5GB"
```

### IOWriter - IO 写入器

```go
func NewIOWriter(w io.Writer) *IOWriter
```

创建一个便捷的写入器包装。

**方法：**

- `Write(data []byte) error`：写入字节
- `WriteLn(data []byte) error`：写入字节并添加换行
- `WriteStr(format string, values ...interface{}) error`：格式化写入字符串
- `WriteStrLn(format string, values ...interface{}) error`：格式化写入并添加换行

### BytesBuilder - 字节构建器

```go
func NewBytesBuilder() *BytesBuilder
```

创建一个字节构建器。

**方法：**

- `Write(data []byte) error`
- `WriteLn(data []byte) error`
- `WriteStr(format string, values ...interface{}) error`
- `WriteStrLn(format string, values ...interface{}) error`
- `String() string`：获取构建的字符串
- `Bytes() []byte`：获取字节切片

## 使用场景

### 循环缓冲区

1. **日志缓存**：保留最新的 N 字节日志
2. **实时数据流**：多个消费者读取实时数据
3. **命令输出捕获**：捕获命令的最新输出
4. **网络数据缓存**：缓存最新的网络数据

### 字节池

1. **高频内存分配**：减少频繁的字节切片分配
2. **网络 IO**：复用读写缓冲区
3. **编解码**：复用编解码缓冲区

### 字节大小处理

1. **配置文件**：解析配置中的大小设置
2. **存储管理**：显示和处理存储容量
3. **带宽计算**：处理网络带宽值

## 注意事项

### 循环缓冲区

1. **数据覆盖**：缓冲区满时会覆盖最旧的数据
2. **读取器管理**：及时关闭不再使用的读取器
3. **线程安全**：所有操作都是线程安全的
4. **截止时间**：读取阻塞时可设置截止时间避免永久等待

### 字节池

1. **大小固定**：每个池的字节切片大小固定
2. **不要修改容量**：放回池中的切片不应被修改容量
3. **及时归还**：使用完及时放回池中

### 字节大小

1. **单位大小写**：注意单位的大小写（"Kb" 表示 bits 会报错）
2. **精度损失**：浮点数转换可能有精度损失
3. **支持小数**：支持 "1.5GB" 这样的小数表示

## 性能考虑

- 循环缓冲区使用固定内存，不会无限增长
- 字节池基于 sync.Pool，自动管理生命周期
- ByteSize 计算不涉及复杂运算，性能开销小
- IOWriter 和 BytesBuilder 适合频繁写入场景

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
