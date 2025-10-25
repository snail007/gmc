# GMC Controller

## 简介

GMC Controller 是 MVC 模式中的控制器组件，负责处理 HTTP 请求、调用业务逻辑、渲染视图等。控制器提供了丰富的辅助方法和自动初始化的对象，简化 Web 开发。

## 功能特性

- **自动路由绑定**：方法名自动映射为路由路径
- **生命周期钩子**：`Before()` 和 `After()` 方法
- **内置对象**：Request、Response、Session、Cookie 等自动初始化
- **便捷方法**：输出、跳转、错误处理等辅助方法
- **会话管理**：内置会话支持
- **模板渲染**：简化视图渲染
- **国际化支持**：内置 i18n 功能
- **流程控制**：`Stop()`、`Die()` 等流程控制方法

## 安装

```bash
go get github.com/snail007/gmc
```

## 快速开始

### 基本控制器

```go
package controller

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type UserController struct {
    gmc.Controller
}

// 方法会自动映射为路由
// 访问: /user/list
func (c *UserController) List() {
    c.Write("User List")
}

// 访问: /user/detail
func (c *UserController) Detail() {
    id := c.Query("id")
    c.JSON(gcore.M{
        "id":   id,
        "name": "Alice",
        "age":  25,
    })
}
```

### 路由绑定

```go
package main

import (
    "github.com/snail007/gmc"
    "myapp/controller"
)

func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    // 绑定控制器
    r := s.Router()
    r.Controller("/user", new(controller.UserController))
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

## 控制器生命周期

每个请求的处理流程：

```
1. MethodCallPre() - 框架自动调用，初始化内置对象
2. Before()        - 用户定义，前置处理（如权限验证）
3. YourMethod()    - 用户定义，实际业务方法
4. After()         - 用户定义，后置处理（如日志记录）
5. MethodCallPost()- 框架自动调用，清理工作
```

### Before() - 前置方法

在实际方法执行前调用，常用于：
- 权限验证
- 参数预处理
- 初始化数据

```go
type AdminController struct {
    gmc.Controller
}

func (c *AdminController) Before() {
    // 权限验证
    token := c.Header("Authorization")
    if token == "" {
        c.WriteHeader(401)
        c.JSON(gcore.M{"error": "Unauthorized"})
        c.Stop() // 停止后续执行
        return
    }
    
    // 验证 token...
    if !isValidToken(token) {
        c.WriteHeader(403)
        c.JSON(gcore.M{"error": "Forbidden"})
        c.Stop()
        return
    }
}

