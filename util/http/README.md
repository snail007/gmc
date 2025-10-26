# ghttp 包

GMC 框架的 HTTP 客户端工具库，提供简洁易用且功能强大的 HTTP 请求能力。

## ✨ 功能特性

- 🚀 **简单易用**：链式调用、默认全局客户端
- 🔄 **智能重试**：支持失败自动重试，可配置重试次数
- 📦 **批量请求**：并发发送多个请求，支持等待全部完成或首个成功
- 📤 **文件传输**：支持文件上传、下载、流式上传
- 🍪 **Cookie 管理**：自动保存和发送 Cookie
- 🔐 **TLS 配置**：支持自定义证书、客户端证书认证、跳过验证
- 🌐 **代理支持**：HTTP/HTTPS/SOCKS5 代理，支持环境变量
- 🔌 **连接管理**：Keep-Alive、连接池、自定义 DNS
- 🎣 **钩子函数**：请求前后拦截器
- ⚡ **高性能**：连接复用、并发控制

## 📦 安装

```bash
go get github.com/snail007/gmc/util/http
```

## 🎯 核心概念

### 全局客户端 vs 独立客户端

```go
// 方式1：使用全局客户端（推荐，自动保持 Cookie 和连接池）
body, code, resp, err := ghttp.Get("https://example.com", 5*time.Second, nil, nil)

// 方式2：创建独立客户端（需要自定义配置时使用）
client := ghttp.NewHTTPClient()
client.SetProxy("http://proxy.com:8080")
body, code, resp, err := client.Get("https://example.com", 5*time.Second, nil, nil)
```

## 🚀 快速开始

### 基础 GET 请求

```go
import (
    "fmt"
    "time"
    ghttp "github.com/snail007/gmc/util/http"
)

// 最简单的 GET 请求
body, code, resp, err := ghttp.Get(
    "https://api.example.com/users",  // URL
    5*time.Second,                     // 超时时间
    nil,                               // 查询参数（可选）
    nil,                               // 请求头（可选）
)

// 带查询参数和请求头
body, code, resp, err := ghttp.Get(
    "https://api.example.com/users",
    5*time.Second,
    map[string]string{                 // 查询参数
        "page": "1",
        "size": "20",
    },
    map[string]string{                 // 请求头
        "Authorization": "Bearer your_token",
        "User-Agent": "MyApp/1.0",
    },
)

if err != nil {
    fmt.Printf("请求失败: %v\n", err)
    return
}
fmt.Printf("状态码: %d\n", code)
fmt.Printf("响应: %s\n", string(body))
```

### 基础 POST 请求

```go
// POST 表单数据
body, code, resp, err := ghttp.Post(
    "https://api.example.com/users",
    map[string]string{                 // 表单数据
        "name": "张三",
        "email": "zhangsan@example.com",
        "age": "25",
    },
    5*time.Second,                     // 超时时间
    map[string]string{                 // 请求头（可选）
        "Authorization": "Bearer token",
    },
)

// POST JSON 数据
jsonData := `{"name":"张三","email":"zhangsan@example.com"}`
body, code, resp, err := ghttp.PostOfReader(
    "https://api.example.com/users",
    strings.NewReader(jsonData),
    5*time.Second,
    map[string]string{
        "Content-Type": "application/json",
    },
)
```

## 📤 文件上传

### 上传本地文件

```go
// 参数说明：
// - url: 上传接口地址
// - fieldName: 表单字段名（文件字段）
// - filename: 本地文件路径
// - data: 附加的表单数据
// - timeout: 超时时间
// - header: 自定义请求头
body, resp, err := ghttp.Upload(
    "https://api.example.com/upload",
    "file",                            // 表单字段名
    "/path/to/document.pdf",           // 本地文件路径
    map[string]string{                 // 附加表单数据
        "description": "重要文档",
        "category": "文档",
    },
    30*time.Second,
    map[string]string{
        "Authorization": "Bearer token",
    },
)

if err != nil {
    fmt.Printf("上传失败: %v\n", err)
    return
}
fmt.Printf("服务器响应: %s\n", body)
```

