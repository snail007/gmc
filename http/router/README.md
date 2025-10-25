# GMC HTTP Router

## 简介

GMC HTTP Router 是一个基于 [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) 的高性能 HTTP 路由器，提供了丰富的路由功能和控制器绑定支持。采用 Radix Tree 算法，性能优异。

**注意：** 本包基于 httprouter，使用时需遵守其 LICENSE。

## 功能特性

- **高性能**：基于 Radix Tree，O(log n) 复杂度
- **RESTful 支持**：支持所有 HTTP 方法（GET、POST、PUT、PATCH、DELETE 等）
- **路径参数**：支持命名参数（:name）和通配符（*name）
- **路由分组**：支持路由分组和命名空间
- **控制器绑定**：自动绑定控制器方法到路由
- **中间件支持**：多级中间件系统
- **方法后缀**：支持在 URL 中添加文件后缀（如 .html、.json）
- **自动重定向**：智能处理尾部斜杠和路径修正
- **静态文件**：内置静态文件服务支持

## 安装

```bash
go get github.com/snail007/gmc
```

## 快速开始

### 基本路由

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    // 获取路由器
    r := s.Router()
    
    // 基本路由
    r.HandlerFunc("GET", "/", func(c gmc.C) {
        c.Write("Welcome!")
    })
    
    r.HandlerFunc("GET", "/hello", func(c gmc.C) {
        c.Write("Hello, GMC!")
    })
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

### RESTful 路由

```go
r := s.Router()

// RESTful 方法
r.GET("/users", listUsers)
r.POST("/users", createUser)
r.GET("/users/:id", getUser)
r.PUT("/users/:id", updateUser)
r.DELETE("/users/:id", deleteUser)
r.PATCH("/users/:id", patchUser)

func getUser(c gmc.C) {
    id := c.Param("id")
    c.Write("User ID: " + id)
}
```

### 路径参数

```go
// 命名参数（匹配单个路径段）
r.GET("/blog/:category/:post", func(c gmc.C) {
    category := c.Param("category")
    post := c.Param("post")
    c.Write("Category: " + category + ", Post: " + post)
})

// 通配符参数（匹配所有剩余路径）
r.GET("/files/*filepath", func(c gmc.C) {
    filepath := c.Param("filepath")
    c.Write("File path: " + filepath)
})
```

**路径参数规则：**

| 路径模式 | 请求 URL | 是否匹配 | 参数 |
|---------|---------|---------|------|
| `/blog/:category/:post` | `/blog/go/tutorial` | ✅ | category="go", post="tutorial" |
| `/blog/:category/:post` | `/blog/go/` | ❌ | |
| `/blog/:category/:post` | `/blog/go` | ❌ | |
| `/files/*filepath` | `/files/doc/test.txt` | ✅ | filepath="/doc/test.txt" |
| `/files/*filepath` | `/files/` | ✅ | filepath="/" |

## 控制器绑定

### 基本控制器

```go
package main

import (
    "github.com/snail007/gmc"
)

type UserController struct {
    gmc.Controller
}

// 方法名会自动转换为小写路由
// 访问: /user/list
func (c *UserController) List() {
    c.JSON(gcore.M{
        "users": []string{"Alice", "Bob", "Charlie"},
    })
}

// 访问: /user/detail
func (c *UserController) Detail() {
    id := c.Query("id")
    c.JSON(gcore.M{
        "id":   id,
        "name": "Alice",
    })
}

func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    r := s.Router()
    
    // 绑定控制器到 /user 路径
    r.Controller("/user", new(UserController))
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

### 控制器方法规则

1. **方法名自动转换**：方法名自动转为小写作为路由路径
   - `List()` → `/list`
   - `GetUser()` → `/getuser`
   - `CreateOrder()` → `/createorder`

2. **特殊方法**：以下方法有特殊含义，不会被路由
   - `Before()`：在实际方法执行前调用（构造函数）
   - `After()`：在实际方法执行后调用（析构函数）
   - 以 `__` 或 `_` 结尾的方法：会被忽略

3. **所有 HTTP 方法**：控制器方法接受所有 HTTP 方法（GET、POST 等）

### 控制器生命周期

```go
type OrderController struct {
    gmc.Controller
}

// Before 在每个方法前执行
func (c *OrderController) Before() {
    // 身份验证
    token := c.Header("Authorization")
    if token == "" {
        c.WriteHeader(401)
        c.JSON(gcore.M{"error": "Unauthorized"})
        c.Stop() // 停止后续执行
    }
}

// After 在每个方法后执行
func (c *OrderController) After() {
    // 清理工作
    c.Logger().Info("Request completed")
}

