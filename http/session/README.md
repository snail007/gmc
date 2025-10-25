# GMC Session 模块

## 简介

GMC Session 模块提供了完整的 HTTP 会话管理功能，支持多种存储后端（Memory、File、Redis），提供简单易用的 API 来存储和管理用户会话数据。

## 功能特性

- **多种存储后端**：支持 Memory、File、Redis 三种存储方式
- **自动 GC**：自动清理过期会话
- **线程安全**：内置并发安全机制
- **灵活配置**：支持自定义 TTL、存储路径等
- **序列化支持**：自动序列化/反序列化会话数据
- **高性能**：基于 gob 编码，性能优异

## 安装

```bash
go get github.com/snail007/gmc
```

## 快速开始

### 在控制器中使用

```go
package controller

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type UserController struct {
    gmc.Controller
}

// 登录
func (c *UserController) Login() {
    username := c.PostForm("username")
    password := c.PostForm("password")
    
    // 验证用户名密码...
    
    // 启动会话
    err := c.SessionStart()
    if err != nil {
        c.JSON(gcore.M{"error": err.Error()})
        c.Stop()
        return
    }
    
    // 设置会话数据
    c.Session.Set("user_id", 123)
    c.Session.Set("username", username)
    c.Session.Set("login_time", time.Now().Unix())
    
    c.JSON(gcore.M{"status": "logged in"})
}

// 获取当前用户信息
func (c *UserController) Profile() {
    err := c.SessionStart()
    if err != nil {
        c.JSON(gcore.M{"error": err.Error()})
        c.Stop()
        return
    }
    
    // 读取会话数据
    userID := c.Session.Get("user_id")
    username := c.Session.Get("username")
    
    if userID == nil {
        c.WriteHeader(401)
        c.JSON(gcore.M{"error": "Not logged in"})
        c.Stop()
        return
    }
    
    c.JSON(gcore.M{
        "user_id":  userID,
        "username": username,
    })
}

// 登出
func (c *UserController) Logout() {
    err := c.SessionStart()
    if err != nil {
        c.JSON(gcore.M{"error": err.Error()})
        c.Stop()
        return
    }
    
    // 销毁会话
    err = c.SessionDestroy()
    if err != nil {
        c.JSON(gcore.M{"error": err.Error()})
        c.Stop()
        return
    }
    
    c.JSON(gcore.M{"status": "logged out"})
}
```

## 配置

### app.toml 配置示例

#### Memory 存储（默认）

```toml
[session]
# 启用会话
enable = true

# 存储类型: memory, file, redis
store = "memory"

# 会话过期时间（秒）
ttl = 3600

[session.memory]
# GC 执行间隔（秒）
gctime = 60
```

#### File 存储

```toml
[session]
enable = true
store = "file"
ttl = 3600

[session.file]
# 会话文件存储目录
dir = "./sessions"

# 文件名前缀
prefix = "gmcsess_"

# GC 执行间隔（秒）
gctime = 60
```

#### Redis 存储

```toml
[session]
enable = true
store = "redis"
ttl = 3600

[session.redis]
# Redis 服务器地址
address = "127.0.0.1:6379"

# Redis 密码（可选）
password = ""

# 数据库编号
dbnum = 0

# 键前缀
prefix = "gmcsess:"

# 连接超时（秒）
timeout = 10

# 最大空闲连接数
maxidle = 10

# 最大活动连接数
maxactive = 100

# 连接最大生命周期（秒）
maxconnlifetime = 300

# 是否等待可用连接
wait = true

# 调试模式
debug = false
```

## API 参考

### Session 对象

#### Set() - 设置会话值

```go
func (s *Session) Set(k interface{}, v interface{})
```

**示例：**
```go
c.Session.Set("user_id", 123)
c.Session.Set("username", "alice")
c.Session.Set("roles", []string{"admin", "user"})
c.Session.Set("preferences", map[string]string{
    "theme": "dark",
    "lang":  "zh-CN",
})
```

