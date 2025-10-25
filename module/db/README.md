# GMC Database 模块

## 简介

GMC Database 模块提供了强大的数据库抽象层，支持 MySQL 和 SQLite3 数据库。提供 ActiveRecord 模式的 ORM 功能、连接池管理、查询缓存、事务支持等特性。

## 功能特性

- **多数据库支持**：MySQL、SQLite3
- **多数据源管理**：支持同时连接多个数据库
- **ActiveRecord 模式**：类似 Ruby on Rails 的 ORM
- **查询构建器**：链式调用构建 SQL 查询
- **连接池**：自动管理数据库连接池
- **查询缓存**：可选的查询结果缓存
- **事务支持**：完整的事务功能
- **表前缀**：支持表名前缀
- **SQLite3 加密**：支持加密的 SQLite3 数据库
- **灵活的操作方式**：可以使用 ActiveRecord 直接操作，也可以使用 Model 进行 ORM 映射

## 安装

```bash
go get github.com/snail007/gmc/module/db
```

## 快速开始

### 从配置初始化

```go
package main

import (
    "github.com/snail007/gmc"
)

func main() {
    // 加载配置
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    err := cfg.ReadInConfig()
    if err != nil {
        panic(err)
    }
    
    // 初始化数据库
    gmc.DB.Init(cfg)
    
    // 获取默认数据库
    db := gmc.DB.DB()
    
    // 使用数据库
    ar := db.AR()
    ar.Table("users")
    // ...
}
```

### 基本 CRUD 操作（使用 ActiveRecord）

**注意：不使用 Model 也可以直接操作数据库**

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    db := gmc.DB.DB()
    
    // 插入数据
    ar := db.AR()
    ar.Insert("users", gdb.M{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   30,
    })
    result, err := db.Exec(ar)
    if err != nil {
        panic(err)
    }
    fmt.Println("Inserted ID:", result.LastInsertId())
    fmt.Println("Rows affected:", result.RowsAffected())
    
    // 查询数据
    ar = db.AR()
    ar.From("users").Where(gdb.M{"id": 1})
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }
    
    // 遍历结果
    for _, row := range rs.Rows() {
        fmt.Printf("User: %s, Email: %s\n", row["name"], row["email"])
    }
    
    // 更新数据
    ar = db.AR()
    ar.Update("users", gdb.M{"age": 31}, gdb.M{"id": 1})
    db.Exec(ar)
    
    // 删除数据
    ar = db.AR()
    ar.Delete("users", gdb.M{"id": 1})
    db.Exec(ar)
}
```

### 使用 Model（可选的 ORM 方式）

除了 ActiveRecord，也可以使用 Model 进行 ORM 映射操作：

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

type User struct {
    gdb.Model
    ID         int    `column:"id" table:"users"`
    Name       string `column:"name"`
    Email      string `column:"email"`
    Age        int    `column:"age"`
    CreatedAt  string `column:"created_at"`
}

func main() {
    db := gmc.DB.DB()
    
    // 插入
    user := &User{
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   25,
    }
    err := user.Insert(db)
    if err != nil {
        panic(err)
    }
    fmt.Println("Inserted ID:", user.ID)
    
    // 查询单条
    user2 := &User{}
    err = user2.Load(db, "email = ?", "alice@example.com")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found: %s, Age: %d\n", user2.Name, user2.Age)
    
    // 更新
    user2.Age = 26
    err = user2.Update(db, "id = ?", user2.ID)
    if err != nil {
        panic(err)
    }
    
    // 删除
    err = user2.Delete(db, "id = ?", user2.ID)
    if err != nil {
        panic(err)
    }
}
```

### ActiveRecord 查询构建器

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    db := gmc.DB.DB()
    
    // 链式查询
    ar := db.AR()
    ar.From("users").
        Select("id, name, email").
        Where(gdb.M{"age >": 18}).
        Where(gdb.M{"status": "active"}).
        OrderBy("created_at", "DESC").
        Limit(0, 10)
    
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }
    
    // 获取所有行
    rows := rs.Rows()
    for _, row := range rows {
        fmt.Printf("User: %s, Email: %s\n", row["name"], row["email"])
    }
    
    // 使用原始 SQL
    rs2, err := db.QuerySQL("SELECT * FROM users WHERE age > ?", 18)
    if err != nil {
        panic(err)
    }
    fmt.Println("Found:", rs2.Len(), "users")
}
```

### 事务处理

```go
package main

import (
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    db := gmc.DB.DB()
    
    // 开始事务
    tx, err := db.Begin()
    if err != nil {
        panic(err)
    }
    
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    
    // 在事务中插入数据
    ar := db.AR()
    ar.Insert("users", gdb.M{
        "name":  "Bob",
        "email": "bob@example.com",
    })
    _, err = db.ExecTx(ar, tx)
    if err != nil {
        tx.Rollback()
        panic(err)
    }
    
    ar = db.AR()
    ar.Insert("orders", gdb.M{
        "user_id": 1,
        "amount":  100,
    })
    _, err = db.ExecTx(ar, tx)
    if err != nil {
        tx.Rollback()
        panic(err)
    }
    
    // 提交事务
    err = tx.Commit()
    if err != nil {
        panic(err)
    }
}
```

### 使用查询缓存

```go
package main