### 流式上传

```go
// 从任意 io.Reader 上传
file, _ := os.Open("/path/to/large-file.zip")
defer file.Close()

body, resp, err := ghttp.UploadOfReader(
    "https://api.example.com/upload",
    "file",                            // 表单字段名
    "large-file.zip",                  // 文件名（显示用）
    file,                              // io.ReadCloser
    map[string]string{
        "description": "大文件",
    },
    60*time.Second,
    nil,
)
```

## 📥 文件下载

### 下载到内存

```go
// 适用于小文件
data, resp, err := ghttp.Download(
    "https://example.com/small-image.png",
    10*time.Second,
    nil,  // 查询参数
    nil,  // 请求头
)

if err != nil {
    fmt.Printf("下载失败: %v\n", err)
    return
}

// 保存到文件
err = ioutil.WriteFile("/tmp/image.png", data, 0644)
```

### 下载到文件

```go
// 适用于大文件，流式写入磁盘
resp, err := ghttp.DownloadToFile(
    "https://example.com/large-file.zip",
    60*time.Second,
    nil,                               // 查询参数
    nil,                               // 请求头
    "/tmp/downloaded.zip",             // 保存路径
)

if err != nil {
    fmt.Printf("下载失败: %v\n", err)
    return
}
fmt.Println("下载成功")
```

### 下载到自定义 Writer

```go
// 下载到任意 io.Writer（如缓冲区、管道等）
var buf bytes.Buffer
resp, err := ghttp.DownloadToWriter(
    "https://example.com/data.json",
    10*time.Second,
    nil,
    nil,
    &buf,  // 任意 io.Writer
)

fmt.Printf("下载内容: %s\n", buf.String())
```

## 🔄 重试请求

### 基础重试

```go
// 创建可重试的 GET 请求
// 参数：URL, 最大重试次数, 超时, 查询参数, 请求头
tr, err := ghttp.NewTriableGet(
    "https://api.example.com/unstable",
    3,                                 // 最多重试 3 次
    5*time.Second,
    map[string]string{"id": "123"},
    map[string]string{"Authorization": "Bearer token"},
)

if err != nil {
    panic(err)
}

// 执行请求（失败会自动重试）
resp := tr.Execute()

if resp.Err() != nil {
    fmt.Printf("所有重试均失败: %v\n", resp.Err())
    fmt.Printf("错误列表: %v\n", tr.ErrAll())
    return
}

body := resp.Body()
fmt.Printf("请求成功: %s\n", string(body))
fmt.Printf("状态码: %d\n", resp.StatusCode)
fmt.Printf("耗时: %v\n", resp.UsedTime())
```

### 可重试的 POST 请求

```go
tr, err := ghttp.NewTriablePost(
    "https://api.example.com/submit",
    3,  // 最多重试 3 次
    5*time.Second,
    map[string]string{
        "data": "important",
    },
    nil,
)

resp := tr.Execute()
if resp.Err() != nil {
    fmt.Printf("失败: %v\n", resp.Err())
    return
}
```

### 自定义错误检查

```go
tr, _ := ghttp.NewTriableGet("https://api.example.com/data", 3, 5*time.Second, nil, nil)

// 自定义失败判断逻辑（返回 error 表示请求失败需要重试）
tr.CheckErrorFunc(func(idx int, req *http.Request, resp *http.Response) error {
    // 例如：状态码 500-599 视为失败
    if resp.StatusCode >= 500 {
        return fmt.Errorf("服务器错误: %d", resp.StatusCode)
    }
    return nil
})

resp := tr.Execute()
```

### 请求拦截器

```go
tr, _ := ghttp.NewTriableGet("https://api.example.com/data", 3, 5*time.Second, nil, nil)

// 请求前拦截器
tr.SetBeforeDo(func(idx int, req *http.Request) {
    fmt.Printf("第 %d 次尝试，URL: %s\n", idx+1, req.URL)
})

// 请求后拦截器
tr.AfterDo(func(resp *ghttp.Response) {
    if resp.Err() != nil {
        fmt.Printf("请求失败: %v，耗时: %v\n", resp.Err(), resp.UsedTime())
    } else {
        fmt.Printf("请求成功，状态码: %d\n", resp.StatusCode)
    }
})

resp := tr.Execute()
```

