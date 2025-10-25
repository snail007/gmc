# ghttp 包

## 简介

ghttp 包提供了功能强大的 HTTP 客户端，支持 GET/POST、文件上传下载、批量请求、重试机制、代理等功能。

## 功能特性

- **基本请求**：GET、POST、PUT、DELETE 等
- **文件操作**：上传、下载文件
- **批量请求**：同时发送多个请求
- **重试机制**：自动重试失败的请求
- **Cookie 管理**：自动保存和发送 Cookie
- **代理支持**：HTTP/SOCKS5 代理
- **TLS 配置**：自定义证书、跳过验证等
- **连接池**：复用 HTTP 连接

## 安装

```bash
go get github.com/snail007/gmc/util/http
```

## 快速开始

### GET 请求

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/http"
)

func main() {
    body, code, resp, err := ghttp.Get(
        "https://api.example.com/users",
        5*time.Second,
        map[string]string{"page": "1"},
        map[string]string{"Authorization": "Bearer token"},
    )
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Status: %d\n", code)
    fmt.Printf("Body: %s\n", string(body))
}
```

### POST 请求

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/http"
)

func main() {
    body, code, resp, err := ghttp.Post(
        "https://api.example.com/users",
        map[string]string{
            "name": "John",
            "email": "john@example.com",
        },
        5*time.Second,
        map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
    )
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Status: %d\n", code)
    fmt.Printf("Response: %s\n", string(body))
}
```

### 文件上传

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/http"
)

func main() {
    body, resp, err := ghttp.Upload(
        "https://api.example.com/upload",
        "file",
        "/path/to/file.txt",
        map[string]string{"description": "test file"},
        30*time.Second,
        nil,
    )
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Response:", body)
}
```

### 文件下载

```go
package main

import (
    "time"
    "github.com/snail007/gmc/util/http"
)

func main() {
    resp, err := ghttp.DownloadToFile(
        "https://example.com/file.zip",
        60*time.Second,
        nil,
        nil,
        "/tmp/downloaded.zip",
    )
    
    if err != nil {
        panic(err)
    }
    
    println("Downloaded successfully")
}
```

### 重试请求

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/http"
)

func main() {
    // 创建可重试的 GET 请求，最多重试 3 次
    tr, err := ghttp.NewTriableGet(
        "https://api.example.com/data",
        3,
        5*time.Second,
        nil,
        nil,
    )
    
    if err != nil {
        panic(err)
    }
    
    body, code, resp, err := tr.Do()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Status: %d\n", code)
    fmt.Println("Body:", string(body))
}
```

### 批量请求

```go
package main

import (
    "fmt"
    "time"
    "github.com/snail007/gmc/util/http"
)

func main() {
    urls := []string{
        "https://api.example.com/users/1",
        "https://api.example.com/users/2",
        "https://api.example.com/users/3",
    }
    
    br, err := ghttp.NewBatchGet(urls, 5*time.Second, nil, nil)
    if err != nil {
        panic(err)
    }
    
    results := br.Do()
    
    for i, result := range results {
        if result.Err != nil {
            fmt.Printf("Request %d failed: %v\n", i, result.Err)
        } else {
            fmt.Printf("Request %d: %d - %s\n", 
                i, result.Code, string(result.Body))
        }
    }
}
```

## API 参考

### 全局函数

- `Get(url, timeout, queryData, header)`
- `Post(url, data, timeout, header)`
- `PostOfReader(url, reader, timeout, header)`
- `Upload(url, fieldName, filename, data, timeout, header)`
- `Download(url, timeout, queryData, header)`
- `DownloadToFile(url, timeout, queryData, header, file)`
- `NewTriableGet(url, maxTry, timeout, queryData, header)`
- `NewTriablePost(url, maxTry, timeout, data, header)`
- `NewBatchGet(urls, timeout, data, header)`
- `NewBatchPost(urls, timeout, data, header)`

### HTTPClient 类型

创建自定义客户端：

```go
client := ghttp.NewHTTPClient()
client.SetProxy("http://proxy.example.com:8080")
client.SetTimeout(10 * time.Second)
```

## 使用场景

1. **API 调用**：调用 RESTful API
2. **文件传输**：上传下载文件
3. **网络爬虫**：批量抓取网页
4. **服务监控**：定期检查服务可用性
5. **数据同步**：与远程服务同步数据

## 注意事项

1. 默认全局客户端保持 Cookie 和连接池
2. 超时时间包括连接和读取时间
3. 批量请求并发执行
4. 重试请求间隔递增

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