import (
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    db := gmc.DB.DB()
    
    // 启用缓存的查询
    ar := db.AR()
    ar.From("users").
        Select("*").
        Where(gdb.M{"id": 1}).
        Cache("user_1", 300) // 缓存键和过期时间（秒）
    
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }
    
    // 第二次查询会从缓存读取
    rs2, _ := db.Query(ar)
    _ = rs2
}
```

## 配置文件

### 完整配置示例

```toml
[database]
# 默认数据库 ID
default = "default"

# MySQL 配置
[[database.mysql]]
enable = true
id = "default"
host = "127.0.0.1"
port = 3306
username = "root"
password = "password"
database = "myapp"
# 表前缀
prefix = ""
# 表前缀占位符
prefix_sql_holder = "__PREFIX__"
# 字符集
charset = "utf8mb4"
# 排序规则
collate = "utf8mb4_general_ci"
# 最大空闲连接数
maxidle = 30
# 最大连接数
maxconns = 200
# 连接超时（毫秒）
timeout = 3000
# 读超时（毫秒）
readtimeout = 5000
# 写超时（毫秒）
writetimeout = 5000
# 连接最大生命周期（秒）
maxlifetimeseconds = 1800

# 第二个 MySQL 数据库（用于分离读写）
[[database.mysql]]
enable = true
id = "read_replica"
host = "192.168.1.100"
port = 3306
username = "readonly"
password = "password"
database = "myapp"
prefix = ""
prefix_sql_holder = "__PREFIX__"
charset = "utf8mb4"
collate = "utf8mb4_general_ci"
maxidle = 20
maxconns = 100
timeout = 3000
readtimeout = 5000
writetimeout = 5000
maxlifetimeseconds = 1800

# SQLite3 配置
[[database.sqlite3]]
enable = true
id = "local"
database = "./data/app.db"
# 密码（为空则不加密）
password = ""
# 表前缀
prefix = ""
prefix_sql_holder = "__PREFIX__"
# 同步模式: 0=OFF, 1=NORMAL, 2=FULL, 3=EXTRA
syncmode = 1
# 打开模式: ro, rw, rwc, memory
openmode = "rwc"
# 缓存模式: shared, private
cachemode = "shared"
```

## API 参考

### 全局函数

```go
// 初始化数据库系统
func Init(cfg gcore.Config) error

// 获取默认数据库或指定 ID 的数据库
func DB(id ...string) gcore.Database

// 获取 MySQL 数据库
func MySQL(id ...string) gcore.Database

// 获取 SQLite3 数据库
func SQLite3(id ...string) gcore.Database

// 设置日志记录器
func SetLogger(logger gcore.Logger)
```

### Database 接口

```go
type Database interface {
    // 创建 ActiveRecord
    AR() gcore.ActiveRecord
    
    // 执行查询
    Query(ar gcore.ActiveRecord) (gcore.ResultSet, error)
    QuerySQL(sql string, values ...interface{}) (gcore.ResultSet, error)
    
    // 执行更新/插入/删除
    Exec(ar gcore.ActiveRecord) (gcore.Result, error)
    ExecSQL(sql string, values ...interface{}) (gcore.Result, error)
    
    // 事务
    Begin() (*sql.Tx, error)
    ExecTx(ar gcore.ActiveRecord, tx *sql.Tx) (gcore.Result, error)
    ExecSQLTx(tx *sql.Tx, sql string, values ...interface{}) (gcore.Result, error)
    
    // 连接池统计
    Stats() sql.DBStats
}
```

### ActiveRecord 接口

```go
type ActiveRecord interface {
    // 表操作
    From(tableName string) gcore.ActiveRecord
    FromAs(from, as string) gcore.ActiveRecord
    
    // 查询构建
    Select(fields string) gcore.ActiveRecord
    SelectNoWrap(fields string) gcore.ActiveRecord
    Where(where gmap.M) gcore.ActiveRecord
    WhereRaw(where string) gcore.ActiveRecord
    WhereWrap(where gmap.M, leftWrap, rightWrap string) gcore.ActiveRecord
    GroupBy(fields string) gcore.ActiveRecord
    Having(having string) gcore.ActiveRecord
    HavingWrap(having, leftWrap, rightWrap string) gcore.ActiveRecord
    OrderBy(field, typ string) gcore.ActiveRecord
    Limit(limit ...int) gcore.ActiveRecord
    Join(table, as, on, typ string) gcore.ActiveRecord
    
    // 数据操作
    Insert(table string, data gmap.M) gcore.ActiveRecord
    InsertBatch(table string, data []gmap.M) gcore.ActiveRecord
    Replace(table string, data gmap.M) gcore.ActiveRecord
    ReplaceBatch(table string, data []gmap.M) gcore.ActiveRecord
    Update(table string, data, where gmap.M) gcore.ActiveRecord
    UpdateBatch(table string, values []gmap.M, whereColumn []string) gcore.ActiveRecord
    Set(column string, value interface{}) gcore.ActiveRecord
    SetNoWrap(column string, value interface{}) gcore.ActiveRecord
    Delete(table string, where gmap.M) gcore.ActiveRecord
    
    // 原始 SQL
    Raw(sql string, values ...interface{}) gcore.ActiveRecord
    
    // 缓存
    Cache(key string, seconds uint) gcore.ActiveRecord
    
    // 获取 SQL
    SQL() string
    Values() []interface{}
    
    // 重置
    Reset()
}
```

### Model 结构（可选）

如果使用 Model 方式，结构体需要嵌入 `gdb.Model`：

```go
type Model struct {
    gdb.Model
}