func (c *AdminController) Dashboard() {
    // 只有通过 Before() 验证才会执行
    c.Write("Admin Dashboard")
}
```

### After() - 后置方法

在实际方法执行后调用，常用于：
- 日志记录
- 资源清理
- 统一响应处理

```go
func (c *UserController) After() {
    // 记录访问日志
    c.Logger().Infof("User %s accessed %s", 
        c.Session.Get("username"), 
        c.Request.URL.Path)
    
    // 添加统一响应头
    c.SetHeader("X-Response-Time", 
        fmt.Sprintf("%dms", c.Ctx.TimeUsed()/1000000))
}
```

## 内置对象

控制器自动初始化以下对象，可直接使用：

```go
type Controller struct {
    Response     http.ResponseWriter    // HTTP 响应对象
    Request      *http.Request         // HTTP 请求对象
    Param        gcore.Params          // 路由参数
    Session      gcore.Session         // 会话对象
    Tpl          gcore.Template        // 模板引擎
    I18n         gcore.I18n            // 国际化
    SessionStore gcore.SessionStorage  // 会话存储
    Router       gcore.HTTPRouter      // 路由器
    Config       gcore.Config          // 配置对象
    Cookie       gcore.Cookies         // Cookie 操作
    Ctx          gcore.Ctx             // 请求上下文
    View         gcore.View            // 视图对象
    Lang         string                // 当前语言
    Logger       gcore.Logger          // 日志对象
}
```

## 辅助方法

### 输出方法

#### Write() - 输出内容

```go
func (c *UserController) Hello() {
    c.Write("Hello, World!")
    c.Write("Multiple ", "arguments ", "supported")
}
```

#### WriteE() - 带错误处理的输出

```go
func (c *UserController) Data() {
    n, err := c.WriteE("Some data")
    if err != nil {
        c.Logger().Error(err)
    }
}
```

#### JSON() - 输出 JSON

```go
func (c *UserController) GetUser() {
    c.JSON(gcore.M{
        "id":   123,
        "name": "Alice",
        "age":  25,
    })
}
```

#### JSONP() - 输出 JSONP

```go
func (c *UserController) GetUserJSONP() {
    callback := c.Query("callback")
    c.JSONP(callback, gcore.M{
        "id":   123,
        "name": "Alice",
    })
}
```

### 流程控制

#### Stop() - 停止执行

停止当前方法执行，但会继续执行 `After()` 方法：

```go
func (c *UserController) CheckPermission() {
    if !c.hasPermission() {
        c.Write("Access Denied")
        c.Stop() // 停止当前方法，但 After() 仍会执行
        return
    }
    
    // 这里不会执行
    c.Write("Welcome")
}
```

#### Die() - 立即终止

立即终止，不执行 `After()` 方法：

```go
func (c *UserController) Fatal() {
    c.Write("Fatal error occurred")
    c.Die() // 立即终止，After() 不会执行
}
```

#### StopE() - 错误检查

错误检查辅助方法，简化错误处理：

```go
func (c *UserController) CreateUser() {
    err := saveUser()
    
    // 如果 err 不为 nil，执行失败回调并停止
    c.StopE(err, func() {
        c.JSON(gcore.M{"error": err.Error()})
    }, func() {
        // err 为 nil 时执行成功回调
        c.JSON(gcore.M{"status": "success"})
    })
}
```

### 请求数据获取

#### Query() - 获取 URL 参数

```go
func (c *UserController) Search() {
    keyword := c.Query("keyword")
    page := c.QueryInt("page", 1)  // 默认值 1
    
    c.JSON(gcore.M{
        "keyword": keyword,
        "page":    page,
    })
}
```

#### PostForm() - 获取 POST 数据

```go
func (c *UserController) Login() {
    username := c.PostForm("username")
    password := c.PostForm("password")
    
    // 验证...
}
```

#### Param() - 获取路由参数

```go
// 路由: /user/:id
func (c *UserController) GetUserByID() {
    id := c.Param("id")
    c.Write("User ID: " + id)
}
```

#### Header() - 获取请求头

```go
func (c *UserController) CheckAuth() {
    token := c.Header("Authorization")
    userAgent := c.Header("User-Agent")
}
```

### 响应设置

#### WriteHeader() - 设置状态码

```go
func (c *UserController) NotFound() {
    c.WriteHeader(404)
    c.JSON(gcore.M{"error": "Not Found"})
}
```

#### SetHeader() - 设置响应头

```go
func (c *UserController) Download() {
    c.SetHeader("Content-Type", "application/octet-stream")
    c.SetHeader("Content-Disposition", "attachment; filename=file.zip")
    // 输出文件内容...
}
```

#### Redirect() - 重定向

```go
func (c *UserController) Old() {
    c.Redirect("/new-url", 302)
}
```

### 会话管理

#### SessionStart() - 启动会话

```go
func (c *UserController) Login() {
    // 必须先启动会话
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    // 设置会话数据
    c.Session.Set("user_id", 123)
    c.Session.Set("username", "alice")
    
    c.JSON(gcore.M{"status": "logged in"})
}
```

#### SessionDestroy() - 销毁会话

```go
func (c *UserController) Logout() {
    err := c.SessionStart()
    if err != nil {
        c.Stop(err)
        return
    }
    
    // 销毁会话
    err = c.SessionDestroy()
    if err != nil {
        c.Stop(err)
        return
    }
    
    c.JSON(gcore.M{"status": "logged out"})
}
```

### 视图渲染

#### View.Render() - 渲染模板

```go
func (c *UserController) Profile() {
    // 设置模板变量
    c.View.Set("username", "Alice")
    c.View.Set("age", 25)
    
    // 渲染模板
    c.View.Render("user/profile")
}
```

#### View.Layout() - 使用布局

```go
func (c *UserController) Index() {
    c.View.Set("title", "Home Page")
    c.View.Layout("layout").Render("home")
}
```

### Cookie 操作

#### Cookie.Set() - 设置 Cookie

```go
func (c *UserController) SetCookie() {
    c.Cookie.Set("theme", "dark", &gcore.CookieOptions{
        MaxAge:   3600,
        Path:     "/",
        HttpOnly: true,
    })
}
```

#### Cookie.Get() - 获取 Cookie

```go
func (c *UserController) GetCookie() {
    theme, err := c.Cookie.Get("theme")
    if err == nil {
        c.Write("Theme: " + theme)
    }
}
```

### 国际化

#### Tr() - 翻译

```go
func (c *UserController) Welcome() {
    // 第一个参数是 key，第二个是默认文本
    msg := c.Tr("welcome.message", "Welcome to our site")
    c.Write(msg)
}
```

## 方法规则

### 路由映射规则

1. **方法名自动转小写**
   - `List()` → `/list`
   - `UserDetail()` → `/userdetail`

2. **特殊后缀的方法不会路由**
   - `Index__()` - 双下划线结尾
   - `Helper_()` - 单下划线结尾

3. **特殊方法名**
   - `Before()` - 前置方法，不会路由
   - `After()` - 后置方法，不会路由

### 示例

```go
type ProductController struct {
    gmc.Controller
}

