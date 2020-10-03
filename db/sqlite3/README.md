# GMC SQLITE3 Driver
GMC SQLITE3 Driver is designed chained style, to build your database CURD in easy way.

## OPEN SQLITE3 DATABASE

```go
// all of db config are :
// Database                 string
// TablePrefix              string
// TablePrefixSqlIdentifier string
// SyncMode                 int
// OpenMode                 string
// CacheMode                string
// Password                 string

var dbCfg = gmcsqlite3.NewDBConfig()
dbCfg.Database = "test.db"
db, err := gmcsqlite3.NewDB(dbCfg)
if err != nil {
        fmt.Printf("ERR:%s", err)
        return
}
```

## CREATE & OPEN SQLITE3 DATABASE
OpenMode is OPEN_MODE_READ_WRITE_CREATE, database will be created if it not exists.

```go
var dbCfg = sqlite3.NewDBConfig()
dbCfg.OpenMode = gmcsqlite3.OPEN_MODE_READ_WRITE_CREATE
dbCfg.Database = "test.db"
db, err := gmcsqlite3.NewDB(dbCfg)
if err != nil {
        fmt.Printf("ERR:%s", err)
        return
}
```

## CREATE & OPEN ENCRYPTED SQLITE3 DATABASE
If Password is set, database will be created and encrypted if it not exists.

```go
var dbCfg = sqlite3.NewDBConfig()
dbCfg.OpenMode = gmcsqlite3.OPEN_MODE_READ_WRITE_CREATE
dbCfg.Database = "test.db"
dbCfg.Password = "pass123"
db, err := gmcsqlite3.NewDB(dbCfg)
if err != nil {
        fmt.Printf("ERR:%s", err)
        return
}
```

## Using DB GROUP
you can connect to multiple database source in one db group.

```go
group := gmcsqlite3.NewDBGroup("default")
group.Regist("default", gmcsqlite3.NewDBConfigWith("test1.db", "", gmcsqlite3.OPEN_MODE_READ_WRITE, gmcsqlite3.CACHE_MODE_SHARED,gmcsqlite3.SYNC_MODE_OFF))
group.Regist("blog", gmcsqlite3.NewDBConfigWith("test2.db", "", gmcsqlite3.OPEN_MODE_READ_WRITE, gmcsqlite3.CACHE_MODE_SHARED, gmcsqlite3.SYNC_MODE_OFF))
group.Regist("www", gmcsqlite3.NewDBConfigWith("test3.db", "", gmcsqlite3.OPEN_MODE_READ_WRITE,gmcsqlite3.CACHE_MODE_SHARED, gmcsqlite3.SYNC_MODE_OFF))

//group.DB() equal to group.DB("default")
db := group.DB("www")
if db != nil {
    rs, err := db.Query(db.AR().From("test"))
    if err != nil {
        t.Errorf("ERR:%s", err)
    } else {
        fmt.Println(rs.Rows())
    }
} else {
    fmt.Printf("db group config of name %s not found", "www")
}
```

## SELECT

```go
type User{
    ID int `column:"id"`
    Name string `column:"name"`
}
rs, err := db.Query(db.AR().
            Select("*").
            From("log").
            Where(map[string]interface{}{
                "id": 11,
            })
)
if err != nil {
    fmt.Printf("ERR:%s", err)
} else {
    fmt.Println(rs.Rows())
}
```

### RESULT TO STRUCT
```go
//struct 
_user :=User{}
user,err:=rs.Struct(_user)
if err != nil {
    fmt.Printf("ERR:%s", err)
} else {
    fmt.Println(user)
}
//structs
_user :=User{}
users,err:=rs.Structs(_user)
if err != nil {
    fmt.Printf("ERR:%s", err)
} else {
    fmt.Println(users)
}
//Map structs
_user :=User{}
usersMap,err=rs.MapStructs("id",_user)
if err != nil {
    fmt.Printf("ERR:%s", err)
} else {
    fmt.Println(usersMap)
}
```