#### Get() - 获取会话值

```go
func (s *Session) Get(k interface{}) (value interface{})
```

**示例：**
```go
userID := c.Session.Get("user_id")
if userID != nil {
    // 类型断言
    id, ok := userID.(int)
    if ok {
        fmt.Println("User ID:", id)
    }
}
```

#### Delete() - 删除会话值

```go
func (s *Session) Delete(k interface{}) error
```

**示例：**
```go
err := c.Session.Delete("temp_data")
if err != nil {
    log.Println("Failed to delete:", err)
}
```

#### Destroy() - 销毁会话

```go
func (s *Session) Destroy() error
```

清空所有会话数据并标记为已销毁。

**示例：**
```go
err := c.Session.Destroy()
```

#### Values() - 获取所有会话数据

```go
func (s *Session) Values() map[interface{}]interface{}
```

**示例：**
```go
allData := c.Session.Values()
for key, value := range allData {
    fmt.Printf("%v: %v\n", key, value)
}
```

#### SessionID() - 获取会话 ID

```go
func (s *Session) SessionID() string
```

**示例：**
```go
sessionID := c.Session.SessionID()
fmt.Println("Session ID:", sessionID)
```

#### TouchTime() - 获取最后访问时间

```go
func (s *Session) TouchTime() int64
```

返回最后一次访问的 Unix 时间戳。

#### IsDestroy() - 检查是否已销毁

```go
func (s *Session) IsDestroy() bool
```

### 存储后端初始化

#### 手动初始化（不使用配置文件）

##### Memory Store

```go
package main

import (
    gsession "github.com/snail007/gmc/http/session"
)

func main() {
    cfg := gsession.NewMemoryStoreConfig()
    cfg.TTL = 3600    // 1 小时
    cfg.GCtime = 60   // 每 60 秒 GC 一次
    
    store, err := gsession.NewMemoryStore(cfg)
    if err != nil {
        panic(err)
    }
    
    // 使用 store...
}
```

##### File Store

```go
package main

import (
    gsession "github.com/snail007/gmc/http/session"
)

func main() {
    cfg := gsession.NewFileStoreConfig()
    cfg.TTL = 3600
    cfg.Dir = "./sessions"
    cfg.Prefix = "sess_"
    cfg.GCtime = 60
    
    store, err := gsession.NewFileStore(cfg)
    if err != nil {
        panic(err)
    }
    
    // 使用 store...
}
```

##### Redis Store

```go
package main

import (
    gsession "github.com/snail007/gmc/http/session"
    "time"
)

func main() {
    cfg := gsession.NewRedisStoreConfig()
    cfg.RedisCfg.Addr = "127.0.0.1:6379"
    cfg.RedisCfg.Password = ""
    cfg.RedisCfg.DBNum = 0
    cfg.RedisCfg.Prefix = "sess:"
    cfg.RedisCfg.Timeout = 10 * time.Second
    cfg.RedisCfg.MaxIdle = 10
    cfg.RedisCfg.MaxActive = 100
    cfg.TTL = 3600
    
    store, err := gsession.NewRedisStore(cfg)
    if err != nil {
        panic(err)
    }
    
    // 使用 store...
}
```

## 使用场景

### 场景 1：用户认证

```go
type AuthController struct {
    gmc.Controller
}

func (c *AuthController) Login() {
    username := c.PostForm("username")
    password := c.PostForm("password")
    
    // 验证凭据
    user := verifyCredentials(username, password)
    if user == nil {
        c.WriteHeader(401)
        c.JSON(gcore.M{"error": "Invalid credentials"})
        c.Stop()
        return
    }
    
    // 创建会话
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    // 保存用户信息到会话
    c.Session.Set("user_id", user.ID)
    c.Session.Set("username", user.Username)
    c.Session.Set("roles", user.Roles)
    c.Session.Set("login_time", time.Now().Unix())
    
    c.JSON(gcore.M{
        "status":  "success",
        "user_id": user.ID,
    })
}

func (c *AuthController) CheckAuth() bool {
    err := c.SessionStart()
    if err != nil {
        return false
    }
    
    userID := c.Session.Get("user_id")
    return userID != nil
}
```

