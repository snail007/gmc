# GMC I18n 模块

## 简介

GMC I18n（国际化）模块提供了一套完整的多语言解决方案，可以轻松地为你的 Web 应用或 API 服务实现国际化和本地化。

## 功能特性

-   **TOML 格式**：语言文件使用简单直观的 `key = "value"` TOML 格式。
-   **自动加载**：可根据 `app.toml` 配置自动加载语言目录。
-   **嵌入式支持**：支持使用 `go:embed` 将语言文件直接打包到二进制文件中，实现单文件部署。
-   **自动语言检测**：能自动从 HTTP 请求的 `Accept-Language` 头中解析出用户的首选语言列表。
-   **回退机制**：当指定语言的翻译不存在时，会自动尝试使用备用语言（Fallback）或默认语言。
-   **模板集成**：在模板中可以方便地调用翻译函数。

## 快速开始

### 1. 创建语言文件

在你的项目 `i18n` 目录下（目录名可在 `app.toml` 中配置），创建语言文件。文件名应遵循 `IETF BCP 47` 语言标签规范，如 `en-US.toml`, `zh-CN.toml`。

**`i18n/en-US.toml`:**

```toml
hello = "Hello GMC!"
welcome = "Welcome, %s"
```

**`i18n/zh-CN.toml`:**

```toml
hello = "你好 GMC！"
welcome = "欢迎, %s"
```

### 2. 配置 `app.toml`

在 `app.toml` 中启用 `i18n` 并指定语言文件目录。

```toml
[i18n]
# 启用 i18n
enable=true
# 语言文件所在目录
dir="i18n"
# 默认回退语言
default="en-US"
```

框架在启动时会自动加载 `dir` 目录下的所有 `.toml` 文件。

### 3. 在控制器中使用

```go
package controller

import (
    "github.com/snail007/gmc"
)

type DemoController struct {
    gmc.Controller
}

func (this *DemoController) Hello() {
    // Tr 方法会自动从请求头中解析语言，并按优先级查找翻译
    // 如果找不到，则使用在 app.toml 中配置的默认语言
    // "Guest" 会替换掉翻译文本中的 %s
    this.Write(this.Tr("welcome", "Guest"))
}
```

### 4. 在模板中使用

在控制器中，通过 `this.View.Set` 传递翻译后的文本。

```go
func (this *DemoController) Index() {
    this.View.Set("title", this.Tr("hello"))
    this.View.Render("index.html")
}
```

**`views/index.html`:**

```html
<h1>{{.title}}</h1>
```

## 打包到二进制文件 (go:embed)

你可以将所有语言文件嵌入到最终的二进制程序中，实现真正的单文件部署。

### 方法一：使用 InitEmbedFS（推荐）

这是最简单的方式，直接使用框架提供的 `InitEmbedFS` 方法。

**`i18n/i18n.go`:**

```go
package i18n

import "embed"

//go:embed *.toml
var I18nFS embed.FS
```

**`main.go` - 方式A（在main中初始化）：**

```go
package main

import (
    "github.com/snail007/gmc"
    gi18n "github.com/snail007/gmc/module/i18n"
    
    // 导入你的 i18n 包
    "myapp/i18n"
)

func main() {
    // 在 app.Run() 之前初始化嵌入的 i18n 文件
    err := gi18n.InitEmbedFS(i18n.I18nFS, "zh-CN")
    if err != nil {
        panic(err)
    }
    
    app := gmc.New.AppDefault()
    app.Run()
}
```

**`main.go` - 方式B（在OnRun钩子中初始化，推荐）：**

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
    gi18n "github.com/snail007/gmc/module/i18n"
    
    // 导入你的 i18n 包
    "myapp/i18n"
)

func main() {
    app := gmc.New.AppDefault()
    
    // 在 OnRun 钩子中初始化（推荐方式）
    app.OnRun(func(cfg gcore.Config) error {
        // 可以从配置中读取默认语言
        defaultLang := cfg.GetString("i18n.default")
        if defaultLang == "" {
            defaultLang = "zh-CN"
        }
        return gi18n.InitEmbedFS(i18n.I18nFS, defaultLang)
    })
    
    app.Run()
}
```

**重要提示**：
1. 使用 `InitEmbedFS` 时，必须将 `app.toml` 中 `[i18n]` 的 `enable` 设置为 `false` 或 `dir` 设置为空 (`dir=""`)，以避免框架尝试从文件系统加载。
2. 推荐使用方式B（OnRun钩子），因为可以访问配置对象，更加灵活。
3. 方式A更简单直接，适合不需要配置的简单场景。

### 方法二：手动加载（高级用法）

如果你需要更多的控制，可以手动遍历和加载嵌入的文件。

**`i18n/i18n.go`:**

```go
package i18n

import "embed"

//go:embed *.toml
var I18nFS embed.FS
```

**`main.go`:**

```go
package main

import (
    "bytes"
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
    gconfig "github.com/snail007/gmc/module/config"
    gi18n "github.com/snail007/gmc/module/i18n"
    "io/fs"
    "path/filepath"
    "strings"

    // 显式导入你的 i18n 包
    "myapp/i18n"
)

func main() {
    app := gmc.New.AppDefault()

    // 在 AfterInit 钩子中加载嵌入的 i18n 文件
    app.AddService(gcore.ServiceItem{
        Service: gmc.New.HTTPServer(app.Ctx()).(gcore.Service),
        AfterInit: func(s *gcore.ServiceItem) (err error) {
            // 获取 i18n 服务实例
            i18nService := gcore.ProviderI18n()(s.Service.(gcore.HTTPServer).Ctx()).(gcore.I18n)

            // 遍历 embed.FS 并加载语言文件
            err = fs.WalkDir(i18n.I18nFS, ".", func(path string, d fs.DirEntry, err error) error {
                if err != nil {
                    return err
                }
                if d.IsDir() {
                    return nil
                }
                if strings.HasSuffix(path, ".toml") {
                    content, err := i18n.I18nFS.ReadFile(path)
                    if err != nil {
                        return err
                    }
                    // 使用 gconfig 解析 TOML 内容
                    c := gconfig.New()
                    c.SetConfigType("toml")
                    err = c.ReadConfig(bytes.NewReader(content))
                    if err != nil {
                        return err
                    }
                    // 添加到 i18n 服务
                    lang := strings.TrimSuffix(filepath.Base(path), ".toml")
                    data := map[string]string{}
                    for _, k := range c.AllKeys() {
                        data[k] = c.GetString(k)
                    }
                    i18nService.Add(lang, data)
                }
                return nil
            })
            return
        },
    })

    app.Run()
}
```

## API 参考

### `gcore.I18n` 接口

-   `Tr(lang, key string, defaultMessage ...string) string`: 翻译一个键。`lang` 参数通常可以传空字符串，此时会自动从 HTTP 请求头或默认配置中获取语言。
-   `TrLangs(langs []string, key string, defaultMessage ...string) string`: 按指定的语言列表优先级进行翻译。
-   `Add(lang string, data map[string]interface{})`: 手动添加一个语言的翻译数据。
-   `ParseAcceptLanguage(r *http.Request) ([]string, error)`: 从 HTTP 请求头解析出语言列表。

### `Controller` 中的方法

-   `this.Tr(key string, args ...interface{}) string`: 在控制器中使用的便捷方法。它会自动处理语言检测，并使用 `fmt.Sprintf` 格式化翻译结果。

    ```go
    // i18n/en-US.toml: welcome = "Welcome, %s"
    this.Tr("welcome", "John") // 输出: "Welcome, John"
    ```
