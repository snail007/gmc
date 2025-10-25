# GMC Config 模块

## 简介

GMC Config 模块基于 Viper 封装，提供强大的配置管理功能，支持多种配置文件格式、环境变量、配置搜索等特性。

## 功能特性

- **多种格式支持**：支持 TOML、YAML、JSON、HCL、INI 等格式
- **环境变量绑定**：自动绑定环境变量
- **配置搜索**：在多个路径中搜索配置文件
- **配置热加载**：支持监听配置文件变化
- **默认值**：支持设置默认配置值
- **类型转换**：自动进行类型转换
- **子配置**：支持获取配置的子节点

## 安装

```bash
go get github.com/snail007/gmc/module/config
```

## 快速开始

### 从文件加载配置

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/module/config"
)

func main() {
    // 从文件加载配置
    cfg, err := gconfig.NewFromFile("app.toml")
    if err != nil {
        panic(err)
    }
    
    // 读取配置项
    appName := cfg.GetString("app.name")
    port := cfg.GetInt("server.port")
    debug := cfg.GetBool("app.debug")
    
    fmt.Printf("App: %s, Port: %d, Debug: %v\n", appName, port, debug)
}
```

### 从字节数组加载配置

```go
package main

import (
    "github.com/snail007/gmc/module/config"
)

func main() {
    configData := []byte(`
[app]
name = "myapp"
version = "1.0.0"

[server]
port = 8080
`)
    
    // 从字节数组加载配置（TOML 格式）
    cfg, err := gconfig.NewConfigBytes(configData, "toml")
    if err != nil {
        panic(err)
    }
    
    name := cfg.GetString("app.name")
    port := cfg.GetInt("server.port")
}
```

### 搜索配置文件

```go
package main

import (
    "github.com/snail007/gmc/module/config"
)

func main() {
    // 在多个路径中搜索配置文件
    paths := []string{".", "./conf", "./config", "/etc/myapp"}
    cfg, err := gconfig.NewFromSearch(paths, "app.toml")
    if err != nil {
        panic(err)
    }
    
    // 使用配置
    _ = cfg
}
```

### 使用默认配置

```go
package main

import (
    "github.com/snail007/gmc/module/config"
)

func main() {
    cfg := gconfig.New()
    
    // 设置配置搜索路径
    cfg.AddConfigPath(".")
    cfg.AddConfigPath("./conf")
    cfg.AddConfigPath("./config")
    
    // 设置配置文件名（不含扩展名）
    cfg.SetConfigName("app")
    
    // 设置配置文件类型
    cfg.SetConfigType("toml")
    
    // 读取配置
    err := cfg.ReadInConfig()
    if err != nil {
        panic(err)
    }
    
    // 使用配置
    port := cfg.GetInt("server.port")
}
```

### 环境变量

```go
package main

import (
    "os"
    "github.com/snail007/gmc/module/config"
)

func main() {
    // 设置环境变量前缀
    os.Setenv("ENV_PREFIX", "MYAPP")
    
    // 设置环境变量
    os.Setenv("MYAPP_SERVER_PORT", "9000")
    os.Setenv("MYAPP_APP_DEBUG", "true")
    
    cfg := gconfig.New()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    
    // 环境变量会覆盖配置文件中的值
    // MYAPP_SERVER_PORT 对应 server.port
    port := cfg.GetInt("server.port") // 返回 9000
}
```

## API 参考

### 创建配置

```go
// 创建新配置对象
func New() *Config

// 从文件加载配置
func NewFromFile(file string, typ ...string) (*Config, error)

// 从字节数组加载配置
func NewConfigBytes(b []byte, typ ...string) (*Config, error)

// 搜索并加载配置文件
func NewFromSearch(paths []string, filename string, typ ...string) (*Config, error)
```

### 读取配置

```go
// 读取配置文件
ReadInConfig() error

// 读取配置（从 io.Reader）
ReadConfig(in io.Reader) error

// 设置配置文件路径
SetConfigFile(file string)

// 设置配置文件名（不含扩展名）
SetConfigName(name string)

// 设置配置文件类型
SetConfigType(typ string)

// 添加配置搜索路径
AddConfigPath(path string)
```

### 获取配置值

```go
// 获取字符串
GetString(key string) string