## RESULT SET
ResultSet is a CRUD result set,all of ResultSet method and properties are below.
```text
ResultSet.Len()
    how many rows of select
ResultSet.MapRows(keyColumn string) (rowsMap map[string]map[string]string)
    get a map which key is each value of row[keyColumn]
ResultSet.MapStructs(keyColumn string, strucT interface{}) (structsMap map[string]interface{},
 err error)
    get a map which key is row[keyColumn],value is strucT
ResultSet.Rows() (rows []map[string]string)
    get rows of select
ResultSet.Structs(strucT interface{}) (structs []interface{}, err error)
    get array of strucT of select
ResultSet.Row() (row map[string]string)
    get first of rows
ResultSet.Struct(strucT interface{}) (Struct interface{}, err error)
    get first strucT of select 
ResultSet.Values(column string) (values []string)
    get an array contains each row[column] 
ResultSet.MapValues(keyColumn, valueColumn string) (values map[string]string)
    get a map key is each row[column],value is row[valueColumn]
ResultSet.Value(column string) (value string)
    get first row[column] of rows
ResultSet.LastInsertId
    if sql type is insert , this is the last insert id
ResultSet.RowsAffected
    if sql type is write , this is the count of rows affected
```

## ActiveRecord
db.AR() return a new *gmcsqlite3.ActiveRecord,you can use it to build you sql.

## INSERT & INSERT BATCH
insert example:

```go
rs, err := db.Exec(db.AR().Insert("test", map[string]interface{}{
    "id":   "id11122",
    "name": "333",
}))
```

insert batch example:

```go
rs, err := db.Exec(db.AR().InsertBatch("test", []map[string]interface{}{
    map[string]interface{}{
        "id":   "id11122",
        "name": "333",
    },
    map[string]interface{}{
        "id":   "id11122",
        "name": "4444",
    },
}))
```
last insert id and rows affected:

```go
lastInsertId:=rs.LastInsertId
rowsAffected:=rs.RowsAffected
fmt.printf("last insert id : %d,rows affected : %d",lastInsertId,rowsAffected)

```

## UPDATE & UPDATE BATCH
attention the where map.
### basic update

```go
rs, err := db.Exec(db.AR().Update("test", map[string]interface{}{
        "id":   "id11122",
        "name": "333",
    }),map[string]interface{}{
        "pid":   223,
    }))
```

equal to sql :

```sql
UPDATE  `test` 
SET `id` = ? , `name` =  ?
WHERE `pid` = ?
```
### column operate

```go
rs, err := db.Exec(db.AR().Update("test", map[string]interface{}{
    "id":   "id11122",
    "score +": "333",
}),map[string]interface{}{
    "pid":   223,
}))
```

equal to sql:

```sql
UPDATE  `test` 
SET `id` = ? , `score` = `score` + ?
WHERE `pid` = ?
```

### update batch

#### basic

```go
rs, err := db.Exec(db.AR().UpdateBatch("test", []map[string]interface{}{
    map[string]interface{}{
        "id":   "id1",
        "name": "333",
    },
    map[string]interface{}{
        "id":   "id2",
        "name": "4444",
    },
}, []string{"id"}))
rowsAffected:=rs.RowsAffected
fmt.printf("rows affected : %d",rowsAffected)
```

equal to sql :

```sql
UPDATE  `test` 
SET `name` = CASE 
WHEN `id` = ? THEN  ? 
WHEN `id` = ? THEN  ? 
ELSE `score` END 
WHERE id IN (?,?)
```

#### column operate

```go
rs, err := db.Exec(db.AR().UpdateBatch("test", []map[string]interface{}{
    map[string]interface{}{
        "id":   "id11",
        "score +": 10,
    },
    map[string]interface{}{
        "id":   "id22",
        "score +": 20,
    },
}, , []string{"id"}))
rowsAffected:=rs.RowsAffected
fmt.printf("rows affected : %d",rowsAffected)
```