// 会路由：/product/list
func (c *ProductController) List() {
    c.Write("Product List")
}

// 不会路由（双下划线）
func (c *ProductController) Helper__() {
    // 辅助方法，不对外暴露
}

// 不会路由（Before 是特殊方法）
func (c *ProductController) Before() {
    // 前置处理
}
```

## 完整示例

### 示例 1：用户管理控制器

```go
package controller

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type UserController struct {
    gmc.Controller
}

// 前置验证
func (c *UserController) Before() {
    // 检查登录状态
    err := c.SessionStart()
    if err != nil {
        c.Die(err)
        return
    }
    
    userID := c.Session.Get("user_id")
    if userID == nil {
        c.WriteHeader(401)
        c.JSON(gcore.M{"error": "Please login"})
        c.Stop()
        return
    }
}

// 用户列表
func (c *UserController) List() {
    page := c.QueryInt("page", 1)
    pageSize := c.QueryInt("page_size", 20)
    
    // 从数据库获取用户...
    users := getUserList(page, pageSize)
    
    c.JSON(gcore.M{
        "users": users,
        "page":  page,
        "total": getTotalUsers(),
    })
}

// 用户详情
func (c *UserController) Detail() {
    id := c.QueryInt("id", 0)
    if id == 0 {
        c.WriteHeader(400)
        c.JSON(gcore.M{"error": "Invalid user ID"})
        c.Stop()
        return
    }
    
    user := getUserByID(id)
    if user == nil {
        c.WriteHeader(404)
        c.JSON(gcore.M{"error": "User not found"})
        c.Stop()
        return
    }
    
    c.JSON(gcore.M{"user": user})
}

// 创建用户
func (c *UserController) Create() {
    username := c.PostForm("username")
    email := c.PostForm("email")
    
    // 验证输入
    if username == "" || email == "" {
        c.WriteHeader(400)
        c.JSON(gcore.M{"error": "Username and email required"})
        c.Stop()
        return
    }
    
    // 创建用户
    err := createUser(username, email)
    c.StopE(err, func() {
        c.WriteHeader(500)
        c.JSON(gcore.M{"error": err.Error()})
    }, func() {
        c.JSON(gcore.M{"status": "success"})
    })
}

// 后置日志
func (c *UserController) After() {
    c.Logger().Infof("User action: %s, Time: %dms",
        c.Request.URL.Path,
        c.Ctx.TimeUsed()/1000000)
}
```

### 示例 2：博客控制器

```go
package controller

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type BlogController struct {
    gmc.Controller
}

// 文章列表页面
func (c *BlogController) Index() {
    page := c.QueryInt("page", 1)
    
    // 获取文章列表
    posts := getPostList(page, 10)
    
    // 设置模板变量
    c.View.Set("posts", posts)
    c.View.Set("page", page)
    c.View.Set("title", "博客首页")
    
    // 渲染模板
    c.View.Layout("layout").Render("blog/index")
}

// 文章详情页面
func (c *BlogController) Post() {
    id := c.QueryInt("id", 0)
    if id == 0 {
        c.Redirect("/blog", 302)
        return
    }
    
    post := getPostByID(id)
    if post == nil {
        c.WriteHeader(404)
        c.View.Layout("layout").Render("404")
        return
    }
    
    c.View.Set("post", post)
    c.View.Set("title", post.Title)
    c.View.Layout("layout").Render("blog/post")
}

