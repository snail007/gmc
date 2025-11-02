# GMC Template 模块

## 简介

GMC Template 模块是对 Go 标准库 `text/template` 的封装和扩展，提供了更强大的模板功能。支持模板继承、自定义函数、嵌入式文件、Sprig 函数库等特性。

## 功能特性

- **基于 text/template**：兼容 Go 标准模板语法
- **扩展函数库**：集成 Sprig 函数库子集
- **嵌入式文件支持**：支持从内存加载模板
- **自定义函数**：支持添加自定义模板函数
- **灵活配置**：支持自定义定界符、扩展名等
- **错误处理**：友好的错误提示

## 安装

```bash
go get github.com/snail007/gmc
```

## 快速开始

### 在 GMC 应用中使用

模板功能已集成在 GMC 框架中，通过 View 模块使用：

```go
package controller

import (
    "github.com/snail007/gmc"
)

type HomeController struct {
    gmc.Controller
}

func (c *HomeController) Index() {
    // 模板自动加载和渲染
    c.View.Set("title", "首页")
    c.View.Render("home/index")
}
```

### 直接使用 Template

如果需要独立使用模板引擎：

```go
package main

import (
    "os"
    "github.com/snail007/gmc"
    gtemplate "github.com/snail007/gmc/http/template"
)

func main() {
    // 创建模板实例，需要提供一个 Ctx 对象
    ctx := gmc.New.Ctx()
    tpl, err := gtemplate.NewTemplate(ctx, "./views")
    if err != nil {
        panic(err)
    }
    
    // 渲染模板
    data := map[string]interface{}{
        "title":   "Hello",
        "content": "World",
    }
    
    result, err := tpl.Execute("index", data)
    if err != nil {
        panic(err)
    }
    
    os.Stdout.Write(result)
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

# 左定界符
delimiterleft = "{{"

# 右定界符
delimiterright = "}}"

# 布局文件目录（相对于 dir）
layout = "layout"
```

## 模板语法

GMC Template 基于 Go 的 `text/template`，支持所有标准语法。

### 基本语法

#### 变量输出

```html
{{.title}}
{{.user.Name}}
{{.items.0}}
```

#### 条件判断

```html
{{if .condition}}
    条件为真
{{else if .other}}
    其他条件
{{else}}
    条件为假
{{end}}
```

#### 循环

```html
<!-- range 遍历 -->
{{range .items}}
    <li>{{.}}</li>
{{end}}

<!-- 带索引的 range -->
{{range $index, $item := .items}}
    <li>{{$index}}: {{$item}}</li>
{{end}}

<!-- 带键值的 range -->
{{range $key, $value := .map}}
    <p>{{$key}}: {{$value}}</p>
{{end}}
```

#### 变量定义

```html
{{$name := .user.Name}}
{{$age := .user.Age}}
<p>{{$name}} is {{$age}} years old</p>
```

#### 管道操作

```html
{{.content | tohtml}}
{{.number | printf "%d"}}
{{.text | upper | trim}}
```

### 标准函数

Go 模板引擎内置函数：

#### 输出函数

```html
{{print .value}}
{{printf "%s: %d" .name .count}}
{{println .text}}
```

#### HTML 处理

```html
<!-- HTML 转义 -->
{{html .content}}

<!-- JavaScript 转义 -->
{{js .script}}

<!-- URL 查询转义 -->
{{urlquery .param}}
```

#### 比较函数

```html
{{if eq .status "active"}}激活{{end}}
{{if ne .count 0}}有数据{{end}}
{{if lt .age 18}}未成年{{end}}
{{if le .score 60}}不及格{{end}}
{{if gt .price 100}}贵{{end}}
{{if ge .stock 10}}充足{{end}}
```

#### 逻辑函数

```html
{{if and .isLogin .isAdmin}}管理员{{end}}
{{if or .isMember .isVIP}}会员{{end}}
{{if not .isDeleted}}正常{{end}}
```

#### 其他函数

```html
<!-- 长度 -->
{{len .items}}

<!-- 索引访问 -->
{{index .array 0}}
{{index .map "key"}}

<!-- 切片 -->
{{slice .text 0 10}}

<!-- 函数调用 -->
{{call .func .arg1 .arg2}}
```

## GMC 扩展函数

### tr - 国际化翻译

翻译指定的键，返回 `template.HTML` 类型。

```html
{{tr .Lang "welcome.message" "欢迎"}}
{{tr .Lang "menu.home" "首页"}}
```

**参数：**
1. `.Lang` - 当前语言标识
2. 翻译键
3. 默认文本（可选）

### trs - 翻译（字符串）

类似 `tr`，但返回 `string` 类型。

