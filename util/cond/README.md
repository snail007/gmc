# gcond 包

## 简介

gcond 包提供了类似三元运算符的条件表达式函数，使代码更加简洁。返回值包装在 `gvalue.Value` 类型中，支持多种类型转换。

## 功能特性

- **条件选择**：根据布尔值选择返回不同的值
- **惰性求值**：支持函数形式的惰性求值，避免不必要的计算
- **类型灵活**：返回值支持多种类型转换

## 安装

```bash
go get github.com/snail007/gmc/util/cond
```

## 快速开始

### 基本条件选择

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/cond"
)

func main() {
    age := 20
    
    // 使用 Cond 进行条件选择
    status := gcond.Cond(age >= 18, "成年人", "未成年人")
    
    fmt.Println(status.String()) // 输出: 成年人
}
```

### 惰性求值

使用 `CondFn` 可以延迟计算，只有在需要时才执行函数：

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/cond"
)

func expensiveTrue() interface{} {
    fmt.Println("Computing true value...")
    return "expensive true"
}

func expensiveFalse() interface{} {
    fmt.Println("Computing false value...")
    return "expensive false"
}

func main() {
    // 只会执行 expensiveTrue，不会执行 expensiveFalse
    result := gcond.CondFn(true, expensiveTrue, expensiveFalse)
    fmt.Println(result.String())
    
    // 输出:
    // Computing true value...
    // expensive true
}
```

### 类型转换

返回值是 `gvalue.Value` 类型，支持多种类型转换：

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/cond"
)

func main() {
    // 数字类型
    num := gcond.Cond(true, 42, 0)
    fmt.Println(num.Int())     // 42
    fmt.Println(num.String())  // "42"
    
    // 字符串类型
    str := gcond.Cond(false, "yes", "no")
    fmt.Println(str.String())  // "no"
    
    // 布尔类型
    flag := gcond.Cond(1 > 0, 1, 0)
    fmt.Println(flag.Bool())   // true
}
```

## API 参考

### Cond

```go
func Cond(check bool, ok interface{}, fail interface{}) *gvalue.Value
```

根据条件选择返回值。

**参数：**
- `check`：条件表达式
- `ok`：条件为 true 时返回的值
- `fail`：条件为 false 时返回的值

**返回值：**
- `*gvalue.Value`：包装后的值，支持多种类型转换

**示例：**

```go
// 简单值选择
result := gcond.Cond(score >= 60, "及格", "不及格")

// 数字选择
max := gcond.Cond(a > b, a, b)

// 结构体选择
user := gcond.Cond(isAdmin, adminUser, normalUser)
```

### CondFn

```go
func CondFn(check bool, ok func() interface{}, fail func() interface{}) *gvalue.Value
```

根据条件惰性求值返回结果。只有被选中的函数会被执行。

**参数：**
- `check`：条件表达式
- `ok`：条件为 true 时执行的函数
- `fail`：条件为 false 时执行的函数

**返回值：**
- `*gvalue.Value`：包装后的值

**示例：**

```go
result := gcond.CondFn(
    user.IsVIP(),
    func() interface{} {
        return calculateVIPPrice(item)
    },
    func() interface{} {
        return item.RegularPrice
    },
)
```

## 使用场景

### 1. 替代 if-else 进行赋值

```go
// 传统方式
var message string
if isSuccess {
    message = "操作成功"
} else {
    message = "操作失败"
}

// 使用 Cond
message := gcond.Cond(isSuccess, "操作成功", "操作失败").String()
```

### 2. 内联条件判断

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/cond"
)

func main() {
    items := []string{"apple", "banana", "cherry"}
    
    for i, item := range items {
        prefix := gcond.Cond(i == 0, "首个", "其他").String()
        fmt.Printf("%s: %s\n", prefix, item)
    }
}
```

### 3. 配置选择