## 📦 批量请求

### 批量 GET 请求

```go
// 同时请求多个 URL
urls := []string{
    "https://api.example.com/users/1",
    "https://api.example.com/users/2",
    "https://api.example.com/users/3",
    "https://api.example.com/users/4",
    "https://api.example.com/users/5",
}

// 参数：URL 列表, 超时, 查询参数, 请求头
br, err := ghttp.NewBatchGet(
    urls,
    5*time.Second,
    map[string]string{"format": "json"},  // 所有请求共用的查询参数
    map[string]string{"Authorization": "Bearer token"},
)

if err != nil {
    panic(err)
}

// 执行批量请求（等待所有请求完成）
br.Execute()

// 检查是否全部成功
if br.Success() {
    fmt.Println("所有请求均成功")
} else {
    fmt.Printf("失败数量: %d\n", br.ErrorCount())
}

// 获取所有响应
responses := br.RespAll()
for i, resp := range responses {
    if resp.Err() != nil {
        fmt.Printf("请求 %d 失败: %v\n", i, resp.Err())
    } else {
        fmt.Printf("请求 %d 成功: %d - %s\n", 
            i, resp.StatusCode, string(resp.Body()))
    }
}
```

### 批量 POST 请求

```go
urls := []string{
    "https://api.example.com/logs",
    "https://api.example.com/analytics",
    "https://api.example.com/metrics",
}

br, _ := ghttp.NewBatchPost(
    urls,
    5*time.Second,
    map[string]string{               // 所有 POST 请求共用的数据
        "timestamp": time.Now().String(),
        "app": "myapp",
    },
    nil,
)

br.Execute()
```

### 等待首个成功响应

```go
// 多个镜像源，只需要一个成功即可
urls := []string{
    "https://mirror1.example.com/data",
    "https://mirror2.example.com/data",
    "https://mirror3.example.com/data",
}

br, _ := ghttp.NewBatchGet(urls, 5*time.Second, nil, nil)

// 设置：获得首个成功响应后立即返回
br.WaitFirstSuccess().Execute()

// 获取首个成功的响应
resp := br.Resp()
if resp != nil {
    fmt.Printf("首个成功响应: %s\n", string(resp.Body()))
} else {
    fmt.Println("所有请求均失败")
}
```

### 批量请求配置

```go
br, _ := ghttp.NewBatchGet(urls, 5*time.Second, nil, nil)

// 启用重试（每个请求最多重试 2 次）
br.MaxTry(2)

// 自定义错误检查
br.CheckErrorFunc(func(idx int, req *http.Request, resp *http.Response) error {
    if resp.StatusCode != 200 {
        return fmt.Errorf("状态码异常: %d", resp.StatusCode)
    }
    return nil
})

// 请求前拦截器
br.SetBeforeDo(func(idx int, req *http.Request) {
    fmt.Printf("正在请求: %s\n", req.URL)
})

// 请求后拦截器
br.AppendAfterDo(func(resp *ghttp.Response) {
    fmt.Printf("完成请求，耗时: %v\n", resp.UsedTime())
})

br.Execute()
```

### 使用协程池控制并发

```go
import "github.com/snail007/gmc/util/gpool"

// 方式1：使用标准 BasicPool
pool := gpool.New(5)  // 限制并发数为 5，返回 *BasicPool
defer pool.Stop()

br, _ := ghttp.NewBatchGet(urls, 5*time.Second, nil, nil)
br.Pool(pool).Execute()

// 方式2：使用优化版 OptimizedPool（性能更好）
optimizedPool := gpool.NewOptimized(5)  // 返回 *OptimizedPool
defer optimizedPool.Stop()

br2, _ := ghttp.NewBatchGet(urls, 5*time.Second, nil, nil)
br2.Pool(optimizedPool).Execute()
```

