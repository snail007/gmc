<div align="center">

# GMC

<img src="/doc/images/logo2.png" width="200" alt="GMC Logo"/>

### 🚀 Modern Go Web & API Development Framework

A smart, flexible, and high-performance Golang Web and API development framework

[![Actions Status](https://github.com/snail007/gmc/workflows/build/badge.svg)](https://github.com/snail007/gmc/actions)
[![codecov](https://codecov.io/gh/snail007/gmc/branch/master/graph/badge.svg)](https://codecov.io/gh/snail007/gmc)
[![Go Report](https://goreportcard.com/badge/github.com/snail007/gmc)](https://goreportcard.com/report/github.com/snail007/gmc)
[![API Reference](https://img.shields.io/badge/go.dev-reference-blue)](https://pkg.go.dev/github.com/snail007/gmc)
[![LICENSE](https://img.shields.io/github/license/snail007/gmc)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/snail007/gmc)](go.mod)

English | [简体中文](README.md)

[📖 Documentation](https://snail007.github.io/gmc/) | [🎯 Quick Start](#-quick-start) | [💡 Features](#-core-features) | [🔧 Examples](#-code-examples)

</div>

---

## 📋 Table of Contents

- [Introduction](#-introduction)
- [Core Features](#-core-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Architecture](#-architecture)
- [Code Examples](#-code-examples)
- [Core Components](#-core-components)
- [Utility Packages](#-utility-packages)
- [Configuration](#-configuration)
- [Performance](#-performance)
- [Project Structure](#-project-structure)
- [Contributing](#-contributing)
- [License](#-license)
- [Contact](#-contact)

---

## 🎯 Introduction

**GMC** (Go Micro Container) is a full-stack Golang framework designed for modern web development. It is committed to providing:

- 🎨 **High Productivity** - Write less code to accomplish more
- ⚡ **High Performance** - Built on high-performance routing and optimized middleware
- 🧩 **Modular Design** - Clean architecture with comprehensive dependency injection
- 🛠️ **Rich Toolset** - 60+ out-of-the-box utility packages
- 📦 **Easy to Use** - Intuitive API design with detailed documentation

GMC is not just a web framework, but a complete development toolkit suitable for various scenarios from small APIs to large enterprise applications.

---

## ✨ Core Features

### 🌐 Web & API Development
- **RESTful API** - Quickly build RESTful style API services
- **MVC Architecture** - Complete MVC pattern support with clear code organization
- **Routing System** - High-performance routing engine with groups, parameters, and middleware
- **Controllers** - Elegant controller design with dependency injection
- **Template Engine** - Built-in template engine with layouts, inheritance, and custom functions

### 🗄️ Data Management
- **Multi-Database Support** - MySQL, SQLite3 out of the box
- **ORM Integration** - Elegant database operation interface
- **Cache System** - Multiple cache backends: Memory, Redis, File
- **Session Management** - Flexible session management mechanism

### 🔧 Development Tools
- **Configuration Management** - Support for TOML, JSON, YAML and more
- **Logging System** - Leveled logging, async writing, auto-rotation
- **Error Handling** - Complete error stack and error chain
- **Internationalization** - i18n support for easy multi-language implementation
- **CAPTCHA** - Built-in CAPTCHA generator
- **Paginator** - Ready-to-use pagination component

### ⚙️ Advanced Features
- **Middleware** - Flexible middleware system
- **Goroutine Pool** - High-performance goroutine pool management
- **Rate Limiting** - Built-in rate limiting and circuit breaker
- **Performance Profiling** - pprof integration for convenient performance analysis
- **Process Management** - Daemon process and graceful restart support
- **Dependency Injection** - Clear dependency injection mechanism
- **Hot Compilation** - Auto compile and restart during development (gmct run)
- **Resource Packaging** - Pack static files, templates, i18n into binary (gmct)

### 🛠️ Utility Libraries (60+)
Covering file operations, network tools, encryption/hashing, type conversion, collections, compression, JSON processing, and more.

### 🔨 GMCT Toolchain
- **Project Generation** - Generate Web/API project scaffolding with one command
- **Hot Compilation** - Auto compile and restart during development
- **Resource Packaging** - Pack static files, templates, i18n into binary
- **Project Management** - Various tools to simplify development workflow

---

## 📦 Installation

### Requirements

- Go 1.16 or higher

### Install Framework

```bash
go get -u github.com/snail007/gmc
```

### Install GMCT Toolchain

**GMCT** is the official CLI tool for GMC, providing project scaffolding, hot compilation, resource packaging and more:

```bash
# Install gmct
go install github.com/snail007/gmct@latest

# Verify installation
gmct version
```

#### Quick Install GMCT (Linux/macOS)

```bash
# Linux AMD64
bash -c "$(curl -L https://github.com/snail007/gmct/raw/master/install.sh)" @ linux-amd64

# Linux ARM64
bash -c "$(curl -L https://github.com/snail007/gmct/raw/master/install.sh)" @ linux-arm64

# macOS - Please download from Release page
# https://github.com/snail007/gmct/releases
```

📖 **GMCT Full Documentation**: [https://github.com/snail007/gmct](https://github.com/snail007/gmct)

### Verify Installation

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
)

func main() {
    fmt.Println("GMC framework installed successfully!")
}
```

---

## 🚀 Quick Start

### Create Project with GMCT (Recommended)

GMCT is the official toolchain for GMC that quickly generates project scaffolding:

```bash
# Create Web project
mkdir myapp && cd myapp
gmct new web

# Or create API project
gmct new api

# Run with hot compilation (recommended for development)
gmct run

# Visit http://localhost:7080
```

Generated project structure:
```
myapp/
├── conf/
│   └── app.toml          # Configuration
├── controller/
│   └── demo.go           # Controller
├── initialize/
│   └── initialize.go     # Initialization
├── router/
│   └── router.go         # Routes
├── static/               # Static files
├── views/                # Templates
├── grun.toml            # GMCT config
└── main.go              # Entry point
```

### Manual Project Creation

### 1. Create a Simple API Service

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
    gmap "github.com/snail007/gmc/util/map"
)

func main() {
    // Create API server
    api, _ := gmc.New.APIServer(gmc.New.Ctx(), ":8080")

    // Register route
    api.API("/", func(c gmc.C) {
        c.Write(gmap.M{
            "code":    0,
            "message": "Hello GMC!",
            "data":    nil,
        })
    })

    // Create app and run
    app := gmc.New.App()
    app.AddService(gcore.ServiceItem{
        Service: api.(gcore.Service),
    })
    
    app.Run()
}
```

After running, visit `http://localhost:8080/` to see the JSON response.

### 2. Create a Web Application

```go
package main

import (
    "github.com/snail007/gmc"
    gcore "github.com/snail007/gmc/core"
)

type HomeController struct {
    gmc.Controller
}

func (c *HomeController) Index() {
    c.Write("Welcome to GMC!")
}

func main() {
    // Create application
    app := gmc.New.App()
    
    // Create HTTP server
    s := gmc.New.HTTPServer(app.Ctx())
    s.Router().Controller("/", new(HomeController))
    
    // Add service and run
    app.AddService(gcore.ServiceItem{
        Service: s,
    })
    
    app.Run()
}
```

---

## 🏗️ Architecture

GMC adopts a clear modular architecture consisting of the following main parts:

```
gmc/
├── core/               # Core interface definitions
├── module/             # Feature module implementations
│   ├── app/           # Application framework
│   ├── cache/         # Cache (Memory, Redis, File)
│   ├── config/        # Configuration management
│   ├── db/            # Database (MySQL, SQLite3)
│   ├── log/           # Logging system
│   ├── i18n/          # Internationalization
│   └── middleware/    # Middleware
├── http/              # HTTP related
│   ├── server/        # HTTP/API server
│   ├── router/        # Router
│   ├── controller/    # Controller
│   ├── session/       # Session management
│   ├── template/      # Template engine
│   └── cookie/        # Cookie handling
├── util/              # Utility packages (60+ independent tools)
│   ├── gpool/         # Goroutine pool
│   ├── captcha/       # CAPTCHA
│   ├── cast/          # Type conversion
│   ├── compress/      # Compression/decompression
│   ├── file/          # File operations
│   ├── http/          # HTTP utilities
│   ├── json/          # JSON utilities
│   ├── rate/          # Rate limiter
│   └── ...            # More tools
└── using/             # Dependency injection registration
```

For detailed architecture description, see [ARCHITECTURE.md](ARCHITECTURE.md)

---

## 🔨 GMCT Toolchain

GMCT is the official CLI tool for GMC framework, providing project scaffolding, hot compilation, resource packaging and other powerful features to greatly improve development efficiency.

### 🎯 Main Features

#### 1. Project Generation

Quickly generate standardized project structure:

```bash
# Generate Web project (MVC architecture)
gmct new web

# Generate API project (lightweight)
gmct new api

# Specify package name
gmct new web --pkg github.com/yourname/myapp
```

#### 2. Hot Compilation Development

Automatically watch file changes, compile and restart during development:

```bash
# Run with hot compilation
gmct run

# Configuration file grun.toml
[run]
# Watch file extensions
watch_ext = [".go", ".toml"]
# Exclude directories
exclude_dir = ["vendor", ".git"]
# Build command
build_cmd = "go build -o tmp/app"
# Run command
run_cmd = "./tmp/app"
```

#### 3. Resource Packaging

Pack static files, templates, i18n files into binary for single-file deployment:

```bash
# Pack template files
gmct tpl --dir ./views

# Pack static files
gmct static --dir ./static

# Pack i18n files
gmct i18n --dir ./i18n

# Clean packed files
gmct static --clean
gmct tpl --clean
gmct i18n --clean
```

After packaging, your application can be compiled into a single binary without carrying any resource files.

#### 4. Project Information

```bash
# Show version
gmct version

# Show help
gmct help

# Show command help
gmct new --help
gmct run --help
```

### 📋 GMCT Command List

| Command | Description | Example |
|---------|-------------|---------|
| `gmct new` | Create new project | `gmct new web` |
| `gmct run` | Run with hot compilation | `gmct run` |
| `gmct tpl` | Pack templates | `gmct tpl --dir ./views` |
| `gmct static` | Pack static files | `gmct static --dir ./static` |
| `gmct i18n` | Pack i18n files | `gmct i18n --dir ./i18n` |
| `gmct version` | Show version | `gmct version` |
| `gmct help` | Show help | `gmct help` |

### 🎬 Complete Development Workflow Example

```bash
# 1. Install GMCT
go install github.com/snail007/gmct@latest

# 2. Create new project
mkdir mywebapp && cd mywebapp
gmct new web --pkg github.com/me/mywebapp

# 3. Hot compilation development
gmct run
# Auto recompile and restart after code changes

# 4. Pack resources (production)
gmct static --dir ./static
gmct tpl --dir ./views
gmct i18n --dir ./i18n

# 5. Build for release
go build -ldflags "-s -w" -o myapp

# 6. Deploy
./myapp
# Single binary file with all resources included
```

### ⚙️ Configuration File grun.toml

GMCT run configuration example:

```toml
[run]
# Watch file extensions
watch_ext = [".go", ".toml", ".html", ".js", ".css"]

# Exclude directories
exclude_dir = [
    "vendor",
    ".git",
    ".idea",
    "tmp",
    "bin",
]

# Commands before build
before_build = []

# Build command
build_cmd = "go build -o tmp/app"

# Run command
run_cmd = "./tmp/app"

# Commands after run
after_run = []

# Restart delay (milliseconds)
restart_delay = 1000
```

### 🌟 GMCT Advantages

1. **Improve Development Efficiency** - Hot compilation saves manual restart hassle
2. **Standardized Projects** - Unified project structure for better team collaboration
3. **Simplified Deployment** - Single file deployment after resource packaging
4. **Lower Learning Curve** - Out-of-the-box best practices
5. **Flexible Configuration** - Customizable build and run process

📖 **Full Documentation**: [GMCT Toolchain Repository](https://github.com/snail007/gmct)

---

## 💡 Code Examples

### API Routing

```go
api, _ := gmc.New.APIServer(gmc.New.Ctx(), ":8080")

// GET request
api.API("/user/:id", func(c gmc.C) {
    id := c.Param().ByName("id")
    c.Write(gmap.M{
        "user_id": id,
        "name":    "John Doe",
    })
})

// POST request
api.API("/user", func(c gmc.C) {
    name := c.Request().FormValue("name")
    // Handle business logic
    c.Write(gmap.M{"status": "created", "name": name})
}, "POST")
```

### Controller

```go
type UserController struct {
    gmc.Controller
}

func (c *UserController) List() {
    users := []string{"Alice", "Bob", "Charlie"}
    c.Write(users)
}

func (c *UserController) Detail() {
    id := c.Param().ByName("id")
    c.Write("User ID: " + id)
}

// Register in router
router.Controller("/user", new(UserController))
```

### Database Operations

```go
// Initialize database
gmc.DB.Init(cfg)
db := gmc.DB.DB()

// Query using ActiveRecord (Recommended)
ar := db.AR()
ar.From("users").Where(gdb.M{"age >": 18}).OrderBy("created_at", "DESC")
rs, err := db.Query(ar)

// Insert data
ar = db.AR()
ar.Insert("users", gdb.M{
    "name":  "John",
    "email": "john@example.com",
    "age":   25,
})
result, err := db.Exec(ar)
lastID := result.LastInsertId()

// Update data
ar = db.AR()
ar.Update("users", gdb.M{"age": 26}, gdb.M{"id": lastID})
db.Exec(ar)
```

📖 **Full Documentation**: [Database Module Guide](module/db/README.md)

### Cache Usage

```go
// Initialize cache
gmc.Cache.Init(cfg)
cache := gmc.Cache.Cache()

// Set cache (expires in 60 seconds)
cache.Set("key", "value", 60)

// Get cache
value, exists := cache.Get("key")

// Delete cache
cache.Del("key")
```

📖 **Full Documentation**: [Cache Module Guide](module/cache/README.md)

### Goroutine Pool

```go
import "github.com/snail007/gmc/util/gpool"

// Create goroutine pool (max 10 concurrent)
pool := gpool.New(10)

// Submit tasks
for i := 0; i < 100; i++ {
    pool.Submit(func() {
        // Execute task
    })
}

// Wait for all tasks to complete
pool.Wait()
```

📖 **Full Documentation**: [Goroutine Pool Guide](util/gpool/README.md)

### CAPTCHA Generation

```go
import "github.com/snail007/gmc/util/captcha"

// Create CAPTCHA
cap := gcaptcha.NewDefault()
img, code := cap.Create(4, gcaptcha.NUM)

// img is the CAPTCHA image data
// code is the CAPTCHA text
```

📖 **Full Documentation**: [CAPTCHA Utility Guide](util/captcha/README.md)

### Rate Limiter

```go
import "github.com/snail007/gmc/util/rate"

// Create rate limiter (100 requests per second)
limiter := grate.NewLimiter(100, 1)

if limiter.Allow() {
    // Handle request
} else {
    // Request rate limited
}
```

📖 **Full Documentation**: [Rate Limiter Guide](util/rate/README.md)

---

### 🔗 More Examples and Documentation

#### Core Modules
- [Application Framework (App)](module/app/README.md) - Application lifecycle management
- [Configuration (Config)](module/config/README.md) - Multi-format configuration support
- [Logging System (Log)](module/log/README.md) - Powerful logging functionality
- [Error Handling (Error)](module/error/README.md) - Error stack and error chain
- [Internationalization (i18n)](module/i18n/README.md) - Multi-language support
- [Middleware](module/middleware/README.md) - Middleware system

#### Utility Packages (Selected)
- [File Operations (File)](util/file/README.md) - File read/write, copy, move, etc.
- [Type Conversion (Cast)](util/cast/README.md) - Convert between various types
- [JSON Utilities (JSON)](util/json/README.md) - High-performance JSON processing
- [Compression (Compress)](util/compress/README.md) - gzip, tar, zip, etc.
- [HTTP Utilities (HTTP)](util/http/README.md) - HTTP client utilities
- [Network Utilities (Net)](util/net/README.md) - Network-related utility functions
- [Hash Utilities (Hash)](util/hash/README.md) - MD5, SHA, bcrypt, etc.
- [String Utilities (Strings)](util/strings/README.md) - String processing tools
- [Collection Utilities (Collection)](util/collection/README.md) - Collection operations
- [Performance Profiling (Pprof)](util/pprof/README.md) - Performance analysis tools

**📚 View All Packages**: [util/](util/)

**🎓 Complete Examples**: The [demos/](demos/) directory contains complete example code for various use cases

---

## 🧩 Core Components

### HTTP Server

GMC provides two types of HTTP servers:

- **HTTPServer** - Full-featured web server with MVC, templates, sessions, etc.
- **APIServer** - Lightweight API server focused on RESTful API development

### Routing System

- High-performance route matching
- Path parameters support `/user/:id`
- Wildcard support `/files/*filepath`
- Route groups and middleware
- RESTful route design

### Middleware

```go
// Global middleware
api.Use(func(c gmc.C, next func()) {
    // Pre-processing
    start := time.Now()
    
    next() // Call next handler
    
    // Post-processing
    duration := time.Since(start)
    log.Printf("Request took %v", duration)
})
```

### Template Engine

```go
// Render template
c.View().Render("user/profile", gmap.M{
    "name": "John",
    "age":  25,
})
```

---

## 🛠️ Utility Packages

GMC provides 60+ independent utility packages that can be used in any Go project:

| Category | Package | Description |
|----------|---------|-------------|
| 🔢 **Data Processing** | cast | Type conversion |
| | json | JSON operations |
| | collection | Collection operations |
| | set | Set data structure |
| | list | List operations |
| | map | Map utilities |
| 📁 **File & I/O** | file | File operations |
| | compress | Compression (gzip, tar, zip, xz) |
| | bytes | Byte handling |
| 🌐 **Network** | http | HTTP client utilities |
| | net | Network utilities |
| | proxy | Proxy utilities |
| | url | URL processing |
| 🔐 **Security** | hash | Hashing (MD5, SHA, bcrypt) |
| | captcha | CAPTCHA generation |
| ⚡ **Concurrency** | gpool | Goroutine pool |
| | sync | Synchronization utilities |
| | rate | Rate limiter |
| | loop | Loop control |
| 🔧 **System** | process | Process management |
| | exec | Command execution |
| | os | OS utilities |
| | env | Environment variables |
| 📊 **Others** | paginator | Paginator |
| | pprof | Performance profiling |
| | args | Argument parsing |
| | rand | Random numbers |

Using utility packages independently:

```go
import "github.com/snail007/gmc/util/cast"

// Type conversion
str := gcast.ToString(123)
num := gcast.ToInt("456")
```

---

## ⚙️ Configuration

GMC supports multiple configuration formats (TOML, JSON, YAML). TOML format is recommended.

### Basic Configuration Example (app.toml)

```toml
# HTTP server configuration
[httpserver]
listen = ":8080"
tlsenable = false
tlscert = ""
tlskey = ""

# Template configuration
[template]
dir = "views"
ext = ".html"

# Database configuration
[database]
default = "mysql"

[database.mysql]
enable = true
driver = "mysql"
dsn = "root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True"
maxidle = 10
maxconns = 100
maxlifetimeseconds = 3600

# Cache configuration
[cache]
default = "redis"

[cache.redis]
enable = true
address = "127.0.0.1:6379"
password = ""
db = 0
timeout = 10

# Log configuration
[log]
level = "info"
output = "stdout"

# Session configuration
[session]
store = "memory"
ttl = 3600
```

### Load Configuration

```go
cfg := gmc.New.Config()
cfg.SetConfigFile("app.toml")
err := cfg.ReadInConfig()
```

---

## 📊 Performance

GMC performs excellently in performance tests:

```bash
# Run benchmarks
go test -bench=. -benchmem ./...
```

Key performance metrics:

- **Routing Performance** - High-speed route matching, supports tens of thousands of routes
- **Concurrency Handling** - Goroutine pool optimization for efficient concurrent task scheduling
- **Memory Usage** - Optimized memory allocation to reduce GC pressure
- **Throughput** - Maintains stable throughput under high concurrency

---

## 📂 Project Structure

Recommended project structure:

```
myapp/
├── main.go              # Application entry
├── app.toml            # Configuration file
├── controller/         # Controllers
│   ├── home.go
│   └── user.go
├── model/              # Data models
│   └── user.go
├── service/            # Business logic layer
│   └── user_service.go
├── middleware/         # Custom middleware
│   └── auth.go
├── router/             # Route configuration
│   └── router.go
├── initialize/         # Initialization logic
│   └── init.go
├── views/              # Template files
│   ├── layout.html
│   └── home/
│       └── index.html
└── static/             # Static files
    ├── css/
    ├── js/
    └── images/
```

---

## 🤝 Contributing

We welcome all forms of contributions! Before submitting a PR, please ensure:

### Code Standards

1. **Comments** - Add clear comments for public functions and types
2. **Testing** - Test coverage should reach 90% or above
3. **Examples** - Provide usage examples for public functions
4. **Benchmarks** - Add benchmarks for performance-critical code

### Required Package Files

Each package should contain the following files (where `xxx` is the package name):

| File | Description |
|------|-------------|
| xxx.go | Main file |
| xxx_test.go | Unit tests (coverage > 90%) |
| example_test.go | Example code |
| benchmark_test.go | Benchmark tests |
| doc.go | Package documentation |
| README.md | Test and benchmark results |

You can refer to the `util/gpool` package for detailed code standards.

### Contribution Process

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 Documentation

- **Full Documentation**: [https://snail007.github.io/gmc/](https://snail007.github.io/gmc/)
- **API Documentation**: [https://pkg.go.dev/github.com/snail007/gmc](https://pkg.go.dev/github.com/snail007/gmc)
- **Architecture**: [ARCHITECTURE.md](ARCHITECTURE.md)
- **Examples**: [demos/](demos/)

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 💬 Contact

- **GitHub Issues**: [https://github.com/snail007/gmc/issues](https://github.com/snail007/gmc/issues)
- **GitHub Discussions**: [https://github.com/snail007/gmc/discussions](https://github.com/snail007/gmc/discussions)

---

## ⭐ Star History

If this project helps you, please give us a Star ⭐

[![Star History Chart](https://api.star-history.com/svg?repos=snail007/gmc&type=Date)](https://star-history.com/#snail007/gmc&Date)

---

## 🙏 Acknowledgments

Thanks to all developers who have contributed to GMC!

---

<div align="center">

**[⬆ Back to Top](#gmc)**

Made with ❤️ by the GMC Team

</div>