```html
<input placeholder="{{trs .Lang "form.username" "用户名"}}">
```

### string - 转换为字符串

将任意类型转换为字符串。

```html
{{string .count}}
{{string .object}}
```

### tohtml - 转换为 HTML

将字符串转换为 `template.HTML` 类型，不进行 HTML 转义。

```html
<!-- 输出原始 HTML -->
{{tohtml .htmlContent}}
```

**⚠️ 警告：** 仅对可信内容使用此函数，以防 XSS 攻击。

### val - 安全获取变量

获取模板变量，如果不存在返回空字符串而非 `<no value>`。

```html
{{val .maybeUndefined}}
{{val .optional.field}}
```

## Sprig 函数库

GMC 集成了 [Sprig](https://masterminds.github.io/sprig/) 函数库的子集，提供丰富的模板函数。

### 字符串函数

```html
<!-- 大小写转换 -->
{{upper "hello"}}  <!-- HELLO -->
{{lower "WORLD"}}  <!-- world -->
{{title "hello world"}}  <!-- Hello World -->

<!-- 字符串操作 -->
{{trim " text "}}  <!-- text -->
{{trimPrefix "hello" "he"}}  <!-- llo -->
{{trimSuffix "world" "ld"}}  <!-- wor -->
{{repeat 3 "abc"}}  <!-- abcabcabc -->

<!-- 替换 -->
{{replace "hello" "l" "L"}}  <!-- heLLo -->
```

### 数值函数

```html
<!-- 数学运算 -->
{{add 1 2}}  <!-- 3 -->
{{sub 5 3}}  <!-- 2 -->
{{mul 4 5}}  <!-- 20 -->
{{div 10 2}}  <!-- 5 -->
{{mod 10 3}}  <!-- 1 -->

<!-- 最大最小值 -->
{{max 1 5 3}}  <!-- 5 -->
{{min 1 5 3}}  <!-- 1 -->
```

### 日期函数

```html
<!-- 格式化日期 -->
{{now | date "2006-01-02"}}
{{now | date "2006-01-02 15:04:05"}}

<!-- 日期计算 -->
{{now | dateModify "+24h" | date "2006-01-02"}}
```

### 列表函数

```html
<!-- 列表操作 -->
{{$list := slice 1 2 3}}  <!-- 使用 slice 创建列表 -->
{{append $list 4}}  <!-- 追加元素 -->
{{prepend $list 0}}  <!-- 前置元素 -->
{{first $list}}  <!-- 1 -->
{{rest $list}}  <!-- [2 3] -->
{{last $list}}  <!-- 3 -->
```

**注意**：GMC 模板没有提供 `dict` 函数。如需创建字典，请在控制器中准备数据：

```go
// 在控制器中
c.View.Set("myDict", map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
})
```

```html
<!-- 在模板中使用 -->
{{.myDict.key1}}
```

### 默认值函数

```html
<!-- 提供默认值 -->
{{default "默认值" .maybeEmpty}}
{{coalesce .val1 .val2 "fallback"}}
```

### 完整 Sprig 文档

更多函数请参考：[template/sprig/docs](./sprig/docs)

## 添加自定义函数

### 在 HTTP Server 中添加

```go
func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    // 添加自定义函数
    s.AddFuncMap(map[string]interface{}{
        // 格式化价格
        "formatPrice": func(price float64) string {
            return fmt.Sprintf("¥%.2f", price)
        },
        
        // 计算折扣价
        "discount": func(price float64, percent int) float64 {
            return price * float64(100-percent) / 100
        },
        
        // 日期格式化
        "formatDate": func(t time.Time) string {
            return t.Format("2006-01-02")
        },
    })
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

### 在模板中使用

```html
<p>原价: {{formatPrice .originalPrice}}</p>
<p>折扣价: {{formatPrice (discount .originalPrice 20)}}</p>
<p>日期: {{formatDate .createdAt}}</p>
```

## 嵌入式模板

### 使用 go:embed

```go
package main

import (
    "embed"
    "github.com/snail007/gmc"
    gtemplate "github.com/snail007/gmc/http/template"
)

//go:embed views/*
var viewsFS embed.FS

func main() {
    app := gmc.New.App()
    s := gmc.New.HTTPServer(app.Ctx())
    
    // 从嵌入式文件系统读取模板
    files := make(map[string][]byte)
    entries, _ := viewsFS.ReadDir("views")
    for _, entry := range entries {
        if !entry.IsDir() {
            content, _ := viewsFS.ReadFile("views/" + entry.Name())
            files[entry.Name()] = content
        }
    }
    
    // 设置嵌入式模板
    gtemplate.SetBinBytes(files)
    
    app.AddService(gmc.ServiceItem{Service: s})
    app.Run()
}
```

### 使用 Base64 编码

```go
// 工具函数：将模板文件转换为 Base64
func encodeTemplate(filepath string) string {
    data, _ := ioutil.ReadFile(filepath)
    return base64.StdEncoding.EncodeToString(data)
}

// 设置 Base64 编码的模板
gtemplate.SetBinBase64(map[string]string{
    "index.html": "PGh0bWw+...",  // Base64 编码的内容
    "about.html": "PGh0bWw+...",
})
```

## 高级用法

### 自定义定界符

如果模板中使用了 Vue.js、Angular 等前端框架，可能需要修改定界符：

```toml
[template]
delimiterleft = "<%"
delimiterright = "%>"
```

模板中使用：
```html
<h1><%".title"%></h1>
<div id="app">
    <!-- Vue.js 可以正常使用 {{}} -->
    {{message}}
</div>
```

### 模板注释

```html
{{/* 这是模板注释，不会出现在输出中 */}}
{{- /* 去除前面的空白 */ -}}
```

### 去除空白

```html
{{- .value -}}  <!-- 去除前后空白 -->
{{- .value }}   <!-- 仅去除前面空白 -->
{{ .value -}}   <!-- 仅去除后面空白 -->
```

### 模板继承

通过 Layout 实现模板继承，详见 [View 模块文档](../view/README.md)。

## 性能优化

### 模板缓存

模板在首次使用时会被解析和编译，之后会被缓存。

### 预编译模板

```go
// 在应用启动时预编译所有模板
ctx := gmc.New.Ctx()
tpl, err := gtemplate.NewTemplate(ctx, "./views")
if err != nil {
    panic(err)
}

// 解析所有模板文件
tpl.Parse()
```

## 最佳实践

### 1. 模板组织

```
views/
├── layout/          # 布局模板
│   ├── main.html
│   └── admin.html
├── components/      # 组件模板
│   ├── header.html
│   └── footer.html
├── home/           # 页面模板
│   └── index.html
└── user/
    ├── list.html
    └── detail.html
```

### 2. 命名规范

- 使用小写字母和下划线：`user_list.html`
- 或使用小写字母和连字符：`user-list.html`
- 布局文件前缀：`layout_*.html`
- 组件文件前缀：`component_*.html`

### 3. 避免复杂逻辑

```html
<!-- ❌ 不推荐：复杂逻辑 -->
{{if and (gt .age 18) (eq .status "active") (ne .role "guest")}}
    ...
{{end}}

<!-- ✅ 推荐：在控制器中处理 -->
{{if .canAccess}}
    ...
{{end}}
```

### 4. 使用自定义函数

将常用的格式化逻辑封装为自定义函数：

```go
s.AddFuncMap(map[string]interface{}{
    "formatMoney": formatMoney,
    "formatDate":  formatDate,
    "avatarURL":   getAvatarURL,
    "truncate":    truncateString,
})
```

### 5. XSS 防护

```html
<!-- ✅ 自动转义（安全） -->
{{.userInput}}

<!-- ⚠️ 不转义（仅用于可信内容） -->
{{tohtml .trustedHTML}}

<!-- ✅ 在控制器中清理用户输入 -->
content := sanitizeHTML(c.PostForm("content"))
c.View.Set("content", content)
```

## 错误处理

### 模板错误

模板语法错误会在解析时抛出：

```go
ctx := gmc.New.Ctx()
tpl, err := gtemplate.NewTemplate(ctx, "./views")
if err != nil {
    log.Fatal("模板解析错误:", err)
}
```

### 渲染错误

渲染时的错误可以通过 View.Err() 获取：

```go
c.View.Render("home/index")
if err := c.View.Err(); err != nil {
    c.Logger().Error("渲染错误:", err)
}
```

## 注意事项

1. **定界符**
   - 默认使用 `{{` 和 `}}`
   - 与前端框架冲突时需要修改

2. **HTML 转义**
   - 默认自动转义，防止 XSS
   - 使用 `tohtml` 函数输出原始 HTML

3. **模板路径**
   - 相对于配置的 `template.dir`
   - 不需要包含文件扩展名

4. **性能**
   - 模板会被缓存
   - 开发时可能需要重启查看更改

5. **嵌入式文件**
   - 用于生产部署，减少文件依赖
   - 文件路径不含前缀斜杠

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [GMC View 模块](../view/README.md)
- [GMC Controller](../controller/README.md)
- [Sprig 函数文档](./sprig/docs)
- [Go text/template 文档](https://pkg.go.dev/text/template)
