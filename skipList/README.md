# 跳表（Skip List）示例

基于 [huandu/skiplist](https://github.com/huandu/skiplist) 演示 Go 语言中跳表（有序 map）的常见用法：

- 基本增删改查与遍历
- 区间查询（`Find` + `Next`）
- 并发读写的封装（`sync.RWMutex`）

## 运行

```bash
cd skipList
go mod tidy
go run .
```

## 输出示例

```text
== 1. 基本操作 ==
Get(10) -> value=world
GetValue(20) -> hello
遍历全部: {10:world} {20:hello} {40:1000}
Remove(40) -> removed=1000, 剩余 Len=2

== 2. 区间查询 ==
区间 [25, 70] 内的元素:
  {30: score-30}
  {40: score-40}
  {50: score-50}
  {60: score-60}
  {70: score-70}

== 3. 并发读写 ==
并发读写完成，最终元素数 Len = 1000
```

## API 要点

| 方法 | 说明 |
| --- | --- |
| `skiplist.New(skiplist.Int)` | 新建跳表，`skiplist.Int` 指定 key 为 int（另有 `String`、`Float64`、`Int64` 等内置类型） |
| `list.Set(key, value)` | 写入，key 已存在时覆盖 value |
| `list.Get(key)` | 返回元素节点 `*Element`（未找到返回 nil） |
| `list.GetValue(key)` | 直接返回 `(value, ok)`，类似 map 取值 |
| `list.Find(key)` | 返回第一个 key ≥ 入参的元素，是区间查询的下界定位 |
| `list.Remove(key)` | 删除并返回被删元素 |
| `elem.Key()` / `elem.Value` | 取节点的 key（方法）与 value（导出字段） |
| `elem.Next()` / `elem.Prev()` | 相邻节点遍历 |
| `list.Front()` / `list.Back()` / `list.Len()` | 首元素 / 尾元素 / 元素个数 |

> 注意：`skiplist.Int` 的 key 是 `int`，取值时需类型断言，例如 `e.Key().(int)`。

## 关于并发安全

`huandu/skiplist` 本身**不是并发安全**的，多 goroutine 读写需自行加锁。示例中的
`ConcurrentSkipList` 用 `sync.RWMutex` 做了封装（`main.go` 中的 `ConcurrentSkipList`），
其中读操作用读锁、写操作用写锁。可用 `go run -race .` 验证无数据竞争。

## 区间查询的写法

跳表本身不提供 `Range(start, end)` 方法，通用做法是：

```go
for e := list.Find(start); e != nil; e = e.Next() {
    if e.Key().(int) > end {
        break
    }
    // 处理 e.Value
}
```

## 内存与清理策略

`huandu/skiplist` 是一个纯有序 map 数据结构，**没有任何自动清理 / 淘汰 / TTL / 容量上限策略**。
它不关心内存水位，也没有「最多 N 个元素」之类的约束，生命周期和内存治理完全由调用方负责。

全部「清理」能力只有显式调用：

| 方法 | 说明 |
| --- | --- |
| `list.Remove(key)` | 删除指定 key |
| `list.RemoveFront()` / `list.RemoveBack()` | 删除最小 / 最大 key，均为 O(1) / O(log n) |
| `list.Init()` | 一次性清空全部元素 |

### 内存会不会爆？

会。内存是 **O(n) 线性增长、无上限**，`Len()` 只增不减，一直 `Set` 不删最终必然 OOM。
每个节点大致占 `Element 结构体 ≈ 80+ 字节` + `levels 切片 ≈ 24 + level×8 字节`（`level` 由
`randLevel` 随机，期望约 2 层），再加上你存入的 `value`。

### 需要「有界」时自己封装

思路是写入后检查数量，超了就按策略淘汰。例如淘汰最小 key：

```go
func (c *ConcurrentSkipList) SetWithCap(key int, value interface{}, cap int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.list.Set(key, value)
    for c.list.Len() > cap {
        c.list.RemoveFront() // 淘汰最小 key
    }
}
```

几点建议：

- 跳表适合做**基于顺序 / 分数**的有界结构（Top-K、滑动窗口、延迟队列等），`Front` / `Back` 是 O(1)。
- 若只想做「最近最少使用」淘汰，跳表本身不维护访问时间，强行用它做 LRU 不如直接用
  `github.com/hashicorp/golang-lru` 这类现成实现。
- 需要「定时过期」时，要在 entry 里存过期时间戳，再由调用方惰性（读取 / 写入时）或后台 goroutine 清理，库本身不提供。