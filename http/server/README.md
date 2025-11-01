# GMC HTTP Server 模块

## 简介

GMC HTTP Server 模块提供了两种核心服务器实现：

-   **`HTTPServer`**: 一个功能完备的 Web 服务器，适用于构建传统的、包含页面渲染的 MVC 应用。它内置了模板引擎、会话管理、静态文件服务等功能。
-   **`APIServer`**: 一个轻量级的 API 服务器，移除了模板和会话等功能，专注于提供高性能的 RESTful API 服务。

两者共享相同的底层请求处理流程和中间件架构。

## 核心概念：请求生命周期与中间件

无论是 `HTTPServer` 还是 `APIServer`，一个 HTTP 请求的处理都遵循着一个清晰、分层的生命周期。理解这个流程对于高效使用中间件至关重要。

### 请求处理流程图

```
+--------------------------------------------------------------------+
|                                                                    |
|    请求进入 -> [ Middleware0 ] -> 路由匹配 -> [ Middleware1 ]      |
|                                    |                               |
|                                    |                               |
|                                    |         +---------------------v-+
|                                    |         |                      |
|                                    |         |   执行控制器/处理器    |
|                                    |         |                      |
|                                    |         +-----------+-----------+
|                                    |                     |            
|                                    |                     |            
|   [ 404 处理器 ] <-----------------+---------------------+            
|        ^                                               |            
|        |                                               v            
|        +------------ [ 未找到路由 ] <------------+     [ Middleware2 ]
|                                                                    |
|                                                                    |
| <----------------------------- [ Middleware3 ] <-------------------+
|                                (defer中执行, 总会执行)               |
|                                                                    |
+--------------------------------------------------------------------+
```

### 各阶段详解

1.  **`Middleware0`**: **路由前中间件**
    -   **触发时机**：服务器接收到请求后，在进行任何路由匹配之前。
    -   **主要用途**：执行最前置的全局逻辑，如 IP 黑名单、全局限流、CORS 预检请求处理等。如果此中间件返回 `true`，请求将直接中断，后续所有步骤（包括 `Middleware3`）都不会执行。

2.  **路由匹配 (Router Lookup)**
    -   框架根据请求的 `METHOD` 和 `PATH` 查找匹配的路由规则。
    -   如果**未找到**匹配的路由，则直接跳转到 **404 处理器**。
    -   如果**找到**匹配的路由，则继续下一步。

3.  **`Middleware1`**: **路由后、处理器前中间件**
    -   **触发时机**：路由匹配成功后，在执行用户定义的控制器方法或处理器之前。
    -   **主要用途**：执行需要路由信息的逻辑，最典型的场景是**身份认证**和**权限校验**。因为此时已经知道请求要访问哪个具体的路由，可以进行精确的权限控制。

4.  **执行控制器/处理器 (Handler Execution)**
    -   执行用户在路由中注册的最终处理逻辑。
    -   如果在此过程中发生 `panic`，执行将中断，并跳转到 **500 错误处理器**。

5.  **`Middleware2`**: **处理器后中间件**
    -   **触发时机**：在控制器/处理器**成功执行完毕**后（即没有发生 `panic`）。
    -   **主要用途**：对成功的请求进行后置处理，例如修改响应内容、添加统一的响应头等。

6.  **`Middleware3`**: **最终中间件**
    -   **触发时机**：在 `defer` 语句中执行，因此**总是在请求的最后阶段被调用**，无论请求是成功、404 还是 500。
    -   **主要用途**：进行最终的清理工作或日志记录，最典型的场景是**记录访问日志**（Access Log），因为它能获取到请求的完整信息，包括最终的响应状态码和处理耗时。

7.  **错误处理器**
    -   **404 处理器**: 当路由未匹配时触发。可通过 `SetNotFoundHandler` 自定义。
    -   **500 处理器**: 当执行处理器时发生 `panic` 触发。可通过 `SetErrorHandler` 自定义。

---

## `HTTPServer` 详解

`HTTPServer` 是为构建传统 Web 应用设计的全功能服务器。

### 生命周期与初始化

1.  **创建**: `gmc.New.HTTPServer(ctx)`
2.  **`Init(cfg)`**: 由 `gmc.App` 调用，此方法会根据 `app.toml` 的配置自动初始化以下组件：
    -   `[httpserver]`：服务器监听地址、TLS 等。
    -   `[template]`：模板引擎。
    -   `[session]`：会话管理器。
    -   `[static]`：静态文件服务。
    -   `[i18n]`：国际化支持。
3.  **`Start()`**: 启动服务器，开始监听端口。
4.  **`Stop()` / `GracefulStop()`**: 关闭服务器。

### 核心钩子与示例

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
    "strings"
)

