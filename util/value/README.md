# gvalue 包

## 简介

gvalue 包提供了通用值类型包装器，支持将任意类型的值转换为各种 Go 基本类型。

## 功能特性

- **类型转换**：支持转换为 int、string、bool、float 等类型
- **安全转换**：转换失败返回零值而不是 panic
- **链式调用**：支持链式操作
- **空值处理**：智能处理 nil 和空值

## 安装

```bash
go get github.com/snail007/gmc/util/value
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/value"
)

func main() {
    // 创建 Value
    v := gvalue.New(42)
    
    // 转换为不同类型
    fmt.Println(v.Int())     // 42
    fmt.Println(v.String())  // "42"
    fmt.Println(v.Float64()) // 42.0
    fmt.Println(v.Bool())    // true（非零值为 true）
    
    // 字符串值
    s := gvalue.New("123")
    fmt.Println(s.Int())     // 123
    fmt.Println(s.String())  // "123"
    
    // AnyValue 类型
    any := gvalue.NewAny("hello")
    fmt.Println(any.String())   // "hello"
    fmt.Println(any.IsEmpty())  // false
    
    // 处理 nil
    nilVal := gvalue.New(nil)
    fmt.Println(nilVal.String()) // ""
    fmt.Println(nilVal.Int())    // 0
}
```

## API 参考

### Value 类型

```go
func New(v interface{}) *Value
```

**方法：**
- `String() string`：转换为字符串
- `Int() int`：转换为 int
- `Int64() int64`：转换为 int64
- `Uint() uint`：转换为 uint
- `Uint64() uint64`：转换为 uint64
- `Float32() float32`：转换为 float32
- `Float64() float64`：转换为 float64
- `Bool() bool`：转换为 bool
- `Interface() interface{}`：获取原始值
- `Bytes() []byte`：转换为字节数组

### AnyValue 类型

```go
func NewAny(v interface{}) *AnyValue
```

**方法：**
- 包含 Value 的所有方法
- `IsEmpty() bool`：是否为空
- `IsNil() bool`：是否为 nil

## 类型转换规则

### 转 String
- 数字：转换为字符串表示
- 布尔：转换为 "true" 或 "false"
- nil：返回空字符串

### 转 Int/Int64
- 字符串：解析数字字符串
- 布尔：true=1, false=0
- 浮点：向下取整
- nil：返回 0

### 转 Bool
- 数字：非零为 true
- 字符串："true", "1", "yes", "on" 为 true
- nil：返回 false

### 转 Float
- 字符串：解析浮点数字符串
- 布尔：true=1.0, false=0.0
- 整数：转换为浮点数
- nil：返回 0.0

## 使用场景

1. **配置解析**：处理配置文件中的值
2. **HTTP 参数**：转换 URL 参数和表单数据
3. **环境变量**：处理环境变量值
4. **数据库结果**：转换数据库查询结果
5. **JSON 数据**：处理动态 JSON 数据

## 示例

### 处理配置值

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/value"
)

func main() {
    config := map[string]interface{}{
        "port":    "8080",
        "debug":   "true",
        "timeout": 30,
    }
    
    port := gvalue.New(config["port"]).Int()
    debug := gvalue.New(config["debug"]).Bool()
    timeout := gvalue.New(config["timeout"]).Int()
    
    fmt.Printf("Port: %d, Debug: %v, Timeout: %d\n", 
        port, debug, timeout)
}
```

### 安全的类型转换

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/value"
)

func main() {
    // 即使值无法转换，也不会 panic
    v1 := gvalue.New("not a number")
    fmt.Println(v1.Int()) // 0
    
    v2 := gvalue.New(nil)
    fmt.Println(v2.String()) // ""
    
    v3 := gvalue.New("123abc")
    fmt.Println(v3.Int()) // 123（尽可能解析）
}
```

## 注意事项

1. **零值返回**：转换失败返回零值而不是错误
2. **精度损失**：浮点数转整数会丢失小数部分
3. **字符串解析**：字符串转数字会尽可能解析前缀数字
4. **布尔转换**：字符串转布尔不区分大小写

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
