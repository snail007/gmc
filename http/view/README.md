# GMC View 模块

## 简介

GMC View 模块负责视图渲染，提供模板数据绑定和渲染功能。支持布局（Layout）、数据传递、链式调用等特性，简化视图层开发。

## 功能特性

- **模板渲染**：渲染 HTML 模板并输出
- **数据绑定**：将数据传递给模板
- **布局支持**：支持模板布局（Layout）
- **链式调用**：支持方法链式调用
- **自动变量注入**：自动注入 GET、POST、Session、Cookie 等数据
- **错误处理**：内置错误处理机制

## 安装

```bash
go get github.com/snail007/gmc
```

## 快速开始

### 基本使用

```go
package controller

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type HomeController struct {
    gmc.Controller
}

func (c *HomeController) Index() {
    // 设置模板变量
    c.View.Set("title", "首页")
    c.View.Set("username", "Alice")
    
    // 渲染模板 views/home/index.html
    c.View.Render("home/index")
}
```

### 使用布局

```go
func (c *HomeController) About() {
    c.View.Set("title", "关于我们")
    c.View.Set("content", "这是关于页面")
    
    // 使用 layout.html 布局渲染 about.html
    c.View.Layout("layout").Render("about")
}
```

## 配置

### app.toml 模板配置

```toml
[template]
# 模板文件目录
dir = "views"

# 模板文件扩展名
ext = ".html"

# 模板定界符（左）
delimiterleft = "{{"

# 模板定界符（右）
delimiterright = "}}"

# 布局文件目录
layout = "layout"
```

## API 参考

### Set() - 设置单个变量

```go
func (v *View) Set(key string, val interface{}) gcore.View
```

**示例：**
```go
c.View.Set("title", "用户列表")
c.View.Set("users", []User{...})
c.View.Set("page", 1)
```

### SetMap() - 批量设置变量

```go
func (v *View) SetMap(d map[string]interface{}) gcore.View
```

**示例：**
```go
c.View.SetMap(gcore.M{
    "title": "产品详情",
    "product": product,
    "price": 99.99,
    "inStock": true,
})
```

### Render() - 渲染模板

```go
func (v *View) Render(tpl string, data ...map[string]interface{}) gcore.View
```

**参数：**
- `tpl`: 模板文件路径（相对于模板目录，不含扩展名）
- `data`: 可选的额外数据

**示例：**
```go
// 渲染 views/user/profile.html
c.View.Render("user/profile")

// 渲染时传递额外数据
c.View.Render("user/profile", gcore.M{
    "extra": "some data",
})
```

### RenderR() - 渲染并返回结果

```go
func (v *View) RenderR(tpl string, data ...map[string]interface{}) []byte
```

返回渲染后的字节数组，不直接输出。

**示例：**
```go
html := c.View.RenderR("email/welcome")
sendEmail(user.Email, html)
```

### Layout() - 设置布局

```go
func (v *View) Layout(l string) gcore.View
```

**示例：**
```go
// 使用 views/layout/main.html 作为布局
c.View.Layout("main").Render("home")

// 使用 views/layout/admin.html 作为布局
c.View.Layout("admin").Render("dashboard")
```

### Stop() - 停止渲染

```go
func (v *View) Stop()
```

停止控制器方法继续执行（等同于 `c.Stop()`）。

**示例：**
```go
c.View.Render("user/list").Stop()
// 后面的代码不会执行
```

### Err() - 获取错误

```go
func (v *View) Err() error
```

获取渲染过程中的错误。

**示例：**
```go
c.View.Render("some/template")
if err := c.View.Err(); err != nil {
    c.Logger().Error(err)
}
```

## 模板目录结构

推荐的模板目录结构：

```
views/
├── layout/              # 布局文件
│   ├── main.html       # 主布局
│   └── admin.html      # 管理后台布局
├── home/               # 首页模板
│   └── index.html
├── user/               # 用户相关模板
│   ├── list.html
│   ├── detail.html
│   └── edit.html
├── product/            # 产品相关模板
│   ├── list.html
│   └── detail.html
└── error/              # 错误页面
    ├── 404.html
    └── 500.html
```

