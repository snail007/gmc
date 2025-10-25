# collection 包

## 简介

collection 包提供了集合运算相关的工具函数，目前主要实现了笛卡尔积（Cartesian Product）运算。

## 功能特性

- **笛卡尔积**：计算多个集合的笛卡尔积

## 安装

```bash
go get github.com/snail007/gmc/util/collection
```

## 快速开始

### 笛卡尔积

笛卡尔积是集合论中的基本运算，用于计算多个集合的所有可能组合。

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/collection"
)

func main() {
    // 定义集合
    set1 := []interface{}{1, 2}
    set2 := []interface{}{7, 8}
    
    // 计算笛卡尔积
    result := collection.CartesianProduct(set1, set2)
    
    // 输出: [[1 7] [1 8] [2 7] [2 8]]
    fmt.Println(result)
}
```

### 多个集合的笛卡尔积

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/collection"
)

func main() {
    colors := []interface{}{"red", "blue"}
    sizes := []interface{}{"S", "M", "L"}
    materials := []interface{}{"cotton", "polyester"}
    
    // 计算三个集合的笛卡尔积
    result := collection.CartesianProduct(colors, sizes, materials)
    
    // 输出所有可能的组合
    for _, combo := range result {
        fmt.Printf("Color: %v, Size: %v, Material: %v\n", 
            combo[0], combo[1], combo[2])
    }
    
    // 输出:
    // Color: red, Size: S, Material: cotton
    // Color: red, Size: S, Material: polyester
    // Color: red, Size: M, Material: cotton
    // Color: red, Size: M, Material: polyester
    // ...
}
```

## API 参考

### CartesianProduct

```go
func CartesianProduct(sets ...[]interface{}) [][]interface{}
```

计算多个集合的笛卡尔积。

**数学定义：**
```
CartesianProduct(A, B) = {(x, y) | ∀ x ∈ A, ∀ y ∈ B}
```

**参数：**
- `sets`：一个或多个集合（interface{} 切片）

**返回值：**
- `[][]interface{}`：笛卡尔积结果，每个元素是一个组合

**示例：**

```go
// 两个集合
A := []interface{}{1, 2}
B := []interface{}{7, 8}
result := collection.CartesianProduct(A, B)
// result: [[1 7] [1 8] [2 7] [2 8]]

// 三个集合
C := []interface{}{"x", "y"}
result := collection.CartesianProduct(A, B, C)
// result: [[1 7 x] [1 7 y] [1 8 x] [1 8 y] [2 7 x] [2 7 y] [2 8 x] [2 8 y]]

// 空集合
result := collection.CartesianProduct()
// result: nil
```

## 使用场景

### 1. 电商商品 SKU 生成

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/collection"
)

func main() {
    colors := []interface{}{"红色", "蓝色", "黑色"}
    sizes := []interface{}{"S", "M", "L", "XL"}
    
    skus := collection.CartesianProduct(colors, sizes)
    
    for i, sku := range skus {
        fmt.Printf("SKU-%03d: 颜色=%v, 尺码=%v\n", 
            i+1, sku[0], sku[1])
    }
}
```

### 2. 参数组合测试

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/collection"
)

func main() {
    browsers := []interface{}{"Chrome", "Firefox", "Safari"}
    os := []interface{}{"Windows", "macOS", "Linux"}
    versions := []interface{}{"v1.0", "v2.0"}
    
    testCombos := collection.CartesianProduct(browsers, os, versions)
    
    fmt.Printf("Total test combinations: %d\n", len(testCombos))
    
    for _, combo := range testCombos {
        fmt.Printf("Test: Browser=%v, OS=%v, Version=%v\n",
            combo[0], combo[1], combo[2])
    }
}
```

### 3. 配置选项生成

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/collection"
)

func main() {
    protocols := []interface{}{"http", "https"}
    methods := []interface{}{"GET", "POST", "PUT"}
    formats := []interface{}{"json", "xml"}
    
    configs := collection.CartesianProduct(protocols, methods, formats)
    
    for _, config := range configs {
        fmt.Printf("Config: %v://.../endpoint [%v] (format: %v)\n",
            config[0], config[1], config[2])
    }
}
```

### 4. 游戏道具组合

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/collection"
)

func main() {
    weapons := []interface{}{"剑", "斧", "弓"}
    armors := []interface{}{"轻甲", "重甲"}
    skills := []interface{}{"火球", "冰冻"}
    
    combinations := collection.CartesianProduct(weapons, armors, skills)
    
    fmt.Printf("Total character builds: %d\n", len(combinations))
    for _, combo := range combinations {
        fmt.Printf("Build: 武器=%v, 护甲=%v, 技能=%v\n",
            combo[0], combo[1], combo[2])
    }
}
```

## 算法说明

笛卡尔积的实现使用了多维索引遍历算法：

1. 为每个集合维护一个索引位置
2. 从最后一个集合开始，逐个递增索引
3. 当某个集合索引到达末尾时，重置为 0 并进位到前一个集合
4. 直到第一个集合索引到达末尾

**时间复杂度：** O(n₁ × n₂ × ... × nₖ)，其中 nᵢ 是第 i 个集合的大小

**空间复杂度：** O(n₁ × n₂ × ... × nₖ)，存储所有组合

## 注意事项

1. **结果数量**：笛卡尔积的结果数量是所有集合大小的乘积，可能会非常大
2. **内存占用**：所有组合都存储在内存中，大规模集合可能导致内存问题
3. **空集合**：如果传入空参数，返回 nil
4. **空元素集合**：如果任何一个集合为空，结果也为空
5. **类型安全**：使用 interface{} 类型，需要在使用时进行类型断言

## 性能考虑

### 结果集大小估算

```go
// 计算结果集大小
setSize := 1
for _, set := range sets {
    setSize *= len(set)
}
fmt.Printf("Result will have %d combinations\n", setSize)
```

### 大数据集处理建议

对于大数据集，考虑：

1. **分批处理**：如果可能，将大集合拆分成小集合分别处理
2. **生成器模式**：如果不需要一次性获取所有组合，考虑实现迭代器
3. **过滤条件**：在生成组合时添加过滤条件，减少结果数量

## 扩展示例

### 带过滤的笛卡尔积

```go
package main

import (
    "fmt"
    "github.com/snail007/gmc/util/collection"
)

func main() {
    numbers := []interface{}{1, 2, 3}
    letters := []interface{}{"A", "B", "C"}
    
    all := collection.CartesianProduct(numbers, letters)
    
    // 过滤：只保留数字为奇数的组合
    var filtered [][]interface{}
    for _, combo := range all {
        if num, ok := combo[0].(int); ok && num%2 == 1 {
            filtered = append(filtered, combo)
        }
    }
    
    fmt.Println("Filtered combinations:", filtered)
    // 输出: [[1 A] [1 B] [1 C] [3 A] [3 B] [3 C]]
}
```

## 相关链接

- [GMC 框架主页](https://github.com/snail007/gmc)
- [笛卡尔积 - 维基百科](https://zh.wikipedia.org/wiki/%E7%AC%9B%E5%8D%A1%E5%84%BF%E7%A7%AF)