func main() {
    s := gmc.New.HTTPServer(gmc.New.CtxDefault())

    // 1. 添加模板函数
    s.AddFuncMap(map[string]interface{}{
        "upper": func(s string) string {
            return strings.ToUpper(s)
        },
    })

    // 2. 自定义 404 处理器
    s.SetNotFoundHandler(func(c gmc.C, tpl gcore.Template) {
        c.WriteHeader(404)
        c.Write("Custom 404 Page: " + c.Request().URL.Path)
    })

    // 3. 自定义 500 处理器
    s.SetErrorHandler(func(c gmc.C, tpl gcore.Template, err interface{}) {
        c.WriteHeader(500)
        c.Write(fmt.Sprintf("Custom 500 Page, Error: %v", err))
    })

    // 4. 添加各级中间件
    s.AddMiddleware0(func(c gmc.C) bool { // CORS 或全局日志
        fmt.Println("Middleware 0")
        return false
    })
    s.AddMiddleware1(func(c gmc.C) bool { // 身份认证
        fmt.Println("Middleware 1")
        return false
    })
    s.AddMiddleware2(func(c gmc.C) bool { // 成功响应后处理
        fmt.Println("Middleware 2")
        return false
    })
    s.AddMiddleware3(func(c gmc.C) bool { // 最终日志记录
        fmt.Println("Middleware 3")
        return false
    })

    // 5. 注册路由
    s.Router().GET("/", func(c gmc.C) {
        c.Write("Hello from HTTPServer!")
    })
    s.Router().GET("/panic", func(c gmc.C) {
        panic("test error")
    })

    // 6. 启动服务
    s.Listen(":7080")
    select {}
}
```

---

## `APIServer` 详解

`APIServer` 是为构建纯 API 服务而优化的轻量级服务器。

### 生命周期与初始化

`APIServer` 的生命周期与 `HTTPServer` 类似，但 `Init(cfg)` 方法是一个空实现，其配置主要在创建时通过 `gmc.New.APIServerDefault(ctx, cfg)` 或 `gmc.New.APIServer(ctx, address)` 完成，只关心 `[apiserver]` 配置块。

### 核心钩子与示例

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

func main() {
    api, _ := gmc.New.APIServer(gmc.New.Ctx(), ":7081")

    // 1. 自定义 404 处理器
    api.SetNotFoundHandler(func(c gmc.C) {
        c.WriteHeader(404)
        c.JSON(gcore.M{"error": "Not Found"})
    })

    // 2. 自定义 500 处理器
    api.SetErrorHandler(func(c gmc.C, err interface{}) {
        c.WriteHeader(500)
        c.JSON(gcore.M{"error": fmt.Sprintf("Internal Server Error: %v", err)})
    })

    // 3. 添加中间件 (与 HTTPServer 相同)
    api.AddMiddleware1(func(c gmc.C) bool {
        fmt.Println("API Auth Check")
        return false
    })

    // 4. 注册路由
    api.API("/hello", func(c gmc.C) {
        c.JSON(gcore.M{"message": "Hello from APIServer!"})
    })
    api.API("/panic", func(c gmc.C) {
        panic("api test error")
    })

    // 5. 启动服务
    api.Run()
    select {}
}
```

## 高级功能

### 优雅关闭与热重载

当 `HTTPServer` 或 `APIServer` 由 `gmc.App` 管理时，它们会自动支持优雅关闭和热重载。`gmc.App` 会在相应时机调用服务的 `GracefulStop()`、`Listeners()` 和 `InjectListeners()` 方法。

### 文件服务

两者都支持通过 `ServeFiles` 和 `ServeEmbedFS` 方法提供静态文件或嵌入式文件服务。

```go
// 提供本地 ./static 目录下的文件, URL 以 /s/ 开头
s.ServeFiles("./static", "/s")

// 提供嵌入式文件系统的文件, URL 以 /assets/ 开头
//go:embed assets/*
var assetsFS embed.FS
s.ServeEmbedFS(assetsFS, "/assets")
```

### 嵌入资源文件

GMC 支持使用 Go 1.16+ 的 `embed` 功能将资源文件打包到二进制中：

**静态文件和模板：**

```go
import "embed"

//go:embed static/*
var staticFS embed.FS

//go:embed views/*
var viewsFS embed.FS

func main() {
    s := ghttp.NewHTTPServer()
    s.ServeEmbedFS(staticFS, "/static")
    // 配置模板使用 viewsFS...
}
```

**i18n 文件：**

在 `i18n` 目录下创建 `i18n.go`：

```go
package i18n

import "embed"

//go:embed *.toml
var I18nFS embed.FS
```

在 `main.go` 中：

```go
import (
    "embed"
    gi18n "github.com/snail007/gmc/module/i18n"
    "myapp/i18n"  // 导入你的 i18n 包
)

func main() {
    // 初始化嵌入的 i18n 文件
    gi18n.InitEmbedFS(i18n.I18nFS, "zh-CN")
    
    s := ghttp.NewHTTPServer()
    s.Run(":8080")
}
```

查看 [i18n 模块文档](../../module/i18n/README.md) 了解更多详情。