### 场景 2：购物车

```go
type CartController struct {
    gmc.Controller
}

func (c *CartController) AddItem() {
    productID := c.PostFormInt("product_id")
    quantity := c.PostFormInt("quantity", 1)
    
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    // 获取购物车
    cart := c.Session.Get("cart")
    if cart == nil {
        cart = make(map[int]int)
    }
    
    cartMap := cart.(map[int]int)
    cartMap[productID] += quantity
    
    // 保存购物车
    c.Session.Set("cart", cartMap)
    
    c.JSON(gcore.M{
        "status": "success",
        "cart":   cartMap,
    })
}

func (c *CartController) GetCart() {
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    cart := c.Session.Get("cart")
    if cart == nil {
        cart = make(map[int]int)
    }
    
    c.JSON(gcore.M{"cart": cart})
}

func (c *CartController) ClearCart() {
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    c.Session.Delete("cart")
    c.JSON(gcore.M{"status": "cleared"})
}
```

### 场景 3：临时数据存储

```go
type WizardController struct {
    gmc.Controller
}

// 多步骤表单 - 第一步
func (c *WizardController) Step1() {
    if c.Request.Method == "POST" {
        err := c.SessionStart()
        if err != nil {
            c.Stop(err)
            return
        }
        
        // 保存第一步的数据
        c.Session.Set("step1_data", map[string]string{
            "name":  c.PostForm("name"),
            "email": c.PostForm("email"),
        })
        
        c.Redirect("/wizard/step2", 302)
        return
    }
    
    c.View.Render("wizard/step1")
}

// 第二步
func (c *WizardController) Step2() {
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    // 获取第一步的数据
    step1Data := c.Session.Get("step1_data")
    if step1Data == nil {
        c.Redirect("/wizard/step1", 302)
        return
    }
    
    if c.Request.Method == "POST" {
        // 保存第二步数据
        c.Session.Set("step2_data", map[string]string{
            "address": c.PostForm("address"),
            "phone":   c.PostForm("phone"),
        })
        
        c.Redirect("/wizard/confirm", 302)
        return
    }
    
    c.View.Set("step1", step1Data)
    c.View.Render("wizard/step2")
}

// 确认和提交
func (c *WizardController) Confirm() {
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    step1Data := c.Session.Get("step1_data")
    step2Data := c.Session.Get("step2_data")
    
    if step1Data == nil || step2Data == nil {
        c.Redirect("/wizard/step1", 302)
        return
    }
    
    if c.Request.Method == "POST" {
        // 处理完整数据
        processWizardData(step1Data, step2Data)
        
        // 清理临时数据
        c.Session.Delete("step1_data")
        c.Session.Delete("step2_data")
        
        c.Redirect("/wizard/success", 302)
        return
    }
    
    c.View.Set("step1", step1Data)
    c.View.Set("step2", step2Data)
    c.View.Render("wizard/confirm")
}
```

### 场景 4：Flash 消息

```go
type FlashHelper struct {
    session gcore.Session
}

func (f *FlashHelper) SetFlash(key, message string) {
    f.session.Set("flash_"+key, message)
}

func (f *FlashHelper) GetFlash(key string) string {
    msg := f.session.Get("flash_" + key)
    if msg != nil {
        f.session.Delete("flash_" + key)
        return msg.(string)
    }
    return ""
}

// 在控制器中使用
func (c *UserController) Create() {
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    // 保存用户...
    
    // 设置成功消息
    c.Session.Set("flash_success", "User created successfully")
    
    c.Redirect("/users", 302)
}

func (c *UserController) List() {
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    // 获取并显示 flash 消息
    successMsg := c.Session.Get("flash_success")
    if successMsg != nil {
        c.View.Set("success", successMsg)
        c.Session.Delete("flash_success")
    }
    
    c.View.Render("users/list")
}
```