```go
package main

import (
    "os"
    "github.com/snail007/gmc/util/cond"
)

func main() {
    env := os.Getenv("ENV")
    
    // 根据环境选择配置
    dbHost := gcond.Cond(env == "production", 
        "prod-db.example.com", 
        "dev-db.example.com").String()
    
    maxConns := gcond.Cond(env == "production", 100, 10).Int()
    
    println("DB Host:", dbHost)
    println("Max Connections:", maxConns)
}
```

### 4. 避免重复计算

使用 `CondFn` 避免不必要的函数调用：

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/cond"
)

func queryDatabase() interface{} {
    fmt.Println("Querying database...")
    // 模拟耗时操作
    return "database result"
}

func getFromCache() interface{} {
    fmt.Println("Getting from cache...")
    return "cache result"
}

func main() {
    hasCache := true
    
    // 如果有缓存，不会执行 queryDatabase
    result := gcond.CondFn(hasCache, getFromCache, queryDatabase)
    fmt.Println(result.String())
    
    // 输出:
    // Getting from cache...
    // cache result
}
```

### 5. 权限检查

```go
package main

import (
    "github.com/snail007/gmc/util/cond"
)

type User struct {
    Role string
}

func (u *User) IsAdmin() bool {
    return u.Role == "admin"
}

func main() {
    user := &User{Role: "user"}
    
    // 根据权限返回不同的数据
    data := gcond.CondFn(
        user.IsAdmin(),
        func() interface{} {
            return fetchAllData()
        },
        func() interface{} {
            return fetchUserData(user)
        },
    )
    
    _ = data
}

func fetchAllData() interface{} {
    return "all data"
}

func fetchUserData(user *User) interface{} {
    return "user data"
}
```

### 6. 默认值处理

```go
package main

import (
    "github.com/snail007/gmc/util/cond"
)

func main() {
    var username string
    // username 为空时使用默认值
    displayName := gcond.Cond(username != "", username, "访客").String()
    
    println("欢迎,", displayName)
}
```

## gvalue.Value 类型

返回值是 `gvalue.Value` 类型，支持以下转换方法（部分列举）：

- `String() string`：转换为字符串
- `Int() int`：转换为整数
- `Int64() int64`：转换为 int64
- `Float64() float64`：转换为 float64
- `Bool() bool`：转换为布尔值
- `Interface() interface{}`：获取原始值

详细 API 请参考 `gvalue` 包文档。

## 性能考虑

### Cond vs CondFn

- **Cond**：两个分支的值都会被计算（即使不会被使用）
- **CondFn**：只有被选中的分支会被执行（惰性求值）

**示例：**

```go
// 不推荐：两个函数都会被执行
result := gcond.Cond(condition, expensiveFunc1(), expensiveFunc2())

// 推荐：只执行需要的函数
result := gcond.CondFn(condition, expensiveFunc1, expensiveFunc2)
```

### 性能建议

1. **简单值**：使用 `Cond`，性能开销小
2. **函数调用**：使用 `CondFn`，避免不必要的计算
3. **复杂逻辑**：如果条件判断很复杂，传统 if-else 可能更清晰

## 注意事项

1. **类型安全**：interface{} 类型需要在使用时进行适当的类型转换
2. **nil 值**：支持返回 nil 值
3. **惰性求值**：`CondFn` 的函数不会预先执行
4. **错误处理**：gvalue.Value 的转换方法可能返回零值，需要注意

## 与其他语言对比

### JavaScript/TypeScript

```javascript
// JavaScript
const result = condition ? "yes" : "no"

// Go with gcond
result := gcond.Cond(condition, "yes", "no").String()
```

### Python

```python
# Python
result = "yes" if condition else "no"

# Go with gcond
result := gcond.Cond(condition, "yes", "no").String()
```

### C/C++/Java

```c
// C/C++/Java
String result = condition ? "yes" : "no";

// Go with gcond
result := gcond.Cond(condition, "yes", "no").String()
```

## 依赖

- `github.com/snail007/gmc/util/value`：值包装和类型转换

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [gvalue 包文档](../value/)