## 模板语法

GMC 使用 Go 的 `text/template` 引擎，支持标准的模板语法。

### 基本语法

```html
<!-- 输出变量 -->
<h1>{{.title}}</h1>
<p>{{.content}}</p>

<!-- 条件判断 -->
{{if .isLoggedIn}}
    <p>欢迎，{{.username}}</p>
{{else}}
    <p>请登录</p>
{{end}}

<!-- 循环 -->
{{range .users}}
    <li>{{.Name}} - {{.Email}}</li>
{{end}}

<!-- 管道操作 -->
<p>{{.content | tohtml}}</p>
```

### 自动注入的变量

GMC 自动注入以下变量到模板：

#### .G - GET 参数

```html
<!-- 访问 URL: /search?keyword=go&page=1 -->
<p>关键词: {{.G.keyword}}</p>
<p>页码: {{.G.page}}</p>
```

#### .P - POST 参数

```html
<!-- POST 表单数据 -->
<p>用户名: {{.P.username}}</p>
<p>邮箱: {{.P.email}}</p>
```

#### .S - Session 数据

```html
<!-- Session 数据（需先调用 SessionStart） -->
{{if .S.user_id}}
    <p>用户ID: {{.S.user_id}}</p>
    <p>用户名: {{.S.username}}</p>
{{end}}
```

#### .C - Cookie 数据

```html
<!-- Cookie 数据 -->
<p>主题: {{.C.theme}}</p>
<p>语言: {{.C.language}}</p>
```

#### .U - URL 信息

```html
<!-- URL 信息 -->
<p>Host: {{.U.HOST}}</p>
<p>Path: {{.U.PATH}}</p>
<p>Scheme: {{.U.SCHEME}}</p>
<p>Query: {{.U.RAW_QUERY}}</p>
```

URL 变量包含的字段：
- `HOST`: 主机名（含端口）
- `HOSTNAME`: 主机名（不含端口）
- `PORT`: 端口号
- `PATH`: 路径
- `SCHEME`: 协议（http/https）
- `RAW_QUERY`: 原始查询字符串
- `URI`: 完整 URI
- `URL`: 完整 URL

#### .H - HTTP 头信息

```html
<!-- HTTP 头信息 -->
<p>User-Agent: {{.H.User-Agent}}</p>
<p>Accept-Language: {{.H.Accept-Language}}</p>
```

#### .Lang - 当前语言

```html
<!-- 当前语言（i18n） -->
<p>当前语言: {{.Lang}}</p>
```

### 内置模板函数

#### 标准函数

Go 模板引擎内置函数：

- **输出**: `print`, `printf`, `println`
- **HTML**: `html`, `js`, `urlquery`
- **比较**: `eq`, `ne`, `lt`, `le`, `gt`, `ge`
- **逻辑**: `and`, `or`, `not`
- **其他**: `call`, `index`, `len`, `slice`

#### GMC 扩展函数

##### tr - 国际化翻译

```html
{{tr .Lang "welcome.message" "欢迎访问"}}
```

- 第一个参数: `.Lang` 当前语言
- 第二个参数: 翻译键
- 第三个参数: 默认文本（可选）

##### trs - 翻译（返回 string）

类似 `tr`，但返回 string 类型而非 template.HTML。

```html
{{trs .Lang "button.submit" "提交"}}
```

##### string - 转换为字符串

```html
{{string .count}}
```

##### tohtml - 转换为 HTML

将字符串转换为 template.HTML 类型，不进行 HTML 转义。

```html
{{tohtml .htmlContent}}
```

##### val - 安全获取变量

获取模板变量，如果不存在返回空字符串而非 `<no value>`。

```html
{{val .maybeUndefined}}
```

### Sprig 函数库