// Model 方法
func (m *Model) Insert(db gcore.Database) error
func (m *Model) Update(db gcore.Database, where string, args ...interface{}) error
func (m *Model) Delete(db gcore.Database, where string, args ...interface{}) error
func (m *Model) Load(db gcore.Database, where string, args ...interface{}) error
func (m *Model) LoadAll(db gcore.Database, where string, args ...interface{}) ([]interface{}, error)
```

### ResultSet 接口

```go
type ResultSet interface {
    // 遍历结果
    Next() bool
    
    // 获取当前行
    Row() map[string]interface{}
    
    // 获取所有行
    Rows() []map[string]interface{}
    
    // 映射到结构体
    MapRows(v interface{}) error
    MapRow(v interface{}) error
    
    // 列信息
    Columns() []string
    
    // 关闭结果集
    Close() error
}
```

## MySQL 特性

### 连接池配置

```go
// 在配置文件中设置
maxidle = 30              // 最大空闲连接
maxconns = 200            // 最大连接数
maxlifetimeseconds = 1800 // 连接最大生命周期
```

### 超时配置

```go
timeout = 3000      // 连接超时 3 秒
readtimeout = 5000  // 读超时 5 秒
writetimeout = 5000 // 写超时 5 秒
```

### 字符集配置

```go
charset = "utf8mb4"
collate = "utf8mb4_general_ci"
```

## SQLite3 特性

### 加密数据库

```toml
[[database.sqlite3]]
database = "./data/secure.db"
password = "my-secret-password"  # 设置密码启用加密
```

### 同步模式

```toml
# 0 = OFF     - 最快，有数据丢失风险
# 1 = NORMAL  - 平衡性能和安全（推荐）
# 2 = FULL    - 最安全，性能较低
# 3 = EXTRA   - 超级安全，性能最低
syncmode = 1
```

### 打开模式

```toml
# ro    - 只读
# rw    - 读写（数据库必须存在）
# rwc   - 读写（不存在则创建）
# memory - 内存数据库
openmode = "rwc"
```

## 使用场景

1. **Web 应用**：用户、文章、评论等数据存储
2. **API 服务**：RESTful API 的数据层
3. **微服务**：各微服务的数据存储
4. **数据分析**：数据查询和聚合
5. **缓存层**：配合缓存使用提升性能

## 最佳实践

### 1. 使用表前缀

```toml
prefix = "app_"
prefix_sql_holder = "__PREFIX__"
```

```go
// SQL 中使用占位符（如果配置了前缀）
ar.From("__PREFIX__users")  // 自动替换为 app_users

// 或者直接使用表名（会自动添加前缀）
ar.From("users")  // 也会替换为 app_users
```

### 2. 连接池优化

```go
// 根据并发量调整
maxidle = min(10, maxconns/10)    // 空闲连接约为最大连接的 10%
maxconns = 并发请求数 * 2           // 最大连接数为并发的 2 倍
maxlifetimeseconds = 1800          // 30 分钟回收连接
```

### 3. 使用事务保证一致性

```go
tx, _ := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 执行多个操作
// ...

tx.Commit()
```

### 4. 使用 Model 简化代码（可选）

```go
type User struct {
    gdb.Model
    ID   int    `column:"id" table:"users"`
    Name string `column:"name"`
}

// Model 方式
user := &User{Name: "John"}
user.Insert(db)

// 或使用 ActiveRecord 直接操作（推荐）
ar := db.AR()
ar.Insert("users", gdb.M{"name": "John"})
db.Exec(ar)
```

### 5. 查询缓存优化性能

```go
// 频繁查询且不常变化的数据使用缓存
ar.Cache("cache_key", 300) // 缓存 300 秒
```

## 性能优化

1. **连接池**：合理设置连接池大小
2. **索引**：为查询字段添加索引
3. **批量操作**：使用批量插入/更新
4. **查询缓存**：缓存不常变化的查询结果
5. **读写分离**：配置多个数据库实现读写分离

## 注意事项

1. **SQL 注入**：使用参数化查询，避免直接拼接 SQL
2. **连接泄漏**：确保 ResultSet 正确关闭
3. **事务管理**：及时提交或回滚事务
4. **字符集**：建议使用 utf8mb4
5. **连接数限制**：注意数据库服务器的最大连接数限制

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)