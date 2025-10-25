# GMC Cookie 包

## 简介

GMC Cookie 包提供了简单易用的 HTTP Cookie 操作接口，用于在 HTTP 请求和响应中设置、获取和删除 Cookie。

## 功能特性

- **设置 Cookie**：支持设置带各种选项的 Cookie
- **获取 Cookie**：从请求中读取 Cookie 值
- **删除 Cookie**：安全删除 Cookie
- **灵活配置**：支持 HttpOnly、Secure、Domain、Path、MaxAge 等选项

## 安装

```bash
go get github.com/snail007/gmc/http/cookie
```

## 快速开始

### 在 GMC Web 应用中使用

当使用 GMC 框架时，Cookie 功能已经集成在 Context 中：

```go
package main

import (
    "github.com/snail007/gmc"
)

type DemoController struct {
    gmc.Controller
}

func (c *DemoController) Index() {
    // 设置 Cookie
    c.Cookie().Set("username", "john")
    
    // 带选项设置 Cookie
    c.Cookie().Set("session_id", "abc123", &gcore.CookieOptions{
        Path:     "/",
        MaxAge:   3600,      // 1 小时
        HttpOnly: true,      // 防止 JavaScript 访问
        Secure:   true,      // 仅 HTTPS 传输
        Domain:   ".example.com",
    })
    
    // 获取 Cookie
    username, err := c.Cookie().Get("username")
    if err != nil {
        c.Write("Cookie not found")
        return
    }
    c.Write("Username: " + username)
    
    // 删除 Cookie
    c.Cookie().Remove("username")
}
```

### 直接使用 Cookie 包

如果需要在非 GMC 框架环境中使用：

```go
package main

import (
    "net/http"
    gcookie "github.com/snail007/gmc/http/cookie"
    gcore "github.com/snail007/gmc/core"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // 创建 Cookies 对象
    cookies := gcookie.New(w, r)
    
    // 设置 Cookie
    cookies.Set("user_id", "12345")
    
    // 带选项设置
    cookies.Set("token", "xyz789", &gcore.CookieOptions{
        Path:     "/api",
        MaxAge:   7200,
        HttpOnly: true,
    })
    
    // 获取 Cookie
    userId, err := cookies.Get("user_id")
    if err == nil {
        w.Write([]byte("User ID: " + userId))
    }
    
    // 删除 Cookie
    cookies.Remove("token")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

## API 参考

### 创建 Cookies 对象

```go
func New(w http.ResponseWriter, r *http.Request) *Cookies
```

### 方法

#### Set - 设置 Cookie

```go
func (c *Cookies) Set(name, val string, options ...*gcore.CookieOptions)
```

**参数：**
- `name`: Cookie 名称
- `val`: Cookie 值
- `options`: 可选的 Cookie 选项

**CookieOptions 结构：**
```go
type CookieOptions struct {
    Path     string  // Cookie 路径，默认 "/"
    Domain   string  // Cookie 域名
    MaxAge   int     // 最大生存时间（秒），0 表示会话 Cookie
    Secure   bool    // 是否仅通过 HTTPS 传输
    HTTPOnly bool    // 是否禁止 JavaScript 访问
}
```

#### Get - 获取 Cookie

```go
func (c *Cookies) Get(name string) (value string, err error)
```

**参数：**
- `name`: Cookie 名称

**返回：**
- `value`: Cookie 值
- `err`: 如果 Cookie 不存在，返回错误

#### Remove - 删除 Cookie

```go
func (c *Cookies) Remove(name string, options ...*gcore.CookieOptions)
```

**参数：**
- `name`: Cookie 名称
- `options`: 可选的 Cookie 选项（需要匹配 Path 和 Domain）

**注意：** 删除 Cookie 时，Path 和 Domain 必须与设置时相同。

## 使用示例

### 示例 1：用户登录状态

```go
func (c *LoginController) Login() {
    // 验证用户名密码...
    
    // 设置登录状态 Cookie（30 天过期）
    c.Cookie().Set("logged_in", "true", &gcore.CookieOptions{
        Path:     "/",
        MaxAge:   30 * 24 * 3600, // 30 天
        HttpOnly: true,
        Secure:   true,
    })
    
    // 设置用户信息
    c.Cookie().Set("username", username, &gcore.CookieOptions{
        Path:   "/",
        MaxAge: 30 * 24 * 3600,
    })
    
    c.Write("Login successful")
}

func (c *LoginController) Logout() {
    // 删除登录状态
    c.Cookie().Remove("logged_in", &gcore.CookieOptions{Path: "/"})
    c.Cookie().Remove("username", &gcore.CookieOptions{Path: "/"})
    
    c.Write("Logout successful")
}
```

### 示例 2：记住用户偏好

```go
func (c *SettingsController) SavePreference() {
    theme := c.Query("theme") // "dark" or "light"
    lang := c.Query("lang")   // "zh-CN", "en-US", etc.
    
    // 保存主题偏好（1 年）
    c.Cookie().Set("theme", theme, &gcore.CookieOptions{
        Path:   "/",
        MaxAge: 365 * 24 * 3600,
    })
    
    // 保存语言偏好（1 年）
    c.Cookie().Set("language", lang, &gcore.CookieOptions{
        Path:   "/",
        MaxAge: 365 * 24 * 3600,
    })
    
    c.JSON(gcore.M{"status": "ok"})
}

func (c *SettingsController) LoadPreference() {
    theme, _ := c.Cookie().Get("theme")
    if theme == "" {
        theme = "light" // 默认主题
    }
    
    lang, _ := c.Cookie().Get("language")
    if lang == "" {
        lang = "zh-CN" // 默认语言
    }
    
    c.JSON(gcore.M{
        "theme":    theme,
        "language": lang,
    })
}
```

### 示例 3：会话 Cookie

```go
func (c *ShoppingCartController) AddItem() {
    // 创建临时购物车 ID（会话 Cookie，浏览器关闭后失效）
    c.Cookie().Set("cart_id", generateCartID(), &gcore.CookieOptions{
        Path:     "/",
        MaxAge:   0, // 0 表示会话 Cookie
        HttpOnly: true,
    })
    
    c.Write("Item added to cart")
}
```

## 安全建议

1. **使用 HttpOnly**：对于敏感的 Cookie（如会话 ID），务必设置 `HttpOnly: true` 防止 XSS 攻击
2. **使用 Secure**：在生产环境（HTTPS）中，敏感 Cookie 应设置 `Secure: true`
3. **设置合适的 Path**：限制 Cookie 的可访问路径，避免不必要的暴露
4. **设置合适的 Domain**：仅在需要跨子域共享时才设置 Domain
5. **设置合理的过期时间**：根据实际需求设置 MaxAge，避免 Cookie 长期有效
6. **不要存储敏感信息**：Cookie 值可以被用户看到和修改，不要存储密码等敏感信息

## 默认配置

如果不提供 CookieOptions，将使用默认配置：

```go
var DefaultCookieOptions = &CookieOptions{
    Path:     "/",
    MaxAge:   0,     // 会话 Cookie
    Secure:   false,
    HTTPOnly: false,
    Domain:   "",
}
```

## 注意事项

1. **删除 Cookie**：删除 Cookie 时，Path 和 Domain 必须与设置时完全一致
2. **Cookie 大小限制**：单个 Cookie 大小不应超过 4KB
3. **Cookie 数量限制**：每个域名的 Cookie 数量有限制（通常约 50 个）
4. **编码问题**：Cookie 值会自动进行 URL 编码，无需手动编码

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [GMC Session 模块](../session/README.md)
- [GMC HTTP Server](../server/README.md)