GMC 模板包含 [Sprig](https://masterminds.github.io/sprig/) 函数库的子集，提供丰富的模板函数。

详细文档：[template/sprig/docs](../template/sprig/docs)

## 布局系统

### 布局文件示例

**views/layout/main.html**:
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.title}} - 我的网站</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <header>
        <h1>我的网站</h1>
        <nav>
            <a href="/">首页</a>
            <a href="/about">关于</a>
            {{if .S.user_id}}
                <a href="/profile">个人中心</a>
                <a href="/logout">退出</a>
            {{else}}
                <a href="/login">登录</a>
            {{end}}
        </nav>
    </header>
    
    <main>
        <!-- 内容区域，由 Render 的模板填充 -->
        {{.GMC_LAYOUT_CONTENT}}
    </main>
    
    <footer>
        <p>&copy; 2024 我的网站</p>
    </footer>
</body>
</html>
```

### 使用布局

**控制器**:
```go
func (c *HomeController) Index() {
    c.View.Set("title", "首页")
    c.View.Set("message", "欢迎访问")
    
    // 使用 main.html 布局
    c.View.Layout("main").Render("home/index")
}
```

**views/home/index.html**:
```html
<h2>{{.title}}</h2>
<p>{{.message}}</p>
<p>这是首页内容</p>
```

**最终输出**将是 main.html 布局，其中 `{{.GMC_LAYOUT_CONTENT}}` 被替换为 home/index.html 的内容。

## 完整示例

### 示例 1：用户列表

**控制器**:
```go
package controller

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type UserController struct {
    gmc.Controller
}

func (c *UserController) List() {
    page := c.QueryInt("page", 1)
    pageSize := 20
    
    // 从数据库获取用户
    users := getUserList(page, pageSize)
    total := getTotalUsers()
    
    // 设置模板变量
    c.View.Set("title", "用户列表")
    c.View.Set("users", users)
    c.View.Set("page", page)
    c.View.Set("pageSize", pageSize)
    c.View.Set("total", total)
    c.View.Set("totalPages", (total+pageSize-1)/pageSize)
    
    // 使用布局渲染
    c.View.Layout("main").Render("user/list")
}
```

**views/user/list.html**:
```html
<div class="user-list">
    <h2>{{.title}}</h2>
    
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>用户名</th>
                <th>邮箱</th>
                <th>操作</th>
            </tr>
        </thead>
        <tbody>
            {{range .users}}
            <tr>
                <td>{{.ID}}</td>
                <td>{{.Username}}</td>
                <td>{{.Email}}</td>
                <td>
                    <a href="/user/detail?id={{.ID}}">查看</a>
                    <a href="/user/edit?id={{.ID}}">编辑</a>
                </td>
            </tr>
            {{else}}
            <tr>
                <td colspan="4">暂无用户</td>
            </tr>
            {{end}}
        </tbody>
    </table>
    
    <!-- 分页 -->
    <div class="pagination">
        {{if gt .page 1}}
            <a href="?page={{sub .page 1}}">上一页</a>
        {{end}}
        
        <span>第 {{.page}} / {{.totalPages}} 页</span>
        
        {{if lt .page .totalPages}}
            <a href="?page={{add .page 1}}">下一页</a>
        {{end}}
    </div>
</div>
```

### 示例 2：表单

**控制器**:
```go
func (c *UserController) Edit() {
    id := c.QueryInt("id", 0)
    
    if c.Request.Method == "GET" {
        user := getUserByID(id)
        if user == nil {
            c.WriteHeader(404)
            c.View.Layout("main").Render("error/404")
            return
        }
        
        c.View.Set("title", "编辑用户")
        c.View.Set("user", user)
        c.View.Layout("main").Render("user/edit")
        return
    }
    
    // POST 处理
    username := c.PostForm("username")
    email := c.PostForm("email")
    
    err := updateUser(id, username, email)
    if err != nil {
        c.View.Set("error", err.Error())
        c.View.Set("title", "编辑用户")
        c.View.Set("user", map[string]interface{}{
            "ID":       id,
            "Username": username,
            "Email":    email,
        })
        c.View.Layout("main").Render("user/edit")
        return
    }
    
    c.Redirect("/user/list", 302)
}
```

**views/user/edit.html**:
```html
<div class="user-edit">
    <h2>{{.title}}</h2>
    
    {{if .error}}
        <div class="alert alert-error">{{.error}}</div>
    {{end}}
    
    <form method="POST">
        <div class="form-group">
            <label>用户名:</label>
            <input type="text" name="username" value="{{.user.Username}}" required>
        </div>
        
        <div class="form-group">
            <label>邮箱:</label>
            <input type="email" name="email" value="{{.user.Email}}" required>
        </div>
        
        <button type="submit">保存</button>
        <a href="/user/list">取消</a>
    </form>