**说明**：
- `Pool()` 方法接受 `gpool.Pool` 接口类型
- 可以传入 `gpool.BasicPool` 或 `gpool.OptimizedPool` 的任意实现
- `OptimizedPool` 性能更好，适合高并发场景
- 不设置 Pool 时，默认使用独立 goroutine（无并发限制）

## 🔧 高级配置

### 创建自定义客户端

```go
client := ghttp.NewHTTPClient()

// 设置代理
client.SetProxy("http://proxy.example.com:8080")
// 或使用 SOCKS5
client.SetProxy("socks5://127.0.0.1:1080")

// 设置自定义 DNS
client.SetDNS("8.8.8.8:53", "8.8.4.4:53")

// 设置 Basic 认证
client.SetBasicAuth("username", "password")

// 禁用 Keep-Alive
client.SetKeepalive(false)

// 使用环境变量代理
client.SetProxyFromEnv(true)

// 发起请求
body, code, resp, err := client.Get("https://example.com", 5*time.Second, nil, nil)
```

### TLS/SSL 配置

```go
client := ghttp.NewHTTPClient()

// 1. 固定证书（Certificate Pinning）
certPEM, _ := ioutil.ReadFile("server-cert.pem")
client.SetPinCert(certPEM)

// 2. 设置根证书
caCert, _ := ioutil.ReadFile("ca-cert.pem")
client.SetRootCaCerts(caCert)

// 3. 客户端证书认证
clientCert, _ := ioutil.ReadFile("client-cert.pem")
clientKey, _ := ioutil.ReadFile("client-key.pem")
client.SetClientCert(clientCert, clientKey)
```

### 请求拦截器

```go
client := ghttp.NewHTTPClient()

// 设置请求前拦截器
client.SetBeforeDo(func(req *http.Request) {
    // 在请求发送前修改请求
    req.Header.Set("X-Request-Time", time.Now().String())
    fmt.Printf("发送请求: %s\n", req.URL)
})

// 添加请求后拦截器
client.AppendAfterDo(func(req *http.Request, resp *http.Response, err error) {
    if err != nil {
        fmt.Printf("请求失败: %v\n", err)
    } else {
        fmt.Printf("响应状态: %d\n", resp.StatusCode)
    }
})

// 链式添加多个拦截器
client.AppendBeforeDo(func(req *http.Request) {
    // 第二个前置拦截器
}).AppendBeforeDo(func(req *http.Request) {
    // 第三个前置拦截器
})
```

### 自定义连接处理

```go
client := ghttp.NewHTTPClient()

// 自定义连接包装（可用于连接加密、日志等）
client.SetConnWrap(func(conn net.Conn) (net.Conn, error) {
    fmt.Printf("建立连接: %s -> %s\n", conn.LocalAddr(), conn.RemoteAddr())
    return conn, nil
})

// 自定义拨号器
client.SetDialer(func(network, address string, timeout time.Duration) (net.Conn, error) {
    fmt.Printf("拨号: %s %s\n", network, address)
    return net.DialTimeout(network, address, timeout)
})
```

### 自定义 HTTP 客户端工厂

```go
client := ghttp.NewHTTPClient()

// 完全自定义 http.Client
client.SetHttpClientFactory(func(r *http.Request) *http.Client {
    return &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
        },
    }
})
```

### 请求预处理

```go
client := ghttp.NewHTTPClient()

// 在请求发送前最后修改
client.SetPreHandler(func(r *http.Request) {
    // 添加签名
    signature := generateSignature(r)
    r.Header.Set("X-Signature", signature)
})
```

## 📋 API 参考

### 全局函数