// 获取整数
GetInt(key string) int
GetInt32(key string) int32
GetInt64(key string) int64

// 获取无符号整数
GetUint(key string) uint
GetUint32(key string) uint32
GetUint64(key string) uint64

// 获取浮点数
GetFloat64(key string) float64

// 获取布尔值
GetBool(key string) bool

// 获取时间
GetTime(key string) time.Time
GetDuration(key string) time.Duration

// 获取数组
GetStringSlice(key string) []string
GetIntSlice(key string) []int

// 获取 Map
GetStringMap(key string) map[string]interface{}
GetStringMapString(key string) map[string]string

// 获取子配置
Sub(key string) gcore.SubConfig

// 获取所有配置
AllSettings() map[string]interface{}

// 检查键是否存在
IsSet(key string) bool
```

### 设置配置值

```go
// 设置配置值
Set(key string, value interface{})

// 设置默认值
SetDefault(key string, value interface{})
```

## 配置文件示例

### TOML 格式

```toml
[app]
name = "myapp"
version = "1.0.0"
debug = false

[server]
host = "0.0.0.0"
port = 8080
timeout = 30

[database]
driver = "mysql"
host = "localhost"
port = 3306
username = "root"
password = "password"
database = "mydb"

[[database.read_replicas]]
host = "replica1.example.com"
port = 3306

[[database.read_replicas]]
host = "replica2.example.com"
port = 3306
```

### YAML 格式

```yaml
app:
  name: myapp
  version: 1.0.0
  debug: false

server:
  host: 0.0.0.0
  port: 8080
  timeout: 30

database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  database: mydb
  read_replicas:
    - host: replica1.example.com
      port: 3306
    - host: replica2.example.com
      port: 3306
```

## 使用场景

1. **应用配置**：管理应用的各种配置项
2. **多环境配置**：通过环境变量区分开发/测试/生产环境
3. **功能开关**：通过配置控制功能的开启和关闭
4. **服务配置**：配置数据库、缓存、消息队列等服务连接
5. **动态配置**：监听配置文件变化，动态更新配置

## 最佳实践

### 1. 使用结构体绑定配置

```go
type AppConfig struct {
    Name    string
    Version string
    Debug   bool
}

type ServerConfig struct {
    Host    string
    Port    int
    Timeout int
}

type Config struct {
    App    AppConfig
    Server ServerConfig
}

func LoadConfig() (*Config, error) {
    cfg := gconfig.New()
    cfg.SetConfigFile("app.toml")
    err := cfg.ReadInConfig()
    if err != nil {
        return nil, err
    }
    
    var config Config
    err = cfg.Unmarshal(&config)
    if err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### 2. 使用子配置

```go
cfg := gconfig.New()
cfg.SetConfigFile("app.toml")
cfg.ReadInConfig()

// 获取数据库子配置
dbConfig := cfg.Sub("database")
if dbConfig != nil {
    host := dbConfig.GetString("host")
    port := dbConfig.GetInt("port")
}
```

### 3. 设置默认值

```go
cfg := gconfig.New()
cfg.SetDefault("server.port", 8080)
cfg.SetDefault("server.host", "0.0.0.0")
cfg.SetDefault("app.debug", false)

cfg.SetConfigFile("app.toml")
cfg.ReadInConfig()

// 如果配置文件中没有这些值，使用默认值
port := cfg.GetInt("server.port")
```

### 4. 环境变量覆盖

```go
// 设置环境变量前缀
os.Setenv("ENV_PREFIX", "MYAPP")

// 环境变量命名规则：
// MYAPP_SERVER_PORT 对应 server.port
// MYAPP_DATABASE_HOST 对应 database.host

cfg := gconfig.New()
cfg.SetConfigFile("app.toml")
cfg.ReadInConfig()

// 环境变量会自动覆盖配置文件中的值
```

## 注意事项

1. **配置文件格式**：需要指定正确的配置文件类型
2. **环境变量命名**：使用下划线分隔，自动转换为配置键的点号
3. **配置搜索顺序**：按添加顺序搜索配置路径
4. **默认前缀**：默认环境变量前缀为 "GMC"
5. **配置优先级**：环境变量 > 配置文件 > 默认值

## 依赖

- `github.com/spf13/viper`：底层配置管理库

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [Viper 文档](https://github.com/spf13/viper)