</div>
```

### 示例 3：多语言

**控制器**:
```go
func (c *HomeController) Index() {
    c.View.Set("title", "home.title")
    c.View.Layout("main").Render("home/index")
}
```

**views/home/index.html**:
```html
<h2>{{tr .Lang "home.title" "首页"}}</h2>
<p>{{tr .Lang "home.welcome" "欢迎访问我们的网站"}}</p>

<div class="features">
    <h3>{{tr .Lang "home.features" "特性"}}</h3>
    <ul>
        <li>{{tr .Lang "home.feature1" "高性能"}}</li>
        <li>{{tr .Lang "home.feature2" "易用性"}}</li>
        <li>{{tr .Lang "home.feature3" "可扩展"}}</li>
    </ul>
</div>
```

## 最佳实践

### 1. 使用链式调用

```go
// ✅ 推荐：链式调用
c.View.Set("title", "标题")
    .Set("content", "内容")
    .Layout("main")
    .Render("page")

// ✅ 也可以
c.View.SetMap(gcore.M{
    "title":   "标题",
    "content": "内容",
}).Layout("main").Render("page")
```

### 2. 统一布局管理

```go
// 基础控制器设置默认布局
type BaseController struct {
    gmc.Controller
}

func (c *BaseController) Before() {
    // 设置通用变量
    c.View.Set("siteName", "我的网站")
    c.View.Set("year", time.Now().Year())
}

// 其他控制器继承
type UserController struct {
    BaseController
}
```

### 3. 错误处理

```go
func (c *HomeController) Index() {
    c.View.Render("home/index")
    
    // 检查错误
    if err := c.View.Err(); err != nil {
        c.Logger().Error(err)
        c.View.Layout("main").Render("error/500")
    }
}
```

### 4. 使用 RenderR 发送邮件

```go
func sendWelcomeEmail(user *User) error {
    // 使用模板渲染邮件内容
    view := gview.New(nil, tpl)
    view.Set("user", user)
    view.Set("activationLink", generateLink(user))
    
    html := view.RenderR("email/welcome")
    
    return sendEmail(user.Email, "欢迎注册", html)
}
```

### 5. 避免在模板中使用复杂逻辑

```go
// ❌ 不推荐：在模板中计算
<p>总价: {{mul .price .quantity}}</p>

// ✅ 推荐：在控制器中计算
c.View.Set("price", price)
c.View.Set("quantity", quantity)
c.View.Set("total", price * quantity)

// 模板中:
<p>总价: {{.total}}</p>
```

## 注意事项

1. **模板路径**
   - 相对于配置的模板目录
   - 不需要包含文件扩展名
   - 使用 `/` 作为路径分隔符

2. **布局变量**
   - 布局模板使用 `{{.GMC_LAYOUT_CONTENT}}` 作为内容占位符
   - 布局文件存放在 `layout` 子目录（可配置）

3. **自动变量注入**
   - `.G`, `.P`, `.S`, `.C`, `.U`, `.H` 变量在渲染时自动注入
   - `.S` 需要先调用 `SessionStart()`

4. **HTML 转义**
   - 默认会进行 HTML 转义
   - 使用 `tohtml` 函数输出原始 HTML

5. **性能考虑**
   - 模板在首次使用时解析和编译
   - 生产环境建议启用模板缓存

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [GMC Template](../template/README.md)
- [GMC Controller](../controller/README.md)
- [GMC HTTP Server](../server/README.md)
- [Sprig 函数文档](../template/sprig/docs)