func (c *OrderController) Create() {
    // 创建订单逻辑
    c.JSON(gcore.M{"status": "created"})
}
```

### 绑定单个方法

```go
r.ControllerMethod("/", new(HomeController), "Index")
r.ControllerMethod("/about", new(HomeController), "About")

// 访问 / 时调用 HomeController.Index()
// 访问 /about 时调用 HomeController.About()
```

## 路由分组

### 创建路由组

```go
r := s.Router()

// API 路由组
api := r.Group("/api")
{
    // /api/users
    api.GET("/users", listUsers)
    
    // /api/products
    api.Controller("/products", new(ProductController))
    
    // 嵌套分组
    v1 := api.Group("/v1")
    {
        // /api/v1/orders
        v1.Controller("/orders", new(OrderController))
    }
}

// 管理后台路由组
admin := r.Group("/admin")
{
    // /admin/dashboard
    admin.Controller("/dashboard", new(DashboardController))
    
    // /admin/users
    admin.Controller("/users", new(AdminUserController))
}
```

### 路由组示例

```go
func InitRouter(s *gmc.HTTPServer) {
    r := s.Router()
    
    // 前台路由
    r.Controller("/", new(HomeController))
    r.Controller("/user", new(UserController))
    
    // API v1
    v1 := r.Group("/api/v1")
    v1.Controller("/users", new(APIUserController))
    v1.Controller("/products", new(APIProductController))
    
    // API v2
    v2 := r.Group("/api/v2")
    v2.Controller("/users", new(APIUserV2Controller))
    
    // 管理后台
    admin := r.Group("/admin")
    admin.Controller("/dashboard", new(AdminDashboardController))
    admin.Controller("/settings", new(AdminSettingsController))
}
```

## URL 后缀

### 设置默认后缀

```go
r := s.Router()

// 为所有路由添加 .html 后缀
r.Ext(".html")

// 控制器绑定
r.Controller("/user", new(UserController))

// 现在访问时需要加 .html
// /user/list.html
// /user/detail.html
```

### 不同分组不同后缀

```go
// HTML 页面
web := r.Group("/web")
web.Ext(".html")
web.Controller("/user", new(WebUserController))
// 访问: /web/user/list.html

// JSON API
api := r.Group("/api")
api.Ext(".json")
api.Controller("/user", new(APIUserController))
// 访问: /api/user/list.json
```

## 静态文件服务

### 使用 ServeFiles

```go
r.ServeFiles("/static/*filepath", http.Dir("./public"))

// 访问:
// /static/css/style.css → ./public/css/style.css
// /static/js/app.js     → ./public/js/app.js
```

### 使用 HandlerFunc

```go
r.GET("/download/*filepath", func(c gmc.C) {
    filepath := c.Param("filepath")
    c.WriteHeader(200)
    c.SetHeader("Content-Type", "application/octet-stream")
    // 读取文件并发送...
})
```

## 高级功能

### HandlerAny - 匹配所有方法

```go
// 接受所有 HTTP 方法
r.HandlerAny("/api/webhook", func(c gmc.C) {
    method := c.Request().Method
    c.Write("Method: " + method)
})

// 等价于：
for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"} {
    r.HandlerFunc(method, "/api/webhook", handler)
}
```

### 路由调试

```go
// 打印所有注册的路由
r.PrintRouteTable(nil)

// 输出示例：
// GET    /
// GET    /user/list
// POST   /user/create
// GET    /api/v1/products
```

### URL 构建

```go
// 在控制器中构建 URL
url := c.Router().URL("/user/detail", gcore.M{
    "id": "123",
})
// 输出: /user/detail?id=123
```

## 配置选项

Router 支持多个配置选项：

```go
r := s.Router()

// 自动重定向尾部斜杠
r.RedirectTrailingSlash = true
// /foo/ 重定向到 /foo

// 自动修正路径
r.RedirectFixedPath = true
// /FOO 重定向到 /foo

// 处理 Method Not Allowed
r.HandleMethodNotAllowed = true
// 如果路由存在但方法不匹配，返回 405

// 自动处理 OPTIONS 请求
r.HandleOPTIONS = true
```

## 性能优化

1. **路由数量**：支持数千条路由而不影响性能
2. **参数提取**：使用对象池减少内存分配
3. **路径查找**：O(log n) 时间复杂度
4. **零内存分配**：路由匹配过程零内存分配（无参数路由）

## 路由匹配优先级

1. **静态路由** > **命名参数** > **通配符**

```go
r.GET("/users/admin", handler1)      // 优先级 1（静态）
r.GET("/users/:id", handler2)        // 优先级 2（命名参数）
r.GET("/users/*action", handler3)    // 优先级 3（通配符）