| 函数 | 说明 | 参数 | 返回值 |
|------|------|------|--------|
| `Get(url, timeout, queryData, header)` | GET 请求 | url:请求地址<br>timeout:超时时间<br>queryData:查询参数<br>header:请求头 | body:响应体<br>code:状态码<br>resp:响应对象<br>err:错误 |
| `Post(url, data, timeout, header)` | POST 表单请求 | url:请求地址<br>data:表单数据<br>timeout:超时时间<br>header:请求头 | 同上 |
| `PostOfReader(url, reader, timeout, header)` | POST 流数据 | url:请求地址<br>reader:数据流<br>timeout:超时时间<br>header:请求头 | 同上 |
| `Upload(url, fieldName, filename, data, timeout, header)` | 上传文件 | url:上传地址<br>fieldName:字段名<br>filename:文件路径<br>data:附加数据<br>timeout:超时<br>header:请求头 | body:响应体<br>resp:响应对象<br>err:错误 |
| `UploadOfReader(url, fieldName, filename, reader, data, timeout, header)` | 流式上传 | 参数同上，reader 为 io.ReadCloser | 同上 |
| `Download(url, timeout, queryData, header)` | 下载到内存 | 同 Get 参数 | data:文件数据<br>resp:响应对象<br>err:错误 |
| `DownloadToFile(url, timeout, queryData, header, file)` | 下载到文件 | 同上 + file:保存路径 | resp:响应对象<br>err:错误 |
| `DownloadToWriter(url, timeout, queryData, header, writer)` | 下载到 Writer | 同上 + writer:io.Writer | 同上 |
| `NewTriableGet(url, maxTry, timeout, queryData, header)` | 创建可重试 GET | maxTry:最大重试次数 | tr:重试请求对象<br>err:错误 |
| `NewTriablePost(url, maxTry, timeout, data, header)` | 创建可重试 POST | 同上 | 同上 |
| `NewBatchGet(urls, timeout, data, header)` | 创建批量 GET | urls:URL 列表 | br:批量请求对象<br>err:错误 |
| `NewBatchPost(urls, timeout, data, header)` | 创建批量 POST | 同上 | 同上 |

### HTTPClient 方法

| 方法 | 说明 |
|------|------|
| `SetProxy(proxyURL)` | 设置代理服务器 |
| `SetProxyFromEnv(bool)` | 从环境变量获取代理 |
| `SetDNS(dns...)` | 设置自定义 DNS |
| `SetBasicAuth(user, pass)` | 设置 Basic 认证 |
| `SetKeepalive(bool)` | 启用/禁用 Keep-Alive |
| `SetPinCert(pemBytes)` | 设置固定证书 |
| `SetClientCert(cert, key)` | 设置客户端证书 |
| `SetRootCaCerts(caPem...)` | 设置根证书 |
| `SetBeforeDo(func)` | 设置请求前拦截器 |
| `AppendBeforeDo(func)` | 追加请求前拦截器 |
| `SetAfterDo(func)` | 设置请求后拦截器 |
| `AppendAfterDo(func)` | 追加请求后拦截器 |
| `SetHttpClientFactory(func)` | 自定义 HTTP 客户端工厂 |
| `SetPreHandler(func)` | 设置请求预处理器 |
| `SetConnWrap(func)` | 设置连接包装器 |
| `SetDialer(func)` | 设置自定义拨号器 |

### TriableRequest 方法

| 方法 | 说明 |
|------|------|
| `Execute()` | 执行请求（返回 Response） |
| `Keepalive(bool)` | 设置 Keep-Alive |
| `CheckErrorFunc(func)` | 自定义错误检查 |
| `SetBeforeDo(func)` | 设置前置拦截器 |
| `AfterDo(func)` | 添加后置拦截器 |
| `Err()` | 获取首个错误 |
| `ErrAll()` | 获取所有错误 |

### BatchRequest 方法

| 方法 | 说明 |
|------|------|
| `Execute()` | 执行批量请求 |
| `MaxTry(int)` | 设置重试次数 |
| `WaitFirstSuccess()` | 等待首个成功 |
| `Pool(gpool.Pool)` | 使用协程池（支持 BasicPool 或 OptimizedPool） |
| `Success()` | 是否全部成功 |
| `Resp()` | 获取首个成功响应 |
| `RespAll()` | 获取所有响应 |
| `Err()` | 获取首个错误 |
| `ErrAll()` | 获取所有错误 |
| `ErrorCount()` | 失败数量 |

### Response 对象

