# GMC 数据库 AR 使用完整分析文档

## 目录
1. [GMC 数据库 AR 概述](#概述)
2. [ActiveRecord (AR) 详细用法](#activerecord-ar-详细用法)
3. [gdb 包单独使用配置](#gdb-包单独使用配置)
4. [框架集成使用配置](#框架集成使用配置)
5. [MySQL 配置详解](#mysql-配置详解)
6. [SQLite3 配置详解](#sqlite3-配置详解)
7. [完整示例代码](#完整示例代码)

---

## 概述

GMC 数据库模块 (`github.com/snail007/gmc/module/db`) 提供了强大的数据库抽象层，支持 MySQL 和 SQLite3 数据库。核心特性包括：

- **ActiveRecord 模式**：链式调用构建 SQL 查询
- **多数据库支持**：MySQL、SQLite3
- **多数据源管理**：支持同时连接多个数据库
- **连接池管理**：自动管理数据库连接池
- **查询缓存**：可选的查询结果缓存
- **事务支持**：完整的事务功能
- **表前缀**：支持表名前缀
- **SQLite3 加密**：支持加密的 SQLite3 数据库

---

## ActiveRecord (AR) 详细用法

### AR 接口定义

```go
type ActiveRecord interface {
    // 表操作
    From(tableName string) ActiveRecord
    FromAs(from, as string) ActiveRecord
    
    // 查询构建
    Select(fields string) ActiveRecord
    SelectNoWrap(fields string) ActiveRecord
    Where(where map[string]interface{}) ActiveRecord
    WhereRaw(where string) ActiveRecord
    WhereWrap(where map[string]interface{}, leftWrap, rightWrap string) ActiveRecord
    GroupBy(fields string) ActiveRecord
    Having(having string) ActiveRecord
    HavingWrap(having, leftWrap, rightWrap string) ActiveRecord
    OrderBy(field, typ string) ActiveRecord
    Limit(limit ...int) ActiveRecord
    Join(table, as, on, typ string) ActiveRecord
    
    // 数据操作
    Insert(table string, data map[string]interface{}) ActiveRecord
    InsertBatch(table string, data []map[string]interface{}) ActiveRecord
    Replace(table string, data map[string]interface{}) ActiveRecord
    ReplaceBatch(table string, data []map[string]interface{}) ActiveRecord
    Update(table string, data, where map[string]interface{}) ActiveRecord
    UpdateBatch(table string, values []map[string]interface{}, whereColumn []string) ActiveRecord
    Set(column string, value interface{}) ActiveRecord
    SetNoWrap(column string, value interface{}) ActiveRecord
    Delete(table string, where map[string]interface{}) ActiveRecord
    
    // 原始 SQL
    Raw(sql string, values ...interface{}) ActiveRecord
    
    // 缓存
    Cache(key string, seconds uint) ActiveRecord
    
    // 获取 SQL
    SQL() string
    Values() []interface{}
    
    // 重置
    Reset()
}
```

### 1. 查询操作 (SELECT)

#### 基本查询

```go
// 获取 AR 对象
ar := db.AR()

// 简单查询
ar.From("users")
rs, err := db.Query(ar)

// 指定字段查询
ar.From("users").Select("id, name, email")
rs, err := db.Query(ar)

// 查询所有字段
ar.From("users").Select("*")
rs, err := db.Query(ar)

// 使用 SelectNoWrap 不包装字段名（用于函数、表达式等）
ar.From("users").SelectNoWrap("COUNT(*) as total, MAX(age) as max_age")
rs, err := db.Query(ar)
```

#### WHERE 条件

```go
import gdb "github.com/snail007/gmc/module/db"

// 等值查询
ar.From("users").Where(gdb.M{"id": 1})
ar.From("users").Where(gdb.M{"status": "active"})

// 多条件 AND
ar.From("users").
    Where(gdb.M{"age >": 18}).
    Where(gdb.M{"status": "active"})

// 多条件 OR（使用 WhereWrap）
ar.From("users").
    WhereWrap(gdb.M{"age >": 18}, "OR", "").
    WhereWrap(gdb.M{"age <": 65}, "OR", "")

// 复杂条件组合
ar.From("users").
    Where(gdb.M{"status": "active"}).
    WhereWrap(gdb.M{"age >": 18}, "AND", "(").
    WhereWrap(gdb.M{"age <": 65}, "OR", "").
    WhereWrap(gdb.M{"vip": true}, "OR", ")")

// 原始 SQL 条件
ar.From("users").WhereRaw("created_at > DATE_SUB(NOW(), INTERVAL 7 DAY)")

// IN 查询
ar.From("users").Where(gdb.M{"id": []int{1, 2, 3, 4, 5}})

// 范围查询
ar.From("users").Where(gdb.M{"age >=": 18, "age <=": 65})
```

#### JOIN 操作

```go
// INNER JOIN
ar.From("users").
    Join("orders", "o", "users.id = o.user_id", "INNER").
    Select("users.id, users.name, o.order_id, o.amount")

// LEFT JOIN
ar.From("users").
    Join("profiles", "p", "users.id = p.user_id", "LEFT").
    Select("users.*, p.bio, p.avatar")

// RIGHT JOIN
ar.From("users").
    Join("departments", "d", "users.dept_id = d.id", "RIGHT")

// 多表 JOIN
ar.From("users").
    Join("orders", "o", "users.id = o.user_id", "LEFT").
    Join("products", "p", "o.product_id = p.id", "LEFT").
    Select("users.name, o.order_id, p.product_name")
```

#### GROUP BY 和 HAVING

```go
// GROUP BY
ar.From("orders").
    Select("user_id, SUM(amount) as total").
    GroupBy("user_id")

// HAVING
ar.From("orders").
    Select("user_id, SUM(amount) as total").
    GroupBy("user_id").
    Having("SUM(amount) > 1000")

// HAVING 复杂条件
ar.From("orders").
    Select("user_id, COUNT(*) as order_count, SUM(amount) as total").
    GroupBy("user_id").
    HavingWrap("SUM(amount) > 1000", "AND", "").
    HavingWrap("COUNT(*) > 5", "AND", "")
```

#### ORDER BY 和 LIMIT

```go
// 单字段排序
ar.From("users").OrderBy("created_at", "DESC")

// 多字段排序
ar.From("users").
    OrderBy("status", "ASC").
    OrderBy("created_at", "DESC")

// LIMIT（限制数量）
ar.From("users").Limit(10)

// LIMIT（偏移量和数量）
ar.From("users").Limit(0, 10)  // 前 10 条
ar.From("users").Limit(10, 10) // 第 11-20 条（分页）

// 完整查询示例
ar.From("users").
    Select("id, name, email").
    Where(gdb.M{"status": "active"}).
    OrderBy("created_at", "DESC").
    Limit(0, 20)
```

### 2. 插入操作 (INSERT)

#### 单条插入

```go
ar := db.AR()
ar.Insert("users", gdb.M{
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   30,
    "status": "active",
})
result, err := db.Exec(ar)
if err != nil {
    panic(err)
}
fmt.Println("插入ID:", result.LastInsertID())
fmt.Println("影响行数:", result.RowsAffected())
```

#### 批量插入

```go
users := []gdb.M{
    {"name": "Alice", "email": "alice@example.com", "age": 25},
    {"name": "Bob", "email": "bob@example.com", "age": 30},
    {"name": "Charlie", "email": "charlie@example.com", "age": 35},
}

ar := db.AR()
ar.InsertBatch("users", users)
result, err := db.Exec(ar)
if err != nil {
    panic(err)
}
fmt.Println("插入ID:", result.LastInsertID())
fmt.Println("影响行数:", result.RowsAffected())
```

#### REPLACE 操作

```go
// REPLACE INTO（如果存在则替换）
ar := db.AR()
ar.Replace("users", gdb.M{
    "id":    1,
    "name":  "John Updated",
    "email": "john@example.com",
})
result, err := db.Exec(ar)

// 批量 REPLACE
ar = db.AR()
ar.ReplaceBatch("users", users)
result, err = db.Exec(ar)
```

### 3. 更新操作 (UPDATE)

#### 单条更新

```go
ar := db.AR()
ar.Update("users", 
    gdb.M{"age": 31, "status": "inactive"},  // 更新数据
    gdb.M{"id": 1},                          // WHERE 条件
)
result, err := db.Exec(ar)
if err != nil {
    panic(err)
}
fmt.Println("影响行数:", result.RowsAffected())
```

#### 使用 Set 方法

```go
ar := db.AR()
ar.From("users").
    Set("age", 31).
    Set("status", "inactive").
    Where(gdb.M{"id": 1})
result, err := db.Exec(ar)

// SetNoWrap 用于设置表达式（不转义）
ar = db.AR()
ar.From("users").
    Set("age", 31).
    SetNoWrap("updated_at", "NOW()").
    Where(gdb.M{"id": 1})
result, err = db.Exec(ar)
```

#### 批量更新

```go
users := []gdb.M{
    {"id": 1, "age": 31, "status": "active"},
    {"id": 2, "age": 32, "status": "active"},
    {"id": 3, "age": 33, "status": "inactive"},
}

ar := db.AR()
ar.UpdateBatch("users", users, []string{"id"})  // 根据 id 字段批量更新
result, err := db.Exec(ar)
```

### 4. 删除操作 (DELETE)

```go
// 单条删除
ar := db.AR()
ar.Delete("users", gdb.M{"id": 1})
result, err := db.Exec(ar)

// 条件删除
ar = db.AR()
ar.Delete("users", gdb.M{"status": "inactive", "age <": 18})
result, err = db.Exec(ar)

// 带排序和限制的删除
ar = db.AR()
ar.From("users").
    Where(gdb.M{"status": "inactive"}).
    OrderBy("created_at", "ASC").
    Limit(10)
ar.sqlType = "delete"  // 注意：需要设置 sqlType
result, err = db.Exec(ar)
```

### 5. 原始 SQL

```go
// 使用原始 SQL 查询
rs, err := db.QuerySQL("SELECT * FROM users WHERE age > ? AND status = ?", 18, "active")

// 使用 AR 的 Raw 方法
ar := db.AR()
ar.Raw("SELECT * FROM users WHERE age > ? AND status = ?", 18, "active")
rs, err := db.Query(ar)

// 使用原始 SQL 执行
result, err := db.ExecSQL("UPDATE users SET age = ? WHERE id = ?", 31, 1)

// 带表前缀的原始 SQL（会自动替换 __PREFIX__）
ar = db.AR()
ar.Raw("SELECT * FROM __PREFIX__users WHERE id = ?", 1)
rs, err = db.Query(ar)
```

### 6. 查询结果处理 (ResultSet)

```go
rs, err := db.Query(ar)
if err != nil {
    panic(err)
}

// 获取所有行（map 格式）
rows := rs.Rows()
for _, row := range rows {
    fmt.Printf("ID: %s, Name: %s, Email: %s\n", row["id"], row["name"], row["email"])
}

// 获取单行
row := rs.Row()
if row != nil {
    fmt.Printf("Name: %s\n", row["name"])
}

// 获取单个字段值
name := rs.Value("name")

// 获取某列的所有值
names := rs.Values("name")

// 获取键值对映射
nameEmailMap := rs.MapValues("id", "email")  // map[id]email

// 映射到结构体
type User struct {
    ID    int    `column:"id"`
    Name  string `column:"name"`
    Email string `column:"email"`
    Age   int    `column:"age"`
}

// 映射单条记录
var user User
userStruct, err := rs.Struct(&user)
if err == nil {
    u := userStruct.(*User)
    fmt.Printf("User: %+v\n", u)
}

// 映射多条记录
users, err := rs.Structs(&User{})
if err == nil {
    for _, u := range users {
        user := u.(*User)
        fmt.Printf("User: %+v\n", user)
    }
}

// 映射为 map（以某个字段为 key）
userMap := rs.MapRows("id")  // map[id]map[field]value

// 获取统计信息
fmt.Printf("查询耗时: %v\n", rs.TimeUsed())
fmt.Printf("结果数量: %d\n", rs.Len())
fmt.Printf("执行的SQL: %s\n", rs.SQL())
```

### 7. 查询缓存

```go
// 启用缓存的查询
ar := db.AR()
ar.From("users").
    Select("*").
    Where(gdb.M{"id": 1}).
    Cache("user_1", 300)  // 缓存键和过期时间（秒）

rs, err := db.Query(ar)
// 第二次查询会从缓存读取（如果缓存未过期）
rs2, err := db.Query(ar)
```

**注意**：使用缓存需要配置 DBCache，详见配置部分。

### 8. 表前缀

```go
// 如果配置了表前缀（如 "app_"），以下两种方式等价：

// 方式1：使用占位符（推荐）
ar.From("__PREFIX__users")  // 自动替换为 app_users

// 方式2：直接使用表名（会自动添加前缀）
ar.From("users")  // 也会替换为 app_users
```

### 9. AR 重置

```go
ar := db.AR()
ar.From("users").Where(gdb.M{"id": 1})
rs, _ := db.Query(ar)

// 重置 AR 对象以便重用
ar.Reset()

// 现在可以重新使用
ar.From("orders").Where(gdb.M{"user_id": 1})
rs2, _ := db.Query(ar)
```

---

## gdb 包单独使用配置

### 方式一：直接使用代码配置（不依赖配置文件）

#### MySQL 单独使用

```go
package main

import (
    "fmt"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    // 创建 MySQL 配置
    config := gdb.NewMySQLDBConfig()
    config.Host = "127.0.0.1"
    config.Port = 3306
    config.Username = "root"
    config.Password = "password"
    config.Database = "testdb"
    config.Charset = "utf8mb4"
    config.Collate = "utf8mb4_general_ci"
    config.MaxOpenConns = 200
    config.MaxIdleConns = 50
    config.Timeout = 3000        // 连接超时（毫秒）
    config.ReadTimeout = 5000     // 读超时（毫秒）
    config.WriteTimeout = 5000    // 写超时（毫秒）
    config.TablePrefix = "app_"   // 表前缀（可选）
    config.TablePrefixSQLIdentifier = "__PREFIX__"  // 前缀占位符（可选）

    // 创建数据库连接
    db, err := gdb.NewMySQLDB(config)
    if err != nil {
        panic(err)
    }

    // 使用数据库
    ar := db.AR()
    ar.From("users").Select("*")
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }

    rows := rs.Rows()
    for _, row := range rows {
        fmt.Printf("User: %s\n", row["name"])
    }
}
```

#### SQLite3 单独使用

```go
package main

import (
    "fmt"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    // 创建 SQLite3 配置
    config := gdb.NewSQLite3DBConfig()
    config.Database = "./data/app.db"
    config.OpenMode = "rwc"        // ro, rw, rwc, memory
    config.CacheMode = "shared"     // shared, private
    config.SyncMode = 1             // 0=OFF, 1=NORMAL, 2=FULL, 3=EXTRA
    config.TablePrefix = "app_"     // 表前缀（可选）
    config.TablePrefixSQLIdentifier = "__PREFIX__"  // 前缀占位符（可选）

    // 创建数据库连接
    db, err := gdb.NewSQLite3DB(config)
    if err != nil {
        panic(err)
    }

    // 使用数据库
    ar := db.AR()
    ar.From("users").Select("*")
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }

    rows := rs.Rows()
    for _, row := range rows {
        fmt.Printf("User: %s\n", row["name"])
    }
}
```

#### 使用 DBGroup 管理多个数据库

```go
package main

import (
    "fmt"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    // 创建 MySQL 组
    mysqlGroup := gdb.NewMySQLDBGroup("default")

    // 注册多个 MySQL 数据库
    config1 := gdb.NewMySQLDBConfig()
    config1.Host = "127.0.0.1"
    config1.Port = 3306
    config1.Username = "root"
    config1.Password = "password"
    config1.Database = "db1"
    err := mysqlGroup.Regist("db1", config1)
    if err != nil {
        panic(err)
    }

    config2 := gdb.NewMySQLDBConfig()
    config2.Host = "127.0.0.1"
    config2.Port = 3306
    config2.Username = "root"
    config2.Password = "password"
    config2.Database = "db2"
    err = mysqlGroup.Regist("db2", config2)
    if err != nil {
        panic(err)
    }

    // 使用不同的数据库
    db1 := mysqlGroup.DB("db1")
    db2 := mysqlGroup.DB("db2")

    // 在 db1 中查询
    ar1 := db1.AR()
    ar1.From("users").Select("*")
    rs1, _ := db1.Query(ar1)

    // 在 db2 中查询
    ar2 := db2.AR()
    ar2.From("products").Select("*")
    rs2, _ := db2.Query(ar2)
}
```

### 方式二：使用配置文件（但不使用框架）

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    // 加载配置文件
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    err := cfg.ReadInConfig()
    if err != nil {
        panic(err)
    }

    // 初始化数据库（从配置文件读取）
    err = gdb.Init(cfg)
    if err != nil {
        panic(err)
    }

    // 获取默认数据库
    db := gdb.DB()
    if db == nil {
        panic("数据库未初始化")
    }

    // 或者获取指定 ID 的数据库
    db = gdb.DB("default")

    // 使用数据库
    ar := db.AR()
    ar.From("users").Select("*")
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }

    rows := rs.Rows()
    for _, row := range rows {
        fmt.Printf("User: %s\n", row["name"])
    }
}
```

---

## 框架集成使用配置

### 使用 gmc 框架（推荐方式）

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

    // 初始化数据库（框架会自动读取配置文件中的 [database] 部分）
    gmc.DB.Init(cfg)

    // 获取默认数据库
    db := gmc.DB.DB()

    // 使用数据库
    ar := db.AR()
    ar.From("users").Select("*")
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }

    // 获取指定 ID 的数据库
    db2 := gmc.DB.DB("read_replica")

    // 获取 MySQL 数据库
    mysqlDB := gmc.DB.MySQL("default")

    // 获取 SQLite3 数据库
    sqliteDB := gmc.DB.SQLite3("default")
}
```

---

## MySQL 配置详解

### 配置文件格式（TOML）

```toml
[database]
# 默认数据库类型（mysql 或 sqlite3）
default = "mysql"

# MySQL 数据库配置
[[database.mysql]]
enable = true                    # 是否启用
id = "default"                   # 数据库唯一标识
host = "127.0.0.1"               # 主机地址
port = 3306                      # 端口号
username = "root"                # 用户名
password = "password"            # 密码
database = "testdb"              # 数据库名
prefix = ""                      # 表前缀（可选）
prefix_sql_holder = "__PREFIX__" # 表前缀占位符（可选）
charset = "utf8mb4"              # 字符集
collate = "utf8mb4_general_ci"   # 排序规则
maxidle = 30                     # 最大空闲连接数
maxconns = 200                   # 最大连接数
timeout = 3000                   # 连接超时（毫秒）
readtimeout = 5000               # 读超时（毫秒）
writetimeout = 5000              # 写超时（毫秒）
maxlifetimeseconds = 1800        # 连接最大生命周期（秒）

# 第二个 MySQL 数据库（用于读写分离等场景）
[[database.mysql]]
enable = true
id = "read_replica"
host = "192.168.1.100"
port = 3306
username = "readonly"
password = "password"
database = "testdb"
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
```

### 配置参数说明

| 参数 | 类型 | 说明 | 默认值 |
|------|------|------|--------|
| `enable` | bool | 是否启用此数据库配置 | false |
| `id` | string | 数据库唯一标识，用于获取指定数据库 | "default" |
| `host` | string | MySQL 服务器地址 | "127.0.0.1" |
| `port` | int | MySQL 服务器端口 | 3306 |
| `username` | string | 数据库用户名 | "root" |
| `password` | string | 数据库密码 | "" |
| `database` | string | 数据库名称 | "test" |
| `prefix` | string | 表名前缀，如 "app_" | "" |
| `prefix_sql_holder` | string | SQL 中的前缀占位符 | "" |
| `charset` | string | 字符集 | "utf8" |
| `collate` | string | 排序规则 | "utf8_general_ci" |
| `maxidle` | int | 最大空闲连接数 | 50 |
| `maxconns` | int | 最大连接数 | 500 |
| `timeout` | int | 连接超时时间（毫秒） | 3000 |
| `readtimeout` | int | 读超时时间（毫秒） | 5000 |
| `writetimeout` | int | 写超时时间（毫秒） | 5000 |
| `maxlifetimeseconds` | int | 连接最大生命周期（秒） | 1800 |

### 代码配置方式

```go
config := gdb.NewMySQLDBConfig()
// 或者使用便捷方法
config := gdb.NewMySQLDBConfigWith("127.0.0.1", 3306, "testdb", "root", "password")

// 设置所有参数
config.Host = "127.0.0.1"
config.Port = 3306
config.Username = "root"
config.Password = "password"
config.Database = "testdb"
config.Charset = "utf8mb4"
config.Collate = "utf8mb4_general_ci"
config.MaxOpenConns = 200
config.MaxIdleConns = 30
config.Timeout = 3000
config.ReadTimeout = 5000
config.WriteTimeout = 5000
config.TablePrefix = "app_"
config.TablePrefixSQLIdentifier = "__PREFIX__"
```

---

## SQLite3 配置详解

### 配置文件格式（TOML）

```toml
[database]
default = "sqlite3"

# SQLite3 数据库配置
[[database.sqlite3]]
enable = true                    # 是否启用
id = "default"                   # 数据库唯一标识
database = "./data/app.db"       # 数据库文件路径
password = ""                    # 密码（为空则不加密，不为空则加密）
prefix = ""                      # 表前缀（可选）
prefix_sql_holder = "__PREFIX__" # 表前缀占位符（可选）
syncmode = 1                     # 同步模式: 0=OFF, 1=NORMAL, 2=FULL, 3=EXTRA
openmode = "rwc"                 # 打开模式: ro, rw, rwc, memory
cachemode = "shared"             # 缓存模式: shared, private
```

### 配置参数说明

| 参数 | 类型 | 说明 | 可选值 | 默认值 |
|------|------|------|--------|--------|
| `enable` | bool | 是否启用此数据库配置 | true/false | false |
| `id` | string | 数据库唯一标识 | - | "default" |
| `database` | string | 数据库文件路径 | - | "test" |
| `password` | string | 数据库密码（加密） | - | "" |
| `prefix` | string | 表名前缀 | - | "" |
| `prefix_sql_holder` | string | SQL 中的前缀占位符 | - | "" |
| `syncmode` | int | 同步模式 | 0, 1, 2, 3 | 0 |
| `openmode` | string | 打开模式 | ro, rw, rwc, memory | "rw" |
| `cachemode` | string | 缓存模式 | shared, private | "shared" |

#### syncmode 详解

- `0 (OFF)`: 最快，但数据丢失风险最高（不推荐生产环境）
- `1 (NORMAL)`: 平衡性能和安全（推荐）
- `2 (FULL)`: 最安全，性能较低
- `3 (EXTRA)`: 超级安全，性能最低

#### openmode 详解

- `ro`: 只读模式，数据库必须存在
- `rw`: 读写模式，数据库必须存在
- `rwc`: 读写模式，不存在则创建（推荐）
- `memory`: 内存数据库，程序退出后数据丢失

#### cachemode 详解

- `shared`: 共享缓存模式（推荐）
- `private`: 私有缓存模式

### 代码配置方式

```go
config := gdb.NewSQLite3DBConfig()
// 或者使用便捷方法
config := gdb.NewSQLite3DBConfigWith("./data/app.db", "rwc", "shared", 1)

// 设置所有参数
config.Database = "./data/app.db"
config.OpenMode = gdb.OpenModeReadWriteCreate  // "rwc"
config.CacheMode = gdb.CacheModeShared         // "shared"
config.SyncMode = gdb.SyncModeNormal           // 1
config.TablePrefix = "app_"
config.TablePrefixSQLIdentifier = "__PREFIX__"
```

### SQLite3 加密配置

```toml
[[database.sqlite3]]
enable = true
id = "secure"
database = "./data/secure.db"
password = "my-secret-password"  # 设置密码启用加密
prefix = ""
prefix_sql_holder = "__PREFIX__"
syncmode = 1
openmode = "rwc"
cachemode = "shared"
```

**注意**：SQLite3 加密功能需要数据库文件支持加密格式。

---

## 完整示例代码

### 示例 1：MySQL 完整 CRUD 操作

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    // 初始化
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    gmc.DB.Init(cfg)
    db := gmc.DB.DB()

    // 1. 插入
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
    userId := result.LastInsertID()
    fmt.Printf("插入成功，ID: %d\n", userId)

    // 2. 查询
    ar = db.AR()
    ar.From("users").Where(gdb.M{"id": userId})
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }
    user := rs.Row()
    fmt.Printf("查询结果: %+v\n", user)

    // 3. 更新
    ar = db.AR()
    ar.Update("users", 
        gdb.M{"age": 31, "email": "john.updated@example.com"},
        gdb.M{"id": userId},
    )
    result, err = db.Exec(ar)
    if err != nil {
        panic(err)
    }
    fmt.Printf("更新成功，影响行数: %d\n", result.RowsAffected())

    // 4. 删除
    ar = db.AR()
    ar.Delete("users", gdb.M{"id": userId})
    result, err = db.Exec(ar)
    if err != nil {
        panic(err)
    }
    fmt.Printf("删除成功，影响行数: %d\n", result.RowsAffected())
}
```

### 示例 2：事务处理

```go
package main

import (
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    gmc.DB.Init(cfg)
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

    // 在事务中插入用户
    ar := db.AR()
    ar.Insert("users", gdb.M{
        "name":  "Alice",
        "email": "alice@example.com",
    })
    _, err = db.ExecTx(ar, tx)
    if err != nil {
        tx.Rollback()
        panic(err)
    }

    // 在事务中插入订单
    ar = db.AR()
    ar.Insert("orders", gdb.M{
        "user_id": 1,
        "amount":  100.50,
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

### 示例 3：复杂查询

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    gmc.DB.Init(cfg)
    db := gmc.DB.DB()

    // 复杂查询：JOIN + WHERE + GROUP BY + HAVING + ORDER BY + LIMIT
    ar := db.AR()
    ar.From("users").
        Join("orders", "o", "users.id = o.user_id", "LEFT").
        Select("users.id, users.name, COUNT(o.id) as order_count, SUM(o.amount) as total_amount").
        Where(gdb.M{"users.status": "active"}).
        Where(gdb.M{"users.age >=": 18}).
        GroupBy("users.id").
        Having("COUNT(o.id) > 0").
        OrderBy("total_amount", "DESC").
        Limit(0, 10)

    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }

    rows := rs.Rows()
    for _, row := range rows {
        fmt.Printf("用户: %s, 订单数: %s, 总金额: %s\n",
            row["name"], row["order_count"], row["total_amount"])
    }
}
```

### 示例 4：SQLite3 使用

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    gmc.DB.Init(cfg)

    // 获取 SQLite3 数据库
    db := gmc.DB.SQLite3("default")
    if db == nil {
        panic("SQLite3 数据库未配置")
    }

    // 使用方式与 MySQL 完全相同
    ar := db.AR()
    ar.From("users").Select("*")
    rs, err := db.Query(ar)
    if err != nil {
        panic(err)
    }

    rows := rs.Rows()
    for _, row := range rows {
        fmt.Printf("User: %s\n", row["name"])
    }
}
```

### 示例 5：多数据库使用

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc"
    gdb "github.com/snail007/gmc/module/db"
)

func main() {
    cfg := gmc.New.Config()
    cfg.SetConfigFile("app.toml")
    cfg.ReadInConfig()
    gmc.DB.Init(cfg)

    // 主数据库（写操作）
    writeDB := gmc.DB.MySQL("default")

    // 从数据库（读操作）
    readDB := gmc.DB.MySQL("read_replica")

    // 写入主库
    ar := writeDB.AR()
    ar.Insert("users", gdb.M{
        "name":  "Bob",
        "email": "bob@example.com",
    })
    result, err := writeDB.Exec(ar)
    if err != nil {
        panic(err)
    }
    userId := result.LastInsertID()

    // 从读库查询
    ar = readDB.AR()
    ar.From("users").Where(gdb.M{"id": userId})
    rs, err := readDB.Query(ar)
    if err != nil {
        panic(err)
    }
    user := rs.Row()
    fmt.Printf("从读库查询: %+v\n", user)
}
```

---

## 总结

### 关键点

1. **AR 是链式调用**：所有方法都返回 `ActiveRecord` 接口，支持链式调用
2. **配置灵活**：支持代码配置和配置文件两种方式
3. **多数据库支持**：可以同时使用多个 MySQL 或 SQLite3 数据库
4. **结果集处理**：`ResultSet` 提供了丰富的结果处理方法
5. **事务支持**：完整的事务 API
6. **查询缓存**：可选的查询结果缓存功能

### 最佳实践

1. **使用配置文件**：生产环境推荐使用配置文件管理数据库连接
2. **连接池优化**：根据并发量合理设置连接池大小
3. **使用事务**：涉及多表操作时使用事务保证一致性
4. **参数化查询**：AR 自动使用参数化查询，防止 SQL 注入
5. **错误处理**：始终检查错误返回值
6. **结果集关闭**：虽然 ResultSet 会自动关闭，但建议及时处理结果

### 注意事项

1. **表前缀**：如果配置了表前缀，使用 `__PREFIX__` 占位符或直接使用表名
2. **SQLite3 加密**：需要数据库文件支持加密格式
3. **连接数限制**：注意数据库服务器的最大连接数限制
4. **字符集**：MySQL 推荐使用 `utf8mb4`
5. **SQLite3 同步模式**：生产环境推荐使用 `NORMAL` (1) 模式

---
