# ghash 包

## 简介

ghash 包提供了多种哈希算法的便捷封装，支持字符串、字节、文件和 Reader 的哈希计算。

## 功能特性

- **MD5**：MD5 哈希算法
- **SHA1**：SHA1 哈希算法  
- **SHA256**：SHA256 哈希算法
- **CRC32**：CRC32 校验
- **Blake2s**：Blake2s 哈希算法
- **多种输入**：支持字符串、字节数组、文件、Reader

## 安装

```bash
go get github.com/snail007/gmc/util/hash
```

## 快速开始

### MD5 哈希

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/hash"
)

func main() {
    // 字符串 MD5
    hash1 := ghash.MD5("hello world")
    fmt.Println(hash1)
    
    // 字节数组 MD5
    hash2 := ghash.Md5Bytes([]byte("hello world"))
    fmt.Println(hash2)
    
    // 文件 MD5
    hash3, err := ghash.MD5File("/path/to/file")
    if err != nil {
        panic(err)
    }
    fmt.Println(hash3)
}
```

### SHA256 哈希

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/hash"
)

func main() {
    // 字符串 SHA256
    hash1 := ghash.SHA256("hello world")
    fmt.Println(hash1)
    
    // 字节数组 SHA256
    hash2 := ghash.SHA256Bytes([]byte("hello world"))
    fmt.Println(hash2)
    
    // 文件 SHA256
    hash3, err := ghash.SHA256File("/path/to/file")
    if err != nil {
        panic(err)
    }
    fmt.Println(hash3)
}
```

### 其他哈希算法

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/hash"
)

func main() {
    str := "hello world"
    
    // SHA1
    fmt.Println("SHA1:", ghash.SHA1(str))
    
    // CRC32
    fmt.Println("CRC32:", ghash.CRC32(str))
    
    // Blake2s
    fmt.Println("Blake2s:", ghash.Blake2s(str))
}
```

### 使用 Reader

```go
package main

import (
    "fmt"
    "os"
    "github.com/snail007/gmc/util/hash"
)

func main() {
    file, err := os.Open("/path/to/file")
    if err != nil {
        panic(err)
    }
    
    hash, err := ghash.MD5Reader(file)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("MD5:", hash)
}
```

## API 参考

### MD5

```go
func MD5(str string) string
func Md5Bytes(d []byte) string
func MD5File(filename string) (string, error)
func MD5Reader(reader io.ReadCloser) (string, error)
```

### SHA1

```go
func SHA1(str string) string
func SHA1Bytes(b []byte) string
func SHA1File(filename string) (string, error)
func SHA1Reader(reader io.ReadCloser) (string, error)
```

### SHA256

```go
func SHA256(str string) string
func SHA256Bytes(b []byte) string
func SHA256File(filename string) (string, error)
func SHA256Reader(reader io.ReadCloser) (string, error)
```

### CRC32

```go
func CRC32(str string) string
func CRC32Bytes(b []byte) string
func CRC32File(filename string) (string, error)
func CRC32Reader(reader io.ReadCloser) (string, error)
```

### Blake2s

```go
func Blake2s(str string) string
func Blake2sBytes(b []byte) string
func Blake2sFile(filename string) (string, error)
func Blake2sReader(reader io.ReadCloser) (string, error)
```

## 使用场景

1. **文件完整性校验**
2. **密码哈希**（注意：MD5 和 SHA1 不安全，仅用于非安全场景）
3. **数据去重**
4. **缓存键生成**
5. **数字签名**

## 注意事项

1. MD5 和 SHA1 已被认为不安全，不应用于安全敏感场景
2. 推荐使用 SHA256 或 Blake2s 用于安全场景
3. 文件哈希计算会读取整个文件
4. Reader 版本会关闭 Reader

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