// 创建文章（需要登录）
func (c *BlogController) Create() {
    // 检查登录
    err := c.SessionStart()
    if err != nil || c.Session.Get("user_id") == nil {
        c.Redirect("/login", 302)
        return
    }
    
    if c.Request.Method == "GET" {
        // 显示创建表单
        c.View.Layout("layout").Render("blog/create")
        return
    }
    
    // 处理 POST 请求
    title := c.PostForm("title")
    content := c.PostForm("content")
    
    err = createPost(title, content)
    c.StopE(err, func() {
        c.View.Set("error", err.Error())
        c.View.Layout("layout").Render("blog/create")
    }, func() {
        c.Redirect("/blog", 302)
    })
}
```

### 示例 3：API 控制器

```go
package controller

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type APIController struct {
    gmc.Controller
}

// 统一 JSON 响应
func (c *APIController) After() {
    // 如果还没有输出，添加默认响应头
    c.SetHeader("Content-Type", "application/json")
    c.SetHeader("X-API-Version", "1.0")
}

// 获取用户信息
func (c *APIController) User() {
    id := c.Query("id")
    
    user, err := getUserByID(id)
    if err != nil {
        c.WriteHeader(404)
        c.JSON(gcore.M{
            "error": "User not found",
            "code":  404,
        })
        c.Stop()
        return
    }
    
    c.JSON(gcore.M{
        "code": 200,
        "data": user,
    })
}

// 搜索
func (c *APIController) Search() {
    keyword := c.Query("q")
    page := c.QueryInt("page", 1)
    
    results := search(keyword, page)
    
    c.JSON(gcore.M{
        "code":    200,
        "keyword": keyword,
        "page":    page,
        "results": results,
    })
}
```

## 内部方法

以下方法由框架自动调用，**不要在代码中手动调用**：

- `MethodCallPre()` - 在 `Before()` 之前调用，初始化内置对象
- `MethodCallPost()` - 在 `After()` 之后调用，执行清理工作

## 最佳实践

### 1. 使用 Before() 进行权限验证

```go
func (c *AdminController) Before() {
    // 统一的权限验证
    if !c.isAdmin() {
        c.WriteHeader(403)
        c.JSON(gcore.M{"error": "Forbidden"})
        c.Stop()
    }
}
```

### 2. 使用 StopE() 简化错误处理

```go
// ❌ 不推荐
err := doSomething()
if err != nil {
    c.JSON(gcore.M{"error": err.Error()})
    c.Stop()
    return
}
c.JSON(gcore.M{"status": "ok"})

// ✅ 推荐
err := doSomething()
c.StopE(err, func() {
    c.JSON(gcore.M{"error": err.Error()})
}, func() {
    c.JSON(gcore.M{"status": "ok"})
})
```

### 3. 使用 After() 记录日志

```go
func (c *BaseController) After() {
    // 统一的日志记录
    c.Logger().Infof("%s %s %d %dms",
        c.Request.Method,
        c.Request.URL.Path,
        c.Ctx.StatusCode(),
        c.Ctx.TimeUsed()/1000000)
}
```

### 4. 合理使用 Stop() 和 Die()

- 使用 `Stop()` - 需要执行 `After()` 时（如记录日志）
- 使用 `Die()` - 发生致命错误时

### 5. 分离业务逻辑

```go
// ✅ 控制器只处理 HTTP 相关逻辑
func (c *UserController) Create() {
    username := c.PostForm("username")
    
    // 业务逻辑放在 service 层
    user, err := userService.CreateUser(username)
    
    c.StopE(err, func() {
        c.JSON(gcore.M{"error": err.Error()})
    }, func() {
        c.JSON(gcore.M{"user": user})
    })
}
```

## 注意事项

1. **不要定义与框架方法同名的方法**
2. **Before() 中使用 Stop() 会阻止实际方法执行**
3. **Die() 会跳过 After() 执行**
4. **SessionStart() 必须在访问 Session 之前调用**
5. **双下划线和单下划线结尾的方法不会路由**

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [GMC Router](../router/README.md)
- [GMC HTTP Server](../server/README.md)
- [GMC Session](../session/README.md)
- [GMC View](../view/README.md)
- [GMC Cookie](../cookie/README.md)