equal to sql :

```sql
UPDATE  `test` 
SET `score` = CASE 
WHEN `id` = ? THEN `score` + ? 
WHEN `id` = ? THEN `score` + ? 
ELSE `score` END 
WHERE id IN (?,?)
```

#### where on more column

```go
rs, err := db.Exec(db.AR().UpdateBatch("test", []map[string]interface{}{
    map[string]interface{}{
        "id":      "id1",
        "gid":     22,
        "name":    "test1",
        "score +": 1,
    }, map[string]interface{}{
        "id":      "id2",
        "gid":     33,
        "name":    "test2",
        "score +": 1,
    },
}, []string{"id", "gid"})
rowsAffected:=rs.RowsAffected
fmt.printf("rows affected : %d",rowsAffected)
```

equal to sql :

```text
UPDATE  `test` 
SET `name` = CASE 
WHEN `id` = ? AND `gid` = ? THEN ? 
WHEN `id` = ? AND `gid` = ? THEN ? 
ELSE `name` END,`score` = CASE 
WHEN `id` = ? AND `gid` = ? THEN `score` + ? 
WHEN `id` = ? AND `gid` = ? THEN `score` + ? 
ELSE `score` END 
WHERE id IN (?,?)   AND gid IN (?,?)
```

## DELETE

```go
rs, err := db.Exec(db.AR().Delete("test", map[string]interface{}{
    "pid":   223,
}))
rowsAffected:=rs.RowsAffected
fmt.printf("rows affected : %d",rowsAffected)
```

### RAW SQL

```go
rs, err := db.Exec(db.AR().Raw("insert into test(id,name) values (?,?)", 555,"6666"))
if err != nil {
    fmt.Printf("ERR:%s", err)
} else {
    fmt.Println(rs.RowsAffected, rs.LastInsertId)
}
```
notice:

```text
if  dbCfg.TablePrefix=="user_" 
    dbCfg.TablePrefixSqlIdentifier="{__PREFIX__}" 
then
    db.AR().Raw("insert into {__PREFIX__}test(id,name) values (?,?)
when execute sql,{__PREFIX__} will be replaced with "user_"
```

## Cache QUERY
gmc sqlite3 driver can caching your query automatically.

you must to set a cache handler to store data.

example:

```go
var (
    cacheData = map[string][]byte{}
)
// MyCache is an example cache handler to set or get cache data . 
type MyCache struct {
}

func (c *MyCache) Set(key string, data []byte, expire uint) (err error) {
    cacheData[key] = data
    log.Println("set cache")
    return
}
func (c *MyCache) Get(key string) (data []byte, err error) {
    if v, ok := cacheData[key]; ok {
        log.Println("form cache")
        return v, nil
    }
    return nil, errors.New("key not found or expired")
}
func main() {
    g := gmcsqlite3.NewDBGroupCache("default", &MyCache{})
    g.Regist("default", gmcsqlite3.NewDBConfigWith("test.db",  "", gmcsqlite3.OPEN_MODE_READ_WRITE, gmcsqlite3.CACHE_MODE_SHARED,gmcsqlite3.SYNC_MODE_OFF))
    // turn on cache on the query by call Cache()    
    fmt.Println(g.DB().Query(g.DB().AR().Cache("testkey", 30).From("test")))
    // get query data from cache directly
    rs, _ := g.DB().Query(g.DB().AR().Cache("testkey", 30).From("test"))
    fmt.Println(rs.Row())
}
```

```text
output like:
2017/11/21 18:12:01 set cache
&{0xc42000d340 0 0} <nil>
2017/11/21 18:12:01 form cache
map[pid:1 id:a1 name:a1111 gid:11]
```

## MORE ABOUT ActiveRecord

More using of ActiveRecord , please refer to the `mysql_test.go`.