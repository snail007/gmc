# gargs 包

## 简介

gargs 包实现了一个简单的词法分析器，使用 Shell 风格的引号和注释规则将输入分割成 token。该包基于 Google 的开源实现，遵循 Apache License 2.0 协议。

## 功能特性

- **Shell 风格解析**：支持单引号、双引号、转义字符和注释
- **流式处理**：支持从 io.Reader 流式读取和解析
- **灵活的 API**：提供简单的字符串分割和高级的 token 流处理

## 安装

```bash
go get github.com/snail007/gmc/util/args
```

## 快速开始

### 基本字符串分割

最简单的使用方式是使用 `Split` 函数分割字符串：

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/args"
)

func main() {
    result, err := gargs.Split(`one "two three" four`)
    if err != nil {
        panic(err)
    }
    // 输出: []string{"one", "two three", "four"}
    fmt.Println(result)
}
```

### 流式处理字符串

使用 Lexer 处理字符串流：

```go
package main

import (
    "fmt"
    "io"
    "strings"
    "github.com/snail007/gmc/util/args"
)

func main() {
    input := strings.NewReader(`command "arg with spaces" --flag=value`)
    lexer := gargs.NewLexer(input)
    
    for {
        token, err := lexer.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }
        fmt.Println(token)
    }
}
```

### 访问原始 Token 流

使用 Tokenizer 获取包括注释在内的所有 token：

```go
package main

import (
    "fmt"
    "io"
    "strings"
    "github.com/snail007/gmc/util/args"
)

func main() {
    input := strings.NewReader(`command # this is a comment`)
    tokenizer := gargs.NewTokenizer(input)
    
    for {
        token, err := tokenizer.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }
        // token 包含类型和值信息
        fmt.Printf("Type: %v, Value: %s\n", token.tokenType, token.value)
    }
}
```

## API 参考

### 主要类型

#### Token

表示一个词法 token，包含类型和值：

```go
type Token struct {
    tokenType TokenType
    value     string
}
```

**方法：**
- `Equal(b *Token) bool`：比较两个 token 是否相等

#### TokenType

Token 的类型常量：

- `UnknownToken`：未知类型
- `WordToken`：单词 token
- `SpaceToken`：空格 token
- `CommentToken`：注释 token

### 主要函数和类型

#### Split

```go
func Split(s string) ([]string, error)
```

将字符串分割成字符串切片。这是最常用的函数。

**参数：**
- `s`：要分割的字符串

**返回值：**
- `[]string`：分割后的字符串数组
- `error`：错误信息

**示例：**

```go
result, _ := gargs.Split(`echo "hello world" 'single quote'`)
// result: ["echo", "hello world", "single quote"]
```

#### Lexer

```go
type Lexer struct { ... }

func NewLexer(r io.Reader) *Lexer
func (l *Lexer) Next() (string, error)
```

词法分析器，跳过空白和注释，只返回单词 token。

**示例：**

```go
lexer := gargs.NewLexer(strings.NewReader("one two three"))
word, err := lexer.Next() // "one"
```

#### Tokenizer

```go
type Tokenizer struct { ... }

func NewTokenizer(r io.Reader) *Tokenizer
func (t *Tokenizer) Next() (*Token, error)
```

Token 流处理器，返回所有类型的 token，包括注释。

## 解析规则

### 引号规则

1. **双引号 (`"`)**：支持转义字符
   ```
   "hello \"world\"" → hello "world"
   ```

2. **单引号 (`'`)**：不支持转义，所有字符都是字面量
   ```
   'hello "world"' → hello "world"
   ```

### 转义字符

在双引号内或未引用的文本中，反斜杠 `\` 用作转义字符：

```
echo \"hello\" → echo "hello"
path\ with\ spaces → path with spaces
```

### 注释

井号 `#` 开始的内容直到行尾被视为注释：

```
command arg1 # this is a comment
```

### 空白字符

以下字符被视为空白分隔符：
- 空格 (` `)
- 制表符 (`\t`)
- 回车 (`\r`)
- 换行 (`\n`)

## 使用场景

1. **命令行解析**：解析命令行参数字符串
2. **配置文件**：解析 Shell 风格的配置文件
3. **脚本解释器**：构建简单的脚本解释器
4. **参数处理**：处理包含空格和特殊字符的参数

## 注意事项

1. Token 流中的注释默认会被 Lexer 过滤，如需保留注释请使用 Tokenizer
2. 转义字符只在双引号内和未引用文本中有效
3. 单引号内的所有字符都是字面量，包括反斜杠
4. 未闭合的引号会返回错误
5. EOF 后的转义字符会返回错误

## 错误处理

常见错误情况：

```go
// 未闭合的引号
_, err := gargs.Split(`"unclosed quote`)
// err: EOF found when expecting closing quote

// EOF 后的转义字符
_, err := gargs.Split(`escape at end\`)
// err: EOF found after escape character
```

## 性能考虑

- Lexer 和 Tokenizer 使用 `bufio.Reader` 进行流式读取，适合处理大文件
- `Split` 函数适合处理短字符串，会一次性加载所有内容到内存
- 对于大量数据，建议使用流式 API（Lexer/Tokenizer）

## 许可证

本包基于 Google 的开源实现，遵循 Apache License 2.0 协议。

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [Apache License 2.0](http://www.apache.org/licenses/LICENSE-2.0)