| 属性/方法 | 说明 |
|----------|------|
| `StatusCode` | HTTP 状态码 |
| `Body()` | 响应体字节数组 |
| `Err()` | 请求错误 |
| `UsedTime()` | 请求耗时 |
| `StartTime` | 开始时间 |
| `EndTime` | 结束时间 |
| `Close()` | 关闭响应 |

## 💡 使用场景

### 1. RESTful API 调用

```go
// GET 请求
users, _, _, _ := ghttp.Get("https://api.example.com/users", 5*time.Second, nil, nil)

// POST 创建资源
body, code, _, err := ghttp.Post(
    "https://api.example.com/users",
    map[string]string{"name": "张三"},
    5*time.Second,
    nil,
)
```

### 2. 网络爬虫

```go
// 批量抓取网页
urls := []string{
    "https://news.example.com/page1",
    "https://news.example.com/page2",
    // ... 更多 URL
}

br, _ := ghttp.NewBatchGet(urls, 10*time.Second, nil, nil)
br.Execute()

for _, resp := range br.RespAll() {
    if resp.Err() == nil {
        // 解析网页内容
        parseHTML(resp.Body())
    }
}
```

### 3. 文件同步

```go
// 下载远程文件
resp, err := ghttp.DownloadToFile(
    "https://cdn.example.com/data.zip",
    60*time.Second,
    nil, nil,
    "/tmp/data.zip",
)

// 上传文件
ghttp.Upload(
    "https://api.example.com/sync",
    "file",
    "/tmp/data.zip",
    nil,
    60*time.Second,
    nil,
)
```

### 4. 服务健康检查

```go
tr, _ := ghttp.NewTriableGet(
    "https://service.example.com/health",
    3,  // 重试 3 次
    2*time.Second,
    nil, nil,
)

resp := tr.Execute()
if resp.Err() != nil {
    // 服务不可用，发送告警
    sendAlert("服务异常")
}
```

### 5. 微服务调用

```go
client := ghttp.NewHTTPClient()

// 配置服务发现
client.SetDNS("consul.example.com:8600")

// 配置超时和重试
client.SetKeepalive(true)

// 调用服务
body, _, _, _ := client.Get(
    "http://user-service/api/users/123",
    3*time.Second,
    nil,
    map[string]string{"X-Request-ID": uuid.New().String()},
)
```

### 6. 数据采集

```go
// 定时采集多个数据源
ticker := time.NewTicker(1 * time.Minute)
for range ticker.C {
    sources := []string{
        "https://data1.example.com/metrics",
        "https://data2.example.com/metrics",
        "https://data3.example.com/metrics",
    }
    
    br, _ := ghttp.NewBatchGet(sources, 5*time.Second, nil, nil)
    br.Execute()
    
    for _, resp := range br.RespAll() {
        if resp.Err() == nil {
            processMetrics(resp.Body())
        }
    }
}
```

### 7. 代理请求

```go
client := ghttp.NewHTTPClient()
client.SetProxy("socks5://127.0.0.1:1080")

// 通过代理访问
body, _, _, _ := client.Get("https://blocked-site.com", 10*time.Second, nil, nil)
```

## ⚠️ 注意事项

1. **超时设置**：超时时间包括连接建立、请求发送和响应接收的总时间
2. **Cookie 管理**：全局客户端自动管理 Cookie，独立客户端需要手动管理
3. **连接复用**：默认启用 Keep-Alive，可提高性能
4. **并发控制**：批量请求默认并发执行，使用 `Pool()` 可限制并发数
5. **重试机制**：重试不会自动递增延迟，如需要请在拦截器中实现
6. **资源释放**：大文件下载建议使用 `DownloadToFile` 而非 `Download`
7. **错误处理**：重试请求的 `Err()` 返回首个错误，`ErrAll()` 返回所有错误
8. **响应关闭**：Response 对象使用完记得调用 `Close()` 释放资源

## 📚 相关链接

- 主项目：[GMC Framework](https://github.com/snail007/gmc)
- 问题反馈：[GitHub Issues](https://github.com/snail007/gmc/issues)
- 协程池：[gpool](https://github.com/snail007/gmc/tree/master/util/gpool)

## 📄 许可证

MIT License - 详见项目 LICENSE 文件
