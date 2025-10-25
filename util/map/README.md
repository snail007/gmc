# gmap 包

## 简介

gmap 包提供了线程安全的 Map 容器，支持保持键的插入顺序，并提供丰富的操作方法。

## 功能特性

- **线程安全**：所有操作都是线程安全的
- **有序遍历**：保持键的插入顺序
- **丰富操作**：存储、删除、查找、遍历等
- **内存管理**：支持 GC 回收内存

## 安装

```bash
go get github.com/snail007/gmc/util/map
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/map"
)

func main() {
    m := gmap.New()
    
    // 存储键值对
    m.Store("name", "John")
    m.Store("age", 30)
    m.Store("city", "New York")
    
    // 获取值
    if value, ok := m.Load("name"); ok {
        fmt.Println("Name:", value)
    }
    
    // 检查是否存在
    if m.Has("age") {
        fmt.Println("Age exists")
    }
    
    // 遍历（按插入顺序）
    m.Range(func(key, value interface{}) bool {
        fmt.Printf("%v: %v\n", key, value)
        return true
    })
    
    // 删除
    m.Delete("city")
    
    // 获取所有键
    keys := m.Keys()
    fmt.Println("Keys:", keys)
    
    // 获取长度
    fmt.Println("Length:", m.Len())
    
    // 转换为普通 map
    normalMap := m.ToMap()
    fmt.Println(normalMap)
}
```

## API 参考

### 创建和基本操作

- `New() *Map`：创建新 Map
- `Store(key, value)`：存储键值对
- `Load(key) (value, bool)`：获取值
- `Delete(key)`：删除键
- `Has(key) bool`：检查键是否存在
- `Len() int`：获取长度

### 遍历

- `Range(func(key, value) bool)`：遍历所有键值对
- `RangeFast(func(key, value) bool)`：快速遍历（无序）
- `Keys() []interface{}`：获取所有键
- `Values() []interface{}`：获取所有值

### 栈和队列操作

- `Shift() (key, value, bool)`：移除并返回第一个元素
- `Pop() (key, value, bool)`：移除并返回最后一个元素

### 其他操作

- `Clear()`：清空
- `Clone() *Map`：克隆
- `ToMap() map[interface{}]interface{}`：转换为普通 map
- `GC()`：垃圾回收，释放内存

## 类型别名

```go
// M 是 map[string]interface{} 的别名
type M = map[string]interface{}

// Mii 是 map[interface{}]interface{} 的别名
type Mii = map[interface{}]interface{}

// Mss 是 map[string]string 的别名
type Mss = map[string]string
```

## 使用场景

1. **配置管理**：存储应用配置
2. **缓存**：内存缓存
3. **会话存储**：HTTP 会话数据
4. **有序字典**：需要保持插入顺序的场景
5. **并发访问**：多 goroutine 安全访问

## 注意事项

1. 遍历时按插入顺序
2. `RangeFast` 不保证顺序但性能更好
3. `GC()` 可以释放已删除键占用的内存
4. 键可以是任意类型

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
