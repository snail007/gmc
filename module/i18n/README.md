# GMC I18n 模块

## 简介

GMC I18n（国际化）模块提供多语言支持，可轻松实现应用的国际化和本地化。

## 功能特性

- **多语言支持**：支持任意数量的语言
- **语言解析**：自动解析 Accept-Language header
- **回退机制**：支持语言回退
- **HTML 安全**：支持 HTML 模板安全输出
- **键值翻译**：基于键值对的翻译机制

## 快速开始

### 初始化 I18n

```go
package main

import (
    "github.com/snail007/gmc/module/i18n"
)

func main() {
    i18n := gi18n.NewI18n()
    
    // 添加英语翻译
    i18n.Add("en", map[string]string{
        "hello":   "Hello",
        "welcome": "Welcome to our site",
        "goodbye": "Goodbye",
    })
    
    // 添加中文翻译
    i18n.Add("zh", map[string]string{
        "hello":   "你好",
        "welcome": "欢迎访问我们的网站",
        "goodbye": "再见",
    })
    
    // 设置默认语言
    i18n.Lang("en")
    
    // 翻译
    msg := i18n.Tr("zh", "hello")
    println(msg) // 输出: 你好
}
```

### 在 Web 应用中使用

```go
func Handler(ctx gcore.Ctx) {
    i18n := ctx.I18n()
    
    // 从请求中解析用户语言
    langs, _ := i18n.ParseAcceptLanguage(ctx.Request())
    
    // 翻译文本
    msg := i18n.TrLangs(langs, "welcome", "Welcome")
    
    ctx.Write(msg)
}
```

### 模板中使用

```go
func Handler(ctx gcore.Ctx) {
    data := gcore.M{
        "title": ctx.I18n().Tr("en", "welcome"),
    }
    
    ctx.View("index.html", data)
}
```

## API 参考

### 初始化

```go
func NewI18n() gcore.I18n
```

### 方法

```go
// 添加语言翻译
Add(lang string, data map[string]string)

// 设置默认语言
Lang(lang string)

// 翻译文本
Tr(lang, key string, defaultMessage ...string) string

// 多语言翻译（按优先级）
TrLangs(langs []string, key string, defaultMessage ...string) string

// HTML 模板安全翻译
TrV(lang, key string, defaultMessage ...string) template.HTML

// 解析 Accept-Language
ParseAcceptLanguage(r *http.Request) ([]string, error)
ParseAcceptLanguageStr(s string) (string, error)

// 克隆 i18n 实例
Clone(lang string) gcore.I18n
```

## 配置示例

### 从文件加载翻译

```go
package main

import (
    "encoding/json"
    "io/ioutil"
    "github.com/snail007/gmc/module/i18n"
)

func LoadI18n() gcore.I18n {
    i18n := gi18n.NewI18n()
    
    // 加载英语翻译
    enData, _ := ioutil.ReadFile("lang/en.json")
    var enTrans map[string]string
    json.Unmarshal(enData, &enTrans)
    i18n.Add("en", enTrans)
    
    // 加载中文翻译
    zhData, _ := ioutil.ReadFile("lang/zh.json")
    var zhTrans map[string]string
    json.Unmarshal(zhData, &zhTrans)
    i18n.Add("zh", zhTrans)
    
    i18n.Lang("en") // 默认英语
    
    return i18n
}
```

### 翻译文件示例

**lang/en.json:**
```json
{
  "hello": "Hello",
  "welcome": "Welcome to our site",
  "user.login": "Login",
  "user.logout": "Logout",
  "error.not_found": "Page not found"
}
```

**lang/zh.json:**
```json
{
  "hello": "你好",
  "welcome": "欢迎访问我们的网站",
  "user.login": "登录",
  "user.logout": "登出",
  "error.not_found": "页面未找到"
}
```

## 使用场景

1. **多语言网站**：支持多种语言的网站
2. **国际化应用**：面向全球用户的应用
3. **本地化内容**：根据用户语言显示不同内容
4. **多语言 API**：提供多语言的 API 响应

## 最佳实践

### 1. 自动检测用户语言

```go
func GetUserLang(ctx gcore.Ctx) string {
    i18n := ctx.I18n()
    langs, err := i18n.ParseAcceptLanguage(ctx.Request())
    if err == nil && len(langs) > 0 {
        return langs[0]
    }
    return "en" // 默认语言
}
```

### 2. 使用语言回退

```go
// 尝试多个语言，直到找到翻译
msg := i18n.TrLangs([]string{"zh-CN", "zh", "en"}, "key", "Default")
```

### 3. 组织翻译键

```go
// 使用命名空间组织翻译键
translations := map[string]string{
    "user.login":    "Login",
    "user.register": "Register",
    "error.400":     "Bad Request",
    "error.404":     "Not Found",
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
