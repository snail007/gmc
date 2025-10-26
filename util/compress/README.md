# gcompress

GMC 框架的压缩包处理工具库，支持自动识别和解压多种压缩格式。

## 特性

- 🔍 **智能识别**：基于魔数（magic number）自动识别压缩格式，无需指定文件类型
- 📦 **多格式支持**：支持 tar.gz、tar.bz2、tar.xz、zip、tar 等常见压缩格式
- 🌊 **流式处理**：支持从文件、HTTP 响应流等任意 io.Reader 直接解压
- 🔒 **安全防护**：内置路径清理机制，防止目录遍历攻击
- ⚡ **高性能**：使用并行 gzip (pgzip) 提升解压速度
- 🚫 **零 CGO**：纯 Go 实现，无 CGO 依赖，跨平台编译友好

## 支持的压缩格式

| 格式 | 扩展名 | 魔数 | 说明 |
|------|--------|------|------|
| Gzip | .gz, .tar.gz, .tgz | `1f 8b` | GNU zip 压缩 |
| Bzip2 | .bz2, .tar.bz2 | `42 5a` | Bzip2 压缩 |
| XZ | .xz, .tar.xz | `fd 37 7a 58 5a 00` | XZ 压缩 |
| Zip | .zip | `50 4b 03 04` | ZIP 归档 |
| Tar | .tar | `75 73 74 61 72` (偏移 257) | TAR 归档 |

## 安装

```bash
go get github.com/snail007/gmc/util/compress
```

## 使用方法

### 基础用法

```go
import gcompress "github.com/snail007/gmc/util/compress"

// 从文件解压
file, err := os.Open("archive.tar.gz")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

// 解压到指定目录
destPath, err := gcompress.Unpack(file, "/path/to/dest")
if err != nil {
    log.Fatal(err)
}
fmt.Println("解压到:", destPath)
```

### 从 HTTP 流解压

```go
import (
    "net/http"
    gcompress "github.com/snail007/gmc/util/compress"
)

// 从网络下载并直接解压
resp, err := http.Get("https://example.com/archive.tar.gz")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

destPath, err := gcompress.Unpack(resp.Body, "/path/to/dest")
if err != nil {
    log.Fatal(err)
}
```

### 解压到临时目录

```go
// 传入空字符串，自动创建临时目录
destPath, err := gcompress.Unpack(reader, "")
if err != nil {
    log.Fatal(err)
}
// destPath 将是 /tmp/unpackit-xxx/ 格式的临时目录
defer os.RemoveAll(destPath) // 使用完记得清理
```

### 单独使用 Unzip 和 Untar

```go
// 仅解压 ZIP 文件
file, _ := os.Open("archive.zip")
destPath, err := gcompress.Unzip(file, "/path/to/dest")

// 仅解压 TAR 文件
file, _ := os.Open("archive.tar")
destPath, err := gcompress.Untar(file, "/path/to/dest")
```

## 工作原理

1. **魔数检测**：读取文件头的魔数识别压缩格式
2. **解压缩**：根据格式使用相应的解压器（gzip/bzip2/xz）
3. **解归档**：如果是归档格式（tar/zip），继续解包提取文件
4. **路径清理**：清理文件路径，防止安全问题（如 `../` 攻击）
5. **权限保持**：保留原始文件的权限和时间戳

## API 文档

### Unpack(reader io.Reader, destPath string) (string, error)

通用解压函数，自动识别格式并解压。

**参数：**
- `reader`：输入流，可以是文件、HTTP 响应体等任意 io.Reader
- `destPath`：目标目录路径，空字符串则自动创建临时目录

**返回：**
- `string`：实际解压的目标路径
- `error`：错误信息

### Unzip(reader io.Reader, destPath string) (string, error)

解压 ZIP 文件。

### Untar(reader io.Reader, destPath string) (string, error)

解包 TAR 归档文件。

## 依赖库

- `github.com/klauspost/pgzip` - 并行 gzip 实现，提升性能
- `github.com/dsnet/compress/bzip2` - bzip2 压缩支持
- `github.com/ulikunitz/xz` - XZ 压缩支持

## 注意事项

1. **权限处理**：在 Unix 系统上会尝试设置原始文件权限，Windows 上可能部分失效
2. **符号链接**：当前实现不处理符号链接
3. **大文件**：内存使用优化良好，支持流式处理大文件
4. **错误处理**：部分非关键错误（如权限设置失败）仅记录日志不中断流程

## 安全性

- ✅ 自动清理恶意路径（如 `../../../etc/passwd`）
- ✅ Windows 盘符处理（去除 `C:` 等盘符前缀）
- ✅ 路径规范化，防止目录遍历
- ⚠️ 建议在隔离环境中处理不受信任的压缩包

## 许可证

本项目采用 MIT 和 MPL 2.0 双重许可证。

- MIT License - 详见项目 LICENSE 文件
- Mozilla Public License 2.0 - 详见源码头部声明

## 相关链接

- 主项目：[GMC Framework](https://github.com/snail007/gmc)
- 问题反馈：[GitHub Issues](https://github.com/snail007/gmc/issues)