## 性能考虑

### Memory Store

**优点：**
- 最快的读写性能
- 零外部依赖

**缺点：**
- 服务重启会丢失数据
- 不支持分布式部署
- 内存占用随会话数量增长

**适用场景：**
- 单机部署
- 开发环境
- 会话数量可控

### File Store

**优点：**
- 服务重启后数据保留
- 无外部依赖

**缺点：**
- I/O 性能瓶颈
- 不支持分布式部署
- 需要定期清理文件

**适用场景：**
- 单机部署
- 中小型应用
- 需要持久化会话

### Redis Store

**优点：**
- 支持分布式部署
- 高性能
- 数据持久化
- 支持集群

**缺点：**
- 需要 Redis 服务
- 网络延迟

**适用场景：**
- 分布式部署
- 大型应用
- 高并发场景

## 最佳实践

### 1. 会话安全

```go
// ✅ 登录后重新生成会话 ID
func (c *UserController) Login() {
    // 验证用户...
    
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    // 销毁旧会话
    c.SessionDestroy()
    
    // 重新启动新会话
    err = c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    c.Session.Set("user_id", userID)
}
```

### 2. 设置合理的 TTL

```toml
[session]
# 根据实际需求设置
ttl = 1800  # 30 分钟（常规应用）
# ttl = 3600  # 1 小时（需要较长会话）
# ttl = 7200  # 2 小时（管理后台）
```

### 3. 使用类型断言

```go
// ✅ 安全的类型断言
userID := c.Session.Get("user_id")
if userID != nil {
    if id, ok := userID.(int); ok {
        // 使用 id
    }
}

// ❌ 不安全
id := c.Session.Get("user_id").(int) // 可能 panic
```

### 4. 及时销毁会话

```go
func (c *UserController) Logout() {
    err := c.SessionStart()
    if err == nil {
        c.SessionDestroy()
    }
    c.Redirect("/login", 302)
}
```

### 5. 在模板中访问会话数据

会话数据会自动注入到模板变量 `.S` 中：

```html
<!-- 模板文件 -->
{{if .S.username}}
    <p>Welcome, {{.S.username}}</p>
{{else}}
    <p>Please login</p>
{{end}}
```

**注意：** 必须先调用 `SessionStart()` 才能在模板中访问。

## 注意事项

1. **必须先调用 SessionStart()**
   - 在访问 Session 对象之前必须调用 `SessionStart()`
   - 每个请求只需调用一次

2. **会话 Cookie**
   - 会话 ID 通过 Cookie 传递
   - Cookie 名称为 "gmcsessionid"
   - 默认为会话 Cookie（浏览器关闭即失效）

3. **数据类型**
   - 会话可以存储任何 Go 类型
   - 获取时需要进行类型断言
   - 复杂类型会自动序列化

4. **并发安全**
   - Session 对象内部有锁保护
   - 多个 goroutine 可以安全访问

5. **分布式部署**
   - Memory 和 File Store 不支持分布式
   - 分布式环境必须使用 Redis Store

## 测试覆盖率

```text
# go test -cover -v
=== RUN   TestNewSession
--- PASS: TestNewSession (1.00s)
=== RUN   TestSerialize
--- PASS: TestSerialize (0.00s)
=== RUN   TestUnserialize
--- PASS: TestUnserialize (0.00s)
PASS
coverage: 95.6% of statements
ok  	github.com/snail007/gmc/http/session	1.017s
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [GMC Controller](../controller/README.md)
- [GMC HTTP Server](../server/README.md)
- [GMC Cookie](../cookie/README.md)
