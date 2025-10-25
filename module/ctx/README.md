# GMC Ctx 模块

## 简介

GMC Ctx（Context）模块提供 HTTP 请求上下文，封装了请求和响应操作，是 Web 和 API 开发的核心组件。

## 功能特性

- **请求处理**：获取请求参数、headers、cookies 等
- **响应输出**：支持 JSON、XML、HTML、文件等多种响应格式
- **会话管理**：集成会话支持
- **模板渲染**：集成模板引擎
- **文件上传**：处理文件上传
- **国际化**：集成 i18n 支持
- **元数据存储**：存储请求级别的临时数据
- **分页支持**：内置分页工具

## 快速开始

### 获取请求参数

```go
func Handler(ctx gcore.Ctx) {
    // GET 参数
    name := ctx.Query("name", "default")
    
    // POST 参数
    email := ctx.POST("email", "")
    
    // 路由参数
    id := ctx.Param("id")
    
    // 所有参数（优先级：POST > GET > 路由）
    value := ctx.GetParam("key")
}
```

### 输出响应

```go
func Handler(ctx gcore.Ctx) {
    // JSON 响应
    ctx.WriteJSON(gcore.M{
        "status": "success",
        "data":   []string{"item1", "item2"},
    })
    
    // 文本响应
    ctx.Write("Hello World")
    
    // HTML 响应
    ctx.WriteHTML("<h1>Title</h1>")
    
    // XML 响应
    ctx.WriteXML(data)
    
    // 文件下载
    ctx.WriteFile("path/to/file.pdf", "downloaded.pdf")
}
```

### 模板渲染

```go
func Handler(ctx gcore.Ctx) {
    data := gcore.M{
        "title": "My Page",
        "user":  "John",
    }
    
    // 渲染模板
    ctx.View("template.html", data)
    
    // 或
    ctx.Template().Display("template.html", data)
}
```

### 获取请求信息

```go
func Handler(ctx gcore.Ctx) {
    // 请求方法
    method := ctx.Method()
    
    // 请求路径
    path := ctx.Request().URL.Path
    
    // Headers
    userAgent := ctx.GetHeader("User-Agent")
    
    // Cookies
    value, _ := ctx.GetCookie("session_id")
    
    // 客户端 IP
    ip := ctx.ClientIP()
    
    // 请求体
    body, _ := ctx.Body()
}
```

### 文件上传

```go
func Handler(ctx gcore.Ctx) {
    // 获取上传的文件
    file, header, err := ctx.GetPostFile("upload")
    if err != nil {
        ctx.WriteE(err)
        return
    }
    defer file.Close()
    
    // 保存文件
    err = ctx.SaveFile(header, "./uploads/"+header.Filename)
    if err != nil {
        ctx.WriteE(err)
        return
    }
    
    ctx.WriteJSON(gcore.M{
        "status":   "success",
        "filename": header.Filename,
    })
}
```

### 会话操作

```go
func Handler(ctx gcore.Ctx) {
    sess := ctx.Session()
    
    // 设置会话
    sess.Set("user_id", 123)
    
    // 获取会话
    userID := sess.Get("user_id")
    
    // 删除会话
    sess.Delete("temp_key")
    
    // 销毁会话
    sess.Destroy()
}
```

## API 参考

### 请求方法

- `Request() *http.Request`：获取原始请求
- `Response() http.ResponseWriter`：获取原始响应
- `Method() string`：获取请求方法
- `Query(key, defaultValue) string`：获取 GET 参数
- `POST(key, defaultValue) string`：获取 POST 参数
- `Param(key) string`：获取路由参数
- `GetParam(key) string`：获取参数（POST > GET > 路由）
- `Body() ([]byte, error)`：获取请求体
- `BodyJSON(v interface{}) error`：解析 JSON 请求体
- `GetHeader(key) string`：获取 Header
- `GetCookie(name) (string, error)`：获取 Cookie
- `ClientIP() string`：获取客户端 IP

### 响应方法

- `Write(data ...interface{})`：输出文本
- `WriteE(err error)`：输出错误
- `WriteJSON(data interface{})`：输出 JSON
- `WriteJSONP(callback string, data interface{})`：输出 JSONP
- `WriteXML(data interface{})`：输出 XML
- `WriteHTML(html string)`：输出 HTML
- `WriteFile(filepath, downloadName)`：下载文件
- `SetHeader(key, value string)`：设置响应 Header
- `SetCookie(cookie *http.Cookie)`：设置 Cookie
- `Redirect(url string, code int)`：重定向
- `StatusCode(code int)`：设置状态码

### 模板和国际化

- `View(tpl string, data)`：渲染模板
- `Template() gcore.Template`：获取模板引擎
- `I18n() gcore.I18n`：获取国际化工具
- `Tr(lang, key, defaultMsg)`：翻译文本

### 文件上传

- `GetPostFile(name) (multipart.File, *multipart.FileHeader, error)`
- `SaveFile(header, dst) error`

### 其他

- `Session() gcore.Session`：获取会话
- `NewPager(count, perPage) *gutil.Pager`：创建分页器
- `Logger() gcore.Logger`：获取日志记录器
- `Config() gcore.Config`：获取配置
- `App() gcore.App`：获取应用实例
- `Metadata() *gmap.Map`：获取元数据存储

## 使用场景

1. **Web 应用**：处理 HTTP 请求和响应
2. **RESTful API**：构建 API 接口
3. **文件上传**：处理文件上传功能
4. **用户认证**：基于会话的用户认证
5. **多语言支持**：国际化应用

## 最佳实践

### 1. 参数验证

```go
func Handler(ctx gcore.Ctx) {
    email := ctx.POST("email", "")
    if email == "" {
        ctx.WriteJSON(gcore.M{
            "error": "email is required",
        })
        return
    }
    
    // 处理业务逻辑
}
```

### 2. 统一错误处理

```go
func Handler(ctx gcore.Ctx) {
    data, err := doSomething()
    if err != nil {
        ctx.WriteE(err)
        return
    }
    
    ctx.WriteJSON(data)
}
```

### 3. 使用元数据传递数据

```go
// 在中间件中设置
func AuthMiddleware(ctx gcore.Ctx, next func()) {
    user := authenticate(ctx)
    ctx.Metadata().Store("user", user)
    next()
}

// 在处理器中使用
func Handler(ctx gcore.Ctx) {
    user, _ := ctx.Metadata().Load("user")
    // 使用 user
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
