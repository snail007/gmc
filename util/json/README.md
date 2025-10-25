# gjson 包

## 简介

gjson 包提供了强大的 JSON 解析和操作功能，基于 tidwall/gjson 和 tidwall/sjson 封装，支持快速查询、修改和遍历 JSON 数据。

## 功能特性

- **快速查询**：使用路径语法快速查询 JSON 值
- **类型转换**：自动转换为各种 Go 类型
- **修改 JSON**：动态修改 JSON 数据
- **数组遍历**：便捷的数组和对象遍历
- **JSON Lines**：支持 JSON Lines 格式

## 安装

```bash
go get github.com/snail007/gmc/util/json
```

## 快速开始

### 基本查询

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/json"
)

func main() {
    json := `{
        "name": "John",
        "age": 30,
        "address": {
            "city": "New York",
            "zip": "10001"
        }
    }`
    
    // 获取值
    name := gjson.Get(json, "name")
    fmt.Println(name.String()) // John
    
    age := gjson.Get(json, "age")
    fmt.Println(age.Int()) // 30
    
    city := gjson.Get(json, "address.city")
    fmt.Println(city.String()) // New York
}
```

### 数组查询

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/json"
)

func main() {
    json := `{
        "users": [
            {"name": "Alice", "age": 25},
            {"name": "Bob", "age": 30},
            {"name": "Charlie", "age": 35}
        ]
    }`
    
    // 获取数组元素
    first := gjson.Get(json, "users.0.name")
    fmt.Println(first.String()) // Alice
    
    // 获取数组长度
    users := gjson.Get(json, "users")
    fmt.Println("Length:", users.Len()) // 3
    
    // 遍历数组
    users.ForEach(func(key, value gjson.Result) bool {
        fmt.Printf("%s: %d\n", value.Get("name").String(), value.Get("age").Int())
        return true // 继续遍历
    })
}
```

### 修改 JSON

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/json"
)

func main() {
    json := `{"name":"John","age":30}`
    
    // 设置值
    result, _ := gjson.Set(json, "age", 31)
    fmt.Println(result) // {"name":"John","age":31}
    
    // 添加新字段
    result, _ = gjson.Set(result, "city", "New York")
    fmt.Println(result) // {"name":"John","age":31,"city":"New York"}
    
    // 删除字段
    result, _ = gjson.Delete(result, "age")
    fmt.Println(result) // {"name":"John","city":"New York"}
}
```

### 类型转换

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/json"
)

func main() {
    json := `{
        "string": "hello",
        "int": 42,
        "float": 3.14,
        "bool": true,
        "array": [1, 2, 3],
        "object": {"key": "value"}
    }`
    
    fmt.Println(gjson.Get(json, "string").String())
    fmt.Println(gjson.Get(json, "int").Int())
    fmt.Println(gjson.Get(json, "float").Float())
    fmt.Println(gjson.Get(json, "bool").Bool())
    fmt.Println(gjson.Get(json, "array").Array())
    fmt.Println(gjson.Get(json, "object").Map())
}
```

## API 参考

### 查询函数

- `Get(json, path) Result`：获取 JSON 值
- `GetBytes(json, path) Result`：从字节数组获取
- `GetMany(json, ...paths) []Result`：获取多个值
- `Parse(json) Result`：解析整个 JSON

### 修改函数

- `Set(json, path, value) (string, error)`：设置值
- `SetBytes(json, path, value) ([]byte, error)`：设置值（字节）
- `Delete(json, path) (string, error)`：删除字段
- `DeleteBytes(json, path) ([]byte, error)`：删除字段（字节）

### 验证函数

- `Valid(json) bool`：验证 JSON 是否有效
- `ValidBytes(json) bool`：验证字节数组 JSON

### Result 方法

- `String() string`
- `Int() int`
- `Int64() int64`
- `Float() float64`
- `Bool() bool`
- `Array() []Result`
- `Map() map[string]Result`
- `Exists() bool`
- `Type() Type`
- `ForEach(func(key, value Result) bool)`

## 路径语法

```go
// 基本路径
"name"
"address.city"
"users.0.name"

// 数组索引
"items.0"      // 第一个元素
"items.-1"     // 最后一个元素

// 通配符
"users.*.name" // 所有用户的 name

// 查询
"users.#.name"              // 所有 name 数组
"users.#(age>25).name"      // age > 25 的 name
"users.#(name%\"*li*\")#"   // name 包含 "li" 的数量

// 修饰符
"@reverse"     // 反转数组
"@ugly"        // 压缩 JSON
"@pretty"      // 格式化 JSON
```

## 使用场景

1. **配置文件**：解析 JSON 配置
2. **API 响应**：处理 HTTP API 返回的 JSON
3. **数据提取**：从复杂 JSON 提取特定字段
4. **JSON 转换**：修改 JSON 结构
5. **日志解析**：解析 JSON 格式的日志

## 注意事项

1. 路径不存在时返回空 Result
2. 类型转换失败返回零值
3. Set/Delete 会创建新的 JSON 字符串
4. 性能优于标准库的 encoding/json

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [gjson 文档](https://github.com/tidwall/gjson)
- [sjson 文档](https://github.com/tidwall/sjson)
