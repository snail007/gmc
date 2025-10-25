# genv 包

## 简介

genv 包提供了一个便捷的环境变量访问器（Accessor），支持自动前缀功能和多种数据类型转换。

## 功能特性

- **前缀支持**：自动为环境变量添加前缀
- **类型转换**：返回 `gvalue.AnyValue` 类型，支持多种类型转换
- **链式调用**：支持链式设置环境变量
- **查找功能**：支持检查环境变量是否存在

## 安装

```bash
go get github.com/snail007/gmc/util/env
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/env"
)

func main() {
    // 创建不带前缀的访问器
    env := genv.NewAccessor("")
    
    // 获取环境变量
    home := env.Get("HOME")
    fmt.Println("Home:", home.String())
    
    // 设置环境变量
    env.Set("MY_VAR", "hello")
    
    // 获取并转换类型
    value := env.Get("MY_VAR")
    fmt.Println(value.String())
}
```

### 使用前缀

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/env"
)

func main() {
    // 创建带前缀的访问器
    env := genv.NewAccessor("MYAPP_")
    
    // 设置环境变量 MYAPP_PORT
    env.Set("PORT", "8080")
    
    // 获取环境变量 MYAPP_PORT
    port := env.Get("PORT")
    fmt.Println("Port:", port.Int())
    
    // 设置环境变量 MYAPP_DEBUG
    env.Set("DEBUG", "true")
    debug := env.Get("DEBUG")
    fmt.Println("Debug:", debug.Bool())
}
```

### 检查环境变量是否存在

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/env"
)

func main() {
    env := genv.NewAccessor("APP_")
    
    // 检查环境变量是否存在
    if value, found := env.Lookup("DATABASE_URL"); found {
        fmt.Println("Database URL:", value.String())
    } else {
        fmt.Println("DATABASE_URL not set")
    }
}
```

### 链式设置

```go
package main

import (
    "github.com/snail007/gmc/util/env"
)

func main() {
    env := genv.NewAccessor("CONFIG_")
    
    // 链式设置多个环境变量
    env.Set("HOST", "localhost").
        Set("PORT", "3000").
        Set("TIMEOUT", "30")
}
```

## API 参考

### NewAccessor

```go
func NewAccessor(prefix string) *Accessor
```

创建一个新的环境变量访问器。

**参数：**
- `prefix`：环境变量前缀，会自动添加到所有操作的变量名前

**返回值：**
- `*Accessor`：访问器实例

**示例：**

```go
// 不使用前缀
env := genv.NewAccessor("")

// 使用前缀 "MYAPP_"
env := genv.NewAccessor("MYAPP_")
```

### Accessor 方法

#### Get

```go
func (s *Accessor) Get(key string) *gvalue.AnyValue
```

获取环境变量的值，会自动添加前缀。

**参数：**
- `key`：环境变量名（不包含前缀）

**返回值：**
- `*gvalue.AnyValue`：环境变量值，支持多种类型转换

**示例：**

```go
env := genv.NewAccessor("APP_")

// 获取 APP_PORT
port := env.Get("PORT")
fmt.Println(port.Int())

// 获取 APP_DEBUG
debug := env.Get("DEBUG")
fmt.Println(debug.Bool())
```

#### Lookup

```go
func (s *Accessor) Lookup(key string) (*gvalue.AnyValue, bool)
```

查找环境变量，返回值和是否存在的标志。

**参数：**
- `key`：环境变量名（不包含前缀）

**返回值：**
- `*gvalue.AnyValue`：环境变量值
- `bool`：是否存在

**示例：**

```go
env := genv.NewAccessor("APP_")

if value, found := env.Lookup("PORT"); found {
    fmt.Println("Port found:", value.Int())
} else {
    fmt.Println("Port not set, using default")
    port := 8080
}
```

#### Set

```go
func (s *Accessor) Set(key, value string) *Accessor
```

设置环境变量，会自动添加前缀。支持链式调用。

**参数：**
- `key`：环境变量名（不包含前缀）
- `value`：环境变量值

**返回值：**
- `*Accessor`：返回自身，支持链式调用

**示例：**

```go
env := genv.NewAccessor("APP_")

// 单个设置
env.Set("PORT", "8080")

// 链式设置
env.Set("HOST", "localhost").
    Set("PORT", "8080").
    Set("DEBUG", "true")
```

#### Unset

```go
func (s *Accessor) Unset(key string) *Accessor
```

删除环境变量，会自动添加前缀。支持链式调用。

**参数：**
- `key`：环境变量名（不包含前缀）

**返回值：**
- `*Accessor`：返回自身，支持链式调用

**示例：**

```go
env := genv.NewAccessor("APP_")

// 删除环境变量
env.Unset("PORT")

// 链式删除
env.Unset("PORT").Unset("DEBUG")
```

## 使用场景

### 1. 应用配置管理

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/env"
)

type Config struct {
    Host    string
    Port    int
    Debug   bool
    Timeout int
}

func LoadConfig() *Config {
    env := genv.NewAccessor("MYAPP_")
    
    return &Config{
        Host:    env.Get("HOST").String(),
        Port:    env.Get("PORT").Int(),
        Debug:   env.Get("DEBUG").Bool(),
        Timeout: env.Get("TIMEOUT").Int(),
    }
}