// 访问 /users/admin → handler1
// 访问 /users/123 → handler2
// 访问 /users/foo/bar → handler3
```

## API 参考

### HTTPRouter 主要方法

```go
// HTTP 方法路由
func (r *HTTPRouter) GET(path string, handle gcore.Handle)
func (r *HTTPRouter) POST(path string, handle gcore.Handle)
func (r *HTTPRouter) PUT(path string, handle gcore.Handle)
func (r *HTTPRouter) PATCH(path string, handle gcore.Handle)
func (r *HTTPRouter) DELETE(path string, handle gcore.Handle)
func (r *HTTPRouter) OPTIONS(path string, handle gcore.Handle)
func (r *HTTPRouter) HEAD(path string, handle gcore.Handle)

// 通用路由
func (r *HTTPRouter) Handle(method, path string, handle gcore.Handle)
func (r *HTTPRouter) HandlerFunc(method, path string, handler http.HandlerFunc)
func (r *HTTPRouter) HandlerAny(path string, handler interface{})

// 控制器绑定
func (r *HTTPRouter) Controller(urlPath string, controller interface{})
func (r *HTTPRouter) ControllerMethod(urlPath string, controller interface{}, method string)

// 路由分组
func (r *HTTPRouter) Group(namespace string) gcore.HTTPRouter

// 静态文件
func (r *HTTPRouter) ServeFiles(path string, root http.FileSystem)

// URL 后缀
func (r *HTTPRouter) Ext(ext string)

// 其他
func (r *HTTPRouter) PrintRouteTable(w io.Writer)
func (r *HTTPRouter) URL(path string, params gcore.M) string
```

## 使用示例

### 示例 1：博客应用

```go
type BlogController struct {
    gmc.Controller
}

func (c *BlogController) List() {
    // 列表页
    c.View("blog/list", nil)
}

func (c *BlogController) Detail() {
    // 详情页
    id := c.Query("id")
    c.View("blog/detail", gcore.M{"id": id})
}

func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    r := s.Router()
    r.Ext(".html")
    r.Controller("/blog", new(BlogController))
    
    // 访问:
    // /blog/list.html
    // /blog/detail.html?id=123
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

### 示例 2：RESTful API

```go
type APIController struct {
    gmc.Controller
}

func (c *APIController) List() {
    c.JSON(gcore.M{"users": []string{"Alice", "Bob"}})
}

func (c *APIController) Create() {
    // 创建资源
    c.JSON(gcore.M{"status": "created"})
}

func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    r := s.Router()
    api := r.Group("/api/v1")
    api.Ext(".json")
    api.Controller("/users", new(APIController))
    
    // 访问:
    // GET  /api/v1/users/list.json
    // POST /api/v1/users/create.json
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

### 示例 3：混合路由

```go
func InitRouter(s *gmc.HTTPServer) {
    r := s.Router()
    
    // 静态页面
    r.GET("/", homeHandler)
    r.GET("/about", aboutHandler)
    
    // 控制器路由
    r.Controller("/user", new(UserController))
    
    // API 路由
    api := r.Group("/api")
    {
        api.GET("/status", apiStatus)
        api.Controller("/products", new(ProductController))
    }
    
    // 静态文件
    r.ServeFiles("/static/*filepath", http.Dir("./public"))
    
    // 文件下载
    r.GET("/download/:filename", downloadHandler)
}
```

## 最佳实践

1. **分离路由配置**：将路由配置放在单独的 `router` 包中
2. **使用路由组**：按功能模块组织路由
3. **RESTful 设计**：遵循 RESTful API 设计规范
4. **控制器命名**：使用清晰的控制器和方法命名
5. **URL 设计**：使用语义化的 URL 路径
6. **参数验证**：在控制器中验证路径参数和查询参数

## 性能基准

基于 httprouter 的性能基准（供参考）：

```
BenchmarkRouter-8           20000000    75.8 ns/op    0 B/op    0 allocs/op
BenchmarkRouterWithParams-8 10000000    140 ns/op     32 B/op   1 allocs/op
```

## 注意事项

1. **路径冲突**：避免定义冲突的路由模式
   - ❌ `/users/:id` 和 `/users/:name`（冲突）
   - ✅ `/users/:id` 和 `/posts/:id`（不冲突）

2. **通配符位置**：通配符必须在路径末尾
   - ✅ `/files/*filepath`
   - ❌ `/files/*filepath/info`

3. **参数命名**：参数名必须唯一
   - ✅ `/blog/:category/:post`
   - ❌ `/blog/:id/:id`

4. **性能考虑**：静态路由比参数路由更快

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [httprouter 项目](https://github.com/julienschmidt/httprouter)
- [GMC HTTP Server](../server/README.md)
- [GMC Controller](../controller/README.md)

