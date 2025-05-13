
---

## ✅ 项目简要描述（一句话版）

**一个使用 Go 编写的高效内存管理模块，支持内存复用，减少频繁分配以提升性能。**

---

## 📄 README.md 文档

````markdown
# Go Memory Pool Manager

> 🚀 一个使用 Go 编写的高性能内存管理库，支持内存复用，旨在减少频繁的内存分配，提高程序性能。

## ✨ 项目特点

- 🧠 **内存池管理**：统一管理内存块，避免频繁创建/销毁。
- ♻️ **支持复用**：已释放的内存对象可快速重用，减少 GC 压力。
- ⚡ **高性能**：通过预分配与对象池技术，加快内存使用效率。
- 🛠️ **可定制**：支持按需调整池大小、内存块大小等参数。

## 📦 安装

```bash
go get [Jonlson/derectBuffer](https://github.com/Jonlson/derectBuffer)
````

## 🧪 示例代码

```go
package main

import (
    "fmt"
    "github.com/Jonlson/derectBuffer"
)

func main() {
    pool := memmanager.NewMemoryPool(1024, 100) // 每块 1KB，总共预分配 100 块

    buf := pool.Get()
    copy(buf, []byte("Hello, Memory Pool!"))
    fmt.Println(string(buf[:20]))

    pool.Put(buf) // 归还到内存池
}
```

## ⚙️ API 说明

```go
// 创建新的内存池
NewMemoryPool(blockSize int, poolSize int) *MemoryPool

// 从内存池获取内存块
(pool *MemoryPool) Get() []byte

// 归还内存块到池中
(pool *MemoryPool) Put([]byte)
```

## 🔍 使用场景

* 高频请求处理（如网络服务器、消息队列）
* 图像/视频处理等对内存频繁操作的应用
* 替代 Go 原生 sync.Pool 时对内存结构更可控的场景

## 📈 性能对比

| 操作    | 原生内存分配 | 使用 MemoryPool |
| ----- | ------ | ------------- |
| 分配耗时  | 高      | 低             |
| GC 压力 | 高      | 低             |
| 内存复用  | 无      | 有             |

## 📄 License

MIT License

---

欢迎 Star ⭐ 和 Issue 💬 提建议！