func main() {
    env := genv.NewAccessor("MYAPP_")
    env.Set("HOST", "localhost").
        Set("PORT", "8080").
        Set("DEBUG", "true").
        Set("TIMEOUT", "30")
    
    config := LoadConfig()
    fmt.Printf("%+v\n", config)
}
```

### 2. 多环境配置

```go
package main

import (
    "github.com/snail007/gmc/util/env"
)

func main() {
    // 开发环境配置
    devEnv := genv.NewAccessor("DEV_")
    devEnv.Set("DB_HOST", "localhost").
           Set("DB_PORT", "5432")
    
    // 生产环境配置
    prodEnv := genv.NewAccessor("PROD_")
    prodEnv.Set("DB_HOST", "prod-db.example.com").
            Set("DB_PORT", "5432")
    
    // 根据环境选择配置
    currentEnv := prodEnv
    dbHost := currentEnv.Get("DB_HOST").String()
    
    println("Database Host:", dbHost)
}
```

### 3. 带默认值的配置

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/env"
)

func getConfigWithDefault(env *genv.Accessor, key, defaultValue string) string {
    if value, found := env.Lookup(key); found && value.String() != "" {
        return value.String()
    }
    return defaultValue
}

func main() {
    env := genv.NewAccessor("APP_")
    
    host := getConfigWithDefault(env, "HOST", "localhost")
    port := getConfigWithDefault(env, "PORT", "8080")
    
    fmt.Printf("Server: %s:%s\n", host, port)
}
```

### 4. 配置验证

```go
package main

import (
    "fmt"
    "os"
    "github.com/snail007/gmc/util/env"
)

func validateConfig() error {
    env := genv.NewAccessor("APP_")
    
    required := []string{"DATABASE_URL", "API_KEY", "SECRET_KEY"}
    
    for _, key := range required {
        if _, found := env.Lookup(key); !found {
            return fmt.Errorf("required environment variable not set: APP_%s", key)
        }
    }
    
    return nil
}

func main() {
    if err := validateConfig(); err != nil {
        fmt.Println("Config validation failed:", err)
        os.Exit(1)
    }
    
    fmt.Println("Config validated successfully")
}
```

### 5. 测试环境隔离

```go
package main

import (
    "testing"
    "github.com/snail007/gmc/util/env"
)

func TestWithEnv(t *testing.T) {
    env := genv.NewAccessor("TEST_")
    
    // 设置测试环境变量
    env.Set("MODE", "test").
        Set("DB_HOST", "localhost")
    
    // 运行测试
    mode := env.Get("MODE").String()
    if mode != "test" {
        t.Errorf("Expected test mode, got %s", mode)
    }
    
    // 清理
    env.Unset("MODE").Unset("DB_HOST")
}
```

## gvalue.AnyValue 类型

返回值是 `gvalue.AnyValue` 类型，支持以下转换方法（部分列举）：

- `String() string`：转换为字符串
- `Int() int`：转换为整数
- `Int64() int64`：转换为 int64
- `Float64() float64`：转换为 float64
- `Bool() bool`：转换为布尔值
- `IsEmpty() bool`：检查是否为空

详细 API 请参考 `gvalue` 包文档。

## 注意事项

1. **前缀管理**：前缀会自动添加到所有操作中，无需手动添加
2. **环境变量作用域**：设置的环境变量对当前进程及其子进程有效
3. **类型转换**：环境变量本质是字符串，类型转换可能失败，注意检查
4. **线程安全**：底层使用 `os.Getenv` 等函数，是线程安全的
5. **大小写敏感**：环境变量名在不同操作系统可能有不同的大小写敏感性

## 最佳实践

### 1. 统一前缀

为应用的所有环境变量使用统一前缀，避免冲突：

```go
env := genv.NewAccessor("MYAPP_")
```

### 2. 配置结构化

使用结构体管理配置：

```go
type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
}

func LoadDBConfig(env *genv.Accessor) *DatabaseConfig {
    return &DatabaseConfig{
        Host:     env.Get("DB_HOST").String(),
        Port:     env.Get("DB_PORT").Int(),
        User:     env.Get("DB_USER").String(),
        Password: env.Get("DB_PASSWORD").String(),
    }
}
```

### 3. 提供默认值

总是为配置提供合理的默认值：

```go
func getPort(env *genv.Accessor) int {
    if port, found := env.Lookup("PORT"); found {
        if p := port.Int(); p > 0 {
            return p
        }
    }
    return 8080 // 默认端口
}
```

### 4. 配置文档化

为每个环境变量编写文档说明：

```go
// Environment Variables:
// APP_HOST - Server host (default: localhost)
// APP_PORT - Server port (default: 8080)
// APP_DEBUG - Enable debug mode (default: false)
// APP_LOG_LEVEL - Log level: debug, info, warn, error (default: info)
```

## 依赖

- `github.com/snail007/gmc/util/value`：值包装和类型转换
- `os`：操作系统环境变量操作

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [gvalue 包文档](../value/)
- [12-Factor App - Config](https://12factor.net/config)
